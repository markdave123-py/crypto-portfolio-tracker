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

func (h *PricesHandler) GetPrices(w http.ResponseWriter, r *http.Request) {
	var req PricesRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Assets) == 0 {
		http.Error(w, "assets cannot be empty", http.StatusBadRequest)
		return
	}

	assets := make([]pricing.AssetRef, 0, len(req.Assets))
	for _, a := range req.Assets {
		if a.Chain == "" || a.ContractAddress == "" {
			http.Error(w, "invalid asset entry", http.StatusBadRequest)
			return
		}
		assets = append(assets, a.ToAssetRef())
	}

	prices, err := h.pricing.GetPrices(r.Context(), assets)
	if err != nil {
		h.logger.Error("pricing-failed", zap.Error(err))
		http.Error(w, "failed to fetch prices", http.StatusInternalServerError)
		return
	}

	resp := PricesResponse{
		Prices: make(map[string]float64),
	}

	for asset, price := range prices {
		key := asset.Chain + ":" + asset.ContractAddress
		resp.Prices[key] = price
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

