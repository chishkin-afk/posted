package services

import (
	"context"
	"errors"
	"log/slog"

	authpb "github.com/chishkin-afk/posted/auth-service/api/auth/v1"
	"github.com/chishkin-afk/posted/auth-service/internal/domain/session"
	"github.com/chishkin-afk/posted/auth-service/internal/domain/user"
	"github.com/chishkin-afk/posted/auth-service/internal/infrastructure/config"
	"github.com/chishkin-afk/posted/auth-service/pkg/errs"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type authService struct {
	cfg        *config.Config
	log        *slog.Logger
	userRepo   user.UserPersistenceRepository
	userCache  user.UserCacheRepository
	jwtManager session.JWTManager
}

func New(
	cfg *config.Config,
	log *slog.Logger,
	userRepo user.UserPersistenceRepository,
	userCache user.UserCacheRepository,
	jwtManager session.JWTManager,
) *authService {
	return &authService{
		cfg:        cfg,
		log:        log,
		userRepo:   userRepo,
		userCache:  userCache,
		jwtManager: jwtManager,
	}
}

func (as *authService) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.Token, error) {
	user, err := user.New(
		user.Email(req.Email),
		req.Password,
		req.Nickname,
	)
	if err != nil {
		return nil, err
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, as.cfg.Server.Timeout)
	defer cancel()

	savedUser, err := as.userRepo.Save(ctxTimeout, user)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrUserAlreadyExists) {
			return nil, err
		}

		as.log.Error("failed to save user",
			slog.String("error", err.Error()),
			slog.String("user_email", req.Email),
		)
		return nil, errs.ErrInternalServer
	}

	access, err := as.jwtManager.GenerateAccess(savedUser.ID())
	if err != nil {
		as.log.Error("failed to generate access token",
			slog.String("error", err.Error()),
			slog.String("user_id", user.ID().String()),
		)
		return nil, errs.ErrInternalServer
	}

	return &authpb.Token{
		Access:    access,
		AccessTtl: as.cfg.Session.AccessTTL.Milliseconds(),
	}, nil
}

func (as *authService) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.Token, error) {
	email := user.Email(req.Email)
	if !email.IsValid() {
		return nil, errs.ErrInvalidEmail
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, as.cfg.Server.Timeout)
	defer cancel()

	user, err := as.userRepo.GetByEmail(ctxTimeout, email)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrUserDoesntExist) {
			return nil, err
		}

		as.log.Error("failed to get user by email",
			slog.String("error", err.Error()),
			slog.String("user_email", req.Email),
		)
		return nil, errs.ErrInternalServer
	}

	if !user.PasswordHash().Compare(req.Password) {
		return nil, errs.ErrInvalidCredentials
	}

	access, err := as.jwtManager.GenerateAccess(user.ID())
	if err != nil {
		as.log.Error("failed to generate access token",
			slog.String("error", err.Error()),
			slog.String("user_id", user.ID().String()),
		)
		return nil, errs.ErrInternalServer
	}

	return &authpb.Token{
		Access:    access,
		AccessTtl: as.cfg.Session.AccessTTL.Milliseconds(),
	}, nil
}

func (as *authService) UpdateUser(ctx context.Context, req *authpb.UpdateRequest) (*authpb.User, error) {
	userID, err := as.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, as.cfg.Server.Timeout)
	defer cancel()

	userToUpdate, err := as.userRepo.GetByID(ctxTimeout, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrUserDoesntExist) {
			return nil, err
		}

		as.log.Error("failed to get user by id",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
		return nil, errs.ErrInternalServer
	}

	if err := as.applyUpdates(userToUpdate, req); err != nil {
		return nil, err
	}

	updatedUser, err := as.userRepo.Update(ctxTimeout, userToUpdate)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrUserDoesntExist) {
			return nil, err
		}

		as.log.Error("failed to update user",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
		return nil, errs.ErrInternalServer
	}

	go as.deleteFromCache(context.Background(), updatedUser.ID())
	return &authpb.User{
		Id:        updatedUser.ID().String(),
		Email:     updatedUser.Email().String(),
		Nickname:  updatedUser.Nickname(),
		CreatedAt: timestamppb.New(updatedUser.CreatedAt()),
		UpdatedAt: timestamppb.New(updatedUser.UpdatedAt()),
	}, nil
}

