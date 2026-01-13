package api

import (
	"context"
	"errors"
	"github.com/shu-bham/go-url-shortener/internal/shortener"
	"github.com/shu-bham/go-url-shortener/internal/validator"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockStorage is a mock of storage.Storage
type mockStorage struct {
	urls map[string]string
}

func (m *mockStorage) SaveURL(ctx context.Context, longURL, shortURL string) error {
	m.urls[shortURL] = longURL
	return nil
}

func (m *mockStorage) GetURL(ctx context.Context, shortURL string) (string, error) {
	longURL, ok := m.urls[shortURL]
	if !ok {
		return "", errors.New("not found")
	}
	return longURL, nil
}

func (m *mockStorage) DeleteURL(ctx context.Context, shortURL string) error {
	delete(m.urls, shortURL)
	return nil
}

func TestShortenURL(t *testing.T) {
	t.Run("should return short url", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/shorten", strings.NewReader(`{"url":"https://www.google.com"}`))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		log := logrus.New()
		storage := &mockStorage{urls: make(map[string]string)}
		s := shortener.NewShortener()
		v := validator.NewValidator(storage)
		handler := NewHandler(log, storage, s, "http://short.url", v)

		httpHandler := http.HandlerFunc(handler.ShortenURL)
		httpHandler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return error for invalid request body", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/shorten", strings.NewReader(`{"url":}`))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		log := logrus.New()
		storage := &mockStorage{urls: make(map[string]string)}
		s := shortener.NewShortener()
		v := validator.NewValidator(storage)
		handler := NewHandler(log, storage, s, "http://short.url", v)
		httpHandler := http.HandlerFunc(handler.ShortenURL)

		httpHandler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return error for empty url", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/shorten", strings.NewReader(`{"url":""}`))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		log := logrus.New()
		storage := &mockStorage{urls: make(map[string]string)}
		s := shortener.NewShortener()
		v := validator.NewValidator(storage)
		handler := NewHandler(log, storage, s, "http://short.url", v)
		httpHandler := http.HandlerFunc(handler.ShortenURL)
		httpHandler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestRedirectURL(t *testing.T) {
	t.Run("should redirect to long url", func(t *testing.T) {
		// Setup
		log := logrus.New()
		storage := &mockStorage{urls: make(map[string]string)}
		s := shortener.NewShortener()
		v := validator.NewValidator(storage)
		handler := NewHandler(log, storage, s, "http://short.url", v)

		// Create a short URL
		storage.SaveURL(context.Background(), "https://www.google.com", "12345")

		req, err := http.NewRequest("GET", "/12345", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		httpHandler := http.HandlerFunc(handler.RedirectURL)
		httpHandler.ServeHTTP(rr, req)

		// Assertions
		assert.Equal(t, http.StatusFound, rr.Code)
		assert.Equal(t, "https://www.google.com", rr.Header().Get("Location"))
	})

	t.Run("should return not found for invalid short url", func(t *testing.T) {
		// Setup
		log := logrus.New()
		storage := &mockStorage{urls: make(map[string]string)}
		s := shortener.NewShortener()
		v := validator.NewValidator(storage)
		handler := NewHandler(log, storage, s, "http://short.url", v)

		req, err := http.NewRequest("GET", "/invalid", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		httpHandler := http.HandlerFunc(handler.RedirectURL)
		httpHandler.ServeHTTP(rr, req)

		// Assertions
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}
