package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"api-gateway/handler"
	"api-gateway/kafka"
	"api-gateway/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/golkity/Calc_2.0/pkg/tokens"
)

func main() {
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	prod := kafka.New(brokers)
	// exprRepo := repository.NewExpr(dbConn)
	defer prod.Close()

	tm := tokens.NewManager([]byte(os.Getenv("JWT_SECRET")), 0, 0)

	r := chi.NewRouter()
	r.Get("/healthz", handler.Healthz(prod))
	r.Group(func(rt chi.Router) {
		rt.Use(middleware.Auth(tm))
		rt.Post("/api/v1/calculate", handler.Calculate(prod))
		// rt.Get("/expressions/{id}", handler.GetExpression(exprRepo))
	})

	log.Println("api-gateway :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

//TODO: add expression get by ID!
