package roles

import (
	"fmt"
	"net/http"

	"github.com/deva-labs/univia/api/gin/src/config"
	"github.com/deva-labs/univia/api/gin/src/modules/role/services"
	"github.com/deva-labs/univia/api/gin/src/utils"
	"github.com/deva-labs/univia/api/gin/src/utils/cache"
	"github.com/deva-labs/univia/api/gin/src/utils/types"

	"github.com/gin-gonic/gin"
)

// CreateRole godoc
// @Summary Create a new role
// @Description Admin creates a new role with a unique name
// @Tags Roles
// @Accept       json
// @Produce      json
// @Param request body types.CreateRoleRequest true "Role Name"
// @Success 201 {object} types.SuccessCreateRoleResponse "Role Created Successfully"
// @Failure 400 {object} types.StatusBadRequest "Invalid Input"
// @Failure 500 {object} types.StatusInternalError "Internal server error"
// @Router /api/v1/roles [post]
func CreateRole(c *gin.Context) {
	var request types.CreateRoleRequest

	// Parse JSON input
	if err := utils.BindJson(c, &request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	results, err := roles.CreateRoleByAdmin(request.Name)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error: "+err.Error(), nil)
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Role Created Successfully", results)
}

// ListRoles godoc
// @Summary List all roles
// @Description Retrieve all roles in the system
// @Tags Roles
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization header string true "Bearer <access_token>"
// @Success 200 {object} types.SuccessListRolesResponse "Successfully List All Roles"
// @Failure 500 {object} types.StatusInternalError "Internal server error"
// @Router /api/v1/roles [get]
func ListRoles(c *gin.Context) {
	// Cache
	cacheKey := "listRoles"
	if results, err := cache.GetJSON[[]map[string]interface{}](config.Redis, cacheKey); err != nil && results != nil {
		utils.SendSuccessResponse(c, http.StatusOK, "Success", nil)
		return
	} else if err != nil {
		fmt.Printf("Cache miss: %s", err)
	}
	// Continue
	results, err := roles.ListAllRoles()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error ", err)
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Success", results)
}
