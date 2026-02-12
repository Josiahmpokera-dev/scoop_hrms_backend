package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole represents user roles
type UserRole string

const (
	RoleSuperAdmin UserRole = "super_admin"
	RoleAdmin      UserRole = "admin"
	RoleHR         UserRole = "hr"
	RoleIT         UserRole = "it"
	RoleManager    UserRole = "manager"
	RoleEmployee   UserRole = "employee"
	RoleUser       UserRole = "user" // Legacy base role — all users have this
)

// User status values (Status field). Blocked/suspended users cannot login.
const (
	UserStatusActive    = "active"
	UserStatusSuspended = "suspended"
	UserStatusBlocked   = "blocked"
)

// User represents a user in the system
type User struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Username         string         `json:"username" gorm:"not null;size:100"` // Username field
	Email            string         `json:"email" gorm:"uniqueIndex:idx_users_email;not null;size:255"`
	Password         string         `json:"-" gorm:"column:password_hash;type:varchar(255);not null"` // Hidden from JSON, maps to password_hash column
	FirstName        string         `json:"first_name" gorm:"not null;size:100"`
	LastName         string         `json:"last_name" gorm:"not null;size:100"`
	PhoneNumber      *string        `json:"phone_number,omitempty" gorm:"size:50"`
	Role             UserRole       `json:"role" gorm:"column:user_type;type:varchar(20);default:'user';not null"` // Maps to user_type column (legacy support)
	Status           string         `json:"status" gorm:"type:varchar(20);default:'active';not null"`              // active, suspended, blocked
	EmailVerified    bool           `json:"email_verified" gorm:"default:false"`
	IsActive         bool           `json:"is_active" gorm:"default:true"`
	LastLogin        *time.Time     `json:"last_login,omitempty" gorm:"column:last_login_at"`
	FailedLoginCount int            `json:"-" gorm:"default:0"`                                // Reset on success; used for rate limiting
	LockedUntil      *time.Time     `json:"locked_until,omitempty" gorm:"column:locked_until"` // Temporary lock (optional)
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	UpdatedBy        *uint          `json:"updated_by,omitempty" gorm:"index"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships (using interface{} to avoid circular imports)
	// Tenant and Roles will be loaded separately when needed
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// IsAdmin checks if user is an admin (includes super_admin)
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin || u.Role == RoleSuperAdmin
}

// IsSuperAdmin checks if user is a super admin
func (u *User) IsSuperAdmin() bool {
	return u.Role == RoleSuperAdmin
}

// IsHR checks if user is HR staff
func (u *User) IsHR() bool {
	return u.Role == RoleHR
}

// IsEmployee checks if user has the employee role
func (u *User) IsEmployee() bool {
	return u.Role == RoleEmployee || u.Role == RoleUser
}

// IsManager checks if user has the manager role
func (u *User) IsManager() bool {
	return u.Role == RoleManager
}

// FullName returns the full name of the user
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// CanLogin returns true if the user is allowed to login (active and not suspended/blocked).
func (u *User) CanLogin() bool {
	if !u.IsActive {
		return false
	}
	return u.Status == UserStatusActive
}
