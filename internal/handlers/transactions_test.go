package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/transactions"
)

type mockTxService struct{}

func (m *mockTxService) List(
	ctx context.Context,
	chain string,
	wallet string,
	page int,
	limit int,
	filters transactions.Filters,
) ([]transactions.Transaction, error) {
	return []transactions.Transaction{
		{
			Hash:  "tx1",
			Chain: "ethereum",
		},
	}, nil
}

func TestTransactionsHandler_List(t *testing.T) {
	r := chi.NewRouter()

	handler := NewTransactionsHandler(&mockTxService{})
	r.Get("/wallets/{wallet}/transactions", handler.List)

	req := httptest.NewRequest(
		http.MethodGet,
		"/wallets/0xabc/transactions?chain=ethereum&page=1&limit=10",
		nil,
	)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp []transactions.Transaction
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Len(t, resp, 1)
	require.Equal(t, "tx1", resp[0].Hash)
}
