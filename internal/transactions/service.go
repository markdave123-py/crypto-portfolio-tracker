package transactions

import (
	"context"
	"strings"

	"go.uber.org/zap"
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
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) List(
	ctx context.Context,
	chain string,
	wallet string,
	page int,
	limit int,
	filters Filters,
) ([]Transaction, error) {
	s.logger.Info("list-transactions",
		zap.String("wallet", wallet),
	)
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
