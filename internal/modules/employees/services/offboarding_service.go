package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	assetRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/repositories"
	departmentRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
	employeeModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"gorm.io/gorm"
)

// OffboardingService handles offboarding business logic
type OffboardingService struct {
	workflowRepo    *employeeRepos.OffboardingWorkflowRepository
	clearanceRepo   *employeeRepos.OffboardingClearanceRepository
	assetReturnRepo *employeeRepos.OffboardingAssetReturnRepository
	settlementRepo  *employeeRepos.FinalSettlementRepository
	employeeRepo    *employeeRepos.EmployeeRepository
	assetRepo        *assetRepos.AssetRepository
	departmentRepo   *departmentRepos.DepartmentRepository
}

// NewOffboardingService creates a new offboarding service
func NewOffboardingService() *OffboardingService {
	return &OffboardingService{
		workflowRepo:    employeeRepos.NewOffboardingWorkflowRepository(),
		clearanceRepo:   employeeRepos.NewOffboardingClearanceRepository(),
		assetReturnRepo: employeeRepos.NewOffboardingAssetReturnRepository(),
		settlementRepo:  employeeRepos.NewFinalSettlementRepository(),
		employeeRepo:    employeeRepos.NewEmployeeRepository(),
		assetRepo:        assetRepos.NewAssetRepository(),
		departmentRepo:   departmentRepos.NewDepartmentRepository(),
	}
}

// InitiateSeparation initiates a new offboarding workflow for an employee
func (s *OffboardingService) InitiateSeparation(req *employeeModels.InitiateSeparationRequest, tenantID *uint, createdBy *uint) (*employeeModels.OffboardingWorkflow, error) {
	// Check if employee exists
	employee, err := s.employeeRepo.FindByEmployeeID(req.EmployeeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("employee not found")
		}
		return nil, err
	}

	// Check if there's already an active offboarding workflow for this employee
	existing, _ := s.workflowRepo.FindByEmployeeID(req.EmployeeID, tenantID)
	if existing != nil && existing.Status != "Completed" {
		return nil, errors.New("an active offboarding workflow already exists for this employee")
	}

	// Parse dates
	resignationDate, err := time.Parse("2006-01-02T15:04:05Z07:00", req.ResignationDate)
	if err != nil {
		// Try alternative format
		resignationDate, err = time.Parse("2006-01-02", req.ResignationDate)
		if err != nil {
			return nil, errors.New("invalid resignation_date format, expected ISO 8601")
		}
	}

	lastWorkingDate, err := time.Parse("2006-01-02T15:04:05Z07:00", req.LastWorkingDate)
	if err != nil {
		// Try alternative format
		lastWorkingDate, err = time.Parse("2006-01-02", req.LastWorkingDate)
		if err != nil {
			return nil, errors.New("invalid last_working_date format, expected ISO 8601")
		}
	}

	// Validate dates
	if lastWorkingDate.Before(resignationDate) {
		return nil, errors.New("last working date must be after resignation date")
	}

	// Calculate served notice period
	servedNoticePeriod := int(lastWorkingDate.Sub(resignationDate).Hours() / 24)

	// Generate offboarding ID
	offboardingID, err := s.workflowRepo.GenerateOffboardingID()
	if err != nil {
		return nil, err
	}

	// Set exit interview required (default: true)
	exitInterviewRequired := true
	if req.ExitInterviewRequired != nil {
		exitInterviewRequired = *req.ExitInterviewRequired
	}

	// Create workflow
	workflow := &employeeModels.OffboardingWorkflow{
		OffboardingID:         offboardingID,
		TenantID:              tenantID,
		EmployeeID:            req.EmployeeID,
		EmployeeDBID:          &employee.ID,
		SeparationType:        req.SeparationType,
		ResignationDate:       &resignationDate,
		LastWorkingDate:       &lastWorkingDate,
		NoticePeriodDays:      &req.NoticePeriodDays,
		ServedNoticePeriod:    &servedNoticePeriod,
		Reason:                &req.Reason,
		ReasonCode:            req.ReasonCode,
		ExitInterviewRequired: &exitInterviewRequired,
		ExitInterviewStatus:  "Not Scheduled",
		ClearanceStatus:       "Pending",
		FinalSettlementStatus: "Pending",
		Status:                "In Progress",
		ClearanceProgress:     new(float64), // 0
		CompletedClearances:   new(int),      // 0
		TotalClearances:       intPtr(5),     // 5 default clearances
		CreatedBy:             createdBy,
		UpdatedBy:             createdBy,
	}

	if req.InitiatedBy != nil {
		workflow.InitiatedByID = req.InitiatedBy
	}
	if req.AdditionalNotes != nil {
		workflow.AdditionalNotes = req.AdditionalNotes
	}

	if err := s.workflowRepo.Create(workflow); err != nil {
		return nil, fmt.Errorf("failed to create offboarding workflow: %w", err)
	}

	// Create default clearances
	clearances := []employeeModels.OffboardingClearance{
		{OffboardingID: offboardingID, TenantID: tenantID, Department: "Manager", Status: "Pending"},
		{OffboardingID: offboardingID, TenantID: tenantID, Department: "IT", Status: "Pending"},
		{OffboardingID: offboardingID, TenantID: tenantID, Department: "HR", Status: "Pending"},
		{OffboardingID: offboardingID, TenantID: tenantID, Department: "Finance", Status: "Pending"},
		{OffboardingID: offboardingID, TenantID: tenantID, Department: "Assets", Status: "Pending"},
	}

	// Generate unique clearance IDs for each clearance
	// Get the base count first, then increment for each clearance
	baseCount, err := s.clearanceRepo.GetClearanceCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get clearance count: %w", err)
	}

	for i := range clearances {
		clearances[i].ClearanceID = fmt.Sprintf("clear-%03d", baseCount+int64(i)+1)
	}

	if err := s.clearanceRepo.CreateBatch(clearances); err != nil {
		return nil, fmt.Errorf("failed to create clearances: %w", err)
	}

	// Create asset return records for all assigned assets
	if err := s.createAssetReturnRecords(offboardingID, req.EmployeeID, tenantID, &lastWorkingDate); err != nil {
		// Log error but don't fail the workflow creation
		fmt.Printf("Warning: Failed to create asset return records: %v\n", err)
	}

	return workflow, nil
}

