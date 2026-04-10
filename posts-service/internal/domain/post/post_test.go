package post

import (
	"testing"
	"time"

	"github.com/chishkin-afk/posted/posts-service/pkg/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Success(t *testing.T) {
	ownerID := uuid.New()
	title := "title"
	body := "body"

	now := time.Now().UTC()
	post, err := New(ownerID, title, body)

	require.NoError(t, err)
	assert.NotEmpty(t, post.ID())
	assert.Equal(t, ownerID, post.OwnerID())
	assert.Equal(t, title, post.Title())
	assert.Equal(t, body, post.Body())
	assert.WithinDuration(t, now, post.PostedAt(), 100*time.Millisecond)
	assert.WithinDuration(t, now, post.UpdatedAt(), 100*time.Millisecond)
}

func TestNew_Invalid(t *testing.T) {
	type data struct {
		ownerID uuid.UUID
		title   string
		body    string
	}

	testCases := []struct {
		name     string
		input    data
		expected error
	}{
		{
			name: "empty_title",
			input: data{
				ownerID: uuid.New(),
				title:   "",
				body:    "body",
			},
			expected: errs.ErrInvalidTitle,
		},
		{
			name: "empty_body",
			input: data{
				ownerID: uuid.New(),
				title:   "title",
				body:    "",
			},
			expected: errs.ErrInvalidBody,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := New(tc.input.ownerID, tc.input.title, tc.input.body)

			require.Error(t, err)
			assert.EqualError(t, err, tc.expected.Error())
		})
	}
}
