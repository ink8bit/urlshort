package storage

import (
	"errors"
	"urlshort/internal/shorten"
)

// NOTE: it's a temporary solution.
//
// We don't have DB for our solution yet, so the best decision for now
// is to use two maps: shortUrls and origUrls.
// shortUrls map structure: website_url -> short_url, i.e. https://go.dev -> http://localhost/abcdefg
// origUrls map structure: short_url -> website_url, i.e. http://localhost/abcdefg -> https://go.dev
// It gives us O(1) time complexity for searching original and shortened urls.
var (
	shortUrls = make(map[string]string)
	origUrls  = make(map[string]string)
)

// SaveURL saves original and short urls to storage.
func SaveURL(origURL string) string {
	id := shorten.GenStr()
	shortUrls[origURL] = id
	origUrls[id] = origURL
	return id
}

// FindURL returns original url by a given short one.
func FindURL(shortURL string) (string, error) {
	origURL, ok := origUrls[shortURL]
	if !ok {
		return "", errors.New("original url not found")
	}
	return origURL, nil
}

// FindShortURL returns short url by a given original url.
func FindShortURL(origURL string) (string, error) {
	id, ok := shortUrls[origURL]
	if !ok {
		return "", errors.New("short url not found")
	}
	return id, nil
}
