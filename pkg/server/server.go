package server

import (
	"context"
	"net/http"
	"time"
)

// HTTPServer wraps the built-in net/http server.
type HTTPServer struct {
	server *http.Server
}

// NewHTTPServer creates a new HTTPServer instance.
func NewHTTPServer(addr string, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

// Start runs the HTTP server.
func (s *HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully stops the HTTP server using the provided context.
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// ForceShutdown forcefully stops the HTTP server with a default timeout.
func (s *HTTPServer) ForceShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.Shutdown(ctx)
}
