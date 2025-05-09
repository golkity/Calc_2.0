package main

import (
	"auth_service/internal/http/server"
	"log"
)

func main() {
	srv, err := server.New()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(srv.Start())
}
