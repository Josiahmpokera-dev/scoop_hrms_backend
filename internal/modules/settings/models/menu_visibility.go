package models

import (
	"time"

	"gorm.io/gorm"
)

// MenuVisibilitySetting stores per-key visibility overrides.
// Items not present in this table are visible by default.
type MenuVisibilitySetting struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	MenuKey   string         `json:"menu_key" gorm:"uniqueIndex;not null;size:100"`
	Visible   bool           `json:"visible" gorm:"not null;default:true;index"`
	UpdatedBy *uint          `json:"updated_by,omitempty" gorm:"index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (MenuVisibilitySetting) TableName() string { return "menu_visibility_settings" }

// ---------- Request DTOs ----------

type BulkUpdateVisibilityRequest struct {
	Items []VisibilityItem `json:"items" binding:"required,min=1,dive"`
}

type VisibilityItem struct {
	MenuKey string `json:"menu_key" binding:"required"`
	Visible bool   `json:"visible"`
}

type ToggleVisibilityRequest struct {
	Visible bool `json:"visible"`
}

// ---------- Menu tree definition ----------

type MenuNode struct {
	Key      string      `json:"menu_key"`
	Title    string      `json:"title"`
	Type     string      `json:"type"` // title, collapse, item
	Path     string      `json:"path,omitempty"`
	Visible  bool        `json:"visible"`
	Children []*MenuNode `json:"children,omitempty"`
}

// protectedKeys that can never be hidden.
var ProtectedKeys = map[string]bool{
	"category.overview":      true,
	"dashboard":              true,
	"category.administration": true,
	"settingsAdmin":          true,
	"settings.system":        true,
}

// DefaultMenuTree returns the full menu hierarchy with all items defaulting to visible.
func DefaultMenuTree() []*MenuNode {
	return []*MenuNode{
		{Key: "category.overview", Title: "Overview", Type: "title", Children: []*MenuNode{
			{Key: "dashboard", Title: "Dashboard", Type: "item", Path: "/dashboard"},
		}},
		{Key: "category.hrPeople", Title: "HR & People", Type: "title", Children: []*MenuNode{
			{Key: "employeeManagement", Title: "Employee Management", Type: "collapse", Children: []*MenuNode{
				{Key: "employees.directory", Title: "Employee Directory", Type: "item", Path: "/employees/directory"},
				{Key: "employees.organization", Title: "Organization Structure", Type: "item", Path: "/employees/organization"},
				{Key: "employees.assets", Title: "Assets & Equipment", Type: "item", Path: "/employees/assets"},
				{Key: "employees.documents", Title: "Documents", Type: "item", Path: "/employees/documents"},
				{Key: "employees.activeOnboarding", Title: "Active Onboarding", Type: "item", Path: "/employees/active-onboarding"},
				{Key: "employees.offboarding", Title: "Offboarding", Type: "item", Path: "/employees/offboarding"},
			}},
		}},
		{Key: "category.timeLeave", Title: "Time & Leave", Type: "title", Children: []*MenuNode{
			{Key: "timeAttendance", Title: "Time & Attendance", Type: "collapse", Children: []*MenuNode{
				{Key: "attendance.daily", Title: "Daily Attendance", Type: "item", Path: "/attendance/daily"},
				{Key: "attendance.shifts", Title: "Shift Management", Type: "item", Path: "/attendance/shifts"},
				{Key: "attendance.timesheets", Title: "Timesheets", Type: "item", Path: "/attendance/timesheets"},
				{Key: "attendance.overtime", Title: "Overtime", Type: "item", Path: "/attendance/overtime"},
				{Key: "attendance.reports", Title: "Attendance Reports", Type: "item", Path: "/attendance/reports"},
			}},
			{Key: "leaveManagement", Title: "Leave Management", Type: "collapse", Children: []*MenuNode{
				{Key: "leave.apply", Title: "Apply Leave", Type: "item", Path: "/leave/apply"},
				{Key: "leave.requests", Title: "Leave Requests", Type: "item", Path: "/leave/requests"},
				{Key: "leave.calendar", Title: "Leave Calendar", Type: "item", Path: "/leave/calendar"},
				{Key: "leave.policies", Title: "Leave Policies", Type: "item", Path: "/leave/policies"},
				{Key: "leave.holidays", Title: "Holidays", Type: "item", Path: "/leave/holidays"},
				{Key: "leave.reports", Title: "Leave Reports", Type: "item", Path: "/leave/reports"},
			}},
		}},
		{Key: "category.payroll", Title: "Payroll & Finance", Type: "title", Children: []*MenuNode{
			{Key: "payrollManagement", Title: "Payroll", Type: "collapse", Children: []*MenuNode{
				{Key: "payroll.dashboard", Title: "Dashboard", Type: "item", Path: "/payroll/dashboard"},
				{Key: "payroll.run", Title: "Run Payroll", Type: "item", Path: "/payroll/run"},
				{Key: "payroll.salaryStructure", Title: "Salary Structure", Type: "item", Path: "/payroll/salary-structure"},
				{Key: "payroll.payslips", Title: "Payslips", Type: "item", Path: "/payroll/payslips"},
				{Key: "payroll.loans", Title: "Loans & Advances", Type: "item", Path: "/payroll/loans"},
				{Key: "payroll.compliance", Title: "Tax & Compliance", Type: "item", Path: "/payroll/compliance"},
				{Key: "payroll.reports", Title: "Payroll Reports", Type: "item", Path: "/payroll/reports"},
			}},
		}},
		{Key: "category.talent", Title: "Talent Management", Type: "title", Children: []*MenuNode{
			{Key: "recruitmentManagement", Title: "Recruitment", Type: "collapse", Children: []*MenuNode{
				{Key: "recruitment.dashboard", Title: "Dashboard", Type: "item", Path: "/recruitment/dashboard"},
				{Key: "recruitment.requisitions", Title: "Requisitions", Type: "item", Path: "/recruitment/requisitions"},
				{Key: "recruitment.openings", Title: "Job Openings", Type: "item", Path: "/recruitment/openings"},
				{Key: "recruitment.candidates", Title: "Candidates", Type: "item", Path: "/recruitment/candidates"},
				{Key: "recruitment.interviews", Title: "Interviews", Type: "item", Path: "/recruitment/interviews"},
				{Key: "recruitment.offers", Title: "Offers", Type: "item", Path: "/recruitment/offers"},
				{Key: "recruitment.talentPool", Title: "Talent Pool", Type: "item", Path: "/recruitment/talent-pool"},
			}},
			{Key: "performanceManagement", Title: "Manage Performance", Type: "collapse", Children: []*MenuNode{
				{Key: "performance.dashboard", Title: "Dashboard", Type: "item", Path: "/performance/dashboard"},
				{Key: "performance.goals", Title: "Goals & OKRs", Type: "item", Path: "/performance/goals"},
				{Key: "performance.appraisals", Title: "Appraisals", Type: "item", Path: "/performance/appraisals"},
				{Key: "performance.feedback360", Title: "360° Feedback", Type: "item", Path: "/performance/feedback-360"},
				{Key: "performance.talentReview", Title: "Talent Review", Type: "item", Path: "/performance/talent-review"},
				{Key: "performance.reports", Title: "Reports", Type: "item", Path: "/performance/reports"},
			}},
			{Key: "trainingDevelopment", Title: "Training & Development", Type: "collapse", Children: []*MenuNode{
				{Key: "training.dashboard", Title: "Learning Dashboard", Type: "item", Path: "/training/dashboard"},
				{Key: "training.catalog", Title: "Course Catalog", Type: "item", Path: "/training/catalog"},
				{Key: "training.myLearning", Title: "My Learning", Type: "item", Path: "/training/my-learning"},
				{Key: "training.sessions", Title: "Training Sessions", Type: "item", Path: "/training/sessions"},
				{Key: "training.certifications", Title: "Certifications", Type: "item", Path: "/training/certifications"},
				{Key: "training.reports", Title: "Training Reports", Type: "item", Path: "/training/reports"},
			}},
			{Key: "projectManagement", Title: "Manage Projects", Type: "collapse", Children: []*MenuNode{
				{Key: "projects.allProjects", Title: "All Projects", Type: "item", Path: "/projects/all"},
				{Key: "projects.myProjects", Title: "My Projects", Type: "item", Path: "/projects/mine"},
				{Key: "projects.dailyTasks", Title: "Daily Tasks", Type: "item", Path: "/projects/daily-tasks"},
			}},
		}},
		{Key: "category.mySpace", Title: "My Space", Type: "title", Children: []*MenuNode{
			{Key: "selfService", Title: "Self Service", Type: "collapse", Children: []*MenuNode{
				{Key: "selfService.myProfile", Title: "My Profile", Type: "item", Path: "/self-service/profile"},
				{Key: "selfService.myPayslips", Title: "My Payslips", Type: "item", Path: "/self-service/payslips"},
				{Key: "selfService.requests", Title: "My Requests", Type: "item", Path: "/self-service/requests"},
				{Key: "selfService.directory", Title: "People Directory", Type: "item", Path: "/self-service/directory"},
				{Key: "selfService.myAssets", Title: "My Assets", Type: "item", Path: "/self-service/assets"},
				{Key: "selfService.requestLeave", Title: "Request Leave", Type: "item", Path: "/self-service/request-leave"},
			}},
		}},
		{Key: "category.supportReports", Title: "Support & Reports", Type: "title", Children: []*MenuNode{
			{Key: "helpdeskSupport", Title: "Helpdesk & Support", Type: "collapse", Children: []*MenuNode{
				{Key: "helpdesk.myTickets", Title: "My Tickets", Type: "item", Path: "/helpdesk/my-tickets"},
				{Key: "helpdesk.knowledgeBase", Title: "Knowledge Base", Type: "item", Path: "/helpdesk/knowledge-base"},
				{Key: "helpdesk.dashboard", Title: "Helpdesk Dashboard", Type: "item", Path: "/helpdesk/dashboard"},
			}},
			{Key: "reportsAnalytics", Title: "Reports & Analytics", Type: "collapse", Children: []*MenuNode{
				{Key: "reports.analytics", Title: "Analytics Dashboard", Type: "item", Path: "/reports/analytics"},
				{Key: "reports.standard", Title: "Standard Reports", Type: "item", Path: "/reports/standard"},
				{Key: "reports.builder", Title: "Report Builder", Type: "item", Path: "/reports/builder"},
				{Key: "reports.scheduled", Title: "Scheduled Reports", Type: "item", Path: "/reports/scheduled"},
			}},
		}},
		{Key: "category.administration", Title: "Administration", Type: "title", Children: []*MenuNode{
			{Key: "settingsAdmin", Title: "Settings & Admin", Type: "collapse", Children: []*MenuNode{
				{Key: "settings.organization", Title: "Organization", Type: "item", Path: "/settings/organization"},
				{Key: "settings.users", Title: "Users & Roles", Type: "item", Path: "/settings/users"},
				{Key: "settings.policies", Title: "Policies", Type: "item", Path: "/settings/policies"},
				{Key: "settings.integrations", Title: "Integrations", Type: "item", Path: "/settings/integrations"},
				{Key: "settings.security", Title: "Security", Type: "item", Path: "/settings/security"},
				{Key: "settings.system", Title: "System Settings", Type: "item", Path: "/settings/system"},
			}},
		}},
	}
}

// AllMenuKeys returns a flat set of every known menu key.
func AllMenuKeys() map[string]bool {
	keys := make(map[string]bool)
	var walk func(nodes []*MenuNode)
	walk = func(nodes []*MenuNode) {
		for _, n := range nodes {
			keys[n.Key] = true
			walk(n.Children)
		}
	}
	walk(DefaultMenuTree())
	return keys
}

// ParentMap returns child→parent mapping and parent→children mapping.
func ParentMap() (childToParent map[string]string, parentToChildren map[string][]string) {
	childToParent = make(map[string]string)
	parentToChildren = make(map[string][]string)
	var walk func(nodes []*MenuNode, parentKey string)
	walk = func(nodes []*MenuNode, parentKey string) {
		for _, n := range nodes {
			if parentKey != "" {
				childToParent[n.Key] = parentKey
				parentToChildren[parentKey] = append(parentToChildren[parentKey], n.Key)
			}
			walk(n.Children, n.Key)
		}
	}
	walk(DefaultMenuTree(), "")
	return
}
