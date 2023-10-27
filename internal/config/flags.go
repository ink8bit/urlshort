package config

import (
	"flag"
	"os"
)

var (
	Addr    string
	BaseURL string
)

func ParseFlags() {
	flag.StringVar(&Addr, "a", ":8080", "address and port to run server")
	flag.StringVar(&BaseURL, "b", "http://localhost:8080", "short url prefix")
	flag.Parse()

	if env := os.Getenv("SERVER_ADDRESS"); env != "" {
		Addr = env
	}
	if env := os.Getenv("BASE_URL"); env != "" {
		BaseURL = env
	}
}
