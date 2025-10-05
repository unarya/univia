-- +migrate Up
CREATE TABLE IF NOT EXISTS verification_codes (
    id CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
    email VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL,
    expires_at DATETIME NOT NULL,
    input_count INT(5) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Indexing
CREATE INDEX idx_verification_codes_email ON verification_codes (email);
