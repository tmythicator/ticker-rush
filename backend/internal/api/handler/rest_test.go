package handler_test

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
	"github.com/tmythicator/ticker-rush/server/internal/apperrors"
	"github.com/tmythicator/ticker-rush/server/internal/config"
	"github.com/tmythicator/ticker-rush/server/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/server/internal/proto/user/v1"
	postgreRepo "github.com/tmythicator/ticker-rush/server/internal/repository/postgres"
	redisRepo "github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/internal/service"
	"golang.org/x/crypto/bcrypt"
)

const testEmail = "userTest@example.com"

var (
	ctx           = context.Background()
	valkeyClient  *redis.Client
	dbPool        *pgxpool.Pool
	userRepo      *postgreRepo.UserRepository
	portfolioRepo *postgreRepo.PortfolioRepository
	marketRepo    *redisRepo.MarketRepository
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
		termErr := postgresContainer.Terminate(ctx)
		if termErr != nil {
			t.Fatalf("failed to terminate postgres container: %s", termErr)
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
	userRepo = postgreRepo.NewUserRepository(dbPool)
	portfolioRepo = postgreRepo.NewPortfolioRepository(dbPool)
	marketRepo = redisRepo.NewMarketRepository(valkeyClient)
	transactor := postgreRepo.NewPgxTransactor(dbPool)

	userService = service.NewUserService(userRepo, portfolioRepo)
	tradeService = service.NewTradeService(userRepo, portfolioRepo, marketRepo, transactor)
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

	reqBody := fmt.Sprintf(
		`{"email": "%s", "password": "password123", "first_name": "Test", "last_name": "User"}`,
		testEmail,
	)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBufferString(reqBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseUser user.User

	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	assert.NoError(t, err)

	user, _, err := userRepo.GetUserByEmail(ctx, responseUser.GetEmail())
	assert.NoError(t, err)
	assert.Equal(t, testEmail, user.GetEmail())
	assert.Equal(t, testEmail, user.GetEmail())
}

func TestLogin(t *testing.T) {
	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	// Create User first
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	_, err := userRepo.CreateUser(ctx, testEmail, string(hashedPassword), "Test", "User", 100.0)
	assert.NoError(t, err)

	// Perform Login
	reqBody := fmt.Sprintf(`{"email": "%s", "password": "password123"}`, testEmail)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(reqBody))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify Cookie
	cookies := w.Result().Cookies()
	found := false

	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			found = true

			assert.True(t, cookie.HttpOnly, "Cookie should be HttpOnly")
			assert.Equal(t, "/", cookie.Path, "Cookie path should be /")
			assert.NotEmpty(t, cookie.Value, "Cookie value should not be empty")
		}
	}

	assert.True(t, found, "auth_token cookie should be present")

	// Verify Response Body (Should NOT have token)
	var response map[string]any

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	_, hasToken := response["token"]
	assert.False(t, hasToken, "Response body should NOT contain token")
}

