package service

import (
	"context"
	"slices"

	"github.com/redis/go-redis/v9"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
	"github.com/tmythicator/ticker-rush/backend/internal/domain"
)

// MarketRepository defines the interface for market data persistence.
type MarketRepository interface {
	GetQuote(ctx context.Context, symbol string) (*domain.Quote, error)
	SaveQuote(ctx context.Context, quote *domain.Quote) error
	SubscribeToQuotes(ctx context.Context, symbol string) *redis.PubSub
}

// HistoryRepository defines the interface for historical market data persistence.
type HistoryRepository interface {
	SaveQuote(ctx context.Context, quote *domain.Quote) error
	GetHistory(ctx context.Context, symbol string, limit int) ([]*domain.Quote, error)
}

// Market handles stock market data operations.
type Market struct {
	marketRepo  MarketRepository
	historyRepo HistoryRepository
	ladderRepo  LadderRepository
}

// NewMarket creates a new instance of Market.
func NewMarket(
	marketRepo MarketRepository,
	historyRepo HistoryRepository,
	ladderRepo LadderRepository,
) *Market {
	return &Market{
		marketRepo:  marketRepo,
		historyRepo: historyRepo,
		ladderRepo:  ladderRepo,
	}
}

// GetQuote gets a quote for a symbol, if allowed.
func (s *Market) GetQuote(ctx context.Context, symbol string) (*domain.Quote, error) {
	allowed, err := s.isSymbolAllowed(ctx, symbol)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, apperrors.ErrSymbolNotAllowed
	}

	return s.marketRepo.GetQuote(ctx, symbol)
}

// SubscribeToQuotes returns a PubSub for real-time quotes, if allowed.
func (s *Market) SubscribeToQuotes(
	ctx context.Context,
	symbol string,
) (*redis.PubSub, error) {
	allowed, err := s.isSymbolAllowed(ctx, symbol)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, apperrors.ErrSymbolNotAllowed
	}

	return s.marketRepo.SubscribeToQuotes(ctx, symbol), nil
}

// GetHistory retrieves historical quotes for a symbol.
func (s *Market) GetHistory(ctx context.Context, symbol string, limit int) ([]*domain.Quote, error) {
	allowed, err := s.isSymbolAllowed(ctx, symbol)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, apperrors.ErrSymbolNotAllowed
	}

	return s.historyRepo.GetHistory(ctx, symbol, limit)
}

func (s *Market) isSymbolAllowed(ctx context.Context, symbol string) (bool, error) {
	ladderID, err := s.ladderRepo.GetActiveLadder(ctx)
	if err != nil {
		return false, err
	}

	allowedTickers, err := s.ladderRepo.GetAllowedTickers(ctx, ladderID)
	if err != nil {
		return false, err
	}

	return slices.ContainsFunc(allowedTickers, func(t *domain.TickerInfo) bool {
		return t.Symbol == symbol
	}), nil
}
