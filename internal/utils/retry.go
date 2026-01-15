package utils

import (
	"context"
	"time"
)

type RetryConfig struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

// Retry: Retries a function fn for MaxRetries times with exponential backoff
func Retry(
	ctx context.Context,
	cfg RetryConfig,
	fn func() error,
) error {
	var err error
	delay := cfg.BaseDelay

	for i := 0; i <= cfg.MaxRetries; i++ {
		if err = fn(); err == nil {
			return nil
		}

		if i == cfg.MaxRetries {
			break
		}

		select {
		case <-time.After(delay):
			delay *= 2
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return err
}
