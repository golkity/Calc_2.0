package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	kaf "calc-worker/kafka"
	"calc-worker/worker"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	prod := kaf.NewProducer(brokers)
	defer prod.Close()

	svc := worker.New(prod)

	wp := 1
	if v, _ := strconv.Atoi(os.Getenv("COMPUTING_POWER")); v > 0 {
		wp = v
	}

	for i := 0; i < wp; i++ {
		go kaf.Consume(ctx, brokers, func(t kaf.Task) {
			svc.Handle(ctx, t)
		})
	}

	go func() {
		r := chi.NewRouter()
		r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
			if err := prod.Ping(ctx); err != nil {
				http.Error(w, "kafka down", http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusOK)
		})
		log.Fatal(http.ListenAndServe(":8091", r))
	}()

	log.Printf("calc-worker started (%d goroutines)", wp)
	<-ctx.Done()
	log.Println("calc-worker graceful shutdown")
}
