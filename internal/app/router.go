package app

import "github.com/go-chi/chi/v5"

func Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/{id:^[0-9A-Za-z]+}", originURLHandler)
	r.Post("/", shortURLHandler)
	return r
}
