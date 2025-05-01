-- Migration: create_galleries_table
-- Created at: 2025-04-28T06:19:06+07:00
-- Up

CREATE TABLE IF NOT EXISTS galleries (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image_url VARCHAR(2048),
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_galleries_user_id ON galleries(user_id);

-- Down
-- Uncomment if you want to use down migrations

-- Write your down migration here
