package storage

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/model"
)

type MarketRepository struct {
	rdb *redis.Client
}

func NewMarketRepository(rdb *redis.Client) *MarketRepository {
	return &MarketRepository{rdb: rdb}
}

func (r *MarketRepository) GetQuote(ctx context.Context, symbol string) (*model.Quote, error) {
	val, err := r.rdb.Get(ctx, "market:"+symbol).Result()
	if err != nil {
		return nil, err
	}

	var quote model.Quote
	if err := json.Unmarshal([]byte(val), &quote); err != nil {
		return nil, err
	}

	return &quote, nil
}
