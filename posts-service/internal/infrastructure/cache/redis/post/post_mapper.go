package postredis

import (
	"encoding/json"

	"github.com/chishkin-afk/posted/posts-service/internal/domain/post"
)

func ToBytes(domain *post.Post) ([]byte, error) {
	bytes, err := json.Marshal(PostModel{
		ID:        domain.ID(),
		OwnerID:   domain.OwnerID(),
		Title:     domain.Title(),
		Body:      domain.Body(),
		PostedAt:  domain.PostedAt(),
		UpdatedAt: domain.UpdatedAt(),
	})
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func ToDomain(bytes []byte) (*post.Post, error) {
	var model PostModel
	if err := json.Unmarshal(bytes, &model); err != nil {
		return nil, err
	}

	return post.From(
		model.ID,
		model.OwnerID,
		model.Title,
		model.Body,
		model.PostedAt,
		model.UpdatedAt,
	)
}
