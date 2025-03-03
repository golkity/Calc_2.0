package server

import (
	"lms/internal/http/handler"
	"lms/pkg/logger"
	"net/http"
)

func RegRoutes(mux *http.ServeMux, log *logger.Logger) {
	h := handler.NewHandler(log)

	mux.HandleFunc("/api/v1/calculate", h.RegRequest)
	mux.HandleFunc("/api/v1/expressions", h.ListExpressionsHandler)
	mux.HandleFunc("/api/v1/expressions/{id}", h.GetExpressionHandler)

	mux.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:

			h.GetTaskHandler(w, r)
		case http.MethodPost:

			h.PostTaskResultHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
