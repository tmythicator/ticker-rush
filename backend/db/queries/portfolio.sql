-- name: GetPortfolio :many
SELECT user_id, stock_symbol, quantity, average_price
FROM portfolio_items
WHERE user_id = $1;

-- name: AddPortfolioItem :exec
INSERT INTO portfolio_items (user_id, stock_symbol, quantity, average_price)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, stock_symbol) DO UPDATE SET
    quantity = portfolio_items.quantity + EXCLUDED.quantity,
    average_price = (portfolio_items.average_price * portfolio_items.quantity + EXCLUDED.average_price * EXCLUDED.quantity) / (portfolio_items.quantity + EXCLUDED.quantity);

-- name: SetPortfolioItem :exec
INSERT INTO portfolio_items (user_id, stock_symbol, quantity, average_price)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, stock_symbol) DO UPDATE SET
    quantity = EXCLUDED.quantity,
    average_price = EXCLUDED.average_price;

-- name: RemovePortfolioItem :exec
DELETE FROM portfolio_items
WHERE user_id = $1 AND stock_symbol = $2;
