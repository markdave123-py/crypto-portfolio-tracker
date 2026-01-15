package cache

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(url string) (*redis.Client, error) {
	if url == "" {
		err := errors.New("REDIS_URL is not set")
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: url,
	})

	return redisClient, nil
}
