-- name: UpsertStock :one
INSERT INTO stocks (symbol, name)
VALUES ($1, $2)
ON CONFLICT (symbol) DO UPDATE SET
    name = EXCLUDED.name
RETURNING symbol, name, created_at;

-- name: GetStock :one
SELECT symbol, name, created_at
FROM stocks
WHERE symbol = $1 LIMIT 1;

-- name: ListStocks :many
SELECT symbol, name, created_at
FROM stocks
ORDER BY symbol;
