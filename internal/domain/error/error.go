package errs

import (
	"errors"
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
	ErrInvalidMethod  = &Error{Code: http.StatusBadRequest, Status: "fail", Message: "Invalid payment method."}
	ErrUserNotFound   = &Error{Code: http.StatusBadRequest, Status: "fail", Message: "User not found."}
	ErrTxNotFound     = &Error{Code: http.StatusBadRequest, Status: "fail", Message: "Transaction not found."}
	ErrTxCompleted    = &Error{Code: http.StatusBadRequest, Status: "fail", Message: "Transaction has been completed."}
	ErrTxExpired      = &Error{Code: http.StatusBadRequest, Status: "fail", Message: "Transaction has been expired."}
	ErrTxNotVerified  = &Error{Code: http.StatusBadRequest, Status: "fail", Message: "Transaction not verified."}
	ErrWalletNotFound = &Error{Code: http.StatusBadRequest, Status: "fail", Message: "Wallet not found."}

	ErrCacheMiss = errors.New("cache miss")
)
