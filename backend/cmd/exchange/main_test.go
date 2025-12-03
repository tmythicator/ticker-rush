package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/tmythicator/ticker-rush/server/internal/storage"
	"github.com/tmythicator/ticker-rush/server/model"
)

func setupTestRouter() (*gin.Engine, *miniredis.Miniredis) {
	mr, _ := miniredis.Run()
	valkeyClient = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	userRepo = storage.NewUserRepository(valkeyClient)
	marketRepo = storage.NewMarketRepository(valkeyClient)
	ctx = context.Background()

	r := gin.Default()
	r.POST("/api/buy", buyStock)
	r.POST("/api/sell", sellStock)
	r.POST("/api/newUser", createUser)
	return r, mr
}

func TestCreateUser(t *testing.T) {
	r, mr := setupTestRouter()
	defer mr.Close()

	reqBody := `{"user_id": 1, "password": "password123"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/newUser", bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	user, _ := userRepo.GetUser(ctx, 1)
	assert.Equal(t, 10000.0, user.Balance)
}

func TestBuyStock(t *testing.T) {
	r, mr := setupTestRouter()
	defer mr.Close()

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
	assert.Equal(t, 2, updatedUser.Portfolio["AAPL"])
}

func TestSellStock(t *testing.T) {
	r, mr := setupTestRouter()
	defer mr.Close()

	// Setup Market Data
	quote := model.Quote{Symbol: "AAPL", Price: 150.0, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:AAPL", quoteBytes, 0)

	// Setup User
	userRepo.CreateUser(ctx, 1, "password123")
	user, _ := userRepo.GetUser(ctx, 1)
	user.Balance = 0.0
	user.Portfolio = map[string]int{"AAPL": 5}
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
	assert.Equal(t, 3, updatedUser.Portfolio["AAPL"])
}

func TestInsufficientFunds(t *testing.T) {
	r, mr := setupTestRouter()
	defer mr.Close()

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
