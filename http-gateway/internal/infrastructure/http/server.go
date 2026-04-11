package httpserver

import (
	"context"
	"net/http"

	"github.com/chishkin-afk/posted/http-gateway/internal/infrastructure/config"
)

type Server struct {
	cfg *config.Config
	srv *http.Server
}

func New(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		cfg: cfg,
		srv: &http.Server{
			Addr:         cfg.Server.HTTP.Addr,
			Handler:      handler,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
	}
}

func (s *Server) Start() error {
	if s.cfg.Server.HTTP.TLS.Enable {
		return s.srv.ListenAndServeTLS(
			s.cfg.Server.HTTP.TLS.ServerCertPath,
			s.cfg.Server.HTTP.TLS.ServerKeyPath,
		)
	}

	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
