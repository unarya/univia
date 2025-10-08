package functions

import (
	"net/http"

	"github.com/unarya/univia/internal/api/modules/post/models"
	"github.com/unarya/univia/internal/infrastructure/mysql"
	"github.com/unarya/univia/pkg/utils"

	"github.com/google/uuid"
)

// CheckIsLiked is a function need userID and postID and return bool
func CheckIsLiked(userID, postID uuid.UUID) (bool, *utils.ServiceError) {
	db := mysql.DB
	var liked bool
	err := db.Model(&posts.PostLike{}).
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
