package dtosAuth

type SignupUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto *SignupUserDTO) Validate() error {
	if dto.Email == "" {
		return ErrInvalidEmail
	}

	if dto.Password == "" || len(dto.Password) < 8 {
		return ErrInvalidPassword
	}

	return nil
}

var (
	ErrInvalidEmail    = &ValidationError{Message: "El email es inválido o está vacío"}
	ErrInvalidPassword = &ValidationError{Message: "La contraseña debe tener al menos 8 caracteres"}
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
