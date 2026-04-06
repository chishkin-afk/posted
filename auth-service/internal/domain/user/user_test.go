package user

import (
	"testing"
	"time"

	"github.com/chishkin-afk/posted/auth-service/pkg/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser_Success(t *testing.T) {
	email := Email("mail@example.com")
	password := "qwerty123"
	nickname := "nickname"

	now := time.Now().UTC()
	user, err := New(email, password, nickname)

	require.NoError(t, err)
	assert.NotEmpty(t, user.ID())
	assert.Equal(t, email, user.Email())
	assert.NotEqual(t, password, user.PasswordHash().String())
	assert.Equal(t, nickname, user.Nickname())
	assert.WithinDuration(t, now, user.CreatedAt(), 100*time.Millisecond)
	assert.WithinDuration(t, now, user.UpdatedAt(), 100*time.Millisecond)
}

func TestNewUser_Invalid(t *testing.T) {
	type data struct {
		email    Email
		password string
		nickname string
	}

	testCases := []struct {
		name     string
		input    data
		expected error
	}{
		{
			name: "empty_email",
			input: data{
				email:    "",
				password: "qwerty123",
				nickname: "nickname",
			},
			expected: errs.ErrInvalidEmail,
		},
		{
			name: "empty_password",
			input: data{
				email:    "mail@example.com",
				password: "",
				nickname: "nickname",
			},
			expected: errs.ErrInvalidPassword,
		},
		{
			name: "empty_nickname",
			input: data{
				email:    "mail@example.com",
				password: "qwerty123",
				nickname: "",
			},
			expected: errs.ErrInvalidNickname,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := New(tc.input.email, tc.input.password, tc.input.nickname)

			assert.EqualError(t, err, tc.expected.Error())
		})
	}
}
