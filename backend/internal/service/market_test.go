package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
	"github.com/tmythicator/ticker-rush/backend/internal/domain"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
	"github.com/tmythicator/ticker-rush/backend/internal/service/mocks"
)

func TestMarketService_GetQuote(t *testing.T) {
	ctx := context.Background()
	symbol := "AAPL"

	t.Run("Success", func(t *testing.T) {
		mockMarketRepo := new(mocks.MockMarketRepository)
		mockHistoryRepo := new(mocks.MockHistoryRepository)
		mockLadderRepo := new(mocks.MockLadderRepository)

		expectedQuote := &domain.Quote{Symbol: symbol}

		mockLadderRepo.On("GetActiveLadder", ctx).Return(int64(1), nil)
		mockLadderRepo.On("GetAllowedTickers", ctx, int64(1)).Return([]*domain.TickerInfo{
			{Symbol: symbol},
		}, nil)
		mockMarketRepo.On("GetQuote", ctx, symbol).Return(expectedQuote, nil)

		s := service.NewMarketService(mockMarketRepo, mockHistoryRepo, mockLadderRepo)
		q, err := s.GetQuote(ctx, symbol)

		assert.NoError(t, err)
		assert.Equal(t, expectedQuote, q)
		mockLadderRepo.AssertExpectations(t)
		mockMarketRepo.AssertExpectations(t)
	})

	t.Run("SymbolNotAllowed", func(t *testing.T) {
		mockMarketRepo := new(mocks.MockMarketRepository)
		mockHistoryRepo := new(mocks.MockHistoryRepository)
		mockLadderRepo := new(mocks.MockLadderRepository)

		mockLadderRepo.On("GetActiveLadder", ctx).Return(int64(1), nil)
		mockLadderRepo.On("GetAllowedTickers", ctx, int64(1)).Return([]*domain.TickerInfo{
			{Symbol: "GOOG"},
		}, nil)

		s := service.NewMarketService(mockMarketRepo, mockHistoryRepo, mockLadderRepo)
		q, err := s.GetQuote(ctx, symbol)

		assert.ErrorIs(t, err, apperrors.ErrSymbolNotAllowed)
		assert.Nil(t, q)
		mockLadderRepo.AssertExpectations(t)
	})

	t.Run("GetActiveLadderError", func(t *testing.T) {
		mockMarketRepo := new(mocks.MockMarketRepository)
		mockHistoryRepo := new(mocks.MockHistoryRepository)
		mockLadderRepo := new(mocks.MockLadderRepository)

		expectedErr := errors.New("db error")
		mockLadderRepo.On("GetActiveLadder", ctx).Return(int64(0), expectedErr)

		s := service.NewMarketService(mockMarketRepo, mockHistoryRepo, mockLadderRepo)
		q, err := s.GetQuote(ctx, symbol)

		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, q)
		mockLadderRepo.AssertExpectations(t)
	})
}

func TestMarketService_SubscribeToQuotes(t *testing.T) {
	ctx := context.Background()
	symbol := "AAPL"

	t.Run("Success", func(t *testing.T) {
		mockMarketRepo := new(mocks.MockMarketRepository)
		mockHistoryRepo := new(mocks.MockHistoryRepository)
		mockLadderRepo := new(mocks.MockLadderRepository)

		expectedPubSub := &redis.PubSub{}

		mockLadderRepo.On("GetActiveLadder", ctx).Return(int64(1), nil)
		mockLadderRepo.On("GetAllowedTickers", ctx, int64(1)).Return([]*domain.TickerInfo{
			{Symbol: symbol},
		}, nil)
		mockMarketRepo.On("SubscribeToQuotes", ctx, symbol).Return(expectedPubSub)

		s := service.NewMarketService(mockMarketRepo, mockHistoryRepo, mockLadderRepo)
		pb, err := s.SubscribeToQuotes(ctx, symbol)

		assert.NoError(t, err)
		assert.Equal(t, expectedPubSub, pb)
		mockLadderRepo.AssertExpectations(t)
		mockMarketRepo.AssertExpectations(t)
	})

	t.Run("SymbolNotAllowed", func(t *testing.T) {
		mockMarketRepo := new(mocks.MockMarketRepository)
		mockHistoryRepo := new(mocks.MockHistoryRepository)
		mockLadderRepo := new(mocks.MockLadderRepository)

		mockLadderRepo.On("GetActiveLadder", ctx).Return(int64(1), nil)
		mockLadderRepo.On("GetAllowedTickers", ctx, int64(1)).Return([]*domain.TickerInfo{
			{Symbol: "MSFT"},
		}, nil)

		s := service.NewMarketService(mockMarketRepo, mockHistoryRepo, mockLadderRepo)
		pb, err := s.SubscribeToQuotes(ctx, symbol)

		assert.ErrorIs(t, err, apperrors.ErrSymbolNotAllowed)
		assert.Nil(t, pb)
		mockLadderRepo.AssertExpectations(t)
	})
}

func TestMarketService_GetHistory(t *testing.T) {
	ctx := context.Background()
	symbol := "AAPL"
	limit := 10

	t.Run("Success", func(t *testing.T) {
		mockMarketRepo := new(mocks.MockMarketRepository)
		mockHistoryRepo := new(mocks.MockHistoryRepository)
		mockLadderRepo := new(mocks.MockLadderRepository)

		expectedHistory := []*domain.Quote{
			{Symbol: symbol},
		}

		mockLadderRepo.On("GetActiveLadder", ctx).Return(int64(1), nil)
		mockLadderRepo.On("GetAllowedTickers", ctx, int64(1)).Return([]*domain.TickerInfo{
			{Symbol: symbol},
		}, nil)
		mockHistoryRepo.On("GetHistory", ctx, symbol, limit).Return(expectedHistory, nil)

		s := service.NewMarketService(mockMarketRepo, mockHistoryRepo, mockLadderRepo)
		h, err := s.GetHistory(ctx, symbol, limit)

		assert.NoError(t, err)
		assert.Equal(t, expectedHistory, h)
		mockLadderRepo.AssertExpectations(t)
		mockHistoryRepo.AssertExpectations(t)
	})

	t.Run("SymbolNotAllowed", func(t *testing.T) {
		mockMarketRepo := new(mocks.MockMarketRepository)
		mockHistoryRepo := new(mocks.MockHistoryRepository)
		mockLadderRepo := new(mocks.MockLadderRepository)

		mockLadderRepo.On("GetActiveLadder", ctx).Return(int64(1), nil)
		mockLadderRepo.On("GetAllowedTickers", ctx, int64(1)).Return([]*domain.TickerInfo{
			{Symbol: "GOOG"},
		}, nil)

		s := service.NewMarketService(mockMarketRepo, mockHistoryRepo, mockLadderRepo)
		h, err := s.GetHistory(ctx, symbol, limit)

		assert.ErrorIs(t, err, apperrors.ErrSymbolNotAllowed)
		assert.Nil(t, h)
		mockLadderRepo.AssertExpectations(t)
	})
}
