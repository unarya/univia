package notifications

import (
	"fmt"
	"net/http"
	"time"

	"github.com/deva-labs/univia/internal/api/functions"
	notifications "github.com/deva-labs/univia/internal/api/modules/notification/services"
	"github.com/deva-labs/univia/internal/infrastructure/redis"
	_ "github.com/deva-labs/univia/pkg/types"
	"github.com/deva-labs/univia/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// List godoc
// @Summary      List all notifications of the current user
// @Description  Retrieve notifications with filters (pagination, seen/unseen, search)
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization header string true "Bearer <access_token>"
// @Param        request body types.ListNotificationRequest true "Notification filter options"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/v1/notifications [post]
func List(c *gin.Context) {
	var request struct {
		CurrentPage  int    `json:"current_page"`
		ItemsPerPage int    `json:"items_per_page"`
		OrderBy      string `json:"order_by"`
		SortBy       string `json:"sort_by"`
		SearchValue  string `json:"search_value"`
		IsSeen       bool   `json:"is_seen"`
		All          bool   `json:"all"`
	}

	if bindErr := utils.BindJson(c, &request); bindErr != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid Input", bindErr)
		return
	}

	currentUser, getUserErr := functions.GetCurrentUser(c)
	if getUserErr != nil {
		utils.SendErrorResponse(c, getUserErr.StatusCode, getUserErr.Message, nil)
		return
	}

	// Try cache
	cacheKey := fmt.Sprintf("notifications_by_user_%v_%d_%d_%s_%s_%s_%v_%v",
		currentUser.ID, request.CurrentPage, request.ItemsPerPage, request.OrderBy, request.SortBy, request.SearchValue, request.IsSeen, request.All)
	if results, err := redis.GetJSON[map[string]interface{}](redis.Redis, cacheKey); err == nil && results != nil {
		utils.SendSuccessResponse(c, http.StatusOK, "Successfully list notifications", results)
		return
	}
	response, err := notifications.GetNotificationsByUserID(
		currentUser.ID,
		request.CurrentPage,
		request.ItemsPerPage,
		request.OrderBy,
		request.SortBy,
		request.SearchValue,
		request.IsSeen,
		request.All,
	)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch notifications", err)
		return
	}
	_ = redis.Redis.SetJSON(cacheKey, response, 3*time.Hour)
	utils.SendSuccessResponse(c, http.StatusOK, "List all notifications successfully", response)
}

// UpdateSeen godoc
// @Summary      Mark a notification as seen
// @Description  Update notification status (seen) for the current user
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization header string true "Bearer <access_token>"
// @Param        request body types.UpdateSeenRequest true "Notification ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/v1/notifications/seen [put]
func UpdateSeen(c *gin.Context) {
	var request struct {
		NotificationID uuid.UUID `json:"notification_id"`
	}

	if bindErr := utils.BindJson(c, &request); bindErr != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid Input", bindErr)
		return
	}

	currentUser, getUserErr := functions.GetCurrentUser(c)
	if getUserErr != nil {
		utils.SendErrorResponse(c, getUserErr.StatusCode, getUserErr.Message, nil)
		return
	}

	serviceErr := notifications.UpdateIsSeen(request.NotificationID, currentUser.ID)
	if serviceErr != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update notification", serviceErr)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Notification updated successfully", nil)
}

// UpdateSeenWithUserID godoc
// @Summary      Mark all notifications as seen
// @Description  Update all notifications of the current user to seen
// @Tags         Notifications
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization header string true "Bearer <access_token>"
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/v1/notifications/seen/all [put]
func UpdateSeenWithUserID(c *gin.Context) {
	currentUser, getUserErr := functions.GetCurrentUser(c)
	if getUserErr != nil {
		utils.SendErrorResponse(c, getUserErr.StatusCode, getUserErr.Message, nil)
		return
	}

	serviceErr := notifications.UpdateIsSeenForAllNotificationByUserID(currentUser.ID)
	if serviceErr != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update all notifications", serviceErr)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "All notifications updated successfully", nil)
}
