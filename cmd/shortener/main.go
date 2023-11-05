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
	return http.ListenAndServe(config.Addr, app.Router())
}
