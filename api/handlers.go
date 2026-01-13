package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shu-bham/go-url-shortener/internal/shortener"
	"github.com/shu-bham/go-url-shortener/internal/storage"
	"github.com/shu-bham/go-url-shortener/internal/validator"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const maxRetries = 10

type Handler struct {
	log        *logrus.Logger
	storage    storage.Storage
	shortener  shortener.URLShortener
	domainName string
	validator  *validator.Validator
}

func NewHandler(log *logrus.Logger, storage storage.Storage, shortener shortener.URLShortener, domainName string, validator *validator.Validator) *Handler {
	return &Handler{
		log:        log,
		storage:    storage,
		shortener:  shortener,
		domainName: domainName,
		validator:  validator,
	}
}

func sanitizeURL(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "http://" + url
	}
	return url
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

	sanitizedURL := sanitizeURL(req.URL)

	var shortCode string
	var err error
	for i := 0; i < maxRetries; i++ {
		shortCode, err = h.shortener.GenerateShortURL()
		if err != nil {
			h.log.WithError(err).Error("Failed to generate short URL")
			http.Error(w, "Failed to generate short URL", http.StatusInternalServerError)
			return
		}
		isUnique, err := h.validator.IsShortURLUnique(ctx, shortCode)
		if err != nil {
			h.log.WithError(err).Error("Failed to validate short URL")
			http.Error(w, "Failed to validate short URL", http.StatusInternalServerError)
			return
		}
		if isUnique {
			break
		}
		if i == maxRetries-1 {
			err := errors.New("failed to generate unique short URL after max retries")
			h.log.WithError(err).Error("Failed to generate unique short URL")
			http.Error(w, "Failed to generate unique short URL", http.StatusInternalServerError)
			return
		}
	}

	if err := h.storage.SaveURL(ctx, sanitizedURL, shortCode); err != nil {
		h.log.WithError(err).Error("Failed to save URL")
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	shortURL := fmt.Sprintf("%s/%s", h.domainName, shortCode)

	h.log.WithFields(logrus.Fields{
		"long_url":  sanitizedURL,
		"short_url": shortURL,
	}).Info("Successfully shortened URL")

	res := ShortenResponse{ShortURL: shortURL}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		h.log.WithError(err).Error("Failed to encode response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	longURL, err := h.storage.GetURL(ctx, shortCode)
	if err != nil {
		h.log.WithError(err).WithField("short_code", shortCode).Error("Failed to get URL for redirection")
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	h.log.WithFields(logrus.Fields{
		"short_code": shortCode,
		"long_url":   longURL,
	}).Info("Redirecting to original URL")

	http.Redirect(w, r, longURL, http.StatusFound)
}
