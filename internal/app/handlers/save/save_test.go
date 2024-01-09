package save_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"urlshort/internal/app/handlers/save"
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

func TestSaveURLHandler(t *testing.T) {
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Bad request",
			want: want{
				code:        http.StatusBadRequest,
				body:        "invalid_url",
				response:    "Bad Request",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Created",
			want: want{
				code:        http.StatusCreated,
				body:        "http://example.com",
				response:    "Bad Request",
				contentType: "text/plain",
			},
		},
		{
			name: "OK",
			want: want{
				code:        http.StatusOK,
				body:        "http://example.com",
				response:    "Bad Request",
				contentType: "text/plain",
			},
		},
	}

	storage, _ := memory.New("records.json")

	t.Cleanup(func() {
		err := storage.Cleanup()
		require.NoError(t, err)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/",
				strings.NewReader(tt.want.body))
			w := httptest.NewRecorder()

			handler := save.SaveURLHandler(baseURL, storage)
			handler(w, r)

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
