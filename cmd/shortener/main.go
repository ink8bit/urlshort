package main

import (
	"fmt"
	"log"
	"net/http"

	"urlshort/internal/app"
	"urlshort/internal/config"
)

func main() {
	config.ParseFlags()

	log.Fatal(run())
}

func run() error {
	fmt.Printf("Running server on address %s", config.Addr)
	return http.ListenAndServe(config.Addr, app.Router())
}
