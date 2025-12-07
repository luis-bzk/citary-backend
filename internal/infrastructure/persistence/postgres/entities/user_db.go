package entities

import (
	"database/sql"
	"time"
)

// UserDB represents the user table structure in PostgreSQL
type UserDB struct {
	UseID                         int            `db:"use_id"`
	IdRole                        int            `db:"id_role"`
	UseEmail                      string         `db:"use_email"`
	UsePasswordHash               string         `db:"use_password_hash"`
	UseEmailVerified              bool           `db:"use_email_verified"`
	UseVerificationToken          sql.NullString `db:"use_verification_token"`
	UseVerificationTokenExpiresAt sql.NullTime   `db:"use_verification_token_expires_at"`
	UseLastLogin                  sql.NullTime   `db:"use_last_login"`
	UseLoginAttempts              int            `db:"use_login_attempts"`
	UseLockedUntil                sql.NullTime   `db:"use_locked_until"`
	UseTermsAcceptedAt            sql.NullTime   `db:"use_terms_accepted_at"`
	UsePrivacyAcceptedAt          sql.NullTime   `db:"use_privacy_accepted_at"`
	UseCreatedDate                time.Time      `db:"use_created_date"`
	UseRecordStatus               string         `db:"use_record_status"`
}
