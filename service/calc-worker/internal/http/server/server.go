package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"calc-worker/internal/http/routes"
	kaf "calc-worker/internal/kafka"
	"calc-worker/internal/worker"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	http *http.Server
	stop context.CancelFunc
}

func New() (*Server, error) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	prod := kaf.NewProducer(brokers)

	wp := 1
	if v, _ := strconv.Atoi(os.Getenv("COMPUTING_POWER")); v > 0 {
		wp = v
	}

	svc := worker.New(prod)

	for i := 0; i < wp; i++ {
		go kaf.Consume(ctx, brokers, func(t kaf.Task) {
			svc.Handle(ctx, t)
		})
	}

	r := chi.NewRouter()
	routes.Register(r, ctx, prod)

	s := &http.Server{
		Addr:              ":8091",
		Handler:           r,
		ReadHeaderTimeout: 3 * time.Second,
	}

	return &Server{http: s, stop: func() {
		_ = prod.Close()
		stop()
	}}, nil
}

func (s *Server) Start() error {
	defer s.stop()
	log.Printf("calc-worker listening on %s", s.http.Addr)
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
