package access_token

import (
	"fmt"
	"time"

	AccessTokens "github.com/unarya/univia/internal/api/modules/key_token/access_token/models"
	Users "github.com/unarya/univia/internal/api/modules/user/models"
	"github.com/unarya/univia/internal/infrastructure/mysql"
	"github.com/unarya/univia/internal/infrastructure/redis"
)

func VerifyToken(token string) (*Users.User, error) {
	db := mysql.DB
	// Proceed with token validation
	var tokenRecord AccessTokens.AccessToken
	// Try cache first
	cacheKey := fmt.Sprintf("access_token_%s", token)
	if result, err := redis.GetJSON[Users.User](redis.Redis, cacheKey); err == nil && result != nil {
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
	_ = redis.Redis.SetJSON(cacheKey, user, 2*time.Hour)
	return &user, nil
}
