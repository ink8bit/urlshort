package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"urlshort/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(conn string) (*Storage, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot create a connection to the databse: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, storage.ErrDBConnection
	}

	s := Storage{db: db}

	return &s, nil
}

// SaveURL saves original and shortened urls to the storage.
func (s *Storage) SaveURL(origURL string) (string, error) {
	// TODO: add logic
	return "id", nil
}

// FindURL returns original url by a given short one.
func (s *Storage) FindURL(shortURL string) (string, error) {
	// TODO: add logic
	return "origURL", nil
}

// FindShortURL returns short url by a given original url.
func (s *Storage) FindShortURL(origURL string) (string, error) {
	// TODO: add logic
	return "id", nil
}

func (s *Storage) Cleanup() error {
	// TODO: add logic
	return nil
}

func (s *Storage) Ping() error {
	err := s.db.Ping()
	if err != nil {
		return storage.ErrDBConnection
	}
	return nil
}
