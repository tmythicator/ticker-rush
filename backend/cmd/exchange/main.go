package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tmythicator/ticker-rush/server/db"
	"github.com/tmythicator/ticker-rush/server/internal/api"
	"github.com/tmythicator/ticker-rush/server/internal/api/handler"
	"github.com/tmythicator/ticker-rush/server/internal/config"
	"github.com/tmythicator/ticker-rush/server/internal/repository/postgres"
	valkey "github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/internal/service"
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
	postgreConnStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)

	// Run embedded migrations
	log.Println("Running database migrations...")
	if err := db.Migrate(postgreConnStr); err != nil {
		log.Fatalf("Migration failed: %v", err)
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

	// Start server in a goroutine
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	go func() {
		log.Printf("Exchange API running on :%d\n", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
