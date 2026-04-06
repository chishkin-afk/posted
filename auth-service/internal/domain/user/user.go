package user

import (
	"strings"
	"time"

	"github.com/chishkin-afk/posted/auth-service/pkg/errs"
	"github.com/google/uuid"
)

type User struct {
	id           uuid.UUID
	email        Email
	passwordHash PasswordHash
	nickname     string
	createdAt    time.Time
	updatedAt    time.Time
}

func New(
	email Email,
	password string,
	nickname string,
) (*User, error) {
	nickname = strings.TrimSpace(nickname)
	if len([]rune(nickname)) < 3 || len([]rune(nickname)) > 64 {
		return nil, errs.ErrInvalidNickname
	}

	if !email.IsValid() {
		return nil, errs.ErrInvalidEmail
	}

	passwordHash, err := NewPasswordHash(password)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &User{
		id:           uuid.New(),
		email:        email,
		passwordHash: passwordHash,
		nickname:     nickname,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func (u *User) ChangeEmail(email Email) error {
	if !email.IsValid() {
		return errs.ErrInvalidEmail
	}

	u.email = email
	u.updatedAt = time.Now().UTC()

	return nil
}

func (u *User) ChangeNickname(nickname string) error {
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		return errs.ErrInvalidNickname
	}

	u.nickname = nickname
	u.updatedAt = time.Now().UTC()

	return nil
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) PasswordHash() PasswordHash {
	return u.passwordHash
}

func (u *User) Nickname() string {
	return u.nickname
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}
