package posts

import (
	"time"

	"github.com/google/uuid"
)

type PostCategory struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null"`
	PostID     uuid.UUID `gorm:"type:uuid;not null"`

	// References
	Category Category `gorm:"foreignKey:CategoryID;references:ID"`
	Post     Post     `gorm:"foreignKey:PostID;references:ID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
