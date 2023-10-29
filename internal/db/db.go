package db

import (
	"errors"
	"strconv"
)

// NOTE: it's a temporary solution.
//
// We don't have DB for our solution yet, so the best decision for now
// is to use two maps: encoded and decoded.
// Encoded map structure: website_url -> id_value, i.e. https://go.dev -> 1
// Decoded map structure: id_value -> website_url, i.e. 1 -> https://go.dev
// It gives us O(1) time complexity for searching original and decoded urls.
var (
	encoded = make(map[string]string)
	decoded = make(map[string]string)
)

// SaveURL saves original url and id into two maps:
// encoded and decoded.
func SaveURL(origURL string) string {
	id := strconv.Itoa(len(encoded) + 1)
	encoded[origURL] = id
	decoded[id] = origURL
	return id
}

// FindURL returns original url string by a given id.
func FindURL(id string) (string, error) {
	origURL, ok := decoded[id]
	if !ok {
		return "", errors.New("url not found")
	}
	return origURL, nil
}

// FindID returns id by a given original url.
func FindID(origURL string) (string, error) {
	id, ok := encoded[origURL]
	if !ok {
		return "", errors.New("id not found")
	}
	return id, nil
}
