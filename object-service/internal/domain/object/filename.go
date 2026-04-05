package object

import (
	"regexp"
	"strings"
)

var safeFilenameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-\. ]+$`)

type Filename string

func (f Filename) IsValid() bool {
	file := string(f)
	if strings.Contains(file, "/") || strings.Contains(file, "\\") {
		return false
	}

	if file == "." || file == ".." {
		return false
	}

	if !safeFilenameRegex.MatchString(file) {
		return false
	}

	return true
}

func (f *Filename) Norm() {
	*f = Filename(strings.TrimSpace(string(*f)))
}

func (f Filename) String() string {
	return string(f)
}
