package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/tmythicator/ticker-rush/backend/internal/proto/exchange/v1"
)

const (
	marketQuotePrefix   = "market"
	marketChannelPrefix = "market:quote"
)

func marketQuoteKey(symbol string) string {
	return fmt.Sprintf("%s:%s", marketQuotePrefix, symbol)
}

func marketQuoteChannel(symbol string) string {
	return fmt.Sprintf("%s:%s", marketChannelPrefix, symbol)
}

// MarketRepository handles market data storage in Redis.
type MarketRepository struct {
	valkey *redis.Client
}

// NewMarketRepository creates a new instance of MarketRepository.
func NewMarketRepository(valkey *redis.Client) *MarketRepository {
	return &MarketRepository{valkey: valkey}
}

// GetQuote retrieves the latest quote for a symbol from Redis.
func (r *MarketRepository) GetQuote(ctx context.Context, symbol string) (*exchange.Quote, error) {
	key := marketQuoteKey(symbol)
	val, err := r.valkey.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var quote exchange.Quote
	if err := json.Unmarshal([]byte(val), &quote); err != nil {
		return nil, err
	}

	// Check if market is closed (data older than 30 minutes)
	if time.Since(time.Unix(quote.Timestamp, 0)) > 30*time.Minute {
		quote.IsClosed = true
	}

	return &quote, nil
}

// SaveQuote saves a quote to Redis and publishes it to the channel.
func (r *MarketRepository) SaveQuote(ctx context.Context, quote *exchange.Quote) error {
	data, err := json.Marshal(quote)
	if err != nil {
		return err
	}

	key := marketQuoteKey(quote.GetSymbol())
	channel := marketQuoteChannel(quote.GetSymbol())

	pipe := r.valkey.Pipeline()
	pipe.Set(ctx, key, data, 0)
	pipe.Publish(ctx, channel, data)
	_, err = pipe.Exec(ctx)

	return err
}

// SubscribeToQuotes subscribes to real-time quote updates for a symbol.
func (r *MarketRepository) SubscribeToQuotes(ctx context.Context, symbol string) *redis.PubSub {
	channel := marketQuoteChannel(symbol)

	return r.valkey.Subscribe(ctx, channel)
}
