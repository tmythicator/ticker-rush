// Package redis provides Redis repository implementations.
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var rateLimitLuaScript = redis.NewScript(`
	local current = redis.call('INCR', KEYS[1])
	if current == 1 then
		redis.call('EXPIRE', KEYS[1], ARGV[1])
	end
	return current
`)

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
	val, err := rateLimitLuaScript.Run(ctx, rl.valkey, []string{key}, int(window.Seconds())).Int64()
	if err != nil {
		return 0, err
	}

	return val, nil
}
