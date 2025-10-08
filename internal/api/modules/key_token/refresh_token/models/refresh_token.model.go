package refresh_token

import (
	"time"

	Users "github.com/unarya/univia/internal/api/modules/user/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshToken represents the refresh token model.
type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null"`
	User      Users.User `gorm:"foreignKey:UserID;references:ID"`
	Token     string     `gorm:"type:varchar(256);not null"`
	Status    bool       `gorm:"default:true"` // true = active, false = revoked
	ExpiresAt time.Time  `gorm:"not null"`     // Token expiration time
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
}

// BeforeCreate sets the default expiration time for the refresh token.
func (r *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	r.ExpiresAt = time.Now().Add(90 * 24 * time.Hour) // 30 days from now
	return
}
