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
	"github.com/tmythicator/ticker-rush/server/internal/worker"
)

func main() {
	apiKey, err := config.GetAPIKey()
	if err != nil {
		log.Fatalf("❌ Fetcher failed to start: API key error: %v", err)
	}

	rdb, err := storage.NewRedisClient(config.REDIS_ADDR)
	if err != nil {
		log.Fatalf("❌ Valkey client Error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	finnhubClient := finnhub.NewClient(apiKey, 5*time.Second)
	marketWorker := worker.NewMarketFetcher(finnhubClient, rdb)

	fmt.Printf("✅ Worker service started. Tracking %d tickers...\n", len(config.Tickers))

	for _, symbol := range config.Tickers {
		marketWorker.Start(ctx, symbol)
		time.Sleep(config.FETCH_INTERVAL)
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Fetcher shutting down...")
	cancel()
	time.Sleep(1 * time.Second)
}
