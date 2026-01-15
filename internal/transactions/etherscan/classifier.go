package etherscan

import (
	"math/big"
	"strings"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/transactions"
)

func weiToEther(value string) float64 {
	wei, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return 0
	}

	eth := new(big.Rat).SetFrac(
		wei,
		big.NewInt(1e18),
	)

	f, _ := eth.Float64()
	return f
}

func classifyType(item txListItem) transactions.TransactionType {
	fn := strings.ToLower(item.FunctionName)

	switch {
	case strings.Contains(fn, "swap"):
		return transactions.TypeSwap
	case strings.Contains(fn, "stake"), strings.Contains(fn, "deposit"), strings.Contains(fn, "delegate"):
		return transactions.TypeStake
	case item.Value != "0":
		return transactions.TypeSend
	default:
		return transactions.TypeReceive
	}
}

func classifyStatus(item txListItem) transactions.TransactionStatus {
	if item.IsError == "1" || item.TxReceiptStatus == "0" {
		return transactions.StatusFailed
	}
	return transactions.StatusSuccess
}
