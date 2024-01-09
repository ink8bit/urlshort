package memory

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"urlshort/internal/shorten"
	"urlshort/internal/storage"

	"github.com/google/uuid"
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
	file      *os.File
}

func New(filename string) (*Memory, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644) //nolint:gosec,gomnd // false positive
	if err != nil {
		return nil, fmt.Errorf("cannot create or open file: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("cannot get file info: %w", err)
	}

	shortUrls := make(map[string]string)
	origUrls := make(map[string]string)

	if stat.Size() == 0 {
		m := Memory{
			shortUrls: shortUrls,
			origUrls:  origUrls,
			file:      file,
		}
		return &m, nil
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var rec Record
		if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
			return nil, fmt.Errorf("error while unmarshaling record: %w", err)
		}
		shortUrls[rec.OriginalURL] = rec.ShortURL
		origUrls[rec.ShortURL] = rec.OriginalURL
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while reading file: %w", err)
	}

	m := Memory{
		shortUrls: shortUrls,
		origUrls:  origUrls,
		file:      file,
	}
	return &m, nil
}

// SaveURL saves original and shortened urls to the storage.
func (m *Memory) SaveURL(origURL string) (string, error) {
	id := shorten.GenRandomStr()
	uuidStr := uuid.New().String()
	m.shortUrls[origURL] = id
	m.origUrls[id] = origURL
	rec := Record{
		UUID:        uuidStr,
		ShortURL:    id,
		OriginalURL: origURL,
	}
	if err := rec.Save(m.file.Name()); err != nil {
		return "", fmt.Errorf("cannot save data to file: %w", err)
	}
	return id, nil
}

// FindURL returns original url by a given short one.
func (m *Memory) FindURL(shortURL string) (string, error) {
	origURL, ok := m.origUrls[shortURL]
	if !ok {
		return "", storage.ErrOrigURLNotFound
	}
	return origURL, nil
}

// FindShortURL returns short url by a given original url.
func (m *Memory) FindShortURL(origURL string) (string, error) {
	id, ok := m.shortUrls[origURL]
	if !ok {
		return "", storage.ErrShortURLNotFound
	}
	return id, nil
}

// Cleanup removes all records from in-memory storage.
func (m *Memory) Cleanup() error {
	m.shortUrls = make(map[string]string)
	m.origUrls = make(map[string]string)
	if err := os.Remove(m.file.Name()); err != nil {
		return fmt.Errorf("cannot remove file: %w", err)
	}
	return nil
}
