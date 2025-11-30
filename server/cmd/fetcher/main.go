package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/tmythicator/ticker-rush/server/internal/clients/finnhub"
	"github.com/tmythicator/ticker-rush/server/internal/config"
	"github.com/tmythicator/ticker-rush/server/internal/storage"
)

func main() {
	apiKey, err := config.GetAPIKey()
	if err != nil {
		log.Fatalf("❌ Fetcher failed to start: API key error: %v", err)
	}
	rdb, err := storage.NewRedisClient(config.REDIS_ADDR)
	if err != nil {
		log.Fatalf("❌ Fetcher failed to start: DB connection error: %v", err)
	} else {
		log.Println("✅ Connected to Valkey")
	}
	ctx := context.Background()

	fmt.Printf("✅ Worker service started. Tracking %d tickers...\n", len(config.Tickers))

	for _, symbol := range config.Tickers {
		finnhub.UpdateMarketData(ctx, symbol, apiKey, rdb)
		time.Sleep(config.FETCH_INTERVAL)
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Fetcher shutting down...")
}
