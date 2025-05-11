package handler

import (
	"api-gateway/internal/middleware"
	"context"
	"encoding/json"
	"net/http"

	"github.com/golkity/Calc_2.0/service/calc-orchenstrator/repository"
)

type ExpressionLister interface {
	List(ctx context.Context, userID int64) ([]repository.Expression, error)
}

func ListExpressions(repo ExpressionLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := middleware.UserID(r.Context())
		list, err := repo.List(r.Context(), uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(struct {
			Expressions []repository.Expression `json:"expressions"`
		}{list})
	}
}
