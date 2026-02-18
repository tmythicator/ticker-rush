-- name: CreateUser :one
INSERT INTO users (username, password_hash, first_name, last_name, balance, website, created_at, is_public)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, username, first_name, last_name, balance, website, created_at, is_public;


-- name: UpdateUserProfile :exec
UPDATE users
SET first_name = $2,
    last_name = $3,
    website = $4,
    is_public = $5
WHERE id = $1;


-- name: UpdateUserBalance :exec
UPDATE users
SET balance = $2
WHERE id = $1;

-- name: CheckUserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);

-- name: GetUser :one
SELECT id, username, first_name, last_name, balance, website, created_at, is_public
FROM users
WHERE id = $1 LIMIT 1;


-- name: GetUserForUpdate :one
SELECT id, username, first_name, last_name, balance, website, created_at, is_public
FROM users
WHERE id = $1 LIMIT 1 FOR UPDATE;


-- name: GetUserByUsername :one
SELECT id, username, password_hash, first_name, last_name, balance, website, created_at, is_public
FROM users
WHERE username = $1 LIMIT 1;


-- name: GetUsers :many
SELECT id, username, first_name, last_name, balance, website, created_at, is_public
FROM users;