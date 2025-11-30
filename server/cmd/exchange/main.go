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
)

var (
	ctx = context.Background()
	rdb *redis.Client
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

	val, err := rdb.Get(ctx, "market:"+symbol).Result()

	if err == redis.Nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Warming up..."})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Valkey error"})
		return
	}

	c.Data(http.StatusOK, "application/json", []byte(val))
}

func main() {
	var err error
	rdb, err = storage.NewRedisClient(config.REDIS_ADDR)
	if err != nil {
		log.Fatalf("❌ API failed to start: DB connection error: %v", err)
	} else {
		log.Println("✅ Connected to Valkey")
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://" + config.CLIENT_ADDR},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/api/quote", getQuote)

	log.Printf("✅ Exchange API running on %s\n", config.SERVER_PORT)
	r.Run(config.SERVER_PORT)
}
