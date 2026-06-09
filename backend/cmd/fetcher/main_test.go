package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	go_redis "github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/tmythicator/ticker-rush/backend/internal/clients/finnhub"
	"github.com/tmythicator/ticker-rush/backend/internal/domain"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/backend/internal/worker"
)

// MockFinnhubClient mocks the Finnhub API.
type MockFinnhubClient struct {
	FinnhubQuote finnhub.Response
}

func (m *MockFinnhubClient) GetQuote(ctx context.Context, symbol string) (*exchange.Quote, error) {
	return &exchange.Quote{
		Symbol:    symbol,
		Price:     m.FinnhubQuote.CurrentPrice,
		Timestamp: m.FinnhubQuote.Timestamp,
	}, nil
}

// MockHistoryRepository mocks the history storage.
type MockHistoryRepository struct{}

func (m *MockHistoryRepository) SaveQuote(ctx context.Context, quote *domain.Quote) error {
	return nil
}

func (m *MockHistoryRepository) GetHistory(ctx context.Context, symbol string, limit int) ([]*domain.Quote, error) {
	return nil, nil
}

// MockLadderRepository mocks the ladder management.
type MockLadderRepository struct {
	ActiveLadderID int64
	Tickers        []*domain.TickerInfo
}

func (m *MockLadderRepository) GetActiveLadder(ctx context.Context) (int64, error) {
	return m.ActiveLadderID, nil
}

func (m *MockLadderRepository) GetLadder(ctx context.Context, id int64) (*domain.Ladder, error) {
	return &domain.Ladder{ID: id, IsActive: true}, nil
}

func (m *MockLadderRepository) GetAllowedTickers(ctx context.Context, ladderID int64) ([]*domain.TickerInfo, error) {
	return m.Tickers, nil
}

func (m *MockLadderRepository) JoinLadder(ctx context.Context, ladderID int64, userID int64) error {
	return nil
}

func (m *MockLadderRepository) IsUserInLadder(ctx context.Context, ladderID int64, userID int64) (bool, error) {
	return true, nil
}

func (m *MockLadderRepository) GetExpiredActiveLadders(ctx context.Context, now time.Time) ([]*domain.Ladder, error) {
	return nil, nil
}

func (m *MockLadderRepository) GetPendingLaddersToActivate(ctx context.Context, now time.Time) ([]*domain.Ladder, error) {
	return nil, nil
}

func (m *MockLadderRepository) UpdateLadderStatus(ctx context.Context, id int64, isActive bool) error {
	return nil
}

func (m *MockLadderRepository) GetLadderParticipants(ctx context.Context, ladderID int64) ([]domain.LadderParticipant, error) {
	return nil, nil
}

func (m *MockLadderRepository) InsertLadderParticipant(ctx context.Context, ladderID int64, userID int64, finalBalance decimal.Decimal, finalRank int32) error {
	return nil
}

func (m *MockLadderRepository) PruneLadderParticipants(ctx context.Context, ladderID int64, rankThreshold int32) error {
	return nil
}

func (m *MockLadderRepository) DeleteLadderPortfolioItemsByLadder(ctx context.Context, ladderID int64) error {
	return nil
}

func TestMarketFetcher(t *testing.T) {
	// 1. Setup Miniredis
	mr, err := miniredis.Run()
	assert.NoError(t, err)

	defer mr.Close()

	rdb := go_redis.NewClient(&go_redis.Options{
		Addr: mr.Addr(),
	})

	// 2. Setup Mock Finnhub
	mockQuote := finnhub.Response{
		Change:       10.0,
		CurrentPrice: 150.0,
		Timestamp:    time.Now().Unix(),
	}
	mockClient := &MockFinnhubClient{FinnhubQuote: mockQuote}

	// 3. Setup Worker
	marketRepo := redis.NewMarketRepository(rdb)
	historyRepo := &MockHistoryRepository{}
	ladderRepo := &MockLadderRepository{
		ActiveLadderID: 1,
		Tickers: []*domain.TickerInfo{
			{Symbol: "AAPL", Source: "Finnhub"},
		},
	}
	marketWorker := worker.NewMarketFetcher("Finnhub", mockClient, marketRepo, historyRepo, ladderRepo, &worker.FetcherConfig{
		FetchInterval:   100 * time.Millisecond,
		RefreshInterval: 1 * time.Minute,
		RequestTimeout:  5 * time.Second,
	})

	// 4. Run Worker logic (simulate one tick)
	ctx := t.Context()

	// Start worker
	go func() {
		if errWork := marketWorker.Start(ctx); errWork != nil && errWork != context.Canceled {
			t.Errorf("Worker failed: %v", errWork)
		}
	}()

	// Wait for Redis update
	time.Sleep(200 * time.Millisecond)

	// 5. Verify Redis
	val, err := rdb.Get(ctx, "market:AAPL").Result()
	assert.NoError(t, err)

	var quote exchange.Quote

	err = json.Unmarshal([]byte(val), &quote)
	assert.NoError(t, err)

	assert.Equal(t, 150.0, quote.GetPrice())
	assert.Equal(t, "AAPL", quote.GetSymbol())
}
