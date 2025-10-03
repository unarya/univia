package posts

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"univia/src/config"
	"univia/src/functions"
	posts "univia/src/modules/post/services"
	model "univia/src/modules/user/models"
	"univia/src/utils"
	"univia/src/utils/cache"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreatePost godoc
// @Summary Create a new post
// @Description Create a post with content, categories, and media uploads
// @Tags Social Routes
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param content formData string true "Post content"
// @Param category_ids formData []string true "List of category UUIDs"
// @Param media formData file true "Media files (multiple allowed)"
// @Success 201 {object} map[string]interface{} "Post Created Successfully"
// @Failure 400 {object} types.StatusBadRequest "Bad Request"
// @Failure 401 {object} types.StatusUnauthorized "Unauthorized"
// @Failure 500 {object} types.StatusInternalError "Internal server error"
// @Router /api/v1/posts/create [post]
func CreatePost(c *gin.Context) {
	// Step 1: Parse form data
	content := strings.TrimSpace(c.PostForm("content"))
	categoryIDs, err := utils.ParseUUIDs(c.PostFormArray("category_ids"))

	if len(categoryIDs) == 0 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "At least one category must be selected", nil)
		return
	}

	// Step 2: Handle multiple file uploads
	form, err := c.MultipartForm()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Failed to parse form data", err)
		return
	}
	files := form.File["media"]

	// Step 3: Get current user from context
	user, exists := c.Get("user")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized: user not found", nil)
		return
	}
	currentUser, _ := user.(*model.User)

	// Step 4: Call service to create post
	result, serviceError := posts.CreatePost(content, categoryIDs, files, currentUser.ID)
	if serviceError != nil {
		utils.SendErrorResponse(c, serviceError.StatusCode, "Failed to create post", serviceError)
		return
	}

	// Step 5: Return success response
	utils.SendSuccessResponse(c, http.StatusCreated, "Post Created Successfully", result)
}

// ListAllPost godoc
// @Summary List all posts
// @Description Get all posts with pagination, search, and sorting
// @Tags Social Routes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body struct {
//     CurrentPage  int    `json:"current_page"`
//     ItemsPerPage int    `json:"items_per_page"`
//     OrderBy      string `json:"order_by"`
//     SortBy       string `json:"sort_by"`
//     SearchValue  string `json:"search_value"`
// } true "Pagination and filter"
// @Success 200 {object} map[string]interface{} "List all posts successfully"
// @Failure 400 {object} types.StatusBadRequest "Invalid input"
// @Failure 401 {object} types.StatusUnauthorized "Unauthorized"
// @Failure 500 {object} types.StatusInternalError "Internal server error"
// @Router /api/v1/posts [post]

func ListAllPost(c *gin.Context) {
	var request struct {
		CurrentPage  int    `json:"current_page"`
		ItemsPerPage int    `json:"items_per_page"`
		OrderBy      string `json:"order_by"`
		SortBy       string `json:"sort_by"`
		SearchValue  string `json:"search_value"`
	}
	if err := c.ShouldBind(&request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}
	currentUser, getUserErr := functions.GetCurrentUser(c)
	if getUserErr != nil {
		utils.SendErrorResponse(c, getUserErr.StatusCode, "An error occurred during execution", getUserErr)
		return
	}

	// Try cache
	cacheKey := fmt.Sprintf("listPost_%d_%d_%s_%s_%s:", request.CurrentPage, request.ItemsPerPage, request.OrderBy, request.SortBy, request.SearchValue)
	if results, err := cache.GetJSON[map[string]interface{}](config.Redis, cacheKey); err == nil && results != nil {
		utils.SendSuccessResponse(c, http.StatusOK, "List all posts successfully", results)
		return
	}
	response, err := posts.List(request.CurrentPage, request.ItemsPerPage, request.OrderBy, request.SortBy, request.SearchValue, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = config.Redis.SetJSON(cacheKey, response, 3*time.Minute)
	utils.SendSuccessResponse(c, http.StatusOK, "List all posts successfully", response)
}

// GetDetailsPost godoc
// @Summary Get post details
// @Description Get detailed information of a single post by ID
// @Tags Social Routes
// @Produce json
// @Security BearerAuth
// @Param id query string true "Post ID"
// @Success 200 {object} map[string]interface{} "Successfully get details of this post"
// @Failure 400 {object} types.StatusBadRequest "ID is required"
// @Failure 401 {object} types.StatusUnauthorized "Unauthorized"
// @Failure 500 {object} types.StatusInternalError "Internal server error"
// @Router /api/v1/posts [get]
func GetDetailsPost(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Id is required", nil)
		return
	}
	// Try cache
	cacheKey := fmt.Sprintf("detailPost_%s", id)
	if results, err := cache.GetJSON[map[string]interface{}](config.Redis, cacheKey); err == nil && results != nil {
		utils.SendSuccessResponse(c, http.StatusOK, "Successfully get details of this post", results)
		return
	}
	response, err := posts.GetDetails(id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get details", err)
		return
	}
	_ = config.Redis.SetJSON(cacheKey, response, 3*time.Minute)
	utils.SendSuccessResponse(c, http.StatusOK, "Successfully get details successfully", response)
}

// UpdatePost godoc
// @Summary Update a post
// @Description Update post content, categories, and media files
// @Tags Social Routes
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id formData string true "Post ID"
// @Param content formData string true "Post content"
// @Param category_ids formData []string true "List of category UUIDs"
// @Param media formData file false "Media files (multiple allowed)"
// @Success 200 {object} map[string]interface{} "Updated post successfully"
// @Failure 400 {object} types.StatusBadRequest "Bad Request"
// @Failure 401 {object} types.StatusUnauthorized "Unauthorized"
// @Failure 500 {object} types.StatusInternalError "Internal server error"
// @Router /api/v1/posts [put]
func UpdatePost(c *gin.Context) {
	// Step 1: Parse form data
	content := strings.TrimSpace(c.PostForm("content"))
	postID, _ := uuid.Parse(c.PostForm("id"))
	categoryIDs, err := utils.ParseUUIDs(c.PostFormArray("category_ids"))

	if len(categoryIDs) == 0 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "At least one category must be selected", nil)
		return
	}

	// Step 2: Handle multiple file uploads
	form, err := c.MultipartForm()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Failed to parse form data", err)
		return
	}
	// Handle Receive multiple media
	files := form.File["media"]

	// Step 3: Get user
	// Retrieve the user from the context (set by Authorization middleware)
	currentUser, getUserErr := functions.GetCurrentUser(c)
	if getUserErr != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized: user not found", nil)
		return
	}
	postInfo := posts.PostInfo{
		UserID:      currentUser.ID,
		PostID:      postID,
		Content:     content,
		CategoryIDs: categoryIDs,
		Media:       files,
	}
	serviceError := posts.EditPostByUserID(postInfo)
	if serviceError != nil {
		utils.SendErrorResponse(c, serviceError.StatusCode, "Failed to update post", serviceError)
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Post Updated Successfully", postInfo)
}
