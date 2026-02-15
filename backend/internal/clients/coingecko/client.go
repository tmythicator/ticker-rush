// Package coingecko provides a client for the CoinGecko API.
package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tmythicator/ticker-rush/server/internal/proto/exchange/v1"
)

// Client is a client for the CoinGecko API.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new CoinGecko client.
func NewClient(apiKey string, timeout time.Duration) *Client {
	return &Client{
		baseURL: "https://api.coingecko.com/api/v3",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Response represents the JSON response from CoinGecko.
// Since the key is dynamic (the coin ID), we use a map, but the value is structured.
type Response map[string]struct {
	USD float64 `json:"usd"`
}

// GetQuote fetches the price of a crypto asset.
func (c *Client) GetQuote(ctx context.Context, symbol string) (*exchange.Quote, error) {
	id := strings.TrimPrefix(symbol, "CG:")

	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd", c.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "TickerRush/1.0")
	if c.apiKey != "" {
		req.Header.Set("x-cg-demo-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("rate limited (429)")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	data, ok := result[id]
	if !ok {
		return nil, fmt.Errorf("symbol not found or no USD price: %s", id)
	}

	timestamp := time.Now().Unix()

	return &exchange.Quote{
		Symbol:    symbol,
		Price:     data.USD,
		Timestamp: timestamp,
	}, nil
}
