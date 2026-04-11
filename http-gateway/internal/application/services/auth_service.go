package services

import (
	"context"
	"log/slog"

	authpb "github.com/chishkin-afk/posted/http-gateway/api/auth/v1"
	"github.com/chishkin-afk/posted/http-gateway/internal/application/dtos"
	"github.com/chishkin-afk/posted/http-gateway/internal/infrastructure/config"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/metadata"
)

type authService struct {
	cfg        *config.Config
	log        *slog.Logger
	authClient authpb.AuthServiceClient
	validate   *validator.Validate
}

func NewAuthService(
	cfg *config.Config,
	log *slog.Logger,
	authClient authpb.AuthServiceClient,
) *authService {
	return &authService{
		cfg:        cfg,
		log:        log,
		authClient: authClient,
		validate:   validator.New(),
	}
}

func (as *authService) Register(ctx context.Context, req *dtos.RegisterRequest) (*dtos.Token, error) {
	if err := as.validate.Struct(req); err != nil {
		return nil, err
	}

	resp, err := as.authClient.Register(ctx, &authpb.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
	})
	if err != nil {
		as.log.Warn("failed to register user",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &dtos.Token{
		Access:    resp.Access,
		AccessTTL: resp.AccessTtl,
	}, nil
}

func (as *authService) Login(ctx context.Context, req *dtos.LoginRequest) (*dtos.Token, error) {
	if err := as.validate.Struct(req); err != nil {
		return nil, err
	}

	resp, err := as.authClient.Login(ctx, &authpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		as.log.Warn("failed to login user",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &dtos.Token{
		Access:    resp.Access,
		AccessTTL: resp.AccessTtl,
	}, nil
}

func (as *authService) UpdateUser(ctx context.Context, req *dtos.UpdateUserRequest, token string) (*dtos.User, error) {
	ctx = as.ctxWithMD(ctx, token)

	resp, err := as.authClient.UpdateUser(ctx, &authpb.UpdateRequest{
		Updates: &authpb.Updates{
			Email:    req.Updates.Email,
			Nickname: req.Updates.Nickname,
		},
		UpdateMask: req.UpdateMask,
	})
	if err != nil {
		as.log.Warn("failed to update user",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &dtos.User{
		ID:        resp.Id,
		Email:     resp.Email,
		Nickname:  resp.Nickname,
		CreatedAt: resp.CreatedAt.AsTime().UnixMilli(),
		UpdatedAt: resp.UpdatedAt.AsTime().UnixMilli(),
	}, nil
}

func (as *authService) DeleteUser(ctx context.Context, token string) error {
	ctx = as.ctxWithMD(ctx, token)

	if _, err := as.authClient.DeleteUser(ctx, nil); err != nil {
		as.log.Warn("failed to delete user",
			slog.String("error", err.Error()),
		)
		return err
	}

	return nil
}

func (as *authService) GetUserSelf(ctx context.Context, token string) (*dtos.User, error) {
	ctx = as.ctxWithMD(ctx, token)

	resp, err := as.authClient.GetUserSelf(ctx, nil)
	if err != nil {
		as.log.Warn("failed to get self user",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &dtos.User{
		ID:        resp.Id,
		Email:     resp.Email,
		Nickname:  resp.Nickname,
		CreatedAt: resp.CreatedAt.AsTime().UnixMilli(),
		UpdatedAt: resp.UpdatedAt.AsTime().UnixMilli(),
	}, nil
}

func (as *authService) GetUserByID(ctx context.Context, id string) (*dtos.User, error) {
	resp, err := as.authClient.GetUserByID(ctx, &authpb.GetByIDRequest{
		Id: id,
	})
	if err != nil {
		as.log.Warn("failed to get user by id",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &dtos.User{
		ID:        resp.Id,
		Email:     resp.Email,
		Nickname:  resp.Nickname,
		CreatedAt: resp.CreatedAt.AsTime().UnixMilli(),
		UpdatedAt: resp.UpdatedAt.AsTime().UnixMilli(),
	}, nil
}

func (as *authService) ctxWithMD(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{
		"authorization": token,
	})

	return metadata.NewOutgoingContext(ctx, md)
}
