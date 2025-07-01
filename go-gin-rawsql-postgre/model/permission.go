package model

// Role constants
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// Permission constants
const (
	// Film permissions
	PermissionFilmRead   = "film:read"
	PermissionFilmCreate = "film:create"
	PermissionFilmUpdate = "film:update"
	PermissionFilmDelete = "film:delete"

	// Staff permissions
	PermissionStaffRead   = "staff:read"
	PermissionStaffCreate = "staff:create"
	PermissionStaffUpdate = "staff:update"
	PermissionStaffDelete = "staff:delete"

	// User management permissions
	PermissionUserRead   = "user:read"
	PermissionUserCreate = "user:create"
	PermissionUserUpdate = "user:update"
	PermissionUserDelete = "user:delete"
)

// RolePermissions maps roles to their allowed permissions
var RolePermissions = map[string][]string{
	RoleAdmin: {
		// Admin has all permissions
		PermissionFilmRead, PermissionFilmCreate, PermissionFilmUpdate, PermissionFilmDelete,
		PermissionStaffRead, PermissionStaffCreate, PermissionStaffUpdate, PermissionStaffDelete,
		PermissionUserRead, PermissionUserCreate, PermissionUserUpdate, PermissionUserDelete,
	},
	RoleUser: {
		// User has limited permissions
		PermissionFilmRead,
		PermissionStaffRead,
		PermissionUserRead,
	},
}

// HasPermission checks if a role has a specific permission
func HasPermission(role, permission string) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// GetRolePermissions returns all permissions for a given role
func GetRolePermissions(role string) []string {
	permissions, exists := RolePermissions[role]
	if !exists {
		return []string{}
	}
	return permissions
}

// IsValidRole checks if a role is valid
func IsValidRole(role string) bool {
	_, exists := RolePermissions[role]
	return exists
}
