// Package middleware provides HTTP middlewares.
package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
)

// RateLimit defines the interface for rate limiting operations.
type RateLimit interface {
	Increment(ctx context.Context, key string, window time.Duration) (int64, error)
}

// RateLimitter is a middleware for rate limiting.
type RateLimitter struct {
	store  RateLimit
	limit  int
	window time.Duration
}

// NewRateLimitter creates a new RateLimitter.
func NewRateLimitter(store RateLimit, limit int, window time.Duration) *RateLimitter {
	return &RateLimitter{
		store:  store,
		limit:  limit,
		window: window,
	}
}

// Limit returns a gin.HandlerFunc that implements rate limiting.
func (rl *RateLimitter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := GetUserID(c)
		var idKey string
		if ok {
			idKey = fmt.Sprintf("user:%d", userID)
		} else {
			idKey = fmt.Sprintf("ip:%s", c.ClientIP())
		}

		timeBucket := time.Now().Unix() / int64(rl.window.Seconds())

		rlKey := fmt.Sprintf("rate_limit:%s:%d", idKey, timeBucket)

		rateVal, err := rl.store.Increment(c, rlKey, rl.window)
		if err != nil {
			log.Printf("[RateLimitter] store error: %v", err)
			c.Next()

			return
		}

		resetTime := (timeBucket + 1) * int64(rl.window.Seconds())
		rem := max(int64(rl.limit)-rateVal, 0)
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", rem))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))

		if rateVal > int64(rl.limit) {
			log.Printf("[RateLimitter] rate limit exceeded for client %s (requests: %d/%d)", idKey, rateVal, rl.limit)
			retryAfter := resetTime - time.Now().Unix()
			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			respondWithProblemDirectly(c,
				http.StatusTooManyRequests,
				apperrors.TypeRateLimitExceeded,
				"API rate limit exceeded. Please try again later.")

			c.Abort()

			return
		}

		c.Next()
	}

}
