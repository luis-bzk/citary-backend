package usecase_auth

import (
	dtosAuth "citary-backend/src/domain/dtos/auth"
	"citary-backend/src/domain/entities"
	"citary-backend/src/domain/errors"
	portsDrivens "citary-backend/src/ports/drivens"
	"citary-backend/src/shared/constants"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type SignupUserUseCase struct {
	userRepository portsDrivens.UserRepository
}

func NewSignupUserUseCase(userRepository portsDrivens.UserRepository) *SignupUserUseCase {
	return &SignupUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *SignupUserUseCase) Execute(dto dtosAuth.SignupUserDTO) (*entities.User, error) {
	// 1. Validate input data
	if err := dto.Validate(); err != nil {
		return nil, errors.ErrBadRequest(err.Error())
	}

	// 2. Verify if user exists
	existingUser, err := uc.userRepository.FindByEmail(dto.Email)
	if err == nil && existingUser != nil {
		return nil, errors.ErrConflict(constants.ErrorMessages.UserAlreadyExists)
	}

	// 3. Hash password
	hashedPassword, err := hashPassword(dto.Password)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}

	// 4. Create user entity
	user := &entities.User{
		Email:            dto.Email,
		PasswordHash:     hashedPassword,
		EmailVerified:    false,
		PhoneVerified:    false,
		TwoFactorEnabled: false,
		LoginAttempts:    0,
		CreatedDate:      time.Now(),
		RecordStatus:     constants.RecordStatus.Active,
	}

	// 5. Persistir el usuario
	if err := uc.userRepository.Create(user); err != nil {
		return nil, errors.ErrInternal(err)
	}

	return user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
