-- +migrate Up
CREATE TABLE IF NOT EXISTS notifications (
     id CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
     sender_id CHAR(36) NOT NULL,
     receiver_id CHAR(36) NOT NULL,
     message TEXT DEFAULT NULL,
     is_seen BOOLEAN DEFAULT FALSE,
     noti_type VARCHAR(50) DEFAULT NULL,
     created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
     updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- Foreign keys
     CONSTRAINT fk_notifications_receiver FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE,
     CONSTRAINT fk_notifications_sender FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Indexing
-- BTREE indexing
CREATE INDEX idx_notifications_receiver_id ON notifications(receiver_id);
CREATE INDEX idx_notifications_sender_id ON notifications(sender_id);

-- FULLTEXT indexing
CREATE FULLTEXT INDEX idx_notifications_message ON notifications(message);