-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    balance DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS stocks (
    symbol TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS portfolio_items (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stock_symbol TEXT NOT NULL REFERENCES stocks(symbol) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 0,
    average_price DOUBLE PRECISION NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, stock_symbol)
);

-- +goose Down
DROP TABLE IF EXISTS portfolio_items;
DROP TABLE IF EXISTS stocks;
DROP TABLE IF EXISTS users;
