package app

import "urlshort/internal/storage"

const baseURL = "http://localhost:8080"

type Server struct {
	baseURL string
	storage storage.Storage
}

func NewServer(baseURL string, storage storage.Storage) *Server {
	s := Server{
		baseURL: baseURL,
		storage: storage,
	}
	return &s
}
