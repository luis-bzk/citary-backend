package repositories

import (
	"citary-backend/internal/domain/entities"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

// FindByEmail mocks the FindByEmail method
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entities.User), args.Error(1)
}

// Create mocks the Create method
func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