// createAssetReturnRecords creates asset return records for all assets assigned to the employee
func (s *OffboardingService) createAssetReturnRecords(offboardingID, employeeID string, tenantID *uint, expectedReturnDate *time.Time) error {
	// Find all assets assigned to this employee using employee_id string
	filters := map[string]interface{}{
		"assigned_to": employeeID, // Use employee_id string (e.g., "EMP001")
	}
	assets, _, err := s.assetRepo.List(tenantID, 1, 1000, filters)
	if err != nil {
		return err
	}

	// Create asset return records
	for _, asset := range assets {
		assetReturn := &employeeModels.OffboardingAssetReturn{
			OffboardingID:      offboardingID,
			AssetID:           asset.ID,
			TenantID:          tenantID,
			ReturnStatus:       "Pending",
			ExpectedReturnDate: expectedReturnDate,
		}
		if err := s.assetReturnRepo.Create(assetReturn); err != nil {
			// Log but continue
			fmt.Printf("Warning: Failed to create asset return record for asset %d: %v\n", asset.ID, err)
		}
	}

	return nil
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}

// ListWorkflows retrieves a paginated list of offboarding workflows
func (s *OffboardingService) ListWorkflows(req *employeeModels.ListOffboardingWorkflowsRequest, tenantID *uint) ([]employeeModels.OffboardingWorkflow, int64, error) {
	return s.workflowRepo.List(req, tenantID)
}

// GetWorkflowByID retrieves an offboarding workflow by offboarding_id
func (s *OffboardingService) GetWorkflowByID(offboardingID string, tenantID *uint) (*employeeModels.OffboardingWorkflow, error) {
	return s.workflowRepo.FindByOffboardingID(offboardingID, tenantID)
}

