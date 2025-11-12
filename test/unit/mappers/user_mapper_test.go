package mappers_test

import (
	"citary-backend/internal/domain/entities"
	dbEntities "citary-backend/internal/infrastructure/persistence/postgres/entities"
	"citary-backend/internal/infrastructure/persistence/postgres/mappers"
	"citary-backend/pkg/constants"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ==========================================
// ToDomainEntity TESTS (Error/Edge cases first)
// ==========================================

func TestUserMapper_ToDomainEntity_NullOptionalFields(t *testing.T) {
	// Arrange
	mapper := mappers.NewUserMapper()
	createdDate := time.Now()

	dbEntity := &dbEntities.UserDB{
		UseID:                1,
		UseEmail:             "test@example.com",
		UsePasswordHash:      "hashed_password",
		UseEmailVerified:     false,
		UsePhoneVerified:     false,
		UseTwoFactorEnabled:  false,
		UseTwoFactorSecret:   sql.NullString{Valid: false}, // NULL
		UseLastLogin:         sql.NullTime{Valid: false},   // NULL
		UseLoginAttempts:     0,
		UseLockedUntil:       sql.NullTime{Valid: false}, // NULL
		UseTermsAcceptedAt:   sql.NullTime{Valid: false}, // NULL
		UsePrivacyAcceptedAt: sql.NullTime{Valid: false}, // NULL
		UseCreatedDate:       createdDate,
		UseRecordStatus:      constants.RecordStatus.Active,
	}

	// Act
	domainEntity := mapper.ToDomainEntity(dbEntity)

	// Assert
	assert.NotNil(t, domainEntity)
	assert.Equal(t, 1, domainEntity.ID)
	assert.Equal(t, "test@example.com", domainEntity.Email)
	assert.Equal(t, "hashed_password", domainEntity.PasswordHash)
	assert.False(t, domainEntity.EmailVerified)
	assert.False(t, domainEntity.PhoneVerified)
	assert.False(t, domainEntity.TwoFactorEnabled)
	assert.Nil(t, domainEntity.TwoFactorSecret, "TwoFactorSecret should be nil")
	assert.Nil(t, domainEntity.LastLogin, "LastLogin should be nil")
	assert.Equal(t, 0, domainEntity.LoginAttempts)
	assert.Nil(t, domainEntity.LockedUntil, "LockedUntil should be nil")
	assert.Nil(t, domainEntity.TermsAcceptedAt, "TermsAcceptedAt should be nil")
	assert.Nil(t, domainEntity.PrivacyAcceptedAt, "PrivacyAcceptedAt should be nil")
	assert.Equal(t, createdDate, domainEntity.CreatedDate)
	assert.Equal(t, constants.RecordStatus.Active, domainEntity.RecordStatus)
}

func TestUserMapper_ToDomainEntity_WithAllOptionalFields(t *testing.T) {
	// Arrange
	mapper := mappers.NewUserMapper()
	createdDate := time.Now()
	lastLogin := time.Now().Add(-1 * time.Hour)
	lockedUntil := time.Now().Add(1 * time.Hour)
	termsAccepted := time.Now().Add(-2 * time.Hour)
	privacyAccepted := time.Now().Add(-2 * time.Hour)
	twoFactorSecret := "SECRET123"

	dbEntity := &dbEntities.UserDB{
		UseID:                2,
		UseEmail:             "full@example.com",
		UsePasswordHash:      "hashed_password",
		UseEmailVerified:     true,
		UsePhoneVerified:     true,
		UseTwoFactorEnabled:  true,
		UseTwoFactorSecret:   sql.NullString{String: twoFactorSecret, Valid: true},
		UseLastLogin:         sql.NullTime{Time: lastLogin, Valid: true},
		UseLoginAttempts:     3,
		UseLockedUntil:       sql.NullTime{Time: lockedUntil, Valid: true},
		UseTermsAcceptedAt:   sql.NullTime{Time: termsAccepted, Valid: true},
		UsePrivacyAcceptedAt: sql.NullTime{Time: privacyAccepted, Valid: true},
		UseCreatedDate:       createdDate,
		UseRecordStatus:      constants.RecordStatus.Active,
	}

	// Act
	domainEntity := mapper.ToDomainEntity(dbEntity)

	// Assert
	assert.NotNil(t, domainEntity)
	assert.Equal(t, 2, domainEntity.ID)
	assert.Equal(t, "full@example.com", domainEntity.Email)
	assert.True(t, domainEntity.EmailVerified)
	assert.True(t, domainEntity.PhoneVerified)
	assert.True(t, domainEntity.TwoFactorEnabled)
	assert.NotNil(t, domainEntity.TwoFactorSecret)
	assert.Equal(t, twoFactorSecret, *domainEntity.TwoFactorSecret)
	assert.NotNil(t, domainEntity.LastLogin)
	assert.Equal(t, lastLogin, *domainEntity.LastLogin)
	assert.Equal(t, 3, domainEntity.LoginAttempts)
	assert.NotNil(t, domainEntity.LockedUntil)
	assert.Equal(t, lockedUntil, *domainEntity.LockedUntil)
	assert.NotNil(t, domainEntity.TermsAcceptedAt)
	assert.Equal(t, termsAccepted, *domainEntity.TermsAcceptedAt)
	assert.NotNil(t, domainEntity.PrivacyAcceptedAt)
	assert.Equal(t, privacyAccepted, *domainEntity.PrivacyAcceptedAt)
}

// ==========================================
// ToDBEntity TESTS (Edge cases first)
// ==========================================

func TestUserMapper_ToDBEntity_WithNilPointers(t *testing.T) {
	// Arrange
	mapper := mappers.NewUserMapper()
	createdDate := time.Now()

	domainEntity := &entities.User{
		ID:                1,
		Email:             "test@example.com",
		PasswordHash:      "hashed_password",
		EmailVerified:     false,
		PhoneVerified:     false,
		TwoFactorEnabled:  false,
		TwoFactorSecret:   nil, // nil pointer
		LastLogin:         nil, // nil pointer
		LoginAttempts:     0,
		LockedUntil:       nil, // nil pointer
		TermsAcceptedAt:   nil, // nil pointer
		PrivacyAcceptedAt: nil, // nil pointer
		CreatedDate:       createdDate,
		RecordStatus:      constants.RecordStatus.Active,
	}

	// Act
	dbEntity := mapper.ToDBEntity(domainEntity)

	// Assert
	assert.NotNil(t, dbEntity)
	assert.Equal(t, 1, dbEntity.UseID)
	assert.Equal(t, "test@example.com", dbEntity.UseEmail)
	assert.Equal(t, "hashed_password", dbEntity.UsePasswordHash)
	assert.False(t, dbEntity.UseEmailVerified)
	assert.False(t, dbEntity.UsePhoneVerified)
	assert.False(t, dbEntity.UseTwoFactorEnabled)
	assert.False(t, dbEntity.UseTwoFactorSecret.Valid, "TwoFactorSecret should be invalid")
	assert.False(t, dbEntity.UseLastLogin.Valid, "LastLogin should be invalid")
	assert.Equal(t, 0, dbEntity.UseLoginAttempts)
	assert.False(t, dbEntity.UseLockedUntil.Valid, "LockedUntil should be invalid")
	assert.False(t, dbEntity.UseTermsAcceptedAt.Valid, "TermsAcceptedAt should be invalid")
	assert.False(t, dbEntity.UsePrivacyAcceptedAt.Valid, "PrivacyAcceptedAt should be invalid")
	assert.Equal(t, createdDate, dbEntity.UseCreatedDate)
	assert.Equal(t, constants.RecordStatus.Active, dbEntity.UseRecordStatus)
}

func TestUserMapper_ToDBEntity_WithAllPointers(t *testing.T) {
	// Arrange
	mapper := mappers.NewUserMapper()
	createdDate := time.Now()
	lastLogin := time.Now().Add(-1 * time.Hour)
	lockedUntil := time.Now().Add(1 * time.Hour)
	termsAccepted := time.Now().Add(-2 * time.Hour)
	privacyAccepted := time.Now().Add(-2 * time.Hour)
	twoFactorSecret := "SECRET123"

	domainEntity := &entities.User{
		ID:                2,
		Email:             "full@example.com",
		PasswordHash:      "hashed_password",
		EmailVerified:     true,
		PhoneVerified:     true,
		TwoFactorEnabled:  true,
		TwoFactorSecret:   &twoFactorSecret,
		LastLogin:         &lastLogin,
		LoginAttempts:     5,
		LockedUntil:       &lockedUntil,
		TermsAcceptedAt:   &termsAccepted,
		PrivacyAcceptedAt: &privacyAccepted,
		CreatedDate:       createdDate,
		RecordStatus:      constants.RecordStatus.Active,
	}

	// Act
	dbEntity := mapper.ToDBEntity(domainEntity)

	// Assert
	assert.NotNil(t, dbEntity)
	assert.Equal(t, 2, dbEntity.UseID)
	assert.Equal(t, "full@example.com", dbEntity.UseEmail)
	assert.True(t, dbEntity.UseEmailVerified)
	assert.True(t, dbEntity.UsePhoneVerified)
	assert.True(t, dbEntity.UseTwoFactorEnabled)
	assert.True(t, dbEntity.UseTwoFactorSecret.Valid)
	assert.Equal(t, twoFactorSecret, dbEntity.UseTwoFactorSecret.String)
	assert.True(t, dbEntity.UseLastLogin.Valid)
	assert.Equal(t, lastLogin, dbEntity.UseLastLogin.Time)
	assert.Equal(t, 5, dbEntity.UseLoginAttempts)
	assert.True(t, dbEntity.UseLockedUntil.Valid)
	assert.Equal(t, lockedUntil, dbEntity.UseLockedUntil.Time)
	assert.True(t, dbEntity.UseTermsAcceptedAt.Valid)
	assert.Equal(t, termsAccepted, dbEntity.UseTermsAcceptedAt.Time)
	assert.True(t, dbEntity.UsePrivacyAcceptedAt.Valid)
	assert.Equal(t, privacyAccepted, dbEntity.UsePrivacyAcceptedAt.Time)
}

// ==========================================
// ROUND-TRIP TESTS (Success cases)
// ==========================================

func TestUserMapper_RoundTrip_PreservesData(t *testing.T) {
	// Arrange
	mapper := mappers.NewUserMapper()
	createdDate := time.Now()
	lastLogin := time.Now().Add(-1 * time.Hour)
	twoFactorSecret := "SECRET123"

	originalDomainEntity := &entities.User{
		ID:                3,
		Email:             "roundtrip@example.com",
		PasswordHash:      "hashed_password",
		EmailVerified:     true,
		PhoneVerified:     false,
		TwoFactorEnabled:  true,
		TwoFactorSecret:   &twoFactorSecret,
		LastLogin:         &lastLogin,
		LoginAttempts:     2,
		LockedUntil:       nil,
		TermsAcceptedAt:   nil,
		PrivacyAcceptedAt: nil,
		CreatedDate:       createdDate,
		RecordStatus:      constants.RecordStatus.Active,
	}

	// Act - Convert Domain -> DB -> Domain
	dbEntity := mapper.ToDBEntity(originalDomainEntity)
	resultDomainEntity := mapper.ToDomainEntity(dbEntity)

	// Assert - All fields should be preserved
	assert.Equal(t, originalDomainEntity.ID, resultDomainEntity.ID)
	assert.Equal(t, originalDomainEntity.Email, resultDomainEntity.Email)
	assert.Equal(t, originalDomainEntity.PasswordHash, resultDomainEntity.PasswordHash)
	assert.Equal(t, originalDomainEntity.EmailVerified, resultDomainEntity.EmailVerified)
	assert.Equal(t, originalDomainEntity.PhoneVerified, resultDomainEntity.PhoneVerified)
	assert.Equal(t, originalDomainEntity.TwoFactorEnabled, resultDomainEntity.TwoFactorEnabled)
	assert.Equal(t, *originalDomainEntity.TwoFactorSecret, *resultDomainEntity.TwoFactorSecret)
	assert.Equal(t, *originalDomainEntity.LastLogin, *resultDomainEntity.LastLogin)
	assert.Equal(t, originalDomainEntity.LoginAttempts, resultDomainEntity.LoginAttempts)
	assert.Nil(t, resultDomainEntity.LockedUntil)
	assert.Nil(t, resultDomainEntity.TermsAcceptedAt)
	assert.Nil(t, resultDomainEntity.PrivacyAcceptedAt)
	assert.Equal(t, originalDomainEntity.CreatedDate, resultDomainEntity.CreatedDate)
	assert.Equal(t, originalDomainEntity.RecordStatus, resultDomainEntity.RecordStatus)
}

func TestUserMapper_RoundTrip_WithAllFieldsNull(t *testing.T) {
	// Arrange
	mapper := mappers.NewUserMapper()
	createdDate := time.Now()

	originalDomainEntity := &entities.User{
		ID:                4,
		Email:             "minimal@example.com",
		PasswordHash:      "hashed_password",
		EmailVerified:     false,
		PhoneVerified:     false,
		TwoFactorEnabled:  false,
		TwoFactorSecret:   nil,
		LastLogin:         nil,
		LoginAttempts:     0,
		LockedUntil:       nil,
		TermsAcceptedAt:   nil,
		PrivacyAcceptedAt: nil,
		CreatedDate:       createdDate,
		RecordStatus:      constants.RecordStatus.Active,
	}

	// Act - Convert Domain -> DB -> Domain
	dbEntity := mapper.ToDBEntity(originalDomainEntity)
	resultDomainEntity := mapper.ToDomainEntity(dbEntity)

	// Assert - All fields should be preserved
	assert.Equal(t, originalDomainEntity.ID, resultDomainEntity.ID)
	assert.Equal(t, originalDomainEntity.Email, resultDomainEntity.Email)
	assert.Nil(t, resultDomainEntity.TwoFactorSecret)
	assert.Nil(t, resultDomainEntity.LastLogin)
	assert.Nil(t, resultDomainEntity.LockedUntil)
	assert.Nil(t, resultDomainEntity.TermsAcceptedAt)
	assert.Nil(t, resultDomainEntity.PrivacyAcceptedAt)
}

func TestNewUserMapper_ReturnsValidInstance(t *testing.T) {
	// Act
	mapper := mappers.NewUserMapper()

	// Assert
	assert.NotNil(t, mapper, "NewUserMapper should return a non-nil instance")
}
