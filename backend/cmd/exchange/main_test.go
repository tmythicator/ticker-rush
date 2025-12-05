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
)

const testEmail = "userTest@example.com"
const testUserID = 333

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

	reqBody := fmt.Sprintf(`{"user_id": %d, "password": "password123", "email": "%s"}`, testUserID, testEmail)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/newUser", bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	user, _ := userRepo.GetUser(ctx, testUserID)
	assert.Equal(t, testEmail, user.Email)
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
	userRepo.CreateUser(ctx, testUserID, "password123", testEmail)
	// Update balance for test
	user, _ := userRepo.GetUser(ctx, testUserID)
	user.Balance = 1000.0
	userRepo.SaveUser(ctx, user)

	// Perform Buy
	reqBody := fmt.Sprintf(`{"user_id": %d, "symbol": "AAPL", "count": 2}`, testUserID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/buy", bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify User State
	updatedUser, _ := userRepo.GetUser(ctx, testUserID)
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
	userRepo.CreateUser(ctx, testUserID, "password123", testEmail)
	user, _ := userRepo.GetUser(ctx, testUserID)
	user.Balance = 0.0
	user.Portfolio = map[string]int32{"AAPL": 5}
	userRepo.SaveUser(ctx, user)

	// Perform Sell
	reqBody := fmt.Sprintf(`{"user_id": %d, "symbol": "AAPL", "count": 2}`, testUserID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sell", bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify User State
	updatedUser, _ := userRepo.GetUser(ctx, testUserID)
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

	userRepo.CreateUser(ctx, testUserID, "password123", testEmail)
	user, _ := userRepo.GetUser(ctx, testUserID)
	user.Balance = 100.0
	userRepo.SaveUser(ctx, user)

	reqBody := fmt.Sprintf(`{"user_id": %d, "symbol": "AAPL", "count": 1}`, testUserID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/buy", bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
