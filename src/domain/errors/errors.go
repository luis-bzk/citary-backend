package errors

import (
	"fmt"

	"citary-backend/src/shared/constants"
)

// DomainError representa un error del dominio con código de estado HTTP
type DomainError struct {
	Message    string
	StatusCode int
	Err        error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewDomainError crea un nuevo error de dominio
func NewDomainError(message string, statusCode int, err error) *DomainError {
	return &DomainError{
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// Errores predefinidos

func ErrNotFound(message string) *DomainError {
	if message == "" {
		message = constants.ErrorMessages.NotFound
	}
	return NewDomainError(message, constants.StatusCode.NotFound, nil)
}

func ErrBadRequest(message string) *DomainError {
	if message == "" {
		message = constants.ErrorMessages.BadRequest
	}
	return NewDomainError(message, constants.StatusCode.BadRequest, nil)
}

func ErrInternal(err error) *DomainError {
	return NewDomainError(constants.ErrorMessages.InternalError, constants.StatusCode.InternalServerError, err)
}

func ErrUnauthorized(message string) *DomainError {
	if message == "" {
		message = constants.ErrorMessages.Unauthorized
	}
	return NewDomainError(message, constants.StatusCode.Unauthorized, nil)
}

func ErrConflict(message string) *DomainError {
	if message == "" {
		message = constants.ErrorMessages.AlreadyExists
	}
	return NewDomainError(message, constants.StatusCode.Conflict, nil)
}
