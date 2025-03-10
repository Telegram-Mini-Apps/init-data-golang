package initdata

import "errors"

var (
	ErrAuthDateMissing  = errors.New("auth_date is missing")
	ErrAuthDateInvalid  = errors.New("auth_date is invalid")
	ErrSignMissing      = errors.New("sign is missing")
	ErrSignInvalid      = errors.New("sign is invalid")
	ErrUnexpectedFormat = errors.New("init data has unexpected format")
	ErrExpired          = errors.New("init data is expired")
)
