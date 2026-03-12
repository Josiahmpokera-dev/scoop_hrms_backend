package seed

import (
	"log"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	deptModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/models"
	deptRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
	empModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	empRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	leaveModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	leaveRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/repositories"
	locModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/models"
	locRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/repositories"
	posModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/models"
	posRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/repositories"
	teamModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/teams/models"
	teamRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/teams/repositories"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	userRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories"
	"golang.org/x/crypto/bcrypt"
)

// RunDevelopmentData seeds comprehensive test data for development & testing.
// This creates: locations, departments, positions, teams, users, employees, and leave balances.
// All seeding is idempotent — existing records are skipped.
func RunDevelopmentData() {
	log.Println("Seeding development test data...")

	// 1. Seed Locations
	locations := seedLocations()

	// 2. Seed Departments (with hierarchy)
	departments := seedDepartments(locations)

	// 3. Seed Job Positions (linked to departments)
	positions := seedPositions(departments)

	// 4. Seed Teams (linked to departments)
	teams := seedTeams(departments, locations)

	// 5. Seed Test Users + Employees
	seedTestEmployees(locations, departments, positions, teams)

	// 6. Seed Leave Types
	seedLeaveTypes()

	// 7. Seed Leave Policies
	seedLeavePolicies()

	// 8. Seed Leave Balances for all test employees
	seedLeaveBalances()

	log.Println("✅ Development test data seeding completed.")
}

// ─── Locations ──────────────────────────────────────────────────────────────

func seedLocations() map[string]*locModels.Location {
	repo := locRepos.NewLocationRepository()
	result := make(map[string]*locModels.Location)

	defs := []struct {
		Name         string
		Type         string
		City         string
		Country      string
		Timezone     string
		IsHeadOffice bool
	}{
		{"Dar es Salaam HQ", "office", "Dar es Salaam", "Tanzania", "Africa/Dar_es_Salaam", true},
		{"Dodoma Branch", "branch", "Dodoma", "Tanzania", "Africa/Dar_es_Salaam", false},
		{"Arusha Branch", "branch", "Arusha", "Tanzania", "Africa/Dar_es_Salaam", false},
		{"Mwanza Branch", "branch", "Mwanza", "Tanzania", "Africa/Dar_es_Salaam", false},
		{"Remote Office", "remote", "Remote", "Tanzania", "Africa/Dar_es_Salaam", false},
	}

	for _, d := range defs {
		existing, _ := repo.FindByName(d.Name, nil, nil)
		if existing != nil {
			result[d.Name] = existing
			continue
		}

		locType := d.Type
		city := d.City
		country := d.Country
		tz := d.Timezone

		loc := &locModels.Location{
			Name:         d.Name,
			LocationType: &locType,
			City:         &city,
			Country:      &country,
			Timezone:     &tz,
			IsHeadOffice: d.IsHeadOffice,
			IsActive:     true,
		}
		if err := repo.Create(loc); err != nil {
			log.Printf("  Warning: Failed to create location %s: %v", d.Name, err)
			continue
		}
		result[d.Name] = loc
		log.Printf("  ✅ Location created: %s", d.Name)
	}

	return result
}

// ─── Departments ────────────────────────────────────────────────────────────

func seedDepartments(locations map[string]*locModels.Location) map[string]*deptModels.Department {
	repo := deptRepos.NewDepartmentRepository()
	result := make(map[string]*deptModels.Department)

	var hqLocationID *uint
	if loc, ok := locations["Dar es Salaam HQ"]; ok {
		hqLocationID = &loc.ID
	}

	// Top-level departments
	topLevel := []struct {
		Code string
		Name string
		Desc string
		Type string
	}{
		{"EXEC", "Executive Management", "C-suite and executive leadership team", "executive"},
		{"HR", "Human Resources", "People operations, recruitment, and employee welfare", "support"},
		{"IT", "Information Technology", "Software development, infrastructure, and IT support", "core"},
		{"FIN", "Finance & Accounting", "Financial planning, accounting, and treasury", "support"},
		{"MKT", "Marketing & Communications", "Brand management, digital marketing, and PR", "core"},
		{"OPS", "Operations", "Business operations and logistics management", "core"},
		{"SALES", "Sales & Business Development", "Revenue generation and client relationships", "core"},
		{"LEGAL", "Legal & Compliance", "Corporate governance, legal affairs, and compliance", "support"},
		{"QA", "Quality Assurance", "Testing, quality control, and standards", "support"},
	}

	for _, d := range topLevel {
		existing, _ := repo.FindByCode(d.Code)
		if existing != nil {
			result[d.Code] = existing
			continue
		}

		desc := d.Desc
		deptType := d.Type
		dept := &deptModels.Department{
			Code:           d.Code,
			Name:           d.Name,
			Description:    &desc,
			DepartmentType: &deptType,
			LocationID:     hqLocationID,
			IsActive:       true,
		}
		if err := repo.Create(dept); err != nil {
			log.Printf("  Warning: Failed to create department %s: %v", d.Code, err)
			continue
		}
		result[d.Code] = dept
		log.Printf("  ✅ Department created: %s (%s)", d.Name, d.Code)
	}

	// Sub-departments (children)
	subDepts := []struct {
		Code     string
		Name     string
		Desc     string
		ParentID string // Code of parent
	}{
		{"IT-DEV", "Software Development", "Application and platform development team", "IT"},
		{"IT-INFRA", "Infrastructure & DevOps", "Server, cloud, and network infrastructure", "IT"},
		{"IT-SUP", "IT Support & Helpdesk", "End-user support and internal helpdesk", "IT"},
		{"HR-REC", "Recruitment & Talent", "Hiring, onboarding, and talent acquisition", "HR"},
		{"HR-PA", "Payroll & Administration", "Payroll processing and HR administration", "HR"},
		{"FIN-ACC", "Accounting", "General ledger, AP/AR, and financial reporting", "FIN"},
		{"FIN-BUD", "Budgeting & Planning", "Financial planning, budgets, and forecasting", "FIN"},
		{"MKT-DIG", "Digital Marketing", "Online campaigns, SEO, and social media", "MKT"},
		{"OPS-LOG", "Logistics", "Supply chain and distribution management", "OPS"},
		{"SALES-BD", "Business Development", "New market and partnership opportunities", "SALES"},
	}

	for _, d := range subDepts {
		existing, _ := repo.FindByCode(d.Code)
		if existing != nil {
			result[d.Code] = existing
			continue
		}

		parent, ok := result[d.ParentID]
		if !ok {
			log.Printf("  Warning: Parent dept %s not found for %s, skipping", d.ParentID, d.Code)
			continue
		}

		desc := d.Desc
		dept := &deptModels.Department{
			Code:               d.Code,
			Name:               d.Name,
			Description:        &desc,
			ParentDepartmentID: &parent.ID,
			LocationID:         hqLocationID,
			IsActive:           true,
		}
		if err := repo.Create(dept); err != nil {
			log.Printf("  Warning: Failed to create sub-department %s: %v", d.Code, err)
			continue
		}
		result[d.Code] = dept
		log.Printf("  ✅ Sub-department created: %s (%s) under %s", d.Name, d.Code, d.ParentID)
	}

	return result
}

