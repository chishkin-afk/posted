package errs

import "errors"

var (
	// Domain's
	ErrInvalidPassword = errors.New("len of password must be more than 6 and less than 64")
	ErrInvalidEmail    = errors.New("user's email is invalid")
	ErrInvalidNickname = errors.New("len of nick must be more than 3 and less than 64")
)
