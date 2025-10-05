package posts

import (
	"errors"

	posts "github.com/deva-labs/univia/internal/api/modules/post/models"
	"github.com/deva-labs/univia/internal/infrastructure/mysql"
)

// ListAllCategories retrieves all categories from the mysql
func ListAllCategories() ([]map[string]interface{}, error) {
	db := mysql.DB
	var categories []posts.Category

	// Fetch categories from DB
	if err := db.Find(&categories).Error; err != nil {
		return nil, errors.New("failed to list all categories")
	}

	// Convert categories to map format
	var result []map[string]interface{}
	for _, category := range categories {
		result = append(result, map[string]interface{}{
			"id":   category.ID,
			"name": category.Name,
		})
	}

	return result, nil
}
