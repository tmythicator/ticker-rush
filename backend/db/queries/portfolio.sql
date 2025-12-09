-- name: GetPortfolio :many
SELECT * FROM portfolio_items
WHERE user_id = $1;

-- name: SetPortfolioItem :exec
INSERT INTO portfolio_items (user_id, stock_symbol, quantity, average_price)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, stock_symbol)
DO UPDATE SET
    quantity = EXCLUDED.quantity,
    average_price = EXCLUDED.average_price;

-- name: DeleteUserPortfolio :exec
DELETE FROM portfolio_items
WHERE user_id = $1;

-- name: DeletePortfolioItem :exec
DELETE FROM portfolio_items
WHERE user_id = $1 AND stock_symbol = $2;

-- name: GetPortfolioItem :one
SELECT user_id, stock_symbol, quantity, average_price 
FROM portfolio_items
WHERE user_id = $1 AND stock_symbol = $2;

-- name: GetPortfolioItemForUpdate :one
SELECT user_id, stock_symbol, quantity, average_price 
FROM portfolio_items
WHERE user_id = $1 AND stock_symbol = $2 FOR UPDATE;
