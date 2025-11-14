package repositories

import (
	"citary-backend/internal/domain/entities"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockRoleRepository is a mock implementation of RoleRepository for testing
type MockRoleRepository struct {
	mock.Mock
}

// FindByCode mocks the FindByCode method
func (m *MockRoleRepository) FindByCode(ctx context.Context, code string) (*entities.Role, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Role), args.Error(1)
}
