package models

// CreateUserRequest is the request body for creating a new system user linked to an existing employee.
// Admin or HR can call this to create login accounts for employees with a specific role (IT, HR, Manager, etc.).
// The user's name and email are taken from the employee record automatically.
type CreateUserRequest struct {
	EmployeeID uint   `json:"employee_id" binding:"required"`                                    // ID of the existing employee to create a user for
	Password   string `json:"password" binding:"required,min=6"`                                 // Login password (minimum 6 characters)
	Role       string `json:"role" binding:"required,oneof=admin hr it manager employee user"`   // Role to assign: admin, hr, it, manager, employee, or user
	Username   string `json:"username,omitempty" binding:"omitempty,min=3"`                       // Optional username (auto-generated from employee email if not provided)
}

// TransferRoleRequest is the request body for transferring a special role (admin, hr, it, manager) to another user.
// When the caller is Admin: from_user_id is required (the user who currently has the role).
// When the caller is transferring their own role: from_user_id is ignored; the caller is the source.
type TransferRoleRequest struct {
	Role         string `json:"role" binding:"required,oneof=admin hr it manager"` // Role to transfer: admin, hr, it, or manager
	TargetUserID uint   `json:"target_user_id" binding:"required"`                 // User who will receive the role (typically a normal employee)
	FromUserID   *uint  `json:"from_user_id,omitempty"`                            // Required when caller is Admin: the user who currently has the role
}

// AssignRoleRequest is the request body for assigning a role to a user.
// Only an Admin can call this. The target user must exist and be active.
// Use role "employee" to set/confirm user as normal employee (legacy user_type stays "user", RBAC gets "employee").
// Use role "manager" to grant manager access with team approval permissions.
// Optionally set position_id to update the user's employee record position (job position) at the same time.
type AssignRoleRequest struct {
	UserID     uint   `json:"user_id" binding:"required"`                                           // User to assign role to
	Role       string `json:"role" binding:"required,oneof=super_admin admin hr it manager employee"` // Role to assign
	PositionID *uint  `json:"position_id,omitempty"`                                                // Optional: set employee's job position
}

// AddRoleRequest is the request body for adding an additional role to a user.
// This allows users to have multiple roles like ["admin", "employee"] or ["hr", "manager"].
// Only an Admin can call this.
type AddRoleRequest struct {
	UserID uint   `json:"user_id" binding:"required"`                                        // User to add role to
	Role   string `json:"role" binding:"required,oneof=super_admin admin hr it manager employee user"` // Role to add
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

// ChangeRolesRequest is the request body for the checkbox-style "Change Roles" UI.
// Send the full desired set of roles (checked items). The backend computes what to add/remove.
// At least one role must be checked. Only an Admin can call this.
type ChangeRolesRequest struct {
	UserID uint     `json:"user_id" binding:"required"` // Target user ID
	Roles  []string `json:"roles" binding:"required"`   // Desired roles (checked items), e.g. ["employee", "manager", "hr"]
}

// ResetPasswordRequest is the request body for resetting a user's password.
// Admin, Super Admin, or HR can call this.
// If password is omitted, the default "GreenTelecom@2026" is used.
type ResetPasswordRequest struct {
	UserID   uint    `json:"user_id" binding:"required"`              // User whose password to reset
	Password *string `json:"password,omitempty" binding:"omitempty,min=6"` // New password (optional, default: GreenTelecom@2026)
}
