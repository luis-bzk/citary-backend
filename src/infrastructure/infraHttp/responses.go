package infraHttp

import (
	"citary-backend/src/domain/errors"
	"citary-backend/src/shared/constants"
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func SendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Success: false,
		Message: message,
		Error: &ErrorInfo{
			Code:    statusCode,
			Message: message,
		},
	}

	json.NewEncoder(w).Encode(response)
}

func HandleDomainError(w http.ResponseWriter, err error) {
	if domainErr, ok := err.(*errors.DomainError); ok {
		SendError(w, domainErr.StatusCode, domainErr.Message)
		return
	}

	SendError(w, constants.StatusCode.InternalServerError, constants.ErrorMessages.InternalError)
}