// ─── Positions ──────────────────────────────────────────────────────────────

func seedPositions(departments map[string]*deptModels.Department) map[string]*posModels.JobPosition {
	repo := posRepos.NewJobPositionRepository()
	result := make(map[string]*posModels.JobPosition)

	defs := []struct {
		Code   string
		Title  string
		Grade  string
		DeptID string // Department code
	}{
		// Executive
		{"CEO", "Chief Executive Officer", "E1", "EXEC"},
		{"CFO", "Chief Financial Officer", "E2", "EXEC"},
		{"CTO", "Chief Technology Officer", "E2", "EXEC"},
		{"COO", "Chief Operations Officer", "E2", "EXEC"},

		// HR
		{"HR-DIR", "HR Director", "D1", "HR"},
		{"HR-MGR", "HR Manager", "M1", "HR"},
		{"HR-OFF", "HR Officer", "S2", "HR"},
		{"HR-REC", "Recruiter", "S2", "HR-REC"},
		{"HR-AST", "HR Assistant", "J1", "HR"},

		// IT
		{"IT-DIR", "IT Director", "D1", "IT"},
		{"IT-MGR", "IT Manager", "M1", "IT"},
		{"DEV-SR", "Senior Software Developer", "S3", "IT-DEV"},
		{"DEV-MID", "Software Developer", "S2", "IT-DEV"},
		{"DEV-JR", "Junior Software Developer", "J2", "IT-DEV"},
		{"DEVOPS", "DevOps Engineer", "S2", "IT-INFRA"},
		{"SYSADM", "System Administrator", "S2", "IT-INFRA"},
		{"IT-SUP", "IT Support Specialist", "S1", "IT-SUP"},

		// Finance
		{"FIN-DIR", "Finance Director", "D1", "FIN"},
		{"FIN-MGR", "Finance Manager", "M1", "FIN"},
		{"ACC-SR", "Senior Accountant", "S3", "FIN-ACC"},
		{"ACC", "Accountant", "S2", "FIN-ACC"},

		// Marketing
		{"MKT-DIR", "Marketing Director", "D1", "MKT"},
		{"MKT-MGR", "Marketing Manager", "M1", "MKT"},
		{"MKT-SPE", "Marketing Specialist", "S2", "MKT-DIG"},
		{"GRA-DES", "Graphic Designer", "S2", "MKT"},
		{"CNT-WRT", "Content Writer", "S1", "MKT-DIG"},

		// Operations
		{"OPS-DIR", "Operations Director", "D1", "OPS"},
		{"OPS-MGR", "Operations Manager", "M1", "OPS"},
		{"LOG-MGR", "Logistics Manager", "M1", "OPS-LOG"},
		{"LOG-OFF", "Logistics Officer", "S2", "OPS-LOG"},

		// Sales
		{"SLS-DIR", "Sales Director", "D1", "SALES"},
		{"SLS-MGR", "Sales Manager", "M1", "SALES"},
		{"SLS-REP", "Sales Representative", "S1", "SALES"},
		{"BD-MGR", "Business Development Manager", "M1", "SALES-BD"},

		// Legal
		{"LEG-DIR", "Legal Director", "D1", "LEGAL"},
		{"LEG-ADV", "Legal Advisor", "S3", "LEGAL"},
		{"CMP-OFF", "Compliance Officer", "S2", "LEGAL"},

		// QA
		{"QA-MGR", "QA Manager", "M1", "QA"},
		{"QA-ENG", "QA Engineer", "S2", "QA"},
		{"QA-ANL", "QA Analyst", "S1", "QA"},
	}

	for _, d := range defs {
		existing, _ := repo.FindByCode(d.Code)
		if existing != nil {
			result[d.Code] = existing
			continue
		}

		grade := d.Grade
		pos := &posModels.JobPosition{
			Code:     d.Code,
			Title:    d.Title,
			Grade:    &grade,
			IsActive: true,
		}
		if dept, ok := departments[d.DeptID]; ok {
			pos.DepartmentID = &dept.ID
		}

		if err := repo.Create(pos); err != nil {
			log.Printf("  Warning: Failed to create position %s: %v", d.Code, err)
			continue
		}
		result[d.Code] = pos
		log.Printf("  ✅ Position created: %s (%s)", d.Title, d.Code)
	}

	return result
}

// ─── Teams ──────────────────────────────────────────────────────────────────

func seedTeams(departments map[string]*deptModels.Department, locations map[string]*locModels.Location) map[string]*teamModels.Team {
	repo := teamRepos.NewTeamRepository()
	result := make(map[string]*teamModels.Team)

	var hqLocationID *uint
	if loc, ok := locations["Dar es Salaam HQ"]; ok {
		hqLocationID = &loc.ID
	}

	defs := []struct {
		Code       string
		Name       string
		Desc       string
		DeptCode   string
		TeamType   string
		MaxMembers int
	}{
		{"BACKEND", "Backend Team", "Server-side development and API team", "IT-DEV", "development", 8},
		{"FRONTEND", "Frontend Team", "Client-side UI/UX development team", "IT-DEV", "development", 6},
		{"MOBILE", "Mobile Team", "iOS and Android development team", "IT-DEV", "development", 5},
		{"DEVOPS-T", "DevOps Team", "CI/CD, cloud infrastructure, and deployment", "IT-INFRA", "operations", 4},
		{"HELPDESK", "IT Helpdesk Team", "End-user technical support team", "IT-SUP", "support", 5},
		{"RECRUIT", "Recruitment Team", "Hiring and talent acquisition team", "HR-REC", "support", 4},
		{"PAYROLL", "Payroll Team", "Salary processing and benefits team", "HR-PA", "support", 3},
		{"DIGITAL", "Digital Marketing Team", "Online campaigns and SEO team", "MKT-DIG", "marketing", 5},
		{"ACCT", "Accounting Team", "Financial records and reporting team", "FIN-ACC", "support", 4},
		{"BD-TEAM", "Business Development Team", "New market expansion team", "SALES-BD", "sales", 5},
		{"QA-AUTO", "QA Automation Team", "Automated testing and quality assurance", "QA", "development", 4},
	}

	for _, d := range defs {
		existing, _ := repo.FindByCode(d.Code)
		if existing != nil {
			result[d.Code] = existing
			continue
		}

		desc := d.Desc
		teamType := d.TeamType
		maxMem := d.MaxMembers
		team := &teamModels.Team{
			Code:        d.Code,
			Name:        d.Name,
			Description: &desc,
			TeamType:    &teamType,
			MaxMembers:  &maxMem,
			LocationID:  hqLocationID,
			IsActive:    true,
		}
		if dept, ok := departments[d.DeptCode]; ok {
			team.DepartmentID = &dept.ID
		}

		if err := repo.Create(team); err != nil {
			log.Printf("  Warning: Failed to create team %s: %v", d.Code, err)
			continue
		}
		result[d.Code] = team
		log.Printf("  ✅ Team created: %s (%s)", d.Name, d.Code)
	}

	return result
}

// ─── Test Employees (User + Employee) ───────────────────────────────────────

