package myerrors

import "errors"

var (
	ErrDuplicated     = errors.New("user with this email already exists used")
	ErrEmpty          = errors.New("the field is empty")
	ErrEmailing       = errors.New("error sending email")
	ErrInvalidToken   = errors.New("invalid token")
	ErrInvalidCreds   = errors.New("invalid credentials")
	ErrInternal       = errors.New("internal error")
	ErrDeletingTokens = errors.New("error deleting tokens")
	ErrSession        = errors.New("session error")
)
