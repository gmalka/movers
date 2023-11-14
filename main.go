package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gmalka/movers/transport/rest"
)

func main() {
	h := rest.NewHandler(nil, nil, nil, nil, rest.Log{})

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: h.Init(),
	}

	form := url.Values{}
	form.Add("login", "login")
	form.Add("password", "1234")

	fmt.Println(form.Encode())

	server.ListenAndServe()
}
