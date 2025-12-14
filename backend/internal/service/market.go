package service

import (
	"context"
	"slices"

	"github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/internal/apperrors"
	"github.com/tmythicator/ticker-rush/server/internal/proto/exchange"
	valkey "github.com/tmythicator/ticker-rush/server/internal/repository/redis"
)

type MarketService struct {
	marketRepo     *valkey.MarketRepository
	allowedTickers []string
}

func NewMarketService(marketRepo *valkey.MarketRepository, allowedTickers []string) *MarketService {
	return &MarketService{
		marketRepo:     marketRepo,
		allowedTickers: allowedTickers,
	}
}

func (s *MarketService) GetQuote(ctx context.Context, symbol string) (*exchange.Quote, error) {
	if !slices.Contains(s.allowedTickers, symbol) {
		return nil, apperrors.ErrSymbolNotAllowed
	}
	return s.marketRepo.GetQuote(ctx, symbol)
}

func (s *MarketService) SubscribeToQuotes(ctx context.Context, symbol string) (*redis.PubSub, error) {
	if !slices.Contains(s.allowedTickers, symbol) {
		return nil, apperrors.ErrSymbolNotAllowed
	}
	return s.marketRepo.SubscribeToQuotes(ctx, symbol), nil
}
