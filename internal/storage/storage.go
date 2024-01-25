package storage

import "errors"

var (
	ErrOrigURLNotFound  = errors.New("original url not found")
	ErrShortURLNotFound = errors.New("short url not found")
	ErrDBConnection     = errors.New("cannot connect to the database")
	ErrDBNoRows         = errors.New("no rows returned")
)

type Storager interface {
	SaveURL(origURL string) (string, error)
	FindURL(shortURL string) (string, error)
	FindShortURL(origURL string) (string, error)
	Cleanup() error
	Ping() error
}
