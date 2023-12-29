package redirect_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshort/internal/app/handlers/redirect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fakeURL = "a1b2c3"
	origURL = "https://go.dev"
	baseURL = "http://localhost"
)

type want struct {
	code           int
	url            string
	contentType    string
	locationHeader string
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

func TestRedirectHandler(t *testing.T) {
	storage := NewFakeMemory()
	shortURL, err := storage.SaveURL(origURL)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := storage.Cleanup()
		require.NoError(t, err)
	})

	tt := want{
		code:        http.StatusTemporaryRedirect,
		url:         baseURL + "/" + shortURL,
		contentType: "",
	}

	r := httptest.NewRequest(http.MethodGet, tt.url, http.NoBody)
	w := httptest.NewRecorder()

	handler := redirect.RedirectHandler(storage)
	handler(w, r)

	res := w.Result()

	defer res.Body.Close()

	assert.Equal(t, tt.code, res.StatusCode)
	assert.Equal(t, origURL, res.Header.Get("Location"))
	assert.Equal(t, tt.contentType, res.Header.Get("Content-Type"))
}
