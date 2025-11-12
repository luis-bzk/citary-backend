package auth_test

import (
	"citary-backend/internal/domain/dtos/auth"
	"citary-backend/internal/domain/entities"
	"citary-backend/internal/domain/errors"
	authUseCase "citary-backend/internal/domain/usecases/auth"
	"citary-backend/pkg/constants"
	mockRepo "citary-backend/test/mocks/repositories"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ==========================================
// ERROR CASES (TDD - Red phase first)
// ==========================================

func TestSignupUserUseCase_Execute_InvalidEmail_Empty(t *testing.T) {
	// Arrange
	mockRepository := new(mockRepo.MockUserRepository)
	useCase := authUseCase.NewSignupUserUseCase(mockRepository)
	ctx := context.Background()

	request := auth.SignupRequest{
		Email:    "", // Invalid: empty
		Password: "ValidPass123!",
	}

	// Act
	user, err := useCase.Execute(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "Email cannot be empty")
	mockRepository.AssertNotCalled(t, "FindByEmail")
	mockRepository.AssertNotCalled(t, "Create")
}

func TestSignupUserUseCase_Execute_InvalidEmail_BadFormat(t *testing.T) {
	// Arrange
	mockRepository := new(mockRepo.MockUserRepository)
	useCase := authUseCase.NewSignupUserUseCase(mockRepository)
	ctx := context.Background()

	request := auth.SignupRequest{
		Email:    "not-an-email", // Invalid: bad format
		Password: "ValidPass123!",
	}

	// Act
	user, err := useCase.Execute(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "Email format is invalid")
	mockRepository.AssertNotCalled(t, "FindByEmail")
	mockRepository.AssertNotCalled(t, "Create")
}

func TestSignupUserUseCase_Execute_InvalidPassword_TooShort(t *testing.T) {
	// Arrange
	mockRepository := new(mockRepo.MockUserRepository)
	useCase := authUseCase.NewSignupUserUseCase(mockRepository)
	ctx := context.Background()

	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "Pass1!", // Invalid: too short
	}

	// Act
	user, err := useCase.Execute(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "Password must be at least 8 characters")
	mockRepository.AssertNotCalled(t, "FindByEmail")
	mockRepository.AssertNotCalled(t, "Create")
}

func TestSignupUserUseCase_Execute_InvalidPassword_NoUppercase(t *testing.T) {
	// Arrange
	mockRepository := new(mockRepo.MockUserRepository)
	useCase := authUseCase.NewSignupUserUseCase(mockRepository)
	ctx := context.Background()

	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "validpass123!", // Invalid: no uppercase
	}

	// Act
	user, err := useCase.Execute(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "Password must contain at least one uppercase letter")
	mockRepository.AssertNotCalled(t, "FindByEmail")
	mockRepository.AssertNotCalled(t, "Create")
}

func TestSignupUserUseCase_Execute_InvalidPassword_NoDigit(t *testing.T) {
	// Arrange
	mockRepository := new(mockRepo.MockUserRepository)
	useCase := authUseCase.NewSignupUserUseCase(mockRepository)
	ctx := context.Background()

	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "ValidPass!", // Invalid: no digit
	}

	// Act
	user, err := useCase.Execute(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "Password must contain at least one digit")
	mockRepository.AssertNotCalled(t, "FindByEmail")
	mockRepository.AssertNotCalled(t, "Create")
}

func TestSignupUserUseCase_Execute_InvalidPassword_NoSpecialChar(t *testing.T) {
	// Arrange
	mockRepository := new(mockRepo.MockUserRepository)
	useCase := authUseCase.NewSignupUserUseCase(mockRepository)
	ctx := context.Background()

	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "ValidPass123", // Invalid: no special char
	}

	// Act
	user, err := useCase.Execute(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "Password must contain at least one special character")
	mockRepository.AssertNotCalled(t, "FindByEmail")
	mockRepository.AssertNotCalled(t, "Create")
}

func TestSignupUserUseCase_Execute_UserAlreadyExists(t *testing.T) {
	// Arrange
	mockRepository := new(mockRepo.MockUserRepository)
	useCase := authUseCase.NewSignupUserUseCase(mockRepository)
	ctx := context.Background()

	request := auth.SignupRequest{
		Email:    "existing@example.com",
		Password: "ValidPass123!",
	}

	existingUser := &entities.User{
		ID:    1,
		Email: "existing@example.com",
	}

	// Mock: User already exists
	mockRepository.On("FindByEmail", ctx, "existing@example.com").Return(existingUser, nil)

	// Act
	user, err := useCase.Execute(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), constants.ErrorMessages.UserAlreadyExists)
	mockRepository.AssertExpectations(t)
	mockRepository.AssertNotCalled(t, "Create")
}

