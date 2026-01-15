package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/transactions"
)

type TransactionsHandler struct {
	service transactions.ServiceAPI
}

func NewTransactionsHandler(s transactions.ServiceAPI) *TransactionsHandler {
	return &TransactionsHandler{service: s}
}

func (h *TransactionsHandler) List(w http.ResponseWriter, r *http.Request) {
	wallet := chi.URLParam(r, "wallet")
	chain := r.URL.Query().Get("chain")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page == 0 {
		page = 1
	}
	if limit == 0 {
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(txs)
}
