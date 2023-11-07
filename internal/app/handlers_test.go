package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"urlshort/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOriginURLHandler(t *testing.T) {
	fakeOrigURL := "http://example.com"
	type want struct {
		code           int
		url            string
		contentType    string
		locationHeader string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Not found",
			want: want{
				code:           404,
				url:            "http://localhost:8080/1",
				contentType:    "text/plain; charset=utf-8",
				locationHeader: "",
			},
		},
		{
			name: "Temporary Redirect",
			want: want{
				code:           307,
				url:            "http://localhost:8080/1",
				contentType:    "",
				locationHeader: fakeOrigURL,
			},
		},
	}
	for _, tt := range tests {
		if tt.want.code == http.StatusTemporaryRedirect {
			findURL = func(id string) (string, error) {
				return fakeOrigURL, nil
			}
			defer func() {
				findURL = storage.FindURL
			}()
		}
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tt.want.url,
				http.NoBody)
			w := httptest.NewRecorder()
			originURLHandler(w, r)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.locationHeader,
				res.Header.Get("Location"))
			assert.Equal(t, tt.want.contentType,
				res.Header.Get("Content-Type"))
		})
	}
}

func TestShortURLHandler(t *testing.T) {
	type want struct {
		code        int
		body        string
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Bad request",
			want: want{
				code:        400,
				body:        "invalid_url",
				response:    "Bad Request",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Created",
			want: want{
				code:        201,
				body:        "http://example.com",
				response:    "Bad Request",
				contentType: "text/plain",
			},
		},
		{
			name: "OK",
			want: want{
				code:        200,
				body:        "http://example.com",
				response:    "Bad Request",
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		if tt.want.code == http.StatusOK {
			findURL = func(id string) (string, error) {
				return "1", nil
			}
			saveURL = func(origURL string) (string, error) {
				return "1", nil
			}
			defer func() {
				findURL = storage.FindURL
				saveURL = storage.SaveURL
			}()
		}
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/",
				strings.NewReader(tt.want.body))
			w := httptest.NewRecorder()
			shortURLHandler(w, r)

			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.contentType,
				res.Header.Get("Content-Type"))
		})
	}
}
