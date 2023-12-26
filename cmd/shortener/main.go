package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"urlshort/internal/app"
	"urlshort/internal/app/logger"
	"urlshort/internal/app/router"
	"urlshort/internal/config"
	"urlshort/internal/storage/memory"
)

func main() {
	config.ParseFlags()
	log.Fatal(run())
}

func run() error {
	logger, err := logger.New()
	if err != nil {
		return fmt.Errorf("cannot initialize zap logger: %w", err)
	}

	logger.Info("Created a logger")

	// Init storage
	storage := memory.New()

	// Setup and run server
	server := app.NewServer(config.BaseURL, storage)

	// Init routes
	r := router.New(server)

	logger.Infow(
		"Starting a server",
		"addr", strings.TrimPrefix(config.Addr, ":"),
	)

	err = http.ListenAndServe(config.Addr, r)
	if err != nil {
		return fmt.Errorf("can't start the server: %w", err)
	}

	return nil
}
