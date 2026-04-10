package errs

import "errors"

var (
	// Domain's
	ErrInvalidTitle = errors.New("len of title must be more than 3 and less than 64")
	ErrInvalidBody  = errors.New("len of body must be more thna 3 and less than 512")

	// Repostiry's
	ErrPostDoesntExist = errors.New("post doesn't exist")

	// General's
	ErrInvalidToken    = errors.New("invalid token")
	ErrInternalServer  = errors.New("internal server error")
	ErrInvalidPostID   = errors.New("invalid post's id")
	ErrTooLargeUpdates = errors.New("too large updates")
	ErrNoEnoughRights  = errors.New("no enough rights")
	ErrInvalidSize     = errors.New("size must be less than 100")
)
