package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	repos "github.com/tmythicator/ticker-rush/server/internal/repository/postgres"
	app_redis "github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/model"
	pb "github.com/tmythicator/ticker-rush/server/proto/user"
)

const testEmail = "userTest@example.com"

func setupTestDB(t *testing.T) string {
	ctx := context.Background()

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

	// Run Migrations
	cmd := exec.Command("goose", "-dir", "../../db/migrations", "postgres", connStr, "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to run migrations: %s", err)
	}

	return connStr
}

func setupTestRouter(t *testing.T) (*gin.Engine, *miniredis.Miniredis, *pgxpool.Pool) {
	mr, _ := miniredis.Run()
	valkeyClient = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	connStr := setupTestDB(t)
	var err error
	dbPool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	userRepo = repos.NewUserRepository(dbPool)
	portfolioRepo = repos.NewPortfolioRepository(dbPool)
	marketRepo = app_redis.NewMarketRepository(valkeyClient)
	ctx = context.Background()

	r := gin.Default()
	r.POST("/api/buy", buyStock)
	r.POST("/api/sell", sellStock)
	r.POST("/api/newUser", createUser)
	return r, mr, dbPool
}

func TestCreateUser(t *testing.T) {
	r, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	reqBody := fmt.Sprintf(`{"email": "%s", "password": "password123", "first_name": "Test", "last_name": "User"}`, testEmail)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/newUser", bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var responseUser pb.User
	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	assert.NoError(t, err)

	user, _ := userRepo.GetUser(ctx, responseUser.Id)
	assert.Equal(t, testEmail, user.Email)
}

func TestBuyStock(t *testing.T) {
	const symbol = "AAPL"
	const balance = 1000.0
	const price = 150.0
	const quantity = 2
	const expectedBalance = balance - price*float64(quantity)

	r, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	// Setup Market Data
	quote := model.Quote{Symbol: symbol, Price: price, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:"+symbol, quoteBytes, 0)

	// Setup User
	createdUser, _ := userRepo.CreateUser(ctx, testEmail, "password123", "Marcel", "Schulz")
	// Update balance for test
	user, _ := userRepo.GetUser(ctx, createdUser.Id)
	user.Balance = balance
	_ = userRepo.SaveUser(ctx, user)

	// Perform Buy
	reqBody := fmt.Sprintf(`{"user_id": %d, "symbol": "%s", "count": %d}`, user.Id, symbol, quantity)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/buy", bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

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
	r, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	// Setup Market Data
	quote := model.Quote{Symbol: "AAPL", Price: 150.0, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:AAPL", quoteBytes, 0)

	// Setup User
	createdUser, _ := userRepo.CreateUser(ctx, testEmail, "password123", "Marcel", "Schulz")
	user, _ := userRepo.GetUser(ctx, createdUser.Id)
	user.Balance = 0.0
	_ = userRepo.SaveUser(ctx, user)

	// Setup Portfolio via Repo
	err := portfolioRepo.SetPortfolioItem(ctx, user.Id, "AAPL", 5, 150.0)
	assert.NoError(t, err)

	// Perform Sell
	reqBody := fmt.Sprintf(`{"user_id": %d, "symbol": "AAPL", "count": 2}`, user.Id)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sell", bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify User State
	updatedUser, _ := userRepo.GetUser(ctx, user.Id)
	assert.Equal(t, 300.0, updatedUser.Balance)

	// Verify Portfolio
	item, err := portfolioRepo.GetPortfolioItem(ctx, user.Id, "AAPL")
	assert.NoError(t, err)
	assert.Equal(t, int32(3), item.Quantity)
}

func TestInsufficientFunds(t *testing.T) {
	r, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	quote := model.Quote{Symbol: "AAPL", Price: 150.0, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:AAPL", quoteBytes, 0)

	createdUser, _ := userRepo.CreateUser(ctx, testEmail, "password123", "Marcel", "Schulz")
	user, _ := userRepo.GetUser(ctx, createdUser.Id)
	user.Balance = 100.0
	_ = userRepo.SaveUser(ctx, user)

	reqBody := fmt.Sprintf(`{"user_id": %d, "symbol": "AAPL", "count": 1}`, user.Id)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/buy", bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSellAllStock(t *testing.T) {
	const symbol = "AAPL"
	const balance = 0.0
	const price = 150.0
	const quantity = 5
	const sellQuantity = 5
	const expectedBalance = balance + float64(sellQuantity*price)

	r, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	// Setup Market Data
	quote := model.Quote{Symbol: symbol, Price: price, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:"+symbol, quoteBytes, 0)

	// Setup User - use existing user 1
	var userID int64 = 1
	user, err := userRepo.GetUser(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	user.Balance = balance
	_ = userRepo.SaveUser(ctx, user)

	// Setup Portfolio
	err = portfolioRepo.SetPortfolioItem(ctx, user.Id, symbol, quantity, price)
	assert.NoError(t, err)

	// Perform Sell All
	reqBody := fmt.Sprintf(`{"user_id": %d, "symbol": "%s", "count": %d}`, user.Id, symbol, sellQuantity)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sell", bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify User State
	updatedUser, _ := userRepo.GetUser(ctx, user.Id)
	assert.Equal(t, expectedBalance, updatedUser.Balance)

	// Should be deleted
	_, err = portfolioRepo.GetPortfolioItem(ctx, user.Id, symbol)
	assert.Error(t, err, "Portfolio item should be removed (not found error expected)")
}
