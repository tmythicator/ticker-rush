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

// QuoteProvider defines the interface for fetching quotes.
type QuoteProvider interface {
	GetQuote(ctx context.Context, symbol string) (*exchange.Quote, error)
}

// MarketFetcher is a worker that fetches market data.
type MarketFetcher struct {
	client QuoteProvider
	repo   *redis.MarketRepository
}

// NewMarketFetcher creates a new instance of MarketFetcher.
func NewMarketFetcher(client QuoteProvider, repo *redis.MarketRepository) *MarketFetcher {
	return &MarketFetcher{
		client: client,
		repo:   repo,
	}
}

// RunLoop begins the fetching loop for a list of symbols.
// It processes tickers serially to ensure rate limits are respected.
func (w *MarketFetcher) RunLoop(
	ctx context.Context,
	symbols []string,
	fetchInterval time.Duration,
	wg *sync.WaitGroup,
) {
	if len(symbols) == 0 {
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Local cache for deduplication
		lastQuotes := make(map[string]*exchange.Quote)

		// Create a ticker for the interval between requests
		// Note: This interval is the pause *between* requests, guaranteeing rate limit compliance.
		delay := time.NewTicker(fetchInterval)
		defer delay.Stop()

		for {
			for _, symbol := range symbols {
				// Check for cancellation before each fetch
				select {
				case <-ctx.Done():
					log.Printf("Worker loop for %v stopped", symbols)

					return
				default:
				}

				// Process the ticker
				q, err := w.processTicker(ctx, symbol, lastQuotes[symbol])
				if err != nil {
					log.Printf("[%s] Fetch failed: %v", symbol, err)
					// We continue to the next ticker even if one fails
				} else {
					// Update local cache if successful
					lastQuotes[symbol] = q
				}

				// Wait for the rate limit interval
				select {
				case <-delay.C:
					// Continue
				case <-ctx.Done():
					log.Printf("Worker loop for %v stopped", symbols)

					return
				}
			}
		}
	}()
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
