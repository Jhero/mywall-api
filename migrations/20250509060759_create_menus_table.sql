-- Migration: create_menuss_table
-- Created at: 2025-05-09T06:07:59+07:00
-- Up

-- Write your up migration here
CREATE TABLE IF NOT EXISTS menus (
    id VARCHAR(30) PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    path VARCHAR(50) NOT NULL,
    UNIQUE KEY unique_path (path)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Down
-- Uncomment if you want to use down migrations

-- Write your down migration here
