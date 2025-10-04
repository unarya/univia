-- +migrate Up
CREATE TABLE IF NOT EXISTS comment_likes (
    id CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
    comment_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- Foreign keys
    CONSTRAINT fk_comment_likes_comment FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    CONSTRAINT fk_comment_likes_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Indexing
CREATE INDEX idx_comment_likes_comment_id ON comment_likes (comment_id);
CREATE INDEX idx_comment_likes_user_id ON comment_likes (user_id);
CREATE UNIQUE INDEX uq_comment_likes_comment_user ON comment_likes (comment_id, user_id);
