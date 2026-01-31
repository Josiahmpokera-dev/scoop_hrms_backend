package seed

import (
	"log"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/models"
	roleRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/repositories"
)

// RunRoles seeds default roles and assigns permissions.
func RunRoles() {
	roleRepo := roleRepos.NewRoleRepository()
	permRepo := roleRepos.NewPermissionRepository()

	defaultRoles := []struct {
		Code        string
		Name        string
		Description string
		Permissions []string
	}{
		{
			Code:        "admin",
			Name:        "Administrator",
			Description: "Full system access with all permissions",
			Permissions: []string{
				"employee:read", "employee:create", "employee:update", "employee:delete",
				"department:read", "department:create", "department:update", "department:delete",
				"payroll:run", "attendance:approve",
				"user:read", "user:create", "user:update", "user:delete",
			},
		},
		{
			Code:        "hr",
			Name:        "HR Manager",
			Description: "HR management with employee and department access",
			Permissions: []string{
				"employee:read", "employee:create", "employee:update",
				"department:read", "payroll:run", "attendance:approve",
				"user:read",
			},
		},
		{
			Code:        "it",
			Name:        "IT",
			Description: "IT support access (e.g. employee read, helpdesk)",
			Permissions: []string{
				"employee:read",
				"user:read",
			},
		},
		{
			Code:        "employee",
			Name:        "Employee",
			Description: "Basic employee access",
			Permissions: []string{
				"employee:read",
			},
		},
	}

	for _, roleData := range defaultRoles {
		existingRole, err := roleRepo.FindByCode(roleData.Code)
		if err == nil && existingRole != nil {
			continue
		}

		role := &models.Role{
			Code:        roleData.Code,
			Name:        roleData.Name,
			Description: &roleData.Description,
		}

		if err := roleRepo.Create(role); err != nil {
			log.Printf("Warning: Failed to create role %s: %v", roleData.Code, err)
			continue
		}

		var permissionIDs []uint
		for _, permCode := range roleData.Permissions {
			perm, err := permRepo.FindByCode(permCode)
			if err == nil && perm != nil {
				permissionIDs = append(permissionIDs, perm.ID)
			}
		}

		if len(permissionIDs) > 0 {
			if err := roleRepo.AssignPermissions(role.ID, permissionIDs); err != nil {
				log.Printf("Warning: Failed to assign permissions to role %s: %v", roleData.Code, err)
			}
		}

		log.Printf("✅ Role created: %s", roleData.Code)
	}

	log.Println("✅ Default roles seeded successfully")
}
