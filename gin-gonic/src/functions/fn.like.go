package functions

import (
	"net/http"
	"univia/src/config"
	"univia/src/modules/post/models"
	"univia/src/utils"

	"github.com/google/uuid"
)

// CheckIsLiked is a function need userID and postID and return bool
func CheckIsLiked(userID, postID uuid.UUID) (bool, *utils.ServiceError) {
	db := config.DB
	var liked bool
	err := db.Model(&models.PostLike{}).
		Select("count(*) > 0").
		Where("user_id = ? AND post_id = ?", userID, postID).
		Find(&liked).Error

	if err != nil {
		return true, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Database error while checking like status",
		}
	}

	return liked, nil
}
