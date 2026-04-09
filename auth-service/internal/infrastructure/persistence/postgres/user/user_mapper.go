package userpg

import "github.com/chishkin-afk/posted/auth-service/internal/domain/user"

func ToModel(domain *user.User) *UserModel {
	return &UserModel{
		ID:           domain.ID(),
		Email:        domain.Email().String(),
		PasswordHash: domain.PasswordHash().String(),
		Nickname:     domain.Nickname(),
		CreatedAt:    domain.CreatedAt(),
		UpdatedAt:    domain.UpdatedAt(),
	}
}

func ToDomain(model *UserModel) (*user.User, error) {
	return user.From(
		model.ID,
		user.Email(model.Email),
		user.PasswordHash(model.PasswordHash),
		model.Nickname,
		model.CreatedAt,
		model.UpdatedAt,
	)
}
