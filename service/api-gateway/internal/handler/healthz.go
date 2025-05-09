package handler

import (
	"net/http"

	"api-gateway/internal/kafka"
)

func Healthz(ping kafka.Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := ping.Ping(r.Context()); err != nil {
			http.Error(w, "kafka down", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
