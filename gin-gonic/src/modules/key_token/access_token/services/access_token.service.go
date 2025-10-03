package access_token

import (
	"fmt"
	"time"
	"univia/src/config"
	AccessTokens "univia/src/modules/key_token/access_token/models"
	Users "univia/src/modules/user/models"
	"univia/src/utils/cache"
)

func VerifyToken(token string) (*Users.User, error) {
	db := config.DB
	// Proceed with token validation
	var tokenRecord AccessTokens.AccessToken
	// Try cache first
	cacheKey := fmt.Sprintf("access_token_%s", token)
	if result, err := cache.GetJSON[Users.User](config.Redis, cacheKey); err == nil && result != nil {
		return result, nil
	} else if err != nil {
		fmt.Println(err)
	}
	// Continue
	if err := db.Where("token = ? and status = true", token).First(&tokenRecord).Error; err != nil {
		return nil, err
	}

	// Find User By ID
	var user Users.User
	if err := db.Where("id = ?", tokenRecord.UserID).First(&user).Error; err != nil {
		return nil, err
	}
	_ = config.Redis.SetJSON(cacheKey, user, 2*time.Hour)
	return &user, nil
}
