-- Migration: add_apps
-- Created at: 2026-01-19T06:06:23+07:00
-- Up

-- Write your up migration here
CREATE TABLE apps (
    id VARCHAR(40) PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Down
-- Uncomment if you want to use down migrations

-- Write your down migration here
