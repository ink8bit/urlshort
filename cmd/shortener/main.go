package main

import (
	"fmt"
	"log"
	"net/http"

	"urlshort/internal/app"
)

func main() {
	parseFlags()

	log.Fatal(run())
}

func run() error {
	fmt.Printf("Running server on address %s", addr)
	return http.ListenAndServe(addr, app.Router())
}
