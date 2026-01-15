package app

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/cache"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/config"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/portfolio"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/pricing"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/pricing/coingecko"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/pricing/mock"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/storage"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/transactions"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/transactions/etherscan"
)

type AppContext struct {
	Config *config.Config
	Logger *zap.Logger

	Cache              cache.CacheManager
	DB                 *pgxpool.Pool
	PricingService     *pricing.Service
	TransactionService *transactions.Service
	PortfolioService   portfolio.Service
}

func NewAppContext(ctx context.Context, cfg *config.Config, logger *zap.Logger, cache cache.CacheManager) (*AppContext, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	db, err := storage.NewPostgres(ctx, cfg.DB, logger)
	if err != nil {
		return nil, err
	}

	cgClient := coingecko.NewClient(
		cfg.CoinGecko.APIKey,
		cfg.CoinGecko.BaseURL,
	)

	primary := coingecko.NewProvider(cgClient)
	fallback := mock.NewProvider()

	pricingTTL := time.Duration(cfg.Pricing.CacheTTLSeconds) * time.Second

	pricingService := pricing.NewService(
		cache,
		primary,
		fallback,
		pricingTTL,
		logger,
	)

	etherscanClient := etherscan.NewClient(cfg.EtherScan.APIKey, cfg.EtherScan.BaseURL)
	txRepo := etherscan.NewProvider(etherscanClient)
	txService := transactions.NewService(txRepo)

	// Hard coded snapshot from requirement
	initial := []*portfolio.Portfolio{
		{
			Wallet: "0xabc123",
			Holdings: []portfolio.Holding{
				{
					Chain:           "ethereum",
					ContractAddress: "0x2260fac5e5542a773aa44fbcfedf7c193bc2c599", //btc
					Amount:          1.5,
				},
				{
					Chain:           "ethereum",
					ContractAddress: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", // USDC
					Amount:          500,
				},
			},
		},
	}

	repo := portfolio.NewMemoryRepository(initial)

	portfolioService := portfolio.NewService(repo, pricingService, logger)

	appCtx := &AppContext{
		Config:             cfg,
		Logger:             logger,
		Cache:              cache,
		DB:                 db,
		PricingService:     pricingService,
		TransactionService: txService,
		PortfolioService:   portfolioService,
	}

	return appCtx, nil
}

func (a *AppContext) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