// UpdateWorkflow updates an offboarding workflow
func (s *OffboardingService) UpdateWorkflow(offboardingID string, req *employeeModels.UpdateOffboardingWorkflowRequest, tenantID *uint, updatedBy *uint) (*employeeModels.OffboardingWorkflow, error) {
	workflow, err := s.workflowRepo.FindByOffboardingID(offboardingID, tenantID)
	if err != nil {
		return nil, errors.New("offboarding workflow not found")
	}

	// Update fields
	if req.LastWorkingDate != nil && *req.LastWorkingDate != "" {
		lastWorkingDate, err := time.Parse("2006-01-02T15:04:05Z07:00", *req.LastWorkingDate)
		if err != nil {
			lastWorkingDate, err = time.Parse("2006-01-02", *req.LastWorkingDate)
			if err != nil {
				return nil, errors.New("invalid last_working_date format")
			}
		}
		workflow.LastWorkingDate = &lastWorkingDate
		// Recalculate served notice period
		if workflow.ResignationDate != nil {
			servedNoticePeriod := int(lastWorkingDate.Sub(*workflow.ResignationDate).Hours() / 24)
			workflow.ServedNoticePeriod = &servedNoticePeriod
		}
	}
	if req.Reason != nil {
		workflow.Reason = req.Reason
	}
	if req.AdditionalNotes != nil {
		workflow.AdditionalNotes = req.AdditionalNotes
	}
	if req.Status != nil {
		workflow.Status = *req.Status
	}

	workflow.UpdatedBy = updatedBy
	if err := s.workflowRepo.Update(workflow); err != nil {
		return nil, err
	}

	return workflow, nil
}

// GetClearances retrieves all clearances for an offboarding workflow
func (s *OffboardingService) GetClearances(offboardingID string, tenantID *uint) ([]employeeModels.OffboardingClearance, error) {
	return s.clearanceRepo.FindByOffboardingID(offboardingID, tenantID)
}

// UpdateClearance updates a clearance status
func (s *OffboardingService) UpdateClearance(offboardingID, clearanceID string, req *employeeModels.UpdateClearanceRequest, tenantID *uint) (*employeeModels.OffboardingClearance, error) {
	clearance, err := s.clearanceRepo.FindByClearanceID(clearanceID, tenantID)
	if err != nil {
		return nil, errors.New("clearance not found")
	}

	// Verify it belongs to the offboarding workflow
	if clearance.OffboardingID != offboardingID {
		return nil, errors.New("clearance does not belong to this offboarding workflow")
	}

	// Validate status
	if req.Status == "Cleared" && (req.ClearedByID == nil || *req.ClearedByID == "") {
		return nil, errors.New("cleared_by_id is required when status is Cleared")
	}
	if req.Status == "Issues" && (req.Issues == nil || *req.Issues == "") {
		return nil, errors.New("issues description is required when status is Issues")
	}

	// Update clearance
	clearance.Status = req.Status
	if req.Status == "Cleared" {
		now := time.Now()
		clearance.ClearanceDate = &now
		if req.ClearedByID != nil {
			clearance.ClearedByID = req.ClearedByID
		}
	}
	if req.Notes != nil {
		clearance.Notes = req.Notes
	}
	if req.Issues != nil {
		clearance.Issues = req.Issues
	}

	if err := s.clearanceRepo.Update(clearance); err != nil {
		return nil, err
	}

	// Update workflow clearance status and progress
	if err := s.updateWorkflowClearanceStatus(offboardingID, tenantID); err != nil {
		fmt.Printf("Warning: Failed to update workflow clearance status: %v\n", err)
	}

	return clearance, nil
}

