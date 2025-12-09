-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    balance NUMERIC NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS portfolio_items (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stock_symbol TEXT NOT NULL,
    quantity NUMERIC NOT NULL DEFAULT 0,
    average_price NUMERIC NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, stock_symbol)
);

-- +goose Down
DROP TABLE IF EXISTS portfolio_items;
DROP TABLE IF EXISTS users;
