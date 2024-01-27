package save

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"urlshort/internal/storage"
)

type URLSaver interface {
	SaveURL(origURL string) (string, error)
	FindShortURL(origURL string) (string, error)
}

// SaveURLHandler creates short url for a given full url string
// if no full url found in DB.
func SaveURLHandler(baseURL string, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		id, err := urlSaver.SaveURL(origURL)
		if err != nil {
			if errors.Is(err, storage.ErrOrigURLExists) {
				id, err := urlSaver.FindShortURL(origURL)
				if err != nil {
					http.Error(w, http.StatusText(http.StatusBadRequest),
						http.StatusBadRequest)
					return
				}
				shortURL := baseURL + "/" + id
				w.WriteHeader(http.StatusConflict)
				_, _ = w.Write([]byte(shortURL))
				return
			}
			http.Error(w, http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}

		shortURL := baseURL + "/" + id
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(shortURL))
	}
}
