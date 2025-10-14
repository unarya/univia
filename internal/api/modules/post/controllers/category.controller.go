package posts

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	posts "github.com/unarya/univia/internal/api/modules/post/services"
	"github.com/unarya/univia/internal/infrastructure/redis"
	_ "github.com/unarya/univia/pkg/types"
	"github.com/unarya/univia/pkg/utils"
)

// ListCategories godoc
// @Summary List categories
// @Description Retrieve all categories in the system
// @Tags Social Routes
// @Produce      json
// @Success 200 {object} types.SuccessListCategoriesResponse "Successfully List Categories"
// @Failure 400 {object} types.StatusBadRequest "Bad Request"
// @Failure 500 {object} types.StatusInternalError "Internal orchestrator error"
// @Router /api/v1/posts/categories [get]
func ListCategories(c *gin.Context) {
	// Cache
	cacheKey := "listCategories"
	if results, err := redis.GetJSON[[]map[string]interface{}](redis.Redis, cacheKey); err != nil && results != nil {
		utils.SendSuccessResponse(c, http.StatusOK, "Categories List Successfully", results)
		return
	}
	results, err := posts.ListAllCategories()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Categories List Error", err)
		return
	}
	_ = redis.Redis.SetJSON(cacheKey, results, 12*time.Hour)
	utils.SendSuccessResponse(c, http.StatusOK, "Categories List Successfully", results)
}
