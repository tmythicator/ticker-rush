package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tmythicator/ticker-rush/backend/internal/domain"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/service/mocks"
)

// MockQuoteProvider mocks the QuoteProvider interface.
type MockQuoteProvider struct {
	mock.Mock
}

func (m *MockQuoteProvider) GetQuote(ctx context.Context, symbol string) (*exchange.Quote, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*exchange.Quote), args.Error(1)
}

func TestMarketFetcher_ProcessTicker(t *testing.T) {
	ctx := context.Background()
	symbol := "AAPL"
	source := "Finnhub"

	cfg := &FetcherConfig{
		FetchInterval:   10 * time.Second,
		RefreshInterval: 1 * time.Minute,
		RequestTimeout:  2 * time.Second,
	}

	t.Run("Success - New Quote", func(t *testing.T) {
		mockClient := new(MockQuoteProvider)
		mockMarketRepo := new(mocks.MockMarketRepository)
		mockHistoryRepo := new(mocks.MockHistoryRepository)
		mockLadderRepo := new(mocks.MockLadderRepository)

		fetchedQuote := &exchange.Quote{
			Symbol:    symbol,
			Price:     150.256, // Will be rounded to 150.26
			Timestamp: time.Now().Unix(),
			Source:    source,
			IsClosed:  false,
		}

		mockClient.On("GetQuote", mock.Anything, symbol).Return(fetchedQuote, nil)
		mockMarketRepo.On("SaveQuote", ctx, mock.MatchedBy(func(q *domain.Quote) bool {
			return q.Symbol == symbol && q.Price.InexactFloat64() == 150.26 && q.Source == source
		})).Return(nil)
		mockHistoryRepo.On("SaveQuote", mock.Anything, mock.MatchedBy(func(q *domain.Quote) bool {
			return q.Symbol == symbol && q.Price.InexactFloat64() == 150.26 && q.Source == source
		})).Return(nil)

		w := NewMarketFetcher(source, mockClient, mockMarketRepo, mockHistoryRepo, mockLadderRepo, cfg)
		lastHistorySave := make(map[string]time.Time) // Empty so time.Since is > 1 min

		res, err := w.processTicker(ctx, symbol, nil, lastHistorySave)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 150.26, res.Price)
		mockClient.AssertExpectations(t)
		mockMarketRepo.AssertExpectations(t)
		mockHistoryRepo.AssertExpectations(t)
	})

	t.Run("Success - Duplicate Quote (Skip Save)", func(t *testing.T) {
		mockClient := new(MockQuoteProvider)
		mockMarketRepo := new(mocks.MockMarketRepository)
		mockHistoryRepo := new(mocks.MockHistoryRepository)
		mockLadderRepo := new(mocks.MockLadderRepository)

		timestamp := time.Now().Unix()
		fetchedQuote := &exchange.Quote{
			Symbol:    symbol,
			Price:     150.26,
			Timestamp: timestamp,
			Source:    source,
		}
		lastQuote := &exchange.Quote{
			Symbol:    symbol,
			Price:     150.26,
			Timestamp: timestamp,
			Source:    source,
		}

		mockClient.On("GetQuote", mock.Anything, symbol).Return(fetchedQuote, nil)

		w := NewMarketFetcher(source, mockClient, mockMarketRepo, mockHistoryRepo, mockLadderRepo, cfg)
		lastHistorySave := make(map[string]time.Time)

		res, err := w.processTicker(ctx, symbol, lastQuote, lastHistorySave)

		assert.NoError(t, err)
		assert.Equal(t, lastQuote, res)
		mockClient.AssertExpectations(t)
		mockMarketRepo.AssertNotCalled(t, "SaveQuote")
		mockHistoryRepo.AssertNotCalled(t, "SaveQuote")
	})

	t.Run("Success - New Price Save Current Only", func(t *testing.T) {
		mockClient := new(MockQuoteProvider)
		mockMarketRepo := new(mocks.MockMarketRepository)
		mockHistoryRepo := new(mocks.MockHistoryRepository)
		mockLadderRepo := new(mocks.MockLadderRepository)

		fetchedQuote := &exchange.Quote{
			Symbol:    symbol,
			Price:     151.0,
			Timestamp: time.Now().Unix(),
			Source:    source,
		}
		lastQuote := &exchange.Quote{
			Symbol:    symbol,
			Price:     150.0,
			Timestamp: time.Now().Unix() - 10,
			Source:    source,
		}

		mockClient.On("GetQuote", mock.Anything, symbol).Return(fetchedQuote, nil)
		mockMarketRepo.On("SaveQuote", ctx, mock.Anything).Return(nil)

		w := NewMarketFetcher(source, mockClient, mockMarketRepo, mockHistoryRepo, mockLadderRepo, cfg)
		lastHistorySave := map[string]time.Time{
			symbol: time.Now(), // Less than a minute ago
		}

		res, err := w.processTicker(ctx, symbol, lastQuote, lastHistorySave)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		mockClient.AssertExpectations(t)
		mockMarketRepo.AssertExpectations(t)
		mockHistoryRepo.AssertNotCalled(t, "SaveQuote")
	})

	t.Run("Error - Client Failure", func(t *testing.T) {
		mockClient := new(MockQuoteProvider)
		mockMarketRepo := new(mocks.MockMarketRepository)
		mockHistoryRepo := new(mocks.MockHistoryRepository)
		mockLadderRepo := new(mocks.MockLadderRepository)

		expectedErr := errors.New("client timeout")
		mockClient.On("GetQuote", mock.Anything, symbol).Return(nil, expectedErr)

		w := NewMarketFetcher(source, mockClient, mockMarketRepo, mockHistoryRepo, mockLadderRepo, cfg)
		lastHistorySave := make(map[string]time.Time)

		res, err := w.processTicker(ctx, symbol, nil, lastHistorySave)

		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, res)
		mockClient.AssertExpectations(t)
	})
}

