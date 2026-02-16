-- name: CreateQuote :exec
INSERT INTO market_quotes (symbol, price, source, created_at)
VALUES ($1, $2, $3, $4);

-- name: GetHistoryForSymbol :many
WITH latest_quotes AS (
    SELECT symbol, price, source, created_at
    FROM market_quotes
    WHERE symbol = $1
    ORDER BY created_at DESC
    LIMIT $2
)
SELECT * FROM latest_quotes ORDER BY created_at ASC;
