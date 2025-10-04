package posts

import (
	Users "github.com/deva-labs/univia-api/api/gin-gonic/src/modules/user/models"
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PostID uuid.UUID `gorm:"type:uuid;not null"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	Text   string    `gorm:"type:text;not null"`
	Left   int       `gorm:"not null"`
	Right  int       `gorm:"not null"`

	// References
	Post Post       `gorm:"foreignKey:PostID;references:ID"`
	User Users.User `gorm:"foreignKey:UserID;references:ID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
