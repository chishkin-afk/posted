package errs

import "errors"

var (
	ErrBadGateway = errors.New("bad gateway")

	ErrInvalidEnvironment = errors.New("invalid environment")
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidToken       = errors.New("invalid token")
)
