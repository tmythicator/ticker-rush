// Package worker provides background workers.
package worker

import (
	"context"
	"log"
	"math"
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
	client          QuoteProvider
	currentRepo     *redis.MarketRepository
	historyRepo     service.HistoryRepository
	ladderRepo      service.LadderRepository
	source          string // "Finnhub" or "CoinGecko"
	fetchInterval   time.Duration
	refreshInterval time.Duration
	requestTimeout  time.Duration
}

// NewMarketFetcher creates a new instance of MarketFetcher.
func NewMarketFetcher(
	source string,
	client QuoteProvider,
	currentRepo *redis.MarketRepository,
	historyRepo service.HistoryRepository,
	ladderRepo service.LadderRepository,
	fetchInterval time.Duration,
	refreshInterval time.Duration,
	requestTimeout time.Duration,
) *MarketFetcher {
	return &MarketFetcher{
		source:          source,
		client:          client,
		currentRepo:     currentRepo,
		historyRepo:     historyRepo,
		ladderRepo:      ladderRepo,
		fetchInterval:   fetchInterval,
		refreshInterval: refreshInterval,
		requestTimeout:  requestTimeout,
	}
}

// Start begins the fetching loop for tickers associated with the active ladder.
func (w *MarketFetcher) Start(ctx context.Context) error {
	refreshTicker := time.NewTicker(w.refreshInterval)
	defer refreshTicker.Stop()

	fetchTimer := time.NewTimer(0)
	defer fetchTimer.Stop()

	lastQuotes := make(map[string]*exchange.Quote)

	for {
		symbols := w.refreshTickers(ctx)
		if len(symbols) == 0 {
			select {
			case <-ctx.Done():
				log.Printf("[Fetcher:%s] Stopping...", w.source)

				return ctx.Err()
			case <-refreshTicker.C:
				continue
			}
		}

		for i := 0; i < len(symbols); i++ {
			symbol := symbols[i]

			select {
			case <-ctx.Done():
				log.Printf("[Fetcher:%s] Stopping...", w.source)

				return ctx.Err()

			case <-refreshTicker.C:
				symbols = w.refreshTickers(ctx)
				i = -1

				continue

			case <-fetchTimer.C:
				q, err := w.processTicker(ctx, symbol, lastQuotes[symbol])
				if err != nil {
					log.Printf("[%s] Fetch failed: %v", symbol, err)
				} else if q != nil {
					lastQuotes[symbol] = q
				}
				fetchTimer.Reset(w.fetchInterval)
			}
		}
	}
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
	fetchCtx, cancel := context.WithTimeout(ctx, w.requestTimeout)
	defer cancel()

	quote, err := w.client.GetQuote(fetchCtx, symbol)
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
