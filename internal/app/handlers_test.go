package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"urlshort/internal/db"

	"github.com/stretchr/testify/assert"
)

func TestShortURLHandler(t *testing.T) {
	t.SkipNow()
}

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
