package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/app"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/cache"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/config"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/handlers"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/httpserver"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/logger"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start portfolio API server",
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := logger.SetupLogger()
		if err != nil {
			fmt.Println("error-creating-logger")
			os.Exit(100)
		}

		cfg, err := config.Load()
		if err != nil {
			logger.Fatal("failed-to-load-config", zap.Error(err))
			os.Exit(100)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		redisClient, err := cache.NewRedisClient(cfg.Redis.URL)
		if err != nil {
			logger.Error("error-creating-cache-client", zap.Error(err))
			os.Exit(100)
		}

		cacheManager, err := cache.NewRedisManager(redisClient, logger)

		// AppContext
		appCtx, err := app.NewAppContext(ctx, cfg, logger, cacheManager)
		if err != nil {
			logger.Fatal("failed-to-create-app-context", zap.Error(err))
		}
		defer appCtx.Close()

		pricesHandler := handlers.NewPricesHandler(appCtx.PricingService, logger)

		txHandler := handlers.NewTransactionsHandler(appCtx.TransactionService)

		portfolioHander := handlers.NewPortfolioHandler(appCtx.PortfolioService, logger)

		router := httpserver.NewRouter(pricesHandler, txHandler, portfolioHander)

		go func() {
			if err := http.ListenAndServe(":8080", router); err != nil {
				logger.Fatal("http-server-failed", zap.Error(err))
			}
		}()

		logger.Info("portfolio-service-started")

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logger.Info("portfolio-service-stopping")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		_ = shutdownCtx
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
