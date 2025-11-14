package repositories_test

import (
	"citary-backend/internal/domain/entities"
	"citary-backend/internal/infrastructure/persistence/postgres/repositories"
	"citary-backend/pkg/constants"
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getTestDB returns a test database connection or skips the test
func getTestDB(t *testing.T) *sql.DB {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("Skipping repository integration tests. Set TEST_DATABASE_URL environment variable to run these tests.")
	}

	db, err := sql.Open("postgres", dbURL)
	require.NoError(t, err, "Failed to connect to test database")

	err = db.Ping()
	require.NoError(t, err, "Failed to ping test database")

	return db
}

// cleanupTestData removes all test data from the database
func cleanupTestData(t *testing.T, db *sql.DB, email string) {
	_, err := db.Exec("DELETE FROM data.data_user WHERE use_email = $1", email)
	if err != nil {
		t.Logf("Warning: Failed to cleanup test data: %v", err)
	}
}

// ==========================================
// ERROR CASES (TDD - Red phase first)
// ==========================================

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	// Arrange
	db := getTestDB(t)
	defer db.Close()

	repo := repositories.NewUserRepositoryImpl(db)
	ctx := context.Background()
	nonExistentEmail := fmt.Sprintf("nonexistent_%d@example.com", time.Now().Unix())

	// Act
	user, err := repo.FindByEmail(ctx, nonExistentEmail)

	// Assert - Repository returns (nil, nil) when not found per architecture
	assert.NoError(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_FindByEmail_EmptyEmail(t *testing.T) {
	// Arrange
	db := getTestDB(t)
	defer db.Close()

	repo := repositories.NewUserRepositoryImpl(db)
	ctx := context.Background()

	// Act
	user, err := repo.FindByEmail(ctx, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	// Arrange
	db := getTestDB(t)
	defer db.Close()

	repo := repositories.NewUserRepositoryImpl(db)
	ctx := context.Background()
	testEmail := fmt.Sprintf("duplicate_%d@example.com", time.Now().Unix())
	defer cleanupTestData(t, db, testEmail)

	user1 := &entities.User{
		Email:        testEmail,
		PasswordHash: "hashed_password",
		RoleID:       1, // Assumes default role exists in test DB
		CreatedDate:  time.Now(),
		RecordStatus: constants.RecordStatus.Active,
	}

	// Create first user
	err := repo.Create(ctx, user1)
	require.NoError(t, err)

	// Try to create duplicate
	user2 := &entities.User{
		Email:        testEmail, // Same email
		PasswordHash: "different_hash",
		RoleID:       1,
		CreatedDate:  time.Now(),
		RecordStatus: constants.RecordStatus.Active,
	}

	// Act
	err = repo.Create(ctx, user2)

	// Assert
	assert.Error(t, err, "Should not allow duplicate email")
}

// ==========================================
// SUCCESS CASES (TDD - Green phase)
// ==========================================

func TestUserRepository_Create_Success(t *testing.T) {
	// Arrange
	db := getTestDB(t)
	defer db.Close()

	repo := repositories.NewUserRepositoryImpl(db)
	ctx := context.Background()
	testEmail := fmt.Sprintf("create_success_%d@example.com", time.Now().Unix())
	defer cleanupTestData(t, db, testEmail)

	createdDate := time.Now()
	user := &entities.User{
		Email:            testEmail,
		PasswordHash:     "hashed_password_123",
		RoleID:           1, // Assumes default role exists in test DB
		EmailVerified:    false,
		PhoneVerified:    false,
		TwoFactorEnabled: false,
		LoginAttempts:    0,
		CreatedDate:      createdDate,
		RecordStatus:     constants.RecordStatus.Active,
	}

	// Act
	err := repo.Create(ctx, user)

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, 0, user.ID, "ID should be set after creation")

	// Verify user was actually created in database
	foundUser, err := repo.FindByEmail(ctx, testEmail)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, testEmail, foundUser.Email)
	assert.Equal(t, 1, foundUser.RoleID)
	assert.Equal(t, "hashed_password_123", foundUser.PasswordHash)
	assert.False(t, foundUser.EmailVerified)
	assert.False(t, foundUser.PhoneVerified)
	assert.False(t, foundUser.TwoFactorEnabled)
	assert.Equal(t, 0, foundUser.LoginAttempts)
	assert.Equal(t, constants.RecordStatus.Active, foundUser.RecordStatus)
}

func TestUserRepository_FindByEmail_Success(t *testing.T) {
	// Arrange
	db := getTestDB(t)
	defer db.Close()

	repo := repositories.NewUserRepositoryImpl(db)
	ctx := context.Background()
	testEmail := fmt.Sprintf("find_success_%d@example.com", time.Now().Unix())
	defer cleanupTestData(t, db, testEmail)

	// Create a user first
	user := &entities.User{
		Email:        testEmail,
		PasswordHash: "hashed_password_456",
		RoleID:       1, // Assumes default role exists in test DB
		CreatedDate:  time.Now(),
		RecordStatus: constants.RecordStatus.Active,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Act
	foundUser, err := repo.FindByEmail(ctx, testEmail)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, 1, foundUser.RoleID)
	assert.Equal(t, testEmail, foundUser.Email)
	assert.Equal(t, "hashed_password_456", foundUser.PasswordHash)
}