func TestBuyStock(t *testing.T) {
	const (
		symbol                  = "AAPL"
		balance         float64 = 1000.0
		price           float64 = 150.0
		quantity        float64 = 2.0
		expectedBalance float64 = balance - price*quantity
	)

	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	// Setup Market Data
	quote := &exchange.Quote{Symbol: symbol, Price: price, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:"+symbol, quoteBytes, 0)

	// Setup User
	createdUser, err := userRepo.CreateUser(
		ctx,
		testEmail,
		"password123",
		"Marcel",
		"Schulz",
		balance,
	)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user, err := userRepo.GetUser(ctx, createdUser.GetId())
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Generate Token
	token, _ := service.GenerateToken(user)

	// Perform Buy
	reqBody := fmt.Sprintf(`{"symbol": "%s", "count": %f}`, symbol, quantity)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/buy", bytes.NewBufferString(reqBody))
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify User State
	updatedUser, _ := userRepo.GetUser(ctx, user.GetId())
	assert.Equal(t, testEmail, updatedUser.GetEmail())
	assert.Equal(t, expectedBalance, updatedUser.GetBalance())

	// Verify Portfolio State from Repo
	item, err := portfolioRepo.GetPortfolioItem(ctx, user.GetId(), symbol)
	assert.NoError(t, err)
	assert.Equal(t, quantity, item.GetQuantity())
}

func TestSellStock(t *testing.T) {
	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	const (
		mockPrice                 float64 = 150.0
		mockStartBalance          float64 = 20.0
		mockPortfolioQuantity     float64 = 5.0
		mockSellQuantity          float64 = 2.0
		expectedPortfolioQuantity float64 = mockPortfolioQuantity - mockSellQuantity
		expectedBalance           float64 = mockPrice*mockSellQuantity + mockStartBalance
	)

	// Setup Market Data
	quote := &exchange.Quote{Symbol: "AAPL", Price: mockPrice, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:AAPL", quoteBytes, 0)

	// Setup User
	createdUser, _ := userRepo.CreateUser(
		ctx,
		testEmail,
		"password123",
		"Marcel",
		"Schulz",
		mockStartBalance,
	)
	user, _ := userRepo.GetUser(ctx, createdUser.GetId())

	// Setup Portfolio via Repo
	err := portfolioRepo.SetPortfolioItem(
		ctx,
		user.GetId(),
		"AAPL",
		mockPortfolioQuantity,
		mockPrice,
	)
	assert.NoError(t, err)

	// Generate Token
	token, _ := service.GenerateToken(user)

	// Perform Sell
	reqBody := fmt.Sprintf(`{"symbol": "AAPL", "count": %f}`, mockSellQuantity)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/sell", bytes.NewBufferString(reqBody))
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify User State
	updatedUser, _ := userRepo.GetUser(ctx, user.GetId())
	assert.Equal(t, expectedBalance, updatedUser.GetBalance())

	// Verify Portfolio
	item, err := portfolioRepo.GetPortfolioItem(ctx, user.GetId(), "AAPL")
	assert.NoError(t, err)
	assert.Equal(t, expectedPortfolioQuantity, item.GetQuantity())
}

func TestInsufficientFunds(t *testing.T) {
	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	const (
		mockPrice        = 151.0
		mockStartBalance = 20.0
		mockBuyQuantity  = 1
	)

	quote := &exchange.Quote{Symbol: "AAPL", Price: mockPrice, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:AAPL", quoteBytes, 0)

	createdUser, _ := userRepo.CreateUser(
		ctx,
		testEmail,
		"password123",
		"Marcel",
		"Schulz",
		mockStartBalance,
	)
	user, _ := userRepo.GetUser(ctx, createdUser.GetId())

	// Generate Token
	token, _ := service.GenerateToken(user)

	// balance < cost
	reqBody := fmt.Sprintf(`{"symbol": "AAPL", "count": %d}`, mockBuyQuantity)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/buy", bytes.NewBufferString(reqBody))
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusPaymentRequired, w.Code)
	assert.Error(t, apperrors.ErrInsufficientFunds)
}

func TestSellAllStock(t *testing.T) {
	const (
		symbol                      = "AAPL"
		mockStartBalance    float64 = 0.0
		mockPrice           float64 = 150.0
		mockQuantity        float64 = 5.0
		mockSellQuantity    float64 = 5.0
		mockExpectedBalance float64 = mockStartBalance + mockSellQuantity*mockPrice
	)

	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	// Setup Market Data
	quote := &exchange.Quote{Symbol: symbol, Price: mockPrice, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:"+symbol, quoteBytes, 0)

	// Setup User (FIX: create user implicitly or explicitly)
	createdUser, err := userRepo.CreateUser(
		ctx,
		testEmail,
		"password123",
		"Marcel",
		"Schulz",
		mockStartBalance,
	)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Setup Portfolio
	err = portfolioRepo.SetPortfolioItem(ctx, createdUser.GetId(), symbol, mockQuantity, mockPrice)
	assert.NoError(t, err)

	// Generate Token
	token, _ := service.GenerateToken(createdUser)

	// Perform Sell All
	reqBody := fmt.Sprintf(`{"symbol": "%s", "count": %f}`, symbol, mockSellQuantity)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/sell", bytes.NewBufferString(reqBody))
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify User State
	updatedUser, _ := userRepo.GetUser(ctx, createdUser.GetId())
	assert.Equal(t, mockExpectedBalance, updatedUser.GetBalance())

	// Should be deleted
	_, err = portfolioRepo.GetPortfolioItem(ctx, createdUser.GetId(), symbol)
	assert.Error(t, err, "Portfolio item should be removed (not found error expected)")
}
