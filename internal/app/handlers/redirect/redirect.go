package redirect

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type URLRedirecter interface {
	FindURL(shortURL string) (string, error)
}

// RedirectHandler redirects the user to the original url
// after providing short url, or return "Not found".
func RedirectHandler(urlRedirecter URLRedirecter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		origURL, err := urlRedirecter.FindURL(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound),
				http.StatusNotFound)
			return
		}
		w.Header().Set("Location", origURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
