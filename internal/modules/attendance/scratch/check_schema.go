//go:build ignore
// +build ignore

package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	roleModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/models"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	organizationModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization/models"
	organizationUnitModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization_units/models"
	departmentModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/models"
	teamModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/teams/models"
	positionModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/models"
	locationModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/models"
	costCenterModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/cost_centers/models"
	assetModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/models"
	employeeModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	helpdeskModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
	biometricModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/models"
	shiftModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/models"
	leaveModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	auditModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/audit/models"
	payrollModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	dashboardModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/dashboard/models"
	attendanceModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	projectModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/projects/models"
	performanceModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	settingsModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/settings/models"
	recruitmentModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/models"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=scoop password=my_scoop_@_2026 dbname=scoop_db_v001 port=5432 sslmode=disable TimeZone=UTC"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	modelsToTest := []interface{}{
		&roleModels.Role{},
		&roleModels.Permission{},
		&roleModels.UserRole{},
		&roleModels.RolePermission{},
		&userModels.User{},
		&organizationModels.Organization{},
		&organizationUnitModels.OrganizationUnit{},
		&departmentModels.Department{},
		&teamModels.Team{},
		&positionModels.JobPosition{},
		&locationModels.Location{},
		&costCenterModels.CostCenter{},
		&assetModels.Asset{},
		&employeeModels.Employee{},
		&employeeModels.EmployeeOnboardingDraft{},
		&employeeModels.EmployeeBasicInformation{},
		&employeeModels.EmployeeEmploymentDetails{},
		&employeeModels.EmployeeAddress{},
		&employeeModels.EmployeeSalaryComponent{},
		&employeeModels.EmployeeBankAccount{},
		&employeeModels.EmployeeStatutoryInfo{},
		&employeeModels.EmployeeDocument{},
		&employeeModels.EmployeeAsset{},
		&employeeModels.EmployeeEmergencyContact{},
		&employeeModels.EmployeePolicy{},
		&employeeModels.PostOnboardingTask{},
		&employeeModels.OffboardingWorkflow{},
		&employeeModels.OffboardingClearance{},
		&employeeModels.OffboardingAssetReturn{},
		&employeeModels.FinalSettlement{},
		&employeeModels.ServiceRequest{},
		&employeeModels.ProfileUpdateRequest{},
		&assetModels.AssetRequest{},
		&assetModels.AssetIssue{},
		&helpdeskModels.Ticket{},
		&helpdeskModels.Comment{},
		&helpdeskModels.Attachment{},
		&helpdeskModels.RoutingRule{},
		&helpdeskModels.KnowledgeBaseArticle{},
		&helpdeskModels.KBArticleFeedback{},
		&helpdeskModels.TicketCategory{},
		&biometricModels.BioTimeConfig{},
		&biometricModels.BioTimeTransaction{},
		&biometricModels.BiometricEnrollment{},
		&shiftModels.Shift{},
		&shiftModels.ShiftLocation{}, // Join table for shifts and locations,
		&shiftModels.RosterAssignment{},
		&shiftModels.SwapRequest{},
		&shiftModels.RosterChangeRequest{},
		&leaveModels.LeaveType{},
		&leaveModels.LeavePolicy{},
		&leaveModels.LeaveRequest{},
		&leaveModels.LeaveApproval{},
		&leaveModels.LeaveDocument{},
		&leaveModels.LeaveBalance{},
		&leaveModels.Holiday{},
		&auditModels.AuditLog{},
		&payrollModels.PayrollRun{},
		&payrollModels.PayrollRunEmployee{},
		&payrollModels.SalaryStructure{},
		&payrollModels.SalaryComponent{},
		&payrollModels.Payslip{},
		&payrollModels.PayslipItem{},
		&payrollModels.Loan{},
		&payrollModels.LoanRepayment{},
		&payrollModels.TaxSlab{},
		&payrollModels.StatutoryRule{},
		&payrollModels.NHIFSchedule{},
		&payrollModels.CompliancePayment{},
		&dashboardModels.Announcement{},
		&dashboardModels.AnnouncementAttachment{},
		&dashboardModels.AnnouncementRead{},
		&attendanceModels.TimesheetWeek{},
		&attendanceModels.TimesheetEntry{},
		&attendanceModels.TimesheetApproval{},
		&attendanceModels.OvertimePolicy{},
		&attendanceModels.OvertimeRequest{},
		&attendanceModels.OvertimeApproval{},
		&projectModels.Project{},
		&projectModels.ProjectMember{},
		&projectModels.DailyTask{},
		&performanceModels.Goal{},
		&performanceModels.KeyResult{},
		&performanceModels.GoalCheckIn{},
		&performanceModels.DepartmentTarget{},
		&performanceModels.DepartmentTargetMilestone{},
		&performanceModels.EmployeeTarget{},
		&performanceModels.AppraisalCycle{},
		&performanceModels.AppraisalWorkflowStep{},
		&performanceModels.Appraisal{},
		&performanceModels.Feedback360Campaign{},
		&performanceModels.Feedback360RaterGroup{},
		&performanceModels.TalentReview{},
		&performanceModels.CalibrationSession{},
		&performanceModels.SuccessionPlan{},
		&settingsModels.MenuVisibilitySetting{},
		&recruitmentModels.JobRequisition{},
		&recruitmentModels.JobOpening{},
		&recruitmentModels.Candidate{},
		&recruitmentModels.JobApplication{},
		&recruitmentModels.Interview{},
		&recruitmentModels.Offer{},
		&recruitmentModels.TalentPoolCandidate{},
		&attendanceModels.ManualPunch{},
	}

	for i, model := range modelsToTest {
		log.Printf("Migrating model index %d (%T)...", i, model)
		err = db.AutoMigrate(model)
		if err != nil {
			log.Fatalf("Migration failed for %T: %v", model, err)
		}
	}
	log.Println("All models migrated successfully")
}
