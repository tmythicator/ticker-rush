package service

import (
	"context"

	"github.com/tmythicator/ticker-rush/server/internal/config"
)

// ConfigService handles configuration related operations.
type ConfigService struct {
	cfg *config.Config
}

// NewConfigService creates a new instance of ConfigService.
func NewConfigService(cfg *config.Config) *ConfigService {
	return &ConfigService{
		cfg: cfg,
	}
}

// GetPublicConfig returns the public configuration that can be exposed to the frontend.
func (s *ConfigService) GetPublicConfig(ctx context.Context) map[string]any {
	return map[string]any{
		"tickers": s.cfg.Tickers,
	}
}
