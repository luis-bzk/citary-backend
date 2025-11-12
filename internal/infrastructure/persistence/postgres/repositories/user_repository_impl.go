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
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	start := time.Now()
	log.Printf("[UserRepository] FindByEmail: email=%s", email)

	query := `
		SELECT use_id, use_email, use_password_hash, use_email_verified,
		       use_phone_verified, use_two_factor_enabled, use_two_factor_secret,
		       use_last_login, use_login_attempts, use_locked_until,
		       use_terms_accepted_at, use_privacy_accepted_at, use_created_date, use_record_status
		FROM data.data_user
		WHERE use_email = $1`

	var dbEntity dbEntities.UserDB

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&dbEntity.UseID,
		&dbEntity.UseEmail,
		&dbEntity.UsePasswordHash,
		&dbEntity.UseEmailVerified,
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

	if err == sql.ErrNoRows {
		log.Printf("[UserRepository] FindByEmail: user not found, email=%s, duration=%v", email, duration)
		return nil, errors.ErrNotFound("User not found")
	}

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
	log.Printf("[UserRepository] Create: email=%s, emailVerified=%v", user.Email, user.EmailVerified)

	query := `
		INSERT INTO data.data_user (
			use_email, use_password_hash, use_email_verified, use_phone_verified,
			use_two_factor_enabled, use_login_attempts, use_created_date, use_record_status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING use_id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.EmailVerified,
		user.PhoneVerified,
		user.TwoFactorEnabled,
		user.LoginAttempts,
		user.CreatedDate,
		user.RecordStatus,
	).Scan(&user.ID)

	duration := time.Since(start)

	if err != nil {
		log.Printf("[UserRepository] Create ERROR: email=%s, error=%v, duration=%v", user.Email, err, duration)
		return errors.ErrInternal(err)
	}

	log.Printf("[UserRepository] Create: success, email=%s, userID=%d, duration=%v", user.Email, user.ID, duration)
	return nil
}
