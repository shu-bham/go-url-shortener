package validator

import (
	"context"
	"github.com/shu-bham/go-url-shortener/internal/storage"
)

type Validator struct {
	storage storage.Storage
}

func NewValidator(storage storage.Storage) *Validator {
	return &Validator{storage: storage}
}

func (v *Validator) IsShortURLUnique(ctx context.Context, shortURL string) (bool, error) {
	_, err := v.storage.GetURL(ctx, shortURL)
	if err != nil {
		// Assuming GetURL returns an error when the URL is not found
		return true, nil
	}
	return false, nil
}
