package services

import (
	"context"
	"errors"
	"log/slog"

	postspb "github.com/chishkin-afk/posted/posts-service/api/posts/v1"
	"github.com/chishkin-afk/posted/posts-service/internal/domain/post"
	"github.com/chishkin-afk/posted/posts-service/internal/domain/session"
	"github.com/chishkin-afk/posted/posts-service/internal/infrastructure/config"
	"github.com/chishkin-afk/posted/posts-service/pkg/errs"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct {
	cfg       *config.Config
	log       *slog.Logger
	postRepo  post.PostPersistenceRepository
	postCache post.PostCacheRepository
	jm        session.JWTManager
}

func New(
	cfg *config.Config,
	log *slog.Logger,
	postRepo post.PostPersistenceRepository,
	postCache post.PostCacheRepository,
	jm session.JWTManager,
) *service {
	return &service{
		cfg:       cfg,
		log:       log,
		postRepo:  postRepo,
		postCache: postCache,
		jm:        jm,
	}
}

func (s *service) Create(ctx context.Context, req *postspb.CreateRequest) (*postspb.Post, error) {
	userID, err := s.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	post, err := post.New(userID, req.Title, req.Body)
	if err != nil {
		return nil, err
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, s.cfg.Server.Timeout)
	defer cancel()

	savedPost, err := s.postRepo.Save(ctxTimeout, post)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			return nil, err
		}

		s.log.Error("failed to save post into db",
			slog.String("error", err.Error()),
			slog.String("post_id", post.ID().String()),
		)
		return nil, errs.ErrInternalServer
	}

	return &postspb.Post{
		Id:        savedPost.ID().String(),
		OwnerId:   savedPost.OwnerID().String(),
		Title:     savedPost.Title(),
		Body:      savedPost.Body(),
		PostedAt:  timestamppb.New(savedPost.PostedAt()),
		UpdatedAt: timestamppb.New(savedPost.UpdatedAt()),
	}, nil
}

func (s *service) Update(ctx context.Context, req *postspb.UpdateRequest) (*postspb.Post, error) {
	userID, err := s.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, errs.ErrInvalidPostID
	}

	if len(req.UpdateMask) > 2 {
		return nil, errs.ErrTooLargeUpdates
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, s.cfg.Server.Timeout)
	defer cancel()

	postToUpdate, err := s.postRepo.GetByID(ctxTimeout, postID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrPostDoesntExist) {
			return nil, err
		}

		s.log.Error("failed to get post by id",
			slog.String("error", err.Error()),
			slog.String("post_id", postID.String()),
		)
		return nil, errs.ErrInternalServer
	}

	if postToUpdate.OwnerID() != userID {
		return nil, errs.ErrNoEnoughRights
	}

	updates := make(map[string]bool, len(req.UpdateMask))
	for _, update := range req.UpdateMask {
		updates[update] = true
	}

	if updates["title"] {
		if err := postToUpdate.ChangeTitle(req.Updates.Title); err != nil {
			return nil, err
		}
	}
	if updates["body"] {
		if err := postToUpdate.ChangeBody(req.Updates.Body); err != nil {
			return nil, err
		}
	}

	updatedPost, err := s.postRepo.Update(ctxTimeout, postToUpdate)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrPostDoesntExist) {
			return nil, err
		}

		s.log.Error("failed to update post",
			slog.String("error", err.Error()),
			slog.String("post_id", postID.String()),
		)
		return nil, errs.ErrInternalServer
	}

	go s.savePostCache(context.Background(), updatedPost)
	return &postspb.Post{
		Id:        updatedPost.ID().String(),
		OwnerId:   updatedPost.OwnerID().String(),
		Title:     updatedPost.Title(),
		Body:      updatedPost.Body(),
		PostedAt:  timestamppb.New(updatedPost.PostedAt()),
		UpdatedAt: timestamppb.New(updatedPost.UpdatedAt()),
	}, nil
}

