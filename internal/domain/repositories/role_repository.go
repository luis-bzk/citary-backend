package repositories

import (
	"citary-backend/internal/domain/entities"
	"context"
)

// RoleRepository defines the contract for role data operations
type RoleRepository interface {
	// FindByCode retrieves a role by its code
	FindByCode(ctx context.Context, code string) (*entities.Role, error)
}
