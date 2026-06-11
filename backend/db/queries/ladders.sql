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

-- name: GetExpiredActiveLadders :many
SELECT id, name, type, start_time, end_time, initial_balance, is_active, created_at
FROM ladders
WHERE is_active = TRUE AND end_time <= $1;

-- name: GetPendingLaddersToActivate :many
SELECT id, name, type, start_time, end_time, initial_balance, is_active, created_at
FROM ladders
WHERE is_active = FALSE AND start_time <= $1 AND end_time > $1;


-- name: GetLadder :one
SELECT id, name, type, start_time, end_time, initial_balance, is_active, created_at
FROM ladders
WHERE id = $1;

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
SELECT ladder_id, user_id, balance, final_balance, final_rank, joined_at
FROM ladder_participants
WHERE ladder_id = $1;


-- name: InsertLadderParticipant :exec
INSERT INTO ladder_participants (ladder_id, user_id, final_balance, final_rank)
VALUES ($1, $2, $3, $4)
ON CONFLICT (ladder_id, user_id) DO UPDATE SET
    final_balance = EXCLUDED.final_balance,
    final_rank = EXCLUDED.final_rank;

-- name: JoinLadderParticipant :exec
INSERT INTO ladder_participants (ladder_id, user_id, balance)
SELECT $1, $2, initial_balance FROM ladders WHERE id = $1;

-- name: IsUserInLadder :one
SELECT EXISTS(
    SELECT 1 FROM ladder_participants
    WHERE ladder_id = $1 AND user_id = $2
);

-- name: GetLadderParticipantBalance :one
SELECT balance
FROM ladder_participants
WHERE ladder_id = $1 AND user_id = $2;

-- name: GetLadderParticipantBalanceForUpdate :one
SELECT balance
FROM ladder_participants
WHERE ladder_id = $1 AND user_id = $2
FOR UPDATE;

-- name: UpdateLadderParticipantBalance :exec
UPDATE ladder_participants
SET balance = $3
WHERE ladder_id = $1 AND user_id = $2;

-- name: GetLadderPortfolio :many
SELECT ladder_id, user_id, stock_symbol, quantity, average_price
FROM ladder_portfolio_items
WHERE ladder_id = $1 AND user_id = $2;

-- name: GetLadderPortfolioItem :one
SELECT ladder_id, user_id, stock_symbol, quantity, average_price
FROM ladder_portfolio_items
WHERE ladder_id = $1 AND user_id = $2 AND stock_symbol = $3;

-- name: GetLadderPortfolioItemForUpdate :one
SELECT ladder_id, user_id, stock_symbol, quantity, average_price
FROM ladder_portfolio_items
WHERE ladder_id = $1 AND user_id = $2 AND stock_symbol = $3
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

-- name: PruneLadderParticipants :exec
DELETE FROM ladder_participants
WHERE ladder_id = $1 AND final_rank > $2;


