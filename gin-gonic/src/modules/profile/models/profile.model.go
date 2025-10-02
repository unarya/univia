package models

import (
	"time"
	Users "univia/src/modules/user/models"

	"github.com/google/uuid"
)

type Profile struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null;unique"`
	User            Users.User `gorm:"foreignKey:UserID;references:ID"`
	ProfilePic      string     `gorm:"type:text;default:null"`
	CoverPhoto      string     `gorm:"type:varchar(255);default:null"`
	BackgroundColor string     `gorm:"type:varchar(255);default:#7b2cbf"`
	Gender          string     `gorm:"type:varchar(10);default:null"`
	Birthday        *time.Time `gorm:"type:date"`
	Location        string     `gorm:"type:varchar(255);default:null"`
	Bio             string     `gorm:"type:text;default:null"`
	Interests       []string   `gorm:"type:json;default:null"`
	SocialLinks     []string   `gorm:"type:json;default:null"`

	// Many-to-Many Relationship (Followers & Followings)
	Followers  []Users.User `gorm:"many2many:follows;joinForeignKey:FollowerID;joinReferences:FollowingID"`
	Followings []Users.User `gorm:"many2many:follows;joinForeignKey:FollowingID;joinReferences:FollowerID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
