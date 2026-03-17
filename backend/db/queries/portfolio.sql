-- name: GetPortfolio :many
SELECT * FROM ladder_portfolio_items
WHERE ladder_id = $1 AND user_id = $2;

-- name: SetPortfolioItem :exec
INSERT INTO ladder_portfolio_items (ladder_id, user_id, stock_symbol, quantity, average_price)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (ladder_id, user_id, stock_symbol)
DO UPDATE SET
   quantity = EXCLUDED.quantity,
   average_price = EXCLUDED.average_price;

-- name: DeleteUserPortfolio :exec
DELETE FROM ladder_portfolio_items
WHERE ladder_id = $1 AND user_id = $2;

-- name: DeletePortfolioItem :exec
DELETE FROM ladder_portfolio_items
WHERE ladder_id = $1 AND user_id = $2 AND stock_symbol = $3;

-- name: GetPortfolioItem :one
SELECT ladder_id, user_id, stock_symbol, quantity, average_price
FROM ladder_portfolio_items
WHERE ladder_id = $1 AND user_id = $2 AND stock_symbol = $3;

-- name: GetPortfolioItemForUpdate :one
SELECT ladder_id, user_id, stock_symbol, quantity, average_price
FROM ladder_portfolio_items
WHERE ladder_id = $1 AND user_id = $2 AND stock_symbol = $3 FOR UPDATE;
