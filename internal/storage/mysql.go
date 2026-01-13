package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type MySQLStorage struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewMySQLStorage(dsn string, log *logrus.Logger) (*MySQLStorage, error) {
	log.Info("Connecting to database...")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.WithError(err).Error("Failed to connect to database")
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.WithError(err).Error("Failed to ping database")
		return nil, err
	}

	log.Info("Successfully connected to database")
	return &MySQLStorage{db: db, log: log}, nil
}

func (s *MySQLStorage) SaveURL(longURL, shortURL string) error {
	s.log.WithFields(logrus.Fields{
		"long_url":  longURL,
		"short_url": shortURL,
	}).Info("Saving URL to database")
	query := "INSERT INTO urls (long_url, short_url) VALUES (?, ?)"
	_, err := s.db.Exec(query, longURL, shortURL)
	if err != nil {
		s.log.WithError(err).Error("Failed to save URL")
		return fmt.Errorf("failed to save URL: %w", err)
	}
	s.log.Info("Successfully saved URL")
	return nil
}

func (s *MySQLStorage) GetURL(shortURL string) (string, error) {
	s.log.WithField("short_url", shortURL).Info("Getting URL from database")
	var longURL string
	query := "SELECT long_url FROM urls WHERE short_url = ?"
	err := s.db.QueryRow(query, shortURL).Scan(&longURL)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.WithField("short_url", shortURL).Warn("URL not found")
			return "", fmt.Errorf("URL not found")
		}
		s.log.WithError(err).Error("Failed to get URL")
		return "", fmt.Errorf("failed to get URL: %w", err)
	}
	s.log.WithField("short_url", shortURL).Info("Successfully got URL")
	return longURL, nil
}

func (s *MySQLStorage) DeleteURL(shortURL string) error {
	s.log.WithField("short_url", shortURL).Info("Deleting URL from database")
	query := "DELETE FROM urls WHERE short_url = ?"
	_, err := s.db.Exec(query, shortURL)
	if err != nil {
		s.log.WithError(err).Error("Failed to delete URL")
		return fmt.Errorf("failed to delete URL: %w", err)
	}
	s.log.WithField("short_url", shortURL).Info("Successfully deleted URL")
	return nil
}
