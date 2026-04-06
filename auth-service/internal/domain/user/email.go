package user

import (
	"net/mail"
	"strings"
)

type Email string

func (e Email) String() string {
	return string(e)
}

func (e Email) IsValid() bool {
	e = Email(strings.TrimSpace(e.String()))
	if e == "" {
		return false
	}

	if _, err := mail.ParseAddress(e.String()); err != nil {
		return false
	}

	return true
}
