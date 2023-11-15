package main

import (
	"log"
	"net/http"

	"github.com/gmalka/movers/transport/rest"
)

func main() {
	h := rest.NewHandler(nil, nil, nil, nil, rest.Log{
		Err: log.Default(),
		Inf: log.Default(),
	})

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: h.Init(),
	}

	server.ListenAndServe()
}
