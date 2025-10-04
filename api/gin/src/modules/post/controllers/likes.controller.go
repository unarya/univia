package posts

import (
	"fmt"
	"net/http"
	"time"

	"github.com/deva-labs/univia/api/gin/src/config"
	"github.com/deva-labs/univia/api/gin/src/functions"
	posts "github.com/deva-labs/univia/api/gin/src/modules/post/services"
	"github.com/deva-labs/univia/api/gin/src/utils"
	"github.com/deva-labs/univia/api/gin/src/utils/cache"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Like godoc
// @Summary Like the post
// @Description Like the post
// @Tags Social Routes
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization header string true "Bearer <access_token>"
// @Param        request body types.LikeRequest true "Post ID Required"
// @Success 200 {object} types.SuccessLikeAPostResponse "Successfully Like A Post"
// @Failure 400 {object} types.StatusBadRequest "Bad Request"
// @Failure 500 {object} types.StatusInternalError "Internal server error"
// @Router /api/v1/likes [get]
func Like(c *gin.Context) {
	var request struct {
		PostID uuid.UUID `json:"post_id"`
	}
	if bindErr := utils.BindJson(c, &request); bindErr != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid Input", bindErr)
		return
	}

	// Get current user from context
	currentUser, getUserErr := functions.GetCurrentUser(c)
	if getUserErr != nil {
		utils.SendErrorResponse(c, getUserErr.StatusCode, "An error occurred during calculation", getUserErr)
		return
	}

	// Try cache
	cacheKey := fmt.Sprintf("likes_%s", request.PostID.String())
	if results, err := cache.GetJSON[int64](config.Redis, cacheKey); err == nil && results != nil {
		utils.SendSuccessResponse(c, http.StatusOK, "Successfully liked post", gin.H{
			"totalLikes": results,
		})
		return
	}
	totalLikes, err := posts.Like(currentUser.ID, request.PostID)
	if err != nil {
		utils.SendErrorResponse(c, err.StatusCode, "An error occurred during calculation", err)
		return
	}
	_ = config.Redis.SetJSON(cacheKey, totalLikes, 3*time.Minute)
	utils.SendSuccessResponse(c, http.StatusOK, "Successfully liked post", gin.H{"totalLikes": totalLikes})
}

// DisLike godoc
// @Summary DisLike the post
// @Description DisLike the post
// @Tags Social Routes
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization header string true "Bearer <access_token>"
// @Param        request body types.LikeRequest true "Post ID Required"
// @Success 200 {object} types.SuccessDisLikeAPostResponse "Successfully Like A Post"
// @Failure 400 {object} types.StatusBadRequest "Bad Request"
// @Failure 500 {object} types.StatusInternalError "Internal server error"
// @Router /api/v1/likes/undo [get]
func DisLike(c *gin.Context) {
	var request struct {
		PostID uuid.UUID `json:"post_id"`
	}
	if bindErr := utils.BindJson(c, &request); bindErr != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid Input", bindErr)
		return
	}
	currentUser, getUserErr := functions.GetCurrentUser(c)
	if getUserErr != nil {
		c.JSON(getUserErr.StatusCode, gin.H{"error": getUserErr.Message})
		return
	}

	// Try cache
	cacheKey := fmt.Sprintf("likes_%s", request.PostID.String())
	if results, err := cache.GetJSON[int64](config.Redis, cacheKey); err == nil && results != nil {
		utils.SendSuccessResponse(c, http.StatusOK, "Successfully disliked post", gin.H{
			"totalLikes": results,
		})
		return
	}
	totalLikes, err := posts.DisLike(currentUser.ID, request.PostID)
	if err != nil {
		utils.SendErrorResponse(c, err.StatusCode, "An error occurred during calculation", err)
		return
	}

	// Set to redis
	_ = config.Redis.SetJSON(cacheKey, totalLikes, 3*time.Minute)
	utils.SendSuccessResponse(c, http.StatusOK, "Successfully disliked post", gin.H{
		"totalLikes": totalLikes,
	})
}