// updateWorkflowClearanceStatus updates the clearance status and progress for a workflow
func (s *OffboardingService) updateWorkflowClearanceStatus(offboardingID string, tenantID *uint) error {
	clearances, err := s.clearanceRepo.FindByOffboardingID(offboardingID, tenantID)
	if err != nil {
		return err
	}

	workflow, err := s.workflowRepo.FindByOffboardingID(offboardingID, tenantID)
	if err != nil {
		return err
	}

	// Count clearances
	totalClearances := len(clearances)
	completedClearances := 0
	hasIssues := false
	hasInProgress := false

	for _, clearance := range clearances {
		if clearance.Status == "Cleared" {
			completedClearances++
		} else if clearance.Status == "Issues" {
			hasIssues = true
		} else if clearance.Status == "In Progress" {
			hasInProgress = true
		}
	}

	// Calculate progress
	var progress float64
	if totalClearances > 0 {
		progress = float64(completedClearances) / float64(totalClearances) * 100
	}

	// Update clearance status
	var clearanceStatus string
	if completedClearances == totalClearances {
		clearanceStatus = "Completed"
	} else if hasIssues {
		clearanceStatus = "Issues"
	} else if hasInProgress || completedClearances > 0 {
		clearanceStatus = "In Progress"
	} else {
		clearanceStatus = "Pending"
	}

	// Update workflow
	workflow.ClearanceStatus = clearanceStatus
	workflow.ClearanceProgress = &progress
	completed := int(completedClearances)
	total := int(totalClearances)
	workflow.CompletedClearances = &completed
	workflow.TotalClearances = &total

	return s.workflowRepo.Update(workflow)
}

// GetAssetReturns retrieves all asset returns for an offboarding workflow
func (s *OffboardingService) GetAssetReturns(offboardingID string, tenantID *uint) ([]employeeModels.OffboardingAssetReturn, error) {
	return s.assetReturnRepo.FindByOffboardingID(offboardingID, tenantID)
}

// RecordAssetReturn records the return of an asset
func (s *OffboardingService) RecordAssetReturn(offboardingID string, assetID uint, req *employeeModels.RecordAssetReturnRequest, tenantID *uint) (*employeeModels.OffboardingAssetReturn, error) {
	// Find asset return record
	assetReturn, err := s.assetReturnRepo.FindByAssetIDAndOffboardingID(assetID, offboardingID, tenantID)
	if err != nil {
		return nil, errors.New("asset return record not found")
	}

	// Parse return date
	returnDate, err := time.Parse("2006-01-02T15:04:05Z07:00", req.ReturnDate)
	if err != nil {
		returnDate, err = time.Parse("2006-01-02", req.ReturnDate)
		if err != nil {
			return nil, errors.New("invalid return_date format")
		}
	}

	// Update asset return
	assetReturn.ReturnStatus = "Returned"
	assetReturn.ReturnDate = &returnDate
	assetReturn.ReturnCondition = &req.ReturnCondition
	if req.ReturnNotes != nil {
		assetReturn.ReturnNotes = req.ReturnNotes
	}
	assetReturn.ReturnedByID = &req.ReturnedByID

	if err := s.assetReturnRepo.Update(assetReturn); err != nil {
		return nil, err
	}

	// Update asset status to available
	asset, err := s.assetRepo.FindByID(assetID)
	if err == nil && asset != nil {
		asset.Status = "available"
		asset.EmployeeID = nil
		asset.AssignedTo = nil
		asset.ReturnDate = &returnDate
		s.assetRepo.Update(asset)
	}

	// Update Assets clearance status
	if err := s.updateAssetsClearanceStatus(offboardingID, tenantID); err != nil {
		fmt.Printf("Warning: Failed to update assets clearance status: %v\n", err)
	}

	return assetReturn, nil
}

// RecordAssetIssue records an issue with an asset
func (s *OffboardingService) RecordAssetIssue(offboardingID string, assetID uint, req *employeeModels.RecordAssetIssueRequest, tenantID *uint) (*employeeModels.OffboardingAssetReturn, error) {
	assetReturn, err := s.assetReturnRepo.FindByAssetIDAndOffboardingID(assetID, offboardingID, tenantID)
	if err != nil {
		return nil, errors.New("asset return record not found")
	}

	assetReturn.ReturnStatus = "Issue"
	assetReturn.IssueType = &req.IssueType
	assetReturn.IssueDescription = &req.IssueDescription
	assetReturn.ReportedByID = &req.ReportedByID

	if err := s.assetReturnRepo.Update(assetReturn); err != nil {
		return nil, err
	}

	// Update Assets clearance status to Issues
	if err := s.updateAssetsClearanceStatus(offboardingID, tenantID); err != nil {
		fmt.Printf("Warning: Failed to update assets clearance status: %v\n", err)
	}

	return assetReturn, nil
}

