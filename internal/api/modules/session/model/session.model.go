package sessions

import (
	"time"

	"github.com/google/uuid"
	"github.com/unarya/univia/internal/api/modules/key_token/refresh_token/models"
	"github.com/unarya/univia/internal/api/modules/user/models"
)

type UserSession struct {
	ID             uuid.UUID                  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SessionID      uuid.UUID                  `gorm:"type:char(36);not null"`
	UserID         uuid.UUID                  `gorm:"type:char(36);uniqueIndex;not null"`
	User           users.User                 `gorm:"foreignKey:UserID;references:ID"`
	IP             string                     `gorm:"type:varchar(64)"`
	UserAgent      string                     `gorm:"type:text"`
	RefreshTokenID uuid.UUID                  `gorm:"type:index;not null"`
	RefreshTokens  refresh_token.RefreshToken `gorm:"foreignKey:RefreshTokenID;references:ID"`
	Status         string                     `gorm:"type:varchar(32);default:'active'"`
	LastActive     *time.Time                 `json:"last_active,omitempty"`
	RevokedAt      *time.Time                 `json:"revoked_at,omitempty"`
	CreatedAt      time.Time                  `json:"created_at"`
	UpdatedAt      time.Time                  `json:"updated_at"`
}
