-- name: GetUser :one
SELECT id, email, password_hash, balance, created_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (id, email, password_hash, balance, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, email, password_hash, balance, created_at;

-- name: UpsertUser :exec
INSERT INTO users (id, email, password_hash, balance, created_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO UPDATE SET
    email = EXCLUDED.email,
    password_hash = EXCLUDED.password_hash,
    balance = EXCLUDED.balance,
    created_at = EXCLUDED.created_at;

-- name: CheckUserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);
