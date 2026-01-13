package shortener

import (
	"crypto/rand"
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
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
