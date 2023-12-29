package save

import (
	"io"
	"net/http"
	"net/url"
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
		id, err := urlSaver.FindShortURL(origURL)
		if err != nil {
			id, err := urlSaver.SaveURL(origURL)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest),
					http.StatusBadRequest)
				return
			}
			shortURL := baseURL + "/" + id
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(shortURL))
			return
		}

		shortURL := baseURL + "/" + id
		_, _ = w.Write([]byte(shortURL))
	}
}
