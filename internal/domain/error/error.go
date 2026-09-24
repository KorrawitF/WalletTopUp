package errs

import "errors"

var (
	ErrInvalidMethod = errors.New("Invalid payment method.")
	ErrUserNotFound  = errors.New("User not found.")
)