// updateAssetsClearanceStatus updates the Assets clearance status based on asset returns
func (s *OffboardingService) updateAssetsClearanceStatus(offboardingID string, tenantID *uint) error {
	assetReturns, err := s.assetReturnRepo.FindByOffboardingID(offboardingID, tenantID)
	if err != nil {
		return err
	}

	// Find Assets clearance
	clearances, err := s.clearanceRepo.FindByOffboardingID(offboardingID, tenantID)
	if err != nil {
		return err
	}

	var assetsClearance *employeeModels.OffboardingClearance
	for i := range clearances {
		if clearances[i].Department == "Assets" {
			assetsClearance = &clearances[i]
			break
		}
	}

	if assetsClearance == nil {
		return nil // No Assets clearance found
	}

	// Check asset return statuses
	allReturned := true
	hasIssues := false

	for _, assetReturn := range assetReturns {
		if assetReturn.ReturnStatus == "Issue" {
			hasIssues = true
			allReturned = false
			break
		}
		if assetReturn.ReturnStatus != "Returned" {
			allReturned = false
		}
	}

	// Update clearance status
	if hasIssues {
		assetsClearance.Status = "Issues"
	} else if allReturned && len(assetReturns) > 0 {
		assetsClearance.Status = "Cleared"
		now := time.Now()
		assetsClearance.ClearanceDate = &now
	} else if len(assetReturns) > 0 {
		assetsClearance.Status = "In Progress"
	}

	if err := s.clearanceRepo.Update(assetsClearance); err != nil {
		return err
	}

	// Update workflow clearance status
	return s.updateWorkflowClearanceStatus(offboardingID, tenantID)
}

// ScheduleExitInterview schedules an exit interview
func (s *OffboardingService) ScheduleExitInterview(offboardingID string, req *employeeModels.ScheduleExitInterviewRequest, tenantID *uint) (*employeeModels.OffboardingWorkflow, error) {
	workflow, err := s.workflowRepo.FindByOffboardingID(offboardingID, tenantID)
	if err != nil {
		return nil, errors.New("offboarding workflow not found")
	}

	// Parse interview date
	interviewDate, err := time.Parse("2006-01-02T15:04:05Z07:00", req.InterviewDate)
	if err != nil {
		interviewDate, err = time.Parse("2006-01-02T15:04:05", req.InterviewDate)
		if err != nil {
			return nil, errors.New("invalid interview_date format")
		}
	}

	workflow.ExitInterviewStatus = "Scheduled"
	workflow.ExitInterviewDate = &interviewDate
	workflow.ExitInterviewerID = &req.InterviewerID
	if req.Location != nil {
		workflow.ExitInterviewLocation = req.Location
	}
	if req.Notes != nil {
		// Store notes temporarily (we'll use a separate field or notes)
		workflow.ExitInterviewNotes = req.Notes
	}

	if err := s.workflowRepo.Update(workflow); err != nil {
		return nil, err
	}

	return workflow, nil
}

// CompleteExitInterview completes an exit interview
func (s *OffboardingService) CompleteExitInterview(offboardingID string, req *employeeModels.CompleteExitInterviewRequest, tenantID *uint) (*employeeModels.OffboardingWorkflow, error) {
	workflow, err := s.workflowRepo.FindByOffboardingID(offboardingID, tenantID)
	if err != nil {
		return nil, errors.New("offboarding workflow not found")
	}

	now := time.Now()
	workflow.ExitInterviewStatus = "Completed"
	workflow.ExitInterviewNotes = &req.InterviewNotes
	workflow.ExitInterviewConductedByID = &req.ConductedByID
	workflow.ExitInterviewCompletedAt = &now
	if req.FeedbackRating != nil {
		workflow.ExitInterviewRating = req.FeedbackRating
	}
	if req.WouldRecommend != nil {
		workflow.WouldRecommend = req.WouldRecommend
	}

	if err := s.workflowRepo.Update(workflow); err != nil {
		return nil, err
	}

	return workflow, nil
}

