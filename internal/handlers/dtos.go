package handlers

import (
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/portfolio"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/pricing"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/transactions"
)

// price handler dtos
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

type PriceAPIResponse struct {
	Success bool           `json:"success"`
	Data    PricesResponse `json:"data,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

// transaction handler dtos
type TransactionListResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Page  int                        `json:"page"`
		Limit int                        `json:"limit"`
		Items []transactions.Transaction `json:"items"`
	} `json:"data"`
}

// porfolio handler dtos
type AddHoldingRequest struct {
	Chain           string  `json:"chain"`
	ContractAddress string  `json:"contract_address"`
	Amount          float64 `json:"amount"`
}

type PortfolioResponse struct {
	Status  string                   `json:"status"`            // success | error
	Message string                   `json:"message,omitempty"` // human-readable
	Data    *portfolio.PortfolioView `json:"data,omitempty"`
	Error   string                   `json:"error,omitempty"`
}
