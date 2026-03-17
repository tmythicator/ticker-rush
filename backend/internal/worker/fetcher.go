// Package worker provides background workers.
package worker

import (
	"context"
	"log"
	"math"
	"sync"
	"time"

	"github.com/tmythicator/ticker-rush/backend/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
)

// QuoteProvider defines the interface for fetching quotes.
type QuoteProvider interface {
	GetQuote(ctx context.Context, symbol string) (*exchange.Quote, error)
}

// MarketFetcher is a worker that fetches market data.
type MarketFetcher struct {
	client      QuoteProvider
	currentRepo *redis.MarketRepository
	historyRepo service.HistoryRepository
	ladderRepo  service.LadderRepository
	source      string // "Finnhub" or "CoinGecko"
}

// NewMarketFetcher creates a new instance of MarketFetcher.
func NewMarketFetcher(
	source string,
	client QuoteProvider,
	currentRepo *redis.MarketRepository,
	historyRepo service.HistoryRepository,
	ladderRepo service.LadderRepository,
) *MarketFetcher {
	return &MarketFetcher{
		source:      source,
		client:      client,
		currentRepo: currentRepo,
		historyRepo: historyRepo,
		ladderRepo:  ladderRepo,
	}
}

// RunLoop begins the fetching loop for tickers associated with the active ladder.
// It periodically refreshes the list of tickers from the database.
func (w *MarketFetcher) RunLoop(
	ctx context.Context,
	fetchInterval time.Duration,
	refreshInterval time.Duration,
	wg *sync.WaitGroup,
) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Local cache for deduplication
		lastQuotes := make(map[string]*exchange.Quote)

		// Create tickers for the intervals
		delay := time.NewTicker(fetchInterval)
		defer delay.Stop()

		refreshTimer := time.NewTicker(refreshInterval)
		defer refreshTimer.Stop()

		var symbols []string

		// Initial fetch
		symbols = w.refreshTickers(ctx)

		for {
			if len(symbols) == 0 {
				select {
				case <-ctx.Done():
					return
				case <-refreshTimer.C:
					symbols = w.refreshTickers(ctx)
				case <-time.After(5 * time.Second):
				}
				continue
			}

			for i := 0; i < len(symbols); i++ {
				symbol := symbols[i]

				select {
				case <-ctx.Done():
					log.Printf("[Fetcher:%s] Stopped", w.source)
					return
				case <-refreshTimer.C:
					symbols = w.refreshTickers(ctx)
					i = -1
					continue
				case <-delay.C:
					q, err := w.processTicker(ctx, symbol, lastQuotes[symbol])
					if err != nil {
						log.Printf("[%s] Fetch failed: %v", symbol, err)
					} else if q != nil {
						lastQuotes[symbol] = q
					}
				}
			}
		}
	}()
}

func (w *MarketFetcher) refreshTickers(ctx context.Context) []string {
	activeLadderID, err := w.ladderRepo.GetActiveLadder(ctx)
	if err != nil {
		log.Printf("[Fetcher:%s] Failed to get active ladder: %v", w.source, err)
		return nil
	}

	tickers, err := w.ladderRepo.GetAllowedTickers(ctx, activeLadderID)
	if err != nil {
		log.Printf("[Fetcher:%s] Failed to get ladder tickers: %v", w.source, err)
		return nil
	}

	var filtered []string
	for _, t := range tickers {
		if t.Source == w.source {
			filtered = append(filtered, t.Symbol)
		}
	}

	if len(filtered) > 0 {
		log.Printf("[Fetcher:%s] Tracking %d tickers", w.source, len(filtered))
	}

	return filtered
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

	if err := w.currentRepo.SaveQuote(ctx, quote); err != nil {
		log.Printf("[%s] Redis Error: %v", symbol, err)

		return nil, err
	}

	// Async save to history (don't block the loop)
	go func(q *exchange.Quote) {
		saveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := w.historyRepo.SaveQuote(saveCtx, q); err != nil {
			log.Printf("[%s] History Clean Error: %v", q.Symbol, err)
		}
	}(quote)

	log.Printf(
		"[%s] Updated: $%.2f (ts: %d)",
		quote.GetSymbol(),
		quote.GetPrice(),
		quote.GetTimestamp(),
	)

	return quote, nil
}