// CancelExitInterview cancels an exit interview
func (s *OffboardingService) CancelExitInterview(offboardingID string, req *employeeModels.CancelExitInterviewRequest, tenantID *uint) (*employeeModels.OffboardingWorkflow, error) {
	workflow, err := s.workflowRepo.FindByOffboardingID(offboardingID, tenantID)
	if err != nil {
		return nil, errors.New("offboarding workflow not found")
	}

	workflow.ExitInterviewStatus = "Not Scheduled"
	workflow.ExitInterviewDate = nil
	workflow.ExitInterviewLocation = nil
	workflow.ExitInterviewerID = nil
	if req.CancellationReason != "" {
		notes := fmt.Sprintf("Cancelled: %s", req.CancellationReason)
		workflow.ExitInterviewNotes = &notes
	}

	if err := s.workflowRepo.Update(workflow); err != nil {
		return nil, err
	}

	return workflow, nil
}

// CalculateSettlement calculates the final settlement
func (s *OffboardingService) CalculateSettlement(offboardingID string, req *employeeModels.CalculateSettlementRequest, tenantID *uint, calculatedByID *string) (*employeeModels.FinalSettlement, error) {
	workflow, err := s.workflowRepo.FindByOffboardingID(offboardingID, tenantID)
	if err != nil {
		return nil, errors.New("offboarding workflow not found")
	}

	// Check if settlement already exists
	existing, _ := s.settlementRepo.FindByOffboardingID(offboardingID, tenantID)
	if existing != nil && existing.Status != "Pending" {
		return nil, errors.New("settlement has already been calculated")
	}

	// Build settlement breakdown
	breakdown := employeeModels.SettlementBreakdown{
		Currency: "TZS",
	}

	// Use provided values or calculate defaults
	if req.OutstandingSalary != nil {
		breakdown.Earnings.OutstandingSalary = *req.OutstandingSalary
	}
	if req.LeaveEncashment != nil {
		breakdown.Earnings.LeaveEncashment = *req.LeaveEncashment
	}
	if req.Bonus != nil {
		breakdown.Earnings.Bonus = *req.Bonus
	}
	if req.Incentives != nil {
		breakdown.Earnings.Incentives = *req.Incentives
	}
	if req.OtherEarnings != nil {
		breakdown.Earnings.OtherEarnings = *req.OtherEarnings
	}
	if req.OutstandingLoans != nil {
		breakdown.Deductions.OutstandingLoans = *req.OutstandingLoans
	}
	if req.Advances != nil {
		breakdown.Deductions.Advances = *req.Advances
	}
	if req.AssetDamages != nil {
		breakdown.Deductions.AssetDamages = *req.AssetDamages
	}
	if req.TaxDeductions != nil {
		breakdown.Deductions.TaxDeductions = *req.TaxDeductions
	}
	if req.OtherDeductions != nil {
		breakdown.Deductions.OtherDeductions = *req.OtherDeductions
	}
	if req.Currency != nil {
		breakdown.Currency = *req.Currency
	}

	// Calculate totals
	breakdown.Earnings.TotalEarnings = breakdown.Earnings.OutstandingSalary +
		breakdown.Earnings.LeaveEncashment +
		breakdown.Earnings.Bonus +
		breakdown.Earnings.Incentives +
		breakdown.Earnings.OtherEarnings

	breakdown.Deductions.TotalDeductions = breakdown.Deductions.OutstandingLoans +
		breakdown.Deductions.Advances +
		breakdown.Deductions.AssetDamages +
		breakdown.Deductions.TaxDeductions +
		breakdown.Deductions.OtherDeductions

	breakdown.NetSettlement = breakdown.Earnings.TotalEarnings - breakdown.Deductions.TotalDeductions

	// Serialize breakdown to JSON
	breakdownJSON, err := json.Marshal(breakdown)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize settlement breakdown: %w", err)
	}

	// Create or update settlement
	settlementID := ""
	if existing != nil {
		settlementID = existing.SettlementID
	} else {
		settlementID, err = s.settlementRepo.GenerateSettlementID()
		if err != nil {
			return nil, err
		}
	}

	now := time.Now()
	settlement := &employeeModels.FinalSettlement{
		SettlementID:        settlementID,
		OffboardingID:       offboardingID,
		EmployeeID:          workflow.EmployeeID,
		TenantID:            tenantID,
		Status:              "Calculated",
		CalculatedDate:      &now,
		SettlementBreakdown: string(breakdownJSON),
	}

	if calculatedByID != nil {
		settlement.CalculatedByID = calculatedByID
	}
	if req.Notes != nil {
		settlement.Notes = req.Notes
	}

	if existing != nil {
		if err := s.settlementRepo.Update(settlement); err != nil {
			return nil, err
		}
	} else {
		if err := s.settlementRepo.Create(settlement); err != nil {
			return nil, err
		}
	}

	// Update workflow settlement status
	workflow.FinalSettlementStatus = "Calculated"
	s.workflowRepo.Update(workflow)

	return settlement, nil
}

