package access_token

import (
	"time"

	Users "github.com/deva-labs/univia/api/gin/src/modules/user/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccessToken represents the access token model.
type AccessToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null"`
	User      Users.User `gorm:"foreignKey:UserID;references:ID"`
	Token     string     `gorm:"type:varchar(256);not null;index"`
	Status    bool       `gorm:"default:true"`
	ExpiresAt time.Time  `gorm:"not null"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
}

// BeforeCreate sets the default expiration time for the access token.
func (a *AccessToken) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ExpiresAt.IsZero() {
		a.ExpiresAt = time.Now().Add(7 * 24 * time.Hour) // 7 days from now
	}
	return
}
