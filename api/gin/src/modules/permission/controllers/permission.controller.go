package permissions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/deva-labs/univia/api/gin/src/modules/permission/services"
	"github.com/deva-labs/univia/common/config"
	"github.com/deva-labs/univia/common/utils"
	"github.com/deva-labs/univia/common/utils/cache"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreatePermission godoc
// @Summary      Create a new permission
// @Description  Admin creates a new permission
// @Tags         Permissions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization header string true "Bearer <access_token>"
// @Param        request body  types.CreatePermissionRequest true "Permission name"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/v1/permissions/create [post]
func CreatePermission(c *gin.Context) {
	var request struct {
		PermissionName string `json:"name"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid Input", err)
		return
	}
	isCreated, err := permissions.CreatePermission(request.PermissionName)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Bad Request", err)
		return
	}
	if !isCreated {
		utils.SendErrorResponse(c, http.StatusBadRequest, "An error occurred during creating permission", nil)
		return
	}
	utils.SendSuccessResponse(c, http.StatusCreated, "Successfully created permission", nil)
}

// ListPermissions godoc
// @Summary      List all permissions
// @Description  Get all permissions in the system
// @Tags         Permissions
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization header string true "Bearer <access_token>"
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/v1/permissions/list [post]
func ListPermissions(c *gin.Context) {
	// Try cache
	cacheKey := fmt.Sprintf("listPermissions")
	if results, err := cache.GetJSON[[]map[string]interface{}](config.Redis, cacheKey); err == nil && results != nil {
		utils.SendSuccessResponse(c, http.StatusOK, "Successfully list permissions", results)
		return
	}
	results, err := permissions.ListAllPermissions()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Bad Request", err)
		return
	}
	_ = config.Redis.SetJSON(cacheKey, results, 12*time.Hour)
	utils.SendSuccessResponse(c, http.StatusOK, "Successfully list permissions", results)
}

// AssignPermissionsToRole godoc
// @Summary      Assign permissions to a role
// @Description  Grant multiple permissions to a specific role
// @Tags         Permissions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization header string true "Bearer <access_token>"
// @Param        request body types.AssignPermissionRequest true "Role ID and list of Permission IDs"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/v1/permissions/assign [post]
func AssignPermissionsToRole(c *gin.Context) {
	var request struct {
		RoleID        uuid.UUID   `json:"role_id"`
		PermissionIDs []uuid.UUID `json:"permission_ids"`
	}

	// Validate request payload
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid Input", err)
		return
	}

	// Call service function
	results, err := permissions.AddPermissionsToRole(request.RoleID, request.PermissionIDs)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Bad Request", err)
		return
	}

	// Return response
	utils.SendSuccessResponse(c, http.StatusOK, "Successfully assign permissions to role", results)
}
