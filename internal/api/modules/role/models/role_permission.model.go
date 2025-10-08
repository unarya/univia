package roles

import (
	"time"

	Permissions "github.com/unarya/univia/internal/api/modules/permission/models"

	"github.com/google/uuid"
)

type RolePermission struct {
	ID           uuid.UUID              `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RoleID       uuid.UUID              `gorm:"not null"`
	PermissionID uuid.UUID              `gorm:"not null"`
	Role         Role                   `gorm:"foreignKey:RoleID;references:ID"`
	Permission   Permissions.Permission `gorm:"foreignKey:PermissionID;references:ID"`
	CreatedAt    time.Time              `gorm:"autoCreateTime"`
	UpdatedAt    time.Time              `gorm:"autoUpdateTime"`
}
