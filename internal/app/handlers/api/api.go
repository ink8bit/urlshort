package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"urlshort/internal/app/middleware/auth"
	"urlshort/internal/storage"
)

const (
	contentType     = "Content-Type"
	applicationJSON = "application/json"
)

type Shortener interface {
	SaveURL(origURL string, userID int) (string, error)
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

		userID := auth.CheckAuth(r)

		id, err := shortener.SaveURL(origURL, userID)
		if err != nil {
			if errors.Is(err, storage.ErrOrigURLExists) {
				id, err := shortener.FindShortURL(origURL)
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
				w.Header().Set(contentType, applicationJSON)
				w.WriteHeader(http.StatusConflict)
				_, _ = w.Write(resp)
				return
			}
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

		w.Header().Set(contentType, applicationJSON)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(resp)
	}
}

type payloadUrls []payloadURL

type payloadURL struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}

type responseURL struct {
	ID  string `json:"correlation_id"`
	URL string `json:"short_url"`
}

func ShortenBatchHandler(baseURL string, shortener Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var urls payloadUrls
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(buf.Bytes(), &urls); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(urls) == 0 {
			http.Error(w, http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}

		var responseUrls []responseURL

		for _, rawURL := range urls {
			u, err := url.ParseRequestURI(rawURL.URL)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest),
					http.StatusBadRequest)
				return
			}

			origURL := u.String()
			userID := auth.CheckAuth(r)

			id, err := shortener.SaveURL(origURL, userID)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest),
					http.StatusBadRequest)
				return
			}

			shortURL := baseURL + "/" + id
			respURL := responseURL{
				ID:  rawURL.ID,
				URL: shortURL,
			}

			responseUrls = append(responseUrls, respURL)
		}

		resp, err := json.Marshal(responseUrls)
		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}

		w.Header().Set(contentType, applicationJSON)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(resp)
	}
}
