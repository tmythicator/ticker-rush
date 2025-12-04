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
	"github.com/jackc/pgx/v5/pgxpool"
	go_redis "github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/internal/config"
	"github.com/tmythicator/ticker-rush/server/internal/repository/postgres"
	app_redis "github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/model"
)

var (
	ctx          = context.Background()
	userRepo     *postgres.UserRepository
	marketRepo   *app_redis.MarketRepository
	valkeyClient *go_redis.Client
)

func getQuote(c *gin.Context, cfg *config.Config) {
	symbol := c.DefaultQuery("symbol", "AAPL")

	isTracked := slices.Contains(cfg.Tickers, symbol)

	if !isTracked {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Symbol '%s' is not tracked. Available: %v", symbol, cfg.Tickers),
		})
		return
	}

	quote, err := marketRepo.GetQuote(ctx, symbol)

	if err == go_redis.Nil {
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
	user.Portfolio[req.Symbol] += int32(req.Count)

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
	if user.Portfolio[req.Symbol] < int32(req.Count) {
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
	user.Portfolio[req.Symbol] -= int32(req.Count)
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
	if err := config.LoadEnv(); err != nil {
		log.Printf("⚠️ Failed to load .env: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("⚠️ Failed to load config: %v", err)
	}

	valkeyClient, err = app_redis.NewClient(fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort))
	if err != nil {
		log.Fatalf("❌ Exchange API failed to start: Valkey connection error (port %d): %v", cfg.RedisPort, err)
	} else {
		log.Println("✅ Connected to Valkey")
	}

	// Connect to Postgres
	connStr := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s", cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresPort, cfg.PostgresDB)
	dbPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("❌ Failed to connect to Postgres: %v", err)
	}
	defer dbPool.Close()
	log.Println("✅ Connected to Postgres")

	userRepo = postgres.NewUserRepository(dbPool)
	stockRepo := postgres.NewStockRepository(dbPool)

	// Populate Stocks
	for _, ticker := range cfg.Tickers {
		if err := stockRepo.UpsertStock(ctx, ticker, ticker); err != nil {
			log.Printf("⚠️ Failed to upsert stock %s: %v", ticker, err)
		}
	}

	marketRepo = app_redis.NewMarketRepository(valkeyClient)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{fmt.Sprintf("http://localhost:%d", cfg.ClientPort)},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/api/quote", func(c *gin.Context) {
		getQuote(c, cfg)
	})
	r.POST("/api/buy", buyStock)
	r.POST("/api/sell", sellStock)
	r.POST("/api/newUser", createUser)

	log.Printf("✅ Exchange API running on :%d\n", cfg.ServerPort)
	r.Run(fmt.Sprintf(":%d", cfg.ServerPort))
}
