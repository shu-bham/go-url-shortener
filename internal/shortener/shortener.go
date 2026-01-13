package shortener

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

type URLShortener interface {
	GenerateShortURL() (string, error)
}

type Shortener struct{}

func NewShortener() *Shortener {
	return &Shortener{}
}

func (s *Shortener) GenerateShortURL() (string, error) {
	salt := make([]byte, 8)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write(salt)
	hash := hasher.Sum(nil)

	shortURL := base64.URLEncoding.EncodeToString(hash[:6])

	return shortURL, nil
}
