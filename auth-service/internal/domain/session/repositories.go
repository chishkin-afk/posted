package session

import "github.com/google/uuid"

type JWTManager interface {
	GenerateAccess(userID uuid.UUID) (string, error)
	Validate(tokenString string) (uuid.UUID, error)
}
