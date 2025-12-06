-- +goose Up
INSERT INTO users (id, email, password_hash, first_name, last_name, balance, created_at)
VALUES 
    (1, 'user1@example.com', '$2a$10$eyNvBrLTyoH91AAYnEWfZe0ba6lh/Ta3UBMBGVHegrFG1uGNJ3oji', 'Testo', 'Matesto', 10000.0, NOW()),
    (2, 'user2@example.com', '$2a$10$eyNvBrLTyoH91AAYnEWfZe0ba6lh/Ta3UBMBGVHegrFG1uGNJ3oji', 'Loba', 'Boba', 5000.0, NOW()),
    (3, 'user3@example.com', '$2a$10$eyNvBrLTyoH91AAYnEWfZe0ba6lh/Ta3UBMBGVHegrFG1uGNJ3oji', 'Tian', 'Petyan', 10000.0, NOW())
ON CONFLICT (id) DO NOTHING;

SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));

-- +goose Down
DELETE FROM users WHERE id IN (1, 2, 3);
