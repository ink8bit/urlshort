package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"urlshort/internal/storage"
)

const createQuery = `
CREATE TABLE IF NOT EXISTS urls(
	id INTEGER PRIMARY KEY,
	orig_url TEXT NOT NULL,
	short_url TEXT NOT NULL);
`

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

	stmt, err := db.Prepare(createQuery)
	if err != nil {
		return nil, fmt.Errorf("error while preparing table creation: %w", err)
	}
	defer stmt.Close() //nolint:errcheck // self-explanatory

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("error while executing table creation: %w", err)
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
