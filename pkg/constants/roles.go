package constants

// RoleCodes defines the available user roles in the system
var RoleCodes = struct {
	SuperAdmin        string
	Staff             string
	Patient           string
	OrganizationOwner string
	Doctor            string
	Admin             string
}{
	SuperAdmin:        "super_admin",
	Staff:             "staff",
	Patient:           "patient",
	OrganizationOwner: "organization_owner",
	Doctor:            "doctor",
	Admin:             "admin",
}

// DefaultUserRole defines the default role for new users
const DefaultUserRole = "patient"
