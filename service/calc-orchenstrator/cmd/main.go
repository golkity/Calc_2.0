package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	kaf "calc-orchenstrator/kafka"
	orchestrator "calc-orchenstrator/orchenstrator"
	"calc-orchenstrator/repository"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	db, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	exprRepo := repository.NewExpr(db)
	prod := kaf.New(strings.Split(os.Getenv("KAFKA_BROKERS"), ","))

	svc := orchestrator.New(exprRepo, prod)

	go func() {
		_ = kaf.Consume(ctx, strings.Split(os.Getenv("KAFKA_BROKERS"), ","), func(t kaf.TaskReq) {
			svc.Handle(ctx, t.UserID, t.Expression)
		})
	}()

	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		if err := db.Ping(ctx); err != nil {
			http.Error(w, "db down", 503)
			return
		}
		w.WriteHeader(200)
	})
	log.Println("orchestrator :8090")
	log.Fatal(http.ListenAndServe(":8090", r))
}
