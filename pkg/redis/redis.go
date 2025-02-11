package redis

import (
	"github.com/go-redis/redis/v8"
)

func NewRedisClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	_, err := rdb.Ping(context.Background()).Result()
	return rdb, err
}

