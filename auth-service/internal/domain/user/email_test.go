package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmail_Success(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "basic",
			input: "mail@example.com",
		},
		{
			name:  "with_spaces",
			input: "    mail@example.com   ",
		},
		{
			name:  "left_spaces",
			input: "   mail@example.com",
		},
		{
			name:  "right_spaces",
			input: "mail@example.com  ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			email := Email(tc.input)

			assert.True(t, email.IsValid())
			require.NotEmpty(t, email)
		})
	}
}

func TestEmail_Invalid(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "empty_email",
			input: "",
		},
		{
			name:  "empty_domain",
			input: "mail@",
		},
		{
			name:  "empty_name",
			input: "@example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			email := Email(tc.input)

			assert.False(t, email.IsValid())
		})
	}
}
