package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/transactions"
)

type mockTxService struct{}

type txListResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Page  int                        `json:"page"`
		Limit int                        `json:"limit"`
		Items []transactions.Transaction `json:"items"`
	} `json:"data"`
}

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

	handler := NewTransactionsHandler(&mockTxService{}, zap.NewNop())
	r.Get("/wallets/{wallet}/transactions", handler.List)

	req := httptest.NewRequest(
		http.MethodGet,
		"/wallets/0xabc/transactions?chain=ethereum&page=1&limit=10",
		nil,
	)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp txListResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.True(t, resp.Success)
	require.Len(t, resp.Data.Items, 1)
	require.Equal(t, "tx1", resp.Data.Items[0].Hash)
}
