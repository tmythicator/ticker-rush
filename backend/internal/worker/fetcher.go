package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/internal/clients/finnhub"
)

type MarketFetcher struct {
	client *finnhub.Client
	rdb    *redis.Client
}

func NewMarketFetcher(client *finnhub.Client, rdb *redis.Client) *MarketFetcher {
	return &MarketFetcher{
		client: client,
		rdb:    rdb,
	}
}

func (w *MarketFetcher) Start(ctx context.Context, symbol string, fetchInterval time.Duration) {
	ticker := time.NewTicker(fetchInterval)
	w.processTicker(ctx, symbol)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.processTicker(ctx, symbol)
			case <-ctx.Done():
				log.Printf("🛑 Worker for %s stopped", symbol)
				return
			}
		}
	}()
}

func (w *MarketFetcher) processTicker(ctx context.Context, symbol string) {
	quote, err := w.client.GetQuote(ctx, symbol)
	if err != nil {
		log.Printf("⚠️ [%s] Fetch failed: %v", symbol, err)
		return
	}

	quote.Price = float64(int(quote.Price*100)) / 100

	jsonBytes, _ := json.Marshal(quote)
	if err := w.rdb.Set(ctx, "market:"+symbol, jsonBytes, 0).Err(); err != nil {
		log.Printf("❌ [%s] Redis Error: %v", symbol, err)
		return
	}

	log.Printf("✅ [%s] Updated: $%.2f", quote.Symbol, quote.Price)
}
