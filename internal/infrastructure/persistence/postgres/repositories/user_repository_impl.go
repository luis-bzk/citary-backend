package repositories

import (
	"citary-backend/internal/domain/entities"
	"citary-backend/internal/domain/errors"
	dbEntities "citary-backend/internal/infrastructure/persistence/postgres/entities"
	"citary-backend/internal/infrastructure/persistence/postgres/mappers"
	"context"
	"database/sql"
	"log"
	"time"
)

// UserRepositoryImpl implements the UserRepository interface using PostgreSQL
type UserRepositoryImpl struct {
	db     *sql.DB
	mapper *mappers.UserMapper
}

// NewUserRepositoryImpl creates a new instance of UserRepositoryImpl
func NewUserRepositoryImpl(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db:     db,
		mapper: mappers.NewUserMapper(),
	}
}

// FindByEmail retrieves a user by their email address
// Returns (nil, nil) if not found - business layer decides if that's an error
// Returns (nil, error) only on technical failures (DB connection, query errors, etc.)
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	start := time.Now()
	log.Printf("[UserRepository] FindByEmail: email=%s", email)

	query := `
		SELECT use_id, id_role, use_email, use_password_hash, use_email_verified,
		       use_verification_token, use_verification_token_expires_at,
		       use_phone_verified, use_two_factor_enabled, use_two_factor_secret,
		       use_last_login, use_login_attempts, use_locked_until,
		       use_terms_accepted_at, use_privacy_accepted_at, use_created_date, use_record_status
		FROM data.data_user
		WHERE use_email = $1`

	var dbEntity dbEntities.UserDB

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&dbEntity.UseID,
		&dbEntity.IdRole,
		&dbEntity.UseEmail,
		&dbEntity.UsePasswordHash,
		&dbEntity.UseEmailVerified,
		&dbEntity.UseVerificationToken,
		&dbEntity.UseVerificationTokenExpiresAt,
		&dbEntity.UsePhoneVerified,
		&dbEntity.UseTwoFactorEnabled,
		&dbEntity.UseTwoFactorSecret,
		&dbEntity.UseLastLogin,
		&dbEntity.UseLoginAttempts,
		&dbEntity.UseLockedUntil,
		&dbEntity.UseTermsAcceptedAt,
		&dbEntity.UsePrivacyAcceptedAt,
		&dbEntity.UseCreatedDate,
		&dbEntity.UseRecordStatus,
	)

	duration := time.Since(start)

	// Not found is NOT an error at infrastructure level - it's a valid result
	if err == sql.ErrNoRows {
		log.Printf("[UserRepository] FindByEmail: user not found, email=%s, duration=%v", email, duration)
		return nil, nil
	}

	// Technical errors (DB connection, query syntax, etc.) ARE errors
	if err != nil {
		log.Printf("[UserRepository] FindByEmail ERROR: email=%s, error=%v, duration=%v", email, err, duration)
		return nil, errors.ErrInternal(err)
	}

	log.Printf("[UserRepository] FindByEmail: success, email=%s, userID=%d, duration=%v", email, dbEntity.UseID, duration)
	return r.mapper.ToDomainEntity(&dbEntity), nil
}

// Create persists a new user to the database
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	start := time.Now()
	log.Printf("[UserRepository] Create: email=%s, roleID=%d, emailVerified=%v", user.Email, user.RoleID, user.EmailVerified)

	// Convert domain entity to DB entity to handle nullable fields properly
	dbEntity := r.mapper.ToDBEntity(user)

	query := `
		INSERT INTO data.data_user (
			id_role, use_email, use_password_hash, use_email_verified,
			use_verification_token, use_verification_token_expires_at,
			use_phone_verified, use_two_factor_enabled, use_login_attempts,
			use_created_date, use_record_status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING use_id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		dbEntity.IdRole,
		dbEntity.UseEmail,
		dbEntity.UsePasswordHash,
		dbEntity.UseEmailVerified,
		dbEntity.UseVerificationToken,
		dbEntity.UseVerificationTokenExpiresAt,
		dbEntity.UsePhoneVerified,
		dbEntity.UseTwoFactorEnabled,
		dbEntity.UseLoginAttempts,
		dbEntity.UseCreatedDate,
		dbEntity.UseRecordStatus,
	).Scan(&user.ID)

	duration := time.Since(start)

	if err != nil {
		log.Printf("[UserRepository] Create ERROR: email=%s, roleID=%d, error=%v, duration=%v", user.Email, user.RoleID, err, duration)
		return errors.ErrInternal(err)
	}

	log.Printf("[UserRepository] Create: success, email=%s, roleID=%d, userID=%d, duration=%v", user.Email, user.RoleID, user.ID, duration)
	return nil
}
