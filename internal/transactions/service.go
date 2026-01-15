package transactions

import (
	"context"
	"strings"
)

type ServiceAPI interface {
	List(
		ctx context.Context,
		chain string,
		wallet string,
		page int,
		limit int,
		filters Filters,
	) ([]Transaction, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(
	ctx context.Context,
	chain string,
	wallet string,
	page int,
	limit int,
	filters Filters,
) ([]Transaction, error) {

	txs, err := s.repo.GetTransactions(ctx, chain, wallet, page, limit)
	if err != nil {
		return nil, err
	}

	wallet = strings.ToLower(wallet)

	out := make([]Transaction, 0)

	for _, tx := range txs {
		tx = detectDirection(tx, wallet)

		if applyFilters(tx, filters) {
			out = append(out, tx)
		}
	}

	return out, nil
}
