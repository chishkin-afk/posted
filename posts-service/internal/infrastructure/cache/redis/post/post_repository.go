package postredis

import (
	"context"
	"errors"
	"fmt"

	"github.com/chishkin-afk/posted/posts-service/internal/domain/post"
	"github.com/chishkin-afk/posted/posts-service/internal/infrastructure/config"
	"github.com/chishkin-afk/posted/posts-service/pkg/errs"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type postCacheRepository struct {
	cfg    *config.Config
	client *redis.Client
}

func New(cfg *config.Config, client *redis.Client) *postCacheRepository {
	return &postCacheRepository{
		cfg:    cfg,
		client: client,
	}
}

func (pcr *postCacheRepository) Set(ctx context.Context, post *post.Post) error {
	key := pcr.getKey(post.ID())
	bytes, err := ToBytes(post)
	if err != nil {
		return fmt.Errorf("failed to convert model into bytes: %w", err)
	}

	if err := pcr.client.Set(ctx, key, bytes, pcr.cfg.Cache.PostTTL).Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (pcr *postCacheRepository) Get(ctx context.Context, id uuid.UUID) (*post.Post, error) {
	key := pcr.getKey(id)

	bytes, err := pcr.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, redis.Nil) {
			return nil, errs.ErrPostDoesntExist
		}

		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return ToDomain(bytes)
}

func (pcr *postCacheRepository) Del(ctx context.Context, id uuid.UUID) error {
	key := pcr.getKey(id)
	if err := pcr.client.Del(ctx, key).Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		if errors.Is(err, redis.Nil) {
			return errs.ErrPostDoesntExist
		}

		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

func (pcr *postCacheRepository) getKey(id uuid.UUID) string {
	return fmt.Sprintf("post:%s", id.String())
}
