package types

import "github.com/google/uuid"

// StatusResponse represents the status portion of all responses
// ================== BASE ==================

type StatusOK struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"Success"`
}

type StatusCreated struct {
	Code    int    `json:"code" example:"201"`
	Message string `json:"message" example:"Created"`
}

// ErrorResponse represents an error API response (generic)
type ErrorResponse struct {
	Status StatusBadRequest `json:"status"`
	Error  string           `json:"error,omitempty"`
}
type StatusBadRequest struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad Request"`
}

type StatusUnauthorized struct {
	Code    int    `json:"code" example:"401"`
	Message string `json:"message" example:"Unauthorized"`
}

type StatusForbidden struct {
	Code    int    `json:"code" example:"403"`
	Message string `json:"message" example:"Forbidden"`
}

type StatusInternalError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Internal Server Error"`
}

// ================== SUCCESS RESPONSES ==================

// SuccessResponse Generic 200 OK with data
type SuccessResponse struct {
	Status StatusOK `json:"status"`
	Data   any      `json:"data,omitempty"`
}

// CreatedResponse Generic 201 Created with data
type CreatedResponse struct {
	Status StatusCreated `json:"status"`
	Data   any           `json:"data,omitempty"`
}

// ================== ERROR RESPONSES ==================

// Error400Response 400 Bad Request
type Error400Response struct {
	Status StatusBadRequest `json:"status"`
	Error  string           `json:"error" example:"Invalid input"`
}

// Error401Response 401 Unauthorized
type Error401Response struct {
	Status StatusUnauthorized `json:"status"`
	Error  string             `json:"error" example:"Unauthorized"`
}

// Error403Response 403 Forbidden
type Error403Response struct {
	Status StatusForbidden `json:"status"`
	Error  string          `json:"error" example:"Forbidden: insufficient permissions"`
}

// Error500Response 500 Internal Server Error
type Error500Response struct {
	Status StatusInternalError `json:"status"`
	Error  string              `json:"error" example:"Internal Server Error"`
}

// ===================== Requests =====================

// LoginRequest represents login request body
type LoginRequest struct {
	Email       string `json:"email" example:"user@example.com"`
	Username    string `json:"username" example:"johndoe"`
	PhoneNumber string `json:"phone_number" example:"+84901234567"`
	Password    string `json:"password" example:"password123" binding:"required"`
}

// GoogleLoginRequest represents Google OAuth login request
type GoogleLoginRequest struct {
	Token string `json:"token" binding:"required" example:"google-oauth-token"`
}

// TwitterLoginRequest represents Twitter login request
type TwitterLoginRequest struct {
	Username               string `json:"username" binding:"required" example:"johndoe"`
	Email                  string `json:"email" binding:"required" example:"user@example.com"`
	Image                  string `json:"image" example:"https://example.com/avatar.jpg"`
	ProfileBackgroundImage string `json:"background_image" example:"https://example.com/bg.jpg"`
	ProfileBackgroundColor string `json:"background_color" example:"#1DA1F2"`
	TwitterID              string `json:"twitter_id" binding:"required" example:"123456789"`
}

// ForgotPasswordRequest represents forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// RenewPasswordRequest represents renew password request (without old password)
type RenewPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8" example:"newPassword123"`
	UserID      uint   `json:"user_id" binding:"required" example:"1"`
}

