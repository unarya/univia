package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/deva-labs/univia/api/gin/src/modules/key_token/access_token/services"
	PermissionServices "github.com/deva-labs/univia/api/gin/src/modules/permission/services"
	"github.com/deva-labs/univia/api/gin/src/modules/user/models"
	"github.com/deva-labs/univia/common/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "authorization token missing", nil)
			c.Abort()
			return
		}

		token, err := SplitToken(authHeader)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Not correct form of header", err)
			c.Abort()
			return
		}

		userInfo, err := access_token.VerifyToken(token)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid Token", err)
			c.Abort()
			return
		}

		// Đặt vào context
		c.Set("user", userInfo)
		c.Next()
	}
}

// Authorization ensures the user has the required permission
func Authorization(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve user from context
		userInterface, exists := c.Get("user")
		if !exists {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		// Type assertion
		user, ok := userInterface.(*users.User)
		if !ok {
			utils.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error", nil)
			c.Abort()
			return
		}

		// Check role permissions
		if !PermissionServices.CheckPermission(user.RoleID, requiredPermission) {
			utils.SendErrorResponse(c, http.StatusForbidden, "Forbidden: insufficient permissions", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func SplitToken(token string) (string, error) {
	splitToken := strings.Split(token, " ")
	// Check if the token is in the correct "Bearer <token>" format
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		return "", errors.New("token format error")
	}
	return splitToken[1], nil
}
