package coingecko

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/pricing"
	"github.com/test-go/testify/require"
)

func TestCoinGeckoProvider_ParsesResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"0xabc": { "usd": 123.45 }
		}`))
	}))
	defer ts.Close()

	provider := NewProvider(
		NewClient("test", ts.URL),
	)

	asset := pricing.AssetRef{
		Chain:           "ethereum",
		ContractAddress: "0xabc",
	}

	prices, err := provider.GetPrices(context.Background(), []pricing.AssetRef{asset})

	require.NoError(t, err)
	require.Equal(t, 123.45, prices[asset])
}
