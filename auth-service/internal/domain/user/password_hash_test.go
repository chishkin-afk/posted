package user

import (
	"testing"

	"github.com/chishkin-afk/posted/auth-service/pkg/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPasswordHash_Success(t *testing.T) {
	password := "qwerty123"
	passwordHash, err := NewPasswordHash(password)

	require.NoError(t, err)
	assert.NotEqual(t, password, passwordHash.String())
}

func TestNewPasswordHash_Invalid(t *testing.T) {
	empty := make([]rune, 128)

	testCases := []struct {
		name     string
		input    string
		expected error
	}{
		{
			name:     "empty_password",
			input:    "",
			expected: errs.ErrInvalidPassword,
		},
		{
			name:     "too_little_size",
			input:    "aaa",
			expected: errs.ErrInvalidPassword,
		},
		{
			name:     "too_big_size",
			input:    string(empty),
			expected: errs.ErrInvalidPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewPasswordHash(tc.input)

			require.Error(t, err)
			assert.EqualError(t, err, tc.expected.Error())
		})
	}
}
