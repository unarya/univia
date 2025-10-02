package models

import (
	"time"

	"github.com/google/uuid"
)

type VerificationCode struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email      string    `gorm:"type:varchar(255);not null"`
	Code       string    `gorm:"type:varchar(255);not null"`
	ExpiresAt  time.Time `gorm:"type:datetime;not null"`
	InputCount int       `gorm:"type:int(5);not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime;"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime;"`
}
