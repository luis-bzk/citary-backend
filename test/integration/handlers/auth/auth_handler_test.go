package auth_test

import (
	"bytes"
	"citary-backend/internal/domain/entities"
	authUseCase "citary-backend/internal/domain/usecases/auth"
	authHandler "citary-backend/internal/infrastructure/http/handlers/auth"
	"citary-backend/pkg/constants"
	mockRepo "citary-backend/test/mocks/repositories"
	mockService "citary-backend/test/mocks/services"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ==========================================
// ERROR CASES (TDD - Red phase first)
// ==========================================

func TestAuthHandler_SignupUser_MethodNotAllowed(t *testing.T) {
	// Arrange
	mockUserRepository := new(mockRepo.MockUserRepository)
	mockRoleRepository := new(mockRepo.MockRoleRepository)
	mockEmailService := new(mockService.MockEmailService)
	useCase := authUseCase.NewSignupUserUseCase(mockUserRepository, mockRoleRepository, mockEmailService)
	handler := authHandler.NewAuthHandler(useCase)

	req := httptest.NewRequest(http.MethodGet, "/auth/signup", nil) // Wrong method
	w := httptest.NewRecorder()

	// Act
	handler.SignupUser(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Method not allowed")
}

func TestAuthHandler_SignupUser_InvalidJSON(t *testing.T) {
	// Arrange
	mockUserRepository := new(mockRepo.MockUserRepository)
	mockRoleRepository := new(mockRepo.MockRoleRepository)
	mockEmailService := new(mockService.MockEmailService)
	useCase := authUseCase.NewSignupUserUseCase(mockUserRepository, mockRoleRepository, mockEmailService)
	handler := authHandler.NewAuthHandler(useCase)

	invalidJSON := bytes.NewBufferString(`{"email": "test@example.com", "password": `)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", invalidJSON)
	w := httptest.NewRecorder()

	// Act
	handler.SignupUser(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid JSON")
}

func TestAuthHandler_SignupUser_EmptyBody(t *testing.T) {
	// Arrange
	mockUserRepository := new(mockRepo.MockUserRepository)
	mockRoleRepository := new(mockRepo.MockRoleRepository)
	mockEmailService := new(mockService.MockEmailService)
	useCase := authUseCase.NewSignupUserUseCase(mockUserRepository, mockRoleRepository, mockEmailService)
	handler := authHandler.NewAuthHandler(useCase)

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBufferString(""))
	w := httptest.NewRecorder()

	// Act
	handler.SignupUser(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_SignupUser_InvalidEmail_Empty(t *testing.T) {
	// Arrange
	mockUserRepository := new(mockRepo.MockUserRepository)
	mockRoleRepository := new(mockRepo.MockRoleRepository)
	mockEmailService := new(mockService.MockEmailService)
	useCase := authUseCase.NewSignupUserUseCase(mockUserRepository, mockRoleRepository, mockEmailService)
	handler := authHandler.NewAuthHandler(useCase)

	payload := map[string]string{
		"email":    "", // Empty email
		"password": "ValidPass123!",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Act
	handler.SignupUser(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Email cannot be empty")
}

func TestAuthHandler_SignupUser_InvalidEmail_BadFormat(t *testing.T) {
	// Arrange
	mockUserRepository := new(mockRepo.MockUserRepository)
	mockRoleRepository := new(mockRepo.MockRoleRepository)
	mockEmailService := new(mockService.MockEmailService)
	useCase := authUseCase.NewSignupUserUseCase(mockUserRepository, mockRoleRepository, mockEmailService)
	handler := authHandler.NewAuthHandler(useCase)

	payload := map[string]string{
		"email":    "not-an-email", // Bad format
		"password": "ValidPass123!",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Act
	handler.SignupUser(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Email format is invalid")
}

func TestAuthHandler_SignupUser_InvalidPassword_TooShort(t *testing.T) {
	// Arrange
	mockUserRepository := new(mockRepo.MockUserRepository)
	mockRoleRepository := new(mockRepo.MockRoleRepository)
	mockEmailService := new(mockService.MockEmailService)
	useCase := authUseCase.NewSignupUserUseCase(mockUserRepository, mockRoleRepository, mockEmailService)
	handler := authHandler.NewAuthHandler(useCase)

	payload := map[string]string{
		"email":    "test@example.com",
		"password": "Pass1!", // Too short
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Act
	handler.SignupUser(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Password must be at least 8 characters")
}

func TestAuthHandler_SignupUser_UserAlreadyExists(t *testing.T) {
	// Arrange
	mockUserRepository := new(mockRepo.MockUserRepository)
	mockRoleRepository := new(mockRepo.MockRoleRepository)
	mockEmailService := new(mockService.MockEmailService)
	useCase := authUseCase.NewSignupUserUseCase(mockUserRepository, mockRoleRepository, mockEmailService)
	handler := authHandler.NewAuthHandler(useCase)

	existingUser := &entities.User{
		ID:           1,
		Email:        "existing@example.com",
		RecordStatus: constants.RecordStatus.Active,
	}

	// Mock: User already exists
	mockUserRepository.On("FindByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)

	payload := map[string]string{
		"email":    "existing@example.com",
		"password": "ValidPass123!",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Act
	handler.SignupUser(w, req)

	// Assert
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), constants.ErrorMessages.UserAlreadyExists)
	mockUserRepository.AssertExpectations(t)
}

// ==========================================
// SUCCESS CASES (TDD - Green phase)
// ==========================================

func TestAuthHandler_SignupUser_Success(t *testing.T) {
	// Arrange
	mockUserRepository := new(mockRepo.MockUserRepository)
	mockRoleRepository := new(mockRepo.MockRoleRepository)
	mockEmailService := new(mockService.MockEmailService)
	useCase := authUseCase.NewSignupUserUseCase(mockUserRepository, mockRoleRepository, mockEmailService)
	handler := authHandler.NewAuthHandler(useCase)

	// Mock: User doesn't exist (returns nil, nil per architecture)
	mockUserRepository.On("FindByEmail", mock.Anything, "newuser@example.com").
		Return(nil, nil)

	// Mock: Default role exists
	defaultRole := &entities.Role{
		ID:           1,
		Code:         constants.DefaultUserRole,
		RecordStatus: constants.RecordStatus.Active,
	}
	mockRoleRepository.On("FindByCode", mock.Anything, constants.DefaultUserRole).Return(defaultRole, nil)

	// Mock: Create succeeds
	mockUserRepository.On("Create", mock.Anything, mock.AnythingOfType("*entities.User")).
		Run(func(args mock.Arguments) {
			user := args.Get(1).(*entities.User)
			user.ID = 1 // Simulate database setting the ID
			user.RoleID = 1
		}).
		Return(nil)

	// Mock: Email service succeeds
	mockEmailService.On("SendVerificationEmail", mock.Anything, "newuser@example.com", mock.AnythingOfType("string")).Return(nil)

	payload := map[string]string{
		"email":    "newuser@example.com",
		"password": "ValidPass123!",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Act
	handler.SignupUser(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.Equal(t, constants.SuccessMessages.UserCreated, response["message"])
	assert.NotNil(t, response["data"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["id"].(float64))
	assert.Equal(t, "newuser@example.com", data["email"])
	assert.False(t, data["emailVerified"].(bool))
	assert.NotEmpty(t, data["createdDate"])

	mockUserRepository.AssertExpectations(t)
}

func TestAuthHandler_SignupUser_Success_DifferentEmails(t *testing.T) {
	testCases := []struct {
		name  string
		email string
	}{
		{name: "Simple email", email: "user@example.com"},
		{name: "Email with subdomain", email: "user@mail.example.com"},
		{name: "Email with plus", email: "user+test@example.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockUserRepository := new(mockRepo.MockUserRepository)
			mockRoleRepository := new(mockRepo.MockRoleRepository)
			mockEmailService := new(mockService.MockEmailService)
			useCase := authUseCase.NewSignupUserUseCase(mockUserRepository, mockRoleRepository, mockEmailService)
			handler := authHandler.NewAuthHandler(useCase)

			// Mock: User doesn't exist (returns nil, nil per architecture)
			mockUserRepository.On("FindByEmail", mock.Anything, tc.email).
				Return(nil, nil)

			// Mock: Default role exists
			defaultRole := &entities.Role{
				ID:           1,
				Code:         constants.DefaultUserRole,
				RecordStatus: constants.RecordStatus.Active,
			}
			mockRoleRepository.On("FindByCode", mock.Anything, constants.DefaultUserRole).Return(defaultRole, nil)

			// Mock: Create succeeds
			mockUserRepository.On("Create", mock.Anything, mock.AnythingOfType("*entities.User")).
				Run(func(args mock.Arguments) {
					user := args.Get(1).(*entities.User)
					user.ID = 1
					user.RoleID = 1
				}).
				Return(nil)

			// Mock: Email service succeeds
			mockEmailService.On("SendVerificationEmail", mock.Anything, tc.email, mock.AnythingOfType("string")).Return(nil)

			payload := map[string]string{
				"email":    tc.email,
				"password": "ValidPass123!",
			}
			body, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			// Act
			handler.SignupUser(w, req)

			// Assert
			assert.Equal(t, http.StatusCreated, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.True(t, response["success"].(bool))

			data := response["data"].(map[string]interface{})
			assert.Equal(t, tc.email, data["email"])

			mockUserRepository.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_SignupUser_ResponseFormat(t *testing.T) {
	// Arrange
	mockUserRepository := new(mockRepo.MockUserRepository)
	mockRoleRepository := new(mockRepo.MockRoleRepository)
	mockEmailService := new(mockService.MockEmailService)
	useCase := authUseCase.NewSignupUserUseCase(mockUserRepository, mockRoleRepository, mockEmailService)
	handler := authHandler.NewAuthHandler(useCase)

	ctx := context.Background()

	// Mock: User doesn't exist (returns nil, nil per architecture)
	mockUserRepository.On("FindByEmail", ctx, "format@example.com").
		Return(nil, nil)

	// Mock: Default role exists
	defaultRole := &entities.Role{
		ID:           1,
		Code:         constants.DefaultUserRole,
		RecordStatus: constants.RecordStatus.Active,
	}
	mockRoleRepository.On("FindByCode", ctx, constants.DefaultUserRole).Return(defaultRole, nil)

	// Mock: Create succeeds
	mockUserRepository.On("Create", ctx, mock.AnythingOfType("*entities.User")).
		Run(func(args mock.Arguments) {
			user := args.Get(1).(*entities.User)
			user.ID = 99
			user.RoleID = 1
		}).
		Return(nil)

	// Mock: Email service succeeds
	mockEmailService.On("SendVerificationEmail", ctx, "format@example.com", mock.AnythingOfType("string")).Return(nil)

	payload := map[string]string{
		"email":    "format@example.com",
		"password": "ValidPass123!",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(body))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Act
	handler.SignupUser(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response structure
	assert.Contains(t, response, "success")
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "data")

	// Verify data fields
	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "id")
	assert.Contains(t, data, "email")
	assert.Contains(t, data, "emailVerified")
	assert.Contains(t, data, "createdDate")
	assert.NotContains(t, data, "password", "Password should not be in response")
	assert.NotContains(t, data, "passwordHash", "Password hash should not be in response")

	mockUserRepository.AssertExpectations(t)
}

func TestNewAuthHandler_ReturnsValidInstance(t *testing.T) {
	// Arrange
	mockUserRepository := new(mockRepo.MockUserRepository)
	mockRoleRepository := new(mockRepo.MockRoleRepository)
	mockEmailService := new(mockService.MockEmailService)
	useCase := authUseCase.NewSignupUserUseCase(mockUserRepository, mockRoleRepository, mockEmailService)

	// Act
	handler := authHandler.NewAuthHandler(useCase)

	// Assert
	assert.NotNil(t, handler)
}
