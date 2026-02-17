-- +goose Up
INSERT INTO users (username, password_hash, first_name, last_name, balance, website, created_at)
VALUES
    ('cool_trader', '$2a$10$iehMYy8l.nvHJ.CxRqmtJOhfRJHw1kpnQfHDmJX7Qorq1QVeI.1lK', 'Cool', 'Trader', 10000, 'https://cool.io', NOW()),  -- PW: test123
    ('rich_guy', '$2a$10$iehMYy8l.nvHJ.CxRqmtJOhfRJHw1kpnQfHDmJX7Qorq1QVeI.1lK', 'Rich', 'Guy', 10000, '', NOW());  -- PW: test123

-- +goose Down
DELETE FROM users WHERE username IN ('cool_trader', 'rich_guy');
