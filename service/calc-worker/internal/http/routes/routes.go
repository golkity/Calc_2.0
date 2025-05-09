package routes

import (
	"context"
	"net/http"

	kaf "calc-worker/internal/kafka"

	"github.com/go-chi/chi/v5"
)

func Register(r chi.Router, ctx context.Context, prod *kaf.Producer) {
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		if err := prod.Ping(ctx); err != nil {
			http.Error(w, "kafka down", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
