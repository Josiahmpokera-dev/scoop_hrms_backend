package models

// Feature represents a system module/feature with its associated permissions.
// This is not a database model - it defines the feature catalog for the frontend.
type Feature struct {
	Code        string       `json:"code"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Icon        string       `json:"icon"`
	Category    string       `json:"category"`
	Permissions []FeaturePerm `json:"permissions"`
}

// FeaturePerm represents a permission within a feature
type FeaturePerm struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// FeatureAccess represents a feature with its access status for a specific role/user
type FeatureAccess struct {
	Code        string             `json:"code"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Icon        string             `json:"icon"`
	Category    string             `json:"category"`
	HasAccess   bool               `json:"has_access"`
	Permissions []PermissionAccess `json:"permissions"`
}

// PermissionAccess represents a permission with its granted status
type PermissionAccess struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Granted     bool   `json:"granted"`
}

// Feature categories
const (
	CategoryCore           = "core"
	CategoryHR             = "hr"
	CategoryOrganization   = "organization"
	CategoryFinance        = "finance"
	CategoryOperations     = "operations"
	CategorySupport        = "support"
	CategoryAdministration = "administration"
)

// GetAllFeatures returns the complete feature catalog for the system
func GetAllFeatures() []Feature {
	return []Feature{
		// ======= Core =======
		{
			Code:        "dashboard",
			Name:        "Dashboard",
			Description: "View dashboard statistics, announcements, and quick actions",
			Icon:        "dashboard",
			Category:    CategoryCore,
			Permissions: []FeaturePerm{
				{Code: "dashboard:view", Name: "View Dashboard", Description: "Access the main dashboard"},
				{Code: "dashboard:manage_announcements", Name: "Manage Announcements", Description: "Create, edit, and delete company announcements"},
			},
		},

		// ======= HR =======
		{
			Code:        "employees",
			Name:        "Employee Management",
			Description: "Manage employee records, onboarding, and offboarding",
			Icon:        "people",
			Category:    CategoryHR,
			Permissions: []FeaturePerm{
				{Code: "employee:read", Name: "View Employees", Description: "View employee profiles and records"},
				{Code: "employee:create", Name: "Create Employees", Description: "Add new employee records"},
				{Code: "employee:update", Name: "Update Employees", Description: "Edit employee information"},
				{Code: "employee:delete", Name: "Delete Employees", Description: "Remove employee records"},
				{Code: "employee:onboard", Name: "Onboard Employees", Description: "Run the employee onboarding process"},
				{Code: "employee:offboard", Name: "Offboard Employees", Description: "Run the employee offboarding process"},
			},
		},
		{
			Code:        "leave",
			Name:        "Leave Management",
			Description: "Manage leave types, policies, requests, and approvals",
			Icon:        "event_busy",
			Category:    CategoryHR,
			Permissions: []FeaturePerm{
				{Code: "leave:read", Name: "View Leave", Description: "View leave requests and balances"},
				{Code: "leave:create", Name: "Apply for Leave", Description: "Submit leave applications"},
				{Code: "leave:approve", Name: "Approve Leave", Description: "Approve or reject leave requests"},
				{Code: "leave:manage_types", Name: "Manage Leave Types", Description: "Create and configure leave types"},
				{Code: "leave:manage_policies", Name: "Manage Leave Policies", Description: "Create and configure leave policies"},
				{Code: "leave:manage_holidays", Name: "Manage Holidays", Description: "Create and manage public holidays"},
			},
		},
		{
			Code:        "attendance",
			Name:        "Attendance & Biometric",
			Description: "Track attendance, manage biometric devices, and approve records",
			Icon:        "fingerprint",
			Category:    CategoryHR,
			Permissions: []FeaturePerm{
				{Code: "attendance:read", Name: "View Attendance", Description: "View attendance records and reports"},
				{Code: "attendance:approve", Name: "Approve Attendance", Description: "Approve attendance corrections"},
				{Code: "attendance:manage_config", Name: "Manage Biometric Config", Description: "Configure biometric devices and settings"},
			},
		},
		{
			Code:        "shifts",
			Name:        "Shifts & Rosters",
			Description: "Manage shifts, roster assignments, and swap requests",
			Icon:        "schedule",
			Category:    CategoryHR,
			Permissions: []FeaturePerm{
				{Code: "shift:read", Name: "View Shifts", Description: "View shifts and roster assignments"},
				{Code: "shift:create", Name: "Create Shifts", Description: "Create new shift definitions"},
				{Code: "shift:update", Name: "Update Shifts", Description: "Modify shift details"},
				{Code: "shift:delete", Name: "Delete Shifts", Description: "Remove shift definitions"},
				{Code: "shift:manage_rosters", Name: "Manage Rosters", Description: "Create and manage roster assignments"},
				{Code: "shift:approve_swaps", Name: "Approve Swap Requests", Description: "Approve or reject shift swap requests"},
			},
		},

		// ======= Organization =======
		{
			Code:        "departments",
			Name:        "Departments",
			Description: "Manage organizational departments",
			Icon:        "business",
			Category:    CategoryOrganization,
			Permissions: []FeaturePerm{
				{Code: "department:read", Name: "View Departments", Description: "View department information"},
				{Code: "department:create", Name: "Create Departments", Description: "Add new departments"},
				{Code: "department:update", Name: "Update Departments", Description: "Edit department details"},
				{Code: "department:delete", Name: "Delete Departments", Description: "Remove departments"},
			},
		},
		{
			Code:        "teams",
			Name:        "Teams",
			Description: "Manage teams within departments",
			Icon:        "groups",
			Category:    CategoryOrganization,
			Permissions: []FeaturePerm{
				{Code: "team:read", Name: "View Teams", Description: "View team information"},
				{Code: "team:create", Name: "Create Teams", Description: "Add new teams"},
				{Code: "team:update", Name: "Update Teams", Description: "Edit team details"},
				{Code: "team:delete", Name: "Delete Teams", Description: "Remove teams"},
			},
		},
		{
			Code:        "positions",
			Name:        "Job Positions",
			Description: "Manage job positions and roles",
			Icon:        "work",
			Category:    CategoryOrganization,
			Permissions: []FeaturePerm{
				{Code: "position:read", Name: "View Positions", Description: "View job positions"},
				{Code: "position:create", Name: "Create Positions", Description: "Add new job positions"},
				{Code: "position:update", Name: "Update Positions", Description: "Edit job position details"},
				{Code: "position:delete", Name: "Delete Positions", Description: "Remove job positions"},
			},
		},
		{
			Code:        "organizations",
			Name:        "Organizations",
			Description: "Manage organization details and units",
			Icon:        "corporate_fare",
			Category:    CategoryOrganization,
			Permissions: []FeaturePerm{
				{Code: "organization:read", Name: "View Organizations", Description: "View organization information"},
				{Code: "organization:create", Name: "Create Organizations", Description: "Add new organizations"},
				{Code: "organization:update", Name: "Update Organizations", Description: "Edit organization details"},
				{Code: "organization:delete", Name: "Delete Organizations", Description: "Remove organizations"},
			},
		},
		{
			Code:        "organization_units",
			Name:        "Organization Units",
			Description: "Manage organization unit hierarchy",
			Icon:        "account_tree",
			Category:    CategoryOrganization,
			Permissions: []FeaturePerm{
				{Code: "organization_unit:read", Name: "View Organization Units", Description: "View organization units"},
				{Code: "organization_unit:create", Name: "Create Organization Units", Description: "Add new organization units"},
				{Code: "organization_unit:update", Name: "Update Organization Units", Description: "Edit organization unit details"},
				{Code: "organization_unit:delete", Name: "Delete Organization Units", Description: "Remove organization units"},
			},
		},
		{
			Code:        "locations",
			Name:        "Locations",
			Description: "Manage office locations and branches",
			Icon:        "location_on",
			Category:    CategoryOrganization,
			Permissions: []FeaturePerm{
				{Code: "location:read", Name: "View Locations", Description: "View location information"},
				{Code: "location:create", Name: "Create Locations", Description: "Add new locations"},
				{Code: "location:update", Name: "Update Locations", Description: "Edit location details"},
				{Code: "location:delete", Name: "Delete Locations", Description: "Remove locations"},
			},
		},
		{
			Code:        "cost_centers",
			Name:        "Cost Centers",
			Description: "Manage cost centers for budgeting",
			Icon:        "account_balance",
			Category:    CategoryOrganization,
			Permissions: []FeaturePerm{
				{Code: "cost_center:read", Name: "View Cost Centers", Description: "View cost center information"},
				{Code: "cost_center:create", Name: "Create Cost Centers", Description: "Add new cost centers"},
				{Code: "cost_center:update", Name: "Update Cost Centers", Description: "Edit cost center details"},
				{Code: "cost_center:delete", Name: "Delete Cost Centers", Description: "Remove cost centers"},
			},
		},

		// ======= Finance =======
		{
			Code:        "payroll",
			Name:        "Payroll",
			Description: "Run payroll, manage salary structures, loans, and compliance",
			Icon:        "payments",
			Category:    CategoryFinance,
			Permissions: []FeaturePerm{
				{Code: "payroll:read", Name: "View Payroll", Description: "View payroll runs, payslips, and reports"},
				{Code: "payroll:run", Name: "Run Payroll", Description: "Execute payroll calculations and processing"},
				{Code: "payroll:manage_structures", Name: "Manage Salary Structures", Description: "Create and edit salary structures and components"},
				{Code: "payroll:manage_loans", Name: "Manage Loans", Description: "Approve and manage employee loans"},
				{Code: "payroll:manage_compliance", Name: "Manage Compliance", Description: "Configure tax slabs and statutory rules"},
				{Code: "payroll:manage_reports", Name: "Manage Reports", Description: "Generate and export payroll reports"},
			},
		},

		// ======= Operations =======
		{
			Code:        "assets",
			Name:        "Asset Management",
			Description: "Track and manage company assets",
			Icon:        "devices",
			Category:    CategoryOperations,
			Permissions: []FeaturePerm{
				{Code: "asset:read", Name: "View Assets", Description: "View asset inventory"},
				{Code: "asset:create", Name: "Create Assets", Description: "Add new assets to inventory"},
				{Code: "asset:update", Name: "Update Assets", Description: "Edit asset information"},
				{Code: "asset:delete", Name: "Delete Assets", Description: "Remove assets from inventory"},
				{Code: "asset:assign", Name: "Assign Assets", Description: "Assign and reassign assets to employees"},
			},
		},

		// ======= Support =======
		{
			Code:        "helpdesk",
			Name:        "Helpdesk",
			Description: "Support ticket management and knowledge base",
			Icon:        "support_agent",
			Category:    CategorySupport,
			Permissions: []FeaturePerm{
				{Code: "helpdesk:read", Name: "View Tickets", Description: "View helpdesk tickets"},
				{Code: "helpdesk:create", Name: "Create Tickets", Description: "Submit support tickets"},
				{Code: "helpdesk:manage", Name: "Manage Tickets", Description: "Assign, resolve, and close tickets"},
				{Code: "helpdesk:manage_kb", Name: "Manage Knowledge Base", Description: "Create and manage knowledge base articles"},
			},
		},

		// ======= Administration =======
		{
			Code:        "users",
			Name:        "User Management",
			Description: "Manage user accounts, roles, and permissions",
			Icon:        "manage_accounts",
			Category:    CategoryAdministration,
			Permissions: []FeaturePerm{
				{Code: "user:read", Name: "View Users", Description: "View user accounts"},
				{Code: "user:create", Name: "Create Users", Description: "Add new user accounts"},
				{Code: "user:update", Name: "Update Users", Description: "Edit user account details"},
				{Code: "user:delete", Name: "Delete Users", Description: "Remove user accounts"},
				{Code: "user:manage_roles", Name: "Manage User Roles", Description: "Assign and revoke roles"},
			},
		},
		{
			Code:        "audit",
			Name:        "Audit & Security",
			Description: "View audit logs and security events",
			Icon:        "security",
			Category:    CategoryAdministration,
			Permissions: []FeaturePerm{
				{Code: "audit:read", Name: "View Audit Logs", Description: "View system audit trails"},
				{Code: "audit:export", Name: "Export Audit Logs", Description: "Export audit log data"},
			},
		},
		{
			Code:        "reports",
			Name:        "Reports",
			Description: "View and export system-wide reports",
			Icon:        "assessment",
			Category:    CategoryAdministration,
			Permissions: []FeaturePerm{
				{Code: "report:view", Name: "View Reports", Description: "View system reports and analytics"},
				{Code: "report:export", Name: "Export Reports", Description: "Export report data"},
			},
		},
		{
			Code:        "settings",
			Name:        "Settings",
			Description: "Manage system configuration and settings",
			Icon:        "settings",
			Category:    CategoryAdministration,
			Permissions: []FeaturePerm{
				{Code: "settings:read", Name: "View Settings", Description: "View system settings"},
				{Code: "settings:update", Name: "Update Settings", Description: "Modify system configuration"},
			},
		},
	}
}

// GetAllPermissionCodes returns all permission codes from the feature catalog
func GetAllPermissionCodes() []string {
	var codes []string
	for _, feature := range GetAllFeatures() {
		for _, perm := range feature.Permissions {
			codes = append(codes, perm.Code)
		}
	}
	return codes
}
