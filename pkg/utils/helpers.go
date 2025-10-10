package utils

import (
	"encoding/binary"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	sessions "github.com/unarya/univia/internal/api/modules/session/model"
	users "github.com/unarya/univia/internal/api/modules/user/models"
	"github.com/unarya/univia/internal/infrastructure/redis"
	"github.com/unarya/univia/pkg/types"
	"gorm.io/gorm"
)

func Paginate(total int64, page, perPage int) (map[string]interface{}, error) {
	// Avoid division by zero
	if perPage <= 0 {
		perPage = 1
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	// Determine next and previous pages
	var nextPage, prevPage *int
	if page < totalPages {
		next := page + 1
		nextPage = &next
	}
	if page > 1 {
		prev := page - 1
		prevPage = &prev
	}

	// Construct and return pagination data
	return map[string]interface{}{
		"current_page":   page,
		"items_per_page": perPage,
		"next_page":      nextPage,
		"previous_page":  prevPage,
		"total_count":    total,
		"total_pages":    totalPages,
	}, nil
}

// ConvertStringToInt64 is the function will receive a string and return int64
func ConvertStringToInt64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

// ConvertInt64ToString is the function will receive an int64 and return string
func ConvertInt64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// ConvertInt64ToUUID is the function will receive an int64 and return uuid.UUID format
func ConvertInt64ToUUID(i int64) uuid.UUID {
	b := make([]byte, 16)
	binary.BigEndian.PutUint64(b[8:], uint64(i))
	return uuid.Must(uuid.FromBytes(b))
}

// ConvertStringToUuid is the function to convert string to uuid.UUID format
func ConvertStringToUuid(str string) uuid.UUID {
	uuid, err := uuid.Parse(str)
	if err != nil {
		panic(err)
	}
	return uuid
}

// ServiceError to define return exception for system
type ServiceError struct {
	StatusCode int
	Message    string
}

func (e *ServiceError) Error() string {
	return e.Message
}

// CalculateOffsetStruct is the struct to define return result for calculate service
type CalculateOffsetStruct struct {
	CurrentPage  int
	ItemsPerPage int
	OrderBy      string
	SortBy       string
	Offset       int
}

// CalculateOffset is the function to calculate offset for list service
func CalculateOffset(currentPage, itemsPerPage int, sortBy, orderBy string) CalculateOffsetStruct {
	if orderBy == "" {
		orderBy = "created_at"
	}
	if sortBy != "asc" && sortBy != "desc" {
		sortBy = "desc"
	}

	// Calculate offset for pagination
	offset := (currentPage - 1) * itemsPerPage
	if offset < 0 {
		offset = 0
	}

	return CalculateOffsetStruct{
		CurrentPage:  currentPage,
		ItemsPerPage: itemsPerPage,
		OrderBy:      orderBy,
		SortBy:       sortBy,
		Offset:       offset,
	}
}

// BindJson is a function to bind the json request
func BindJson(c *gin.Context, request interface{}) error {
	if err := c.ShouldBindJSON(request); err != nil {
		return err
	}
	return nil
}

func ParseUUIDs(strs []string) ([]uuid.UUID, error) {
	var uuids []uuid.UUID
	for _, s := range strs {
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, err
		}
		uuids = append(uuids, id)
	}
	return uuids, nil
}

// SendSuccessResponse and SendErrorResponse Helper functions for consistent responses
func SendSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, types.SuccessResponse{
		Status: types.StatusOK{
			Code:    statusCode,
			Message: message,
		},
		Data: data,
	})
}

func SendErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	response := types.ErrorResponse{
		Status: types.StatusBadRequest{
			Code:    statusCode,
			Message: message,
		},
	}
	if err != nil {
		response.Error = err.Error()
	}
	c.JSON(statusCode, response)
}

func NowPtr() *time.Time {
	now := time.Now()
	return &now
}

func SetSessionToRedis(db *gorm.DB, session sessions.UserSession, user users.User, meta types.SessionMetadata) error {
	// Save redis for signal handshaking
	cacheKey := fmt.Sprintf("session:%s", session.SessionID)
	cacheValue := map[string]interface{}{
		"user_id":     user.ID,
		"email":       user.Email,
		"username":    user.Username,
		"session_id":  session.SessionID,
		"ip":          meta.IP,
		"user_agent":  meta.UserAgent,
		"created_at":  session.CreatedAt,
		"last_active": session.LastActive,
	}
	err := redis.Redis.SetJSON(cacheKey, cacheValue, 12*time.Hour)

	// Update session status to active
	session.Status = "active"
	session.LastActive = NowPtr()
	if err := db.Save(&session).Error; err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

func SetSessionToRedisByUserID(c *gin.Context, db *gorm.DB, user users.User) (uuid.UUID, error) {
	var session sessions.UserSession
	err := db.Where("user_id = ?", user.ID).
		First(&session).Error
	if err != nil {
		return uuid.Nil, err
	}
	// Store session in Redis for signaling handshake
	meta, _ := GetSessionMetadata(c)
	if err := SetSessionToRedis(db, session, user, meta); err != nil {
		return uuid.Nil, err
	}
	return session.SessionID, nil
}

// SetHttpOnlyCookieForSession is a function to set http only cookie on a device
func SetHttpOnlyCookieForSession(c *gin.Context, sessionID uuid.UUID) {
	env := os.Getenv("NODE_ENV")
	isProd := env == "production"

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteLaxMode,
	}

	if isProd {
		cookie.SameSite = http.SameSiteNoneMode
	}

	http.SetCookie(c.Writer, cookie)
}

// SetHttpOnlyCookieForUser is a function to set http only cookie on a device
func SetHttpOnlyCookieForUser(c *gin.Context, userID string) {
	env := os.Getenv("NODE_ENV")
	isProd := env == "production"

	cookie := &http.Cookie{
		Name:     "user_id",
		Value:    userID,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteLaxMode,
	}

	if isProd {
		cookie.SameSite = http.SameSiteNoneMode
	}

	http.SetCookie(c.Writer, cookie)
}

func GetHttpOnlyCookieForSession(c *gin.Context) string {
	results, _ := c.Cookie("session_id")
	return results
}

func GetHttpOnlyCookieForUser(c *gin.Context, cookieName string) string {
	results, _ := c.Cookie(cookieName)
	return results
}

func GetSessionMetadata(c *gin.Context) (types.SessionMetadata, error) {
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()
	meta := types.SessionMetadata{IP: ip, UserAgent: userAgent}
	return meta, nil
}
