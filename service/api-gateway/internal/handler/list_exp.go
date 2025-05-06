package handler

import (
	"api-gateway/internal/middleware"
	"encoding/json"
	"net/http"

	"github.com/golkity/Calc_2.0/service/calc-orchenstrator/repository"
)

func ListExpressions(repo *repository.ExpressionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := middleware.UserID(r.Context())
		list, err := repo.List(r.Context(), uid)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		_ = json.NewEncoder(w).Encode(struct {
			Expressions []repository.Expression `json:"expressions"`
		}{list})
	}
}
