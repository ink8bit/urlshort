package main

import (
	"log"
	"net/http"

	"urlshort/internal/app"

	"github.com/go-chi/chi/v5"
)

func main() {
	log.Fatal(run())
}

func run() error {
	r := chi.NewRouter()

	r.Get("/{id:^[0-9]}", app.OriginURLHandler)
	r.Post("/", app.ShortURLHandler)

	return http.ListenAndServe(":8080", r)
}
