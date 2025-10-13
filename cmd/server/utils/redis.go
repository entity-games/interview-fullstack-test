package utils

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedisPool(config Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         config.RedisDsn,
		Password:     "",
		DB:           0,
		PoolSize:     10, // max number of socket connections
		MinIdleConns: 3,  // number of idle connections
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
}

func PingRedis(ctx context.Context, client *redis.Client) bool {
	_, err := client.Set(ctx, "ping", "pong", time.Second).Result()
	return err == nil
}
