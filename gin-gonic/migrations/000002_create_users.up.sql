-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number BIGINT UNSIGNED UNIQUE,
    google_id VARCHAR(255) UNIQUE,
    twitter_id VARCHAR(255) UNIQUE,
    password VARCHAR(255) DEFAULT NULL,
    status BOOLEAN DEFAULT TRUE,
    role_id CHAR(36) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Indexing
CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_role_id ON users (role_id);
