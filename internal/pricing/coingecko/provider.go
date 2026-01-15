package coingecko

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/pricing"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/utils"
)

type Provider struct {
	client *Client
}

func NewProvider(client *Client) *Provider {
	return &Provider{client: client}
}

func (p *Provider) Name() string {
	return "coingecko"
}

func (p *Provider) GetPrices(ctx context.Context, assets []pricing.AssetRef) (map[pricing.AssetRef]float64, error) {

	result := make(map[pricing.AssetRef]float64)

	// group by chain as required by coingecko
	grouped := make(map[string][]pricing.AssetRef)
	for _, a := range assets {
		grouped[a.Chain] = append(grouped[a.Chain], a)
	}

	for chain, group := range grouped {
		// retry with exponential backoff
		err := utils.Retry(ctx, utils.RetryConfig{
			MaxRetries: 3,
			BaseDelay:  500 * time.Millisecond,
			MaxDelay:   4 * time.Second,
		}, func() error {

			contracts := make([]string, 0, len(group))
			for _, a := range group {
				contracts = append(contracts, a.ContractAddress)
			}

			raw, err := p.client.FetchTokenPrices(ctx, chain, contracts)
			if err != nil {
				return err
			}

			for _, a := range group {
				addr := strings.ToLower(a.ContractAddress)
				v, ok := raw[addr]
				if !ok {
					continue
				}
				obj, ok := v.(map[string]any)
				if !ok {
					continue
				}

				usdVal, ok := obj["usd"]
				if !ok {
					continue
				}
				var price float64

				switch val := usdVal.(type) {
				case float64:
					price = val
				case string:
					price, err = strconv.ParseFloat(val, 64)
					if err != nil {
						continue
					}
				default:
					continue
				}
				result[a] = price

			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
