package postpg

import "github.com/chishkin-afk/posted/posts-service/internal/domain/post"

func ToModel(domain *post.Post) *PostModel {
	return &PostModel{
		ID:        domain.ID(),
		OwnerID:   domain.OwnerID(),
		Title:     domain.Title(),
		Body:      domain.Body(),
		PostedAt:  domain.PostedAt(),
		UpdatedAt: domain.UpdatedAt(),
	}
}

func ToDomain(model *PostModel) (*post.Post, error) {
	return post.From(
		model.ID,
		model.OwnerID,
		model.Title,
		model.Body,
		model.PostedAt,
		model.UpdatedAt,
	)
}
