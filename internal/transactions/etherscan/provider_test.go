package etherscan

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/transactions"
)

type mockClient struct {
	resp *txListResponse
	err  error
}

func (m *mockClient) FetchTxList(
	ctx context.Context,
	chain string,
	wallet string,
	page int,
	limit int,
) (*txListResponse, error) {
	return m.resp, m.err
}

func TestProvider_GetTransactions_Success(t *testing.T) {
	mockResp := &txListResponse{
		Result: []txListItem{
			{
				Hash:            "0xtx1",
				From:            "0xfrom",
				To:              "0xto",
				Value:           "1000000000000000000", // 1 ETH
				TimeStamp:       "1700000000",
				FunctionName:    "transfer(address,uint256)",
				IsError:         "0",
				TxReceiptStatus: "1",
			},
		},
	}

	client := &mockClient{resp: mockResp}
	provider := NewProvider(client)

	txs, err := provider.GetTransactions(
		context.Background(),
		"ethereum",
		"0xwallet",
		1,
		10,
	)

	require.NoError(t, err)
	require.Len(t, txs, 1)

	tx := txs[0]
	require.Equal(t, "0xtx1", tx.Hash)
	require.Equal(t, transactions.TypeSend, tx.Type)
	require.Equal(t, transactions.StatusSuccess, tx.Status)
	require.Equal(t, 1.0, tx.Amount)
	require.Equal(t, time.Unix(1700000000, 0), tx.Timestamp)
}

func TestProvider_GetTransactions_ClientError(t *testing.T) {
	client := &mockClient{
		err: assert.AnError,
	}

	provider := NewProvider(client)
	provider.client = client

	_, err := provider.GetTransactions(
		context.Background(),
		"ethereum",
		"0xwallet",
		1,
		10,
	)

	require.Error(t, err)
}
