package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"
	"go.uber.org/zap"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/pricing"
)

type mockPricingService struct {
	result map[pricing.AssetRef]float64
	err    error
}

type pricesResponseTest struct {
	Success bool           `json:"success"`
	Data    PricesResponse `json:"data"`
}

func (m *mockPricingService) GetPrices(
	ctx context.Context,
	assets []pricing.AssetRef,
) (map[pricing.AssetRef]float64, error) {
	return m.result, m.err
}

func TestPricesHandler_GetPrices_Success(t *testing.T) {
	logger := zap.NewNop()

	mockSvc := &mockPricingService{
		result: map[pricing.AssetRef]float64{
			{
				Chain:           "ethereum",
				ContractAddress: "0xabc",
			}: 123.45,
		},
	}

	handler := NewPricesHandler(mockSvc, logger)

	body := `{
		"assets": [
			{
				"chain": "ethereum",
				"contract_address": "0xabc"
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/prices", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.GetPrices(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp pricesResponseTest
	fmt.Println(string(rec.Body.Bytes()))
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Len(t, resp.Data.Prices, 1)
	require.Equal(t, 123.45, resp.Data.Prices["ethereum:0xabc"])
}

func TestPricesHandler_GetPrices_InvalidJSON(t *testing.T) {
	logger := zap.NewNop()
	handler := NewPricesHandler(&mockPricingService{}, logger)

	req := httptest.NewRequest(http.MethodPost, "/prices", bytes.NewBufferString("{"))
	rec := httptest.NewRecorder()

	handler.GetPrices(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPricesHandler_GetPrices_EmptyAssets(t *testing.T) {
	logger := zap.NewNop()
	handler := NewPricesHandler(&mockPricingService{}, logger)

	body := `{"assets":[]}`

	req := httptest.NewRequest(http.MethodPost, "/prices", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.GetPrices(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPricesHandler_GetPrices_ServiceError(t *testing.T) {
	logger := zap.NewNop()

	mockSvc := &mockPricingService{
		err: assert.AnError,
	}

	handler := NewPricesHandler(mockSvc, logger)

	body := `{
		"assets": [
			{
				"chain": "ethereum",
				"contract_address": "0xabc"
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/prices", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.GetPrices(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
