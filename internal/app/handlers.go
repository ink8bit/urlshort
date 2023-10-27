package app

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"urlshort/internal/db"
)

const baseURL = "http://localhost:8080"

// shortURLHandler creates short url for a given full url string
// if no full url found in DB.
func shortURLHandler(w http.ResponseWriter, r *http.Request) {
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	u, err := url.Parse(string(bs))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	origURL := u.String()
	w.Header().Set("Content-Type", "text/plain")

	id, err := db.FindID(origURL)
	if err != nil {
		id := db.Save(origURL)
		shortURL := fmt.Sprintf("%v/%v", baseURL, id)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURL))
		return
	}

	shortURL := fmt.Sprintf("%v/%v", baseURL, id)
	w.Write([]byte(shortURL))
}

// originURLHandler redirects the user to the original url
// after providing short url, or return "Not found".
func originURLHandler(w http.ResponseWriter, r *http.Request) {
	// Valid IDs are /1, /12, etc. therefore we should check
	// the length of url path less than 2 symbols
	urlPath := r.URL.Path
	if len(urlPath) < 2 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	key := urlPath[1:]

	origURL, err := db.FindURL(key)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Header().Set("Location", origURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// Mux is a multiplexer which routes every request to specific handlers.
// Allowed HTTP methods: GET, POST.
func Mux(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		originURLHandler(w, r)
	case http.MethodPost:
		shortURLHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
