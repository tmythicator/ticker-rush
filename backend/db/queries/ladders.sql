-- name: CreateLadder :one
INSERT INTO ladders (name, type, start_time, end_time, initial_balance, is_active)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, type, start_time, end_time, initial_balance, is_active, created_at;

-- name: GetActiveLadder :one
SELECT id, name, type, start_time, end_time, initial_balance, is_active, created_at
FROM ladders
WHERE is_active = TRUE
ORDER BY start_time DESC
LIMIT 1;

-- name: GetLadder :one
SELECT id, name, type, start_time, end_time, initial_balance, is_active, created_at
FROM ladders
WHERE id = $1 LIMIT 1;

-- name: ListLadders :many
SELECT id, name, type, start_time, end_time, initial_balance, is_active, created_at
FROM ladders
ORDER BY start_time DESC;

-- name: GetLadderTickers :many
SELECT stock_symbol, source
FROM ladder_tickers
WHERE ladder_id = $1;

-- name: AddLadderTicker :exec
INSERT INTO ladder_tickers (ladder_id, stock_symbol, source)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;

-- name: UpdateLadderStatus :exec
UPDATE ladders
SET is_active = $2
WHERE id = $1;

-- name: GetLadderLeaderboard :many
SELECT lp.ladder_id, lp.user_id, lp.final_balance, lp.final_rank, lp.joined_at, u.username
FROM ladder_participants lp
JOIN users u ON lp.user_id = u.id
WHERE lp.ladder_id = $1
ORDER BY lp.final_rank ASC
LIMIT $2;

-- name: GetLadderParticipants :many
SELECT ladder_id, user_id, final_balance, final_rank, joined_at
FROM ladder_participants
WHERE ladder_id = $1;

-- name: InsertLadderParticipant :exec
INSERT INTO ladder_participants (ladder_id, user_id, final_balance, final_rank)
VALUES ($1, $2, $3, $4)
ON CONFLICT (ladder_id, user_id) DO UPDATE SET
    final_balance = EXCLUDED.final_balance,
    final_rank = EXCLUDED.final_rank;

-- name: JoinLadderParticipant :exec
INSERT INTO ladder_participants (ladder_id, user_id)
VALUES ($1, $2)
ON CONFLICT (ladder_id, user_id) DO NOTHING;

-- name: IsUserInLadder :one
SELECT EXISTS(
    SELECT 1 FROM ladder_participants
    WHERE ladder_id = $1 AND user_id = $2
);

-- name: GetLadderBalance :one
SELECT balance
FROM ladder_balances
WHERE ladder_id = $1 AND user_id = $2 LIMIT 1;

-- name: GetLadderBalanceForUpdate :one
SELECT balance
FROM ladder_balances
WHERE ladder_id = $1 AND user_id = $2 LIMIT 1
FOR UPDATE;

-- name: InsertLadderBalance :exec
INSERT INTO ladder_balances (ladder_id, user_id, balance)
VALUES ($1, $2, $3)
ON CONFLICT (ladder_id, user_id) DO NOTHING;

-- name: UpdateLadderBalance :exec
UPDATE ladder_balances
SET balance = $3
WHERE ladder_id = $1 AND user_id = $2;

-- name: GetLadderPortfolio :many
SELECT ladder_id, user_id, stock_symbol, quantity, average_price
FROM ladder_portfolio_items
WHERE ladder_id = $1 AND user_id = $2;

-- name: GetLadderPortfolioItem :one
SELECT ladder_id, user_id, stock_symbol, quantity, average_price
FROM ladder_portfolio_items
WHERE ladder_id = $1 AND user_id = $2 AND stock_symbol = $3 LIMIT 1;

-- name: GetLadderPortfolioItemForUpdate :one
SELECT ladder_id, user_id, stock_symbol, quantity, average_price
FROM ladder_portfolio_items
WHERE ladder_id = $1 AND user_id = $2 AND stock_symbol = $3 LIMIT 1
FOR UPDATE;

-- name: SetLadderPortfolioItem :exec
INSERT INTO ladder_portfolio_items (ladder_id, user_id, stock_symbol, quantity, average_price)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (ladder_id, user_id, stock_symbol) DO UPDATE SET
    quantity = EXCLUDED.quantity,
    average_price = EXCLUDED.average_price;

-- name: DeleteLadderPortfolioItem :exec
DELETE FROM ladder_portfolio_items
WHERE ladder_id = $1 AND user_id = $2 AND stock_symbol = $3;

-- name: DeleteLadderPortfolioItemsByLadder :exec
DELETE FROM ladder_portfolio_items
WHERE ladder_id = $1;

-- name: DeleteLadderBalancesByLadder :exec
DELETE FROM ladder_balances
WHERE ladder_id = $1;
