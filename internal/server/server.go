package server

import (
	"github.com/shu-bham/go-url-shortener/internal/config"
	"github.com/shu-bham/go-url-shortener/internal/logger"
	"github.com/shu-bham/go-url-shortener/internal/shortener"
	"github.com/shu-bham/go-url-shortener/internal/storage"
	"net/http"
	"os"

	"github.com/shu-bham/go-url-shortener/api"
)

// NewServer sets up and returns a new HTTP server.
func NewServer(handler *api.Handler, port string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", handler.ShortenURL)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	return server
}

// StartServer starts the HTTP server.
func StartServer() {
	cfg, err := config.LoadConfig(os.Getenv("APP_ENV"))
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(*cfg)

	db, err := storage.NewMySQLStorage(cfg.DB.DSN, log)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}

	shortenerSvc := shortener.NewShortener()
	handler := api.NewHandler(log, db, shortenerSvc, cfg.Server.DomainName)
	server := NewServer(handler, cfg.Server.Port)

	log.Infof("Server is listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %s", err)
	}
}
