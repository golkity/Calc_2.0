package main

import (
	"log"

	"api-gateway/internal/server/http"
)

func main() {
	srv, err := http.New()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(srv.Start())
}
