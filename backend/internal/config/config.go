package config

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	Tickers        []string
	ServerPort     int
	RedisHost      string
	RedisPort      int
	ClientPort     int
	FetchInterval  time.Duration
	SleepInterval  time.Duration
	FinnhubKey     string
	FinnhubTimeout time.Duration
	PostgresUser   string
	PostgresPass   string
	PostgresDB     string
	PostgresPort   int
	PostgresHost   string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Tickers: []string{
			"AAPL",
			"BINANCE:BTCUSDT",
		},
		ServerPort:     getEnvInt("SERVER_PORT", 8081),
		RedisHost:      getEnvString("REDIS_HOST", "localhost"),
		RedisPort:      getEnvInt("REDIS_PORT", 6379),
		ClientPort:     getEnvInt("CLIENT_PORT", 5173),
		FinnhubKey:     getEnvString("FINNHUB_API_KEY", ""),
		FetchInterval:  3 * time.Second,
		SleepInterval:  2 * time.Second,
		FinnhubTimeout: 10 * time.Second,
		PostgresUser:   getEnvString("POSTGRES_USER", "postgres"),
		PostgresPass:   getEnvString("POSTGRES_PASSWORD", "postgres"),
		PostgresDB:     getEnvString("POSTGRES_DB", "ticker_rush"),
		PostgresPort:   getEnvInt("POSTGRES_PORT", 5432),
		PostgresHost:   getEnvString("POSTGRES_HOST", "localhost"),
	}
	log.Printf("config loaded: %s %s %s", cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresDB)

	return cfg, nil
}

func (c *Config) ValidateFetcher() error {
	if c.FinnhubKey == "" {
		return fmt.Errorf("FINNHUB_API_KEY is not set")
	}
	return nil
}

func getEnvInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	// Remove any potential colon prefix if it exists, just in case
	if val[0] == ':' {
		val = val[1:]
	}
	// Parse int
	var i int
	if _, err := fmt.Sscanf(val, "%d", &i); err != nil {
		return defaultValue
	}
	return i
}

func getEnvString(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}