func TestMarketFetcher_RefreshTickers(t *testing.T) {
	ctx := context.Background()
	source := "Finnhub"

	cfg := &FetcherConfig{
		FetchInterval:   10 * time.Second,
		RefreshInterval: 1 * time.Minute,
		RequestTimeout:  2 * time.Second,
	}

	t.Run("Success - Filter by source", func(t *testing.T) {
		mockClient := new(MockQuoteProvider)
		mockMarketRepo := new(mocks.MockMarketRepository)
		mockHistoryRepo := new(mocks.MockHistoryRepository)
		mockLadderRepo := new(mocks.MockLadderRepository)

		mockLadderRepo.On("GetActiveLadder", ctx).Return(int64(1), nil)
		mockLadderRepo.On("GetAllowedTickers", ctx, int64(1)).Return([]*domain.TickerInfo{
			{Symbol: "AAPL", Source: "Finnhub"},
			{Symbol: "BTCUSDT", Source: "CoinGecko"},
			{Symbol: "MSFT", Source: "Finnhub"},
		}, nil)

		w := NewMarketFetcher(source, mockClient, mockMarketRepo, mockHistoryRepo, mockLadderRepo, cfg)
		res := w.refreshTickers(ctx)

		assert.Equal(t, []string{"AAPL", "MSFT"}, res)
		mockLadderRepo.AssertExpectations(t)
	})

	t.Run("Error - Active Ladder Failed", func(t *testing.T) {
		mockClient := new(MockQuoteProvider)
		mockMarketRepo := new(mocks.MockMarketRepository)
		mockHistoryRepo := new(mocks.MockHistoryRepository)
		mockLadderRepo := new(mocks.MockLadderRepository)

		mockLadderRepo.On("GetActiveLadder", ctx).Return(int64(0), errors.New("db error"))

		w := NewMarketFetcher(source, mockClient, mockMarketRepo, mockHistoryRepo, mockLadderRepo, cfg)
		res := w.refreshTickers(ctx)

		assert.Nil(t, res)
		mockLadderRepo.AssertExpectations(t)
	})
}
