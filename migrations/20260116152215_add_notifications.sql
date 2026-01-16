-- Migration: add_notifications
-- Created at: 2026-01-16T15:22:15+07:00
-- Up

-- Write your up migration here
CREATE TABLE IF NOT EXISTS notifications (
  id VARCHAR(40) PRIMARY KEY,
  user_id INT NOT NULL,
  title TEXT NOT NULL,
  body TEXT NOT NULL,
  type VARCHAR(50) NOT NULL,  -- e.g. "system", "message", "alert"
  metadata JSON,
  is_read INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX notifications_user_created_at_idx ON notifications (user_id, created_at);
CREATE INDEX notifications_user_isread_idx ON notifications (user_id, is_read);

-- Down
-- Uncomment if you want to use down migrations

-- Write your down migration here
