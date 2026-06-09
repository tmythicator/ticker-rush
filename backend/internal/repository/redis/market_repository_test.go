package redis_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/tmythicator/ticker-rush/backend/internal/proto/exchange/v1"
	redisRepo "github.com/tmythicator/ticker-rush/backend/internal/repository/redis"
)

func TestMarketRepository(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	rClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer rClient.Close()

	repo := redisRepo.NewMarketRepository(rClient)
	ctx := context.Background()

	const symbol = "AAPL"

	t.Run("Save and Get Active Quote", func(t *testing.T) {
		now := time.Now().Unix()
		quote := &exchange.Quote{
			Symbol:    symbol,
			Price:     150.25,
			Timestamp: now,
			IsClosed:  false,
		}

		err := repo.SaveQuote(ctx, quote)
		assert.NoError(t, err)

		fetched, err := repo.GetQuote(ctx, symbol)
		assert.NoError(t, err)
		assert.Equal(t, symbol, fetched.Symbol)
		assert.Equal(t, 150.25, fetched.Price)
		assert.Equal(t, now, fetched.Timestamp)
		assert.False(t, fetched.IsClosed)
	})

	t.Run("Get Quote Older Than 30 Minutes Sets IsClosed", func(t *testing.T) {
		oldTime := time.Now().Add(-35 * time.Minute).Unix()
		quote := &exchange.Quote{
			Symbol:    "MSFT",
			Price:     320.50,
			Timestamp: oldTime,
			IsClosed:  false,
		}

		err := repo.SaveQuote(ctx, quote)
		assert.NoError(t, err)

		fetched, err := repo.GetQuote(ctx, "MSFT")
		assert.NoError(t, err)
		assert.True(t, fetched.IsClosed, "Market should be marked closed for old data")
	})

	t.Run("Subscribe to Quotes", func(t *testing.T) {
		pubSub := repo.SubscribeToQuotes(ctx, symbol)
		defer pubSub.Close()

		// Wait briefly for subscription to register in miniredis
		time.Sleep(10 * time.Millisecond)

		quote := &exchange.Quote{
			Symbol:    symbol,
			Price:     152.00,
			Timestamp: time.Now().Unix(),
		}

		err := repo.SaveQuote(ctx, quote)
		assert.NoError(t, err)

		msg, err := pubSub.ReceiveMessage(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, msg.Payload)

		var received exchange.Quote
		err = json.Unmarshal([]byte(msg.Payload), &received)
		assert.NoError(t, err)
		assert.Equal(t, 152.00, received.Price)
	})
}
