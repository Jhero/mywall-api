-- Migration: create_image_view.sql
-- Created at: 2025-11-26T23:37:42+07:00
-- Up

CREATE TABLE IF NOT EXISTS image_views (
    gallery_id VARCHAR(20) PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    count INT DEFAULT 0,
    user_id INT NOT NULL,
    FOREIGN KEY (gallery_id) REFERENCES galleries(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_image_views_gallery_id ON image_views(gallery_id);

-- Down
-- Uncomment if you want to use down migrations

-- Write your down migration here
