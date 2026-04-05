package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilename_Success(t *testing.T) {
	name := "kotyata.png"
	filename := Filename(name)

	filename.Norm()
	assert.True(t, filename.IsValid())
	assert.Equal(t, name, filename.String())
}

func TestFilename_Invalid(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"contains_slash", "path/to/file.txt"},
		{"contains_backslash", "path\\to\\file.txt"},
		{"is_dot", "."},
		{"is_double_dot", ".."},
		{"contains_special_chars", "file<name>.txt"},
		{"contains_question_mark", "file?.txt"},
		{"contains_asterisk", "file*.txt"},
		{"contains_colon", "file:name.txt"},
		{"contains_pipe", "file|name.txt"},
		{"empty_string", ""},
		{"only_spaces", "   "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename := Filename(tc.input)

			filename.Norm()
			assert.False(t, filename.IsValid())
		})
	}
}
