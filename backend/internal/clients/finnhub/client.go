package finnhub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tmythicator/ticker-rush/server/model"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string, timeout time.Duration) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *Client) GetQuote(ctx context.Context, symbol string) (*model.Quote, error) {
	url := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", symbol, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API status: %d", resp.StatusCode)
	}

	var fq model.FinnhubQuote
	if err := json.NewDecoder(resp.Body).Decode(&fq); err != nil {
		return nil, fmt.Errorf("json error: %w", err)
	}

	if fq.CurrentPrice == 0 {
		return nil, fmt.Errorf("zero price received")
	}
	ts := fq.Timestamp
	if ts == 0 {
		ts = time.Now().Unix()
	}

	return &model.Quote{
		Symbol:    symbol,
		Price:     fq.CurrentPrice,
		Timestamp: ts,
	}, nil
}
