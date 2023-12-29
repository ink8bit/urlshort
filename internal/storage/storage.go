package storage

import "errors"

var (
	ErrOrigURLNotFound  = errors.New("original url not found")
	ErrShortURLNotFound = errors.New("short url not found")
)
