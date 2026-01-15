package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/handlers"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/portfolio"
)

type mockPortfolioService struct {
	view *portfolio.PortfolioView
}

func (m *mockPortfolioService) Get(ctx context.Context, wallet string) (*portfolio.PortfolioView, error) {
	return m.view, nil
}

func (m *mockPortfolioService) AddHolding(ctx context.Context, wallet string, h portfolio.Holding) error {
	return nil
}

func (m *mockPortfolioService) UpdateHolding(ctx context.Context, wallet string, h portfolio.Holding) error {
	return nil
}

func (m *mockPortfolioService) RemoveHolding(ctx context.Context, wallet, chain, contract string) error {
	return nil
}

func setupRouter(svc portfolio.Service) http.Handler {
	r := chi.NewRouter()
	h := handlers.NewPortfolioHandler(svc, zap.NewNop())

	r.Route("/wallets/{wallet}", func(r chi.Router) {
		r.Get("/portfolio", h.Get)
		r.Post("/holdings", h.AddHolding)
	})

	return r
}

func TestGetPortfolioHandler(t *testing.T) {
	view := &portfolio.PortfolioView{
		Wallet:        "wallet1",
		TotalValueUSD: 1000,
	}

	svc := &mockPortfolioService{view: view}
	router := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/wallets/wallet1/portfolio", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp portfolio.PortfolioView
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, 1000.0, resp.TotalValueUSD)
}

func TestAddHoldingHandler(t *testing.T) {
	svc := &mockPortfolioService{}
	router := setupRouter(svc)

	body := map[string]interface{}{
		"chain":            "ethereum",
		"contract_address": "",
		"amount":           1.5,
	}

	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/wallets/wallet1/holdings", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
}
