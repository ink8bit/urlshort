package storage

import (
	"errors"
	"strconv"
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

// SaveURL saves original url and id into two maps:
// encoded and decoded.
func SaveURL(origURL string) string {
	id := strconv.Itoa(len(shortUrls) + 1)
	shortUrls[origURL] = id
	origUrls[id] = origURL
	return id
}

// FindURL returns original url string by a given id.
func FindURL(id string) (string, error) {
	origURL, ok := origUrls[id]
	if !ok {
		return "", errors.New("url not found")
	}
	return origURL, nil
}

// FindID returns id by a given original url.
func FindID(origURL string) (string, error) {
	id, ok := shortUrls[origURL]
	if !ok {
		return "", errors.New("id not found")
	}
	return id, nil
}
