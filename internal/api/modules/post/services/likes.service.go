package posts

import (
	"fmt"
	"log"
	"net/http"

	"github.com/unarya/univia/internal/api/functions"
	"github.com/unarya/univia/internal/api/modules/notification/services"
	posts "github.com/unarya/univia/internal/api/modules/post/models"
	Users "github.com/unarya/univia/internal/api/modules/user/models"
	"github.com/unarya/univia/internal/infrastructure/mysql"
	"github.com/unarya/univia/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Like is a service function to process the like action.
func Like(userID, postID uuid.UUID) (int64, *utils.ServiceError) {
	db := mysql.DB
	var counts int64

	// 1. Check if the user already liked the post
	liked, err := functions.CheckIsLiked(userID, postID)
	if err != nil {
		return counts, &utils.ServiceError{Message: err.Message, StatusCode: err.StatusCode}
	}

	if liked {
		return counts, &utils.ServiceError{
			Message:    "You have already liked this post",
			StatusCode: http.StatusBadRequest,
		}
	}

	// 2. Create a new like for the post in a single query
	newLike := posts.PostLike{
		UserID: userID,
		PostID: postID,
	}

	// Insert the like, while also counting the total likes in one query.
	// Using db.Transaction to ensure atomicity of the operations.
	if err := db.Transaction(func(tx *gorm.DB) error {
		// Insert the like
		if err := tx.Create(&newLike).Error; err != nil {
			return err
		}

		// Count the total likes after inserting the new like
		if err := tx.Model(&posts.PostLike{}).
			Where("post_id = ?", postID).
			Count(&counts).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return counts, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to like the post",
		}
	}

	// Send Notification

	// 1. Get username of the user who liked the post
	var username string
	usernameErr := db.Model(&Users.User{}).
		Select("username").
		Where("id = ?", userID).
		Scan(&username).Error
	if usernameErr != nil {
		log.Printf("Failed to get username for userID %d: %v", userID, usernameErr)
		return counts, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("Failed to get username for userID %d: %v", userID, usernameErr),
		}
	}

	// 2. Get owner of the post
	var postOwner uuid.UUID
	selectOwnerErr := db.Model(&posts.Post{}).
		Select("user_id").
		Where("id = ?", postID).
		Scan(&postOwner).Error
	if selectOwnerErr != nil {
		log.Printf("Failed to get post owner for postID %d: %v", postID, selectOwnerErr)
		return counts, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("Failed to get post owner for postID %d: %v", postID, selectOwnerErr),
		}
	}

	if postOwner != userID {
		// 3. Send notification to the post owner
		message := fmt.Sprintf("%s just liked your post", username)
		noti_type := "personal_post"
		sendNotiErr := notifications.NotificationHandler(userID, postOwner, message, noti_type)
		if sendNotiErr != nil {
			log.Printf("Failed to send notification: %v", sendNotiErr)
		}
	}

	return counts, nil
}

// DisLike is a service function to delete a like from the PostLike table and update the like count.
func DisLike(userID, postID uuid.UUID) (int64, *utils.ServiceError) {
	db := mysql.DB
	var counts int64

	// 1. Check if the user has already liked the post
	liked, err := functions.CheckIsLiked(userID, postID)
	if err != nil {
		return counts, &utils.ServiceError{Message: err.Message, StatusCode: err.StatusCode}
	}

	if !liked {
		return counts, &utils.ServiceError{
			Message:    "You have already disliked this post",
			StatusCode: http.StatusBadRequest,
		}
	}

	// 2. Delete the like record for the post in a single query
	if err := db.Transaction(func(tx *gorm.DB) error {
		// Delete the like
		if err := tx.Where("post_id = ? AND user_id = ?", postID, userID).
			Delete(&posts.PostLike{}).Error; err != nil {
			return err
		}

		// Count the total likes after deleting the like
		if err := tx.Model(&posts.PostLike{}).
			Where("post_id = ?", postID).
			Count(&counts).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return counts, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete the liked post",
		}
	}

	return counts, nil
}
