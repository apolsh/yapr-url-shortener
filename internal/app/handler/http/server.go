package http

import (
	"context"
	"crypto/tls"
	"net/http"
)

// HTTPServer HTTP имплементация интерфейса Server.
type HTTPServer struct {
	*http.Server
}

// NewHTTPServer возвращает http сервер.
func NewHTTPServer(s *http.Server) HTTPServer {
	return HTTPServer{s}
}

// Start запускает сервер.
func (s *HTTPServer) Start() error {
	return s.ListenAndServe()
}

// StartTLS запускает сервер с TLS шифрованием.
func (s *HTTPServer) StartTLS(cfg *tls.Config) error {
	s.TLSConfig = cfg
	return s.ListenAndServeTLS("", "")
}

// Stop останавливает сервер.
func (s *HTTPServer) Stop(ctx context.Context) error {
	s.SetKeepAlivesEnabled(false)
	return s.Shutdown(ctx)
}
