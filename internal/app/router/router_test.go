package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"urlshort/internal/app"
	"urlshort/internal/storage/memory"

	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080"

func TestRouterGet(t *testing.T) {
	storage := memory.New()
	server := app.NewServer(baseURL, storage)
	r := New(server)

	srv := httptest.NewServer(r)
	defer srv.Close()

	tests := []struct {
		name   string
		code   int
		urlStr string
	}{
		{name: "Method Not Allowed", code: 405, urlStr: srv.URL + "/"},
		{name: "Not Found", code: 404, urlStr: srv.URL + "/1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := http.Get(tt.urlStr)
			require.NoError(t, err)

			defer func() {
				_ = res.Body.Close()
			}()

			if res.StatusCode != tt.code {
				t.Errorf("expected status %q; got %v",
					tt.name, res.Status)
			}
		})
	}
}

func TestRouterPost(t *testing.T) {
	storage := memory.New()
	server := app.NewServer(baseURL, storage)
	r := New(server)

	srv := httptest.NewServer(r)
	defer srv.Close()

	t.Cleanup(func() {
		err := storage.Cleanup()
		require.NoError(t, err)
	})

	tests := []struct {
		name string
		code int
		body string
	}{
		{
			name: "Created",
			code: 201,
			body: "https://go.dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := http.Post(srv.URL, "text/plain", strings.NewReader(tt.body))
			require.NoError(t, err)

			defer func() {
				_ = res.Body.Close()
			}()

			if res.StatusCode != tt.code {
				t.Errorf("expected status %q; got %v",
					tt.name, res.Status)
			}
		})
	}
}
