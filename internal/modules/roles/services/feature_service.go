package services

import (
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/repositories"
)

// FeatureService handles feature access business logic
type FeatureService struct {
	roleRepo *repositories.RoleRepository
	permRepo *repositories.PermissionRepository
}

// NewFeatureService creates a new feature service
func NewFeatureService() *FeatureService {
	return &FeatureService{
		roleRepo: repositories.NewRoleRepository(),
		permRepo: repositories.NewPermissionRepository(),
	}
}

// GetAllFeatures returns the complete feature catalog
func (s *FeatureService) GetAllFeatures() []models.Feature {
	return models.GetAllFeatures()
}

// GetMyFeatureAccess returns the feature access map for a given user.
// It combines permissions from all roles assigned to the user.
func (s *FeatureService) GetMyFeatureAccess(userID uint) ([]models.FeatureAccess, error) {
	// Get all roles for the user (with permissions preloaded)
	userRoles, err := s.roleRepo.GetUserRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Build a set of granted permission codes
	grantedPerms := make(map[string]bool)
	isAdmin := false
	for _, role := range userRoles {
		if role.Code == "admin" {
			isAdmin = true
		}
		for _, perm := range role.Permissions {
			grantedPerms[perm.Code] = true
		}
	}

	// Build feature access list
	features := models.GetAllFeatures()
	result := make([]models.FeatureAccess, 0, len(features))

	for _, feature := range features {
		fa := models.FeatureAccess{
			Code:        feature.Code,
			Name:        feature.Name,
			Description: feature.Description,
			Icon:        feature.Icon,
			Category:    feature.Category,
			HasAccess:   false,
			Permissions: make([]models.PermissionAccess, 0, len(feature.Permissions)),
		}

		for _, perm := range feature.Permissions {
			granted := isAdmin || grantedPerms[perm.Code]
			pa := models.PermissionAccess{
				Code:        perm.Code,
				Name:        perm.Name,
				Description: perm.Description,
				Granted:     granted,
			}
			fa.Permissions = append(fa.Permissions, pa)

			// Feature has access if at least one permission is granted
			if granted {
				fa.HasAccess = true
			}
		}

		result = append(result, fa)
	}

	return result, nil
}

// GetRoleFeatureAccess returns the feature access map for a specific role
func (s *FeatureService) GetRoleFeatureAccess(roleID uint) ([]models.FeatureAccess, error) {
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}

	// Build a set of granted permission codes for this role
	grantedPerms := make(map[string]bool)
	isAdmin := role.Code == "admin"
	for _, perm := range role.Permissions {
		grantedPerms[perm.Code] = true
	}

	// Build feature access list
	features := models.GetAllFeatures()
	result := make([]models.FeatureAccess, 0, len(features))

	for _, feature := range features {
		fa := models.FeatureAccess{
			Code:        feature.Code,
			Name:        feature.Name,
			Description: feature.Description,
			Icon:        feature.Icon,
			Category:    feature.Category,
			HasAccess:   false,
			Permissions: make([]models.PermissionAccess, 0, len(feature.Permissions)),
		}

		for _, perm := range feature.Permissions {
			granted := isAdmin || grantedPerms[perm.Code]
			pa := models.PermissionAccess{
				Code:        perm.Code,
				Name:        perm.Name,
				Description: perm.Description,
				Granted:     granted,
			}
			fa.Permissions = append(fa.Permissions, pa)

			if granted {
				fa.HasAccess = true
			}
		}

		result = append(result, fa)
	}

	return result, nil
}

// UpdateRolePermissions updates the permissions assigned to a role
func (s *FeatureService) UpdateRolePermissions(roleID uint, permissionCodes []string) error {
	// Verify the role exists
	_, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	// Validate all permission codes exist
	var permissionIDs []uint
	for _, code := range permissionCodes {
		perm, err := s.permRepo.FindByCode(code)
		if err != nil {
			return fmt.Errorf("permission not found: %s", code)
		}
		permissionIDs = append(permissionIDs, perm.ID)
	}

	// Replace permissions for the role
	return s.roleRepo.AssignPermissions(roleID, permissionIDs)
}

// GetRoleByCode returns a role by its code
func (s *FeatureService) GetRoleByCode(code string) (*models.Role, error) {
	return s.roleRepo.FindByCode(code)
}

// GetRoleByID returns a role by its ID
func (s *FeatureService) GetRoleByID(id uint) (*models.Role, error) {
	return s.roleRepo.FindByID(id)
}
