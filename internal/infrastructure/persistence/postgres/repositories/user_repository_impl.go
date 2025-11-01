package repositories

import (
	"citary-backend/internal/domain/entities"
	"citary-backend/internal/domain/errors"
	dbEntities "citary-backend/internal/infrastructure/persistence/postgres/entities"
	"citary-backend/internal/infrastructure/persistence/postgres/mappers"
	"context"
	"database/sql"
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

	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound("User not found")
	}

	if err != nil {
		return nil, errors.ErrInternal(err)
	}

	return r.mapper.ToDomainEntity(&dbEntity), nil
}

// Create persists a new user to the database
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
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

	if err != nil {
		return errors.ErrInternal(err)
	}

	return nil
}
