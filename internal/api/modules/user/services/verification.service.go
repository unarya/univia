package users

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	RefreshTokens "github.com/unarya/univia/internal/api/modules/key_token/refresh_token/models"
	refreshTokenService "github.com/unarya/univia/internal/api/modules/key_token/refresh_token/services"
	"github.com/unarya/univia/internal/api/modules/session/queries"
	"github.com/unarya/univia/internal/api/modules/user/models"
	"github.com/unarya/univia/internal/infrastructure/mysql"
	"github.com/unarya/univia/pkg/types"
	"github.com/unarya/univia/pkg/utils"

	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

// Helper function to generate a 6-digit verification code
func generateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// Helper function to save the verification code
func saveVerificationCode(email, code string) error {
	db := mysql.DB
	var user users.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	// Delete all verification code before sent a new one
	if err := db.Where("email = ?", email).Delete(&users.VerificationCode{}).Error; err != nil {
		return fmt.Errorf("failed to delete verification code: %v", err)
	}

	verification := users.VerificationCode{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(3 * time.Minute),
	}
	return db.Create(&verification).Error
}

// sendVerificationEmail sends a verification email with the given code
func sendVerificationEmail(email, code string) error {
	// Email server mysqluration
	smtpHost := os.Getenv("SMTP_HOST")     // E.g., "smtp.gmail.com"
	smtpPort := os.Getenv("SMTP_PORT")     // E.g., "587"
	smtpUsername := os.Getenv("SMTP_USER") // Your email address
	smtpPassword := os.Getenv("SMTP_PASS") // Your email password or app-specific password
	// Email content
	subject := "Your Verification Code"
	body := fmt.Sprintf("Your verification code is: %s", code)
	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	// Set up authentication information
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// Send the email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", smtpHost, smtpPort),
		auth,
		smtpUsername,
		[]string{email},
		[]byte(message),
	)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

// VerifyCodeAndGenerateTokens Function to verify the code and generate tokens
func VerifyCodeAndGenerateTokens(c *gin.Context, code users.VerificationCode, meta types.SessionMetadata, sessionID string) (types.ResponseSession, int, error) {
	db := mysql.DB

	var user users.User
	if err := db.Where("email = ?", code.Email).First(&user).Error; err != nil {
		return types.ResponseSession{}, http.StatusNotFound, err
	}

	var verification users.VerificationCode
	if err := db.Where("email = ?", code.Email).First(&verification).Error; err != nil {
		return types.ResponseSession{}, http.StatusNotFound, err
	}

	// Check valid
	if verification.ExpiresAt.Before(time.Now()) {
		db.Delete(&verification)
		return types.ResponseSession{}, http.StatusUnauthorized, fmt.Errorf("verification code has expired")
	}

	// Check code
	if verification.Code != code.Code {
		verification.InputCount++
		db.Save(&verification)

		if verification.InputCount >= 5 {
			db.Delete(&verification)
			user.Status = false
			db.Save(&user)
			return types.ResponseSession{}, http.StatusUnauthorized, fmt.Errorf("too many failed attempts, your account has been suspended")
		}

		return types.ResponseSession{}, http.StatusUnauthorized, fmt.Errorf("invalid verification code")
	}

	// Init tokens
	accessToken, refreshToken, err := refreshTokenService.GenerateHexTokens(user.ID)
	if err != nil {
		return types.ResponseSession{}, http.StatusInternalServerError, err
	}

	var rID uuid.UUID
	if err := db.Model(&RefreshTokens.RefreshToken{}).
		Select("id").
		Where("token = ?", refreshToken).
		Scan(&rID).Error; err != nil {
		return types.ResponseSession{}, http.StatusInternalServerError, fmt.Errorf("failed to find refresh token: %v", err)
	}
	// Create new session record
	if sessionID == "" {
		// Init new session record
		sID, err := queries.InsertNewSessionByUserID(user.ID, rID, meta)
		if err != nil {
			return types.ResponseSession{}, http.StatusInternalServerError, fmt.Errorf("failed to insert new session: %v", err)
		}
		sessionID = sID.String()
	}

	// Session exists
	err = queries.InsertIntoSessionByValidSession(utils.ConvertStringToUuid(sessionID), user.ID, rID, meta)
	if err != nil {
		return types.ResponseSession{}, http.StatusInternalServerError, fmt.Errorf("failed to insert new session: %v", err)
	}
	sID, err := utils.SetSessionToRedisByUserID(c, db, user)

	// Delete verification record
	db.Delete(&verification)

	return types.ResponseSession{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionID:    sID,
		UserID:       user.ID,
	}, http.StatusOK, nil
}

func VerifyCode(code, email string) error {
	// Import db queries
	db := mysql.DB
	var user users.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	var verification users.VerificationCode
	if err := db.Where("email = ? AND code = ?", email, code).First(&verification).Error; err != nil {
		return err
	}
	// Step 3: Check if the code is expired
	if verification.ExpiresAt.Before(time.Now()) {
		db.Delete(&verification) // Ensure to delete expired verification code
		return fmt.Errorf("verification code has expired")
	}
	// Step 4: Check valid code
	if verification.Code != code {
		verification.InputCount += 1
		db.Save(&verification)

		// Step 5: Lock user after 5 failed attempts
		if verification.InputCount >= 5 {
			db.Delete(&verification)
			user.Status = false
			db.Save(&user)
			return fmt.Errorf("too many time argument, your account had been suspended, please get contact to admin")
		}
		return fmt.Errorf("invalid verification code")
	}
	// Step 6: Delete all existing tokens for the user
	if err := DeleteAllTokensByUserID(user.ID); err != nil {
		return fmt.Errorf("failed to delete existing tokens")
	}

	return nil
}
