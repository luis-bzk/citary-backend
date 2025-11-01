package auth

import (
	"citary-backend/internal/domain/dtos/auth"
	"citary-backend/internal/domain/entities"
	"citary-backend/internal/domain/errors"
	"citary-backend/internal/domain/repositories"
	"citary-backend/pkg/constants"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SignupUserUseCase handles the business logic for user registration
type SignupUserUseCase struct {
	userRepository repositories.UserRepository
}

// NewSignupUserUseCase creates a new instance of SignupUserUseCase
func NewSignupUserUseCase(userRepository repositories.UserRepository) *SignupUserUseCase {
	return &SignupUserUseCase{
		userRepository: userRepository,
	}
}

// Execute processes a user signup request
func (uc *SignupUserUseCase) Execute(ctx context.Context, dto auth.SignupRequest) (*entities.User, error) {
	// 1. Validate input data
	if err := dto.Validate(); err != nil {
		return nil, errors.ErrBadRequest(err.Error())
	}

	// 2. Verify if user exists
	existingUser, err := uc.userRepository.FindByEmail(ctx, dto.Email)
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

	// 5. Persist the user
	if err := uc.userRepository.Create(ctx, user); err != nil {
		return nil, errors.ErrInternal(err)
	}

	return user, nil
}

// hashPassword generates a bcrypt hash of the given password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
