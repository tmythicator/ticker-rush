package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	go_redis "github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/internal/clients/finnhub"
	"github.com/tmythicator/ticker-rush/server/internal/config"
	"github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/internal/worker"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		log.Printf("Failed to load .env: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := cfg.ValidateFetcher(); err != nil {
		log.Fatalf("Fetcher failed to start: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Redis Client
	rdb := go_redis.NewClient(&go_redis.Options{
		Addr: cfg.RedisHost + ":" + strconv.Itoa(cfg.RedisPort),
	})

	// Initialize Market Repository
	marketRepo := redis.NewMarketRepository(rdb)

	// Initialize Finnhub Client
	finnhubClient := finnhub.NewClient(cfg.FinnhubKey, cfg.FinnhubTimeout)

	// Initialize Worker
	marketWorker := worker.NewMarketFetcher(finnhubClient, marketRepo)

	fmt.Printf("Worker service started. Tracking %d tickers...\n", len(cfg.Tickers))

	var wg sync.WaitGroup

	for _, symbol := range cfg.Tickers {
		marketWorker.Start(ctx, symbol, cfg.FetchInterval, &wg)
		time.Sleep(cfg.SleepInterval)
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Fetcher shutting down...")
	cancel()

	log.Println("Waiting for workers to finish...")
	wg.Wait()
	log.Println("All workers stopped. Exiting.")
}
