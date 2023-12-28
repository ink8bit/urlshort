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
	log.Fatal(run())
}

func run() error {
	// Init config
	cfg := config.Load()

	// TODO: Init logger

	// Init storage
	storage := memory.New()

	// Setup and run server
	server := app.NewServer(cfg.BaseURL, storage)

	// Init routes
	r := router.New(server)

	fmt.Println("Running server on address", cfg.Addr)
	err := http.ListenAndServe(cfg.Addr, r)
	if err != nil {
		return fmt.Errorf("can't start the server: %w", err)
	}
	return nil
}
