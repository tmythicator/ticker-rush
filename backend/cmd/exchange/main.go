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
	db "github.com/tmythicator/ticker-rush/server/internal/gen/sqlc"
	"github.com/tmythicator/ticker-rush/server/internal/repository/postgres"
	valkey "github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/model"
	pb "github.com/tmythicator/ticker-rush/server/proto/user"
)

var (
	ctx           = context.Background()
	userRepo      *postgres.UserRepository
	portfolioRepo *postgres.PortfolioRepository
	marketRepo    *valkey.MarketRepository
	valkeyClient  *go_redis.Client
	dbPool        *pgxpool.Pool
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

	// START TRANSACTION
	tx, err := dbPool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction start failed"})
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Create Transactional Repos
	txUserRepo := userRepo.WithTx(tx)
	txPortfolioRepo := portfolioRepo.WithTx(tx)

	// 2. Get User
	user, err := txUserRepo.GetUserForUpdate(ctx, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User error"})
		return
	}

	// 3. Check Balance
	if user.Balance < cost {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		return
	}

	// 4. Get Current Portfolio Item
	item, err := txPortfolioRepo.GetPortfolioItemForUpdate(ctx, user.Id, req.Symbol)
	if err != nil {
		item = db.PortfolioItem{
			UserID:       user.Id,
			StockSymbol:  req.Symbol,
			Quantity:     0,
			AveragePrice: 0,
		}
	}

	// 5. Execute Trade Logic
	user.Balance -= cost

	currentTotalValue := float64(item.Quantity) * item.AveragePrice
	newTotalValue := currentTotalValue + (float64(req.Count) * quote.Price)
	newTotalQty := item.Quantity + int32(req.Count)

	item.Quantity = newTotalQty
	item.AveragePrice = newTotalValue / float64(newTotalQty)

	// 6. Persistence
	if err := txUserRepo.SaveUser(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user balance"})
		return
	}

	if err := txPortfolioRepo.SetPortfolioItem(ctx, user.Id, req.Symbol, item.Quantity, item.AveragePrice); err != nil {
		log.Printf("ERROR: Failed to update portfolio for user %d", user.Id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update portfolio"})
		return
	}

	// COMMIT
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
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

	// START TRANSACTION
	tx, err := dbPool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction start failed"})
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Create Transactional Repos
	txUserRepo := userRepo.WithTx(tx)
	txPortfolioRepo := portfolioRepo.WithTx(tx)

	// 1. Get User
	user, err := txUserRepo.GetUserForUpdate(ctx, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User error"})
		return
	}

	// 2. Check Portfolio Item
	item, err := txPortfolioRepo.GetPortfolioItemForUpdate(ctx, user.Id, req.Symbol)
	if err != nil || item.Quantity < int32(req.Count) {
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

	// 4. Execute Trade Logic
	user.Balance += revenue
	item.Quantity -= int32(req.Count)

	// 5. Persistence
	if err := txUserRepo.SaveUser(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user balance"})
		return
	}

	if item.Quantity == 0 {
		if err := txPortfolioRepo.DeletePortfolioItem(ctx, user.Id, req.Symbol); err != nil {
			log.Printf("ERROR: Failed to delete portfolio item for user %d", user.Id)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update portfolio"})
			return
		}
	} else {
		if err := txPortfolioRepo.SetPortfolioItem(ctx, user.Id, req.Symbol, item.Quantity, item.AveragePrice); err != nil {
			log.Printf("ERROR: Failed to update portfolio item for user %d", user.Id)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update portfolio"})
			return
		}
	}

	// COMMIT
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func getUser(c *gin.Context) {
	userIDStr := c.Param("id")
	var userID int64
	_, err := fmt.Sscan(userIDStr, &userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := userRepo.GetUser(ctx, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Fetch portfolio separately
	items, err := portfolioRepo.GetPortfolio(ctx, userID)
	if err != nil {
		log.Printf("Error fetching portfolio for user %d: %v", userID, err)
		// Non-fatal, return user with empty portfolio
	}

	// We need a response struct that combines User and Portfolio
	type UserResponse struct {
		*pb.User
		Portfolio map[string]*pb.PortfolioItem `json:"portfolio"`
	}

	portfolioMap := make(map[string]*pb.PortfolioItem)
	for _, item := range items {
		portfolioMap[item.StockSymbol] = &pb.PortfolioItem{
			StockSymbol:  item.StockSymbol,
			Quantity:     item.Quantity,
			AveragePrice: item.AveragePrice,
		}
	}

	resp := UserResponse{
		User:      user,
		Portfolio: portfolioMap,
	}

	c.JSON(http.StatusOK, resp)
}

func createUser(c *gin.Context) {
	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Create user with auto-generated ID
	user, err := userRepo.CreateUser(ctx, req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
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

	valkeyClient, err = valkey.NewClient(fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort))
	if err != nil {
		log.Fatalf("❌ Exchange API failed to start: Valkey connection error (port %d): %v", cfg.RedisPort, err)
	} else {
		log.Println("✅ Connected to Valkey")
	}

	// Connect to Postgres
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)
	var errPool error
	dbPool, errPool = pgxpool.New(ctx, connStr)
	if errPool != nil {
		log.Fatalf("❌ Failed to connect to Postgres: %v", errPool)
	}
	defer dbPool.Close()
	log.Println("✅ Connected to Postgres")

	userRepo = postgres.NewUserRepository(dbPool)
	portfolioRepo = postgres.NewPortfolioRepository(dbPool)
	marketRepo = valkey.NewMarketRepository(valkeyClient)

	r := gin.Default()
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Printf("⚠️ Failed to set trusted proxies: %v", err)
	}
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
	r.GET("/api/user/:id", func(c *gin.Context) {
		getUser(c)
	})
	r.POST("/api/buy", func(c *gin.Context) {
		buyStock(c)
	})
	r.POST("/api/sell", func(c *gin.Context) {
		sellStock(c)
	})
	r.POST("/api/newUser", func(c *gin.Context) {
		createUser(c)
	})

	log.Printf("✅ Exchange API running on :%d\n", cfg.ServerPort)
	if err := r.Run(fmt.Sprintf(":%d", cfg.ServerPort)); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
