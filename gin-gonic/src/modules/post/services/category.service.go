package posts

import (
	"errors"
	"github.com/deva-labs/univia-api/api/gin-gonic/src/config"
	posts "github.com/deva-labs/univia-api/api/gin-gonic/src/modules/post/models"
)

// ListAllCategories retrieves all categories from the database
func ListAllCategories() ([]map[string]interface{}, error) {
	db := config.DB
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
