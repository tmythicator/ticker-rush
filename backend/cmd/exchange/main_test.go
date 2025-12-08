package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/tmythicator/ticker-rush/server/db"
	"github.com/tmythicator/ticker-rush/server/internal/api"
	"github.com/tmythicator/ticker-rush/server/internal/api/handler"
	"github.com/tmythicator/ticker-rush/server/internal/config"
	repos "github.com/tmythicator/ticker-rush/server/internal/repository/postgres"
	app_redis "github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/internal/service"
	"github.com/tmythicator/ticker-rush/server/model"
	pb "github.com/tmythicator/ticker-rush/server/proto/user"
)

const testEmail = "userTest@example.com"

var (
	valkeyClient  *redis.Client
	dbPool        *pgxpool.Pool
	userRepo      *repos.UserRepository
	portfolioRepo *repos.PortfolioRepository
	marketRepo    *app_redis.MarketRepository
	userService   *service.UserService
	tradeService  *service.TradeService
	marketService *service.MarketService
	restHandler   *handler.RestHandler
)

func setupTestPostgres(t *testing.T) string {
	dbName := "test_db"
	dbUser := "test_user"
	dbPassword := "test_password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %s", err)
	}

	t.Cleanup(func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate postgres container: %s", err)
		}
	})

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	// Run Migrations (Embedded)
	if err := db.Migrate(connStr); err != nil {
		t.Fatalf("failed to run migrations: %s", err)
	}

	return connStr
}

func setupTestRouter(t *testing.T) (*api.Router, *miniredis.Miniredis, *pgxpool.Pool) {
	mr, _ := miniredis.Run()
	valkeyClient = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	connStr := setupTestPostgres(t)
	var err error
	dbPool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	// Initialize Layers
	userRepo = repos.NewUserRepository(dbPool)
	portfolioRepo = repos.NewPortfolioRepository(dbPool)
	marketRepo = app_redis.NewMarketRepository(valkeyClient)

	userService = service.NewUserService(userRepo, portfolioRepo)
	tradeService = service.NewTradeService(userRepo, portfolioRepo, marketRepo, dbPool)
	// Mock config tickers
	tickers := []string{"AAPL", "GOOG", "BTC", "FAKE"}
	marketService = service.NewMarketService(marketRepo, tickers)

	restHandler = handler.NewRestHandler(userService, tradeService, marketService)

	// Mock Config for Router
	cfg := &config.Config{
		ServerPort: 8080,
		ClientPort: 3000,
	}

	// Initialize Router
	router, err := api.NewRouter(restHandler, cfg)
	if err != nil {
		t.Fatalf("Failed to create router: %v", err)
	}

	return router, mr, dbPool
}

func TestCreateUser(t *testing.T) {
	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	reqBody := fmt.Sprintf(`{"email": "%s", "password": "password123", "first_name": "Test", "last_name": "User"}`, testEmail)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/newUser", bytes.NewBufferString(reqBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseUser pb.User
	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	assert.NoError(t, err)

	user, _ := userRepo.GetUser(ctx, responseUser.Id)
	assert.Equal(t, testEmail, user.Email)
}

