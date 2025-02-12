package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	return rdb, err
}
