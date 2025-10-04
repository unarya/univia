package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
		ctx:    context.Background(),
	}
}

// SetJSON marshal object to JSON and set to Redis
func (rc *RedisCache) SetJSON(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rc.client.Set(rc.ctx, key, data, ttl).Err()
}

// GetJSON unmarshal data from Redis to struct generic
func GetJSON[T any](rc *RedisCache, key string) (*T, error) {
	data, err := rc.client.Get(rc.ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var obj T
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

// Delete xóa key
func (rc *RedisCache) Delete(key string) error {
	return rc.client.Del(rc.ctx, key).Err()
}

// UserCacheKey Key generator helpers
func UserCacheKey(userID uuid.UUID) string {
	return fmt.Sprintf("user:%d:info", userID)
}

func RolePermissionsCacheKey(roleID uint) string {
	return fmt.Sprintf("role:%d:permissions", roleID)
}