type testEmployeeDef struct {
	// User fields
	Username string
	Email    string
	Password string
	UserRole userModels.UserRole
	// Employee fields
	EmployeeID     string
	FirstName      string
	LastName       string
	Gender         string
	DOB            string // YYYY-MM-DD
	Nationality    string
	Phone          string
	EmploymentType string
	HireDate       string // YYYY-MM-DD
	// References (by code)
	DeptCode     string
	PositionCode string
	LocationName string
	TeamCode     string
	// Org hierarchy: manager's EmployeeID code
	ReportsToEmpID string
}

func seedTestEmployees(
	locations map[string]*locModels.Location,
	departments map[string]*deptModels.Department,
	positions map[string]*posModels.JobPosition,
	teams map[string]*teamModels.Team,
) {
	userRepo := userRepos.NewUserRepository()
	empRepo := empRepos.NewEmployeeRepository()

	defaultPassword := "Test@1234"

	employees := []testEmployeeDef{
		// ── Executive Management ──
		{
			Username: "john.mwanga", Email: "john.mwanga@hrms.com", Password: defaultPassword, UserRole: userModels.RoleAdmin,
			EmployeeID: "EMP-001", FirstName: "John", LastName: "Mwanga", Gender: "male",
			DOB: "1975-03-15", Nationality: "Tanzanian", Phone: "+255712000001",
			EmploymentType: "full_time", HireDate: "2018-01-10",
			DeptCode: "EXEC", PositionCode: "CEO", LocationName: "Dar es Salaam HQ",
		},
		{
			Username: "amina.hassan", Email: "amina.hassan@hrms.com", Password: defaultPassword, UserRole: userModels.RoleAdmin,
			EmployeeID: "EMP-002", FirstName: "Amina", LastName: "Hassan", Gender: "female",
			DOB: "1978-07-22", Nationality: "Tanzanian", Phone: "+255712000002",
			EmploymentType: "full_time", HireDate: "2018-06-01",
			DeptCode: "EXEC", PositionCode: "CTO", LocationName: "Dar es Salaam HQ",
			ReportsToEmpID: "EMP-001",
		},
		{
			Username: "peter.kamau", Email: "peter.kamau@hrms.com", Password: defaultPassword, UserRole: userModels.RoleAdmin,
			EmployeeID: "EMP-003", FirstName: "Peter", LastName: "Kamau", Gender: "male",
			DOB: "1976-11-05", Nationality: "Tanzanian", Phone: "+255712000003",
			EmploymentType: "full_time", HireDate: "2018-03-15",
			DeptCode: "EXEC", PositionCode: "CFO", LocationName: "Dar es Salaam HQ",
			ReportsToEmpID: "EMP-001",
		},
		{
			Username: "grace.nyerere", Email: "grace.nyerere@hrms.com", Password: defaultPassword, UserRole: userModels.RoleAdmin,
			EmployeeID: "EMP-004", FirstName: "Grace", LastName: "Nyerere", Gender: "female",
			DOB: "1980-02-14", Nationality: "Tanzanian", Phone: "+255712000004",
			EmploymentType: "full_time", HireDate: "2019-01-02",
			DeptCode: "EXEC", PositionCode: "COO", LocationName: "Dar es Salaam HQ",
			ReportsToEmpID: "EMP-001",
		},

		// ── HR Department ──
		{
			Username: "fatma.salim", Email: "fatma.salim@hrms.com", Password: defaultPassword, UserRole: userModels.RoleHR,
			EmployeeID: "EMP-010", FirstName: "Fatma", LastName: "Salim", Gender: "female",
			DOB: "1982-04-18", Nationality: "Tanzanian", Phone: "+255712000010",
			EmploymentType: "full_time", HireDate: "2019-03-01",
			DeptCode: "HR", PositionCode: "HR-DIR", LocationName: "Dar es Salaam HQ",
			ReportsToEmpID: "EMP-001",
		},
		{
			Username: "james.mushi", Email: "james.mushi@hrms.com", Password: defaultPassword, UserRole: userModels.RoleHR,
			EmployeeID: "EMP-011", FirstName: "James", LastName: "Mushi", Gender: "male",
			DOB: "1985-09-30", Nationality: "Tanzanian", Phone: "+255712000011",
			EmploymentType: "full_time", HireDate: "2020-01-15",
			DeptCode: "HR", PositionCode: "HR-MGR", LocationName: "Dar es Salaam HQ",
			ReportsToEmpID: "EMP-010",
		},
		{
			Username: "sarah.kimaro", Email: "sarah.kimaro@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-012", FirstName: "Sarah", LastName: "Kimaro", Gender: "female",
			DOB: "1990-12-05", Nationality: "Tanzanian", Phone: "+255712000012",
			EmploymentType: "full_time", HireDate: "2021-06-01",
			DeptCode: "HR", PositionCode: "HR-OFF", LocationName: "Dar es Salaam HQ",
			TeamCode: "RECRUIT", ReportsToEmpID: "EMP-011",
		},
		{
			Username: "david.mnali", Email: "david.mnali@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-013", FirstName: "David", LastName: "Mnali", Gender: "male",
			DOB: "1993-06-20", Nationality: "Tanzanian", Phone: "+255712000013",
			EmploymentType: "full_time", HireDate: "2022-02-01",
			DeptCode: "HR-REC", PositionCode: "HR-REC", LocationName: "Dar es Salaam HQ",
			TeamCode: "RECRUIT", ReportsToEmpID: "EMP-011",
		},
		{
			Username: "happiness.charles", Email: "happiness.charles@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-014", FirstName: "Happiness", LastName: "Charles", Gender: "female",
			DOB: "1995-01-10", Nationality: "Tanzanian", Phone: "+255712000014",
			EmploymentType: "full_time", HireDate: "2023-03-01",
			DeptCode: "HR-PA", PositionCode: "HR-AST", LocationName: "Dar es Salaam HQ",
			TeamCode: "PAYROLL", ReportsToEmpID: "EMP-011",
		},

		// ── IT Department ──
		{
			Username: "michael.temba", Email: "michael.temba@hrms.com", Password: defaultPassword, UserRole: userModels.RoleIT,
			EmployeeID: "EMP-020", FirstName: "Michael", LastName: "Temba", Gender: "male",
			DOB: "1981-08-12", Nationality: "Tanzanian", Phone: "+255712000020",
			EmploymentType: "full_time", HireDate: "2019-04-01",
			DeptCode: "IT", PositionCode: "IT-DIR", LocationName: "Dar es Salaam HQ",
			ReportsToEmpID: "EMP-002",
		},
		{
			Username: "rose.andrew", Email: "rose.andrew@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-021", FirstName: "Rose", LastName: "Andrew", Gender: "female",
			DOB: "1987-05-25", Nationality: "Tanzanian", Phone: "+255712000021",
			EmploymentType: "full_time", HireDate: "2020-07-01",
			DeptCode: "IT", PositionCode: "IT-MGR", LocationName: "Dar es Salaam HQ",
			ReportsToEmpID: "EMP-020",
		},
		{
			Username: "joseph.mosha", Email: "joseph.mosha@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-022", FirstName: "Joseph", LastName: "Mosha", Gender: "male",
			DOB: "1988-03-17", Nationality: "Tanzanian", Phone: "+255712000022",
			EmploymentType: "full_time", HireDate: "2020-09-15",
			DeptCode: "IT-DEV", PositionCode: "DEV-SR", LocationName: "Dar es Salaam HQ",
			TeamCode: "BACKEND", ReportsToEmpID: "EMP-021",
		},
		{
			Username: "lucy.mwakasege", Email: "lucy.mwakasege@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-023", FirstName: "Lucy", LastName: "Mwakasege", Gender: "female",
			DOB: "1991-11-08", Nationality: "Tanzanian", Phone: "+255712000023",
			EmploymentType: "full_time", HireDate: "2021-01-10",
			DeptCode: "IT-DEV", PositionCode: "DEV-SR", LocationName: "Dar es Salaam HQ",
			TeamCode: "FRONTEND", ReportsToEmpID: "EMP-021",
		},
		{
			Username: "alex.njau", Email: "alex.njau@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-024", FirstName: "Alex", LastName: "Njau", Gender: "male",
			DOB: "1993-07-12", Nationality: "Tanzanian", Phone: "+255712000024",
			EmploymentType: "full_time", HireDate: "2021-06-01",
			DeptCode: "IT-DEV", PositionCode: "DEV-MID", LocationName: "Dar es Salaam HQ",
			TeamCode: "BACKEND", ReportsToEmpID: "EMP-022",
		},
		{
			Username: "neema.john", Email: "neema.john@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-025", FirstName: "Neema", LastName: "John", Gender: "female",
			DOB: "1995-02-28", Nationality: "Tanzanian", Phone: "+255712000025",
			EmploymentType: "full_time", HireDate: "2022-03-15",
			DeptCode: "IT-DEV", PositionCode: "DEV-MID", LocationName: "Dar es Salaam HQ",
			TeamCode: "FRONTEND", ReportsToEmpID: "EMP-023",
		},
		{
			Username: "baraka.said", Email: "baraka.said@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-026", FirstName: "Baraka", LastName: "Said", Gender: "male",
			DOB: "1997-10-01", Nationality: "Tanzanian", Phone: "+255712000026",
			EmploymentType: "full_time", HireDate: "2023-01-10",
			DeptCode: "IT-DEV", PositionCode: "DEV-JR", LocationName: "Dar es Salaam HQ",
			TeamCode: "BACKEND", ReportsToEmpID: "EMP-022",
		},
		{
			Username: "rehema.ally", Email: "rehema.ally@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-027", FirstName: "Rehema", LastName: "Ally", Gender: "female",
			DOB: "1998-04-15", Nationality: "Tanzanian", Phone: "+255712000027",
			EmploymentType: "full_time", HireDate: "2023-06-01",
			DeptCode: "IT-DEV", PositionCode: "DEV-JR", LocationName: "Dar es Salaam HQ",
			TeamCode: "MOBILE", ReportsToEmpID: "EMP-021",
		},
		{
			Username: "frank.robert", Email: "frank.robert@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-028", FirstName: "Frank", LastName: "Robert", Gender: "male",
			DOB: "1989-06-30", Nationality: "Tanzanian", Phone: "+255712000028",
			EmploymentType: "full_time", HireDate: "2021-04-01",
			DeptCode: "IT-INFRA", PositionCode: "DEVOPS", LocationName: "Dar es Salaam HQ",
			TeamCode: "DEVOPS-T", ReportsToEmpID: "EMP-021",
		},
		{
			Username: "anna.kessy", Email: "anna.kessy@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-029", FirstName: "Anna", LastName: "Kessy", Gender: "female",
			DOB: "1992-08-20", Nationality: "Tanzanian", Phone: "+255712000029",
			EmploymentType: "full_time", HireDate: "2022-01-10",
			DeptCode: "IT-SUP", PositionCode: "IT-SUP", LocationName: "Dar es Salaam HQ",
			TeamCode: "HELPDESK", ReportsToEmpID: "EMP-021",
		},

		// ── Finance Department ──
		{
			Username: "daniel.lupogo", Email: "daniel.lupogo@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-030", FirstName: "Daniel", LastName: "Lupogo", Gender: "male",
			DOB: "1979-01-25", Nationality: "Tanzanian", Phone: "+255712000030",
			EmploymentType: "full_time", HireDate: "2019-02-01",
			DeptCode: "FIN", PositionCode: "FIN-DIR", LocationName: "Dar es Salaam HQ",
			ReportsToEmpID: "EMP-003",
		},
		{
			Username: "esther.mwakalinga", Email: "esther.mwakalinga@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-031", FirstName: "Esther", LastName: "Mwakalinga", Gender: "female",
			DOB: "1984-10-12", Nationality: "Tanzanian", Phone: "+255712000031",
			EmploymentType: "full_time", HireDate: "2020-04-01",
			DeptCode: "FIN", PositionCode: "FIN-MGR", LocationName: "Dar es Salaam HQ",
			TeamCode: "ACCT", ReportsToEmpID: "EMP-030",
		},
		{
			Username: "ibrahim.ally", Email: "ibrahim.ally@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-032", FirstName: "Ibrahim", LastName: "Ally", Gender: "male",
			DOB: "1990-05-18", Nationality: "Tanzanian", Phone: "+255712000032",
			EmploymentType: "full_time", HireDate: "2021-08-01",
			DeptCode: "FIN-ACC", PositionCode: "ACC-SR", LocationName: "Dar es Salaam HQ",
			TeamCode: "ACCT", ReportsToEmpID: "EMP-031",
		},
		{
			Username: "clara.lyimo", Email: "clara.lyimo@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-033", FirstName: "Clara", LastName: "Lyimo", Gender: "female",
			DOB: "1994-09-02", Nationality: "Tanzanian", Phone: "+255712000033",
			EmploymentType: "full_time", HireDate: "2022-11-01",
			DeptCode: "FIN-ACC", PositionCode: "ACC", LocationName: "Dar es Salaam HQ",
			TeamCode: "ACCT", ReportsToEmpID: "EMP-031",
		},

		// ── Marketing Department ──
		{
			Username: "gloria.massawe", Email: "gloria.massawe@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-040", FirstName: "Gloria", LastName: "Massawe", Gender: "female",
			DOB: "1983-06-14", Nationality: "Tanzanian", Phone: "+255712000040",
			EmploymentType: "full_time", HireDate: "2019-05-01",
			DeptCode: "MKT", PositionCode: "MKT-DIR", LocationName: "Dar es Salaam HQ",
			ReportsToEmpID: "EMP-004",
		},
		{
			Username: "george.msemwa", Email: "george.msemwa@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-041", FirstName: "George", LastName: "Msemwa", Gender: "male",
			DOB: "1988-12-01", Nationality: "Tanzanian", Phone: "+255712000041",
			EmploymentType: "full_time", HireDate: "2020-10-01",
			DeptCode: "MKT", PositionCode: "MKT-MGR", LocationName: "Dar es Salaam HQ",
			TeamCode: "DIGITAL", ReportsToEmpID: "EMP-040",
		},
		{
			Username: "mary.mdee", Email: "mary.mdee@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-042", FirstName: "Mary", LastName: "Mdee", Gender: "female",
			DOB: "1992-03-22", Nationality: "Tanzanian", Phone: "+255712000042",
			EmploymentType: "full_time", HireDate: "2021-12-01",
			DeptCode: "MKT-DIG", PositionCode: "MKT-SPE", LocationName: "Dar es Salaam HQ",
			TeamCode: "DIGITAL", ReportsToEmpID: "EMP-041",
		},
		{
			Username: "wilson.kimati", Email: "wilson.kimati@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-043", FirstName: "Wilson", LastName: "Kimati", Gender: "male",
			DOB: "1994-08-10", Nationality: "Tanzanian", Phone: "+255712000043",
			EmploymentType: "full_time", HireDate: "2022-05-01",
			DeptCode: "MKT", PositionCode: "GRA-DES", LocationName: "Dar es Salaam HQ",
			ReportsToEmpID: "EMP-041",
		},
		{
			Username: "lilian.mkongo", Email: "lilian.mkongo@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-044", FirstName: "Lilian", LastName: "Mkongo", Gender: "female",
			DOB: "1996-01-05", Nationality: "Tanzanian", Phone: "+255712000044",
			EmploymentType: "contract", HireDate: "2023-07-01",
			DeptCode: "MKT-DIG", PositionCode: "CNT-WRT", LocationName: "Remote Office",
			TeamCode: "DIGITAL", ReportsToEmpID: "EMP-041",
		},

		// ── Operations Department ──
		{
			Username: "robert.kaaya", Email: "robert.kaaya@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-050", FirstName: "Robert", LastName: "Kaaya", Gender: "male",
			DOB: "1980-05-08", Nationality: "Tanzanian", Phone: "+255712000050",
			EmploymentType: "full_time", HireDate: "2019-01-15",
			DeptCode: "OPS", PositionCode: "OPS-DIR", LocationName: "Dar es Salaam HQ",
			ReportsToEmpID: "EMP-004",
		},
		{
			Username: "mariam.mhando", Email: "mariam.mhando@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-051", FirstName: "Mariam", LastName: "Mhando", Gender: "female",
			DOB: "1986-09-15", Nationality: "Tanzanian", Phone: "+255712000051",
			EmploymentType: "full_time", HireDate: "2020-03-01",
			DeptCode: "OPS-LOG", PositionCode: "LOG-MGR", LocationName: "Dodoma Branch",
			ReportsToEmpID: "EMP-050",
		},
		{
			Username: "charles.sanga", Email: "charles.sanga@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-052", FirstName: "Charles", LastName: "Sanga", Gender: "male",
			DOB: "1991-04-28", Nationality: "Tanzanian", Phone: "+255712000052",
			EmploymentType: "full_time", HireDate: "2021-09-01",
			DeptCode: "OPS-LOG", PositionCode: "LOG-OFF", LocationName: "Arusha Branch",
			ReportsToEmpID: "EMP-051",
		},

		// ── Sales Department ──
		{
			Username: "simon.mbwilo", Email: "simon.mbwilo@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-060", FirstName: "Simon", LastName: "Mbwilo", Gender: "male",
			DOB: "1982-11-20", Nationality: "Tanzanian", Phone: "+255712000060",
			EmploymentType: "full_time", HireDate: "2019-06-01",
			DeptCode: "SALES", PositionCode: "SLS-DIR", LocationName: "Dar es Salaam HQ",
			ReportsToEmpID: "EMP-004",
		},
		{
			Username: "joyce.masanja", Email: "joyce.masanja@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-061", FirstName: "Joyce", LastName: "Masanja", Gender: "female",
			DOB: "1987-02-10", Nationality: "Tanzanian", Phone: "+255712000061",
			EmploymentType: "full_time", HireDate: "2020-08-01",
			DeptCode: "SALES", PositionCode: "SLS-MGR", LocationName: "Dar es Salaam HQ",
			TeamCode: "BD-TEAM", ReportsToEmpID: "EMP-060",
		},
		{
			Username: "elias.mwenda", Email: "elias.mwenda@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-062", FirstName: "Elias", LastName: "Mwenda", Gender: "male",
			DOB: "1993-12-18", Nationality: "Tanzanian", Phone: "+255712000062",
			EmploymentType: "full_time", HireDate: "2022-01-15",
			DeptCode: "SALES", PositionCode: "SLS-REP", LocationName: "Mwanza Branch",
			ReportsToEmpID: "EMP-061",
		},
		{
			Username: "edith.makundi", Email: "edith.makundi@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-063", FirstName: "Edith", LastName: "Makundi", Gender: "female",
			DOB: "1990-07-04", Nationality: "Tanzanian", Phone: "+255712000063",
			EmploymentType: "full_time", HireDate: "2021-03-01",
			DeptCode: "SALES-BD", PositionCode: "BD-MGR", LocationName: "Dar es Salaam HQ",
			TeamCode: "BD-TEAM", ReportsToEmpID: "EMP-060",
		},

		// ── Legal & Compliance ──
		{
			Username: "advocate.paul", Email: "advocate.paul@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-070", FirstName: "Paul", LastName: "Kiswaga", Gender: "male",
			DOB: "1977-03-30", Nationality: "Tanzanian", Phone: "+255712000070",
			EmploymentType: "full_time", HireDate: "2019-08-01",
			DeptCode: "LEGAL", PositionCode: "LEG-DIR", LocationName: "Dar es Salaam HQ",
			ReportsToEmpID: "EMP-001",
		},
		{
			Username: "hellen.shirima", Email: "hellen.shirima@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-071", FirstName: "Hellen", LastName: "Shirima", Gender: "female",
			DOB: "1986-12-15", Nationality: "Tanzanian", Phone: "+255712000071",
			EmploymentType: "full_time", HireDate: "2020-11-01",
			DeptCode: "LEGAL", PositionCode: "CMP-OFF", LocationName: "Dar es Salaam HQ",
			ReportsToEmpID: "EMP-070",
		},

		// ── QA Department ──
		{
			Username: "samuel.mgaya", Email: "samuel.mgaya@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-080", FirstName: "Samuel", LastName: "Mgaya", Gender: "male",
			DOB: "1985-10-10", Nationality: "Tanzanian", Phone: "+255712000080",
			EmploymentType: "full_time", HireDate: "2020-02-01",
			DeptCode: "QA", PositionCode: "QA-MGR", LocationName: "Dar es Salaam HQ",
			TeamCode: "QA-AUTO", ReportsToEmpID: "EMP-002",
		},
		{
			Username: "asha.ramadhani", Email: "asha.ramadhani@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-081", FirstName: "Asha", LastName: "Ramadhani", Gender: "female",
			DOB: "1992-06-25", Nationality: "Tanzanian", Phone: "+255712000081",
			EmploymentType: "full_time", HireDate: "2021-05-01",
			DeptCode: "QA", PositionCode: "QA-ENG", LocationName: "Dar es Salaam HQ",
			TeamCode: "QA-AUTO", ReportsToEmpID: "EMP-080",
		},
		{
			Username: "rashid.mwamba", Email: "rashid.mwamba@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-082", FirstName: "Rashid", LastName: "Mwamba", Gender: "male",
			DOB: "1996-03-12", Nationality: "Tanzanian", Phone: "+255712000082",
			EmploymentType: "part_time", HireDate: "2023-09-01",
			DeptCode: "QA", PositionCode: "QA-ANL", LocationName: "Remote Office",
			TeamCode: "QA-AUTO", ReportsToEmpID: "EMP-080",
		},

		// ── Intern / Probation employees (various departments) ──
		{
			Username: "kevin.massawe", Email: "kevin.massawe@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-090", FirstName: "Kevin", LastName: "Massawe", Gender: "male",
			DOB: "2000-01-20", Nationality: "Tanzanian", Phone: "+255712000090",
			EmploymentType: "intern", HireDate: "2025-09-01",
			DeptCode: "IT-DEV", PositionCode: "DEV-JR", LocationName: "Dar es Salaam HQ",
			TeamCode: "FRONTEND", ReportsToEmpID: "EMP-023",
		},
		{
			Username: "agnes.kazimoto", Email: "agnes.kazimoto@hrms.com", Password: defaultPassword, UserRole: userModels.RoleUser,
			EmployeeID: "EMP-091", FirstName: "Agnes", LastName: "Kazimoto", Gender: "female",
			DOB: "1999-11-15", Nationality: "Tanzanian", Phone: "+255712000091",
			EmploymentType: "intern", HireDate: "2025-10-01",
			DeptCode: "MKT-DIG", PositionCode: "CNT-WRT", LocationName: "Dar es Salaam HQ",
			TeamCode: "DIGITAL", ReportsToEmpID: "EMP-041",
		},
	}

	// Build a map of employee_id -> employee DB ID for ReportsToID linking
	empIDMap := make(map[string]uint) // EmployeeID code -> DB ID

	for _, def := range employees {
		// Always hash the password (needed for both new and existing users)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(def.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("  Warning: Failed to hash password for %s: %v", def.Email, err)
			continue
		}

		// Check if employee already exists
		if empRepo.ExistsByEmployeeID(def.EmployeeID) {
			existing, _ := empRepo.FindByEmployeeID(def.EmployeeID)
			if existing != nil {
				empIDMap[def.EmployeeID] = existing.ID
			}

			// Reset the user account: unblock, reset password, clear failed login count
			// This ensures test accounts are always usable after app restart
			existingUser, _ := userRepo.FindByEmail(def.Email)
			if existingUser != nil {
				needsUpdate := false
				if existingUser.Status != "active" {
					existingUser.Status = "active"
					needsUpdate = true
				}
				if !existingUser.IsActive {
					existingUser.IsActive = true
					needsUpdate = true
				}
				if existingUser.FailedLoginCount > 0 {
					existingUser.FailedLoginCount = 0
					needsUpdate = true
				}
				if existingUser.LockedUntil != nil {
					existingUser.LockedUntil = nil
					needsUpdate = true
				}
				// Always reset password to the seeder default so test credentials always work
				existingUser.Password = string(hashedPassword)
				needsUpdate = true

				if needsUpdate {
					if err := userRepo.Update(existingUser); err != nil {
						log.Printf("  Warning: Failed to reset user %s: %v", def.Email, err)
					} else {
						log.Printf("  ✅ User account reset: %s (password: %s)", def.Email, def.Password)
					}
				}
			}
			continue
		}

		// Create user first
		existingUser, _ := userRepo.FindByEmail(def.Email)
		var userID uint
		if existingUser != nil {
			// User exists but employee doesn't — reset the user and reuse
			existingUser.Password = string(hashedPassword)
			existingUser.Status = "active"
			existingUser.IsActive = true
			existingUser.FailedLoginCount = 0
			existingUser.LockedUntil = nil
			_ = userRepo.Update(existingUser)
			userID = existingUser.ID
		} else {
			user := &userModels.User{
				Username:  def.Username,
				Email:     def.Email,
				Password:  string(hashedPassword),
				FirstName: def.FirstName,
				LastName:  def.LastName,
				Role:      def.UserRole,
				Status:    "active",
				IsActive:  true,
			}
			if err := userRepo.Create(user); err != nil {
				log.Printf("  Warning: Failed to create user %s: %v", def.Email, err)
				continue
			}
			userID = user.ID
		}

		// Build the employee record
		dob, _ := time.Parse("2006-01-02", def.DOB)
		hireDate, _ := time.Parse("2006-01-02", def.HireDate)
		workEmail := def.Email

		emp := &empModels.Employee{
			UserID:         &userID,
			EmployeeID:     def.EmployeeID,
			FirstName:      def.FirstName,
			LastName:       def.LastName,
			Gender:         &def.Gender,
			DateOfBirth:    &dob,
			Nationality:    &def.Nationality,
			WorkEmail:      &workEmail,
			PhoneNumber:    &def.Phone,
			EmploymentType: &def.EmploymentType,
			HireDate:       &hireDate,
			Status:         empModels.StatusActive,
			IsActive:       true,
		}

		// Link department
		if dept, ok := departments[def.DeptCode]; ok {
			emp.DepartmentID = &dept.ID
		}

		// Link position
		if pos, ok := positions[def.PositionCode]; ok {
			emp.PositionID = &pos.ID
		}

		// Link location
		if loc, ok := locations[def.LocationName]; ok {
			emp.LocationID = &loc.ID
		}

		// Link team
		if def.TeamCode != "" {
			if team, ok := teams[def.TeamCode]; ok {
				emp.TeamID = &team.ID
			}
		}

		// Link reporting manager (from earlier employees in this batch)
		if def.ReportsToEmpID != "" {
			if managerDBID, ok := empIDMap[def.ReportsToEmpID]; ok {
				emp.ReportsToID = &managerDBID
			}
		}

		if err := empRepo.Create(emp); err != nil {
			log.Printf("  Warning: Failed to create employee %s (%s %s): %v", def.EmployeeID, def.FirstName, def.LastName, err)
			continue
		}
		empIDMap[def.EmployeeID] = emp.ID
		log.Printf("  ✅ Employee created: %s - %s %s (%s)", def.EmployeeID, def.FirstName, def.LastName, def.Email)
	}

	// Second pass: update ReportsToID for employees that reference a later-created employee
	for _, def := range employees {
		if def.ReportsToEmpID == "" {
			continue
		}
		emp, _ := empRepo.FindByEmployeeID(def.EmployeeID)
		if emp == nil || emp.ReportsToID != nil {
			continue
		}
		if managerDBID, ok := empIDMap[def.ReportsToEmpID]; ok {
			emp.ReportsToID = &managerDBID
			if err := empRepo.Update(emp); err != nil {
				log.Printf("  Warning: Failed to link %s -> %s: %v", def.EmployeeID, def.ReportsToEmpID, err)
			}
		}
	}

	log.Printf("  ✅ Test employees seeded: %d employees created/verified", len(employees))
	log.Println("  ────────────────────────────────────────────────────────────")
	log.Printf("  📋 TEST CREDENTIALS (all accounts use password: %s)", defaultPassword)
	log.Println("  ────────────────────────────────────────────────────────────")
	log.Println("  Admin:  john.mwanga@hrms.com    (CEO, Admin)")
	log.Println("  HR:     fatma.salim@hrms.com    (HR Director)")
	log.Println("  IT:     michael.temba@hrms.com  (IT Director)")
	log.Println("  User:   joseph.mosha@hrms.com   (Sr. Developer)")
	log.Println("  ────────────────────────────────────────────────────────────")
}

