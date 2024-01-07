package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

type Shortener interface {
	SaveURL(origURL string) (string, error)
	FindShortURL(origURL string) (string, error)
}

type payload struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

func ShortenHandler(baseURL string, shortener Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload payload
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u, err := url.ParseRequestURI(payload.URL)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}

		origURL := u.String()
		id, err := shortener.SaveURL(origURL)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}

		shortURL := baseURL + "/" + id
		resp, err := json.Marshal(Response{Result: shortURL})
		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(resp)
	}
}
