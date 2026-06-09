// Ticker Rush
// Copyright (C) 2025-2026 Alexandr Timchenko
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
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	go_redis "github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	"github.com/tmythicator/ticker-rush/backend/internal/clients/coingecko"
	"github.com/tmythicator/ticker-rush/backend/internal/clients/finnhub"
	"github.com/tmythicator/ticker-rush/backend/internal/config"
	"github.com/tmythicator/ticker-rush/backend/internal/repository/postgres"
	"github.com/tmythicator/ticker-rush/backend/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/backend/internal/worker"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err = cfg.ValidateFinnhubKey(); err != nil {
		log.Fatalf("Fetcher failed to start: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize Redis Client
	rdb := go_redis.NewClient(&go_redis.Options{
		Addr: cfg.RedisHost + ":" + strconv.Itoa(cfg.RedisPort),
	})
	defer func() {
		log.Println("Closing Redis client...")
		_ = rdb.Close()
	}()

	// Initialize Market Repository
	marketRepo := redis.NewMarketRepository(rdb)

	// Connect to Postgres
	postgreConnStr := cfg.DatabaseURL()
	pgPool, err := pgxpool.New(ctx, postgreConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		log.Println("Closing Postgres pool...")
		pgPool.Close()
	}()

	historyRepo := postgres.NewHistoryRepository(pgPool)
	ladderRepo := postgres.NewLadderRepository(pgPool)

	// Initialize Finnhub Client
	finnhubClient := finnhub.NewClient(cfg.FinnhubKey, cfg.FinnhubTimeout)

	// Initialize CoinGecko Client
	coingeckoClient := coingecko.NewClient(cfg.CoingeckoKey, cfg.CoingeckoTimeout)

	// Initialize Workers
	finnhubWorker := worker.NewMarketFetcher("Finnhub", finnhubClient, marketRepo, historyRepo, ladderRepo, &worker.FetcherConfig{
		FetchInterval:   cfg.FinnhubFetchInterval,
		RefreshInterval: cfg.MarketFetcherRefreshInterval,
		RequestTimeout:  cfg.FinnhubTimeout,
	})
	coingeckoWorker := worker.NewMarketFetcher("CoinGecko", coingeckoClient, marketRepo, historyRepo, ladderRepo, &worker.FetcherConfig{
		FetchInterval:   cfg.CoingeckoFetchInterval,
		RefreshInterval: cfg.MarketFetcherRefreshInterval,
		RequestTimeout:  cfg.CoingeckoTimeout,
	})

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := finnhubWorker.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			return fmt.Errorf("finnhub worker error: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		if err := coingeckoWorker.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			return fmt.Errorf("coingecko worker error: %w", err)
		}

		return nil
	})

	log.Println("Fetcher service running...")

	if err := g.Wait(); err != nil {
		log.Printf("Fetcher service stopped with error: %v\n", err)
	} else {
		log.Println("Fetcher service stopped cleanly")
	}
}
