-- Migration: add_field_roleid_rbac
-- Created at: 2025-06-29T22:33:55+07:00
-- Up

-- Write your up migration here
ALTER TABLE rbacs
    ADD COLUMN role_id VARCHAR(20) NOT NULL,
    ADD CONSTRAINT fk_rbacs_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;

CREATE INDEX idx_rbacs_role_id ON rbacs(role_id);

-- Down
-- Uncomment if you want to use down migrations

-- Write your down migration here
