package users

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	model "github.com/deva-labs/univia/api/gin/src/modules/user/models"
	users "github.com/deva-labs/univia/api/gin/src/modules/user/services"
	"github.com/deva-labs/univia/common/config"
	"github.com/deva-labs/univia/common/utils"
	"github.com/deva-labs/univia/common/utils/cache"
	"github.com/deva-labs/univia/common/utils/types"

	"github.com/gin-gonic/gin"
)

// =================== GET USER ===================

// GetUser godoc
// @Summary      Get User Info
// @Description  Retrieve the user information by given token
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization header string true "Bearer <access_token>"
// @Success      200 {object} types.SuccessResponse{status=types.StatusOK,data=model.User} "User info"
// @Failure      401 {object} types.StatusUnauthorized "Unauthorized"
// @Failure      403 {object} types.StatusForbidden "Forbidden: insufficient permissions"
// @Failure      500 {object} types.StatusInternalError "Internal Server Error"
// @Router       /api/v1/auth/user-info [get]
func GetUser(c *gin.Context) {
	// Retrieve the user from the context (set by Authorization middleware)
	user, exists := c.Get("user")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized: user not found", nil)
		return
	}

	// Type assertion
	currentUser, ok := user.(*model.User)
	if !ok {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid user data", nil)
		return
	}

	// Redis cache key
	cacheKey := fmt.Sprintf("userInfo:%s", currentUser.Email)

	// Try get from Redis cache
	if cachedUser, err := cache.GetJSON[map[string]interface{}](config.Redis, cacheKey); err == nil && cachedUser != nil {
		utils.SendSuccessResponse(c, http.StatusOK, "Retrieved the profile of user successfully", cachedUser)
		return
	}

	// Fetch user info from service
	userInfo, err := users.GetUserInfo(currentUser.ID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get user", err)
		return
	}

	// Save to Redis
	_ = config.Redis.SetJSON(cacheKey, userInfo, 30*time.Minute)

	utils.SendSuccessResponse(c, http.StatusOK, "Retrieved the profile of user successfully", userInfo)
}

// =================== REGISTER ===================

// RegisterUser godoc
// @Summary      Register new user
// @Description  Create a new user with email/username/password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body model.User true "User data"
// @Success      201 {object} types.SuccessResponse{status=types.StatusOK,data=model.User} "Created user"
// @Failure      400 {object} types.StatusBadRequest "Invalid input"
// @Failure      500 {object} types.StatusInternalError "Internal Server Error"
// @Router       /api/v1/auth/register [post]
func RegisterUser(c *gin.Context) {
	var userData model.User
	if err := utils.BindJson(c, &userData); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	response, err := users.RegisterUser(userData)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "User has been created successfully", response)
}

// =================== LOGIN ===================

