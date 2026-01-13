package server

import (
	"context"
	"github.com/shu-bham/go-url-shortener/api"
	"github.com/shu-bham/go-url-shortener/internal/shortener"
	"github.com/shu-bham/go-url-shortener/internal/validator"
	"github.com/sirupsen/logrus"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockStorage is a mock of storage.Storage
type mockStorage struct{}

func (m *mockStorage) SaveURL(ctx context.Context, longURL, shortURL string) error {
	return nil
}

func (m *mockStorage) GetURL(ctx context.Context, shortURL string) (string, error) {
	return "", nil
}

func (m *mockStorage) DeleteURL(ctx context.Context, shortURL string) error {
	return nil
}

func TestNewServer(t *testing.T) {
	t.Run("should return a new server", func(t *testing.T) {
		log := logrus.New()
		st := &mockStorage{}
		s := shortener.NewShortener()
		v := validator.NewValidator(st)
		handler := api.NewHandler(log, st, s, "http://short.url", v)
		server := NewServer(handler, ":8080")
		assert.NotNil(t, server)
		assert.NotNil(t, server.Handler)
	})
}
