package user

import (
	"context"

	"github.com/google/uuid"
)

type UserPersistenceRepository interface {
	Save(ctx context.Context, user *User) (*User, error)
	GetByEmail(ctx context.Context, email Email) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserCacheRepository interface {
	Set(ctx context.Context, user *User) error
	Get(ctx context.Context, id uuid.UUID) (*User, error)
	Del(ctx context.Context, id uuid.UUID) error
}
