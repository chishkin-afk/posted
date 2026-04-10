package errs

import "errors"

var (
	// Domain's
	ErrInvalidTitle = errors.New("len of title must be more than 3 and less than 64")
	ErrInvalidBody  = errors.New("len of body must be more thna 3 and less than 512")

	// Repostiry's
	ErrPostDoesntExist = errors.New("post doesn't exist")
)
