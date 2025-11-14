package mappers

import (
	domainEntities "citary-backend/internal/domain/entities"
	dbEntities "citary-backend/internal/infrastructure/persistence/postgres/entities"
	"database/sql"
)

// UserMapper handles conversion between domain and database entities
type UserMapper struct{}

// NewUserMapper creates a new UserMapper instance
func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

// ToDBEntity converts a domain User entity to a database UserDB entity
func (m *UserMapper) ToDBEntity(user *domainEntities.User) *dbEntities.UserDB {
	dbEntity := &dbEntities.UserDB{
		UseID:               user.ID,
		IdRole:              user.RoleID,
		UseEmail:            user.Email,
		UsePasswordHash:     user.PasswordHash,
		UseEmailVerified:    user.EmailVerified,
		UsePhoneVerified:    user.PhoneVerified,
		UseTwoFactorEnabled: user.TwoFactorEnabled,
		UseLoginAttempts:    user.LoginAttempts,
		UseCreatedDate:      user.CreatedDate,
		UseRecordStatus:     user.RecordStatus,
	}

	// Handle optional fields
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

// ToDomainEntity converts a database UserDB entity to a domain User entity
func (m *UserMapper) ToDomainEntity(dbEntity *dbEntities.UserDB) *domainEntities.User {
	user := &domainEntities.User{
		ID:               dbEntity.UseID,
		RoleID:           dbEntity.IdRole,
		Email:            dbEntity.UseEmail,
		PasswordHash:     dbEntity.UsePasswordHash,
		EmailVerified:    dbEntity.UseEmailVerified,
		PhoneVerified:    dbEntity.UsePhoneVerified,
		TwoFactorEnabled: dbEntity.UseTwoFactorEnabled,
		LoginAttempts:    dbEntity.UseLoginAttempts,
		CreatedDate:      dbEntity.UseCreatedDate,
		RecordStatus:     dbEntity.UseRecordStatus,
	}

	// Handle optional fields
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
