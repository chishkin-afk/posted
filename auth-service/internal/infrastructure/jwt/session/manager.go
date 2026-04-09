package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/chishkin-afk/posted/auth-service/internal/infrastructure/config"
	"github.com/chishkin-afk/posted/auth-service/pkg/errs"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtManager struct {
	public  *rsa.PublicKey
	private *rsa.PrivateKey
	cfg     *config.Config
}

func New(cfg *config.Config) (*jwtManager, error) {
	public, err := loadPublic(cfg.Session.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load public key: %w", err)
	}

	private, err := loadPrivate(cfg.Session.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	return &jwtManager{
		public:  public,
		private: private,
		cfg:     cfg,
	}, nil
}

func (jm *jwtManager) GenerateAccess(userID uuid.UUID) (string, error) {
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "posted-auth",
			Subject:   "posted-user",
			ExpiresAt: jwt.NewNumericDate(now.Add(jm.cfg.Session.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})

	signedToken, err := token.SignedString(jm.private)
	if err != nil {
		return "", fmt.Errorf("failed to generate access: %w", err)
	}

	return signedToken, nil
}

func (jm *jwtManager) GenerateRefresh() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return hex.EncodeToString(buf), nil
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

func loadPrivate(path string) (*rsa.PrivateKey, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPrivateKeyFromPEM(bytes)
}

func loadPublic(path string) (*rsa.PublicKey, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(bytes)
}
