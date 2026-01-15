package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/portfolio"
)

type PortfolioHandler struct {
	service portfolio.Service
	logger  *zap.Logger
}

func NewPortfolioHandler(
	service portfolio.Service,
	logger *zap.Logger,
) *PortfolioHandler {
	return &PortfolioHandler{
		service: service,
		logger:  logger,
	}
}

func (h *PortfolioHandler) Get(w http.ResponseWriter, r *http.Request) {
	wallet := chi.URLParam(r, "wallet")

	view, err := h.service.Get(r.Context(), wallet)
	if err != nil {
		h.logger.Error("get-portfolio-failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(view)
}

func (h *PortfolioHandler) AddHolding(w http.ResponseWriter, r *http.Request) {
	wallet := chi.URLParam(r, "wallet")

	var req AddHoldingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.AddHolding(r.Context(), wallet, portfolio.Holding{
		Chain:           req.Chain,
		ContractAddress: req.ContractAddress,
		Amount:          req.Amount,
	})

	if err != nil {
		h.logger.Error("add-holding-failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *PortfolioHandler) UpdateHolding(w http.ResponseWriter, r *http.Request) {
	wallet := chi.URLParam(r, "wallet")

	var req AddHoldingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateHolding(r.Context(), wallet, portfolio.Holding{
		Chain:           req.Chain,
		ContractAddress: req.ContractAddress,
		Amount:          req.Amount,
	})

	if err != nil {
		h.logger.Error("update-holding-failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PortfolioHandler) RemoveHolding(w http.ResponseWriter, r *http.Request) {
	wallet := chi.URLParam(r, "wallet")
	chain := r.URL.Query().Get("chain")
	contract := r.URL.Query().Get("contract")

	if chain == "" {
		http.Error(w, "missing chain", http.StatusBadRequest)
		return
	}

	err := h.service.RemoveHolding(r.Context(), wallet, chain, contract)
	if err != nil {
		h.logger.Error("remove-holding-failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