func TestSignupUserUseCase_Execute_RepositoryCreateError(t *testing.T) {
	// Arrange
	mockRepository := new(mockRepo.MockUserRepository)
	useCase := authUseCase.NewSignupUserUseCase(mockRepository)
	ctx := context.Background()

	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "ValidPass123!",
	}

	// Mock: User doesn't exist
	mockRepository.On("FindByEmail", ctx, "test@example.com").Return(nil, errors.ErrNotFound("User not found"))

	// Mock: Create fails
	mockRepository.On("Create", ctx, mock.AnythingOfType("*entities.User")).Return(errors.ErrInternal(nil))

	// Act
	user, err := useCase.Execute(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepository.AssertExpectations(t)
}

// ==========================================
// SUCCESS CASES (TDD - Green phase)
// ==========================================

func TestSignupUserUseCase_Execute_Success(t *testing.T) {
	// Arrange
	mockRepository := new(mockRepo.MockUserRepository)
	useCase := authUseCase.NewSignupUserUseCase(mockRepository)
	ctx := context.Background()

	request := auth.SignupRequest{
		Email:    "newuser@example.com",
		Password: "ValidPass123!",
	}

	// Mock: User doesn't exist
	mockRepository.On("FindByEmail", ctx, "newuser@example.com").Return(nil, errors.ErrNotFound("User not found"))

	// Mock: Create succeeds and sets the ID
	mockRepository.On("Create", ctx, mock.AnythingOfType("*entities.User")).
		Run(func(args mock.Arguments) {
			user := args.Get(1).(*entities.User)
			user.ID = 1 // Simulate database setting the ID
		}).
		Return(nil)

	// Act
	user, err := useCase.Execute(ctx, request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "newuser@example.com", user.Email)
	assert.NotEmpty(t, user.PasswordHash, "Password should be hashed")
	assert.NotEqual(t, "ValidPass123!", user.PasswordHash, "Password should not be stored in plain text")
	assert.False(t, user.EmailVerified)
	assert.False(t, user.PhoneVerified)
	assert.False(t, user.TwoFactorEnabled)
	assert.Equal(t, 0, user.LoginAttempts)
	assert.Equal(t, constants.RecordStatus.Active, user.RecordStatus)
	assert.False(t, user.CreatedDate.IsZero(), "CreatedDate should be set")
	assert.True(t, user.CreatedDate.Before(time.Now().Add(1*time.Second)))
	mockRepository.AssertExpectations(t)
}

func TestSignupUserUseCase_Execute_Success_DifferentEmails(t *testing.T) {
	testCases := []struct {
		name  string
		email string
	}{
		{name: "Simple email", email: "user@example.com"},
		{name: "Email with subdomain", email: "user@mail.example.com"},
		{name: "Email with plus", email: "user+test@example.com"},
		{name: "Email with dots", email: "user.name@example.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepository := new(mockRepo.MockUserRepository)
			useCase := authUseCase.NewSignupUserUseCase(mockRepository)
			ctx := context.Background()

			request := auth.SignupRequest{
				Email:    tc.email,
				Password: "ValidPass123!",
			}

			// Mock: User doesn't exist
			mockRepository.On("FindByEmail", ctx, tc.email).Return(nil, errors.ErrNotFound("User not found"))

			// Mock: Create succeeds
			mockRepository.On("Create", ctx, mock.AnythingOfType("*entities.User")).
				Run(func(args mock.Arguments) {
					user := args.Get(1).(*entities.User)
					user.ID = 1
				}).
				Return(nil)

			// Act
			user, err := useCase.Execute(ctx, request)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, tc.email, user.Email)
			mockRepository.AssertExpectations(t)
		})
	}
}

func TestSignupUserUseCase_Execute_PasswordIsHashed(t *testing.T) {
	// Arrange
	mockRepository := new(mockRepo.MockUserRepository)
	useCase := authUseCase.NewSignupUserUseCase(mockRepository)
	ctx := context.Background()

	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "MySecretPassword123!",
	}

	// Mock: User doesn't exist
	mockRepository.On("FindByEmail", ctx, "test@example.com").Return(nil, errors.ErrNotFound("User not found"))

	// Mock: Create succeeds
	var capturedPasswordHash string
	mockRepository.On("Create", ctx, mock.AnythingOfType("*entities.User")).
		Run(func(args mock.Arguments) {
			user := args.Get(1).(*entities.User)
			capturedPasswordHash = user.PasswordHash
			user.ID = 1
		}).
		Return(nil)

	// Act
	user, err := useCase.Execute(ctx, request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEqual(t, "MySecretPassword123!", capturedPasswordHash, "Password should be hashed")
	assert.NotEmpty(t, capturedPasswordHash)
	assert.True(t, len(capturedPasswordHash) > 50, "Bcrypt hash should be at least 50 characters")
	assert.Contains(t, capturedPasswordHash, "$2a$", "Should be a bcrypt hash")
	mockRepository.AssertExpectations(t)
}

func TestNewSignupUserUseCase_ReturnsValidInstance(t *testing.T) {
	// Arrange
	mockRepository := new(mockRepo.MockUserRepository)

	// Act
	useCase := authUseCase.NewSignupUserUseCase(mockRepository)

	// Assert
	assert.NotNil(t, useCase)
}
