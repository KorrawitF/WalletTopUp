package errs

import (
	"net/http"
)

type Error struct {
	Code    int
	Status  string
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

var _ error = (*Error)(nil)

var (
	ErrInvalidMethod = &Error{Code: http.StatusBadRequest, Status: "fail", Message: "Invalid payment method."}
	ErrUserNotFound  = &Error{Code: http.StatusBadRequest, Status: "fail", Message: "User not found."}
)
