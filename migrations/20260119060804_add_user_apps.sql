-- Migration: add_user_apps
-- Created at: 2026-01-19T06:08:04+07:00
-- Up

-- Write your up migration here
CREATE TABLE user_apps (
    id VARCHAR(40) PRIMARY KEY,
    user_id INT NOT NULL,
    app_id VARCHAR(40) NOT NULL,
    role_id VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_app (user_id, app_id)
);
CREATE INDEX idx_user_apps_user_app ON user_apps(user_id, app_id);

-- Down
-- Uncomment if you want to use down migrations

-- Write your down migration here
