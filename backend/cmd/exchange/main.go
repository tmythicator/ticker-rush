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
	googlegrpc "google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Failed to load .env: %v", err)
	}

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to Valkey
	valkeyClient, err := valkey.NewClient(fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort))
	if err != nil {
		log.Fatalf("Valkey creation failed: %v", err)
	}

	defer func() { _ = valkeyClient.Close() }()

	log.Println("Connected to Valkey")

	// Connect to Postgres
	postgreConnStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.PostgresUser,
		cfg.PostgresPass,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	// Run embedded migrations
	log.Println("Running database migrations...")

	if migrateErr := db.Migrate(postgreConnStr); migrateErr != nil {
		log.Fatalf("Migration failed: %v", migrateErr)
	}

	log.Println("Database migrations applied successfully")

	postgreClient, err := pgxpool.New(ctx, postgreConnStr)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}

	defer postgreClient.Close()

	log.Println("Connected to Postgres")

	// Initialize repositories and services
	userRepo := postgres.NewUserRepository(postgreClient)
	portfolioRepo := postgres.NewPortfolioRepository(postgreClient)
	marketRepo := valkey.NewMarketRepository(valkeyClient)
	transactor := postgres.NewPgxTransactor(postgreClient)

	userService := service.NewUserService(userRepo, portfolioRepo)
	tradeService := service.NewTradeService(userRepo, portfolioRepo, marketRepo, transactor)
	marketService := service.NewMarketService(marketRepo, cfg.Tickers)

	restHandler := handler.NewRestHandler(userService, tradeService, marketService)

	// Initialize router
	router, err := api.NewRouter(restHandler, cfg)
	if err != nil {
		log.Fatalf("Failed to create router: %v", err)
	}

	// Start HTTP server in a goroutine
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	go func() {
		log.Printf("Exchange API running on :%d\n", cfg.ServerPort)

		srvErr := srv.ListenAndServe()

		if srvErr != nil && !errors.Is(srvErr, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", srvErr)
		}
	}()

	// Init gRPC server
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := googlegrpc.NewServer(
		googlegrpc.UnaryInterceptor(middleware.GrpcAuthInterceptor),
	)
	exchangeServer := grpcapi.NewExchangeServer(tradeService, marketService)
	exchange.RegisterExchangeServiceServer(grpcServer, exchangeServer)

	go func() {
		log.Printf("Exchange gRPC running on :%d\n", 50051)

		err := grpcServer.Serve(grpcListener)
		if err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	grpcServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
