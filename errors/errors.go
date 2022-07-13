package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Error struct {
	Message string
	Status  int
}

func (e *Error) Error() string {
	return e.Message
}

func New(message string, status int) *Error {
	return &Error{
		Message: message,
		Status:  status,
	}
}

// InActiveUserError defines an inactive user error
var InActiveUserError = errors.New("user is inactive")

func GetUniqueContraintError(err error) *Error {
	fields := strings.Split(err.Error(), "UNIQUE constraint failed: ")
	return &Error{
		Message: fmt.Sprintf("%s must be unique", strings.Split(fields[1], ".")[1]),
		Status:  http.StatusBadRequest,
	}
}

func GetValidationError(err ValidationError) *Error {
	return &Error{
		Message: err.Error(),
		Status:  http.StatusBadRequest,
	}
}
