package mock

import (
	"context"
	"hash/fnv"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/pricing"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) Name() string {
	return "mock"
}

var _ pricing.PriceProvider = (*Provider)(nil)

func (p *Provider) GetPrices(
	ctx context.Context,
	assets []pricing.AssetRef,
) (map[pricing.AssetRef]float64, error) {

	result := make(map[pricing.AssetRef]float64)
	for _, a := range assets {
		result[a] = deterministicPrice(a.ContractAddress)
	}
	return result, nil
}

func deterministicPrice(input string) float64 {
	h := fnv.New32a()
	h.Write([]byte(input))
	return float64(h.Sum32()%50_000) / 100
}
