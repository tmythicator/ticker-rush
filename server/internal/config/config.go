package config

import (
	"fmt"
	"os"
	"time"
)

const (
	SERVER_PORT = ":8080"
	REDIS_ADDR  = "localhost:6379"
	CLIENT_ADDR = "localhost:5173"

	FETCH_INTERVAL = 3 * time.Second
	SLEEP_INTERVAL = 2 * time.Second
)

var Tickers = []string{
	"AAPL",
	"BINANCE:BTCUSDT",
}

func GetAPIKey() (string, error) {
	LoadEnv()

	key := os.Getenv("FINNHUB_API_KEY")
	if key == "" {
		return "", fmt.Errorf("FINNHUB_API_KEY is not set.")
	}
	return key, nil
}
