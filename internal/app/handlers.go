package app

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"urlshort/internal/config"
	"urlshort/internal/storage"

	"github.com/go-chi/chi/v5"
)

var (
	// For testing purposes
	findURL = storage.FindURL
	saveURL = storage.SaveURL
)

// shortURLHandler creates short url for a given full url string
// if no full url found in DB.
func shortURLHandler(w http.ResponseWriter, r *http.Request) {
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
	id, err := storage.FindShortURL(origURL)
	if err != nil {
		id := saveURL(origURL)
		shortURL := fmt.Sprintf("%v/%v", config.BaseURL, id)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURL))
		return
	}

	shortURL := fmt.Sprintf("%v/%v", config.BaseURL, id)
	w.Write([]byte(shortURL))
}

// originURLHandler redirects the user to the original url
// after providing short url, or return "Not found".
func originURLHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	origURL, err := findURL(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound),
			http.StatusNotFound)
		return
	}
	w.Header().Set("Location", origURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
