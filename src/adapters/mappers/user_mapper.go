package mappers

import (
	"citary-backend/src/data/postgresql/db_entities"
	"citary-backend/src/domain/entities"
	"database/sql"
)

type UserMapper struct{}

func (m *UserMapper) ToDBEntity(user *entities.User) *db_entities.UserDBEntity {
	dbEntity := &db_entities.UserDBEntity{
		UseID:               user.ID,
		UseEmail:            user.Email,
		UsePasswordHash:     user.PasswordHash,
		UseEmailVerified:    user.EmailVerified,
		UsePhoneVerified:    user.PhoneVerified,
		UseTwoFactorEnabled: user.TwoFactorEnabled,
		UseLoginAttempts:    user.LoginAttempts,
		UseCreatedDate:      user.CreatedDate,
		UseRecordStatus:     user.RecordStatus,
	}

	if user.TwoFactorSecret != nil {
		dbEntity.UseTwoFactorSecret = sql.NullString{String: *user.TwoFactorSecret, Valid: true}
	}

	if user.LastLogin != nil {
		dbEntity.UseLastLogin = sql.NullTime{Time: *user.LastLogin, Valid: true}
	}

	if user.LockedUntil != nil {
		dbEntity.UseLockedUntil = sql.NullTime{Time: *user.LockedUntil, Valid: true}
	}

	if user.TermsAcceptedAt != nil {
		dbEntity.UseTermsAcceptedAt = sql.NullTime{Time: *user.TermsAcceptedAt, Valid: true}
	}

	if user.PrivacyAcceptedAt != nil {
		dbEntity.UsePrivacyAcceptedAt = sql.NullTime{Time: *user.PrivacyAcceptedAt, Valid: true}
	}

	return dbEntity
}

func (m *UserMapper) ToDomainEntity(dbEntity *db_entities.UserDBEntity) *entities.User {
	user := &entities.User{
		ID:               dbEntity.UseID,
		Email:            dbEntity.UseEmail,
		PasswordHash:     dbEntity.UsePasswordHash,
		EmailVerified:    dbEntity.UseEmailVerified,
		PhoneVerified:    dbEntity.UsePhoneVerified,
		TwoFactorEnabled: dbEntity.UseTwoFactorEnabled,
		LoginAttempts:    dbEntity.UseLoginAttempts,
		CreatedDate:      dbEntity.UseCreatedDate,
		RecordStatus:     dbEntity.UseRecordStatus,
	}

	// Manejar campos opcionales
	if dbEntity.UseTwoFactorSecret.Valid {
		secret := dbEntity.UseTwoFactorSecret.String
		user.TwoFactorSecret = &secret
	}

	if dbEntity.UseLastLogin.Valid {
		lastLogin := dbEntity.UseLastLogin.Time
		user.LastLogin = &lastLogin
	}

	if dbEntity.UseLockedUntil.Valid {
		lockedUntil := dbEntity.UseLockedUntil.Time
		user.LockedUntil = &lockedUntil
	}

	if dbEntity.UseTermsAcceptedAt.Valid {
		termsAccepted := dbEntity.UseTermsAcceptedAt.Time
		user.TermsAcceptedAt = &termsAccepted
	}

	if dbEntity.UsePrivacyAcceptedAt.Valid {
		privacyAccepted := dbEntity.UsePrivacyAcceptedAt.Time
		user.PrivacyAcceptedAt = &privacyAccepted
	}

	return user
}
