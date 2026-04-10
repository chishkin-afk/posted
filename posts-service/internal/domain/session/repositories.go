package session

import "github.com/google/uuid"

type JWTManager interface {
	Validate(tokenString string) (uuid.UUID, error)
}
