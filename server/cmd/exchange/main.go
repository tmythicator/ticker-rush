package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/internal/config"
	"github.com/tmythicator/ticker-rush/server/internal/storage"
	"github.com/tmythicator/ticker-rush/server/model"
)

var (
	ctx          = context.Background()
	userRepo     *storage.UserRepository
	marketRepo   *storage.MarketRepository
	valkeyClient *redis.Client
)

func getQuote(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "AAPL")

	isTracked := slices.Contains(config.Tickers, symbol)

	if !isTracked {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Symbol '%s' is not tracked. Available: %v", symbol, config.Tickers),
		})
		return
	}

	quote, err := marketRepo.GetQuote(ctx, symbol)

	if err == redis.Nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Warming up..."})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Valkey error"})
		return
	}

	c.JSON(http.StatusOK, quote)
}

func buyStock(c *gin.Context) {
	var req model.TradeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 1. Get current price
	quote, err := marketRepo.GetQuote(ctx, req.Symbol)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol not found"})
		return
	}

	cost := quote.Price * float64(req.Count)

	// 2. Get User
	user, err := userRepo.GetUser(ctx, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User error"})
		return
	}

	// 3. Check Balance
	if user.Balance < cost {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		return
	}

	// 4. Execute Trade
	user.Balance -= cost
	user.Portfolio[req.Symbol] += req.Count

	// 5. Save User
	if err := userRepo.SaveUser(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func sellStock(c *gin.Context) {
	var req model.TradeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 1. Get User
	user, err := userRepo.GetUser(ctx, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User error"})
		return
	}

	// 2. Check Portfolio
	if user.Portfolio[req.Symbol] < req.Count {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough shares"})
		return
	}

	// 3. Get Price
	quote, err := marketRepo.GetQuote(ctx, req.Symbol)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Market closed"})
		return
	}

	revenue := quote.Price * float64(req.Count)

	// 4. Execute Trade
	user.Balance += revenue
	user.Portfolio[req.Symbol] -= req.Count
	if user.Portfolio[req.Symbol] == 0 {
		delete(user.Portfolio, req.Symbol)
	}

	// 5. Save User
	if err := userRepo.SaveUser(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func createUser(c *gin.Context) {
	var req struct {
		UserID   int64  `json:"user_id"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := userRepo.CreateUser(ctx, req.UserID, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func main() {
	var err error
	valkeyClient, err = storage.NewRedisClient(config.REDIS_ADDR)
	if err != nil {
		log.Fatalf("❌ API failed to start: Valkey connection error: %v", err)
	} else {
		log.Println("✅ Connected to Valkey")
	}

	userRepo = storage.NewUserRepository(valkeyClient)
	marketRepo = storage.NewMarketRepository(valkeyClient)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://" + config.CLIENT_ADDR},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/api/quote", getQuote)
	r.POST("/api/buy", buyStock)
	r.POST("/api/sell", sellStock)
	r.POST("/api/newUser", createUser)

	log.Printf("✅ Exchange API running on %s\n", config.SERVER_PORT)
	r.Run(config.SERVER_PORT)
}
