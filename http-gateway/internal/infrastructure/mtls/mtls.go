package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/chishkin-afk/posted/http-gateway/internal/infrastructure/config"
)

func LoadMTLS(cfg *config.MTLS) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(
		cfg.ClientCertPath,
		cfg.ClientKeyPath,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load x509 pair of server: %w", err)
	}

	caBytes, err := os.ReadFile(cfg.ClientCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read root ca path: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caBytes) {
		return nil, fmt.Errorf("failed to append root ca into pool: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caPool,
	}, nil
}