func TestBuyStock(t *testing.T) {
	const symbol = "AAPL"
	const balance float64 = 1000.0
	const price float64 = 150.0
	const quantity int32 = 2
	const expectedBalance float64 = balance - price*float64(quantity)

	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	// Setup Market Data
	quote := model.Quote{Symbol: symbol, Price: price, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:"+symbol, quoteBytes, 0)

	// Setup User
	createdUser, _ := userRepo.CreateUser(ctx, testEmail, "password123", "Marcel", "Schulz", balance)
	user, _ := userRepo.GetUser(ctx, createdUser.Id)

	// Perform Buy
	reqBody := fmt.Sprintf(`{"user_id": %d, "symbol": "%s", "count": %d}`, user.Id, symbol, quantity)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/buy", bytes.NewBufferString(reqBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify User State
	updatedUser, _ := userRepo.GetUser(ctx, user.Id)
	assert.Equal(t, testEmail, updatedUser.Email)
	assert.Equal(t, expectedBalance, updatedUser.Balance)

	// Verify Portfolio State from Repo
	item, err := portfolioRepo.GetPortfolioItem(ctx, user.Id, symbol)
	assert.NoError(t, err)
	assert.Equal(t, int32(quantity), item.Quantity)
}

func TestSellStock(t *testing.T) {
	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	const mockPrice float64 = 150.0
	const mockStartBalance float64 = 20.0
	const mockPortfolioQuantity int32 = 5
	const mockSellQuantity int32 = 2
	const expectedPortfolioQuantity int32 = mockPortfolioQuantity - mockSellQuantity
	const expectedBalance float64 = mockPrice*float64(mockSellQuantity) + mockStartBalance

	// Setup Market Data
	quote := model.Quote{Symbol: "AAPL", Price: mockPrice, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:AAPL", quoteBytes, 0)

	// Setup User
	createdUser, _ := userRepo.CreateUser(ctx, testEmail, "password123", "Marcel", "Schulz", mockStartBalance)
	user, _ := userRepo.GetUser(ctx, createdUser.Id)

	// Setup Portfolio via Repo
	err := portfolioRepo.SetPortfolioItem(ctx, user.Id, "AAPL", mockPortfolioQuantity, mockPrice)
	assert.NoError(t, err)

	// Perform Sell
	reqBody := fmt.Sprintf(`{"user_id": %d, "symbol": "AAPL", "count": %d}`, user.Id, mockSellQuantity)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sell", bytes.NewBufferString(reqBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify User State
	updatedUser, _ := userRepo.GetUser(ctx, user.Id)
	assert.Equal(t, expectedBalance, updatedUser.Balance)

	// Verify Portfolio
	item, err := portfolioRepo.GetPortfolioItem(ctx, user.Id, "AAPL")
	assert.NoError(t, err)
	assert.Equal(t, expectedPortfolioQuantity, item.Quantity)
}

func TestInsufficientFunds(t *testing.T) {
	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	const mockPrice = 151.0
	const mockStartBalance = 20.0
	const mockBuyQuantity = 1

	quote := model.Quote{Symbol: "AAPL", Price: mockPrice, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:AAPL", quoteBytes, 0)

	createdUser, _ := userRepo.CreateUser(ctx, testEmail, "password123", "Marcel", "Schulz", mockStartBalance)
	user, _ := userRepo.GetUser(ctx, createdUser.Id)

	// balance < cost
	reqBody := fmt.Sprintf(`{"user_id": %d, "symbol": "AAPL", "count": %d}`, user.Id, mockBuyQuantity)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/buy", bytes.NewBufferString(reqBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusPaymentRequired, w.Code)
	assert.Error(t, model.ErrInsufficientFunds)
}

func TestSellAllStock(t *testing.T) {
	const symbol = "AAPL"
	const mockStartBalance float64 = 0.0
	const mockPrice float64 = 150.0
	const mockQuantity int32 = 5
	const mockSellQuantity int32 = 5
	const mockExpectedBalance float64 = mockStartBalance + float64(mockSellQuantity)*mockPrice

	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	// Setup Market Data
	quote := model.Quote{Symbol: symbol, Price: mockPrice, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:"+symbol, quoteBytes, 0)

	// Setup User (FIX: create user implicitly or explicitly)
	createdUser, err := userRepo.CreateUser(ctx, testEmail, "password123", "Marcel", "Schulz", mockStartBalance)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Setup Portfolio
	err = portfolioRepo.SetPortfolioItem(ctx, createdUser.Id, symbol, mockQuantity, mockPrice)
	assert.NoError(t, err)

	// Perform Sell All
	reqBody := fmt.Sprintf(`{"user_id": %d, "symbol": "%s", "count": %d}`, createdUser.Id, symbol, mockSellQuantity)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sell", bytes.NewBufferString(reqBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify User State
	updatedUser, _ := userRepo.GetUser(ctx, createdUser.Id)
	assert.Equal(t, mockExpectedBalance, updatedUser.Balance)

	// Should be deleted
	_, err = portfolioRepo.GetPortfolioItem(ctx, createdUser.Id, symbol)
	assert.Error(t, err, "Portfolio item should be removed (not found error expected)")
}
