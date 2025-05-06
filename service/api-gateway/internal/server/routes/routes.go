package routes

import (
	"api-gateway/internal/handler"
	"api-gateway/internal/kafka"
	"api-gateway/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/golkity/Calc_2.0/pkg/tokens"
	orcrepo "github.com/golkity/Calc_2.0/service/calc-orchenstrator/repository"
)

func RegRoutes(prod *kafka.Producer, repo *orcrepo.ExpressionRepo, tm *tokens.Manager) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/healthz", handler.Healthz(prod))

	r.Group(func(rt chi.Router) {
		rt.Use(middleware.Auth(tm))

		rt.Post("/api/v1/calculate", handler.Calculate(prod))
		rt.Get("/api/v1/expressions", handler.ListExpressions(repo))
		rt.Get("/api/v1/expressions/{id}", handler.GetExpression(repo))
	})

	return r
}
