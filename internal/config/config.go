package config

import (
	"flag"
	"os"
)

const (
	envServerAddress   = "SERVER_ADDRESS"
	envBaseURL         = "BASE_URL"
	envFileStoragePath = "FILE_STORAGE_PATH"
	envDBCon           = "DATABASE_DSN"
)

var (
	addr        string
	baseURL     string
	storagePath string
	dbCon       string
)

type Config struct {
	Addr        string
	BaseURL     string
	StoragePath string
	DBCon       string // DB connection string
}

func Load() *Config {
	var cfg Config

	flag.StringVar(&addr, "a", ":8080", "address and port to run server")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "short url prefix")
	flag.StringVar(&storagePath, "f", "/tmp/short-url-db.json", "file storage path")
	flag.StringVar(&dbCon, "d", "", "DB connection string")
	flag.Parse()

	cfg.Addr = addr
	cfg.BaseURL = baseURL
	cfg.StoragePath = storagePath
	cfg.StoragePath = storagePath
	cfg.DBCon = dbCon

	if env := os.Getenv(envServerAddress); env != "" {
		cfg.Addr = env
	}
	if env := os.Getenv(envBaseURL); env != "" {
		cfg.BaseURL = env
	}
	if env := os.Getenv(envFileStoragePath); env != "" {
		cfg.StoragePath = env
	}
	if env := os.Getenv(envDBCon); env != "" {
		cfg.DBCon = env
	}

	return &cfg
}
