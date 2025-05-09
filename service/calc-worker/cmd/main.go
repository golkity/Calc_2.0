package main

import (
	"calc-worker/internal/http/server"
	"log"
)

func main() {
	srv, err := server.New()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(srv.Start())
}
