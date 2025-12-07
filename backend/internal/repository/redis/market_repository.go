package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/model"
)

type MarketRepository struct {
	valkey *redis.Client
}

func NewMarketRepository(valkey *redis.Client) *MarketRepository {
	return &MarketRepository{valkey: valkey}
}

func (r *MarketRepository) GetQuote(ctx context.Context, symbol string) (*model.Quote, error) {
	val, err := r.valkey.Get(ctx, "market:"+symbol).Result()
	if err != nil {
		return nil, err
	}

	var quote model.Quote
	if err := json.Unmarshal([]byte(val), &quote); err != nil {
		return nil, err
	}

	return &quote, nil
}

func (r *MarketRepository) SaveQuote(ctx context.Context, quote *model.Quote) error {
	data, err := json.Marshal(quote)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("market:%s", quote.Symbol)
	return r.valkey.Set(ctx, key, data, 0).Err()
}
