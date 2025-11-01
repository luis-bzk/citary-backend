package entities

import (
	"citary-backend/pkg/constants"
	"time"
)

// User represents the user entity in the domain layer
type User struct {
	ID                int
	Email             string
	PasswordHash      string
	EmailVerified     bool
	PhoneVerified     bool
	TwoFactorEnabled  bool
	TwoFactorSecret   *string
	LastLogin         *time.Time
	LoginAttempts     int
	LockedUntil       *time.Time
	TermsAcceptedAt   *time.Time
	PrivacyAcceptedAt *time.Time
	CreatedDate       time.Time
	RecordStatus      string
}

// IsActive checks if the user account is active
func (u *User) IsActive() bool {
	return u.RecordStatus == constants.RecordStatus.Active
}

// IsLocked checks if the user account is currently locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return u.LockedUntil.After(time.Now())
}
