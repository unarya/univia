package functions

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/unarya/univia/internal/api/modules/post/models"
	"github.com/unarya/univia/internal/infrastructure/mysql"
	"github.com/unarya/univia/pkg/utils"

	"github.com/google/uuid"
)

// CreatePost is the function will create post with userID and content
func CreatePost(content string, userID uuid.UUID) (postID uuid.UUID, errService *utils.ServiceError) {
	db := mysql.DB
	post := posts.Post{UserID: userID, Content: content}
	if err := db.Create(&post).Error; err != nil {
		return uuid.Nil, &utils.ServiceError{StatusCode: http.StatusInternalServerError, Message: "Failed to create post"}
	}
	return post.ID, nil
}

// SaveCategoriesToPost is the function will save post categories to mysql
func SaveCategoriesToPost(categoryIDs []uuid.UUID, postID uuid.UUID) *utils.ServiceError {
	db := mysql.DB
	var postCategories []posts.PostCategory

	for _, categoryID := range categoryIDs {
		postCategories = append(postCategories, posts.PostCategory{
			PostID:     postID,
			CategoryID: categoryID,
		})
	}

	if len(postCategories) > 0 {
		if err := db.Create(&postCategories).Error; err != nil {
			return &utils.ServiceError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to associate categories",
			}
		}
	}
	return nil
}

// SaveMediaRecords is the function will save from media path into mysql
func SaveMediaRecords(savedMedia []posts.Media, postID uuid.UUID) *utils.ServiceError {
	db := mysql.DB
	for _, media := range savedMedia {
		media.PostID = postID
		if err := db.Create(&media).Error; err != nil {
			return &utils.ServiceError{StatusCode: http.StatusInternalServerError, Message: "Failed to save media"}
		}
	}
	return nil
}

// DeletePostRecord is the function to delete post records from mysql
func DeletePostRecord(postID uint) *utils.ServiceError {
	db := mysql.DB
	// Attempt to delete the post record with the given ID
	if err := db.Delete(&posts.Post{}, postID).Error; err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete post",
		}
	}
	return nil
}

// DeleteMediaRecords is the function to delete media records following to the post
func DeleteMediaRecords(postID uuid.UUID) *utils.ServiceError {
	db := mysql.DB

	if err := db.Where("post_id = ?", postID).Delete(&posts.Media{}).Error; err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to remove old media records",
		}
	}
	return nil
}

// DeleteCategoryRecords is the function to delete categories following to the post
func DeleteCategoryRecords(postID uuid.UUID) *utils.ServiceError {
	db := mysql.DB
	if err := db.Where("post_id = ?", postID).Delete(&posts.PostCategory{}).Error; err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to remove old category associations",
		}
	}
	return nil
}

// UpdatePostContent is the function to update the content of post by given postID
func UpdatePostContent(content string, postID uuid.UUID) *utils.ServiceError {
	db := mysql.DB
	if err := db.Model(&posts.Post{}).Where("id = ?", postID).Updates(posts.Post{
		Content: content,
	}).Error; err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update post",
		}
	}
	return nil
}

// CheckPostExits is the function will check post was valid on mysql or not
func CheckPostExits(postID uuid.UUID) *utils.ServiceError {
	db := mysql.DB
	// Check if Post Exists
	var exists bool
	if err := db.Model(&posts.Post{}).
		Select("count(*) > 0").
		Where("id = ?", postID).
		Find(&exists).Error; err != nil {

		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check post existence",
		}
	}

	// If the post does not exist, return an error
	if !exists {
		return &utils.ServiceError{
			StatusCode: http.StatusNotFound,
			Message:    "Post not found",
		}
	}
	return nil
}

// SelectPosts is the function will execute the sql queries with given parameters and return rows
func SelectPosts(searchValue, orderBy, sortBy string, offset, limit int) (*sql.Rows, *utils.ServiceError) {
	rows, err := mysql.DB.Table("posts").
		Select(`
			posts.id, posts.content, posts.created_at, posts.updated_at,
			users.id AS user_id, users.username AS username, profiles.profile_pic,
			GROUP_CONCAT(DISTINCT categories.id ORDER BY categories.id ASC SEPARATOR ',') AS category_ids,
			GROUP_CONCAT(DISTINCT categories.name ORDER BY categories.id ASC SEPARATOR ',') AS category_names,
			media.id AS media_id, media.path, media.type, media.status,
			COUNT(DISTINCT comments.id) AS comment_count,
			COUNT(DISTINCT post_likes.id) AS likes_count,
			COUNT(DISTINCT post_shares.id) AS shares_count,
			COUNT(posts.id) OVER() AS total_count
		`).
		Joins(`
			LEFT JOIN post_categories ON post_categories.post_id = posts.id
			LEFT JOIN categories ON categories.id = post_categories.category_id
			LEFT JOIN media ON media.post_id = posts.id
			LEFT JOIN users ON users.id = posts.user_id
			LEFT JOIN profiles ON profiles.user_id = users.id
			LEFT JOIN comments ON comments.post_id = posts.id
			LEFT JOIN post_likes ON post_likes.post_id = posts.id
			LEFT JOIN post_shares ON post_shares.post_id = posts.id
		`).
		Where("LOWER(posts.content) LIKE LOWER(?)", "%"+searchValue+"%").
		Group("posts.id, users.id, users.username, profiles.profile_pic, media.id").
		Order(fmt.Sprintf("posts.%s %s", orderBy, sortBy)).
		Offset(offset).
		Limit(limit).
		Rows()
	if err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}
	return rows, nil
}

// ListNotifications returns a row query for list of notifications
func ListNotifications(searchValue, orderBy, sortBy string, offset, limit int, isSeen bool, receiverID uuid.UUID, all bool) (*sql.Rows, *utils.ServiceError) {
	query := mysql.DB.Table("notifications").
		Select(`*,
            COUNT(notifications.id) OVER() AS total_count
        `).
		Where("LOWER(notifications.message) LIKE LOWER(?) AND receiver_id = ?", "%"+searchValue+"%", receiverID)

	// Only add is_seen filter if all=false
	if !all {
		query = query.Where("is_seen = ?", isSeen)
	}

	rows, err := query.
		Order(fmt.Sprintf("notifications.%s %s", orderBy, sortBy)).
		Offset(offset).
		Limit(limit).
		Rows()

	if err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}
	return rows, nil
}
