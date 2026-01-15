package pricing

import "context"

type AssetRef struct {
	Chain           string // e.g. "ethereum", "polygon"
	ContractAddress string // lowercase hex
}

type PriceProvider interface {
	GetPrices(ctx context.Context, assets []AssetRef) (map[AssetRef]float64, error)
	Name() string
}
