package object

import (
	"strings"
	"testing"

	"github.com/chishkin-afk/posted/object-service/pkg/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewObject_Invalid(t *testing.T) {
	testCases := []struct {
		name string
		ext  string
		body []byte
		err  error
	}{
		{"ext_empty", "", []byte("data"), errs.ErrInvalidExtension},
		{"ext_no_dot", "txt", []byte("data"), errs.ErrInvalidExtension},
		{"ext_with_slash", ".txt/exe", []byte("data"), errs.ErrInvalidExtension},
		{"ext_with_special_char", ".txt<", []byte("data"), errs.ErrInvalidExtension},
		{"ext_with_colon", ".jp:g", []byte("data"), errs.ErrInvalidExtension},

		{"empty_body", ".txt", []byte{}, errs.ErrInvalidBody},
		{"body_too_large", ".txt", make([]byte, 10*1024*1024+1), errs.ErrInvalidBody},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id := uuid.New()
			obj, err := New(id, tc.ext, tc.body)

			assert.Nil(t, obj)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func TestNewObject_Success(t *testing.T) {
	testCases := []struct {
		name string
		ext  string
		body []byte
	}{
		{"valid_simple", ".txt", []byte("d")},
		{"valid_complex", ".tar.gz", []byte("some content here")},
		{"valid_uppercase", ".JPEG", []byte("image data")},
		{"valid_with_numbers", ".mp4", []byte("video data")},
		{"max_size_valid", ".bin", make([]byte, 10*1024*1024)},
		{"valid_with_dash", ".my-ext", []byte("data")},
		{"valid_with_underscore", ".my_ext", []byte("data")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id := uuid.New()
			obj, err := New(id, tc.ext, tc.body)

			assert.NoError(t, err)
			assert.NotNil(t, obj, "Объект должен быть создан")
			assert.Equal(t, id, obj.ID())
			assert.Equal(t, strings.TrimSpace(tc.ext), obj.Ext())
			assert.Equal(t, tc.body, obj.Body())
		})
	}
}
