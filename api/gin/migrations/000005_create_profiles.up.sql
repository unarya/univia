-- +migrate Up
CREATE TABLE IF NOT EXISTS profiles (
    id CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL UNIQUE,
    profile_pic TEXT DEFAULT NULL,
    cover_photo VARCHAR(255) DEFAULT NULL,
    background_color VARCHAR(255) DEFAULT '#7b2cbf',
    gender VARCHAR(10) DEFAULT NULL,
    birthday DATE DEFAULT NULL,
    location VARCHAR(255) DEFAULT NULL,
    bio TEXT DEFAULT NULL,
    interests JSON DEFAULT NULL,
    social_links JSON DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_profiles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Indexing
CREATE INDEX idx_profiles_user_id ON profiles (user_id);
