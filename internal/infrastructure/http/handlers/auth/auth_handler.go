package auth

import (
	authDTO "citary-backend/internal/domain/dtos/auth"
	"citary-backend/internal/domain/usecases/auth"
	httpDTO "citary-backend/internal/infrastructure/http/dto"
	"citary-backend/internal/infrastructure/http/response"
	"citary-backend/pkg/constants"
	"encoding/json"
	"net/http"
)

// AuthHandler handles HTTP requests for authentication operations
type AuthHandler struct {
	signupUserUseCase *auth.SignupUserUseCase
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(signupUserUseCase *auth.SignupUserUseCase) *AuthHandler {
	return &AuthHandler{
		signupUserUseCase: signupUserUseCase,
	}
}

// SignupUser handles user registration requests
func (h *AuthHandler) SignupUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendError(w, constants.StatusCode.BadRequest, "Method not allowed")
		return
	}

	var req authDTO.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendError(w, constants.StatusCode.BadRequest, "Invalid JSON")
		return
	}

	user, err := h.signupUserUseCase.Execute(r.Context(), req)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	authResponse := httpDTO.SignupResponse{
		ID:            user.ID,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		CreatedDate:   user.CreatedDate,
	}

	response.SendSuccess(w, constants.StatusCode.Created, constants.SuccessMessages.UserCreated, authResponse)
}
