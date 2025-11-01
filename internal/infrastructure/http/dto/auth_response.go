package dto

import "time"

// SignupResponse represents the API response for user signup
type SignupResponse struct {
	ID            int       `json:"id"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"emailVerified"`
	CreatedDate   time.Time `json:"createdDate"`
}
