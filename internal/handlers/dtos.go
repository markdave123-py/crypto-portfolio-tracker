package handlers

import "github.com/markdave123-py/crypto-portfolio-tracker/internal/pricing"

type PricesRequest struct {
	Assets []AssetRequest `json:"assets"`
}

type AssetRequest struct {
	Chain           string `json:"chain"`
	ContractAddress string `json:"contract_address"`
}

func (r AssetRequest) ToAssetRef() pricing.AssetRef {
	return pricing.AssetRef{
		Chain:           r.Chain,
		ContractAddress: r.ContractAddress,
	}
}

type PricesResponse struct {
	Prices map[string]float64 `json:"prices"`
}

type AddHoldingRequest struct {
	Chain           string  `json:"chain"`
	ContractAddress string  `json:"contract_address"`
	Amount          float64 `json:"amount"`
}
