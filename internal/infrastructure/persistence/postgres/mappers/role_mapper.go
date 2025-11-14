package mappers

import (
	"citary-backend/internal/domain/entities"
	dbEntities "citary-backend/internal/infrastructure/persistence/postgres/entities"
	"encoding/json"
)

// RoleMapper handles conversion between domain and database entities
type RoleMapper struct{}

// NewRoleMapper creates a new RoleMapper instance
func NewRoleMapper() *RoleMapper {
	return &RoleMapper{}
}

// ToDomainEntity converts a database RoleDB entity to a domain Role entity
func (m *RoleMapper) ToDomainEntity(dbEntity *dbEntities.RoleDB) *entities.Role {
	role := &entities.Role{
		ID:           dbEntity.RolID,
		Name:         dbEntity.RolName,
		Code:         dbEntity.RolCode,
		CreatedDate:  dbEntity.RolCreatedDate,
		RecordStatus: dbEntity.RolRecordStatus,
	}

	// Handle optional description
	if dbEntity.RolDescription.Valid {
		desc := dbEntity.RolDescription.String
		role.Description = &desc
	}

	// Handle optional permissions (JSONB)
	if dbEntity.RolPermissions != nil && len(dbEntity.RolPermissions) > 0 {
		rawMsg := json.RawMessage(dbEntity.RolPermissions)
		role.Permissions = &rawMsg
	}

	return role
}
