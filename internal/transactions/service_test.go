package transactions

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockRepository struct {
	txs []Transaction
}

func (m *mockRepository) GetTransactions(
	ctx context.Context,
	chain string,
	wallet string,
	page int,
	limit int,
) ([]Transaction, error) {
	return m.txs, nil
}

func TestService_List_WithFiltering(t *testing.T) {
	repo := &mockRepository{
		txs: []Transaction{
			{
				Hash: "tx1",
				Type: TypeSwap,
				From: "0xabc",
			},
			{
				Hash: "tx2",
				Type: TypeSend,
				From: "0xabc",
			},
		},
	}

	svc := NewService(repo)

	f := Filters{
		Type: ptrType(TypeSwap),
	}

	out, err := svc.List(
		context.Background(),
		"ethereum",
		"0xabc",
		1,
		10,
		f,
	)

	require.NoError(t, err)
	require.Len(t, out, 1)
	require.Equal(t, "tx1", out[0].Hash)
	require.Equal(t, DirectionOut, out[0].Direction)
}