// GetSettlement retrieves settlement details
func (s *OffboardingService) GetSettlement(offboardingID string, tenantID *uint) (*employeeModels.FinalSettlement, error) {
	return s.settlementRepo.FindByOffboardingID(offboardingID, tenantID)
}

// ApproveSettlement approves a final settlement
func (s *OffboardingService) ApproveSettlement(offboardingID string, req *employeeModels.ApproveSettlementRequest, tenantID *uint) (*employeeModels.FinalSettlement, error) {
	settlement, err := s.settlementRepo.FindByOffboardingID(offboardingID, tenantID)
	if err != nil || settlement == nil {
		return nil, errors.New("settlement not found")
	}

	if settlement.Status != "Calculated" {
		return nil, errors.New("settlement must be in Calculated status to be approved")
	}

	now := time.Now()
	settlement.Status = "Approved"
	settlement.ApprovedDate = &now
	settlement.ApprovedByID = &req.ApprovedByID
	if req.ApprovalNotes != nil {
		settlement.ApprovalNotes = req.ApprovalNotes
	}

	if err := s.settlementRepo.Update(settlement); err != nil {
		return nil, err
	}

	// Update workflow
	workflow, _ := s.workflowRepo.FindByOffboardingID(offboardingID, tenantID)
	if workflow != nil {
		workflow.FinalSettlementStatus = "Approved"
		s.workflowRepo.Update(workflow)
	}

	return settlement, nil
}

// PaySettlement marks settlement as paid
func (s *OffboardingService) PaySettlement(offboardingID string, req *employeeModels.PaySettlementRequest, tenantID *uint) (*employeeModels.FinalSettlement, error) {
	settlement, err := s.settlementRepo.FindByOffboardingID(offboardingID, tenantID)
	if err != nil || settlement == nil {
		return nil, errors.New("settlement not found")
	}

	if settlement.Status != "Approved" {
		return nil, errors.New("settlement must be approved before payment")
	}

	paidDate, err := time.Parse("2006-01-02T15:04:05Z07:00", req.PaidDate)
	if err != nil {
		paidDate, err = time.Parse("2006-01-02", req.PaidDate)
		if err != nil {
			return nil, errors.New("invalid paid_date format")
		}
	}

	settlement.Status = "Paid"
	settlement.PaidDate = &paidDate
	settlement.PaidByID = &req.PaidByID
	settlement.PaymentMethod = &req.PaymentMethod
	if req.PaymentReference != nil {
		settlement.PaymentReference = req.PaymentReference
	}
	if req.PaymentNotes != nil {
		settlement.PaymentNotes = req.PaymentNotes
	}

	if err := s.settlementRepo.Update(settlement); err != nil {
		return nil, err
	}

	// Update workflow
	workflow, _ := s.workflowRepo.FindByOffboardingID(offboardingID, tenantID)
	if workflow != nil {
		workflow.FinalSettlementStatus = "Paid"
		s.workflowRepo.Update(workflow)
	}

	return settlement, nil
}

