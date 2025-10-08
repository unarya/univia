package roles

import (
	"errors"

	roles "github.com/unarya/univia/internal/api/modules/role/models"
	Users "github.com/unarya/univia/internal/api/modules/user/models"
	"github.com/unarya/univia/internal/infrastructure/mysql"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetRoleByUserID retrieves the role of a user by their user ID
func GetRoleByUserID(userID uint) (*roles.Role, error) {
	db := mysql.DB

	// Fetch user to get RoleID
	var user Users.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Fetch role based on RoleID
	var role roles.Role
	if err := db.Where("id = ?", user.RoleID).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return &role, nil
}

// CreateRoleByAdmin creates a new role with the given name
func CreateRoleByAdmin(roleName string) (map[string]interface{}, error) {
	db := mysql.DB

	// Validate input
	if roleName == "" {
		return nil, errors.New("role name cannot be empty")
	}

	// Check if role already exists
	var existingRole roles.Role
	if err := db.Where("name = ?", roleName).First(&existingRole).Error; err == nil {
		return nil, errors.New("role already exists")
	}

	// Create new role
	role := &roles.Role{
		Name: roleName,
	}

	if err := db.Create(role).Error; err != nil {
		return nil, err
	}

	return gin.H{
		"id":   role.ID,
		"name": roleName,
	}, nil
}

func ListAllRoles() ([]map[string]interface{}, error) {
	db := mysql.DB
	var roles []roles.Role
	if err := db.Find(&roles).Error; err != nil {
		return nil, err
	}
	// Convert to a slice of maps
	var result []map[string]interface{}
	for _, roles := range roles {
		result = append(result, gin.H{
			"id":   roles.ID,
			"name": roles.Name,
		})
	}
	return result, nil
}

// GetRoleID is a function given roleId by name
func GetRoleID(name string) (uuid.UUID, error) {
	var id uuid.UUID
	db := mysql.DB
	if err := db.Model(&roles.Role{}).
		Where("name = ?", name).
		Pluck("id", &id).Error; err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
