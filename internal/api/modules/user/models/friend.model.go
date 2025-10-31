package users

import (
	"time"

	"github.com/google/uuid"
)

// Friend represents the friends table, establishing relationships between users.
type Friend struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      uuid.UUID `gorm:"not null"` // The user who initiated the friend request
	FriendTo    uuid.UUID `gorm:"not null"` // The user who is the friend
	RequestedOn time.Time `gorm:"not null"`
	AcceptedOn  time.Time `gorm:"default:null"`
	Description string    `gorm:"type:text;default:null"`
	Status      bool      `gorm:"default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// Relationships
	User   User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`   // Relation to the User who initiated the request
	Friend User `gorm:"foreignKey:FriendTo;constraint:OnDelete:CASCADE;"` // Relation to the User who is the friend
}
