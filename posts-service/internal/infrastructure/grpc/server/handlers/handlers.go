package handlers

import (
	"context"
	"errors"

	postspb "github.com/chishkin-afk/posted/posts-service/api/posts/v1"
	"github.com/chishkin-afk/posted/posts-service/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type postsService interface {
	Create(ctx context.Context, req *postspb.CreateRequest) (*postspb.Post, error)
	Update(ctx context.Context, req *postspb.UpdateRequest) (*postspb.Post, error)
	Delete(ctx context.Context, req *postspb.DeleteRequest) error
	GetByID(ctx context.Context, req *postspb.GetByIDRequest) (*postspb.Post, error)
	GetSelfPosts(ctx context.Context, req *postspb.GetPostsRequest) (*postspb.Posts, error)
}

type handlers struct {
	postspb.UnimplementedPostsServiceServer
	service postsService
}

func New(service postsService) *handlers {
	return &handlers{
		service: service,
	}
}

func (h *handlers) Create(ctx context.Context, req *postspb.CreateRequest) (*postspb.Post, error) {
	resp, err := h.service.Create(ctx, req)
	if err != nil {
		return nil, status.Error(h.getCode(err), err.Error())
	}

	return resp, nil
}

func (h *handlers) Delete(ctx context.Context, req *postspb.DeleteRequest) (*emptypb.Empty, error) {
	if err := h.service.Delete(ctx, req); err != nil {
		return nil, status.Error(h.getCode(err), err.Error())
	}

	return nil, nil
}

func (h *handlers) GetByID(ctx context.Context, req *postspb.GetByIDRequest) (*postspb.Post, error) {
	resp, err := h.service.GetByID(ctx, req)
	if err != nil {
		return nil, status.Error(h.getCode(err), err.Error())
	}

	return resp, nil
}

func (h *handlers) GetSelfPosts(ctx context.Context, req *postspb.GetPostsRequest) (*postspb.Posts, error) {
	resp, err := h.service.GetSelfPosts(ctx, req)
	if err != nil {
		return nil, status.Error(h.getCode(err), err.Error())
	}

	return resp, nil
}

func (h *handlers) Update(ctx context.Context, req *postspb.UpdateRequest) (*postspb.Post, error) {
	resp, err := h.service.Update(ctx, req)
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
	case errors.Is(err, errs.ErrInvalidTitle),
		errors.Is(err, errs.ErrInvalidBody),
		errors.Is(err, errs.ErrInvalidPostID),
		errors.Is(err, errs.ErrTooLargeUpdates),
		errors.Is(err, errs.ErrInvalidSize):
		return codes.InvalidArgument
	case errors.Is(err, errs.ErrPostDoesntExist):
		return codes.NotFound
	case errors.Is(err, errs.ErrInvalidToken):
		return codes.Unauthenticated
	case errors.Is(err, errs.ErrNoEnoughRights):
		return codes.PermissionDenied
	}

	return codes.Internal
}
