package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth_service/internal/config"
	"auth_service/internal/db"
	h "auth_service/internal/http/handler"
	"auth_service/internal/http/routes"
	"auth_service/internal/kafka"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/golkity/Calc_2.0/pkg/tokens"
)

type Server struct {
	http *http.Server
	stop context.CancelFunc
}

func New() (*Server, error) {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	store, err := db.New(ctx, cfg.PostgresDSN)
	if err != nil {
		stop()
		return nil, err
	}

	tokenMgr := tokens.NewManager(cfg.JWTSecret, cfg.AccessTTLSeconds, cfg.RefreshTTLSeconds)
	producer := kafka.New(cfg.KafkaBrokers)

	handler := h.New(store, tokenMgr, producer)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(5 * time.Second))
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	routes.Register(r, handler)

	addr := ":" + os.Getenv("HTTP_PORT")
	if addr == ":" {
		addr = ":8080"
	}

	s := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Server{http: s, stop: func() {
		_ = producer.Close(ctx)
		stop()
	}}, nil
}

func (s *Server) Start() error {
	defer s.stop()
	log.Printf("auth-service listening on %s", s.http.Addr)
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
