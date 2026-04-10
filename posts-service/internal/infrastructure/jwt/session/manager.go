package jwt

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/chishkin-afk/posted/posts-service/internal/infrastructure/config"
	"github.com/chishkin-afk/posted/posts-service/pkg/errs"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtManager struct {
	public *rsa.PublicKey
	cfg    *config.Config
}

func New(cfg *config.Config) (*jwtManager, error) {
	public, err := loadPublic(cfg.Session.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load public key: %w", err)
	}

	return &jwtManager{
		public: public,
		cfg:    cfg,
	}, nil
}

func (jm *jwtManager) Validate(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errs.ErrInvalidToken
		}

		return jm.public, nil
	})
	if err != nil {
		return uuid.Nil, errs.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims.UserID, nil
	}

	return uuid.Nil, errs.ErrInvalidToken
}

func loadPublic(path string) (*rsa.PublicKey, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(bytes)
}