// ─── Leave Types ────────────────────────────────────────────────────────────

func seedLeaveTypes() {
	repo := leaveRepos.NewLeaveTypeRepository()

	types := []struct {
		Code     string
		Name     string
		Icon     string
		Category string
		Paid     bool
		DocReq   bool
		Desc     string
	}{
		{"AL", "Annual Leave", "🏖️", "annual", true, false, "Regular paid time off for rest and recreation"},
		{"SL", "Sick Leave", "🏥", "medical", true, true, "Leave for illness or medical appointments"},
		{"ML", "Maternity Leave", "👶", "family", true, true, "Leave for expecting/new mothers (84 days)"},
		{"PL", "Paternity Leave", "👨‍👦", "family", true, true, "Leave for new fathers (3 days)"},
		{"CL", "Compassionate Leave", "🙏", "compassionate", true, false, "Leave for bereavement or family emergencies"},
		{"EL", "Emergency Leave", "🚨", "emergency", true, false, "Short-notice leave for unexpected urgent matters"},
		{"UL", "Unpaid Leave", "📋", "unpaid", false, false, "Leave without pay for personal reasons"},
		{"STL", "Study Leave", "📚", "education", true, true, "Leave for exams or educational commitments"},
		{"WFH", "Work From Home", "🏠", "remote", true, false, "Working remotely from home"},
		{"HL", "Half Day Leave", "⏳", "other", true, false, "Half-day leave (morning or afternoon)"},
	}

	for _, t := range types {
		existing, _ := repo.FindByCode(t.Code)
		if existing != nil {
			continue
		}

		icon := t.Icon
		desc := t.Desc
		lt := &leaveModels.LeaveType{
			Code:                  t.Code,
			Name:                  t.Name,
			Icon:                  &icon,
			Category:              t.Category,
			PaidLeave:             t.Paid,
			RequiresDocumentation: t.DocReq,
			Description:           &desc,
			IsActive:              true,
		}
		if err := repo.Create(lt); err != nil {
			log.Printf("  Warning: Failed to create leave type %s: %v", t.Code, err)
			continue
		}
		log.Printf("  ✅ Leave type created: %s (%s)", t.Name, t.Code)
	}
}

