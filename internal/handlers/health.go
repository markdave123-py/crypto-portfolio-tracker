package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func HealthRoute(r chi.Router) {
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
}
