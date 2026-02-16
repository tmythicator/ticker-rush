package service

import (
	"context"
	"slices"

	"github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/internal/apperrors"
	"github.com/tmythicator/ticker-rush/server/internal/proto/exchange/v1"
)

// MarketService handles stock market data operations.
type MarketService struct {
	marketRepo     MarketRepository
	historyRepo    HistoryRepository
	allowedTickers []string
}

// NewMarketService creates a new instance of MarketService.
func NewMarketService(
	marketRepo MarketRepository,
	historyRepo HistoryRepository,
	allowedTickers []string,
) *MarketService {
	return &MarketService{
		marketRepo:     marketRepo,
		historyRepo:    historyRepo,
		allowedTickers: allowedTickers,
	}
}

// GetQuote gets a quote for a symbol, if allowed.
func (s *MarketService) GetQuote(ctx context.Context, symbol string) (*exchange.Quote, error) {
	if !slices.Contains(s.allowedTickers, symbol) {
		return nil, apperrors.ErrSymbolNotAllowed
	}

	return s.marketRepo.GetQuote(ctx, symbol)
}

// SubscribeToQuotes returns a PubSub for real-time quotes, if allowed.
func (s *MarketService) SubscribeToQuotes(
	ctx context.Context,
	symbol string,
) (*redis.PubSub, error) {
	if !slices.Contains(s.allowedTickers, symbol) {
		return nil, apperrors.ErrSymbolNotAllowed
	}

	return s.marketRepo.SubscribeToQuotes(ctx, symbol), nil
}

// GetHistory retrieves historical quotes for a symbol.
func (s *MarketService) GetHistory(ctx context.Context, symbol string, limit int) ([]*exchange.Quote, error) {
	if !slices.Contains(s.allowedTickers, symbol) {
		return nil, apperrors.ErrSymbolNotAllowed
	}

	return s.historyRepo.GetHistory(ctx, symbol, limit)
}
