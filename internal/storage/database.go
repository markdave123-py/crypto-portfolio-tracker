package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/config"
)

func NewPostgres(ctx context.Context, cfg config.DBConfig, logger *zap.Logger) (*pgxpool.Pool, error) {
	// db connection string
	dsn := cfg.ConnectionString()

	// Bound startup time
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// db connection pool configs
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	// Fail fast if DB is unreachable
	if err := pool.Ping(ctx); err != nil {
		logger.Error("postgres-ping-failed", zap.Error(err))
		pool.Close()
		return nil, err
	}

	logger.Info("postgres-connected")
	return pool, nil
}