// CompleteOffboarding completes the offboarding workflow
func (s *OffboardingService) CompleteOffboarding(offboardingID string, req *employeeModels.CompleteOffboardingRequest, tenantID *uint) (*employeeModels.OffboardingWorkflow, error) {
	workflow, err := s.workflowRepo.FindByOffboardingID(offboardingID, tenantID)
	if err != nil {
		return nil, errors.New("offboarding workflow not found")
	}

	// Validate completion requirements
	clearances, err := s.clearanceRepo.FindByOffboardingID(offboardingID, tenantID)
	if err != nil {
		return nil, err
	}

	var pendingClearances []string
	for _, clearance := range clearances {
		if clearance.Status != "Cleared" {
			pendingClearances = append(pendingClearances, clearance.Department)
		}
	}

	if len(pendingClearances) > 0 {
		return nil, fmt.Errorf("cannot complete offboarding: pending clearances: %v", pendingClearances)
	}

	// Check exit interview (if required)
	if workflow.ExitInterviewRequired != nil && *workflow.ExitInterviewRequired {
		if workflow.ExitInterviewStatus != "Completed" {
			return nil, errors.New("exit interview must be completed before offboarding can be completed")
		}
	}

	// Check final settlement
	if workflow.FinalSettlementStatus != "Paid" {
		return nil, errors.New("final settlement must be paid before offboarding can be completed")
	}

	// Complete offboarding
	now := time.Now()
	workflow.Status = "Completed"
	workflow.CompletedAt = &now
	workflow.AccessRevokedDate = &now
	if req.CompletedByID != nil {
		workflow.CompletedByID = req.CompletedByID
	}

	if err := s.workflowRepo.Update(workflow); err != nil {
		return nil, err
	}

	// Update employee status to "Separated"
	if workflow.EmployeeDBID != nil {
		employee, err := s.employeeRepo.FindByID(*workflow.EmployeeDBID)
		if err == nil && employee != nil {
			employee.Status = "terminated" // or "separated" if you have that status
			employee.IsActive = false
			s.employeeRepo.Update(employee)
		}
	}

	return workflow, nil
}

// GetStatistics retrieves offboarding statistics
func (s *OffboardingService) GetStatistics(tenantID *uint) (*employeeModels.OffboardingStatisticsResponse, error) {
	stats := &employeeModels.OffboardingStatisticsResponse{
		SettlementsByStatus: make(map[string]int),
		SeparationsByType:   make(map[string]int),
		WorkflowsByStatus:   make(map[string]int),
	}

	// Get all workflows for tenant
	req := &employeeModels.ListOffboardingWorkflowsRequest{
		Page:     1,
		PageSize: 10000,
	}
	workflows, _, err := s.workflowRepo.List(req, tenantID)
	if err != nil {
		return nil, err
	}

	// Calculate statistics
	for _, workflow := range workflows {
		// Active separations (not completed)
		if workflow.Status != "Completed" {
			stats.ActiveSeparations++
		}

		// Separations by type
		stats.SeparationsByType[workflow.SeparationType]++

		// Workflows by status
		stats.WorkflowsByStatus[workflow.Status]++

		// Exit interviews
		if workflow.ExitInterviewStatus == "Scheduled" {
			stats.ScheduledExitInterviews++
		}
		if workflow.ExitInterviewStatus == "Completed" {
			stats.CompletedExitInterviews++
		}

		// Clearances
		clearances, _ := s.clearanceRepo.FindByOffboardingID(workflow.OffboardingID, tenantID)
		stats.TotalClearances += len(clearances)
		for _, clearance := range clearances {
			if clearance.Status == "Cleared" {
				stats.CompletedClearances++
			}
		}

		// Settlements
		settlement, _ := s.settlementRepo.FindByOffboardingID(workflow.OffboardingID, tenantID)
		if settlement != nil {
			stats.SettlementsByStatus[settlement.Status]++
			if settlement.Status == "Pending" || settlement.Status == "Calculated" {
				stats.PendingSettlements++
			}
		}
	}

	// Calculate clearance completion percentage
	if stats.TotalClearances > 0 {
		stats.ClearanceCompletionPercentage = float64(stats.CompletedClearances) / float64(stats.TotalClearances) * 100
	}

	return stats, nil
}
