-- +goose Up
INSERT INTO users (username, password_hash, first_name, last_name, website, created_at, agb_accepted_at, is_admin)
VALUES
    ('cool_trader', '$2a$10$iehMYy8l.nvHJ.CxRqmtJOhfRJHw1kpnQfHDmJX7Qorq1QVeI.1lK', 'Cool', 'Trader', 'https://example.com', NOW(), NOW(), FALSE),  -- PW: test123
    ('rich_guy', '$2a$10$iehMYy8l.nvHJ.CxRqmtJOhfRJHw1kpnQfHDmJX7Qorq1QVeI.1lK', 'Rich', 'Guy', '', NOW(), NOW(), FALSE),  -- PW: test123
    ('admin', '$2a$10$iehMYy8l.nvHJ.CxRqmtJOhfRJHw1kpnQfHDmJX7Qorq1QVeI.1lK', 'Admin', 'Admin', '', NOW(), NOW(), TRUE);  -- PW: test123

INSERT INTO ladders (name, type, start_time, end_time, initial_balance, is_active)
VALUES
    ('Test game', 'monthly', NOW(), NOW() + INTERVAL '30 days', 10000.0, TRUE);

INSERT INTO ladder_tickers (ladder_id, stock_symbol, source)
SELECT id, 'AAPL', 'Finnhub' FROM ladders WHERE name = 'Test game'
UNION ALL
SELECT id, 'AMZN', 'Finnhub' FROM ladders WHERE name = 'Test game'
UNION ALL
SELECT id, 'MSFT', 'Finnhub' FROM ladders WHERE name = 'Test game'
UNION ALL
SELECT id, 'bitcoin', 'CoinGecko' FROM ladders WHERE name = 'Test game'
UNION ALL
SELECT id, 'ethereum', 'CoinGecko' FROM ladders WHERE name = 'Test game';

-- +goose Down
DELETE FROM ladder_tickers WHERE ladder_id IN (SELECT id FROM ladders WHERE name = 'Test game');
DELETE FROM ladders WHERE name = 'Test game';
DELETE FROM users WHERE username IN ('cool_trader', 'rich_guy', 'admin');
