-- Migration: add_field_userid_categories
-- Created at: 2025-05-04T07:52:10+07:00
-- Up

-- Write your up migration here
ALTER TABLE categories
    ADD COLUMN user_id INT NOT NULL,
    ADD CONSTRAINT fk_categories_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX idx_categories_user_id ON categories(user_id);

-- Down
-- Uncomment if you want to use down migrations

-- Write your down migration here
