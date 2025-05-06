package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"auth_service/internal/config"
	"auth_service/internal/db"
	h "auth_service/internal/http/handler"
	"auth_service/internal/kafka"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/golkity/Calc_2.0/pkg/tokens"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	store, err := db.New(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}

	tokenMgr := tokens.NewManager(cfg.JWTSecret, cfg.AccessTTLSeconds, cfg.RefreshTTLSeconds)
	// refreshStore := tokens.NewStore(cfg.RedisAddr)
	producer := kafka.New(cfg.KafkaBrokers)

	handler := h.New(store, tokenMgr, producer)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(5 * time.Second))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/api/v1", func(rt chi.Router) {
		rt.Post("/register", handler.Register)
		rt.Post("/login", handler.Login)
		rt.Post("/token/refresh", handler.Refresh)
	})

	addr := ":" + os.Getenv("HTTP_PORT")
	if addr == ":" {
		addr = ":8080"
	}

	log.Printf("auth-service listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
