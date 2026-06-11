// Package redis provides Valkey/Redis repositories.
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"

	"github.com/tmythicator/ticker-rush/backend/internal/domain"
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

// valkeyQuote is used for backward-compatible JSON serialization in Redis.
type valkeyQuote struct {
	Symbol        string  `json:"symbol,omitempty"`
	Price         float64 `json:"price,omitempty"`
	Change        float64 `json:"change,omitempty"`
	ChangePercent float64 `json:"change_percent,omitempty"`
	Timestamp     int64   `json:"timestamp,omitempty"`
	Source        string  `json:"source,omitempty"`
	IsClosed      bool    `json:"is_closed,omitempty"`
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
func (r *MarketRepository) GetQuote(ctx context.Context, symbol string) (*domain.Quote, error) {
	key := marketQuoteKey(symbol)
	val, err := r.valkey.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var vq valkeyQuote
	if err := json.Unmarshal([]byte(val), &vq); err != nil {
		return nil, err
	}

	// Check if market is closed (data older than 30 minutes)
	isClosed := vq.IsClosed
	if time.Since(time.Unix(vq.Timestamp, 0)) > 30*time.Minute {
		isClosed = true
	}

	return &domain.Quote{
		Symbol:        vq.Symbol,
		Price:         decimal.NewFromFloat(vq.Price),
		Change:        decimal.NewFromFloat(vq.Change),
		ChangePercent: decimal.NewFromFloat(vq.ChangePercent),
		Timestamp:     time.Unix(vq.Timestamp, 0),
		Source:        vq.Source,
		IsClosed:      isClosed,
	}, nil
}

// SaveQuote saves a quote to Redis and publishes it to the channel.
func (r *MarketRepository) SaveQuote(ctx context.Context, quote *domain.Quote) error {
	vq := valkeyQuote{
		Symbol:        quote.Symbol,
		Price:         quote.Price.InexactFloat64(),
		Change:        quote.Change.InexactFloat64(),
		ChangePercent: quote.ChangePercent.InexactFloat64(),
		Timestamp:     quote.Timestamp.Unix(),
		Source:        quote.Source,
		IsClosed:      quote.IsClosed,
	}

	data, err := json.Marshal(vq)
	if err != nil {
		return err
	}

	key := marketQuoteKey(quote.Symbol)
	channel := marketQuoteChannel(quote.Symbol)

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
