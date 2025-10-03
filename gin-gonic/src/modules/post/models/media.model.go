package posts

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PostID    uuid.UUID `gorm:"type:uuid;not null"`
	Post      Post      `gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE;"`
	Path      string    `gorm:"type:varchar(255);not null"`
	Type      string    `gorm:"type:varchar(255);not null"`
	Status    bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
