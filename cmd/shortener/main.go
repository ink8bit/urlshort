package main

import (
	"fmt"
	"log"
	"net/http"

	"urlshort/internal/app/handlers/api"
	"urlshort/internal/app/handlers/ping"
	"urlshort/internal/app/handlers/redirect"
	"urlshort/internal/app/handlers/save"
	"urlshort/internal/storage"

	mwGzip "urlshort/internal/app/middleware/gzip"
	mwLog "urlshort/internal/app/middleware/logger"

	"urlshort/internal/config"
	"urlshort/internal/logger"
	"urlshort/internal/storage/memory"
	"urlshort/internal/storage/postgresql"

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

	// Init storage: memory or postgres
	var storage storage.Storager

	if cfg.DBCon != "" {
		storage, err = postgresql.New(cfg.DBCon)
		if err != nil {
			return fmt.Errorf(
				"cannot initialize db storage: %w", err)
		}
	} else {
		storage, err = memory.New(cfg.StoragePath)
		if err != nil {
			return fmt.Errorf(
				"cannot initialize memory storage: %w", err)
		}
	}

	// Init router
	r := chi.NewRouter()

	// Middleware
	r.Use(mwLog.Log(logger))
	r.Use(mwGzip.Compress())

	// Handlers
	r.Get("/{id:^[0-9A-Za-z]+$}", redirect.RedirectHandler(storage))
	r.Post("/", save.SaveURLHandler(cfg.BaseURL, storage))
	r.Post("/api/shorten", api.ShortenHandler(cfg.BaseURL, storage))
	r.Get("/ping", ping.DBConHandler(storage))

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("cannot start the server: %w", err)
	}
	return nil
}
