package etherscan

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/transactions"
)

func TestClassifyType_Swap(t *testing.T) {
	item := txListItem{
		FunctionName: "swapExactTokensForTokens",
		Value:        "0",
	}

	typ := classifyType(item)

	require.Equal(t, transactions.TypeSwap, typ)
}

func TestClassifyType_Stake(t *testing.T) {
	item := txListItem{
		FunctionName: "stake(uint256)",
	}

	typ := classifyType(item)

	require.Equal(t, transactions.TypeStake, typ)
}

func TestClassifyType_Send(t *testing.T) {
	item := txListItem{
		FunctionName: "",
		Value:        "1000000000000000000",
	}

	typ := classifyType(item)

	require.Equal(t, transactions.TypeSend, typ)
}

func TestClassifyStatus_Failed(t *testing.T) {
	item := txListItem{
		IsError:         "1",
		TxReceiptStatus: "0",
	}

	status := classifyStatus(item)

	require.Equal(t, transactions.StatusFailed, status)
}

func TestClassifyStatus_Success(t *testing.T) {
	item := txListItem{
		IsError:         "0",
		TxReceiptStatus: "1",
	}

	status := classifyStatus(item)

	require.Equal(t, transactions.StatusSuccess, status)
}
