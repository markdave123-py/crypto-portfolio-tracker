package transactions

import (
	"context"
)

type Repository interface {
	GetTransactions(
		ctx context.Context,
		chain string,
		wallet string,
		page int,
		limit int,
	) ([]Transaction, error)
}
