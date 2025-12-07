package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tmythicator/ticker-rush/server/internal/api"
	"github.com/tmythicator/ticker-rush/server/internal/api/handler"
	"github.com/tmythicator/ticker-rush/server/internal/config"
	"github.com/tmythicator/ticker-rush/server/internal/repository/postgres"
	valkey "github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/internal/service"
)

var (
	ctx = context.Background()
)

func main() {
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
	postgreConnStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)
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

	userService := service.NewUserService(userRepo, portfolioRepo)
	tradeService := service.NewTradeService(userRepo, portfolioRepo, marketRepo, postgreClient)
	marketService := service.NewMarketService(marketRepo, cfg.Tickers)

	restHandler := handler.NewRestHandler(userService, tradeService, marketService)

	// Initialize router
	router, err := api.NewRouter(restHandler, cfg)
	if err != nil {
		log.Fatalf("Failed to create router: %v", err)
	}

	log.Printf("Exchange API running on :%d\n", cfg.ServerPort)
	if err := router.Run(fmt.Sprintf(":%d", cfg.ServerPort)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
