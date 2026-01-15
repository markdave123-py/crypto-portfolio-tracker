package portfolio

import (
	"context"
	"fmt"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/pricing"
	"go.uber.org/zap"
)

type Service interface {
	Get(ctx context.Context, wallet string) (*PortfolioView, error)
	AddHolding(ctx context.Context, wallet string, h Holding) error
	UpdateHolding(ctx context.Context, wallet string, h Holding) error
	RemoveHolding(ctx context.Context, wallet string, chain string, contract string) error
}

type service struct {
	repo    Repository
	pricing pricing.ServiceAPI
	logger  *zap.Logger
}

func NewService(repo Repository, pricing pricing.ServiceAPI, logger *zap.Logger) Service {
	return &service{
		repo:    repo,
		pricing: pricing,
		logger:  logger,
	}
}

func (s *service) AddHolding(ctx context.Context, wallet string, h Holding) error {
	s.logger.Info("add-holding",
		zap.String("wallet", wallet),
		zap.String("chain", h.Chain),
		zap.String("contract", h.ContractAddress),
		zap.Float64("amount", h.Amount),
	)

	p, err := s.repo.Get(ctx, wallet)
	if err != nil {
		s.logger.Warn("portfolio-not-found-creating-new",
			zap.String("wallet", wallet),
		)
		p = &Portfolio{Wallet: wallet}
	}

	for _, existing := range p.Holdings {
		if existing.Chain == h.Chain && existing.ContractAddress == h.ContractAddress {
			s.logger.Warn("holding-already-exists",
				zap.String("wallet", wallet),
				zap.String("chain", h.Chain),
				zap.String("contract", h.ContractAddress),
			)
			return fmt.Errorf("holding already exists")
		}
	}

	p.Holdings = append(p.Holdings, h)
	return s.repo.Save(ctx, p)
}

func (s *service) UpdateHolding(ctx context.Context, wallet string, h Holding) error {
	s.logger.Info("update-holding",
		zap.String("wallet", wallet),
		zap.String("chain", h.Chain),
		zap.String("contract", h.ContractAddress),
		zap.Float64("amount", h.Amount),
	)

	p, err := s.repo.Get(ctx, wallet)
	if err != nil {
		s.logger.Error("portfolio-not-found",
			zap.String("wallet", wallet),
			zap.Error(err),
		)
		return err
	}

	found := false
	for i, existing := range p.Holdings {
		if existing.Chain == h.Chain && existing.ContractAddress == h.ContractAddress {
			p.Holdings[i].Amount = h.Amount
			found = true
			break
		}
	}

	if !found {
		s.logger.Warn("holding-not-found",
			zap.String("wallet", wallet),
			zap.String("chain", h.Chain),
			zap.String("contract", h.ContractAddress),
		)
		return fmt.Errorf("holding not found")
	}

	return s.repo.Save(ctx, p)
}

func (s *service) RemoveHolding(ctx context.Context, wallet, chain, contract string) error {
	s.logger.Info("remove-holding",
		zap.String("wallet", wallet),
		zap.String("chain", chain),
		zap.String("contract", contract),
	)

	p, err := s.repo.Get(ctx, wallet)
	if err != nil {
		s.logger.Error("portfolio-not-found",
			zap.String("wallet", wallet),
			zap.Error(err),
		)
		return err
	}

	out := make([]Holding, 0, len(p.Holdings))
	for _, h := range p.Holdings {
		if h.Chain == chain && h.ContractAddress == contract {
			continue
		}
		out = append(out, h)
	}

	p.Holdings = out
	return s.repo.Save(ctx, p)
}

func (s *service) Get(ctx context.Context, wallet string) (*PortfolioView, error) {
	s.logger.Info("get-portfolio",
		zap.String("wallet", wallet),
	)

	p, err := s.repo.Get(ctx, wallet)
	if err != nil {
		s.logger.Error("portfolio-not-found",
			zap.String("wallet", wallet),
			zap.Error(err),
		)
		return nil, err
	}

	refs := make([]pricing.AssetRef, 0, len(p.Holdings))
	for _, h := range p.Holdings {
		refs = append(refs, pricing.AssetRef{
			Chain:           h.Chain,
			ContractAddress: h.ContractAddress,
		})
	}

	prices, err := s.pricing.GetPrices(ctx, refs)
	fmt.Println(prices)
	if err != nil {
		s.logger.Error("pricing-failed",
			zap.String("wallet", wallet),
			zap.Error(err),
		)
		return nil, err
	}

	var total float64
	views := make([]HoldingView, 0, len(p.Holdings))

	for _, h := range p.Holdings {
		ref := pricing.AssetRef{
			Chain:           h.Chain,
			ContractAddress: h.ContractAddress,
		}

		price := prices[ref]
		value := price * h.Amount
		total += value

		views = append(views, HoldingView{
			Chain:           h.Chain,
			ContractAddress: h.ContractAddress,
			Amount:          h.Amount,
			PriceUSD:        price,
			ValueUSD:        value,
		})
	}

	s.logger.Info("portfolio-valued",
		zap.String("wallet", wallet),
		zap.Int("holdings", len(views)),
		zap.Float64("total_usd", total),
	)

	return &PortfolioView{
		Wallet:        wallet,
		Holdings:      views,
		TotalValueUSD: total,
	}, nil
}
