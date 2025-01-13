package myerrors

import "errors"

var (
	ErrDuplicated     = errors.New("user with this email already exists used")
	ErrEmpty          = errors.New("the field is empty")
	ErrEmailing       = errors.New("error sending email")
	ErrInvalidToken   = errors.New("invalid token")
	ErrExpiredToken   = errors.New("expired token")
	ErrInvalidCreds   = errors.New("invalid credentials")
	ErrInternal       = errors.New("internal error")
	ErrDeletingTokens = errors.New("error deleting tokens")
	ErrLocked         = errors.New("locked")
	ErrCode           = errors.New("invalid code")
	ErrExpiredCode    = errors.New("expired code")
	ErrNotFound       = errors.New("not found")
	ErrDualSession    = errors.New("you've already been logged in with your device. try to login again")
)
