-- Migration: add_field_users
-- Created at: 2025-05-01T10:29:57+07:00
-- Up

-- Write your up migration here
ALTER TABLE users
    ADD COLUMN password VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN api_key VARCHAR(64) UNIQUE,
    ADD COLUMN role VARCHAR(20) DEFAULT 'user',
    ADD COLUMN is_active BOOLEAN DEFAULT true,
    DROP COLUMN provider,
    DROP COLUMN sub;

-- Down
-- Uncomment if you want to use down migrations

-- Write your down migration here
