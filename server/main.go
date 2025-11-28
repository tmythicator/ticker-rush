package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/lib"
)

const (
	SERVER_PORT = ":8080"
	SERVER_ADDR = "localhost" + SERVER_PORT
	REDIS_ADDR  = "localhost:6379"
	CLIENT_ADDR = "localhost:5173"
)

var TICKERS = []string{
	"AAPL",
	"BINANCE:BTCUSDT",
}

var (
	ctx    = context.Background()
	rdb    *redis.Client
	apiKey string
)

func getQuote(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "AAPL")

	isTracked := slices.Contains(TICKERS, symbol)

	if !isTracked {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Market data warming up..."})
		return
	}

	key := "market:" + symbol
	val, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Market data warming up..."})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cache error"})
		return
	}

	c.Data(http.StatusOK, "application/json", []byte(val))
}

func main() {
	lib.LoadEnv()
	apiKey = os.Getenv("FINNHUB_API_KEY")

	if apiKey == "" {
		log.Fatal("FINNHUB_API_KEY is not set.")
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: REDIS_ADDR,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Valkey: %v", err)
	}
	log.Println("Connected to Valkey")

	log.Printf("Starting workers for %d tickers...", len(TICKERS))
	for _, symbol := range TICKERS {
		lib.FetchMarketData(ctx, symbol, apiKey, rdb)
		time.Sleep(2 * time.Second)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://" + CLIENT_ADDR},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/api/quote", getQuote)

	log.Printf("Exchange Server running on %s\n", SERVER_ADDR)
	r.Run(SERVER_PORT)
}
