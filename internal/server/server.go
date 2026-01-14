package server

import (
	"github.com/shu-bham/go-url-shortener/internal/config"
	"github.com/shu-bham/go-url-shortener/internal/logger"
	"github.com/shu-bham/go-url-shortener/internal/shortener"
	"github.com/shu-bham/go-url-shortener/internal/storage"
	"github.com/shu-bham/go-url-shortener/internal/validator"
	"net/http"
	"os"

	"github.com/shu-bham/go-url-shortener/api"
)

func NewServer(handler *api.Handler, port string) *http.Server {
	mux := http.NewServeMux()
	router := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/shorten" {
			handler.ShortenURL(w, r)
			return
		}
		if r.URL.Path == "/" {
			http.NotFound(w, r)
			return
		}
		handler.RedirectURL(w, r)
	}
	mux.HandleFunc("/", router)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	return server
}

func StartServer() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(*cfg)

	db, err := storage.NewMySQLStorage(cfg.DB.DSN, log)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}

	shortenerSvc := shortener.NewShortener()
	validatorSvc := validator.NewValidator(db)
	handler := api.NewHandler(log, db, shortenerSvc, cfg.Server.DomainName, validatorSvc)
	server := NewServer(handler, cfg.Server.Port)

	log.Infof("Server is listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %s", err)
	}
}
