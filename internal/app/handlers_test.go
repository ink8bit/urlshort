package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"urlshort/internal/storage/memory"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fakeURL = "a1b2c3"
	origURL = "https://go.dev"
)

type want struct {
	code           int
	url            string
	contentType    string
	locationHeader string
}

func TestOriginURLHandlerNotFound(t *testing.T) {
	tt := want{
		code:           404,
		url:            baseURL + "/1",
		contentType:    "text/plain; charset=utf-8",
		locationHeader: "",
	}

	storage := memory.New()
	server := NewServer(baseURL, storage)

	r := httptest.NewRequest(http.MethodGet, tt.url, http.NoBody)
	w := httptest.NewRecorder()

	server.OriginURLHandler(w, r)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, tt.code, res.StatusCode)
	assert.Equal(t, tt.locationHeader, res.Header.Get("Location"))
	assert.Equal(t, tt.contentType, res.Header.Get("Content-Type"))
}

type fakeMemory struct {
	shortUrls map[string]string
	origUrls  map[string]string
}

func NewFakeMemory() *fakeMemory {
	m := fakeMemory{
		shortUrls: make(map[string]string),
		origUrls:  make(map[string]string),
	}
	return &m
}

func (m *fakeMemory) SaveURL(u string) (string, error) {
	return fakeURL, nil
}

func (m *fakeMemory) FindURL(u string) (string, error) {
	return origURL, nil
}

func (m *fakeMemory) FindShortURL(u string) (string, error) {
	return fakeURL, nil
}

func (m *fakeMemory) Cleanup() error {
	m.shortUrls = make(map[string]string)
	m.origUrls = make(map[string]string)
	return nil
}

func TestOriginURLHandlerTemporaryRedirect(t *testing.T) {
	storage := NewFakeMemory()
	shortURL, err := storage.SaveURL(origURL)
	require.NoError(t, err)

	server := NewServer(baseURL, storage)

	t.Cleanup(func() {
		err := storage.Cleanup()
		require.NoError(t, err)
	})

	tt := want{
		code:        307,
		url:         baseURL + "/" + shortURL,
		contentType: "",
	}

	r := httptest.NewRequest(http.MethodGet, tt.url, http.NoBody)
	w := httptest.NewRecorder()

	server.OriginURLHandler(w, r)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, tt.code, res.StatusCode)
	assert.Equal(t, origURL, res.Header.Get("Location"))
	assert.Equal(t, tt.contentType, res.Header.Get("Content-Type"))
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

	storage := memory.New()
	server := NewServer(baseURL, storage)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/",
				strings.NewReader(tt.want.body))
			w := httptest.NewRecorder()

			server.ShortURLHandler(w, r)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)

			defer func() {
				_ = res.Body.Close()
			}()

			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.contentType,
				res.Header.Get("Content-Type"))
		})
	}
}
