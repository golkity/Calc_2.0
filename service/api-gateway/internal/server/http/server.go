package http

import (
	"context"
	"net/http"
	"os"
	"strings"

	"api-gateway/internal/kafka"
	"api-gateway/internal/server/routes"

	"github.com/golkity/Calc_2.0/pkg/tokens"
	orcrepo "github.com/golkity/Calc_2.0/service/calc-orchenstrator/repository"
	"github.com/jackc/pgx/v5"
)

type Server struct {
	http *http.Server
	prod *kafka.Producer
}

func New() (*Server, error) {
	prod := kafka.New(strings.Split(os.Getenv("KAFKA_BROKERS"), ","))
	pg, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	exprRepo := orcrepo.NewExpr(pg)
	tm := tokens.NewManager([]byte(os.Getenv("JWT_SECRET")), 0, 0)

	r := routes.RegRoutes(prod, exprRepo, tm)

	s := &Server{
		http: &http.Server{
			Addr:    ":8080",
			Handler: r,
		},
		prod: prod,
	}
	return s, nil
}

func (s *Server) Start() error {
	defer s.prod.Close()
	return s.http.ListenAndServe()
}
