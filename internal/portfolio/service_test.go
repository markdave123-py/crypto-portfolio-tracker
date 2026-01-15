package portfolio_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/portfolio"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/pricing"
)

type mockPricingService struct {
	prices map[pricing.AssetRef]float64
}

func (m *mockPricingService) GetPrices(
	ctx context.Context,
	assets []pricing.AssetRef,
) (map[pricing.AssetRef]float64, error) {
	return m.prices, nil
}

func setupService() portfolio.Service {
	logger := zap.NewNop()

	repo := portfolio.NewMemoryRepository([]*portfolio.Portfolio{
		{
			Wallet: "wallet1",
			Holdings: []portfolio.Holding{
				{
					Chain:           "ethereum",
					ContractAddress: "",
					Amount:          2,
				},
			},
		},
	})

	pricingSvc := &mockPricingService{
		prices: map[pricing.AssetRef]float64{
			{Chain: "ethereum", ContractAddress: ""}: 2000,
		},
	}

	return portfolio.NewService(repo, pricingSvc, logger)
}

func TestGetPortfolio(t *testing.T) {
	svc := setupService()

	view, err := svc.Get(context.Background(), "wallet1")
	require.NoError(t, err)

	require.Equal(t, "wallet1", view.Wallet)
	require.Len(t, view.Holdings, 1)

	h := view.Holdings[0]
	require.Equal(t, 2.0, h.Amount)
	require.Equal(t, 2000.0, h.PriceUSD)
	require.Equal(t, 4000.0, h.ValueUSD)

	require.Equal(t, 4000.0, view.TotalValueUSD)
}

func TestAddHolding(t *testing.T) {
	svc := setupService()

	err := svc.AddHolding(context.Background(), "wallet1", portfolio.Holding{
		Chain:           "ethereum",
		ContractAddress: "0xusdc",
		Amount:          100,
	})

	require.NoError(t, err)

	view, _ := svc.Get(context.Background(), "wallet1")
	require.Len(t, view.Holdings, 2)
}

func TestUpdateHolding(t *testing.T) {
	svc := setupService()

	err := svc.UpdateHolding(context.Background(), "wallet1", portfolio.Holding{
		Chain:           "ethereum",
		ContractAddress: "",
		Amount:          5,
	})

	require.NoError(t, err)

	view, _ := svc.Get(context.Background(), "wallet1")
	require.Equal(t, 5.0, view.Holdings[0].Amount)
}

func TestRemoveHolding(t *testing.T) {
	svc := setupService()

	err := svc.RemoveHolding(context.Background(), "wallet1", "ethereum", "")
	require.NoError(t, err)

	view, _ := svc.Get(context.Background(), "wallet1")
	require.Len(t, view.Holdings, 0)
}
