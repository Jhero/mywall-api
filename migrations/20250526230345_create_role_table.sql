-- Migration: create_role_table
-- Created at: 2025-05-26T23:03:45+07:00
-- Up
CREATE TABLE IF NOT EXISTS roles (
    id VARCHAR(20) PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_roles_user_id ON roles(user_id);

-- Write your up migration here

-- Down
-- Uncomment if you want to use down migrations

-- Write your down migration here
