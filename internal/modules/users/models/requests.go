package models

// TransferRoleRequest is the request body for transferring a special role (admin, hr, it) to another user.
// When the caller is Admin: from_user_id is required (the user who currently has the role).
// When the caller is transferring their own role: from_user_id is ignored; the caller is the source.
type TransferRoleRequest struct {
	Role         string `json:"role" binding:"required,oneof=admin hr it"` // Role to transfer: admin, hr, or it
	TargetUserID uint   `json:"target_user_id" binding:"required"`         // User who will receive the role (typically a normal employee)
	FromUserID   *uint  `json:"from_user_id,omitempty"`                    // Required when caller is Admin: the user who currently has the role
}

// AssignRoleRequest is the request body for assigning a role (admin, hr, it, or employee) to a user.
// Only an Admin can call this. The target user must currently have role "user" (normal employee).
// Use role "employee" to set/confirm user as normal employee (legacy user_type stays "user", RBAC gets "employee").
// Optionally set position_id to update the user's employee record position (job position) at the same time.
type AssignRoleRequest struct {
	UserID     uint   `json:"user_id" binding:"required"`                         // User to assign role to (must have role "user")
	Role       string `json:"role" binding:"required,oneof=admin hr it employee"` // Role to assign: admin, hr, it, or employee
	PositionID *uint  `json:"position_id,omitempty"`                              // Optional: set employee's job position (only if user has an employee record)
}

// AddRoleRequest is the request body for adding an additional role to a user.
// This allows users to have multiple roles like ["admin", "employee"] or ["hr", "it"].
// Only an Admin can call this.
type AddRoleRequest struct {
	UserID uint   `json:"user_id" binding:"required"`                              // User to add role to
	Role   string `json:"role" binding:"required,oneof=admin hr it employee user"` // Role to add
}

// RemoveRoleRequest is the request body for removing a role from a user.
// Cannot remove the last role - user must have at least one role.
// Only an Admin can call this.
type RemoveRoleRequest struct {
	UserID uint   `json:"user_id" binding:"required"` // User to remove role from
	Role   string `json:"role" binding:"required"`    // Role to remove
}

// SetRolesRequest is the request body for setting/replacing all roles for a user.
// At least one role must be specified. Only an Admin can call this.
type SetRolesRequest struct {
	UserID uint     `json:"user_id" binding:"required"` // User to set roles for
	Roles  []string `json:"roles" binding:"required"`   // New roles array (e.g. ["admin", "employee"])
}
