package response

import (
	"citary-backend/internal/domain/errors"
	"citary-backend/internal/infrastructure/http/dto"
	"citary-backend/pkg/constants"
	"encoding/json"
	"net/http"
)

// SendSuccess writes a successful JSON response
func SendSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := dto.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// SendError writes an error JSON response
func SendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := dto.APIResponse{
		Success: false,
		Message: message,
		Error: &dto.ErrorInfo{
			Code:    statusCode,
			Message: message,
		},
	}

	json.NewEncoder(w).Encode(response)
}

// HandleDomainError handles domain-specific errors and sends appropriate HTTP responses
func HandleDomainError(w http.ResponseWriter, err error) {
	if domainErr, ok := err.(*errors.DomainError); ok {
		SendError(w, domainErr.StatusCode, domainErr.Message)
		return
	}

	SendError(w, constants.StatusCode.InternalServerError, constants.ErrorMessages.InternalError)
}
