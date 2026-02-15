// Package config handles application configuration.
package config

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config holds the application configuration.
type Config struct {
	Tickers        []string      `env:"TICKERS" envDefault:"AAPL,BINANCE:BTCUSDT" envSeparator:","`
	ServerPort     int           `env:"SERVER_PORT" envDefault:"8081"`
	RedisHost      string        `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort      int           `env:"REDIS_PORT" envDefault:"6379"`
	ClientPort     int           `env:"CLIENT_PORT" envDefault:"5173"`
	FetchInterval  time.Duration `env:"FETCH_INTERVAL" envDefault:"3s"`
	SleepInterval  time.Duration `env:"SLEEP_INTERVAL" envDefault:"2s"`
	FinnhubKey     string        `env:"FINNHUB_API_KEY"`
	FinnhubTimeout time.Duration `env:"FINNHUB_TIMEOUT" envDefault:"10s"`
	PostgresUser   string        `env:"POSTGRES_USER" envDefault:"postgres"`
	PostgresPass   string        `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	PostgresDB     string        `env:"POSTGRES_DB" envDefault:"ticker_rush"`
	PostgresPort   int           `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresHost   string        `env:"POSTGRES_HOST"`
	JWTSecret      string        `env:"JWT_SECRET" envDefault:"secret"`
}

// LoadConfig loads the configuration from environment variables.
func LoadConfig() (*Config, error) {
	_ = godotenv.Load("../.env")

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	log.Println("Configuration loaded:")
	log.Printf("  TICKERS: %v", cfg.Tickers)
	log.Printf("  SERVER_PORT: %d", cfg.ServerPort)
	log.Printf("  REDIS_HOST: %s", cfg.RedisHost)
	log.Printf("  REDIS_PORT: %d", cfg.RedisPort)
	log.Printf("  CLIENT_PORT: %d", cfg.ClientPort)
	log.Printf("  FETCH_INTERVAL: %s", cfg.FetchInterval)
	log.Printf("  SLEEP_INTERVAL: %s", cfg.SleepInterval)
	log.Printf("  FINNHUB_API_KEY: %s", maskString(cfg.FinnhubKey))
	log.Printf("  FINNHUB_TIMEOUT: %s", cfg.FinnhubTimeout)
	log.Printf("  POSTGRES_USER: %s", cfg.PostgresUser)
	log.Printf("  POSTGRES_PASSWORD: %s", maskString(cfg.PostgresPass))
	log.Printf("  POSTGRES_DB: %s", cfg.PostgresDB)
	log.Printf("  POSTGRES_PORT: %d", cfg.PostgresPort)
	log.Printf("  POSTGRES_HOST: %s", cfg.PostgresHost)
	log.Printf("  JWT_SECRET: %s", maskString(cfg.JWTSecret))

	return cfg, nil
}

func maskString(s string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + "****" + s[len(s)-2:]
}

// ValidateFinnhubKey checks if the finnhub key is valid.
func (c *Config) ValidateFinnhubKey() error {
	if c.FinnhubKey == "" {
		return fmt.Errorf("FINNHUB_API_KEY is not set")
	}

	return nil
}
