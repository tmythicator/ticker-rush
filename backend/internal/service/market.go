package service

import (
	"context"
	"slices"

	valkey "github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/model"
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

func (s *MarketService) GetQuote(ctx context.Context, symbol string) (any, error) {
	if !slices.Contains(s.allowedTickers, symbol) {
		return nil, model.ErrSymbolNotAllowed
	}
	return s.marketRepo.GetQuote(ctx, symbol)
}
