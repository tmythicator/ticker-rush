// Package worker provides background workers.
package worker

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/tmythicator/ticker-rush/backend/internal/domain"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
)

// QuoteProvider defines the interface for fetching quotes.
type QuoteProvider interface {
	GetQuote(ctx context.Context, symbol string) (*exchange.Quote, error)
}

// FetcherConfig holds configuration for the MarketFetcher.
type FetcherConfig struct {
	FetchInterval   time.Duration
	RefreshInterval time.Duration
	RequestTimeout  time.Duration
}

// MarketFetcher is a worker that fetches market data.
type MarketFetcher struct {
	client      QuoteProvider
	currentRepo service.MarketRepository
	historyRepo service.HistoryRepository
	ladderRepo  service.LadderRepository
	source      string // "Finnhub" or "CoinGecko"
	cfg         *FetcherConfig
}

// NewMarketFetcher creates a new instance of MarketFetcher.
func NewMarketFetcher(
	source string,
	client QuoteProvider,
	currentRepo service.MarketRepository,
	historyRepo service.HistoryRepository,
	ladderRepo service.LadderRepository,
	cfg *FetcherConfig,
) *MarketFetcher {
	return &MarketFetcher{
		source:      source,
		client:      client,
		currentRepo: currentRepo,
		historyRepo: historyRepo,
		ladderRepo:  ladderRepo,
		cfg:         cfg,
	}
}

// Start begins the fetching loop for tickers associated with the active ladder.
func (w *MarketFetcher) Start(ctx context.Context) error {
	refreshTicker := time.NewTicker(w.cfg.RefreshInterval)
	defer refreshTicker.Stop()

	fetchTicker := time.NewTicker(w.cfg.FetchInterval)
	defer fetchTicker.Stop()

	lastQuotes := make(map[string]*exchange.Quote)
	lastHistorySave := make(map[string]time.Time)

	symbols := w.refreshTickers(ctx)
	idx := 0

	for {
		select {
		case <-ctx.Done():
			log.Printf("[Fetcher:%s] Stopping...", w.source)

			return ctx.Err()

		case <-refreshTicker.C:
			symbols = w.refreshTickers(ctx)

		case <-fetchTicker.C:
			if len(symbols) == 0 {
				continue
			}

			idx = idx % len(symbols)
			symbol := symbols[idx]
			idx++

			q, err := w.processTicker(ctx, symbol, lastQuotes[symbol], lastHistorySave)
			if err != nil {
				log.Printf("[%s] Fetch failed: %v", symbol, err)
			} else if q != nil {
				lastQuotes[symbol] = q
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
	lastHistorySave map[string]time.Time,
) (*exchange.Quote, error) {
	tracer := otel.Tracer("market-fetcher")
	ctx, span := tracer.Start(ctx, "MarketFetcher.processTicker")
	defer span.End()

	span.SetAttributes(
		attribute.String("fetcher.source", w.source),
		attribute.String("fetcher.symbol", symbol),
	)

	fetchCtx, cancel := context.WithTimeout(ctx, w.cfg.RequestTimeout)
	defer cancel()

	quote, err := w.client.GetQuote(fetchCtx, symbol)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	quote.Price = math.Round(quote.GetPrice()*100) / 100

	isClosed := domain.CalculateIsClosed(quote.IsClosed, quote.GetTimestamp().AsTime())

	span.SetAttributes(
		attribute.Float64("fetcher.price", quote.Price),
		attribute.Int64("fetcher.timestamp", quote.GetTimestamp().GetSeconds()),
		attribute.Bool("fetcher.is_closed", isClosed),
	)

	if lastQuote != nil && quote.GetPrice() == lastQuote.GetPrice() &&
		quote.GetTimestamp().GetSeconds() == lastQuote.GetTimestamp().GetSeconds() &&
		quote.GetTimestamp().GetNanos() == lastQuote.GetTimestamp().GetNanos() {
		span.SetAttributes(attribute.Bool("fetcher.skipped_save", true))

		return lastQuote, nil
	}

	domainQuote := &domain.Quote{
		Symbol:    quote.Symbol,
		Price:     decimal.NewFromFloat(quote.Price),
		Timestamp: quote.GetTimestamp().AsTime(),
		Source:    quote.Source,
		IsClosed:  isClosed,
	}

	if err := w.currentRepo.SaveQuote(ctx, domainQuote); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Printf("[%s] Current Save Error: %v", quote.Symbol, err)

		return nil, err
	}

	// Only save to history if a minute has passed to aggregate fast writes
	if time.Since(lastHistorySave[symbol]) >= time.Minute {
		span.SetAttributes(attribute.Bool("fetcher.saved_history", true))
		saveCtx, historyCancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := w.historyRepo.SaveQuote(saveCtx, domainQuote); err != nil {
			log.Printf("[%s] History Save Error: %v", quote.Symbol, err)
		} else {
			lastHistorySave[symbol] = time.Now()
		}
		historyCancel()
	}

	log.Printf(
		"[%s] Updated: $%.2f (ts: %d)",
		quote.GetSymbol(),
		quote.GetPrice(),
		quote.GetTimestamp().GetSeconds(),
	)

	return quote, nil
}
