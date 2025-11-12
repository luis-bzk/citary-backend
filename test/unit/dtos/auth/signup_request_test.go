package auth_test

import (
	"citary-backend/internal/domain/dtos/auth"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==========================================
// EMAIL VALIDATION TESTS (Error cases first)
// ==========================================

func TestSignupRequest_Validate_EmailEmpty(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "",
		Password: "ValidPass123!",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, auth.ErrEmailEmpty, err)
}

func TestSignupRequest_Validate_EmailInvalidFormat_MissingAt(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "invalidemailexample.com",
		Password: "ValidPass123!",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, auth.ErrEmailInvalidFormat, err)
}

func TestSignupRequest_Validate_EmailInvalidFormat_MissingDomain(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "invalid@",
		Password: "ValidPass123!",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, auth.ErrEmailInvalidFormat, err)
}

func TestSignupRequest_Validate_EmailInvalidFormat_MissingLocalPart(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "@example.com",
		Password: "ValidPass123!",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, auth.ErrEmailInvalidFormat, err)
}

func TestSignupRequest_Validate_EmailInvalidFormat_MissingTLD(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "invalid@example",
		Password: "ValidPass123!",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, auth.ErrEmailInvalidFormat, err)
}

func TestSignupRequest_Validate_EmailInvalidFormat_WithSpaces(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "invalid email@example.com",
		Password: "ValidPass123!",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, auth.ErrEmailInvalidFormat, err)
}

func TestSignupRequest_Validate_EmailInvalidFormat_DoubleAt(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "invalid@@example.com",
		Password: "ValidPass123!",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, auth.ErrEmailInvalidFormat, err)
}

func TestSignupRequest_Validate_EmailTooLong(t *testing.T) {
	// Arrange - Create email with 101+ characters
	longEmail := ""
	for i := 0; i < 101; i++ {
		longEmail += "a"
	}
	longEmail += "@test.com"

	request := auth.SignupRequest{
		Email:    longEmail,
		Password: "ValidPass123!",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, auth.ErrEmailTooLong, err)
}

// ==========================================
// PASSWORD VALIDATION TESTS (Error cases first)
// ==========================================

func TestSignupRequest_Validate_PasswordEmpty(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, auth.ErrPasswordEmpty, err)
}

func TestSignupRequest_Validate_PasswordTooShort(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "Pass1!", // Only 6 characters
	}

	// Act
	err := request.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, auth.ErrPasswordTooShort, err)
}

func TestSignupRequest_Validate_PasswordTooLong(t *testing.T) {
	// Arrange - Create password with 101+ characters
	longPassword := ""
	for i := 0; i < 101; i++ {
		longPassword += "a"
	}

	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: longPassword,
	}

	// Act
	err := request.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, auth.ErrPasswordTooLong, err)
}

func TestSignupRequest_Validate_PasswordNoLowercase(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "VALIDPASS123!", // No lowercase
	}

	// Act
	err := request.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, auth.ErrPasswordNoLowercase, err)
}

func TestSignupRequest_Validate_PasswordNoUppercase(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "validpass123!", // No uppercase
	}

	// Act
	err := request.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, auth.ErrPasswordNoUppercase, err)
}

func TestSignupRequest_Validate_PasswordNoDigit(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "ValidPass!", // No digit
	}

	// Act
	err := request.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, auth.ErrPasswordNoDigit, err)
}

func TestSignupRequest_Validate_PasswordNoSpecialChar(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "ValidPass123", // No special character
	}

	// Act
	err := request.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, auth.ErrPasswordNoSpecialChar, err)
}

// ==========================================
// SUCCESS CASES (Happy path - Green phase)
// ==========================================

func TestSignupRequest_Validate_Success_SimpleEmail(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "user@example.com",
		Password: "ValidPass123!",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestSignupRequest_Validate_Success_EmailWithDots(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "user.name@example.com",
		Password: "ValidPass123!",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestSignupRequest_Validate_Success_EmailWithPlus(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "user+tag@example.com",
		Password: "ValidPass123!",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestSignupRequest_Validate_Success_EmailWithNumbers(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "user123@example456.com",
		Password: "ValidPass123!",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestSignupRequest_Validate_Success_EmailWithSubdomain(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "user@mail.example.com",
		Password: "ValidPass123!",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestSignupRequest_Validate_Success_PasswordWithExclamation(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "ValidPass123!",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestSignupRequest_Validate_Success_PasswordWithAtSign(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "ValidPass123@",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestSignupRequest_Validate_Success_PasswordWithHash(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "ValidPass123#",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestSignupRequest_Validate_Success_PasswordWithDollar(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "ValidPass123$",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestSignupRequest_Validate_Success_PasswordWithPercent(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "ValidPass123%",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestSignupRequest_Validate_Success_PasswordWithAsterisk(t *testing.T) {
	// Arrange
	request := auth.SignupRequest{
		Email:    "test@example.com",
		Password: "ValidPass123*",
	}

	// Act
	err := request.Validate()

	// Assert
	assert.NoError(t, err)
}

// ==========================================
// VALIDATION ERROR HELPER TEST
// ==========================================

func TestValidationError_Error(t *testing.T) {
	// Arrange
	expectedMessage := "Test error message"
	validationErr := &auth.ValidationError{Message: expectedMessage}

	// Act
	errorMessage := validationErr.Error()

	// Assert
	assert.Equal(t, expectedMessage, errorMessage)
}
