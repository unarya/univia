-- +migrate Up
CREATE TABLE IF NOT EXISTS user_sessions (
    id CHAR(36) PRIMARY KEY,
    session_id CHAR(36) NOT NULL UNIQUE,
    user_id CHAR(36) NOT NULL,
    ip VARCHAR(64),
    user_agent TEXT,
    refresh_token_id CHAR(36) NOT NULL,
    status VARCHAR(32) DEFAULT 'active',
    last_active DATETIME DEFAULT NULL,
    revoked_at DATETIME DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_sessions_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_user_sessions_refresh_token FOREIGN KEY (refresh_token_id) REFERENCES refresh_tokens(id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Indexing
CREATE INDEX idx_user_sessions_user_id ON user_sessions (user_id);
CREATE INDEX idx_user_sessions_rtoken_id ON user_sessions (refresh_token_id);

