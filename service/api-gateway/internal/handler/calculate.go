package handler

import (
	"encoding/json"
	"net/http"

	"api-gateway/internal/kafka"
	"api-gateway/internal/middleware"
)

type calcRequest struct {
	Expression string `json:"expression"`
}

func Calculate(prod *kafka.Producer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req calcRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
			http.Error(w, "invalid body", http.StatusUnprocessableEntity)
			return
		}
		userID := middleware.UserID(r.Context())
		if err := prod.PublishCalcRequest(r.Context(), userID, req.Expression); err != nil {
			http.Error(w, "kafka error", http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
	}
}
