// Package finnhub provides a client for the Finnhub API.
package finnhub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tmythicator/ticker-rush/server/internal/proto/exchange/v1"
)

// Response represents a stock quote from the Finnhub API.
type Response struct {
	CurrentPrice  float64 `json:"c"`  // c = Current price
	Change        float64 `json:"d"`  // d = Change
	PercentChange float64 `json:"dp"` // dp = Percent change
	Timestamp     int64   `json:"t"`  // t = Timestamp
}

// Client is a client for the Finnhub API.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new instance of FinnhubClient.
func NewClient(apiKey string, timeout time.Duration) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    "https://finnhub.io/api/v1",
		httpClient: &http.Client{Timeout: timeout},
	}
}

// GetQuote fetches a stock quote for a given symbol.
func (c *Client) GetQuote(ctx context.Context, symbol string) (*exchange.Quote, error) {
	apiSymbol := strings.TrimPrefix(symbol, "FH:")

	url := fmt.Sprintf("%s/quote?symbol=%s&token=%s", c.baseURL, apiSymbol, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API status: %d", resp.StatusCode)
	}

	var fq Response
	if err := json.NewDecoder(resp.Body).Decode(&fq); err != nil {
		return nil, fmt.Errorf("json error: %w", err)
	}

	if fq.CurrentPrice == 0 {
		return nil, errors.New("zero price received")
	}

	ts := fq.Timestamp
	if ts == 0 {
		ts = time.Now().Unix()
	}

	return &exchange.Quote{
		Symbol:        symbol,
		Price:         fq.CurrentPrice,
		Change:        fq.Change,
		ChangePercent: fq.PercentChange,
		Timestamp:     ts,
	}, nil
}
