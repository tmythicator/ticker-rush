package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Clear envs - careful not to disrupt actual env but here we test cleanly
	os.Clearenv()

	cfg, err := LoadConfig()
	assert.NoError(t, err)

	assert.Equal(t, []string{"AAPL", "BINANCE:BTCUSDT"}, cfg.Tickers)
	assert.Equal(t, 8081, cfg.ServerPort)
	assert.Equal(t, "localhost", cfg.RedisHost)
	assert.Equal(t, 6379, cfg.RedisPort)
	assert.Equal(t, 5173, cfg.ClientPort)
	assert.Equal(t, 10*time.Second, cfg.FinnhubFetchInterval)
	assert.Equal(t, 10*time.Second, cfg.FinnhubTimeout)
	assert.Equal(t, "postgres", cfg.PostgresUser)
	assert.Equal(t, "postgres", cfg.PostgresPass)
	assert.Equal(t, "ticker_rush", cfg.PostgresDB)
	assert.Equal(t, 5432, cfg.PostgresPort)
	assert.Equal(t, "secret", cfg.JWTSecret)
	// FinnhubKey is required but has no default, so it will be empty here
	assert.Empty(t, cfg.FinnhubKey)
}

func TestLoadConfig_Overrides(t *testing.T) {
	os.Clearenv()
	require.NoError(t, os.Setenv("TICKERS", "GOOG,MSFT"))
	require.NoError(t, os.Setenv("SERVER_PORT", "9090"))
	require.NoError(t, os.Setenv("REDIS_HOST", "redis-prod"))
	require.NoError(t, os.Setenv("REDIS_PORT", "6380"))
	require.NoError(t, os.Setenv("CLIENT_PORT", "3000"))
	require.NoError(t, os.Setenv("FINNHUB_FETCH_INTERVAL", "500ms")) // Test duration parsing
	require.NoError(t, os.Setenv("FINNHUB_API_KEY", "secret_key"))
	require.NoError(t, os.Setenv("FINNHUB_TIMEOUT", "5s"))
	require.NoError(t, os.Setenv("POSTGRES_USER", "admin"))
	require.NoError(t, os.Setenv("POSTGRES_PASSWORD", "secure"))
	require.NoError(t, os.Setenv("POSTGRES_DB", "prod_db"))
	require.NoError(t, os.Setenv("POSTGRES_PORT", "5433"))
	require.NoError(t, os.Setenv("POSTGRES_HOST", "db-prod"))
	defer os.Clearenv() // Clean up after test

	cfg, err := LoadConfig()
	assert.NoError(t, err)

	assert.Equal(t, []string{"GOOG", "MSFT"}, cfg.Tickers)
	assert.Equal(t, 9090, cfg.ServerPort)
	assert.Equal(t, "redis-prod", cfg.RedisHost)
	assert.Equal(t, 6380, cfg.RedisPort)
	assert.Equal(t, 3000, cfg.ClientPort)
	assert.Equal(t, 500*time.Millisecond, cfg.FinnhubFetchInterval)
	assert.Equal(t, "secret_key", cfg.FinnhubKey)
	assert.Equal(t, 5*time.Second, cfg.FinnhubTimeout)
	assert.Equal(t, "admin", cfg.PostgresUser)
	assert.Equal(t, "secure", cfg.PostgresPass)
	assert.Equal(t, "prod_db", cfg.PostgresDB)
	assert.Equal(t, 5433, cfg.PostgresPort)
	assert.Equal(t, "db-prod", cfg.PostgresHost)
}

func TestValidateFinnhubKey(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		cfg := &Config{
			FinnhubKey: "some_key",
		}
		assert.NoError(t, cfg.ValidateFinnhubKey())
	})

	t.Run("MissingAPIKey", func(t *testing.T) {
		cfg := &Config{
			FinnhubKey: "",
		}
		assert.Error(t, cfg.ValidateFinnhubKey())
		assert.EqualError(t, cfg.ValidateFinnhubKey(), "FINNHUB_API_KEY is not set")
	})
}
