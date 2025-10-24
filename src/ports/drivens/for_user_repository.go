package portsDrivens

import "citary-backend/src/domain/entities"

type UserRepository interface {
	FindByEmail(email string) (*entities.User, error)
	Create(user *entities.User) error
}
