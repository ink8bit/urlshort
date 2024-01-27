package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"

	"urlshort/internal/shorten"
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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS urls(
	id SERIAL PRIMARY KEY,
	orig_url TEXT NOT NULL,
	short_url TEXT NOT NULL);
	`)

	if err != nil {
		return nil, fmt.Errorf("error while executing table creation: %w", err)
	}

	_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS orig_url_idx ON urls (orig_url)`)
	if err != nil {
		return nil, fmt.Errorf("error while creating an index in orig_url field: %w", err)
	}

	s := Storage{db: db}

	return &s, nil
}

var pqErr *pq.Error

// SaveURL saves original and shortened urls to the storage.
func (s *Storage) SaveURL(origURL string) (string, error) {
	id := shorten.GenRandomStr()
	stmt := `INSERT INTO urls (orig_url, short_url)
	VALUES ($1, $2)`
	_, err := s.db.Exec(stmt, origURL, id)
	if err != nil {
		if errors.As(err, &pqErr) {
			if pqErr.Code == pgerrcode.UniqueViolation {
				return "", storage.ErrOrigURLExists
			}
		}
		return "", fmt.Errorf(
			"error while inserting data to db: %w", err)
	}
	return id, nil
}

// FindURL returns original url by a given short one.
func (s *Storage) FindURL(shortURL string) (string, error) {
	stmt := `SELECT orig_url FROM urls WHERE short_url=$1;`
	var origURL string
	row := s.db.QueryRow(stmt, shortURL)
	err := row.Scan(&origURL)
	if err != nil {
		return "", storage.ErrDBNoRows
	}
	return origURL, nil
}

// FindShortURL returns short url by a given original url.
func (s *Storage) FindShortURL(origURL string) (string, error) {
	stmt := `SELECT short_url FROM urls WHERE orig_url=$1;`
	var shortURL string
	row := s.db.QueryRow(stmt, origURL)
	err := row.Scan(&shortURL)
	if err != nil {
		return "", storage.ErrDBNoRows
	}
	return shortURL, nil
}

func (s *Storage) Cleanup() error {
	stmt := `TRUNCATE TABLE urls`
	_, err := s.db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("error while truncating table: %w", err)
	}
	return nil
}

func (s *Storage) Ping() error {
	err := s.db.Ping()
	if err != nil {
		return storage.ErrDBConnection
	}
	return nil
}
