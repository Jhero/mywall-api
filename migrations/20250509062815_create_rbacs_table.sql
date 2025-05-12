-- Migration: create_rbacs_table
-- Created at: 2025-05-09T05:28:15+07:00
-- Up

-- Write your up migration here
CREATE TABLE IF NOT EXISTS rbacs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    menu_id VARCHAR(50) NOT NULL,
    permission VARCHAR(200) NOT NULL,
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_rbacs_user_id ON rbacs(user_id);
CREATE INDEX idx_rbacs_menu_id ON rbacs(menu_id);

-- Down
-- Uncomment if you want to use down migrations

-- Write your down migration here
