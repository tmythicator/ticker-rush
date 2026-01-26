// Package redis provides Redis client and repositories.
package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

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
