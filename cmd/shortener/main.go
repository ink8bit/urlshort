package main

import (
	"log"
	"net/http"

	"urlshort/internal/app"
)

func main() {
	log.Fatal(run())
}

func run() error {
	return http.ListenAndServe(":8080", http.HandlerFunc(app.Mux))
}