func (s *service) Delete(ctx context.Context, req *postspb.DeleteRequest) error {
	userID, err := s.getUserID(ctx)
	if err != nil {
		return err
	}

	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return errs.ErrInvalidPostID
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, s.cfg.Server.Timeout)
	defer cancel()

	postToDelete, err := s.postRepo.GetByID(ctxTimeout, postID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrPostDoesntExist) {
			return err
		}

		s.log.Error("failed to get post by id",
			slog.String("error", err.Error()),
			slog.String("post_id", postID.String()),
		)
		return errs.ErrInternalServer
	}

	if postToDelete.OwnerID() != userID {
		return errs.ErrNoEnoughRights
	}

	if err := s.postRepo.Delete(ctxTimeout, postToDelete.ID()); err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrPostDoesntExist) {
			return err
		}

		s.log.Error("failed to delete post",
			slog.String("error", err.Error()),
			slog.String("post_id", postID.String()),
		)
		return errs.ErrInternalServer
	}

	go s.deletePostCache(context.Background(), postToDelete.ID())
	return nil
}

func (s *service) GetByID(ctx context.Context, req *postspb.GetByIDRequest) (*postspb.Post, error) {
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, errs.ErrInvalidPostID
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, s.cfg.Server.Timeout)
	defer cancel()

	if post, err := s.postCache.Get(ctxTimeout, postID); err == nil {
		return s.returnPost(post), nil
	} else if !errors.Is(err, errs.ErrPostDoesntExist) {
		s.log.Warn("failed to get post from cache",
			slog.String("error", err.Error()),
			slog.String("post_id", postID.String()),
		)
	}

	post, err := s.postRepo.GetByID(ctxTimeout, postID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrPostDoesntExist) {
			return nil, err
		}

		s.log.Error("failed to get post by id",
			slog.String("error", err.Error()),
			slog.String("post_id", postID.String()),
		)
		return nil, errs.ErrInternalServer
	}

	go s.savePostCache(context.Background(), post)
	return s.returnPost(post), nil
}

func (s *service) GetSelfPosts(ctx context.Context, req *postspb.GetPostsRequest) (*postspb.Posts, error) {
	userID, err := s.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.Limit > 100 {
		return nil, errs.ErrInvalidSize
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, s.cfg.Server.Timeout)
	defer cancel()

	posts, err := s.postRepo.GetSelfPosts(ctxTimeout, userID, req.Page, req.Limit)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			return nil, err
		}

		s.log.Error("failed to get self posts",
			slog.String("error", err.Error()),
		)
		return nil, errs.ErrInternalServer
	}

	response := postspb.Posts{
		Posts: make([]*postspb.Post, len(posts)),
	}
	for idx, post := range posts {
		response.Posts[idx] = s.returnPost(post)
	}

	return &response, nil
}

func (s *service) returnPost(post *post.Post) *postspb.Post {
	return &postspb.Post{
		Id:        post.ID().String(),
		OwnerId:   post.OwnerID().String(),
		Title:     post.Title(),
		Body:      post.Body(),
		PostedAt:  timestamppb.New(post.PostedAt()),
		UpdatedAt: timestamppb.New(post.UpdatedAt()),
	}
}

func (s *service) getUserID(ctx context.Context) (uuid.UUID, error) {
	raw := ctx.Value(session.KeyUserID)
	if raw == nil {
		return uuid.Nil, errs.ErrInvalidToken
	}

	if id, ok := raw.(uuid.UUID); ok {
		return id, nil
	}

	return uuid.Nil, errs.ErrInvalidToken
}

func (s *service) savePostCache(ctx context.Context, post *post.Post) {
	ctxTimeout, cancel := context.WithTimeout(ctx, s.cfg.Server.Timeout)
	defer cancel()

	if err := s.postCache.Set(ctxTimeout, post); err != nil {
		s.log.Warn("failed to set post into cache",
			slog.String("error", err.Error()),
			slog.String("post_id", post.ID().String()),
		)
	}
}

func (s *service) deletePostCache(ctx context.Context, postID uuid.UUID) {
	ctxTimeout, cancel := context.WithTimeout(ctx, s.cfg.Server.Timeout)
	defer cancel()

	if err := s.postCache.Del(ctxTimeout, postID); err != nil {
		s.log.Warn("failed to delete post from cache",
			slog.String("error", err.Error()),
			slog.String("post_id", postID.String()),
		)
	}
}
