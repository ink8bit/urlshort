package memory

import (
	"errors"

	"urlshort/internal/shorten"
)

// TODO: it's a temporary solution.
//
// We don't have DB for our solution yet, so the best decision for now
// is to use two maps: shortUrls and origUrls.
// shortUrls map structure: website_url -> short_url, i.e. https://go.dev -> http://localhost/abcdefg
// origUrls map structure: short_url -> website_url, i.e. http://localhost/abcdefg -> https://go.dev
// It gives us O(1) time complexity for searching original and shortened urls.

// Memory is an in-memory storage for
// original and shortened urls.
type Memory struct {
	shortUrls map[string]string
	origUrls  map[string]string
}

func New() *Memory {
	m := Memory{
		shortUrls: make(map[string]string),
		origUrls:  make(map[string]string),
	}
	return &m
}

// SaveURL saves original and shortened urls to the storage.
func (m *Memory) SaveURL(origURL string) (string, error) {
	id := shorten.GenRandomStr()
	m.shortUrls[origURL] = id
	m.origUrls[id] = origURL
	return id, nil
}

// FindURL returns original url by a given short one.
func (m *Memory) FindURL(shortURL string) (string, error) {
	origURL, ok := m.origUrls[shortURL]
	if !ok {
		return "", errors.New("original url not found")
	}
	return origURL, nil
}

// FindShortURL returns short url by a given original url.
func (m *Memory) FindShortURL(origURL string) (string, error) {
	id, ok := m.shortUrls[origURL]
	if !ok {
		return "", errors.New("short url not found")
	}
	return id, nil
}

// Cleanup removes all records from in-memory storage.
func (m *Memory) Cleanup() {
	m.shortUrls = make(map[string]string)
	m.origUrls = make(map[string]string)
}