// ─── Leave Balances ─────────────────────────────────────────────────────────

func seedLeaveBalances() {
	balanceRepo := leaveRepos.NewLeaveBalanceRepository()
	db := database.GetDB()

	year := time.Now().Year()

	// Get all active employees
	var employees []empModels.Employee
	db.Where("status = ? AND is_active = ?", "active", true).Find(&employees)

	if len(employees) == 0 {
		log.Println("  No active employees found for leave balance seeding")
		return
	}

	// Leave entitlements per type
	entitlements := map[string]float64{
		"AL":  28, // 28 days annual leave
		"SL":  21, // 21 days sick leave
		"CL":  5,  // 5 days compassionate
		"EL":  3,  // 3 days emergency
		"UL":  30, // 30 days unpaid (unlimited but capped)
		"STL": 10, // 10 days study
		"WFH": 52, // ~1 per week
		"HL":  12, // 12 half days
	}

	count := 0
	for _, emp := range employees {
		for code, entitlement := range entitlements {
			balance := &leaveModels.LeaveBalance{
				EmployeeID:    emp.EmployeeID,
				LeaveTypeCode: code,
				Year:          year,
				Entitlement:   entitlement,
				Used:          0,
				Pending:       0,
				Available:     entitlement,
			}
			if err := balanceRepo.CreateOrUpdate(balance); err != nil {
				log.Printf("  Warning: Failed to create leave balance for %s/%s: %v", emp.EmployeeID, code, err)
				continue
			}
			count++
		}
	}

	log.Printf("  ✅ Leave balances seeded: %d records for %d employees", count, len(employees))
}