// ChangePasswordRequest represents change password request (with old password)
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"oldPassword123"`
	NewPassword string `json:"new_password" binding:"required,min=8" example:"newPassword123"`
	UserID      uint   `json:"user_id" binding:"required" example:"1"`
}

// GetUserAvatarRequest represents get user avatar request
type GetUserAvatarRequest struct {
	UserID uint `json:"user_id" binding:"required" example:"1"`
}

// ===================== Response DTOs (with examples) =====================

// UserPublicResponse Public user info in responses
type UserPublicResponse struct {
	ID        uint   `json:"id" example:"1"`
	Email     string `json:"email" example:"user@example.com"`
	Username  string `json:"username" example:"johndoe"`
	FullName  string `json:"full_name" example:"John Doe"`
	AvatarURL string `json:"avatar_url" example:"https://cdn.example.com/u/1.png"`
	Role      string `json:"role" example:"USER"`
}

// LoginInitiatedData Data for login step where verification email is sent
type LoginInitiatedData struct {
	Email string `json:"email" example:"user@example.com"`
}

// TokensWithUser Tokens + user info (for Google/Twitter login and refresh token)
type TokensWithUser struct {
	AccessToken  string             `json:"access_token" example:"eyJhbGciOi..."`
	RefreshToken string             `json:"refresh_token" example:"eyJhbGciOi..."`
	ExpiresIn    int64              `json:"expires_in" example:"3600"`
	TokenType    string             `json:"token_type" example:"Bearer"`
	User         UserPublicResponse `json:"user"`
}

// AvatarData Avatar response data
type AvatarData struct {
	AvatarURL string `json:"avatar_url" example:"https://cdn.example.com/u/1.png"`
}

type StatusResponseVerification struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"Verification code sent to your email. Please verify to process."`
}
type SuccessLoginEmailResponse struct {
	Status StatusResponseVerification `json:"status"`
	Data   string                     `json:"data" example:"example@hotmail.com"`
}

// ================== VERIFICATION CONTROLLER TYPES ==================

type VerifyCodeRequest struct {
	Email string `json:"email" example:"ties.node@outlook.com"`
	Code  string `json:"code" example:"257699"`
}
type StatusVerificationCode struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"Verification successful"`
}
type TokenData struct {
	AccessToken  string `json:"access_token" example:"8ef033a06c2035ad5ba7b585918b7455a899e35da086ee0e84c98303d01ba9fc"`
	RefreshToken string `json:"refresh_token" example:"f38c31612b2b6f03caaa8a8fb8e0809ab21e7403d9fe7c9fb358a5b5255c1959"`
}

type SuccessVerifyCodeResponse struct {
	Status StatusVerificationCode `json:"status"`
	Data   TokenData              `json:"data"`
}

// ================== ROLE CONTROLLER TYPES ==================
// ---- Request ----

type CreateRoleRequest struct {
	Name string `json:"name" example:"team_leader"`
}

// ---- Success Responses ----

type SuccessCreateRoleResponse struct {
	Status StatusCreated `json:"status"`
	Data   RoleResponse  `json:"data"`
}

type SuccessListRolesResponse struct {
	Status StatusOK       `json:"status"`
	Data   []RoleResponse `json:"data"`
}

// ---- Role schema ----

type RoleResponse struct {
	ID   string `json:"id" example:"79fb3083-a010-11f0-94f9-362064ef513e"`
	Name string `json:"name" example:"team_leader"`
}

// ================== SOCIAL BLOCK CONTROLLER TYPES ==================

type SuccessListCategoriesResponse struct {
	Status StatusOK `json:"status"`
	Data   []string `json:"data" example:"[]"`
}

type LikeRequest struct {
	PostID uuid.UUID `json:"post_id"`
}

type SuccessLikeAPostResponse struct {
	Status StatusOK `json:"status"`
	Data   struct {
		TotalLiked uint `json:"totalLikes" example:"1"`
	}
}
type SuccessDisLikeAPostResponse struct {
	Status StatusOK `json:"status"`
	Data   struct {
		TotalLiked uint `json:"totalLikes" example:"1"`
	}
}

// ================== NOTIFICATIONS BLOCK CONTROLLER TYPES ==================

type ListNotificationRequest struct {
	CurrentPage  int    `json:"current_page" example:"1"`
	ItemsPerPage int    `json:"items_per_page" example:"10"`
	OrderBy      string `json:"order_by" example:"id"`
	SortBy       string `json:"sort_by" example:"asc"`
	SearchValue  string `json:"search_value" example:"general"`
	IsSeen       bool   `json:"is_seen" example:"true"`
	All          bool   `json:"all" example:"true"`
}

type UpdateSeenRequest struct {
	NotificationID uuid.UUID `json:"notification_id" example:"36byte"`
}

// ================== PERMISSION BLOCK CONTROLLER TYPES ==================

type CreatePermissionRequest struct {
	PermissionName string `json:"name" example:"allow_create_post"`
}

type AssignPermissionRequest struct {
	RoleID        uuid.UUID   `json:"role_id" example:"36byte"`
	PermissionIDs []uuid.UUID `json:"permission_ids" example:"[]"`
}
