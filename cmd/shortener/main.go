package main

import (
	"fmt"
	"log"
	"net/http"

	"urlshort/internal/app/handlers/redirect"
	"urlshort/internal/app/handlers/save"
	mwLog "urlshort/internal/app/middleware/logger"
	"urlshort/internal/config"
	"urlshort/internal/logger"
	"urlshort/internal/storage/memory"

	"github.com/go-chi/chi/v5"
)

func main() {
	log.Fatal(run())
}

func run() error {
	// Init config
	cfg := config.Load()

	// Init logger
	logger, err := logger.New()
	if err != nil {
		return fmt.Errorf("cannot create logger")
	}

	// Init storage
	storage := memory.New()

	// Init router
	r := chi.NewRouter()

	// Middleware
	r.Use(mwLog.Log(logger))

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
