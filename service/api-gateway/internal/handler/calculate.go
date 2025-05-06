package handler

import (
	"encoding/json"
	"net/http"

	"api-gateway/internal/kafka"
	"api-gateway/internal/middleware"

	"github.com/golkity/Calc_2.0/pkg/calc"
	orcrepo "github.com/golkity/Calc_2.0/service/calc-orchenstrator/repository"
)

type calcReq struct {
	Expression string `json:"expression"`
}

type calcResp struct {
	ID int64 `json:"id"`
}

func Calculate(prod *kafka.Producer, repo *orcrepo.ExpressionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req calcReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if _, err := calc.Calc(req.Expression); err != nil {
			http.Error(w, "invalid expression: "+err.Error(), http.StatusInternalServerError)
			return
		}

		uid := middleware.UserID(r.Context())
		id, err := repo.Create(r.Context(), uid, req.Expression)
		if err != nil {
			http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := prod.PublishCalcRequest(r.Context(), uid, req.Expression); err != nil {
			http.Error(w, "kafka error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(calcResp{ID: id})
	}
}
