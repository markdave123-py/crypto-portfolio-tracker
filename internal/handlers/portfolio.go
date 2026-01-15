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

// GetPortfolio godoc
// @Summary Get portfolio
// @Description Fetch wallet portfolio with live valuation
// @Tags Portfolio
// @Produce json
// @Param wallet path string true "Wallet address"
// @Success 200 {object} handlers.PortfolioResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Router /wallets/{wallet}/portfolio [get]
func (h *PortfolioHandler) Get(w http.ResponseWriter, r *http.Request) {
	wallet := chi.URLParam(r, "wallet")

	portfolio, err := h.service.Get(r.Context(), wallet)
	if err != nil {
		h.logger.Error("get-portfolio-failed", zap.Error(err))
		RespondError(
			w,
			http.StatusNotFound,
			"PORTFOLIO_NOT_FOUND",
			"portfolio not found",
		)
		return
	}

	RespondOK(w, http.StatusOK, portfolio)
}

// AddHolding godoc
// @Summary Add holding
// @Tags Portfolio
// @Accept json
// @Produce json
// @Param wallet path string true "Wallet address"
// @Param holding body portfolio.Holding true "Holding"
// @Success 201
// @Failure 400 {object} handlers.ErrorResponse
// @Router /wallets/{wallet}/portfolio/holdings [post]
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
		RespondError(
			w,
			http.StatusNotFound,
			"FAILED_TO_ADD_HOLDING",
			"add holding failed",
		)
		return
	}

	RespondOK(w, http.StatusCreated, nil)
}

// UpdateHolding godoc
// @Summary Update holding
// @Tags Portfolio
// @Accept json
// @Produce json
// @Param wallet path string true "Wallet address"
// @Param holding body portfolio.Holding true "Holding"
// @Success 200
// @Failure 400 {object} handlers.ErrorResponse
// @Router /wallets/{wallet}/portfolio/holdings [put]
func (h *PortfolioHandler) UpdateHolding(w http.ResponseWriter, r *http.Request) {
	wallet := chi.URLParam(r, "wallet")

	var req AddHoldingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request")
		return
	}

	err := h.service.UpdateHolding(r.Context(), wallet, portfolio.Holding{
		Chain:           req.Chain,
		ContractAddress: req.ContractAddress,
		Amount:          req.Amount,
	})

	if err != nil {
		h.logger.Error("update-holding-failed", zap.Error(err))
		RespondError(w, http.StatusNotFound, "NOT_FOUND", "failed to update holding")
		return
	}

	RespondOK(w, http.StatusOK, nil)
}

// RemoveHolding godoc
// @Summary Delete holding
// @Tags Portfolio
// @Accept json
// @Produce json
// @Param wallet path string true "Wallet address"
// @Param holding body portfolio.Holding true "Holding"
// @Success 200
// @Failure 400 {object} handlers.ErrorResponse
// @Router /wallets/{wallet}/portfolio/holdings [delete]
func (h *PortfolioHandler) RemoveHolding(w http.ResponseWriter, r *http.Request) {
	wallet := chi.URLParam(r, "wallet")
	chain := r.URL.Query().Get("chain")
	contract := r.URL.Query().Get("contract")

	if chain == "" {
		RespondError(w, http.StatusNotFound, "NOT_FOUND", "missing chain")
		return
	}

	err := h.service.RemoveHolding(r.Context(), wallet, chain, contract)
	if err != nil {
		h.logger.Error("remove-holding-failed", zap.Error(err))
		RespondError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to remove holding")
		return
	}

	RespondOK(w, http.StatusOK, nil)
}
