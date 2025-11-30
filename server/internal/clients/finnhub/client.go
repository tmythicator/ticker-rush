package finnhub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/internal/config"
	"github.com/tmythicator/ticker-rush/server/model"
)

func UpdateMarketData(ctx context.Context, symbol string, apiKey string, rdb *redis.Client) {
	ticker := time.NewTicker(config.FETCH_INTERVAL)

	// Pattern: Immediate + Interval
	updatePrice(ctx, symbol, apiKey, rdb)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				updatePrice(ctx, symbol, apiKey, rdb)
			case <-ctx.Done():
				log.Printf("Worker for %s stopped", symbol)
				return
			}
		}
	}()
}

func updatePrice(ctx context.Context, symbol string, apiKey string, rdb *redis.Client) {
	url := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", symbol, apiKey)
	httpClient := http.Client{Timeout: config.FETCH_INTERVAL}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error fetching quote for %s: %v", symbol, err)
		return
	}
	defer resp.Body.Close()

	var fq model.FinnhubQuote
	if err := json.NewDecoder(resp.Body).Decode(&fq); err != nil {
		log.Printf("Error decoding finnhub quote for %s: %v", symbol, err)
		return
	}

	if fq.CurrentPrice == 0 {
		return
	}

	ts := fq.Timestamp
	if ts == 0 {
		ts = time.Now().Unix()
	}

	fq.CurrentPrice = float64(int(fq.CurrentPrice*100)) / 100

	quote := model.Quote{
		Symbol:    symbol,
		Price:     fq.CurrentPrice,
		Timestamp: ts,
	}

	jsonBytes, _ := json.Marshal(quote)
	err = rdb.Set(ctx, "market:"+symbol, jsonBytes, 0).Err()
	if err != nil {
		log.Printf("Redis Write Error: %v", err)
		return
	}

	log.Printf("✅ Market Updated: %s @ $%.2f", quote.Symbol, quote.Price)
}
