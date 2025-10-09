package sessions

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/unarya/univia/internal/api/modules/session/model"
	"github.com/unarya/univia/internal/api/modules/user/models"
	"github.com/unarya/univia/internal/infrastructure/mysql"
	"gorm.io/gorm"
)

func RevokeSession(db *gorm.DB, sessionID string) error {
	now := time.Now()
	return db.Model(&sessions.UserSession{}).
		Where("session_id = ?", sessionID).
		Updates(map[string]interface{}{
			"status":     "revoked",
			"revoked_at": now,
		}).Error
}

func CheckValidDevice(email string, sessionID uuid.UUID) bool {
	db := mysql.DB
	var results []struct {
		users.User           `gorm:"embedded"` // Embed struct User
		sessions.UserSession `gorm:"embedded"` // Embed struct UserSession
	}

	err := db.Table("users").
		Select("users.*, user_sessions.*").
		Joins("inner join user_sessions on user_sessions.user_id = users.id").
		Where("users.email = ? AND user_sessions.session_id = ?", email, sessionID).
		Scan(&results).Error
	if err != nil {
		log.Fatal(err)
		return false
	}

	// Process results
	for _, r := range results {
		if r.UserSession.Status != "revoked" {
			return true
		}
	}
	return false
}
