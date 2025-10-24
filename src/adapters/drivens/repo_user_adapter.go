package adaptersDrivens

import (
	"citary-backend/src/adapters/mappers"
	"citary-backend/src/data/postgresql/db_entities"
	"citary-backend/src/domain/entities"
	"citary-backend/src/domain/errors"
	"database/sql"
)

type UserRepositoryAdapter struct {
	db     *sql.DB
	mapper *mappers.UserMapper
}

func NewUserUserRepositoryImpl(db *sql.DB) *UserRepositoryAdapter {
	return &UserRepositoryAdapter{
		db:     db,
		mapper: &mappers.UserMapper{},
	}
}

func (r *UserRepositoryAdapter) FindByEmail(email string) (*entities.User, error) {
	query := `
		SELECT use_id, use_email, use_password_hash, use_email_verified,
		       use_phone_verified, use_two_factor_enabled, use_two_factor_secret,
		       use_last_login, use_login_attempts, use_locked_until,
		       use_terms_accepted_at, use_privacy_accepted_at, use_created_date, use_record_status
		FROM data.data_user
		WHERE use_email = $1`

	var dbEntity db_entities.UserDBEntity

	err := r.db.QueryRow(query, email).Scan(
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
		return nil, errors.ErrNotFound("Usuario no encontrado")
	}

	if err != nil {
		return nil, errors.ErrInternal(err)
	}

	return r.mapper.ToDomainEntity(&dbEntity), nil
}

func (r *UserRepositoryAdapter) Create(user *entities.User) error {
	query := `
		INSERT INTO data.data_user (
			use_email, use_password_hash, use_email_verified, use_phone_verified,
			use_two_factor_enabled, use_login_attempts, use_created_date, use_record_status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING use_id
	`

	err := r.db.QueryRow(
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
