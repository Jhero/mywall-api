-- Migration: add_field_categoryid_galery
-- Created at: 2025-05-04T07:27:44+07:00
-- Up

-- Write your up migration here
ALTER TABLE galleries
    ADD COLUMN category_id INT NOT NULL,
    ADD CONSTRAINT fk_galleries_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE;

CREATE INDEX idx_galleries_category_id ON galleries(category_id);

-- Down
-- Uncomment if you want to use down migrations
ALTER TABLE galleries
    DROP INDEX idx_galleries_category_id,
    DROP FOREIGN KEY fk_galleries_category,
    DROP COLUMN category_id;
-- Write your down migration here
