package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	userRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles user-related business logic.
type UserService struct {
	userRepo     *userRepos.UserRepository
	roleRepo     *repositories.RoleRepository
	employeeRepo *employeeRepos.EmployeeRepository
}

// NewUserService creates a new user service.
func NewUserService() *UserService {
	return &UserService{
		userRepo:     userRepos.NewUserRepository(),
		roleRepo:     repositories.NewRoleRepository(),
		employeeRepo: employeeRepos.NewEmployeeRepository(),
	}
}

// CreateUserResult holds the result of creating a new user linked to an employee.
type CreateUserResult struct {
	User       *models.User `json:"user"`
	EmployeeID string       `json:"employee_id"` // Employee number (e.g. EMP001)
	Roles      []string     `json:"roles"`
	Email      string       `json:"email"`    // Login email (returned so admin can share credentials)
	Password   string       `json:"password"` // Plain password (returned once at creation only)
}

// CreateUser creates a new system user account linked to an existing employee record.
// The employee must exist and must not already have a user account.
// User data (name, email, phone) is pulled from the employee record automatically.
// Only Admin or HR can call this.
func (s *UserService) CreateUser(callerUserID uint, req *models.CreateUserRequest) (*CreateUserResult, error) {
	// Verify caller is admin or HR
	caller, err := s.userRepo.FindByID(callerUserID)
	if err != nil || caller == nil {
		return nil, errors.New("caller user not found")
	}
	if !caller.IsActive {
		return nil, errors.New("caller account is not active")
	}
	if !caller.IsAdmin() && !caller.IsHR() {
		return nil, errors.New("only an admin or HR can create users")
	}

	// Find the employee
	employee, err := s.employeeRepo.FindByID(req.EmployeeID)
	if err != nil || employee == nil {
		return nil, errors.New("employee not found")
	}

	// Check if employee is active
	if !employee.IsActiveStatus() {
		return nil, errors.New("employee is not active")
	}

	// Check if employee already has a user account
	if employee.UserID != nil && *employee.UserID > 0 {
		return nil, errors.New("this employee already has a user account")
	}

	// Determine email: use work_email first, then personal_email
	email := ""
	if employee.WorkEmail != nil && *employee.WorkEmail != "" {
		email = *employee.WorkEmail
	} else if employee.PersonalEmail != nil && *employee.PersonalEmail != "" {
		email = *employee.PersonalEmail
	}
	if email == "" {
		return nil, errors.New("employee does not have an email address (work_email or personal_email required)")
	}

	// Check if email is already used by another user
	if s.userRepo.ExistsByEmail(email) {
		return nil, fmt.Errorf("a user with email %s already exists", email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Generate username from email if not provided
	username := req.Username
	if username == "" {
		emailParts := strings.Split(email, "@")
		if len(emailParts) > 0 {
			username = emailParts[0]
		} else {
			username = strings.ToLower(employee.FirstName + employee.LastName)
		}
		username = strings.ToLower(strings.ReplaceAll(username, ".", ""))
		username = strings.ReplaceAll(username, "_", "")
		username = strings.ReplaceAll(username, "-", "")
	}

	// Determine legacy role (user_type column)
	role := req.Role
	legacyRole := role
	if role == "employee" || role == "manager" {
		legacyRole = "user"
	}

	// Create user from employee data
	user := &models.User{
		Username:    username,
		Email:       email,
		Password:    string(hashedPassword),
		FirstName:   employee.FirstName,
		LastName:    employee.LastName,
		PhoneNumber: employee.PhoneNumber,
		Role:        models.UserRole(legacyRole),
		Status:      models.UserStatusActive,
		IsActive:    true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Link employee to the new user
	employee.UserID = &user.ID
	if err := s.employeeRepo.Update(employee); err != nil {
		// Rollback: delete the created user if linking fails
		_ = s.userRepo.Delete(user.ID)
		return nil, errors.New("failed to link employee to user")
	}

	// Assign the RBAC role so permissions work correctly
	rbacRole, err := s.roleRepo.FindByCode(role)
	if err == nil && rbacRole != nil {
		_ = s.roleRepo.AssignUserRole(user.ID, rbacRole.ID, &callerUserID)
	}

	// Get the final roles list
	roles, _ := s.GetUserRoles(user.ID)

	return &CreateUserResult{
		User:       user,
		EmployeeID: employee.EmployeeID,
		Roles:      roles,
		Email:      email,
		Password:   req.Password,
	}, nil
}

// ListEmployeesWithoutUser returns employees that don't have a linked user account yet.
// This is used by the frontend to populate the employee picker when creating a new user.
func (s *UserService) ListEmployeesWithoutUser(page, pageSize int, search string) ([]map[string]interface{}, int64, error) {
	employees, total, err := s.employeeRepo.ListEmployeesWithoutUser(page, pageSize, search)
	if err != nil {
		return nil, 0, err
	}

	// Build a simplified list for the picker
	var results []map[string]interface{}
	for _, emp := range employees {
		email := ""
		if emp.WorkEmail != nil && *emp.WorkEmail != "" {
			email = *emp.WorkEmail
		} else if emp.PersonalEmail != nil && *emp.PersonalEmail != "" {
			email = *emp.PersonalEmail
		}

		item := map[string]interface{}{
			"id":          emp.ID,
			"employee_id": emp.EmployeeID,
			"first_name":  emp.FirstName,
			"last_name":   emp.LastName,
			"email":       email,
			"department_id": emp.DepartmentID,
			"position_id":   emp.PositionID,
			"status":        emp.Status,
		}
		if emp.PhoneNumber != nil {
			item["phone_number"] = *emp.PhoneNumber
		}
		if emp.PhotoURL != nil {
			item["photo_url"] = *emp.PhotoURL
		}
		results = append(results, item)
	}

	return results, total, nil
}

// ListUsers returns a paginated list of users with optional filters (tenant, role, search).
func (s *UserService) ListUsers(tenantID *uint, page, pageSize int, role, search string) ([]models.User, int64, error) {
	return s.userRepo.List(tenantID, page, pageSize, role, search)
}

// ListSpecialRoleUsers returns users who have non-basic roles (admin, hr, it, manager) for transfer-role "from" picker, etc.
func (s *UserService) ListSpecialRoleUsers(tenantID *uint, page, pageSize int, search string) ([]models.User, int64, error) {
	roles := []string{string(models.RoleSuperAdmin), string(models.RoleAdmin), string(models.RoleHR), string(models.RoleIT), string(models.RoleManager)}
	return s.userRepo.ListByRoles(tenantID, page, pageSize, roles, search)
}

// ListBlockedUsers returns users blocked from login (rate-limit or manual block). Optionally includes suspended.
func (s *UserService) ListBlockedUsers(tenantID *uint, page, pageSize int, includeSuspended bool) ([]models.User, int64, error) {
	return s.userRepo.ListBlockedUsers(tenantID, page, pageSize, includeSuspended)
}

// RateLimitStatus holds login rate-limit state for a user (admin view).
type RateLimitStatus struct {
	Email             string     `json:"email"`
	Blocked           bool       `json:"blocked"`
	Suspended         bool       `json:"suspended"`
	FailedLoginCount  int        `json:"failedLoginCount"`
	MaxAttempts       int        `json:"maxAttempts"`
	AttemptsRemaining int        `json:"attemptsRemaining"`
	LockedUntil       *time.Time `json:"lockedUntil,omitempty"`
	CanLogin          bool       `json:"canLogin"`
}

// GetRateLimitStatus returns rate-limit status for an email (admin only). If user not found, returns generic to avoid enumeration.
func (s *UserService) GetRateLimitStatus(email string) (*RateLimitStatus, error) {
	maxAttempts := 5
	if config.AppConfig != nil && config.AppConfig.LoginRateLimit.MaxAttempts > 0 {
		maxAttempts = config.AppConfig.LoginRateLimit.MaxAttempts
	}
	user, err := s.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return &RateLimitStatus{
			Email:             email,
			MaxAttempts:       maxAttempts,
			AttemptsRemaining: maxAttempts,
			CanLogin:          true,
		}, nil
	}
	attemptsRemaining := maxAttempts - user.FailedLoginCount
	if attemptsRemaining < 0 {
		attemptsRemaining = 0
	}
	return &RateLimitStatus{
		Email:             user.Email,
		Blocked:           user.Status == models.UserStatusBlocked,
		Suspended:         user.Status == models.UserStatusSuspended,
		FailedLoginCount:  user.FailedLoginCount,
		MaxAttempts:       maxAttempts,
		AttemptsRemaining: attemptsRemaining,
		LockedUntil:       user.LockedUntil,
		CanLogin:          user.CanLogin(),
	}, nil
}

// TransferRoleResult holds the result of a role transfer.
type TransferRoleResult struct {
	FromUser   *models.User `json:"from_user"`   // User who lost the role (now has role "user")
	TargetUser *models.User `json:"target_user"` // User who received the role
	Role       string       `json:"role"`        // Role that was transferred
}

// TransferRole transfers a special role (admin, hr, it) from one user to another.
// Caller must be Admin (can transfer any role, from_user_id required) or the current holder of the role (transfers their own role).
func (s *UserService) TransferRole(callerUserID uint, role string, targetUserID uint, fromUserID *uint) (*TransferRoleResult, error) {
	roleType := models.UserRole(role)
	if roleType != models.RoleAdmin && roleType != models.RoleHR && roleType != models.RoleIT {
		return nil, errors.New("role must be admin, hr, or it")
	}

	caller, err := s.userRepo.FindByID(callerUserID)
	if err != nil || caller == nil {
		return nil, errors.New("caller user not found")
	}
	if !caller.IsActive {
		return nil, errors.New("caller account is not active")
	}

	target, err := s.userRepo.FindByID(targetUserID)
	if err != nil || target == nil {
		return nil, errors.New("target user not found")
	}
	if !target.IsActive {
		return nil, errors.New("target user is not active")
	}
	// Target must be a normal employee (role "user") to receive the transferred role
	if target.Role != models.RoleUser {
		return nil, errors.New("role can only be transferred to a normal employee (user with role 'user')")
	}

	var from *models.User
	if caller.IsAdmin() {
		// Admin must specify who to take the role from
		if fromUserID == nil {
			return nil, errors.New("from_user_id is required when transferring as admin")
		}
		from, err = s.userRepo.FindByID(*fromUserID)
		if err != nil || from == nil {
			return nil, errors.New("from user not found")
		}
	} else {
		// Caller is transferring their own role
		from = caller
		if from.Role != roleType {
			return nil, fmt.Errorf("you do not have the role %s to transfer", role)
		}
	}

	if from.ID == target.ID {
		return nil, errors.New("cannot transfer role to the same user")
	}
	if from.Role != roleType {
		return nil, fmt.Errorf("from user does not have role %s", role)
	}

	// Transfer: from becomes "user", target gets the role
	fromOriginalRole := from.Role
	from.Role = models.RoleUser
	target.Role = roleType

	if err := s.userRepo.Update(from); err != nil {
		return nil, fmt.Errorf("failed to update from user: %w", err)
	}
	if err := s.userRepo.Update(target); err != nil {
		// Rollback: restore from user's role
		from.Role = fromOriginalRole
		_ = s.userRepo.Update(from)
		return nil, fmt.Errorf("failed to update target user: %w", err)
	}

	return &TransferRoleResult{
		FromUser:   from,
		TargetUser: target,
		Role:       role,
	}, nil
}

// AssignRoleResult holds the result of assigning a special role to a user.
type AssignRoleResult struct {
	User            *models.User `json:"user"`             // User who was assigned the role
	Role            string       `json:"role"`             // Role that was assigned (admin, hr, or it)
	PositionUpdated bool         `json:"position_updated"` // True if employee's position was updated (optional position_id was provided and user has employee record)
}

// AssignRole assigns or reassigns a role to any user.
// Only an Admin can call this. The target user must exist and be active. Any current role can be reassigned.
// For roles "employee" and "manager", legacy user_type is set to "user" and RBAC gets the respective role.
// If positionID is not nil, the linked employee record (if any) has its position_id updated.
func (s *UserService) AssignRole(callerUserID uint, targetUserID uint, role string, positionID *uint) (*AssignRoleResult, error) {
	allowedRoles := map[string]bool{"super_admin": true, "admin": true, "hr": true, "it": true, "manager": true, "employee": true}
	if !allowedRoles[role] {
		return nil, errors.New("role must be super_admin, admin, hr, it, manager, or employee")
	}

	caller, err := s.userRepo.FindByID(callerUserID)
	if err != nil || caller == nil {
		return nil, errors.New("caller user not found")
	}
	if !caller.IsActive {
		return nil, errors.New("caller account is not active")
	}
	if !caller.IsAdmin() {
		return nil, errors.New("only an admin can assign roles to a user")
	}

	target, err := s.userRepo.FindByID(targetUserID)
	if err != nil || target == nil {
		return nil, errors.New("target user not found")
	}
	if !target.IsActive {
		return nil, errors.New("target user is not active")
	}

	// Legacy user_type: "employee" and "manager" map to "user"; admin/hr/it stay as-is
	legacyRole := role
	if role == "employee" || role == "manager" {
		legacyRole = string(models.RoleUser)
	}
	target.Role = models.UserRole(legacyRole)
	if err := s.userRepo.Update(target); err != nil {
		return nil, fmt.Errorf("failed to update user role: %w", err)
	}

	// Sync RBAC: ensure user_roles has the corresponding role so middleware/auth see it
	rbacRole, err := s.roleRepo.FindByCode(role)
	if err == nil && rbacRole != nil {
		codes, _ := s.roleRepo.GetUserRoleCodes(target.ID)
		hasRole := false
		for _, c := range codes {
			if c == role {
				hasRole = true
				break
			}
		}
		if !hasRole {
			_ = s.roleRepo.AssignUserRole(target.ID, rbacRole.ID, &callerUserID)
		}
	}

	positionUpdated := false
	if positionID != nil {
		emp, err := s.employeeRepo.FindByUserID(target.ID)
		if err == nil && emp != nil {
			emp.PositionID = positionID
			if err := s.employeeRepo.Update(emp); err == nil {
				positionUpdated = true
			}
		}
	}

	return &AssignRoleResult{
		User:            target,
		Role:            role,
		PositionUpdated: positionUpdated,
	}, nil
}

// SuspendUser sets the user's status to suspended; they cannot login until unblocked.
// Caller must be HR or Admin. Cannot suspend yourself.
func (s *UserService) SuspendUser(callerUserID, targetUserID uint) (*models.User, error) {
	return s.setUserStatus(callerUserID, targetUserID, models.UserStatusSuspended, "suspend")
}

// BlockUser sets the user's status to blocked; they cannot login until unblocked.
// Caller must be HR or Admin. Cannot block yourself.
func (s *UserService) BlockUser(callerUserID, targetUserID uint) (*models.User, error) {
	return s.setUserStatus(callerUserID, targetUserID, models.UserStatusBlocked, "block")
}

// UnblockUser sets the user's status to active (restores login). Use for both suspended and blocked users.
// Caller must be HR or Admin.
func (s *UserService) UnblockUser(callerUserID, targetUserID uint) (*models.User, error) {
	return s.setUserStatus(callerUserID, targetUserID, models.UserStatusActive, "unblock")
}

// UnsuspendUser sets the user's status to active (restores login). Use for suspended users; same effect as unblock.
// Caller must be HR or Admin.
func (s *UserService) UnsuspendUser(callerUserID, targetUserID uint) (*models.User, error) {
	return s.setUserStatus(callerUserID, targetUserID, models.UserStatusActive, "unsuspend")
}

func (s *UserService) setUserStatus(callerUserID, targetUserID uint, status, action string) (*models.User, error) {
	if callerUserID == targetUserID && action != "unblock" {
		return nil, errors.New("you cannot " + action + " yourself")
	}
	target, err := s.userRepo.FindByID(targetUserID)
	if err != nil || target == nil {
		return nil, errors.New("user not found")
	}
	target.Status = status
	if status == models.UserStatusActive {
		target.FailedLoginCount = 0
		target.LockedUntil = nil
	}
	if err := s.userRepo.Update(target); err != nil {
		return nil, fmt.Errorf("failed to %s user: %w", action, err)
	}
	return target, nil
}

// AddRoleToUser adds an additional role to a user without changing existing roles.
// This allows a user to have multiple roles like ["admin", "employee"] or ["hr", "manager"].
// Only Admin can call this.
func (s *UserService) AddRoleToUser(callerUserID uint, targetUserID uint, role string) ([]string, error) {
	allowedRoles := map[string]bool{"super_admin": true, "admin": true, "hr": true, "it": true, "manager": true, "employee": true, "user": true}
	role = strings.ToLower(strings.TrimSpace(role))
	if !allowedRoles[role] {
		return nil, errors.New("role must be super_admin, admin, hr, it, manager, employee, or user")
	}

	caller, err := s.userRepo.FindByID(callerUserID)
	if err != nil || caller == nil {
		return nil, errors.New("caller user not found")
	}
	if !caller.IsAdmin() {
		return nil, errors.New("only an admin can add roles to users")
	}

	target, err := s.userRepo.FindByID(targetUserID)
	if err != nil || target == nil {
		return nil, errors.New("target user not found")
	}

	// Check if user already has this role in RBAC
	codes, _ := s.roleRepo.GetUserRoleCodes(target.ID)
	for _, c := range codes {
		if strings.ToLower(c) == role {
			return nil, fmt.Errorf("user already has role %s", role)
		}
	}

	// Find the role in RBAC
	rbacRole, err := s.roleRepo.FindByCode(role)
	if err != nil || rbacRole == nil {
		return nil, fmt.Errorf("role %s not found in system", role)
	}

	// Add the role
	if err := s.roleRepo.AssignUserRole(target.ID, rbacRole.ID, &callerUserID); err != nil {
		return nil, fmt.Errorf("failed to add role: %w", err)
	}

	// Return all roles for the user
	return s.GetUserRoles(targetUserID)
}

// RemoveRoleFromUser removes a role from a user.
// Cannot remove the last role - user must have at least one role.
// Only Admin can call this.
func (s *UserService) RemoveRoleFromUser(callerUserID uint, targetUserID uint, role string) ([]string, error) {
	role = strings.ToLower(strings.TrimSpace(role))

	caller, err := s.userRepo.FindByID(callerUserID)
	if err != nil || caller == nil {
		return nil, errors.New("caller user not found")
	}
	if !caller.IsAdmin() {
		return nil, errors.New("only an admin can remove roles from users")
	}

	target, err := s.userRepo.FindByID(targetUserID)
	if err != nil || target == nil {
		return nil, errors.New("target user not found")
	}

	// Get current roles
	codes, _ := s.roleRepo.GetUserRoleCodes(target.ID)

	// Count total roles (legacy + RBAC)
	totalRoles := len(codes)
	if strings.TrimSpace(string(target.Role)) != "" {
		// Check if legacy role is in codes to avoid double counting
		legacyInCodes := false
		for _, c := range codes {
			if strings.EqualFold(c, string(target.Role)) {
				legacyInCodes = true
				break
			}
		}
		if !legacyInCodes {
			totalRoles++
		}
	}

	if totalRoles <= 1 {
		return nil, errors.New("cannot remove the last role - user must have at least one role")
	}

	// Find and remove the role from RBAC
	rbacRole, err := s.roleRepo.FindByCode(role)
	if err != nil || rbacRole == nil {
		return nil, fmt.Errorf("role %s not found", role)
	}

	if err := s.roleRepo.RemoveUserRole(target.ID, rbacRole.ID); err != nil {
		return nil, fmt.Errorf("failed to remove role: %w", err)
	}

	// If removing the legacy role, change it to "user"
	if strings.ToLower(string(target.Role)) == role {
		target.Role = models.RoleUser
		_ = s.userRepo.Update(target)
	}

	// Return remaining roles
	return s.GetUserRoles(targetUserID)
}

// UserDetail holds user data plus roles for the "view user" API.
type UserDetail struct {
	ID            uint       `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	PhoneNumber   *string    `json:"phone_number,omitempty"`
	Role          string     `json:"role"`    // Legacy primary role
	Roles         []string   `json:"roles"`   // All RBAC roles
	Status        string     `json:"status"`  // active, suspended, blocked
	IsActive      bool       `json:"is_active"`
	EmailVerified bool       `json:"email_verified"`
	LastLogin     *time.Time `json:"last_login,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// GetUserDetail returns a user's full profile with their roles. HR or Admin can call this.
func (s *UserService) GetUserDetail(targetUserID uint) (*UserDetail, error) {
	user, err := s.userRepo.FindByID(targetUserID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	roles, _ := s.GetUserRoles(user.ID)

	return &UserDetail{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		PhoneNumber:   user.PhoneNumber,
		Role:          string(user.Role),
		Roles:         roles,
		Status:        user.Status,
		IsActive:      user.IsActive,
		EmailVerified: user.EmailVerified,
		LastLogin:     user.LastLogin,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}, nil
}

// GetUserRoles returns all roles for a user (legacy + RBAC combined, deduplicated).
func (s *UserService) GetUserRoles(userID uint) ([]string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	seen := make(map[string]bool)
	var roles []string

	// Add legacy role
	legacy := strings.ToLower(strings.TrimSpace(string(user.Role)))
	if legacy != "" && !seen[legacy] {
		seen[legacy] = true
		roles = append(roles, legacy)
	}

	// Add RBAC roles
	codes, _ := s.roleRepo.GetUserRoleCodes(user.ID)
	for _, code := range codes {
		c := strings.ToLower(strings.TrimSpace(code))
		if c != "" && !seen[c] {
			seen[c] = true
			roles = append(roles, c)
		}
	}

	// Always ensure "user" is in the roles array - all system users are treated as users
	if !seen["user"] {
		roles = append(roles, "user")
	}

	return roles, nil
}

// SetUserRoles replaces all roles for a user with the specified roles.
// At least one role must be provided. Only Admin can call this.
func (s *UserService) SetUserRoles(callerUserID uint, targetUserID uint, newRoles []string) ([]string, error) {
	if len(newRoles) == 0 {
		return nil, errors.New("at least one role must be specified")
	}

	allowedRoles := map[string]bool{"super_admin": true, "admin": true, "hr": true, "it": true, "manager": true, "employee": true, "user": true}
	for i, r := range newRoles {
		newRoles[i] = strings.ToLower(strings.TrimSpace(r))
		if !allowedRoles[newRoles[i]] {
			return nil, fmt.Errorf("invalid role: %s", r)
		}
	}

	caller, err := s.userRepo.FindByID(callerUserID)
	if err != nil || caller == nil {
		return nil, errors.New("caller user not found")
	}
	if !caller.IsAdmin() {
		return nil, errors.New("only an admin can set user roles")
	}

	target, err := s.userRepo.FindByID(targetUserID)
	if err != nil || target == nil {
		return nil, errors.New("target user not found")
	}

	// Remove all existing RBAC roles
	existingCodes, _ := s.roleRepo.GetUserRoleCodes(target.ID)
	for _, code := range existingCodes {
		rbacRole, err := s.roleRepo.FindByCode(code)
		if err == nil && rbacRole != nil {
			_ = s.roleRepo.RemoveUserRole(target.ID, rbacRole.ID)
		}
	}

	// Add new roles
	for _, role := range newRoles {
		rbacRole, err := s.roleRepo.FindByCode(role)
		if err == nil && rbacRole != nil {
			_ = s.roleRepo.AssignUserRole(target.ID, rbacRole.ID, &callerUserID)
		}
	}

	// Set legacy role to the first role (for backward compatibility)
	primaryRole := newRoles[0]
	if primaryRole == "employee" || primaryRole == "manager" {
		primaryRole = "user"
	}
	target.Role = models.UserRole(primaryRole)
	_ = s.userRepo.Update(target)

	return s.GetUserRoles(targetUserID)
}

// ChangeRolesResult holds the detailed result of a change-roles operation.
type ChangeRolesResult struct {
	UserID        uint     `json:"user_id"`
	Email         string   `json:"email"`
	Username      string   `json:"username"`
	PreviousRoles []string `json:"previous_roles"` // Roles before the change
	CurrentRoles  []string `json:"current_roles"`  // Roles after the change
	Added         []string `json:"added"`           // Roles that were added (newly checked)
	Removed       []string `json:"removed"`         // Roles that were removed (unchecked)
	Unchanged     []string `json:"unchanged"`       // Roles that stayed the same
}

// ChangeUserRoles applies a checkbox-style role change: the caller sends the full desired set of roles,
// and the service computes what needs to be added/removed. Returns the diff for the frontend.
// Only Admin or Super Admin can call this.
func (s *UserService) ChangeUserRoles(callerUserID uint, targetUserID uint, desiredRoles []string) (*ChangeRolesResult, error) {
	if len(desiredRoles) == 0 {
		return nil, errors.New("at least one role must be selected")
	}

	// Validate all desired roles
	allowedRoles := map[string]bool{"super_admin": true, "admin": true, "hr": true, "it": true, "manager": true, "employee": true, "user": true}
	for i, r := range desiredRoles {
		desiredRoles[i] = strings.ToLower(strings.TrimSpace(r))
		if !allowedRoles[desiredRoles[i]] {
			return nil, fmt.Errorf("invalid role: %s (allowed: super_admin, admin, hr, it, manager, employee, user)", r)
		}
	}

	// Verify caller is admin
	caller, err := s.userRepo.FindByID(callerUserID)
	if err != nil || caller == nil {
		return nil, errors.New("caller user not found")
	}
	if !caller.IsAdmin() {
		return nil, errors.New("only an admin can change user roles")
	}

	// Get target user
	target, err := s.userRepo.FindByID(targetUserID)
	if err != nil || target == nil {
		return nil, errors.New("target user not found")
	}

	// Get current RBAC roles
	previousRoles, err := s.GetUserRoles(targetUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current roles: %w", err)
	}

	// Build sets for comparison
	previousSet := make(map[string]bool)
	for _, r := range previousRoles {
		previousSet[r] = true
	}
	desiredSet := make(map[string]bool)
	for _, r := range desiredRoles {
		desiredSet[r] = true
	}
	// Always ensure "user" is in desired set
	desiredSet["user"] = true

	// Compute diff
	var added, removed, unchanged []string
	for r := range desiredSet {
		if previousSet[r] {
			unchanged = append(unchanged, r)
		} else {
			added = append(added, r)
		}
	}
	for r := range previousSet {
		if !desiredSet[r] {
			removed = append(removed, r)
		}
	}

	// Apply additions
	for _, role := range added {
		rbacRole, err := s.roleRepo.FindByCode(role)
		if err == nil && rbacRole != nil {
			_ = s.roleRepo.AssignUserRole(target.ID, rbacRole.ID, &callerUserID)
		}
	}

	// Apply removals
	for _, role := range removed {
		rbacRole, err := s.roleRepo.FindByCode(role)
		if err == nil && rbacRole != nil {
			_ = s.roleRepo.RemoveUserRole(target.ID, rbacRole.ID)
		}
	}

	// Update legacy user_type column to the highest-priority role
	legacyRole := determineLegacyRole(desiredSet)
	target.Role = models.UserRole(legacyRole)
	_ = s.userRepo.Update(target)

	// Get final roles
	currentRoles, _ := s.GetUserRoles(targetUserID)

	return &ChangeRolesResult{
		UserID:        target.ID,
		Email:         target.Email,
		Username:      target.Username,
		PreviousRoles: previousRoles,
		CurrentRoles:  currentRoles,
		Added:         added,
		Removed:       removed,
		Unchanged:     unchanged,
	}, nil
}

// determineLegacyRole picks the highest-priority role for the legacy user_type column.
// Priority: super_admin > admin > hr > it > user (employee and manager map to "user").
func determineLegacyRole(roles map[string]bool) string {
	if roles["super_admin"] {
		return "super_admin"
	}
	if roles["admin"] {
		return "admin"
	}
	if roles["hr"] {
		return "hr"
	}
	if roles["it"] {
		return "it"
	}
	return "user"
}

// DefaultResetPassword is used when no custom password is provided.
const DefaultResetPassword = "GreenTelecom@2026"

// ResetPasswordResult holds the result of a password reset.
type ResetPasswordResult struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

// ResetUserPassword resets a user's password to the provided value (or default).
// Also clears failed login count, unblocks, and reactivates the account.
// Only Admin, Super Admin, or HR can call this.
func (s *UserService) ResetUserPassword(callerUserID, targetUserID uint, newPassword string) (*ResetPasswordResult, error) {
	// Verify caller
	caller, err := s.userRepo.FindByID(callerUserID)
	if err != nil || caller == nil {
		return nil, errors.New("caller user not found")
	}
	if !caller.IsAdmin() && !caller.IsHR() {
		return nil, errors.New("only Admin, Super Admin, or HR can reset passwords")
	}

	// Prevent self-reset via this endpoint
	if callerUserID == targetUserID {
		return nil, errors.New("cannot reset your own password via this endpoint")
	}

	// Find target user
	target, err := s.userRepo.FindByID(targetUserID)
	if err != nil || target == nil {
		return nil, errors.New("target user not found")
	}

	// Use default password if none provided
	if newPassword == "" {
		newPassword = DefaultResetPassword
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Reset password and clear any locks
	target.Password = string(hashedPassword)
	target.FailedLoginCount = 0
	target.LockedUntil = nil
	target.Status = "active"
	target.IsActive = true

	if err := s.userRepo.Update(target); err != nil {
		return nil, fmt.Errorf("failed to reset password: %w", err)
	}

	return &ResetPasswordResult{
		UserID:   target.ID,
		Email:    target.Email,
		Username: target.Username,
		Message:  "Password reset successfully. Account unblocked and reactivated.",
	}, nil
}
