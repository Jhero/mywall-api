-- Migration: add_field_categoryid_galery
-- Created at: 2025-05-04T07:27:44+07:00
-- Up

-- Write your up migration here
ALTER TABLE rbacs
    ADD COLUMN owner_id INT NOT NULL,
    ADD CONSTRAINT fk_rbacs_owner FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX idx_rbacs_owner_id ON rbacs(owner_id);

-- Down
-- Uncomment if you want to use down migrations
ALTER TABLE rbacs
    DROP INDEX idx_rbacs_owner_id,
    DROP FOREIGN KEY fk_rbacs_owner,
    DROP COLUMN owner_id;
-- Write your down migration here
