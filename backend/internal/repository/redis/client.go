// Package redis provides Redis client and repositories.
package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// а зачем инициализация клиента в репозитории?
// плюс этот метод вызывается из exchange а в fetcher напрямую redis.NewClient
// NewClient creates a new Redis client.
func NewClient(addr string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
