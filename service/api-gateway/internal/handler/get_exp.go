package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"api-gateway/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/golkity/Calc_2.0/service/calc-orchenstrator/repository"
)

type ExpressionGetter interface {
	One(ctx context.Context, id, userID int64) (repository.Expression, error)
}

func GetExpression(repo ExpressionGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exprID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		userID := middleware.UserID(r.Context())

		e, err := repo.One(r.Context(), exprID, userID)
		if err != nil {
			http.Error(w, "expression not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(struct {
			Expression repository.Expression `json:"expression"`
		}{e})
	}
}
