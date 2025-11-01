package errors

import (
	"fmt"

	"citary-backend/pkg/constants"
)

// DomainError represents a domain-level error with HTTP status code mapping
type DomainError struct {
	Message    string
	StatusCode int
	Err        error
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error for error unwrapping
func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates a new domain error with the specified message, status code, and wrapped error
func NewDomainError(message string, statusCode int, err error) *DomainError {
	return &DomainError{
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// ErrNotFound creates a not found error (404)
func ErrNotFound(message string) *DomainError {
	if message == "" {
		message = constants.ErrorMessages.NotFound
	}
	return NewDomainError(message, constants.StatusCode.NotFound, nil)
}

// ErrBadRequest creates a bad request error (400)
func ErrBadRequest(message string) *DomainError {
	if message == "" {
		message = constants.ErrorMessages.BadRequest
	}
	return NewDomainError(message, constants.StatusCode.BadRequest, nil)
}

// ErrInternal creates an internal server error (500)
func ErrInternal(err error) *DomainError {
	return NewDomainError(constants.ErrorMessages.InternalError, constants.StatusCode.InternalServerError, err)
}

// ErrUnauthorized creates an unauthorized error (401)
func ErrUnauthorized(message string) *DomainError {
	if message == "" {
		message = constants.ErrorMessages.Unauthorized
	}
	return NewDomainError(message, constants.StatusCode.Unauthorized, nil)
}

// ErrConflict creates a conflict error (409)
func ErrConflict(message string) *DomainError {
	if message == "" {
		message = constants.ErrorMessages.AlreadyExists
	}
	return NewDomainError(message, constants.StatusCode.Conflict, nil)
}
