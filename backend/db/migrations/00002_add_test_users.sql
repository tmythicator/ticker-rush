-- +goose Up
INSERT INTO users (id, email, password_hash, first_name, last_name, balance, created_at)
VALUES 
    (1, 'user1@example.com', '$2a$10$iehMYy8l.nvHJ.CxRqmtJOhfRJHw1kpnQfHDmJX7Qorq1QVeI.1lK', 'Testo', 'Matesto', 10000.0, NOW()), -- PW: test123
    (2, 'user2@example.com', '$2a$10$iehMYy8l.nvHJ.CxRqmtJOhfRJHw1kpnQfHDmJX7Qorq1QVeI.1lK', 'Loba', 'Boba', 5000.0, NOW()), -- PW: test123
    (3, 'user3@example.com', '$2a$10$iehMYy8l.nvHJ.CxRqmtJOhfRJHw1kpnQfHDmJX7Qorq1QVeI.1lK', 'Tian', 'Petyan', 10000.0, NOW()) -- PW: test123
ON CONFLICT (email) DO NOTHING;

SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));

-- +goose Down
DELETE FROM users WHERE id IN (1, 2, 3);
