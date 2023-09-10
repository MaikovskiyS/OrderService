package server

import (
	"context"
	"net/http"
	"orderservice/internal/config"
	"time"
)

const (
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 5 * time.Second
)

type Server struct {
	server *http.Server
}

// Constructor
func New(handler http.Handler, cfg *config.Config) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         cfg.OrderService.Port(),
	}

	s := &Server{
		server: httpServer,
	}
	return s
}

// Start listen 8080 port
func (s *Server) Run() {
	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			return
		}
	}()
}

// Shutdown server
func (s *Server) Stop() error {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		return err
	}
	return nil
}
