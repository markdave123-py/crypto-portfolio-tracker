package pricing

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/cache"
	"github.com/test-go/testify/require"
	"go.uber.org/zap"
)

type fakeCache struct {
	data map[string]string
}

func newFakeCache() cache.CacheManager {
	return &fakeCache{data: make(map[string]string)}
}

func (f *fakeCache) Get(ctx context.Context, key string) (string, error) {
	v, ok := f.data[key]
	if !ok {
		return "", fmt.Errorf("cache key not found")
	}
	return v, nil
}

func (f *fakeCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	f.data[key] = value
	return nil
}

func (f *fakeCache) Del(ctx context.Context, key string) error {
	delete(f.data, key)
	return nil
}

type fakeProvider struct {
	name   string
	prices map[AssetRef]float64
	err    error
	calls  int
}

func (f *fakeProvider) Name() string {
	return f.name
}

func (f *fakeProvider) GetPrices(
	ctx context.Context,
	assets []AssetRef,
) (map[AssetRef]float64, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	return f.prices, nil
}

func TestPricingService_CacheHit(t *testing.T) {
	cache := newFakeCache()
	asset := AssetRef{Chain: "ethereum", ContractAddress: "0xabc"}
	cache.Set(context.Background(), cacheKey(asset), "123.0", time.Minute)

	provider := &fakeProvider{name: "primary"}

	svc := NewService(
		cache,
		provider,
		nil,
		time.Minute,
		zap.NewNop(),
	)

	prices, err := svc.GetPrices(context.Background(), []AssetRef{asset})

	require.NoError(t, err)
	require.Equal(t, 123.0, prices[asset])
	require.Equal(t, 0, provider.calls)
}

func TestPricingService_PrimarySuccess(t *testing.T) {
	cache := newFakeCache()
	asset := AssetRef{Chain: "ethereum", ContractAddress: "0xabc"}

	primary := &fakeProvider{
		name: "primary",
		prices: map[AssetRef]float64{
			asset: 42.0,
		},
	}

	svc := NewService(
		cache,
		primary,
		nil,
		time.Minute,
		zap.NewNop(),
	)

	prices, err := svc.GetPrices(context.Background(), []AssetRef{asset})

	require.NoError(t, err)
	require.Equal(t, 42.0, prices[asset])
	require.Equal(t, 1, primary.calls)
}

func TestPricingService_FallbackUsed(t *testing.T) {
	cache := newFakeCache()
	asset := AssetRef{Chain: "ethereum", ContractAddress: "0xabc"}

	primary := &fakeProvider{
		name: "primary",
		err:  errors.New("primary down"),
	}
	fallback := &fakeProvider{
		name: "fallback",
		prices: map[AssetRef]float64{
			asset: 99.0,
		},
	}

	svc := NewService(
		cache,
		primary,
		fallback,
		time.Minute,
		zap.NewNop(),
	)

	prices, err := svc.GetPrices(context.Background(), []AssetRef{asset})

	require.NoError(t, err)
	require.Equal(t, 99.0, prices[asset])
	require.Equal(t, 1, primary.calls)
	require.Equal(t, 1, fallback.calls)
}

func TestPricingService_AllProvidersFail(t *testing.T) {
	cache := newFakeCache()
	asset := AssetRef{Chain: "ethereum", ContractAddress: "0xabc"}

	primary := &fakeProvider{name: "primary", err: errors.New("down")}
	fallback := &fakeProvider{name: "fallback", err: errors.New("down")}

	svc := NewService(
		cache,
		primary,
		fallback,
		time.Minute,
		zap.NewNop(),
	)

	_, err := svc.GetPrices(context.Background(), []AssetRef{asset})

	require.Error(t, err)
}
