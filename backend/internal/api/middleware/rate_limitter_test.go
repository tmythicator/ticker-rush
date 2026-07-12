package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
)

// mockRateLimitStore implements middleware.RateLimit for testing.
type mockRateLimitStore struct {
	incrFunc func(ctx context.Context, key string, window time.Duration) (int64, error)
}

func (m *mockRateLimitStore) Increment(ctx context.Context, key string, window time.Duration) (int64, error) {
	if m.incrFunc != nil {
		return m.incrFunc(ctx, key, window)
	}

	return 1, nil
}

func setupTestGin() (*gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	w := httptest.NewRecorder()

	return r, w
}

func TestRateLimitter_Allow(t *testing.T) {
	r, w := setupTestGin()

	mockLimit := 5
	mockWindow := time.Minute

	store := &mockRateLimitStore{
		incrFunc: func(ctx context.Context, key string, window time.Duration) (int64, error) {
			return 3, nil
		},
	}

	limiter := NewRateLimitter(store, mockLimit, mockWindow)

	r.Use(limiter.Limit())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "success" {
		t.Errorf("expected body 'success', got %q", w.Body.String())
	}

	if limit := w.Header().Get("X-RateLimit-Limit"); limit != strconv.Itoa(mockLimit) {
		t.Errorf("expected X-RateLimit-Limit to be %d, got %s", mockLimit, limit)
	}

	if remaining := w.Header().Get("X-RateLimit-Remaining"); remaining != strconv.Itoa(mockLimit-3) {
		t.Errorf("expected X-RateLimit-Remaining to be %d, got %s", mockLimit-3, remaining)
	}

	if reset := w.Header().Get("X-RateLimit-Reset"); reset == "" {
		t.Error("expected X-RateLimit-Reset header to be set, but it was empty")
	}
}

func TestRateLimitter_LimitExceeded(t *testing.T) {
	r, w := setupTestGin()

	mockLimit := 5
	mockWindow := time.Minute

	store := &mockRateLimitStore{
		incrFunc: func(ctx context.Context, key string, window time.Duration) (int64, error) {
			return 6, nil
		},
	}

	limiter := NewRateLimitter(store, mockLimit, mockWindow)

	r.Use(limiter.Limit())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	// Verify HTTP status code is 429 Too Many Requests
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected status 429, got %d", w.Code)
	}

	// Verify headers
	if limit := w.Header().Get("X-RateLimit-Limit"); limit != strconv.Itoa(mockLimit) {
		t.Errorf("expected X-RateLimit-Limit to be %d, got %s", mockLimit, limit)
	}

	if remaining := w.Header().Get("X-RateLimit-Remaining"); remaining != "0" {
		t.Errorf("expected X-RateLimit-Remaining to be '0', got %s", remaining)
	}

	if retryAfter := w.Header().Get("Retry-After"); retryAfter == "" {
		t.Error("expected Retry-After header to be set, but it was empty")
	}

	// Verify RFC 7807 problem details response body
	var problem apperrors.ProblemDetails
	err := json.Unmarshal(w.Body.Bytes(), &problem)
	if err != nil {
		t.Fatalf("failed to unmarshal problem details: %v", err)
	}

	if problem.Type != apperrors.TypeRateLimitExceeded {
		t.Errorf("expected type %s, got %s", apperrors.TypeRateLimitExceeded, problem.Type)
	}
}

func TestRateLimitter_IPKey_Anonymous(t *testing.T) {
	r, w := setupTestGin()

	var capturedKey string
	store := &mockRateLimitStore{
		incrFunc: func(ctx context.Context, key string, window time.Duration) (int64, error) {
			capturedKey = key

			return 1, nil
		},
	}

	limiter := NewRateLimitter(store, 5, time.Minute)

	r.Use(limiter.Limit())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.0.2.1:12345"
	r.ServeHTTP(w, req)

	// Verify that the captured key starts with the ip prefix and contains the IP
	if !strings.Contains(capturedKey, "rate_limit:ip:192.0.2.1:") {
		t.Errorf("expected key to contain 'rate_limit:ip:192.0.2.1:', got %q", capturedKey)
	}
}

func TestRateLimitter_UserKey_Authenticated(t *testing.T) {
	r, w := setupTestGin()

	var capturedKey string
	store := &mockRateLimitStore{
		incrFunc: func(ctx context.Context, key string, window time.Duration) (int64, error) {
			capturedKey = key

			return 1, nil
		},
	}

	limiter := NewRateLimitter(store, 5, time.Minute)

	// Forcefully set UserID in context before invoking the rate limiter
	r.Use(func(c *gin.Context) {
		c.Set(UserIDKey, int64(42))
		c.Next()
	})
	r.Use(limiter.Limit())

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	// Verify that the captured key starts with the user prefix and contains the user ID
	if !strings.Contains(capturedKey, "rate_limit:user:42:") {
		t.Errorf("expected key to contain 'rate_limit:user:42:', got %q", capturedKey)
	}
}

func TestRateLimitter_FailOpen(t *testing.T) {
	r, w := setupTestGin()

	store := &mockRateLimitStore{
		incrFunc: func(ctx context.Context, key string, window time.Duration) (int64, error) {
			return 0, errors.New("redis down")
		},
	}

	limiter := NewRateLimitter(store, 5, time.Minute)

	r.Use(limiter.Limit())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "success" {
		t.Errorf("expected body 'success', got %q", w.Body.String())
	}
}
