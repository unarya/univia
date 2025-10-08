package refresh_token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/unarya/univia/internal/api/modules/key_token/access_token/models"
	"github.com/unarya/univia/internal/api/modules/key_token/refresh_token/models"
	"github.com/unarya/univia/internal/infrastructure/mysql"
	"gorm.io/gorm"
)

func GetRefreshTokenIDByToken(token string) (uuid.UUID, error) {
	db := mysql.DB
	var rID uuid.UUID
	if err := db.Model(&refresh_token.RefreshToken{}).
		Where("token = ?", token).
		Pluck("id", &rID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return uuid.Nil, fmt.Errorf("token not found")
		}
		return uuid.Nil, fmt.Errorf("db error: %w", err)
	}
	return rID, nil
}

func GenerateHexTokens(userID uuid.UUID) (string, string, error) {
	db := mysql.DB

	// Generate random hex strings for tokens
	accessTokenBytes := make([]byte, 32) // 256-bit
	refreshTokenBytes := make([]byte, 32)

	_, err := rand.Read(accessTokenBytes)
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}
	_, err = rand.Read(refreshTokenBytes)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	accessToken := hex.EncodeToString(accessTokenBytes)
	refreshToken := hex.EncodeToString(refreshTokenBytes)

	// Save access token to mysql
	accessTokenEntry := access_token.AccessToken{
		UserID: userID,
		Token:  accessToken,
	}
	if err := db.Create(&accessTokenEntry).Error; err != nil {
		return "", "", err
	}

	// Save refresh token to mysql
	refreshTokenEntry := refresh_token.RefreshToken{
		UserID: userID,
		Token:  refreshToken,
	}
	if err := db.Create(&refreshTokenEntry).Error; err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
