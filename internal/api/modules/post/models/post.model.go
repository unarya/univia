package posts

import (
	"time"

	Users "github.com/unarya/univia/internal/api/modules/user/models"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null"`
	User      Users.User `gorm:"foreignKey:UserID;references:ID"`
	Content   string     `gorm:"type:text;default:null"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
}
