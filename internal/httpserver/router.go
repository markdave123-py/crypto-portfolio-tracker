package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(pricesHandler *handlers.PricesHandler, txHandler *handlers.TransactionsHandler, portfolioHander *handlers.PortfolioHandler) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// health route
	handlers.HealthRoute(r)

	// price route
	r.Post("/prices", pricesHandler.GetPrices)

	r.Get("/wallets/{wallet}/transactions", txHandler.List)

	r.Route("/wallets/{wallet}/portfolio", func(r chi.Router) {
		r.Get("/", portfolioHander.Get)

		r.Route("/holdings", func(r chi.Router) {
			r.Post("/", portfolioHander.AddHolding)
			r.Put("/", portfolioHander.UpdateHolding)
			r.Delete("/", portfolioHander.RemoveHolding)
		})
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}
