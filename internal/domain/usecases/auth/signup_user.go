package auth

import (
	"citary-backend/internal/domain/dtos/auth"
	"citary-backend/internal/domain/entities"
	"citary-backend/internal/domain/errors"
	"citary-backend/internal/domain/repositories"
	"citary-backend/pkg/constants"
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SignupUserUseCase handles the business logic for user registration
type SignupUserUseCase struct {
	userRepository repositories.UserRepository
	roleRepository repositories.RoleRepository
}

// NewSignupUserUseCase creates a new instance of SignupUserUseCase
func NewSignupUserUseCase(
	userRepository repositories.UserRepository,
	roleRepository repositories.RoleRepository,
) *SignupUserUseCase {
	return &SignupUserUseCase{
		userRepository: userRepository,
		roleRepository: roleRepository,
	}
}

// Execute processes a user signup request
func (uc *SignupUserUseCase) Execute(ctx context.Context, dto auth.SignupRequest) (*entities.User, error) {
	log.Printf("[SignupUserUseCase] Execute: email=%s", dto.Email)

	// 1. Validate input data
	if err := dto.Validate(); err != nil {
		log.Printf("[SignupUserUseCase] Validation failed: %v", err)
		return nil, errors.ErrBadRequest(err.Error())
	}

	// 2. Verify if user exists (BUSINESS LOGIC - validate both physical and logical existence)
	existingUser, err := uc.userRepository.FindByEmail(ctx, dto.Email)
	if err != nil {
		// Technical error from repository
		log.Printf("[SignupUserUseCase] Error checking existing user: %v", err)
		return nil, err
	}

	// Business validation: Check physical existence
	if existingUser != nil {
		// Business validation: Check logical existence (is it active?)
		if existingUser.RecordStatus == constants.RecordStatus.Active {
			log.Printf("[SignupUserUseCase] User already exists and is active: email=%s", dto.Email)
			return nil, errors.ErrConflict(constants.ErrorMessages.UserAlreadyExists)
		}
		// User exists but is inactive - could reactivate or return error based on business rules
		log.Printf("[SignupUserUseCase] User exists but is inactive: email=%s, status=%s", dto.Email, existingUser.RecordStatus)
		return nil, errors.ErrConflict("User account exists but is inactive")
	}

	// 3. Get default role (BUSINESS LOGIC - validate both physical and logical existence)
	defaultRole, err := uc.roleRepository.FindByCode(ctx, constants.DefaultUserRole)
	if err != nil {
		// Technical error from repository
		log.Printf("[SignupUserUseCase] Error fetching default role: %v", err)
		return nil, err
	}

	// Business validation: Check if role exists physically
	if defaultRole == nil {
		log.Printf("[SignupUserUseCase] Default role not found: code=%s", constants.DefaultUserRole)
		return nil, errors.ErrInternal(fmt.Errorf("default role '%s' not configured in system", constants.DefaultUserRole))
	}

	// Business validation: Check if role is active logically
	if defaultRole.RecordStatus != constants.RecordStatus.Active {
		log.Printf("[SignupUserUseCase] Default role is inactive: code=%s, status=%s", constants.DefaultUserRole, defaultRole.RecordStatus)
		return nil, errors.ErrInternal(fmt.Errorf("default role '%s' is not active", constants.DefaultUserRole))
	}

	log.Printf("[SignupUserUseCase] Using role: code=%s, id=%d, name=%s", defaultRole.Code, defaultRole.ID, defaultRole.Name)

	// 4. Hash password
	hashedPassword, err := hashPassword(dto.Password)
	if err != nil {
		log.Printf("[SignupUserUseCase] Error hashing password: %v", err)
		return nil, errors.ErrInternal(err)
	}

	// 5. Create user entity
	user := &entities.User{
		RoleID:           defaultRole.ID,
		Email:            dto.Email,
		PasswordHash:     hashedPassword,
		EmailVerified:    false,
		PhoneVerified:    false,
		TwoFactorEnabled: false,
		LoginAttempts:    0,
		CreatedDate:      time.Now(),
		RecordStatus:     constants.RecordStatus.Active,
	}

	// 6. Persist the user
	if err := uc.userRepository.Create(ctx, user); err != nil {
		log.Printf("[SignupUserUseCase] Error creating user: %v", err)
		return nil, err
	}

	log.Printf("[SignupUserUseCase] User created successfully: email=%s, userID=%d, roleID=%d", user.Email, user.ID, user.RoleID)
	return user, nil
}

// hashPassword generates a bcrypt hash of the given password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
