package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/handlers"
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
		r.Get("/holdings", portfolioHander.Get)
		r.Post("/holdings", portfolioHander.AddHolding)
		r.Put("/holdings", portfolioHander.UpdateHolding)
		r.Delete("/holdings", portfolioHander.RemoveHolding)
	})

	return r
}
