package services

import (
	"errors"
	"fmt"

	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	userRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories"
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

// ListUsers returns a paginated list of users with optional filters (tenant, role, search).
func (s *UserService) ListUsers(tenantID *uint, page, pageSize int, role, search string) ([]models.User, int64, error) {
	return s.userRepo.List(tenantID, page, pageSize, role, search)
}

// ListSpecialRoleUsers returns users who have role admin, hr, or it (for transfer-role "from" picker, etc.).
func (s *UserService) ListSpecialRoleUsers(tenantID *uint, page, pageSize int, search string) ([]models.User, int64, error) {
	roles := []string{string(models.RoleAdmin), string(models.RoleHR), string(models.RoleIT)}
	return s.userRepo.ListByRoles(tenantID, page, pageSize, roles, search)
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
	User            *models.User `json:"user"`              // User who was assigned the role
	Role            string       `json:"role"`              // Role that was assigned (admin, hr, or it)
	PositionUpdated bool         `json:"position_updated"`  // True if employee's position was updated (optional position_id was provided and user has employee record)
}

// AssignRole assigns or reassigns a role (admin, hr, it, or employee) to any user.
// Only an Admin can call this. The target user must exist and be active. Any current role can be reassigned.
// For role "employee", legacy user_type is set to "user" and RBAC gets the "employee" role.
// If positionID is not nil, the linked employee record (if any) has its position_id updated.
func (s *UserService) AssignRole(callerUserID uint, targetUserID uint, role string, positionID *uint) (*AssignRoleResult, error) {
	allowedRoles := map[string]bool{"admin": true, "hr": true, "it": true, "employee": true}
	if !allowedRoles[role] {
		return nil, errors.New("role must be admin, hr, it, or employee")
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

	// Legacy user_type: "employee" maps to "user"; admin/hr/it stay as-is
	legacyRole := role
	if role == "employee" {
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
	if err := s.userRepo.Update(target); err != nil {
		return nil, fmt.Errorf("failed to %s user: %w", action, err)
	}
	return target, nil
}
