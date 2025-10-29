package dtosAuth

import (
	"regexp"
	"unicode"
)

type SignupUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto *SignupUserDTO) Validate() error {
	// Validaciones de email
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

	// Validaciones de password
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

var (
	// Errores de email
	ErrEmailEmpty         = &ValidationError{Message: "El email no puede estar vacío"}
	ErrEmailInvalidFormat = &ValidationError{Message: "El email no tiene un formato válido"}
	ErrEmailTooLong       = &ValidationError{Message: "El email no puede tener más de 100 caracteres"}

	// Errores de password
	ErrPasswordEmpty          = &ValidationError{Message: "La contraseña no puede estar vacía"}
	ErrPasswordTooShort       = &ValidationError{Message: "La contraseña debe tener al menos 8 caracteres"}
	ErrPasswordTooLong        = &ValidationError{Message: "La contraseña no puede tener más de 100 caracteres"}
	ErrPasswordNoLowercase    = &ValidationError{Message: "La contraseña debe contener al menos una letra minúscula"}
	ErrPasswordNoUppercase    = &ValidationError{Message: "La contraseña debe contener al menos una letra mayúscula"}
	ErrPasswordNoDigit        = &ValidationError{Message: "La contraseña debe contener al menos un número"}
	ErrPasswordNoSpecialChar  = &ValidationError{Message: "La contraseña debe contener al menos un carácter especial"}
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
