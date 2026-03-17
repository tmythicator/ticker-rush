package service

import (
	"context"
	"slices"

	"github.com/redis/go-redis/v9"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/ladder/v1"
)

// MarketService handles stock market data operations.
type MarketService struct {
	marketRepo  MarketRepository
	historyRepo HistoryRepository
	ladderRepo  LadderRepository
}

// NewMarketService creates a new instance of MarketService.
func NewMarketService(
	marketRepo MarketRepository,
	historyRepo HistoryRepository,
	ladderRepo LadderRepository,
) *MarketService {
	return &MarketService{
		marketRepo:  marketRepo,
		historyRepo: historyRepo,
		ladderRepo:  ladderRepo,
	}
}

// GetQuote gets a quote for a symbol, if allowed.
func (s *MarketService) GetQuote(ctx context.Context, symbol string) (*exchange.Quote, error) {
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
func (s *MarketService) SubscribeToQuotes(
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
func (s *MarketService) GetHistory(ctx context.Context, symbol string, limit int) ([]*exchange.Quote, error) {
	allowed, err := s.isSymbolAllowed(ctx, symbol)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, apperrors.ErrSymbolNotAllowed
	}

	return s.historyRepo.GetHistory(ctx, symbol, limit)
}

func (s *MarketService) isSymbolAllowed(ctx context.Context, symbol string) (bool, error) {
	ladderID, err := s.ladderRepo.GetActiveLadder(ctx)
	if err != nil {
		return false, err
	}
	allowedTickers, err := s.ladderRepo.GetAllowedTickers(ctx, ladderID)
	if err != nil {
		return false, err
	}

	return slices.ContainsFunc(allowedTickers, func(t *ladder.TickerInfo) bool {
		return t.Symbol == symbol
	}), nil
}
