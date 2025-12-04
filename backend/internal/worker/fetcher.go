package worker

import (
	"context"
	"log"
	"time"

	"github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/model"
)

type FinnhubClient interface {
	GetQuote(ctx context.Context, symbol string) (*model.Quote, error)
}

type MarketFetcher struct {
	client FinnhubClient
	repo   *redis.MarketRepository
}

func NewMarketFetcher(client FinnhubClient, repo *redis.MarketRepository) *MarketFetcher {
	return &MarketFetcher{
		client: client,
		repo:   repo,
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

	if err := w.repo.SaveQuote(ctx, quote); err != nil {
		log.Printf("❌ [%s] Redis Error: %v", symbol, err)
		return
	}

	log.Printf("✅ [%s] Updated: $%.2f", quote.Symbol, quote.Price)
}
