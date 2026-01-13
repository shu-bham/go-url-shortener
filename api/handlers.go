package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/shu-bham/go-url-shortener/internal/storage"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	log     *logrus.Logger
	storage storage.Storage
}

func NewHandler(log *logrus.Logger, storage storage.Storage) *Handler {
	return &Handler{
		log:     log,
		storage: storage,
	}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Error("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		h.log.Warn("URL is required but was not provided")
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortURL, err := h.generateShortURL()
	if err != nil {
		h.log.WithError(err).Error("Failed to generate short URL")
		http.Error(w, "Failed to generate short URL", http.StatusInternalServerError)
		return
	}

	if err := h.storage.SaveURL(req.URL, shortURL); err != nil {
		h.log.WithError(err).Error("Failed to save URL")
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	h.log.WithFields(logrus.Fields{
		"long_url":  req.URL,
		"short_url": shortURL,
	}).Info("Successfully shortened URL")

	res := ShortenResponse{ShortURL: shortURL}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		h.log.WithError(err).Error("Failed to encode response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) generateShortURL() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
