package app

import (
	"io"
	"net/http"
	"net/url"

	"urlshort/internal/config"

	"github.com/go-chi/chi/v5"
)

// ShortURLHandler creates short url for a given full url string
// if no full url found in DB.
func (s *Server) ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}
	u, err := url.ParseRequestURI(string(bs))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}
	origURL := u.String()
	w.Header().Set("Content-Type", "text/plain")
	id, err := s.storage.FindShortURL(origURL)
	if err != nil {
		id, err := s.storage.SaveURL(origURL)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}
		shortURL := config.BaseURL + "/" + id
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(shortURL))
		return
	}

	shortURL := config.BaseURL + "/" + id
	_, _ = w.Write([]byte(shortURL))
}

// OriginURLHandler redirects the user to the original url
// after providing short url, or return "Not found".
func (s *Server) OriginURLHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	origURL, err := s.storage.FindURL(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound),
			http.StatusNotFound)
		return
	}
	w.Header().Set("Location", origURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