// ─── Department Head Assignment ─────────────────────────────────────────────

// RunDepartmentHeads assigns department heads after all employees are created.
// This runs as a separate step so all employees exist before linking.
func RunDepartmentHeads() {
	deptRepo := deptRepos.NewDepartmentRepository()
	empRepo := empRepos.NewEmployeeRepository()

	assignments := []struct {
		DeptCode string
		EmpID    string // EmployeeID of the head
	}{
		{"EXEC", "EMP-001"},   // John Mwanga heads Executive
		{"HR", "EMP-010"},     // Fatma Salim heads HR
		{"IT", "EMP-020"},     // Michael Temba heads IT
		{"FIN", "EMP-030"},    // Daniel Lupogo heads Finance
		{"MKT", "EMP-040"},    // Gloria Massawe heads Marketing
		{"OPS", "EMP-050"},    // Robert Kaaya heads Operations
		{"SALES", "EMP-060"},  // Simon Mbwilo heads Sales
		{"LEGAL", "EMP-070"},  // Paul Kiswaga heads Legal
		{"QA", "EMP-080"},     // Samuel Mgaya heads QA
		{"IT-DEV", "EMP-021"}, // Rose Andrew heads Software Dev
		{"HR-REC", "EMP-011"}, // James Mushi heads Recruitment
	}

	for _, a := range assignments {
		dept, err := deptRepo.FindByCode(a.DeptCode)
		if err != nil || dept == nil {
			continue
		}
		if dept.ManagerID != nil {
			continue // Already assigned
		}

		emp, err := empRepo.FindByEmployeeID(a.EmpID)
		if err != nil || emp == nil {
			continue
		}

		dept.ManagerID = &emp.ID
		if err := deptRepo.Update(dept); err != nil {
			log.Printf("  Warning: Failed to assign head for %s: %v", a.DeptCode, err)
			continue
		}
		log.Printf("  ✅ Department head assigned: %s -> %s %s", a.DeptCode, emp.FirstName, emp.LastName)
	}
}

