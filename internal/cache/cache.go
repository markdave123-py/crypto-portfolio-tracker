package cache

import (
	"context"
	"time"
)


type CacheManager interface {
	// Set stores a key-value pair in the cache
	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	// Get retrieves a value from the cache
	Get(ctx context.Context, key string) (string, error)

	// Del deletes a key from the cache
	Del(ctx context.Context, key string) error
}
