// Package redis provides Redis repository implementations.
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimitter implements a rate limiting repository using Redis.
type RateLimitter struct {
	valkey *redis.Client
}

// NewRateLimitter creates an instance of RateLimitter repository.
func NewRateLimitter(valkey *redis.Client) *RateLimitter {
	return &RateLimitter{valkey: valkey}
}

// Increment increases the count for a key and sets an expiration time.
func (rl *RateLimitter) Increment(ctx context.Context, key string, window time.Duration) (int64, error) {
	pipe := rl.valkey.Pipeline()

	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return incr.Val(), nil
}
