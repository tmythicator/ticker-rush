-- +goose Up
CREATE TABLE market_quotes (
    id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR(50) NOT NULL,
    price DECIMAL(20, 8) NOT NULL,
    source VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_market_quotes_symbol_created_at ON market_quotes(symbol, created_at DESC);

-- +goose Down
DROP TABLE market_quotes;
