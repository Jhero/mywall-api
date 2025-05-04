-- Migration: add_field_category_galery
-- Created at: 2025-05-04T05:53:46+07:00
-- Up

ALTER TABLE galleries
    ADD COLUMN category VARCHAR(50) NOT NULL DEFAULT ''

-- Down
ALTER TABLE galleries
    DROP COLUMN category
-- Uncomment if you want to use down migrations

-- Write your down migration here
