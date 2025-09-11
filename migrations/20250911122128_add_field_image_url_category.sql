-- Migration: add_field_image_url_category
-- Created at: 2025-09-11T12:21:28+07:00
-- Up

-- Write your up migration here
ALTER TABLE categories
    ADD COLUMN image_url VARCHAR(255) NOT NULL;

-- Down
-- Uncomment if you want to use down migrations

-- Write your down migration here
