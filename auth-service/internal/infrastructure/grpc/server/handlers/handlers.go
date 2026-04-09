package handlers

import (
	"context"
	"errors"

	authpb "github.com/chishkin-afk/posted/auth-service/api/auth/v1"
	"github.com/chishkin-afk/posted/auth-service/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type authService interface {
	Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.Token, error)
	Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.Token, error)
	UpdateUser(ctx context.Context, req *authpb.UpdateRequest) (*authpb.User, error)
	DeleteUser(ctx context.Context) error
	GetUserSelf(ctx context.Context) (*authpb.User, error)
	GetUserByID(ctx context.Context, req *authpb.GetByIDRequest) (*authpb.User, error)
}

type handlers struct {
	authpb.UnimplementedAuthServiceServer
	service authService
}

func New(service authService) *handlers {
	return &handlers{
		service: service,
	}
}

func (h *handlers) DeleteUser(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := h.service.DeleteUser(ctx); err != nil {
		return nil, status.Error(h.getCode(err), err.Error())
	}

	return nil, nil
}

func (h *handlers) GetUserByID(ctx context.Context, req *authpb.GetByIDRequest) (*authpb.User, error) {
	resp, err := h.service.GetUserByID(ctx, req)
	if err != nil {
		return nil, status.Error(h.getCode(err), err.Error())
	}

	return resp, nil
}

func (h *handlers) GetUserSelf(ctx context.Context, _ *emptypb.Empty) (*authpb.User, error) {
	resp, err := h.service.GetUserSelf(ctx)
	if err != nil {
		return nil, status.Error(h.getCode(err), err.Error())
	}

	return resp, nil
}

func (h *handlers) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.Token, error) {
	resp, err := h.service.Login(ctx, req)
	if err != nil {
		return nil, status.Error(h.getCode(err), err.Error())
	}

	return resp, nil
}

func (h *handlers) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.Token, error) {
	resp, err := h.service.Register(ctx, req)
	if err != nil {
		return nil, status.Error(h.getCode(err), err.Error())
	}

	return resp, nil
}

func (h *handlers) UpdateUser(ctx context.Context, req *authpb.UpdateRequest) (*authpb.User, error) {
	resp, err := h.UpdateUser(ctx, req)
	if err != nil {
		return nil, status.Error(h.getCode(err), err.Error())
	}

	return resp, nil
}

func (h *handlers) getCode(err error) codes.Code {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, errs.ErrUserDoesntExist):
		return codes.NotFound
	case errors.Is(err, errs.ErrUserAlreadyExists):
		return codes.AlreadyExists
	}

	return codes.Internal
}
