package storage

import "context"

type Storage interface {
	SaveURL(ctx context.Context, longURL, shortURL string) error
	GetURL(ctx context.Context, shortURL string) (string, error)
	DeleteURL(ctx context.Context, shortURL string) error
}
