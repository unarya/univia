-- +migrate Up
CREATE TABLE IF NOT EXISTS media (
    id CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
    post_id CHAR(36) NOT NULL,
    path VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    status BOOL DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- Foreign keys
    CONSTRAINT fk_media_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Indexing
CREATE INDEX idx_media_post_id ON media (post_id);