-- +goose Up
INSERT INTO ladders (name, type, start_time, end_time, initial_balance, is_active)
SELECT 'First Challenge', 'monthly', NOW(), NOW() + INTERVAL '30 days', 10000.0, TRUE
WHERE NOT EXISTS (SELECT 1 FROM ladders WHERE name = 'First Challenge');

INSERT INTO ladder_tickers (ladder_id, stock_symbol, source)
SELECT id, 'AAPL', 'Finnhub' FROM ladders WHERE name = 'First Challenge'
ON CONFLICT DO NOTHING;

INSERT INTO ladder_tickers (ladder_id, stock_symbol, source)
SELECT id, 'AMZN', 'Finnhub' FROM ladders WHERE name = 'First Challenge'
ON CONFLICT DO NOTHING;

INSERT INTO ladder_tickers (ladder_id, stock_symbol, source)
SELECT id, 'MSFT', 'Finnhub' FROM ladders WHERE name = 'First Challenge'
ON CONFLICT DO NOTHING;

INSERT INTO ladder_tickers (ladder_id, stock_symbol, source)
SELECT id, 'bitcoin', 'CoinGecko' FROM ladders WHERE name = 'First Challenge'
ON CONFLICT DO NOTHING;

INSERT INTO ladder_tickers (ladder_id, stock_symbol, source)
SELECT id, 'ethereum', 'CoinGecko' FROM ladders WHERE name = 'First Challenge'
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM ladder_tickers WHERE ladder_id IN (SELECT id FROM ladders WHERE name = 'First Challenge');
DELETE FROM ladders WHERE name = 'First Challenge';