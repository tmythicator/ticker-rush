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

// Package main serves the exchange API.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/db"
	"github.com/tmythicator/ticker-rush/server/internal/api"
	grpcapi "github.com/tmythicator/ticker-rush/server/internal/api/grpc"
	"github.com/tmythicator/ticker-rush/server/internal/api/handler"
	"github.com/tmythicator/ticker-rush/server/internal/api/middleware"
	"github.com/tmythicator/ticker-rush/server/internal/config"
	"github.com/tmythicator/ticker-rush/server/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/server/internal/repository/postgres"
	valkey "github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/internal/service"
	"github.com/tmythicator/ticker-rush/server/internal/worker"
	"golang.org/x/sync/errgroup"
	googlegrpc "google.golang.org/grpc"
)

type App struct {
	cfg                *config.Config
	userService        *service.UserService
	tradeService       *service.TradeService
	marketService      *service.MarketService
	leaderboardService *service.LeaderBoardService
	restHandler        *handler.RestHandler
	valkeyClient       *redis.Client
	postgreClient      *pgxpool.Pool
}

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := NewApp(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer app.Close()

	if err := app.Run(ctx); err != nil {
		log.Fatalf("App exited with error: %v", err)
	}
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	// Connect to Valkey
	valkeyClient, err := valkey.NewClient(fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort))
	if err != nil {
		return nil, fmt.Errorf("valkey creation failed: %w", err)
	}

	// Connect to Postgres
	postgreConnStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.PostgresUser,
		cfg.PostgresPass,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	if err = db.Migrate(postgreConnStr); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	postgreClient, err := pgxpool.New(ctx, postgreConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(postgreClient)
	portfolioRepo := postgres.NewPortfolioRepository(postgreClient)
	marketRepo := valkey.NewMarketRepository(valkeyClient)
	transactor := postgres.NewPgxTransactor(postgreClient)

	// Initialize services
	userService := service.NewUserService(userRepo, portfolioRepo)
	tradeService := service.NewTradeService(userRepo, portfolioRepo, marketRepo, transactor)
	marketService := service.NewMarketService(marketRepo, cfg.Tickers)
	leaderboardService := service.NewLeaderBoardService(userRepo, portfolioRepo, marketRepo, valkeyClient)

	restHandler := handler.NewRestHandler(userService, tradeService, marketService, leaderboardService, cfg.JWTSecret)

	return &App{
		cfg:                cfg,
		userService:        userService,
		tradeService:       tradeService,
		marketService:      marketService,
		leaderboardService: leaderboardService,
		restHandler:        restHandler,
		valkeyClient:       valkeyClient,
		postgreClient:      postgreClient,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// HTTP Server
	router, err := api.NewRouter(a.restHandler, a.cfg)
	if err != nil {
		return fmt.Errorf("failed to create router: %w", err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.ServerPort),
		Handler: router,
	}

	g.Go(func() error {
		log.Printf("Exchange API running on :%d\n", a.cfg.ServerPort)
		if srvErr := srv.ListenAndServe(); srvErr != nil && !errors.Is(srvErr, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server error: %w", srvErr)
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		log.Println("Shutting down HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return srv.Shutdown(shutdownCtx)
	})

	// gRPC Server
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		return fmt.Errorf("failed to listen gRPC: %w", err)
	}

	grpcServer := googlegrpc.NewServer(
		googlegrpc.UnaryInterceptor(middleware.GrpcAuthInterceptor(a.cfg.JWTSecret)),
	)
	exchangeServer := grpcapi.NewExchangeServer(a.tradeService, a.marketService)
	exchange.RegisterExchangeServiceServer(grpcServer, exchangeServer)

	g.Go(func() error {
		log.Printf("Exchange gRPC running on :%d\n", 50051)
		if gErr := grpcServer.Serve(grpcListener); gErr != nil {
			return fmt.Errorf("gRPC server error: %w", gErr)
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()

		return nil
	})

	// Leaderboard Worker
	lbWorker := worker.NewLeaderboardWorker(a.leaderboardService, 10*time.Minute)
	g.Go(func() error {
		if lbErr := lbWorker.Start(ctx); lbErr != nil && !errors.Is(lbErr, context.Canceled) {
			return fmt.Errorf("leaderboard worker error: %w", lbErr)
		}

		return nil
	})

	return g.Wait()
}

func (a *App) Close() {
	if a.valkeyClient != nil {
		_ = a.valkeyClient.Close()
	}
	if a.postgreClient != nil {
		a.postgreClient.Close()
	}
}
