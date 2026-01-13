package storage

type Storage interface {
	SaveURL(longURL, shortURL string) error
	GetURL(shortURL string) (string, error)
	DeleteURL(shortURL string) error
}
