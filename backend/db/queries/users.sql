-- name: CreateUser :one
INSERT INTO users (email, password_hash, first_name, last_name, balance, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, email, first_name, last_name, balance, created_at;

-- name: UpdateUser :exec
UPDATE users
SET email = $2,
    first_name = $3,
    last_name = $4,
    balance = $5
WHERE id = $1;

-- name: CheckUserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);

-- name: UpdateUserBalance :exec
UPDATE users
SET balance = $2
WHERE id = $1;

-- name: GetUser :one
SELECT id, email, first_name, last_name, balance, created_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserForUpdate :one
SELECT id, email, first_name, last_name, balance, created_at
FROM users
WHERE id = $1 LIMIT 1 FOR UPDATE;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, first_name, last_name, balance, created_at
FROM users
WHERE email = $1 LIMIT 1;
