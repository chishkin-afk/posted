package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"github.com/chishkin-afk/posted/auth-service/internal/infrastructure/config"
)

func LoadMTLSConfig(cfg *config.Config) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(
		cfg.Server.GRPC.MTLS.ServerCertPath,
		cfg.Server.GRPC.MTLS.ServerKeyPath,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load server cert & key: %w", err)
	}

	caBytes, err := os.ReadFile(cfg.Server.GRPC.MTLS.RootCAPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ca cert: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caBytes) {
		return nil, errors.New("failed to append ca cert into pool")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caPool,
	}, nil
}
