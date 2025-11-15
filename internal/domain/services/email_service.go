package services

import "context"

// EmailService defines the interface for sending emails
type EmailService interface {
	// SendVerificationEmail sends an email verification link to the user
	SendVerificationEmail(ctx context.Context, email, token string) error
}
