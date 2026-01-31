package services

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/repositories"
)

type RoleService struct {
	roleRepo *repositories.RoleRepository
}

func NewRoleService() *RoleService {
	return &RoleService{
		roleRepo: repositories.NewRoleRepository(),
	}
}

// ListRoles returns a paginated list of roles (RBAC roles).
func (s *RoleService) ListRoles(tenantID *uint, page, pageSize int) ([]models.Role, int64, error) {
	return s.roleRepo.List(tenantID, page, pageSize)
}
