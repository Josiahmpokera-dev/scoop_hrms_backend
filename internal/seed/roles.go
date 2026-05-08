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

	// Get all permission codes from the feature catalog
	allPerms := models.GetAllPermissionCodes()

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
			Permissions: allPerms, // Admin gets everything
		},
		{
			Code:        "ceo",
			Name:        "Chief Executive Officer",
			Description: "Executive role with final authority including offboarding final approvals",
			Permissions: allPerms,
		},
		{
			Code:        "managing_director",
			Name:        "Managing Director",
			Description: "Executive role with final authority including offboarding final approvals",
			Permissions: allPerms,
		},
		{
			Code:        "hr",
			Name:        "HR Manager",
			Description: "HR management with employee, leave, attendance, and organizational access",
			Permissions: []string{
				// Dashboard
				"dashboard:view", "dashboard:manage_announcements",
				// Employees (full)
				"employee:read", "employee:create", "employee:update", "employee:delete",
				"employee:onboard", "employee:offboard",
				// Leave Management (full)
				"leave:read", "leave:create", "leave:approve",
				"leave:manage_types", "leave:manage_policies", "leave:manage_holidays",
				// Attendance
				"attendance:read", "attendance:approve",
				// Shifts & Rosters
				"shift:read", "shift:create", "shift:update", "shift:delete",
				"shift:manage_rosters", "shift:approve_swaps",
				// Organization
				"department:read", "department:create", "department:update", "department:delete",
				"team:read", "team:create", "team:update", "team:delete",
				"position:read", "position:create", "position:update", "position:delete",
				"organization:read", "organization:update",
				"organization_unit:read", "organization_unit:create", "organization_unit:update", "organization_unit:delete",
				"location:read", "location:create", "location:update", "location:delete",
				"cost_center:read", "cost_center:create", "cost_center:update", "cost_center:delete",
				// Payroll
				"payroll:read", "payroll:run", "payroll:manage_structures",
				"payroll:manage_loans", "payroll:manage_reports",
				// Assets
				"asset:read", "asset:create", "asset:update", "asset:assign",
				// Helpdesk
				"helpdesk:read", "helpdesk:create", "helpdesk:manage", "helpdesk:manage_kb",
				// Users (read only)
				"user:read",
				// Reports
				"report:view", "report:export",
			},
		},
		{
			Code:        "it",
			Name:        "IT Support",
			Description: "IT support access with helpdesk and asset management",
			Permissions: []string{
				// Dashboard
				"dashboard:view",
				// Employees (read only)
				"employee:read",
				// Assets (full)
				"asset:read", "asset:create", "asset:update", "asset:assign",
				// Helpdesk (full)
				"helpdesk:read", "helpdesk:create", "helpdesk:manage", "helpdesk:manage_kb",
				// Users (read only)
				"user:read",
				// Attendance
				"attendance:read", "attendance:manage_config",
			},
		},
		{
			Code:        "manager",
			Name:        "Manager",
			Description: "Team manager with approval access for leave, attendance, and team oversight",
			Permissions: []string{
				// Dashboard
				"dashboard:view",
				// Employees (read + update for direct reports)
				"employee:read", "employee:update",
				// Leave (read + approve for team)
				"leave:read", "leave:create", "leave:approve",
				// Attendance
				"attendance:read", "attendance:approve",
				// Shifts
				"shift:read", "shift:approve_swaps",
				// Organization (read)
				"department:read", "team:read", "position:read",
				"organization:read", "location:read",
				// Helpdesk
				"helpdesk:read", "helpdesk:create",
				// Reports
				"report:view",
			},
		},
		{
			Code:        "employee",
			Name:        "Employee",
			Description: "Basic employee access with self-service features",
			Permissions: []string{
				// Dashboard
				"dashboard:view",
				// Employees (read own profile)
				"employee:read",
				// Leave (apply + view own)
				"leave:read", "leave:create",
				// Shifts (view own)
				"shift:read",
				// Organization (read)
				"department:read", "team:read", "position:read",
				"organization:read", "location:read",
				// Helpdesk (create own tickets)
				"helpdesk:read", "helpdesk:create",
				// Payroll (view own)
				"payroll:read",
				// Assets (view own)
				"asset:read",
				// Attendance (view own)
				"attendance:read",
			},
		},

		{
			Code:        "hod",
			Name:        "Head of Department",
			Description: "Department head with oversight for department employees, leave, attendance, and approvals",
			Permissions: []string{
				// Dashboard
				"dashboard:view", "dashboard:department_view",
				// Employees (department level access)
				"employee:read", "employee:update", "employee:department_view",
				// Leave (department approvals)
				"leave:read", "leave:create", "leave:approve", "leave:department_view",
				// Attendance (department oversight)
				"attendance:read", "attendance:approve", "attendance:department_view",
				// Shifts (department oversight)
				"shift:read", "shift:approve_swaps", "shift:department_view",
				// Organization (department management)
				"department:read", "department:update",
				"team:read", "team:update",
				"position:read", "position:update",
				"organization:read", "location:read",
				"cost_center:read",
				// Payroll (department reports)
				"payroll:read", "payroll:department_reports",
				// Helpdesk (department tickets)
				"helpdesk:read", "helpdesk:create", "helpdesk:department_view",
				// Reports (department level)
			"report:view", "report:department_export",
			},
		},
	}

	for _, roleData := range defaultRoles {
		existingRole, err := roleRepo.FindByCode(roleData.Code)
		if err == nil && existingRole != nil {
			// Role exists - update permissions if needed
			var permissionIDs []uint
			for _, permCode := range roleData.Permissions {
				perm, err := permRepo.FindByCode(permCode)
				if err == nil && perm != nil {
					permissionIDs = append(permissionIDs, perm.ID)
				}
			}
			if len(permissionIDs) > 0 {
				if err := roleRepo.AssignPermissions(existingRole.ID, permissionIDs); err != nil {
					log.Printf("Warning: Failed to update permissions for role %s: %v", roleData.Code, err)
				}
			}
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

		log.Printf("✅ Role created: %s (%d permissions)", roleData.Code, len(permissionIDs))
	}

	log.Println("✅ Default roles seeded successfully")
}
