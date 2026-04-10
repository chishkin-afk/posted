package post

import (
	"context"

	"github.com/google/uuid"
)

type PostPersistenceRepository interface {
	Save(ctx context.Context, post *Post) (*Post, error)
	Update(ctx context.Context, post *Post) (*Post, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Post, error)
	GetSelfPosts(ctx context.Context, userID uuid.UUID, page, limit uint32) ([]*Post, error)
}

type PostCacheRepository interface {
	Set(ctx context.Context, post *Post) error
	Get(ctx context.Context, id uuid.UUID) (*Post, error)
	Del(ctx context.Context, id uuid.UUID) error
}
