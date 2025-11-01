package auth

import (
	"regexp"
	"unicode"
)

// SignupRequest represents the data required to create a new user account
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate performs validation on the signup request data
func (dto *SignupRequest) Validate() error {
	// Email validations
	if dto.Email == "" {
		return ErrEmailEmpty
	}

	if len(dto.Email) > 100 {
		return ErrEmailTooLong
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(dto.Email) {
		return ErrEmailInvalidFormat
	}

	// Password validations
	if dto.Password == "" {
		return ErrPasswordEmpty
	}

	if len(dto.Password) < 8 {
		return ErrPasswordTooShort
	}

	if len(dto.Password) > 100 {
		return ErrPasswordTooLong
	}

	var hasLower, hasUpper, hasDigit, hasSpecial bool
	for _, char := range dto.Password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasLower {
		return ErrPasswordNoLowercase
	}

	if !hasUpper {
		return ErrPasswordNoUppercase
	}

	if !hasDigit {
		return ErrPasswordNoDigit
	}

	if !hasSpecial {
		return ErrPasswordNoSpecialChar
	}

	return nil
}

// ValidationError represents a validation error with a custom message
type ValidationError struct {
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return e.Message
}

// Validation error definitions
var (
	// Email validation errors
	ErrEmailEmpty         = &ValidationError{Message: "Email cannot be empty"}
	ErrEmailInvalidFormat = &ValidationError{Message: "Email format is invalid"}
	ErrEmailTooLong       = &ValidationError{Message: "Email cannot exceed 100 characters"}

	// Password validation errors
	ErrPasswordEmpty         = &ValidationError{Message: "Password cannot be empty"}
	ErrPasswordTooShort      = &ValidationError{Message: "Password must be at least 8 characters long"}
	ErrPasswordTooLong       = &ValidationError{Message: "Password cannot exceed 100 characters"}
	ErrPasswordNoLowercase   = &ValidationError{Message: "Password must contain at least one lowercase letter"}
	ErrPasswordNoUppercase   = &ValidationError{Message: "Password must contain at least one uppercase letter"}
	ErrPasswordNoDigit       = &ValidationError{Message: "Password must contain at least one digit"}
	ErrPasswordNoSpecialChar = &ValidationError{Message: "Password must contain at least one special character"}
)
