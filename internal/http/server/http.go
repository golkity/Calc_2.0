package server

import (
	"github.com/golkity/Calc_2.0/internal/http/handler"
	"github.com/golkity/Calc_2.0/internal/middleware"
	"github.com/golkity/Calc_2.0/pkg/logger"
	"net/http"
)

func RegRoutes(mux *http.ServeMux, log *logger.Logger) {
	h := handler.NewHandler(log)

	mux.Handle("/api/v1/calculate", middleware.LoggingMiddleware(http.HandlerFunc(h.RegRequest)))
	mux.Handle("/api/v1/expressions", middleware.LoggingMiddleware(http.HandlerFunc(h.ListExpressionsHandler)))
	mux.Handle("/api/v1/expressions/{id}", middleware.LoggingMiddleware(http.HandlerFunc(h.GetExpressionHandler)))

	mux.Handle("/internal/task", middleware.LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetTaskHandler(w, r)
		case http.MethodPost:
			h.PostTaskResultHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}
