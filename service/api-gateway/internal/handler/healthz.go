package handler

import (
	"context"
	"net/http"

	"api-gateway/internal/kafka"
)

func Healthz(ping kafka.Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if err := ping.Ping(context.Background()); err != nil {
			http.Error(w, "kafka down", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
