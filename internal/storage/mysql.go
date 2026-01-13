package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLStorage struct {
	db *sql.DB
}

func NewMySQLStorage(dsn string) (*MySQLStorage, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &MySQLStorage{db: db}, nil
}

func (s *MySQLStorage) SaveURL(longURL, shortURL string) error {
	query := "INSERT INTO urls (long_url, short_url) VALUES (?, ?)"
	_, err := s.db.Exec(query, longURL, shortURL)
	if err != nil {
		return fmt.Errorf("failed to save URL: %w", err)
	}
	return nil
}

func (s *MySQLStorage) GetURL(shortURL string) (string, error) {
	var longURL string
	query := "SELECT long_url FROM urls WHERE short_url = ?"
	err := s.db.QueryRow(query, shortURL).Scan(&longURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("URL not found")
		}
		return "", fmt.Errorf("failed to get URL: %w", err)
	}
	return longURL, nil
}

func (s *MySQLStorage) DeleteURL(shortURL string) error {
	query := "DELETE FROM urls WHERE short_url = ?"
	_, err := s.db.Exec(query, shortURL)
	if err != nil {
		return fmt.Errorf("failed to delete URL: %w", err)
	}
	return nil
}
