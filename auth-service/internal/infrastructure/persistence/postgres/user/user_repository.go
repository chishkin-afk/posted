package userpg

import (
	"context"
	"errors"
	"fmt"

	"github.com/chishkin-afk/posted/auth-service/internal/domain/user"
	"github.com/chishkin-afk/posted/auth-service/pkg/errs"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type userPersistenceRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *userPersistenceRepository {
	return &userPersistenceRepository{
		db: db,
	}
}

func (upr *userPersistenceRepository) Save(ctx context.Context, user *user.User) (*user.User, error) {
	model := ToModel(user)
	if err := upr.db.WithContext(ctx).Create(&model).Error; err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if upr.isUniqueError(err) {
			return nil, errs.ErrUserAlreadyExists
		}

		return nil, fmt.Errorf("failed to save user into db: %w", err)
	}

	return ToDomain(model)
}

func (upr *userPersistenceRepository) GetByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	var model UserModel
	if err := upr.db.WithContext(ctx).First(&model, "email = ?", email.String()).Error; err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserDoesntExist
		}

		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return ToDomain(&model)
}

func (upr *userPersistenceRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var model UserModel
	if err := upr.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserDoesntExist
		}

		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return ToDomain(&model)
}

func (upr *userPersistenceRepository) Update(ctx context.Context, user *user.User) (*user.User, error) {
	updates := map[string]any{
		"nickname": user.Nickname(),
		"email":    user.Email(),
	}

	if err := upr.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", user.ID()).Updates(updates).Error; err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if upr.isUniqueError(err) {
			return nil, errs.ErrUserAlreadyExists
		}

		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (upr *userPersistenceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := upr.db.WithContext(ctx).Delete(&UserModel{}, "id = ?", id)
	if err := result.Error; err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		return fmt.Errorf("failed to delete user from db: %w", err)
	}

	if result.RowsAffected == 0 {
		return errs.ErrUserDoesntExist
	}

	return nil
}

func (upr *userPersistenceRepository) isUniqueError(err error) bool {
	return errors.Is(postgres.Dialector{}.Translate(err), gorm.ErrDuplicatedKey)
}
