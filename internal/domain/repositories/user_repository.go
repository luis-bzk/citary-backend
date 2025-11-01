package repositories

import (
	"citary-backend/internal/domain/entities"
	"context"
)

// UserRepository defines the contract for user data operations
type UserRepository interface {
	// FindByEmail retrieves a user by their email address
	FindByEmail(ctx context.Context, email string) (*entities.User, error)

	// Create persists a new user to the database
	Create(ctx context.Context, user *entities.User) error
}
