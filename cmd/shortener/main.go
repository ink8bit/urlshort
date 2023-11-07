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
	fmt.Println("Running server on address", config.Addr)
	err := http.ListenAndServe(config.Addr, app.Router())
	if err != nil {
		return fmt.Errorf("can't start the server: %w", err)
	}
	return nil
}
