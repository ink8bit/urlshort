package main

import (
	"fmt"
	"log"
	"net/http"

	"urlshort/internal/app"
	"urlshort/internal/app/router"
	"urlshort/internal/config"
	"urlshort/internal/storage/memory"
)

func main() {
	config.ParseFlags()
	log.Fatal(run())
}

func run() error {
	// Init storage
	storage := memory.New()

	// Setup and run server
	server := app.NewServer(config.BaseURL, storage)

	// Init routes
	r := router.New(server)

	fmt.Println("Running server on address", config.Addr)
	err := http.ListenAndServe(config.Addr, r)
	if err != nil {
		return fmt.Errorf("can't start the server: %w", err)
	}
	return nil
}