func TestUserRepository_Create_WithOptionalFields(t *testing.T) {
	// Arrange
	db := getTestDB(t)
	defer db.Close()

	repo := repositories.NewUserRepositoryImpl(db)
	ctx := context.Background()
	testEmail := fmt.Sprintf("optional_fields_%d@example.com", time.Now().Unix())
	defer cleanupTestData(t, db, testEmail)

	twoFactorSecret := "SECRET123"
	lastLogin := time.Now().Add(-1 * time.Hour)
	createdDate := time.Now()

	user := &entities.User{
		Email:            testEmail,
		PasswordHash:     "hashed_password",
		RoleID:           1, // Assumes default role exists in test DB
		EmailVerified:    true,
		PhoneVerified:    true,
		TwoFactorEnabled: true,
		TwoFactorSecret:  &twoFactorSecret,
		LastLogin:        &lastLogin,
		LoginAttempts:    3,
		CreatedDate:      createdDate,
		RecordStatus:     constants.RecordStatus.Active,
	}

	// Act
	err := repo.Create(ctx, user)

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, 0, user.ID)

	// Verify optional fields were saved
	foundUser, err := repo.FindByEmail(ctx, testEmail)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, 1, foundUser.RoleID)
	assert.True(t, foundUser.EmailVerified)
	assert.True(t, foundUser.PhoneVerified)
	assert.True(t, foundUser.TwoFactorEnabled)
	assert.NotNil(t, foundUser.TwoFactorSecret)
	assert.Equal(t, twoFactorSecret, *foundUser.TwoFactorSecret)
	assert.NotNil(t, foundUser.LastLogin)
	assert.Equal(t, 3, foundUser.LoginAttempts)
}

func TestUserRepository_FindByEmail_WithNullOptionalFields(t *testing.T) {
	// Arrange
	db := getTestDB(t)
	defer db.Close()

	repo := repositories.NewUserRepositoryImpl(db)
	ctx := context.Background()
	testEmail := fmt.Sprintf("null_fields_%d@example.com", time.Now().Unix())
	defer cleanupTestData(t, db, testEmail)

	// Create user with no optional fields
	user := &entities.User{
		Email:            testEmail,
		PasswordHash:     "hashed_password",
		RoleID:           1, // Assumes default role exists in test DB
		EmailVerified:    false,
		PhoneVerified:    false,
		TwoFactorEnabled: false,
		TwoFactorSecret:  nil,
		LastLogin:        nil,
		LoginAttempts:    0,
		LockedUntil:      nil,
		TermsAcceptedAt:  nil,
		PrivacyAcceptedAt: nil,
		CreatedDate:      time.Now(),
		RecordStatus:     constants.RecordStatus.Active,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Act
	foundUser, err := repo.FindByEmail(ctx, testEmail)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, 1, foundUser.RoleID)
	assert.Nil(t, foundUser.TwoFactorSecret)
	assert.Nil(t, foundUser.LastLogin)
	assert.Nil(t, foundUser.LockedUntil)
	assert.Nil(t, foundUser.TermsAcceptedAt)
	assert.Nil(t, foundUser.PrivacyAcceptedAt)
}

func TestUserRepository_RoundTrip_PreservesData(t *testing.T) {
	// Arrange
	db := getTestDB(t)
	defer db.Close()

	repo := repositories.NewUserRepositoryImpl(db)
	ctx := context.Background()
	testEmail := fmt.Sprintf("roundtrip_%d@example.com", time.Now().Unix())
	defer cleanupTestData(t, db, testEmail)

	twoFactorSecret := "SECRETABC"
	originalUser := &entities.User{
		Email:            testEmail,
		PasswordHash:     "original_hash_789",
		RoleID:           1, // Assumes default role exists in test DB
		EmailVerified:    true,
		PhoneVerified:    false,
		TwoFactorEnabled: true,
		TwoFactorSecret:  &twoFactorSecret,
		LoginAttempts:    5,
		CreatedDate:      time.Now(),
		RecordStatus:     constants.RecordStatus.Active,
	}

	// Act - Create and then Find
	err := repo.Create(ctx, originalUser)
	require.NoError(t, err)

	foundUser, err := repo.FindByEmail(ctx, testEmail)
	require.NoError(t, err)

	// Assert - All data should be preserved
	assert.Equal(t, originalUser.ID, foundUser.ID)
	assert.Equal(t, originalUser.RoleID, foundUser.RoleID)
	assert.Equal(t, originalUser.Email, foundUser.Email)
	assert.Equal(t, originalUser.PasswordHash, foundUser.PasswordHash)
	assert.Equal(t, originalUser.EmailVerified, foundUser.EmailVerified)
	assert.Equal(t, originalUser.PhoneVerified, foundUser.PhoneVerified)
	assert.Equal(t, originalUser.TwoFactorEnabled, foundUser.TwoFactorEnabled)
	assert.Equal(t, *originalUser.TwoFactorSecret, *foundUser.TwoFactorSecret)
	assert.Equal(t, originalUser.LoginAttempts, foundUser.LoginAttempts)
	assert.Equal(t, originalUser.RecordStatus, foundUser.RecordStatus)
}

func TestNewUserRepositoryImpl_ReturnsValidInstance(t *testing.T) {
	// Arrange
	db := getTestDB(t)
	defer db.Close()

	// Act
	repo := repositories.NewUserRepositoryImpl(db)

	// Assert
	assert.NotNil(t, repo)
}
