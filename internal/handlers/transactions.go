package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/transactions"
	"go.uber.org/zap"
)

type TransactionsHandler struct {
	service transactions.ServiceAPI
	logger  *zap.Logger
}

func NewTransactionsHandler(s transactions.ServiceAPI, logger *zap.Logger) *TransactionsHandler {
	return &TransactionsHandler{service: s, logger: logger}
}


// ListTransactions godoc
// @Summary List wallet transactions
// @Description Fetch paginated transactions for a wallet
// @Tags Transactions
// @Produce json
// @Param wallet path string true "Wallet address"
// @Param chain query string true "Blockchain (ethereum)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(20)
// @Param type query string false "Transaction type"
// @Param status query string false "Transaction status"
// @Param token query string false "Token symbol"
// @Param start_date query string false "Start date RFC3339"
// @Param end_date query string false "End date RFC3339"
// @Success 200 {object} handlers.TransactionListResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /wallets/{wallet}/transactions [get]
func (h *TransactionsHandler) List(w http.ResponseWriter, r *http.Request) {
	wallet := chi.URLParam(r, "wallet")
	chain := r.URL.Query().Get("chain")

	if wallet == "" || chain == "" {
		RespondError(
			w,
			http.StatusBadRequest,
			"MISSING_PARAMS",
			"wallet and chain are required",
		)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	var filters transactions.Filters

	if v := r.URL.Query().Get("type"); v != "" {
		t := transactions.TransactionType(v)
		filters.Type = &t
	}

	if v := r.URL.Query().Get("status"); v != "" {
		s := transactions.TransactionStatus(v)
		filters.Status = &s
	}

	if v := r.URL.Query().Get("token"); v != "" {
		filters.Token = &v
	}

	if v := r.URL.Query().Get("start_date"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filters.StartDate = &t
		}
	}

	if v := r.URL.Query().Get("end_date"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filters.EndDate = &t
		}
	}

	txs, err := h.service.List(r.Context(), chain, wallet, page, limit, filters)
	if err != nil {
		h.logger.Error("get-transactions-failed", zap.Error(err))
		RespondError(
			w,
			http.StatusInternalServerError,
			"TRANSACTIONS_FAILED",
			"failed to fetch transactions",
		)
		return
	}

	RespondOK(w, http.StatusOK, map[string]any{
		"page":  page,
		"limit": limit,
		"items": txs,
	})
}
