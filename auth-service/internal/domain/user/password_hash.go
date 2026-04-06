package user

import (
	"strings"

	"github.com/chishkin-afk/posted/auth-service/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHash string

func (ph PasswordHash) String() string {
	return string(ph)
}

func NewPasswordHash(password string) (PasswordHash, error) {
	password = strings.TrimSpace(password)

	n := len([]rune(password))
	if n < 6 || n > 64 {
		return "", errs.ErrInvalidPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errs.ErrInvalidPassword
	}

	return PasswordHash(hash), nil
}
