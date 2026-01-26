package redis

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/internal/proto/exchange"
)

type MarketRepository struct {
	valkey *redis.Client
}

func NewMarketRepository(valkey *redis.Client) *MarketRepository {
	return &MarketRepository{valkey: valkey}
}

func (r *MarketRepository) GetQuote(ctx context.Context, symbol string) (*exchange.Quote, error) {
	val, err := r.valkey.Get(ctx, "market:"+symbol).Result()
	if err != nil {
		return nil, err
	}

	var quote exchange.Quote
	if err := json.Unmarshal([]byte(val), &quote); err != nil {
		return nil, err
	}

	return &quote, nil
}

func (r *MarketRepository) SaveQuote(ctx context.Context, quote *exchange.Quote) error {
	data, err := json.Marshal(quote)
	if err != nil {
		return err
	}

	key := "market:" + quote.GetSymbol()

	// Publish to specific channel for the symbol
	channel := "market:quote:" + quote.GetSymbol()

	pipe := r.valkey.Pipeline()
	pipe.Set(ctx, key, data, 0)
	pipe.Publish(ctx, channel, data)
	_, err = pipe.Exec(ctx)

	return err
}

func (r *MarketRepository) SubscribeToQuotes(ctx context.Context, symbol string) *redis.PubSub {
	channel := "market:quote:" + symbol

	return r.valkey.Subscribe(ctx, channel)
}
