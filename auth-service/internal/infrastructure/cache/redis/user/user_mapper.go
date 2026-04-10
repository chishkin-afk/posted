package userredis

import (
	"encoding/json"
	"fmt"

	"github.com/chishkin-afk/posted/auth-service/internal/domain/user"
)

func ToBytes(domain *user.User) ([]byte, error) {
	bytes, err := json.Marshal(UserModel{
		ID:        domain.ID(),
		Email:     domain.Email().String(),
		Nickname:  domain.Nickname(),
		CreatedAt: domain.CreatedAt(),
		UpdatedAt: domain.UpdatedAt(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal domain model: %w", err)
	}

	return bytes, nil
}

func ToDomain(bytes []byte) (*user.User, error) {
	var model UserModel
	if err := json.Unmarshal(bytes, &model); err != nil {
		return nil, fmt.Errorf("failed to parse bytes: %w", err)
	}

	return user.From(
		model.ID,
		user.Email(model.Email),
		"",
		model.Nickname,
		model.CreatedAt,
		model.UpdatedAt,
	)
}
