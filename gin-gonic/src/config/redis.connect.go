package config

import (
	"context"
	"fmt"
	"github.com/deva-labs/univia-api/gin-gonic/src/utils/cache"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Redis *cache.RedisCache
	Ctx   = context.Background()
)

func ConnectRedis() *cache.RedisCache {
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

	Redis = cache.NewRedisCache(rdb)
	return Redis
}
