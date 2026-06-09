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
	"golang.org/x/sync/errgroup"
	googlegrpc "google.golang.org/grpc"

	"github.com/tmythicator/ticker-rush/backend/db"
	"github.com/tmythicator/ticker-rush/backend/internal/api"
	grpcapi "github.com/tmythicator/ticker-rush/backend/internal/api/grpc"
	"github.com/tmythicator/ticker-rush/backend/internal/api/handler"
	"github.com/tmythicator/ticker-rush/backend/internal/api/middleware"
	"github.com/tmythicator/ticker-rush/backend/internal/config"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/repository/postgres"
	valkey "github.com/tmythicator/ticker-rush/backend/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
	"github.com/tmythicator/ticker-rush/backend/internal/worker"
)

type App struct {
	cfg                *config.Config
	userService        *service.UserService
	tradeService       *service.TradeService
	marketService      *service.MarketService
	leaderboardService *service.LeaderBoardService
	lifecycleWorker    *worker.LadderLifecycleWorker
	leaderboardWorker  *worker.LeaderboardWorker
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

func NewApp(ctx context.Context, cfg *config.Config) (app *App, err error) {
	// Connect to Valkey
	valkeyClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
	})
	if err = valkeyClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("valkey connection failed: %w", err)
	}
	defer func() {
		if err != nil {
			_ = valkeyClient.Close()
		}
	}()

	// Connect to Postgres
	postgreConnStr := cfg.DatabaseURL()

	if err = db.Migrate(postgreConnStr, cfg.AdminUsername, cfg.AdminPasswordHash); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	postgreClient, err := pgxpool.New(ctx, postgreConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}
	defer func() {
		if err != nil {
			postgreClient.Close()
		}
	}()

	// Initialize repositories
	ladderRepo := postgres.NewLadderRepository(postgreClient)
	userRepo := postgres.NewUserRepository(postgreClient)
	portfolioRepo := postgres.NewPortfolioRepository(postgreClient)
	marketRepo := valkey.NewMarketRepository(valkeyClient)
	leaderboardRepo := valkey.NewLeaderboardRepository(valkeyClient)
	historyRepo := postgres.NewHistoryRepository(postgreClient)
	transactor := postgres.NewPgxTransactor(postgreClient)

	// Initialize services
	userService := service.NewUserService(userRepo, portfolioRepo, ladderRepo)
	tradeService := service.NewTradeService(userRepo, portfolioRepo, marketRepo, ladderRepo, transactor)
	marketService := service.NewMarketService(marketRepo, historyRepo, ladderRepo)
	ladderService := service.NewLadderService(ladderRepo)
	leaderboardService := service.NewLeaderBoardService(userRepo, portfolioRepo, marketRepo, ladderRepo, leaderboardRepo)

	restHandler := handler.NewRestHandler(userService, tradeService, marketService, leaderboardService, ladderService, cfg.JWTSecret)

	// Initialize workers
	leaderboardWorker := worker.NewLeaderboardWorker(leaderboardService, 1*time.Minute)
	lifecycleWorker := worker.NewLadderLifecycleWorker(ladderRepo, portfolioRepo, marketRepo, 1*time.Minute)

	return &App{
		cfg:                cfg,
		userService:        userService,
		tradeService:       tradeService,
		marketService:      marketService,
		leaderboardService: leaderboardService,
		lifecycleWorker:    lifecycleWorker,
		leaderboardWorker:  leaderboardWorker,
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
	g.Go(func() error {
		if lbErr := a.leaderboardWorker.Start(ctx); lbErr != nil && !errors.Is(lbErr, context.Canceled) {
			return fmt.Errorf("leaderboard worker error: %w", lbErr)
		}

		return nil
	})

	// Ladder Lifecycle Worker
	g.Go(func() error {
		if llErr := a.lifecycleWorker.Start(ctx); llErr != nil && !errors.Is(llErr, context.Canceled) {
			return fmt.Errorf("ladder lifecycle worker error: %w", llErr)
		}

		return nil
	})

	return g.Wait()
}

func (a *App) Close() {
	_ = a.valkeyClient.Close()
	a.postgreClient.Close()
}
