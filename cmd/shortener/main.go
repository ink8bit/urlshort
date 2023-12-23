package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"urlshort/internal/app"
	"urlshort/internal/app/router"
	"urlshort/internal/config"
	"urlshort/internal/storage/memory"

	"go.uber.org/zap"
)

func main() {
	config.ParseFlags()
	log.Fatal(run())
}

func run() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("cannot initialize zap logger: %w", err)
	}
	defer func() error {
		err := logger.Sync()
		if err != nil {
			return fmt.Errorf("logger error: %w", err)
		}
		return nil
	}()

	sugar := logger.Sugar()

	sugar.Info("Created a logger")

	// Init storage
	storage := memory.New()

	// Setup and run server
	server := app.NewServer(config.BaseURL, storage)

	// Init routes
	r := router.New(server)

	sugar.Infow(
		"Starting a server",
		"addr", strings.TrimPrefix(config.Addr, ":"),
	)

	err = http.ListenAndServe(config.Addr, r)
	if err != nil {
		return fmt.Errorf("can't start the server: %w", err)
	}

	return nil
}
