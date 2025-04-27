package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"auth_service/internal/config"
	"auth_service/internal/db"
	h "auth_service/internal/http/handler"
	"auth_service/internal/kafka"
	"auth_service/internal/tokens"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	store, err := db.New(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}

	tokenMgr := tokens.NewManager(cfg.JWTSecret, cfg.AccessTTLSeconds, cfg.RefreshTTLSeconds)
	refreshStore := tokens.NewStore(cfg.RedisAddr)
	producer := kafka.New(cfg.KafkaBrokers)

	handler := h.New(store, tokenMgr, refreshStore, producer)

	r := chi.NewRouter()
	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)
	r.Post("/refresh", handler.Refresh)

	r.Route("/api", func(r chi.Router) {
		r.Use(handler.Auth)
		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			if id, ok := h.UserID(r.Context()); ok {
				w.Write([]byte(fmt.Sprintf("hello user %d", id)))
			}
		})
	})

	log.Printf("auth-service listening on :%s", cfg.HTTPPort)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTPPort, r))
}
