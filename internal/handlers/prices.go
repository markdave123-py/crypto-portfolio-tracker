package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/pricing"
)

type PricesHandler struct {
	pricing pricing.ServiceAPI
	logger  *zap.Logger
}

func NewPricesHandler(
	pricing pricing.ServiceAPI,
	logger *zap.Logger,
) *PricesHandler {
	return &PricesHandler{
		pricing: pricing,
		logger:  logger,
	}
}

// GetPrices godoc
// @Summary Get token prices
// @Description Fetch USD prices for tokens by chain + contract address
// @Tags Prices
// @Accept json
// @Produce json
// @Param request body handlers.PricesRequest true "Assets to price"
// @Success 200 {object} handlers.PriceAPIResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /prices [post]
func (h *PricesHandler) GetPrices(w http.ResponseWriter, r *http.Request) {
	var req PricesRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(
			w,
			http.StatusBadRequest,
			"INVALID_JSON",
			"invalid request body",
		)
		return
	}

	if len(req.Assets) == 0 {
		RespondError(
			w,
			http.StatusBadRequest,
			"EMPTY_ASSETS",
			"assets cannot be empty",
		)
		return
	}

	assets := make([]pricing.AssetRef, 0, len(req.Assets))
	for _, a := range req.Assets {
		if a.Chain == "" || a.ContractAddress == "" {
			RespondError(
				w,
				http.StatusBadRequest,
				"INVALID_ASSET",
				"chain and contract_address are required",
			)
			return
		}
		assets = append(assets, a.ToAssetRef())
	}

	prices, err := h.pricing.GetPrices(r.Context(), assets)
	if err != nil {
		h.logger.Error("pricing-failed", zap.Error(err))
		RespondError(
			w,
			http.StatusInternalServerError,
			"PRICING_FAILED",
			"failed to fetch prices",
		)
		return
	}

	resp := PricesResponse{
		Prices: make(map[string]float64),
	}

	for asset, price := range prices {
		key := asset.Chain + ":" + asset.ContractAddress
		resp.Prices[key] = price
	}

	RespondOK(w, http.StatusOK, resp)
}
