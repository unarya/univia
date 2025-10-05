package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

var (
	Redis *RedisCache
	Ctx   = context.Background()
)

func ConnectRedis() *RedisCache {
	mode := os.Getenv("REDIS_MODE")

	var rdb *redis.Client

	if mode == "sentinel" {
		masterName := os.Getenv("REDIS_MASTER_NAME")
		password := os.Getenv("REDIS_PASSWORD")
		sentinels := []string{
			"redis-sentinel-node-0.redis-sentinel-headless.architecture.svc.cluster.local:26379",
			"redis-sentinel-node-1.redis-sentinel-headless.architecture.svc.cluster.local:26379",
			"redis-sentinel-node-2.redis-sentinel-headless.architecture.svc.cluster.local:26379",
			"redis-sentinel-node-3.redis-sentinel-headless.architecture.svc.cluster.local:26379",
			"redis-sentinel-node-4.redis-sentinel-headless.architecture.svc.cluster.local:26379",
		}

		rdb = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       masterName,
			SentinelAddrs:    sentinels,
			Password:         password,
			SentinelPassword: password,
			DB:               0,
		})
	} else {
		host := os.Getenv("REDIS_HOST")
		port := os.Getenv("REDIS_PORT")
		password := os.Getenv("REDIS_PASSWORD")
		addr := fmt.Sprintf("%s:%s", host, port)

		rdb = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       0,
		})
	}

	pong, err := rdb.Ping(Ctx).Result()
	if err != nil {
		fmt.Printf("❌ Failed to connect to Redis (%s mode): %v\n", mode, err)
		return nil
	}

	fmt.Println("✅ Redis connected (mode:", mode, "):", pong)

	Redis = NewRedisCache(rdb)
	return Redis
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
