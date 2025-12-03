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
	if err := config.LoadEnv(); err != nil {
		log.Printf("⚠️ Failed to load .env: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("⚠️ Failed to load config: %v", err)
	}

	if err := cfg.ValidateFetcher(); err != nil {
		log.Fatalf("❌ Fetcher failed to start: %v", err)
	}

	rdb, err := storage.NewRedisClient(fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort))
	if err != nil {
		log.Fatalf("❌ Fetcher failed to start: Valkey connection error (port:%d): %v", cfg.RedisPort, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	finnhubClient := finnhub.NewClient(cfg.FinnhubKey, 5*time.Second)
	marketWorker := worker.NewMarketFetcher(finnhubClient, rdb)

	fmt.Printf("✅ Worker service started. Tracking %d tickers...\n", len(cfg.Tickers))

	for _, symbol := range cfg.Tickers {
		marketWorker.Start(ctx, symbol, cfg.FetchInterval)
		time.Sleep(cfg.SleepInterval)
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Fetcher shutting down...")
	cancel()
	time.Sleep(1 * time.Second)
}
