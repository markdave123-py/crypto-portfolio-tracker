package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type redisImplementation struct {
	logger *zap.Logger
	client *redis.Client
}

// NewRedisManager returns an implementation of the cache interface
func NewRedisManager(redisClient *redis.Client, logger *zap.Logger) (c CacheManager, err error) {
	redisLogger := logger.With(zap.String("service", "redis"))
	c = &redisImplementation{
		logger: redisLogger,
		client: redisClient,
	}

	return c, err
}

// Set stores a key-value pair in the cache
func (r *redisImplementation) Set(ctx context.Context, key string, value string, ttl time.Duration) (err error) {
	if err := r.client.Set(ctx, key, value, ttl).Err(); err != nil {
		r.logger.Error("failed-to-set-value-in-cache", zap.Error(err))
		return err
	}
	return nil
}

// Get retrieves a value from the cache
func (r *redisImplementation) Get(ctx context.Context, key string) (value string, err error) {
	value, err = r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			r.logger.Warn("key-does-not-exist-in-cache", zap.String("key", key))
			return value, fmt.Errorf("key (%s) does not exist in cache: %w", key, errors.New("key not found"))
		}

		r.logger.Error("failed-to-get-value-from-cache", zap.Error(err))
		return "", err
	}
	return value, nil
}

// Del deletes a key from the cache
func (r *redisImplementation) Del(ctx context.Context, key string) (err error) {
	err = r.client.Del(ctx, key).Err()
	if err != nil {
		r.logger.Error("failed-to-delete-key-from-cache", zap.Error(err))
	}
	return err
}
