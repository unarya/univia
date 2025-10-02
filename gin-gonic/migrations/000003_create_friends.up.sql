-- +migrate Up
CREATE TABLE IF NOT EXISTS friends (
    id CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    friend_to CHAR(36) NOT NULL,
    requested_on DATETIME NOT NULL,
    accepted_on DATETIME DEFAULT NULL,
    description TEXT DEFAULT NULL,
    status BOOLEAN DEFAULT FALSE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_friends_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_friends_friend FOREIGN KEY (friend_to) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Indexing
CREATE INDEX idx_friends_user_id ON friends (user_id);
CREATE INDEX idx_friends_friend_to ON friends (friend_to);
CREATE UNIQUE INDEX uq_friend_pair ON friends (user_id, friend_to);
