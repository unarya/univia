package functions

import (
	"net/http"

	model "github.com/deva-labs/univia/api/gin/src/modules/user/models"
	"github.com/deva-labs/univia/common/utils"

	"github.com/gin-gonic/gin"
)

// GetCurrentUser retrieves the current user from the context
func GetCurrentUser(c *gin.Context) (*model.User, *utils.ServiceError) {
	user, exists := c.Get("user")
	if !exists {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized: user not found",
		}
	}

	// Type assertion (ensure user is of type *model.User)
	currentUser, ok := user.(*model.User)
	if !ok {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Type assertion failed: user is not of type *model.User",
		}
	}

	return currentUser, nil
}
