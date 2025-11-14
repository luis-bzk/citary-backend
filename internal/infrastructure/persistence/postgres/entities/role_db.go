package entities

import (
	"database/sql"
	"time"
)

// RoleDB represents the role table structure in PostgreSQL
type RoleDB struct {
	RolID           int            `db:"rol_id"`
	RolName         string         `db:"rol_name"`
	RolCode         string         `db:"rol_code"`
	RolDescription  sql.NullString `db:"rol_description"`
	RolPermissions  []byte         `db:"rol_permissions"` // JSONB
	RolCreatedDate  time.Time      `db:"rol_created_date"`
	RolRecordStatus string         `db:"rol_record_status"`
}
