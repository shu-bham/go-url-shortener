package storage

import (
	"github.com/shu-bham/go-url-shortener/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *MySQLStorage {
	cfg, err := config.LoadConfig("dev")
	require.NoError(t, err)

	storage, err := NewMySQLStorage(cfg.DB.DSN)
	require.NoError(t, err)

	_, err = storage.db.Exec("TRUNCATE TABLE urls")
	require.NoError(t, err)

	return storage
}

func TestMySQLStorage_SaveURL(t *testing.T) {
	storage := setupTestDB(t)

	longURL := "https://www.google.com"
	shortURL := "google"

	err := storage.SaveURL(longURL, shortURL)
	require.NoError(t, err)

	var retrievedLongURL string
	err = storage.db.QueryRow("SELECT long_url FROM urls WHERE short_url = ?", shortURL).Scan(&retrievedLongURL)
	require.NoError(t, err)
	assert.Equal(t, longURL, retrievedLongURL)
}

func TestMySQLStorage_GetURL(t *testing.T) {
	storage := setupTestDB(t)

	longURL := "https://www.example.com"
	shortURL := "example"

	_, err := storage.db.Exec("INSERT INTO urls (long_url, short_url) VALUES (?, ?)", longURL, shortURL)
	require.NoError(t, err)

	retrievedURL, err := storage.GetURL(shortURL)
	require.NoError(t, err)
	assert.Equal(t, longURL, retrievedURL)
}

func TestMySQLStorage_GetURL_NotFound(t *testing.T) {
	storage := setupTestDB(t)

	_, err := storage.GetURL("non-existent-url")
	assert.Error(t, err)
}

func TestMySQLStorage_DeleteURL(t *testing.T) {
	storage := setupTestDB(t)

	longURL := "https://www.yahoo.com"
	shortURL := "yahoo"

	_, err := storage.db.Exec("INSERT INTO urls (long_url, short_url) VALUES (?, ?)", longURL, shortURL)
	require.NoError(t, err)

	err = storage.DeleteURL(shortURL)
	require.NoError(t, err)

	var count int
	err = storage.db.QueryRow("SELECT COUNT(*) FROM urls WHERE short_url = ?", shortURL).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}
