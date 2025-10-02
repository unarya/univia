package models

import (
	"time"
	Roles "univia/src/modules/role/models"

	"github.com/google/uuid"
)

// User represents the user model
type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username    string     `gorm:"type:varchar(255);not null;"`
	Email       string     `gorm:"type:varchar(255);not null;unique"`
	PhoneNumber uint64     `gorm:"default:null;unique"`
	GoogleID    string     `gorm:"type:varchar(255);default:null;unique"`
	TwitterID   string     `gorm:"type:varchar(255);default:null;unique"`
	Password    string     `gorm:"type:varchar(255);default:null"`
	Status      bool       `gorm:"default:true"`
	RoleID      uuid.UUID  `gorm:"type:uuid;not null"`
	Role        Roles.Role `gorm:"foreignKey:RoleID;references:ID"` // Foreign key to Role
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
}
