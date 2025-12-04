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
	"github.com/tmythicator/ticker-rush/server/internal/config"
	repos "github.com/tmythicator/ticker-rush/server/internal/repository/postgres"
	app_redis "github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/model"
)

func setupTestDB(t *testing.T) string {
	ctx := context.Background()
	// Load config to get Postgres credentials
	_ = config.LoadEnv() // Ignore error, might be already loaded or missing
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Connect to default 'postgres' database to create the test database
	defaultConnStr := fmt.Sprintf("postgres://%s:%s@localhost:%d/postgres?sslmode=disable",
		cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresPort)
	conn, err := pgxpool.New(ctx, defaultConnStr)
	if err != nil {
		t.Fatalf("failed to connect to local postgres: %v. Is it running?", err)
	}
	defer conn.Close()

	// Create a unique database name
	dbName := fmt.Sprintf("test_exchange_%d", time.Now().UnixNano())
	_, err = conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		t.Fatalf("failed to create test database %s: %v", dbName, err)
	}

	// Cleanup: Drop the database after test
	t.Cleanup(func() {
		// Reconnect to default DB to drop the test DB
		cleanupConn, err := pgxpool.New(context.Background(), defaultConnStr)
		if err != nil {
			t.Logf("failed to connect to cleanup db: %v", err)
			return
		}
		defer cleanupConn.Close()

		_, err = cleanupConn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE %s WITH (FORCE)", dbName))
		if err != nil {
			t.Logf("failed to drop test database %s: %v", dbName, err)
		}
	})

	// Connection string for the new test database
	testConnStr := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable",
		cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresPort, dbName)

	// Run Migrations
	// We assume the test is running from backend/cmd/exchange, so migrations are in ../../db/migrations
	cmd := exec.Command("goose", "-dir", "../../db/migrations", "postgres", testConnStr, "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to run migrations: %s", err)
	}

	return testConnStr
}

func setupTestRouter(t *testing.T) (*gin.Engine, *miniredis.Miniredis, *pgxpool.Pool) {
	mr, _ := miniredis.Run()
	valkeyClient = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	connStr := setupTestDB(t)
	dbPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	userRepo = repos.NewUserRepository(dbPool)
	stockRepo := repos.NewStockRepository(dbPool)
	stockRepo.UpsertStock(context.Background(), "AAPL", "Apple Inc.")

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

	reqBody := `{"user_id": 1, "password": "password123"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/newUser", bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	user, _ := userRepo.GetUser(ctx, 1)
	assert.Equal(t, 10000.0, user.Balance)
}

func TestBuyStock(t *testing.T) {
	r, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	// Setup Market Data
	quote := model.Quote{Symbol: "AAPL", Price: 150.0, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:AAPL", quoteBytes, 0)

	// Setup User
	userRepo.CreateUser(ctx, 1, "password123")
	// Update balance for test
	user, _ := userRepo.GetUser(ctx, 1)
	user.Balance = 1000.0
	userRepo.SaveUser(ctx, user)

	// Perform Buy
	reqBody := `{"user_id": 1, "symbol": "AAPL", "count": 2}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/buy", bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify User State
	updatedUser, _ := userRepo.GetUser(ctx, 1)
	assert.Equal(t, 700.0, updatedUser.Balance)
	assert.Equal(t, int32(2), updatedUser.Portfolio["AAPL"])
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
	userRepo.CreateUser(ctx, 1, "password123")
	user, _ := userRepo.GetUser(ctx, 1)
	user.Balance = 0.0
	user.Portfolio = map[string]int32{"AAPL": 5}
	userRepo.SaveUser(ctx, user)

	// Perform Sell
	reqBody := `{"user_id": 1, "symbol": "AAPL", "count": 2}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sell", bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify User State
	updatedUser, _ := userRepo.GetUser(ctx, 1)
	assert.Equal(t, 300.0, updatedUser.Balance)
	assert.Equal(t, int32(3), updatedUser.Portfolio["AAPL"])
}

func TestInsufficientFunds(t *testing.T) {
	r, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	quote := model.Quote{Symbol: "AAPL", Price: 150.0, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:AAPL", quoteBytes, 0)

	userRepo.CreateUser(ctx, 1, "password123")
	user, _ := userRepo.GetUser(ctx, 1)
	user.Balance = 100.0
	userRepo.SaveUser(ctx, user)

	reqBody := `{"user_id": 1, "symbol": "AAPL", "count": 1}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/buy", bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
