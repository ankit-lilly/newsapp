package errors

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

type InvalidParamsError struct {
	Param   string
	Message string
}

func (e *InvalidParamsError) Error() string {
	return fmt.Sprintf("invalid parameter %q: %s", e.Param, e.Message)
}

type ApiError struct {
	Message string
	Code    int
}

func (e *ApiError) Error() string {
	return e.Message
}

func NewApiError(message string, code int) error {
	return &ApiError{Message: message, Code: code}
}

func IsApiError(err error) bool {
	_, ok := err.(*ApiError)
	return ok
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) error {
	return &NotFoundError{Message: message}
}

func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}
