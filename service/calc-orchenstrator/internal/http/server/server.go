package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golkity/Calc_2.0/service/calc-orchenstrator/internal/http/routes"
	kaf "github.com/golkity/Calc_2.0/service/calc-orchenstrator/internal/kafka"
	orchestrator "github.com/golkity/Calc_2.0/service/calc-orchenstrator/internal/orchenstrator"
	"github.com/golkity/Calc_2.0/service/calc-orchenstrator/internal/repository"

	"github.com/jackc/pgx/v5"
)

type Server struct {
	http *http.Server
	stop context.CancelFunc
}

func New() (*Server, error) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	db, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		stop()
		return nil, err
	}

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	prod := kaf.New(brokers)

	exprRepo := repository.NewExpr(db)
	svc := orchestrator.New(exprRepo, prod)

	go func() {
		_ = kaf.Consume(ctx, brokers, func(t kaf.TaskReq) {
			svc.Handle(ctx, t.UserID, t.Expression)
		})
	}()

	r := chi.NewRouter()
	routes.Register(r, ctx, db)

	s := &http.Server{
		Addr:              ":8090",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Server{http: s, stop: func() {
		_ = prod.Close()
		_ = db.Close(ctx)
		stop()
	}}, nil
}

func (s *Server) Start() error {
	defer s.stop()
	log.Printf("orchestrator listening on %s", s.http.Addr)
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
