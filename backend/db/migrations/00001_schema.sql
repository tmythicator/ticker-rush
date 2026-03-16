-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name TEXT NOT NULL CHECK (char_length(first_name) > 0),
    last_name TEXT NOT NULL CHECK (char_length(last_name) > 0),
    website TEXT NOT NULL DEFAULT '',
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    is_banned BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    agb_accepted_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT users_username_length_check CHECK (char_length(username) >= 3 AND char_length(username) <= 20)
);

CREATE TABLE IF NOT EXISTS market_quotes (
    symbol TEXT NOT NULL DEFAUlT 'unknown',
    price NUMERIC NOT NULL,
    source TEXT NOT NULL DEFAULT 'unknown',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS market_quotes_symbol_created_at_idx ON market_quotes (symbol, created_at);

CREATE TABLE IF NOT EXISTS ladders (
    id bigserial PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    initial_balance NUMERIC NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS ladder_tickers (
    ladder_id BIGINT NOT NULL REFERENCES ladders(id) ON DELETE CASCADE,
    stock_symbol TEXT NOT NULL,
    PRIMARY KEY (ladder_id, stock_symbol)
);

CREATE TABLE IF NOT EXISTS ladder_participants (
    ladder_id BIGINT NOT NULL REFERENCES ladders(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    final_balance NUMERIC,
    final_rank INT,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (ladder_id, user_id)
);

CREATE TABLE IF NOT EXISTS ladder_balances (
    ladder_id BIGINT NOT NULL REFERENCES ladders(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    balance NUMERIC NOT NULL,
    PRIMARY KEY (ladder_id, user_id)
);

CREATE TABLE IF NOT EXISTS ladder_portfolio_items (
    ladder_id BIGINT NOT NULL REFERENCES ladders(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stock_symbol TEXT NOT NULL,
    quantity NUMERIC NOT NULL DEFAULT 0,
    average_price NUMERIC NOT NULL DEFAULT 0,
    PRIMARY KEY (ladder_id, user_id, stock_symbol)
);

-- +goose Down
DROP TABLE IF EXISTS ladder_portfolio_items;
DROP TABLE IF EXISTS ladder_balances;
DROP TABLE IF EXISTS ladder_participants;
DROP TABLE IF EXISTS ladder_tickers;
DROP TABLE IF EXISTS ladders;
DROP TABLE IF EXISTS market_quotes;
DROP TABLE IF EXISTS portfolio_items;
DROP TABLE IF EXISTS users;
