package main

import (
	"fmt"
	"log"
	"net/http"

	"urlshort/internal/app/handlers/redirect"
	"urlshort/internal/app/handlers/save"
	"urlshort/internal/config"
	"urlshort/internal/storage/memory"

	"github.com/go-chi/chi/v5"
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

	// Init router
	r := chi.NewRouter()

	// Handlers
	r.Get("/{id:^[0-9A-Za-z]+$}", redirect.RedirectHandler(storage))
	r.Post("/", save.SaveURLHandler(cfg.BaseURL, storage))

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("cannot start the server: %w", err)
	}

	return nil
}
