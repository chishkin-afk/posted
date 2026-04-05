package errs

import "errors"

var (
	// Domain's
	ErrInvalidFilename = errors.New("len of filename must be more than 3 and less than 64")
	ErrInvalidBody     = errors.New("body of file must be more than 0 and less than 4 megabytes")
)
