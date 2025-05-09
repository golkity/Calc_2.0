package routes

import (
	h "auth_service/internal/http/handler"

	"github.com/go-chi/chi/v5"
)

func Register(r chi.Router, hnd *h.Handler) {
	r.Route("/api/v1", func(rt chi.Router) {
		rt.Post("/register", hnd.Register)
		rt.Post("/login", hnd.Login)
		rt.Post("/token/refresh", hnd.Refresh)
	})
}
