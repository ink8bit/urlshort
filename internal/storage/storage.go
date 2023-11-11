package storage

type Storage interface {
	StorageSaver
	StorageFinder
}

type StorageSaver interface {
	SaveURL(origURL string) (string, error)
}

type StorageFinder interface {
	FindURL(shortURL string) (string, error)
	FindShortURL(origURL string) (string, error)
}
