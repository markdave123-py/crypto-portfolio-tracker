package transactions

import (
	"strings"

)

func detectDirection(tx Transaction, wallet string) Transaction {
	wallet = strings.ToLower(wallet)

	switch {
	case strings.ToLower(tx.From) == wallet:
		tx.Direction = DirectionOut
	case strings.ToLower(tx.To) == wallet:
		tx.Direction = DirectionIn
	}

	return tx
}
