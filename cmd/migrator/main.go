package main

import (
	"log"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/app"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	assetModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/models"
	attendanceModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	auditModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/audit/models"
	biometricModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/models"
	costCenterModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/cost_centers/models"
	dashboardModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/dashboard/models"
	departmentModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/models"
	employeeModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	helpdeskModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
	leaveModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	locationModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/models"
	organizationUnitModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization_units/models"
	organizationModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organizations/models"
	payrollModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	performanceModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	positionModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/models"
	projectModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/projects/models"
	recruitmentModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/models"
	roleModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/models"
	settingsModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/settings/models"
	shiftModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/models"
	teamModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/teams/models"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
)

func main() {
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	modelsToMigrate := []interface{}{
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
		&assetModels.AssetAssignmentHistory{},
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
		&shiftModels.ShiftLocation{},
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
		&recruitmentModels.ApprovalStep{}, // child table for JobRequisition.approvalFlow (required for GET / approve)
		&recruitmentModels.JobOpening{},
		&recruitmentModels.Candidate{},
		&recruitmentModels.JobApplication{},
		&recruitmentModels.ApplicationSubmissionQueue{},
		&recruitmentModels.Interview{},
		&recruitmentModels.Feedback{},
		&recruitmentModels.InterviewWorkflowDefinition{},
		&recruitmentModels.InterviewWorkflowDefinitionStage{},
		&recruitmentModels.InterviewWorkflowProcess{},
		&recruitmentModels.InterviewWorkflowProcessStageSnapshot{},
		&recruitmentModels.InterviewWorkflowStageAttempt{},
		&recruitmentModels.InterviewWorkflowFinalApproval{},
		&recruitmentModels.InterviewWorkflowAuditEvent{},
		&recruitmentModels.Offer{},
		&recruitmentModels.TalentPoolCandidate{},
		&attendanceModels.ManualPunch{},
	}

	// AutoMigrate all registered models: creates missing tables AND adds new columns
	// (e.g. job_openings.apply_token) on existing databases. The previous "new tables only"
	// mode never altered existing tables, which caused schema drift.
	log.Printf("Syncing schema for %d model(s) (tables + columns)...", len(modelsToMigrate))
	if err := database.Migrate(modelsToMigrate...); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully.")
}
