// Ticker Rush
// Copyright (C) 2025 Alexandr Timchenko
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package main implements the market data fetcher service.
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

	"github.com/jackc/pgx/v5/pgxpool"
	go_redis "github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/internal/clients/coingecko"
	"github.com/tmythicator/ticker-rush/server/internal/clients/finnhub"
	"github.com/tmythicator/ticker-rush/server/internal/config"
	"github.com/tmythicator/ticker-rush/server/internal/repository/postgres"
	"github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/internal/worker"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err = cfg.ValidateFinnhubKey(); err != nil {
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

	// Connect to Postgres
	postgreConnStr := cfg.DatabaseURL()
	pgPool, err := pgxpool.New(ctx, postgreConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pgPool.Close()

	historyRepo := postgres.NewHistoryRepository(pgPool)

	// Initialize Finnhub Client
	finnhubClient := finnhub.NewClient(cfg.FinnhubKey, cfg.FinnhubTimeout)

	// Initialize CoinGecko Client
	coingeckoClient := coingecko.NewClient(cfg.CoingeckoKey, cfg.CoingeckoTimeout)

	// Initialize Workers
	finnhubWorker := worker.NewMarketFetcher(finnhubClient, marketRepo, historyRepo)
	coingeckoWorker := worker.NewMarketFetcher(coingeckoClient, marketRepo, historyRepo)

	fmt.Printf("Worker service started. Tracking %d tickers...\n", len(cfg.Tickers))

	var wg sync.WaitGroup

	// Separate tickers by provider
	var coingeckoTickers []string
	var finnhubTickers []string

	for _, symbol := range cfg.Tickers {
		if len(symbol) > 3 && symbol[:3] == "CG:" {
			coingeckoTickers = append(coingeckoTickers, symbol)
		} else {
			finnhubTickers = append(finnhubTickers, symbol)
		}
	}

	// Start Workers
	if len(coingeckoTickers) > 0 {
		fmt.Printf("Starting CoinGecko worker for %d tickers (Interval: %s)\n", len(coingeckoTickers), cfg.CoingeckoFetchInterval)
		coingeckoWorker.RunLoop(ctx, coingeckoTickers, cfg.CoingeckoFetchInterval, &wg)
	}

	if len(finnhubTickers) > 0 {
		fmt.Printf("Starting Finnhub worker for %d tickers (Interval: %s)\n", len(finnhubTickers), cfg.FinnhubFetchInterval)
		finnhubWorker.RunLoop(ctx, finnhubTickers, cfg.FinnhubFetchInterval, &wg)
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
