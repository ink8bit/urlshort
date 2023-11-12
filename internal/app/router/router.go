package router

import (
	"net/http"

	"urlshort/internal/app"

	"github.com/go-chi/chi/v5"
)

func New(s *app.Server) http.Handler {
	r := chi.NewRouter()

	r.Get("/{id:^[0-9A-Za-z]+$}", s.OriginURLHandler)
	r.Post("/", s.ShortURLHandler)

	return r
}