// LoginUser godoc
// @Summary      Login with email/username/phone
// @Description  Login using email, username or phone + password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body types.LoginRequest true "Login request"
// @Success      200 {object} types.SuccessLoginEmailResponse "Login success with verification code sent (data is email)"
// @Failure      400 {object} types.StatusBadRequest "Invalid input"
// @Failure      401 {object} types.StatusUnauthorized "Unauthorized"
// @Router       /api/v1/auth/login [post]
func LoginUser(c *gin.Context) {
	var request types.LoginRequest

	if err := utils.BindJson(c, &request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	email, status, err := users.LoginUser(request.Email, request.PhoneNumber, request.Password, request.Username)
	if err != nil {
		utils.SendErrorResponse(c, status, err.Error(), nil)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Verification code sent to your email. Please verify to process.", email)
}

// =================== LOGIN GOOGLE ===================

// LoginGoogle godoc
// @Summary      Login with Google
// @Description  Login using Google OAuth token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body types.GoogleLoginRequest true "Google login request"
// @Success      200 {object} types.SuccessResponse "Login success with tokens"
// @Failure      400 {object} types.StatusBadRequest "Invalid input"
// @Router       /api/v1/auth/login/google [post]
func LoginGoogle(c *gin.Context) {
	var request types.GoogleLoginRequest

	if err := utils.BindJson(c, &request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	response, err := users.LoginGoogle(request.Token)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Login successful", response)
}

// =================== LOGIN TWITTER ===================

// LoginTwitter godoc
// @Summary      Login with Twitter
// @Description  Login using Twitter account info
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body types.TwitterLoginRequest true "Twitter login request"
// @Success      200 {object} types.SuccessResponse "Login success with tokens"
// @Failure      400 {object} types.StatusBadRequest "Invalid input"
// @Router       /api/v1/auth/login/twitter [post]
func LoginTwitter(c *gin.Context) {
	// Read and log the raw request body
	body, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var request types.TwitterLoginRequest

	if err := utils.BindJson(c, &request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	response, err := users.LoginTwitter(
		request.Username,
		request.Email,
		request.Image,
		request.ProfileBackgroundImage,
		request.ProfileBackgroundColor,
		request.TwitterID,
	)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Login successful", response)
}

// =================== REFRESH TOKEN ===================

// RefreshAccessToken godoc
// @Summary      Refresh Access Token
// @Description  Refresh JWT token with refresh token + client id
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        x-rtoken-id header string true "Refresh Token (Bearer token format)"
// @Param        x-client-id header string true "Client ID"
// @Success      200 {object} types.SuccessResponse "New tokens"
// @Failure      401 {object} types.StatusUnauthorized "Unauthorized"
// @Failure      500 {object} types.StatusInternalError "Internal Server Error"
// @Router       /api/v1/auth/refresh-access-token [post]
func RefreshAccessToken(c *gin.Context) {
	authHeader := c.GetHeader("x-rtoken-id")
	clientID := c.GetHeader("x-client-id")

	if authHeader == "" || clientID == "" {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Missing token or client id on header", nil)
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid token header format", nil)
		return
	}
	refreshToken := tokenParts[1]

	response, err := users.RefreshAccessToken(refreshToken, clientID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to refresh token", err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Token refreshed successfully", response)
}

// =================== FORGOT PASSWORD ===================

// ForgotPassword godoc
// @Summary      Forgot Password
// @Description  Send a verification email for password reset
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body types.ForgotPasswordRequest true "Email"
// @Success      200 {object} types.SuccessResponse "Verification email sent"
// @Failure      400 {object} types.StatusBadRequest "Invalid input"
// @Failure      500 {object} types.StatusInternalError "Internal Server Error"
// @Router       /api/v1/auth/forgot-password [post]
func ForgotPassword(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var request types.ForgotPasswordRequest

	if err := utils.BindJson(c, &request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	status, err := users.ForgotPassword(request.Email)
	if err != nil {
		utils.SendErrorResponse(c, status, "Failed to send verification email", err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Verification email has been sent", nil)
}

// =================== RENEW PASSWORD ===================

// RenewPassword godoc
// @Summary      Renew Password
// @Description  Reset password without old password (via verification flow)
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body types.RenewPasswordRequest true "Renew password"
// @Success      200 {object} types.SuccessResponse "Password changed successfully"
// @Failure      400 {object} types.StatusBadRequest "Invalid input"
// @Failure      500 {object} types.StatusInternalError "Internal Server Error"
// @Router       /api/v1/auth/reset-password [post]
func RenewPassword(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var request types.RenewPasswordRequest

	if err := utils.BindJson(c, &request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	status, err := users.RenewPassword(request.NewPassword, strconv.Itoa(int(request.UserID)))
	if err != nil {
		utils.SendErrorResponse(c, status, "Failed to change password", err)
		return
	}

	utils.SendSuccessResponse(c, status, "Successfully changed password", nil)
}

// =================== CHANGE PASSWORD ===================

// ChangePassword godoc
// @Summary      Change Password
// @Description  Change password by providing old and new password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body types.ChangePasswordRequest true "Change password"
// @Success      200 {object} types.SuccessResponse "Password changed successfully"
// @Failure      400 {object} types.StatusBadRequest "Invalid input"
// @Failure      500 {object} types.StatusInternalError "Internal Server Error"
// @Router       /api/v1/auth/change-password [post]
func ChangePassword(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var request types.ChangePasswordRequest

	if err := utils.BindJson(c, &request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	status, err := users.ChangePassword(request.OldPassword, request.NewPassword, strconv.Itoa(int(request.UserID)))
	if err != nil {
		utils.SendErrorResponse(c, status, "Failed to change password", err)
		return
	}

	utils.SendSuccessResponse(c, status, "Successfully changed password", nil)
}

// =================== GET AVATAR ===================

// GetUserAvatar godoc
// @Summary      Get User Avatar
// @Description  Get avatar of a user by user_id
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization header string true "Bearer <access_token>"
// @Param        request body types.GetUserAvatarRequest true "User ID"
// @Success      200 {object} types.SuccessResponse "Avatar URL"
// @Failure      400 {object} types.StatusBadRequest "Invalid input"
// @Failure      401 {object} types.StatusUnauthorized "Unauthorized"
// @Failure      500 {object} types.StatusInternalError "Internal Server Error"
// @Router       /api/v1/users/avatar [post]
func GetUserAvatar(c *gin.Context) {
	var request types.GetUserAvatarRequest

	bindErr := utils.BindJson(c, &request)
	if bindErr != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid Input", bindErr)
		return
	}
	cacheKey := fmt.Sprintf("avatar_%d", request.UserID)
	// Try to cache
	if results, err := cache.GetJSON[map[string]interface{}](config.Redis, cacheKey); err == nil && results != nil {
		utils.SendSuccessResponse(c, http.StatusOK, "Retrieved the profile of user successfully", results)
		return
	} else if err != nil {
		fmt.Println(err)
	}
	// Continue
	avatarUser, err := users.GetUserImageByID(request.UserID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Message, nil)
		return
	}
	// Save to Redis
	_ = config.Redis.SetJSON(cacheKey, avatarUser, 30*time.Minute)
	utils.SendSuccessResponse(c, http.StatusOK, "Successfully get user avatar", avatarUser)
}
