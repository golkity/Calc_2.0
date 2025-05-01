package handler

import (
	"context"
	"net/http"

	"api-gateway/kafka"
)

func Healthz(prod *kafka.Producer) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if err := prod.Ping(context.Background()); err != nil {
			http.Error(w, "kafka down", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
