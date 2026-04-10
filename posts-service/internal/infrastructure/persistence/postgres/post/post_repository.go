package postpg

import (
	"context"
	"errors"
	"fmt"

	"github.com/chishkin-afk/posted/posts-service/internal/domain/post"
	"github.com/chishkin-afk/posted/posts-service/pkg/errs"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postPersistenceRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *postPersistenceRepository {
	return &postPersistenceRepository{
		db: db,
	}
}

func (ppr *postPersistenceRepository) Save(ctx context.Context, post *post.Post) (*post.Post, error) {
	model := ToModel(post)
	if err := ppr.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		return nil, fmt.Errorf("failed to save post into db: %w", err)
	}

	return ToDomain(model)
}

func (ppr *postPersistenceRepository) Update(ctx context.Context, post *post.Post) (*post.Post, error) {
	updates := map[string]any{
		"title": post.Title(),
		"body":  post.Body(),
	}

	if err := ppr.db.WithContext(ctx).Model(&PostModel{}).Where("id = ?", post.ID()).Updates(updates).Error; err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if ppr.isUniqueError(err) {
			return nil, errs.ErrPostDoesntExist
		}

		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return post, nil
}

func (ppr *postPersistenceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := ppr.db.WithContext(ctx).Delete(&PostModel{}, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) || errors.Is(result.Error, context.Canceled) {
			return result.Error
		}

		return fmt.Errorf("failed to delete post: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errs.ErrPostDoesntExist
	}

	return nil
}

func (ppr *postPersistenceRepository) GetByID(ctx context.Context, id uuid.UUID) (*post.Post, error) {
	var model PostModel
	if err := ppr.db.WithContext(ctx).Take(&model, id).Error; err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if ppr.isUniqueError(err) {
			return nil, errs.ErrPostDoesntExist
		}

		return nil, fmt.Errorf("failed to get post by id: %w", err)
	}

	return ToDomain(&model)
}

func (ppr *postPersistenceRepository) GetSelfPosts(ctx context.Context, userID uuid.UUID, page, limit uint32) ([]*post.Post, error) {
	offset := limit * (page - 1)
	var models []PostModel
	if err := ppr.db.WithContext(ctx).
		Where("owner_id = ?", userID).
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&models).Error; err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		return nil, fmt.Errorf("failed to get self posts: %w", err)
	}

	result := make([]*post.Post, len(models))
	for idx, model := range models {
		post, err := ToDomain(&model)
		if err != nil {
			return nil, err
		}

		result[idx] = post
	}

	return result, nil
}

func (ppr *postPersistenceRepository) isUniqueError(err error) bool {
	return errors.Is(postgres.Dialector{}.Translate(err), gorm.ErrRecordNotFound)
}
