-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name TEXT NOT NULL CHECK (char_length(first_name) > 0),
    last_name TEXT NOT NULL CHECK (char_length(last_name) > 0),
    website TEXT NOT NULL DEFAULT '',
    balance NUMERIC NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT users_username_length_check CHECK (char_length(username) >= 3 AND char_length(username) <= 20)
);

CREATE TABLE IF NOT EXISTS portfolio_items (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stock_symbol TEXT NOT NULL,
    quantity NUMERIC NOT NULL DEFAULT 0,
    average_price NUMERIC NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, stock_symbol)
);

CREATE TABLE IF NOT EXISTS market_quotes (
    symbol TEXT NOT NULL DEFAUlT 'unknown',
    price NUMERIC NOT NULL,
    source TEXT NOT NULL DEFAULT 'unknown',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS market_quotes_symbol_created_at_idx ON market_quotes (symbol, created_at);

-- +goose Down
DROP TABLE IF EXISTS market_quotes;
DROP TABLE IF EXISTS portfolio_items;
DROP TABLE IF EXISTS users;