// ─── Team Lead Assignment ───────────────────────────────────────────────────

// RunTeamLeads assigns team leads after all employees are created.
func RunTeamLeads() {
	teamRepo := teamRepos.NewTeamRepository()
	empRepo := empRepos.NewEmployeeRepository()

	assignments := []struct {
		TeamCode string
		EmpID    string
	}{
		{"BACKEND", "EMP-022"},  // Joseph Mosha leads Backend
		{"FRONTEND", "EMP-023"}, // Lucy Mwakasege leads Frontend
		{"MOBILE", "EMP-027"},   // Rehema Ally leads Mobile
		{"DEVOPS-T", "EMP-028"}, // Frank Robert leads DevOps
		{"HELPDESK", "EMP-029"}, // Anna Kessy leads Helpdesk
		{"RECRUIT", "EMP-012"},  // Sarah Kimaro leads Recruitment
		{"PAYROLL", "EMP-014"},  // Happiness Charles leads Payroll
		{"DIGITAL", "EMP-041"},  // George Msemwa leads Digital Marketing
		{"ACCT", "EMP-031"},     // Esther Mwakalinga leads Accounting
		{"BD-TEAM", "EMP-061"},  // Joyce Masanja leads BD
		{"QA-AUTO", "EMP-080"},  // Samuel Mgaya leads QA Automation
	}

	for _, a := range assignments {
		team, err := teamRepo.FindByCode(a.TeamCode)
		if err != nil || team == nil {
			continue
		}
		if team.TeamLeadID != nil {
			continue // Already assigned
		}

		emp, err := empRepo.FindByEmployeeID(a.EmpID)
		if err != nil || emp == nil {
			continue
		}

		team.TeamLeadID = &emp.ID
		if err := teamRepo.Update(team); err != nil {
			log.Printf("  Warning: Failed to assign lead for %s: %v", a.TeamCode, err)
			continue
		}
		log.Printf("  ✅ Team lead assigned: %s -> %s %s", a.TeamCode, emp.FirstName, emp.LastName)
	}
}

