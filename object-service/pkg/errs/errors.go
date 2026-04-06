package errs

import "errors"

var (
	// Domain's
	ErrInvalidExtension = errors.New("invalid extension of file")
	ErrInvalidBody      = errors.New("body of file must be more than 0 and less than 4 megabytes")

	// Repository's
	ErrObjectNotFound = errors.New("object not found")
)
