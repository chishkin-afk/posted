package object

import (
	"testing"

	"github.com/chishkin-afk/posted/object-service/pkg/errs"
	"github.com/stretchr/testify/assert"
)

func TestNewObject_Invalid(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		body     []byte
		err      error
	}{
		{"invalid_filename_path", "path/to/file.txt", []byte("data"), errs.ErrInvalidFilename},
		{"invalid_filename_dot", ".", []byte("data"), errs.ErrInvalidFilename},
		{"invalid_filename_double_dot", "..", []byte("data"), errs.ErrInvalidFilename},
		{"invalid_filename_special_chars", "file<name>.txt", []byte("data"), errs.ErrInvalidFilename},
		{"invalid_filename_empty", "", []byte("data"), errs.ErrInvalidFilename},
		{"invalid_filename_spaces_only", "   ", []byte("data"), errs.ErrInvalidFilename},

		{"empty_body", "valid_file.txt", []byte{}, errs.ErrInvalidBody},
		{"body_too_large", "valid_file.txt", make([]byte, 4*1024*1024+1), errs.ErrInvalidBody},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename := Filename(tc.filename)
			obj, err := New(filename, tc.body)

			assert.Nil(t, obj)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func TestNewObject_Success(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		body     []byte
	}{
		{"minimal_valid", "file.txt", []byte("d")},
		{"normal_valid", "my_document_v1.pdf", []byte("some content here")},
		{"max_size_valid", "large_file.bin", make([]byte, 4*1024*1024)},
		{"with_spaces_in_name", "my file name.txt", []byte("data")},
		{"with_underscore_and_dash", "test_file-name.log", []byte("log data")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename := Filename(tc.filename)
			obj, err := New(filename, tc.body)

			assert.NoError(t, err, "Не ожидалось ошибок при валидных данных")
			assert.NotNil(t, obj, "Объект должен быть создан")
			assert.Equal(t, tc.filename, obj.Filename().String())
			assert.Equal(t, tc.body, obj.Body())
		})
	}
}