// ─── Leave Policies ───────────────────────────────────────────────────────

// seedLeavePolicies seeds default leave policies for Tanzania
func seedLeavePolicies() {
	repo := leaveRepos.NewLeavePolicyRepository()

	policies := []struct {
		PolicyName             string
		Country                string
		LeaveTypeCode          string
		Entitlement            int
		AccrualFrequency       string
		ProrationOnJoin        bool
		ProrationOnExit        bool
		CarryForward           bool
		CarryForwardLimit      *int
		CarryForwardExpiry     *time.Time
		EncashmentAllowed      bool
		EncashmentLimit        *int
		NegativeBalanceAllowed bool
		SandwichRules          bool
		HalfDayAllowed         bool
		MinimumNoticeDays      int
		MaximumDaysPerRequest  *int
		IsActive               bool
	}{
		// Annual Leave Policy - 28 days, starts after 3-month probation
		{
			PolicyName:             "Tanzania Annual Leave Policy",
			Country:                "Tanzania",
			LeaveTypeCode:          "AL",
			Entitlement:            28,
			AccrualFrequency:       "Annual",
			ProrationOnJoin:        true,
			ProrationOnExit:        true,
			CarryForward:           true,
			CarryForwardLimit:      intPtr(14), // Max 14 days can be carried forward
			CarryForwardExpiry:     timePtr(time.Date(time.Now().Year(), 12, 31, 0, 0, 0, 0, time.UTC)),
			EncashmentAllowed:      true,
			EncashmentLimit:        intPtr(7), // Max 7 days can be encashed
			NegativeBalanceAllowed: false,
			SandwichRules:          true,
			HalfDayAllowed:         true,
			MinimumNoticeDays:      3,
			MaximumDaysPerRequest:  intPtr(14), // Max 14 days per request
			IsActive:               true,
		},
		// Maternity Leave Policy - 84-100 days paid
		{
			PolicyName:             "Tanzania Maternity Leave Policy",
			Country:                "Tanzania",
			LeaveTypeCode:          "ML",
			Entitlement:            100, // Up to 100 days as requested
			AccrualFrequency:       "Per Pregnancy",
			ProrationOnJoin:        false,
			ProrationOnExit:        false,
			CarryForward:           false,
			CarryForwardLimit:      nil,
			CarryForwardExpiry:     nil,
			EncashmentAllowed:      false,
			EncashmentLimit:        nil,
			NegativeBalanceAllowed: false,
			SandwichRules:          false,
			HalfDayAllowed:         false,
			MinimumNoticeDays:      30, // 30 days notice for maternity leave
			MaximumDaysPerRequest:  intPtr(100),
			IsActive:               true,
		},
		// Paternity Leave Policy - 3 days
		{
			PolicyName:             "Tanzania Paternity Leave Policy",
			Country:                "Tanzania",
			LeaveTypeCode:          "PL",
			Entitlement:            3,
			AccrualFrequency:       "Per Birth Event",
			ProrationOnJoin:        false,
			ProrationOnExit:        false,
			CarryForward:           false,
			CarryForwardLimit:      nil,
			CarryForwardExpiry:     nil,
			EncashmentAllowed:      false,
			EncashmentLimit:        nil,
			NegativeBalanceAllowed: false,
			SandwichRules:          false,
			HalfDayAllowed:         false,
			MinimumNoticeDays:      7,
			MaximumDaysPerRequest:  intPtr(3),
			IsActive:               true,
		},
		// Compassionate Leave Policy - 4 days
		{
			PolicyName:             "Tanzania Compassionate Leave Policy",
			Country:                "Tanzania",
			LeaveTypeCode:          "CL",
			Entitlement:            4,
			AccrualFrequency:       "Annual",
			ProrationOnJoin:        true,
			ProrationOnExit:        false,
			CarryForward:           false,
			CarryForwardLimit:      nil,
			CarryForwardExpiry:     nil,
			EncashmentAllowed:      false,
			EncashmentLimit:        nil,
			NegativeBalanceAllowed: false,
			SandwichRules:          false,
			HalfDayAllowed:         false,
			MinimumNoticeDays:      0, // Immediate leave for emergencies
			MaximumDaysPerRequest:  intPtr(4),
			IsActive:               true,
		},
		// Emergency Leave Policy - Available during probation
		{
			PolicyName:             "Tanzania Emergency Leave Policy",
			Country:                "Tanzania",
			LeaveTypeCode:          "EL",
			Entitlement:            5, // 5 days emergency leave per year
			AccrualFrequency:       "Annual",
			ProrationOnJoin:        true,
			ProrationOnExit:        true,
			CarryForward:           false,
			CarryForwardLimit:      nil,
			CarryForwardExpiry:     nil,
			EncashmentAllowed:      false,
			EncashmentLimit:        nil,
			NegativeBalanceAllowed: true, // Allow negative balance for emergencies
			SandwichRules:          false,
			HalfDayAllowed:         true,
			MinimumNoticeDays:      0,         // No notice required for emergencies
			MaximumDaysPerRequest:  intPtr(3), // Max 3 days per emergency request
			IsActive:               true,
		},
	}

	for _, p := range policies {
		existing, _ := repo.FindByCountryAndLeaveType(p.Country, p.LeaveTypeCode, nil)
		if existing != nil {
			continue
		}

		policy := &leaveModels.LeavePolicy{
			PolicyName:             p.PolicyName,
			Country:                p.Country,
			LeaveTypeCode:          p.LeaveTypeCode,
			Entitlement:            p.Entitlement,
			AccrualFrequency:       p.AccrualFrequency,
			ProrationOnJoin:        p.ProrationOnJoin,
			ProrationOnExit:        p.ProrationOnExit,
			CarryForward:           p.CarryForward,
			CarryForwardLimit:      p.CarryForwardLimit,
			CarryForwardExpiry:     p.CarryForwardExpiry,
			EncashmentAllowed:      p.EncashmentAllowed,
			EncashmentLimit:        p.EncashmentLimit,
			NegativeBalanceAllowed: p.NegativeBalanceAllowed,
			SandwichRules:          p.SandwichRules,
			HalfDayAllowed:         p.HalfDayAllowed,
			MinimumNoticeDays:      p.MinimumNoticeDays,
			MaximumDaysPerRequest:  p.MaximumDaysPerRequest,
			IsActive:               p.IsActive,
		}

		if err := repo.Create(policy); err != nil {
			log.Printf("  Warning: Failed to create leave policy %s: %v", p.PolicyName, err)
			continue
		}
		log.Printf("  ✅ Leave policy created: %s", p.PolicyName)
	}
}

// Helper functions for pointer types
func intPtr(i int) *int {
	return &i
}

func timePtr(t time.Time) *time.Time {
	return &t
}
