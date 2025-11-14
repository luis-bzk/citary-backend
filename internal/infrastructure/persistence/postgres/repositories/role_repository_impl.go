package repositories

import (
	"citary-backend/internal/domain/entities"
	"citary-backend/internal/domain/errors"
	dbEntities "citary-backend/internal/infrastructure/persistence/postgres/entities"
	"citary-backend/internal/infrastructure/persistence/postgres/mappers"
	"context"
	"database/sql"
	"log"
	"time"
)

// RoleRepositoryImpl implements the RoleRepository interface using PostgreSQL
type RoleRepositoryImpl struct {
	db     *sql.DB
	mapper *mappers.RoleMapper
}

// NewRoleRepositoryImpl creates a new instance of RoleRepositoryImpl
func NewRoleRepositoryImpl(db *sql.DB) *RoleRepositoryImpl {
	return &RoleRepositoryImpl{
		db:     db,
		mapper: mappers.NewRoleMapper(),
	}
}

// FindByCode retrieves a role by its code
// Returns (nil, nil) if not found - business layer decides if that's an error
// Returns (nil, error) only on technical failures (DB connection, query errors, etc.)
func (r *RoleRepositoryImpl) FindByCode(ctx context.Context, code string) (*entities.Role, error) {
	start := time.Now()
	log.Printf("[RoleRepository] FindByCode: code=%s", code)

	query := `
		SELECT rol_id, rol_name, rol_code, rol_description, rol_permissions,
		       rol_created_date, rol_record_status
		FROM core.core_role
		WHERE rol_code = $1`

	var dbEntity dbEntities.RoleDB

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&dbEntity.RolID,
		&dbEntity.RolName,
		&dbEntity.RolCode,
		&dbEntity.RolDescription,
		&dbEntity.RolPermissions,
		&dbEntity.RolCreatedDate,
		&dbEntity.RolRecordStatus,
	)

	duration := time.Since(start)

	// Not found is NOT an error at infrastructure level - it's a valid result
	if err == sql.ErrNoRows {
		log.Printf("[RoleRepository] FindByCode: role not found, code=%s, duration=%v", code, duration)
		return nil, nil
	}

	// Technical errors (DB connection, query syntax, etc.) ARE errors
	if err != nil {
		log.Printf("[RoleRepository] FindByCode ERROR: code=%s, error=%v, duration=%v", code, err, duration)
		return nil, errors.ErrInternal(err)
	}

	log.Printf("[RoleRepository] FindByCode: success, code=%s, roleID=%d, roleName=%s, status=%s, duration=%v",
		code, dbEntity.RolID, dbEntity.RolName, dbEntity.RolRecordStatus, duration)
	return r.mapper.ToDomainEntity(&dbEntity), nil
}
