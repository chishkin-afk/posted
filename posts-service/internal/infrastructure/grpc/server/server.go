package server

import (
	"fmt"
	"net"

	"github.com/chishkin-afk/posted/posts-service/internal/infrastructure/config"
	"google.golang.org/grpc"
)

type Server struct {
	cfg *config.Config
	srv *grpc.Server
}

func New(cfg *config.Config, srv *grpc.Server) *Server {
	return &Server{
		cfg: cfg,
		srv: srv,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.cfg.Server.GRPC.Addr)
	if err != nil {
		return fmt.Errorf("failed to open listener: %w", err)
	}

	if err := s.srv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve listener: %w", err)
	}

	return nil
}

func (s *Server) GracefulStop() {
	s.srv.GracefulStop()
}
