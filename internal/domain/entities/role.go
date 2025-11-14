package entities

import (
	"encoding/json"
	"time"
)

// Role represents a user role entity in the domain layer
type Role struct {
	ID           int
	Name         string
	Code         string
	Description  *string
	Permissions  *json.RawMessage
	CreatedDate  time.Time
	RecordStatus string
}
