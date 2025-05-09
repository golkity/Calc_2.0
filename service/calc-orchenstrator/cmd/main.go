package main

import (
	"log"

	"github.com/golkity/Calc_2.0/service/calc-orchenstrator/internal/http/server"
)

func main() {
	srv, err := server.New()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(srv.Start())
}
