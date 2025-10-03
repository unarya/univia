package posts

import (
	"time"
	Users "univia/src/modules/user/models"

	"github.com/google/uuid"
)

type CommentLike struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CommentID uuid.UUID `gorm:"type:uuid;not null"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`

	// References
	Comment Comment    `gorm:"foreignKey:CommentID;references:ID"`
	User    Users.User `gorm:"foreignKey:UserID;references:ID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