func (as *authService) DeleteUser(ctx context.Context) error {
	userID, err := as.getUserID(ctx)
	if err != nil {
		return err
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, as.cfg.Server.Timeout)
	defer cancel()

	if err := as.userRepo.Delete(ctxTimeout, userID); err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrUserDoesntExist) {
			return err
		}

		as.log.Error("failed to delete user",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
		return errs.ErrInternalServer
	}

	go as.deleteFromCache(context.Background(), userID)
	return nil
}

func (as *authService) GetUserSelf(ctx context.Context) (*authpb.User, error) {
	userID, err := as.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, as.cfg.Server.Timeout)
	defer cancel()

	if user, err := as.userCache.Get(ctxTimeout, userID); err == nil {
		return as.returnUser(user), nil
	} else if !errors.Is(err, errs.ErrUserDoesntExist) {
		as.log.Error("failed to get user from cache",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
	}

	user, err := as.userRepo.GetByID(ctxTimeout, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrUserDoesntExist) {
			return nil, err
		}

		as.log.Error("failed to get user by id",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
		return nil, errs.ErrInternalServer
	}

	go as.saveUserCache(context.Background(), user)
	return as.returnUser(user), nil
}

func (as *authService) GetUserByID(ctx context.Context, req *authpb.GetByIDRequest) (*authpb.User, error) {
	userID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, errs.ErrInvalidUserID
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, as.cfg.Server.Timeout)
	defer cancel()

	if user, err := as.userCache.Get(ctxTimeout, userID); err == nil {
		return as.returnUser(user), nil
	} else if !errors.Is(err, errs.ErrUserDoesntExist) {
		as.log.Error("failed to get user from cache",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
	}

	user, err := as.userRepo.GetByID(ctxTimeout, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrUserDoesntExist) {
			return nil, err
		}

		as.log.Error("failed to get user by id",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
		return nil, errs.ErrInternalServer
	}

	go as.saveUserCache(context.Background(), user)
	return as.returnUser(user), nil
}

func (as *authService) returnUser(user *user.User) *authpb.User {
	return &authpb.User{
		Id:        user.ID().String(),
		Email:     user.Email().String(),
		Nickname:  user.Nickname(),
		CreatedAt: timestamppb.New(user.CreatedAt()),
		UpdatedAt: timestamppb.New(user.UpdatedAt()),
	}
}

func (as *authService) getUserID(ctx context.Context) (uuid.UUID, error) {
	raw := ctx.Value(session.KeyUserID)
	if raw == nil {
		return uuid.Nil, errs.ErrInvalidToken
	}

	if id, ok := raw.(uuid.UUID); ok {
		return id, nil
	}

	return uuid.Nil, errs.ErrInvalidToken
}

func (as *authService) applyUpdates(userToUpdate *user.User, req *authpb.UpdateRequest) error {
	updates := make(map[string]bool, len(req.UpdateMask))
	for _, update := range req.UpdateMask {
		updates[update] = true
	}

	if updates["email"] {
		if err := userToUpdate.ChangeEmail(user.Email(req.Updates.Email)); err != nil {
			return err
		}
	}

	if updates["nickname"] {
		if err := userToUpdate.ChangeNickname(req.Updates.Nickname); err != nil {
			return err
		}
	}

	return nil
}

func (as *authService) deleteFromCache(ctx context.Context, id uuid.UUID) {
	ctxTimeout, cancel := context.WithTimeout(ctx, as.cfg.Server.Timeout)
	defer cancel()

	if err := as.userCache.Del(ctxTimeout, id); err != nil {
		as.log.Warn("failed to delete user from cache",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
	}
}

func (as *authService) saveUserCache(ctx context.Context, user *user.User) {
	ctxTimeout, cancel := context.WithTimeout(ctx, as.cfg.Server.Timeout)
	defer cancel()

	if err := as.userCache.Set(ctxTimeout, user); err != nil {
		as.log.Warn("failed to set user into cache",
			slog.String("error", err.Error()),
			slog.String("user_id", user.ID().String()),
		)
	}
}
