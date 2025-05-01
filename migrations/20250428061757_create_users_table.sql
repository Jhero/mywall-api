-- Migration: create_users_table
-- Created at: 2025-04-28T06:17:57+07:00
-- Up

-- Write your up migration here

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    sub VARCHAR(255) NOT NULL,
    UNIQUE KEY unique_email (email),
    UNIQUE KEY unique_sub (sub)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Down
-- Uncomment if you want to use down migrations

-- Write your down migration here
