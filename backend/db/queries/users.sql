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
WHERE id = $1;


-- name: GetUserForUpdate :one
SELECT id, username, first_name, last_name, website, created_at, is_public, is_admin
FROM users
WHERE id = $1 FOR UPDATE;


-- name: GetUserByUsername :one
SELECT id, username, password_hash, first_name, last_name, website, created_at, is_public, is_admin
FROM users
WHERE username = $1;


-- name: GetUsers :many
SELECT id, username, first_name, last_name, website, created_at, is_public, is_admin
FROM users;

-- name: GetUserWithPortfolioForActiveLadder :many
SELECT u.id AS user_id,
       u.username,
       u.first_name,
       u.last_name,
       u.website,
       u.created_at,
       u.is_public,
       u.is_admin,
       u.is_banned,
       COALESCE(al.id, 0)::bigint AS ladder_id,
       COALESCE(lp.balance, al.initial_balance, 0.0) AS balance,
       lpi.stock_symbol,
       COALESCE(lpi.quantity, 0.0)::float8 AS quantity,
       COALESCE(lpi.average_price, 0.0)::float8 AS average_price,
       (lp.user_id IS NOT NULL)::boolean AS is_participating
FROM users u
LEFT JOIN ladders al ON al.is_active = true
LEFT JOIN ladder_participants lp ON u.id = lp.user_id AND lp.ladder_id = al.id
LEFT JOIN ladder_portfolio_items lpi ON u.id = lpi.user_id AND lpi.ladder_id = al.id
WHERE u.id = $1;

-- name: AnonymizeUser :exec
UPDATE users
SET username = 'deleted_' || substring(md5(random()::text) from 1 for 12),
    password_hash = 'DELETED',
    first_name = 'Deleted',
    last_name = 'User',
    website = '',
    is_public = false
WHERE id = $1;