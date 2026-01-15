package pricing

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/cache"
)

type ServiceAPI interface {
	GetPrices(
		ctx context.Context,
		assets []AssetRef,
	) (map[AssetRef]float64, error)
}

type Service struct {
	cache    cache.CacheManager
	primary  PriceProvider
	fallback PriceProvider
	cacheTTL time.Duration
	logger   *zap.Logger
}

func NewService(
	cache cache.CacheManager,
	primary PriceProvider,
	fallback PriceProvider,
	cacheTTL time.Duration,
	logger *zap.Logger,
) *Service {
	return &Service{
		cache:    cache,
		primary:  primary,
		fallback: fallback,
		cacheTTL: cacheTTL,
		logger:   logger,
	}
}

func (s *Service) GetPrices(
	ctx context.Context,
	assets []AssetRef,
) (map[AssetRef]float64, error) {
	s.logger.Info("get-prices")

	results := make(map[AssetRef]float64)
	missing := make([]AssetRef, 0)

	// Cache lookup first
	for _, a := range assets {
		key := cacheKey(a)

		cachedStr, err := s.cache.Get(ctx, key)
		if err == nil {
			price, err := strconv.ParseFloat(cachedStr, 64)
			if err == nil {
				results[a] = price
				continue
			}
		}
		missing = append(missing, a)
	}

	// if all is cached
	if len(missing) == 0 {
		return results, nil
	}

	// Try the primary provider (ciongecko)
	prices, err := s.primary.GetPrices(ctx, missing)
	if err != nil {
		s.logger.Warn("primary-pricing-failed",
			zap.String("provider", s.primary.Name()),
			zap.Error(err),
		)

		// Try Fallback provider (Mock)
		prices, err = s.fallback.GetPrices(ctx, missing)
		if err != nil {
			return nil, fmt.Errorf("pricing failed: %w", err)
		}
	}

	// Populate cache and merge results
	for asset, price := range prices {
		key := cacheKey(asset)
		priceStr := strconv.FormatFloat(price, 'f', -1, 64)
		_ = s.cache.Set(ctx, key, priceStr, s.cacheTTL)
		results[asset] = price
	}

	return results, nil
}

func cacheKey(a AssetRef) string {
	return fmt.Sprintf("price:%s:%s", a.Chain, a.ContractAddress)
}
