-- +migrate Up
CREATE TABLE IF NOT EXISTS access_tokens (
     id CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
     user_id CHAR(36) NOT NULL,
     token VARCHAR(256) NOT NULL,
     status BOOLEAN NOT NULL DEFAULT TRUE,
     expires_at DATETIME NOT NULL,
     created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
     updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Indexing
CREATE INDEX idx_access_tokens_token_status ON access_tokens(token, status);
CREATE INDEX idx_access_token_user_status ON access_tokens(user_id, status);