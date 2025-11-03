package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Server represents the HTTP server
type Server struct {
	port       int
	httpServer *http.Server
}

// NewServer creates a new HTTP server instance
func NewServer(port int, handler http.Handler) *Server {
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		port:       port,
		httpServer: httpServer,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("🚀 HTTP server started on http://localhost:%d", s.port)

	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}
