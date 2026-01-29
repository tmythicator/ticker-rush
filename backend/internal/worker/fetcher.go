// Package worker provides background workers.
package worker

import (
	"context"
	"log"
	"math"
	"sync"
	"time"

	"github.com/tmythicator/ticker-rush/server/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/server/internal/repository/redis"
)

// FinnhubClient defines the interface for fetching quotes.
type FinnhubClient interface {
	GetQuote(ctx context.Context, symbol string) (*exchange.Quote, error)
}

// MarketFetcher is a worker that fetches market data.
type MarketFetcher struct {
	client FinnhubClient
	repo   *redis.MarketRepository
}

// NewMarketFetcher creates a new instance of MarketFetcher.
func NewMarketFetcher(client FinnhubClient, repo *redis.MarketRepository) *MarketFetcher {
	return &MarketFetcher{
		client: client,
		repo:   repo,
	}
}

// Start begins the fetching loop for a symbol.
func (w *MarketFetcher) Start(
	ctx context.Context,
	symbol string,
	fetchInterval time.Duration,
	wg *sync.WaitGroup,
) {
	ticker := time.NewTicker(fetchInterval)

	// Initial fetch
	lastQuote, err := w.processTicker(ctx, symbol, nil)
	if err != nil {
		log.Printf("[%s] Initial fetch failed (will retry in loop): %v", symbol, err)
	}

	wg.Go(func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				q, err := w.processTicker(ctx, symbol, lastQuote)
				if err != nil {
					log.Printf("[%s] Tick skipped: %v", symbol, err)

					continue
				}

				lastQuote = q
			case <-ctx.Done():
				log.Printf("Worker for %s stopped", symbol)

				return
			}
		}
	})
}

func (w *MarketFetcher) processTicker(
	ctx context.Context,
	symbol string,
	lastQuote *exchange.Quote,
) (*exchange.Quote, error) {
	quote, err := w.client.GetQuote(ctx, symbol)
	if err != nil {
		return nil, err
	}

	quote.Price = math.Round(quote.GetPrice()*100) / 100

	if lastQuote != nil && quote.GetPrice() == lastQuote.GetPrice() &&
		quote.GetTimestamp() == lastQuote.GetTimestamp() {
		return lastQuote, nil
	}

	if err := w.repo.SaveQuote(ctx, quote); err != nil {
		log.Printf("[%s] Redis Error: %v", symbol, err)

		return nil, err
	}

	log.Printf(
		"[%s] Updated: $%.2f (ts: %d)",
		quote.GetSymbol(),
		quote.GetPrice(),
		quote.GetTimestamp(),
	)

	return quote, nil
}
