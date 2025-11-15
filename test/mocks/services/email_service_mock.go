package services

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockEmailService is a mock implementation of the EmailService interface
type MockEmailService struct {
	mock.Mock
}

// SendVerificationEmail mocks the SendVerificationEmail method
func (m *MockEmailService) SendVerificationEmail(ctx context.Context, email, token string) error {
	args := m.Called(ctx, email, token)
	return args.Error(0)
}
