package api

import (
	"encoding/json"
	"fmt"
	"github.com/shu-bham/go-url-shortener/internal/shortener"
	"github.com/shu-bham/go-url-shortener/internal/storage"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	log        *logrus.Logger
	storage    storage.Storage
	shortener  shortener.URLShortener
	domainName string
}

func NewHandler(log *logrus.Logger, storage storage.Storage, shortener shortener.URLShortener, domainName string) *Handler {
	return &Handler{
		log:        log,
		storage:    storage,
		shortener:  shortener,
		domainName: domainName,
	}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	shortCode, err := h.shortener.GenerateShortURL()
	if err != nil {
		h.log.WithError(err).Error("Failed to generate short URL")
		http.Error(w, "Failed to generate short URL", http.StatusInternalServerError)
		return
	}

	if err := h.storage.SaveURL(ctx, req.URL, shortCode); err != nil {
		h.log.WithError(err).Error("Failed to save URL")
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	shortURL := fmt.Sprintf("%s/%s", h.domainName, shortCode)

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
