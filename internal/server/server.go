package server

import (
	"fmt"
	"net/http"

	"github.com/shu-bham/go-url-shortener/api"
)

// NewServer sets up and returns a new HTTP server.
func NewServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", api.ShortenURL)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return server
}

// StartServer starts the HTTP server.
func StartServer() {
	server := NewServer()

	fmt.Printf("Server is listening on %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
