-- name: CreateUser :one
INSERT INTO users (username, password_hash, first_name, last_name, website, created_at, is_public, agb_accepted_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, username, first_name, last_name, website, created_at, is_public, is_admin;


-- name: UpdateUserProfile :exec
UPDATE users
SET first_name = $2,
    last_name = $3,
    website = $4,
    is_public = $5
WHERE id = $1;


-- name: BanUser :exec
UPDATE users
SET is_banned = TRUE
WHERE id = $1;

-- name: CheckUserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);

-- name: GetUser :one
SELECT id, username, first_name, last_name, website, created_at, is_public, is_admin
FROM users
WHERE id = $1 LIMIT 1;


-- name: GetUserForUpdate :one
SELECT id, username, first_name, last_name, website, created_at, is_public, is_admin
FROM users
WHERE id = $1 LIMIT 1 FOR UPDATE;


-- name: GetUserByUsername :one
SELECT id, username, password_hash, first_name, last_name, website, created_at, is_public, is_admin
FROM users
WHERE username = $1 LIMIT 1;


-- name: GetUsers :many
SELECT id, username, first_name, last_name, website, created_at, is_public, is_admin
FROM users;