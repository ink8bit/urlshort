package app

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"urlshort/internal/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fakeOrigURL = "http://example.com"

func TestOriginURLHandler(t *testing.T) {
	type want struct {
		code           int
		urlPath        string
		contentType    string
		locationHeader string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Bad request",
			want: want{
				code:           400,
				urlPath:        "/",
				contentType:    "text/plain; charset=utf-8",
				locationHeader: "",
			},
		},
		{
			name: "Not found",
			want: want{
				code:           404,
				urlPath:        "/1",
				contentType:    "text/plain; charset=utf-8",
				locationHeader: "",
			},
		},
		{
			name: "Temporary Redirect",
			want: want{
				code:           307,
				urlPath:        "/1",
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
				findURL = db.FindURL
			}()
		}
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tt.want.urlPath, http.NoBody)
			w := httptest.NewRecorder()
			originURLHandler(w, r)

			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.locationHeader, res.Header.Get("Location"))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
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
				body:        "https://example.com",
				response:    "Bad Request",
				contentType: "text/plain",
			},
		},
		{
			name: "OK",
			want: want{
				code:        200,
				body:        "https://example.com",
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
			saveURL = func(origURL string) string {
				return "1"
			}
			defer func() {
				findURL = db.FindURL
				saveURL = db.SaveURL
			}()
		}
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/",
				bytes.NewBufferString(tt.want.body))
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
