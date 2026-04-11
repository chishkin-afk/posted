package services

import (
	"context"
	"log/slog"

	postspb "github.com/chishkin-afk/posted/http-gateway/api/posts/v1"
	"github.com/chishkin-afk/posted/http-gateway/internal/application/dtos"
	"github.com/chishkin-afk/posted/http-gateway/internal/infrastructure/config"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/metadata"
)

type postsService struct {
	cfg         *config.Config
	log         *slog.Logger
	postsClient postspb.PostsServiceClient
	validate    *validator.Validate
}

func NewPostsService(
	cfg *config.Config,
	log *slog.Logger,
	postsClient postspb.PostsServiceClient,
) *postsService {
	return &postsService{
		cfg:         cfg,
		log:         log,
		postsClient: postsClient,
		validate:    validator.New(),
	}
}

func (ps *postsService) Create(ctx context.Context, req *dtos.CreatePostRequest, token string) (*dtos.Post, error) {
	if err := ps.validate.Struct(req); err != nil {
		return nil, err
	}

	ctx = ps.ctxWithMD(ctx, token)

	resp, err := ps.postsClient.Create(ctx, &postspb.CreateRequest{
		Title: req.Title,
		Body:  req.Body,
	})
	if err != nil {
		ps.log.Warn("failed to create post",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &dtos.Post{
		ID:        resp.Id,
		OwnerID:   resp.OwnerId,
		Title:     resp.Title,
		Body:      resp.Body,
		PostedAt:  resp.PostedAt.AsTime().UnixMilli(),
		UpdatedAt: resp.UpdatedAt.AsTime().UnixMilli(),
	}, nil
}

func (ps *postsService) Update(ctx context.Context, req *dtos.UpdatePostRequest, token string) (*dtos.Post, error) {
	if err := ps.validate.Struct(req); err != nil {
		return nil, err
	}

	ctx = ps.ctxWithMD(ctx, token)

	resp, err := ps.postsClient.Update(ctx, &postspb.UpdateRequest{
		PostId: req.PostID,
		Updates: &postspb.Updates{
			Title: req.Updates.Title,
			Body:  req.Updates.Body,
		},
		UpdateMask: req.UpdateMask,
	})
	if err != nil {
		ps.log.Warn("failed to update post",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &dtos.Post{
		ID:        resp.Id,
		OwnerID:   resp.OwnerId,
		Title:     resp.Title,
		Body:      resp.Body,
		PostedAt:  resp.PostedAt.AsTime().UnixMilli(),
		UpdatedAt: resp.UpdatedAt.AsTime().UnixMilli(),
	}, nil
}

func (ps *postsService) Delete(ctx context.Context, id string, token string) error {
	ctx = ps.ctxWithMD(ctx, token)

	if _, err := ps.postsClient.Delete(ctx, &postspb.DeleteRequest{
		PostId: id,
	}); err != nil {
		ps.log.Warn("failed to delete post",
			slog.String("error", err.Error()),
		)
		return err
	}

	return nil
}

func (ps *postsService) GetByID(ctx context.Context, id string) (*dtos.Post, error) {
	resp, err := ps.postsClient.GetByID(ctx, &postspb.GetByIDRequest{
		PostId: id,
	})
	if err != nil {
		ps.log.Warn("failed to get post by id",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &dtos.Post{
		ID:        resp.Id,
		OwnerID:   resp.OwnerId,
		Title:     resp.Title,
		Body:      resp.Body,
		PostedAt:  resp.PostedAt.AsTime().UnixMilli(),
		UpdatedAt: resp.UpdatedAt.AsTime().UnixMilli(),
	}, nil
}

func (ps *postsService) GetSelfPosts(ctx context.Context, token string, page, limit uint32) (*dtos.Posts, error) {
	ctx = ps.ctxWithMD(ctx, token)

	resp, err := ps.postsClient.GetSelfPosts(ctx, &postspb.GetPostsRequest{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		ps.log.Warn("failed to get self posts",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	response := dtos.Posts{
		Posts: make([]dtos.Post, len(resp.Posts)),
	}
	for idx, post := range resp.Posts {
		response.Posts[idx] = dtos.Post{
			ID:        post.Id,
			OwnerID:   post.OwnerId,
			Title:     post.Title,
			Body:      post.Body,
			PostedAt:  post.PostedAt.AsTime().UnixMilli(),
			UpdatedAt: post.UpdatedAt.AsTime().UnixMilli(),
		}
	}

	return &response, nil
}

func (ps *postsService) ctxWithMD(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{
		"authorization": token,
	})

	return metadata.NewOutgoingContext(ctx, md)
}
