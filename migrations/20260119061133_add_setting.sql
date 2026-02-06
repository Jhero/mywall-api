-- Migration: add_setting
-- Created at: 2026-01-19T06:11:33+07:00
-- Up

-- Write your up migration here
CREATE TABLE settings (
    id VARCHAR(40) PRIMARY KEY,
    app_id VARCHAR(40) NOT NULL,
    user_id INT,
    `key` VARCHAR(100) NOT NULL,
    value JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_settings_app FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE,
    CONSTRAINT fk_settings_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Index untuk mempercepat pencarian setting per app dan user
CREATE UNIQUE INDEX idx_settings_app_user_key
    ON settings(app_id, user_id, `key`);

-- Down
-- Uncomment if you want to use down migrations

-- Write your down migration here
