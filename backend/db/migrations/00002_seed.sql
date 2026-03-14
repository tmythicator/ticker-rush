-- +goose Up
INSERT INTO users (username, password_hash, first_name, last_name, website, created_at, agb_accepted_at, is_admin)
VALUES
    ('cool_trader', '$2a$10$iehMYy8l.nvHJ.CxRqmtJOhfRJHw1kpnQfHDmJX7Qorq1QVeI.1lK', 'Cool', 'Trader', 'https://example.com', NOW(), NOW(), FALSE),  -- PW: test123
    ('rich_guy', '$2a$10$iehMYy8l.nvHJ.CxRqmtJOhfRJHw1kpnQfHDmJX7Qorq1QVeI.1lK', 'Rich', 'Guy', '', NOW(), NOW(), FALSE),  -- PW: test123
    ('admin', '$2a$10$iehMYy8l.nvHJ.CxRqmtJOhfRJHw1kpnQfHDmJX7Qorq1QVeI.1lK', 'Admin', 'Admin', '', NOW(), NOW(), TRUE);  -- PW: test123

-- +goose Down
DELETE FROM users WHERE username IN ('cool_trader', 'rich_guy', 'admin');
