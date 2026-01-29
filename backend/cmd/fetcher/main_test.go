package main

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	go_redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/tmythicator/ticker-rush/server/internal/clients/finnhub"
	"github.com/tmythicator/ticker-rush/server/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/internal/worker"
)

// MockFinnhubClient mocks the Finnhub API.
type MockFinnhubClient struct {
	FinnhubQuote finnhub.APIQuote
}

func (m *MockFinnhubClient) GetQuote(ctx context.Context, symbol string) (*exchange.Quote, error) {
	return &exchange.Quote{
		Symbol:    symbol,
		Price:     m.FinnhubQuote.CurrentPrice,
		Timestamp: m.FinnhubQuote.Timestamp,
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
	mockQuote := finnhub.APIQuote{
		Change:       10.0,
		CurrentPrice: 150.0,
		Timestamp:    time.Now().Unix(),
	}
	mockClient := &MockFinnhubClient{FinnhubQuote: mockQuote}

	// 3. Setup Worker
	marketRepo := redis.NewMarketRepository(rdb)
	marketWorker := worker.NewMarketFetcher(mockClient, marketRepo)

	// 4. Run Worker logic (simulate one tick)
	ctx := t.Context()

	// Start in a goroutine
	go marketWorker.Start(ctx, "AAPL", 100*time.Millisecond, &sync.WaitGroup{})

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
