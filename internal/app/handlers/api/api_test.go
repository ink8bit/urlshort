package api_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"urlshort/internal/app/handlers/api"
	"urlshort/internal/storage/memory"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost"

type want struct {
	code        int
	body        string
	response    string
	contentType string
}

func prepareData(urlStr string) []byte {
	return []byte(fmt.Sprintf(`{ "url": %q }`, urlStr))
}

func TestShortenHandler(t *testing.T) {
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Bad request",
			want: want{
				code:        http.StatusBadRequest,
				body:        "",
				response:    "Bad Request",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Created",
			want: want{
				code:        http.StatusCreated,
				body:        "https://go.dev",
				response:    "Bad Request",
				contentType: "application/json",
			},
		},
	}

	storage := memory.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := prepareData(tt.want.body)

			r := httptest.NewRequest(
				http.MethodPost,
				"/api/shorten",
				bytes.NewBuffer(data),
			)
			w := httptest.NewRecorder()

			handler := api.ShortenHandler(baseURL, storage)
			handler(w, r)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)

			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.contentType,
				res.Header.Get("Content-Type"))
		})
	}
}
