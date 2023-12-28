package config

import (
	"flag"
	"os"
)

var (
	addr    string
	baseURL string
)

type Config struct {
	Addr    string
	BaseURL string
}

func Load() *Config {
	var cfg Config

	flag.StringVar(&addr, "a", ":8080", "address and port to run server")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "short url prefix")
	flag.Parse()

	cfg.Addr = addr
	cfg.BaseURL = baseURL

	if env := os.Getenv("SERVER_ADDRESS"); env != "" {
		cfg.Addr = env
	}
	if env := os.Getenv("BASE_URL"); env != "" {
		cfg.BaseURL = env
	}

	return &cfg
}
