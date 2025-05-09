package routes

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func Register(r chi.Router, ctx context.Context, db *pgx.Conn) {
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		if err := db.Ping(ctx); err != nil {
			http.Error(w, "db down", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
