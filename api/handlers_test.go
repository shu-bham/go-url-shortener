package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenURL(t *testing.T) {
	t.Run("should return short url", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/shorten", strings.NewReader(`{"url":"https://www.google.com"}`))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ShortenURL)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		expected := `{"short_url":"http://short.url/12345"}`
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("should return error for invalid request body", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/shorten", strings.NewReader(`{"url":}`))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ShortenURL)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return error for empty url", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/shorten", strings.NewReader(`{"url":""}`))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ShortenURL)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
