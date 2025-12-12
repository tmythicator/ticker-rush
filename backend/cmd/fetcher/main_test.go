package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	go_redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/tmythicator/ticker-rush/server/internal/clients/finnhub"
	"github.com/tmythicator/ticker-rush/server/internal/model"
	"github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/internal/worker"
)

// MockFinnhubClient mocks the Finnhub API
type MockFinnhubClient struct {
	Quote finnhub.FinnhubQuote
}

func (m *MockFinnhubClient) GetQuote(ctx context.Context, symbol string) (*model.Quote, error) {
	return &model.Quote{
		Symbol:    symbol,
		Price:     m.Quote.CurrentPrice,
		Timestamp: m.Quote.Timestamp,
	}, nil
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
	mockQuote := finnhub.FinnhubQuote{
		Change:       10.0,
		CurrentPrice: 150.0,
		Timestamp:    time.Now().Unix(),
	}
	mockClient := &MockFinnhubClient{Quote: mockQuote}

	// 3. Setup Worker
	marketRepo := redis.NewMarketRepository(rdb)
	marketWorker := worker.NewMarketFetcher(mockClient, marketRepo)

	// 4. Run Worker logic (simulate one tick)
	// Since Start runs in a loop, we can test the internal logic by calling processTicker directly if it was public,
	// or we can run Start in a goroutine and wait.
	// However, processTicker is private. Let's make it public or just run Start and wait a bit.
	// Actually, for testing, it's better if we can trigger it.
	// Let's rely on the fact that Start calls processTicker immediately once.

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start in a goroutine
	go marketWorker.Start(ctx, "AAPL", 100*time.Millisecond)

	// Wait for Redis update
	time.Sleep(200 * time.Millisecond)

	// 5. Verify Redis
	val, err := rdb.Get(ctx, "market:AAPL").Result()
	assert.NoError(t, err)

	var quote model.Quote
	err = json.Unmarshal([]byte(val), &quote)
	assert.NoError(t, err)

	assert.Equal(t, 150.0, quote.Price)
	assert.Equal(t, "AAPL", quote.Symbol)
}
