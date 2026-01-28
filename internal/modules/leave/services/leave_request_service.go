package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
)

type LeaveRequestService struct {
	repo            *repositories.LeaveRequestRepository
	balanceRepo     *repositories.LeaveBalanceRepository
	policyRepo      *repositories.LeavePolicyRepository
	leaveTypeRepo   *repositories.LeaveTypeRepository
	approvalRepo    *repositories.LeaveApprovalRepository
	documentRepo    *repositories.LeaveDocumentRepository
	holidayService  *HolidayService
	balanceService  *LeaveBalanceService
	employeeRepo    *employeeRepos.EmployeeRepository
}

func NewLeaveRequestService() *LeaveRequestService {
	return &LeaveRequestService{
		repo:           repositories.NewLeaveRequestRepository(),
		balanceRepo:    repositories.NewLeaveBalanceRepository(),
		policyRepo:     repositories.NewLeavePolicyRepository(),
		leaveTypeRepo:  repositories.NewLeaveTypeRepository(),
		approvalRepo:   repositories.NewLeaveApprovalRepository(),
		documentRepo:   repositories.NewLeaveDocumentRepository(),
		holidayService: NewHolidayService(),
		balanceService: NewLeaveBalanceService(),
		employeeRepo:   employeeRepos.NewEmployeeRepository(),
	}
}

// CalculateLeaveDays calculates working days excluding weekends and holidays
func (s *LeaveRequestService) CalculateLeaveDays(fromDate, toDate time.Time, halfDay bool, employeeID string, tenantID *uint) (float64, int, int, error) {
	var workingDays float64 = 0
	weekendCount := 0
	holidayCount := 0

	current := fromDate
	for !current.After(toDate) {
		// Check if weekend (Saturday = 6, Sunday = 0)
		weekday := current.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			weekendCount++
			current = current.AddDate(0, 0, 1)
			continue
		}

		// Check if holiday
		isHoliday, err := s.holidayService.IsHoliday(current, tenantID)
		if err != nil {
			return 0, 0, 0, err
		}
		if isHoliday {
			holidayCount++
			current = current.AddDate(0, 0, 1)
			continue
		}

		workingDays++
		current = current.AddDate(0, 0, 1)
	}

	// Handle half day
	if halfDay && workingDays > 0 {
		workingDays = 0.5
	}

	return workingDays, weekendCount, holidayCount, nil
}

// CreateLeaveRequest creates a new leave request
func (s *LeaveRequestService) CreateLeaveRequest(req *models.CreateLeaveRequestRequest, employeeID string, tenantID *uint, createdBy *uint) (*models.LeaveRequest, error) {
	// Parse dates
	fromDate, err := time.Parse("2006-01-02", req.FromDate)
	if err != nil {
		return nil, errors.New("invalid from_date format. Use YYYY-MM-DD")
	}

	toDate, err := time.Parse("2006-01-02", req.ToDate)
	if err != nil {
		return nil, errors.New("invalid to_date format. Use YYYY-MM-DD")
	}

	if toDate.Before(fromDate) {
		return nil, errors.New("to_date must be after from_date")
	}

	applicationDate, err := time.Parse("2006-01-02", req.ApplicationDate)
	if err != nil {
		applicationDate = time.Now()
	}

	reportingBackDate, err := time.Parse("2006-01-02", req.ReportingBackDate)
	if err != nil {
		return nil, errors.New("invalid reporting_back_date format. Use YYYY-MM-DD")
	}

	// Validate leave type
	_, err = s.leaveTypeRepo.FindByCode(req.LeaveTypeCode)
	if err != nil {
		return nil, errors.New("leave type not found")
	}

	// Calculate total days
	halfDay := false
	if req.HalfDay != nil {
		halfDay = *req.HalfDay
	}

	totalDays, _, _, err := s.CalculateLeaveDays(fromDate, toDate, halfDay, employeeID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate leave days: %w", err)
	}

	// Generate application number
	appNumber, err := s.repo.GenerateApplicationNumber(req.LeavePeriodYear)
	if err != nil {
		return nil, fmt.Errorf("failed to generate application number: %w", err)
	}

	// Determine status
	status := "pending"
	saveAsDraft := false
	if req.SaveAsDraft != nil {
		saveAsDraft = *req.SaveAsDraft
		if saveAsDraft {
			status = "draft"
		}
	}

	// Create leave request
	leaveRequest := &models.LeaveRequest{
		TenantID:                  tenantID,
		ApplicationNumber:         appNumber,
		DocumentNumber:            "HR.FO.04.00",
		EmployeeID:                employeeID,
		LeaveTypeCode:             req.LeaveTypeCode,
		FromDate:                  fromDate,
		ToDate:                    toDate,
		TotalDays:                 totalDays,
		HalfDay:                   halfDay,
		Reason:                    req.Reason,
		WorkDelegatedTo:           req.WorkDelegatedTo,
		HandoverNotes:             req.HandoverNotes,
		ApplicationDate:           applicationDate,
		LeavePeriodYear:           req.LeavePeriodYear,
		ReportingBackDate:         reportingBackDate,
		EmployeeContactNumber:     req.EmployeeContactNumber,
		EmergencyContactPerson:    req.EmergencyContactPerson,
		EmergencyContactNumber:    req.EmergencyContactNumber,
		EmployeeSignature:         req.EmployeeSignature,
		Status:                    status,
		PreviousDaysUsed:          req.LeaveTracker.PreviousDaysUsed,
		DaysRemainingAfterRequest: req.LeaveTracker.DaysRemainingAfterRequest,
		LeaveTakenFromPreviousYear: req.LeaveTracker.LeaveTakenFromPreviousYear,
		LeaveBalanceFromPreviousYear: req.LeaveTracker.LeaveBalanceFromPreviousYear,
		CreatedBy:                 createdBy,
		UpdatedBy:                 createdBy,
	}

	if err := s.repo.Create(leaveRequest); err != nil {
		return nil, fmt.Errorf("failed to create leave request: %w", err)
	}

	// Create documents if provided
	if len(req.Documents) > 0 {
		documents := make([]models.LeaveDocument, len(req.Documents))
		for i, doc := range req.Documents {
			documents[i] = models.LeaveDocument{
				TenantID:      tenantID,
				LeaveRequestID: leaveRequest.ID,
				FileName:      doc.FileName,
				FileURL:       doc.FileURL,
				FileType:      doc.FileType,
				FileSize:      doc.FileSize,
				DocumentType:  doc.DocumentType,
			}
		}
		if err := s.documentRepo.CreateBatch(documents); err != nil {
			return nil, fmt.Errorf("failed to create documents: %w", err)
		}
	}

	// Create approval workflow if not draft
	if !saveAsDraft {
		if err := s.approvalRepo.CreateApprovalWorkflow(leaveRequest.ID, tenantID); err != nil {
			return nil, fmt.Errorf("failed to create approval workflow: %w", err)
		}

		// Add to pending balance
		if err := s.balanceService.AddPendingBalance(employeeID, req.LeaveTypeCode, req.LeavePeriodYear, totalDays); err != nil {
			// Log error but don't fail the request creation
			_ = err
		}
	}

	return s.repo.FindByID(leaveRequest.ID)
}

// GetLeaveRequestByID retrieves a leave request by ID
func (s *LeaveRequestService) GetLeaveRequestByID(id uint) (*models.LeaveRequest, error) {
	return s.repo.FindByID(id)
}

// UpdateLeaveRequest updates a leave request (only for draft or returned_for_info)
func (s *LeaveRequestService) UpdateLeaveRequest(id uint, req *models.UpdateLeaveRequestRequest, updatedBy *uint) (*models.LeaveRequest, error) {
	leaveRequest, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("leave request not found")
	}

	if leaveRequest.Status != "draft" && leaveRequest.Status != "returned_for_info" {
		return nil, errors.New("only draft or returned_for_info requests can be updated")
	}

	if req.FromDate != nil {
		fromDate, err := time.Parse("2006-01-02", *req.FromDate)
		if err != nil {
			return nil, errors.New("invalid from_date format")
		}
		leaveRequest.FromDate = fromDate
	}

	if req.ToDate != nil {
		toDate, err := time.Parse("2006-01-02", *req.ToDate)
		if err != nil {
			return nil, errors.New("invalid to_date format")
		}
		leaveRequest.ToDate = toDate
	}

	if req.HalfDay != nil {
		leaveRequest.HalfDay = *req.HalfDay
	}

	if req.Reason != nil {
		leaveRequest.Reason = *req.Reason
	}
	if req.WorkDelegatedTo != nil {
		leaveRequest.WorkDelegatedTo = req.WorkDelegatedTo
	}
	if req.HandoverNotes != nil {
		leaveRequest.HandoverNotes = req.HandoverNotes
	}
	if req.ReportingBackDate != nil {
		reportingBackDate, err := time.Parse("2006-01-02", *req.ReportingBackDate)
		if err != nil {
			return nil, errors.New("invalid reporting_back_date format")
		}
		leaveRequest.ReportingBackDate = reportingBackDate
	}
	if req.EmployeeContactNumber != nil {
		leaveRequest.EmployeeContactNumber = *req.EmployeeContactNumber
	}
	if req.EmergencyContactPerson != nil {
		leaveRequest.EmergencyContactPerson = *req.EmergencyContactPerson
	}
	if req.EmergencyContactNumber != nil {
		leaveRequest.EmergencyContactNumber = *req.EmergencyContactNumber
	}

	// Recalculate days if dates changed
	if req.FromDate != nil || req.ToDate != nil {
		totalDays, _, _, err := s.CalculateLeaveDays(leaveRequest.FromDate, leaveRequest.ToDate, leaveRequest.HalfDay, leaveRequest.EmployeeID, leaveRequest.TenantID)
		if err != nil {
			return nil, fmt.Errorf("failed to recalculate days: %w", err)
		}
		leaveRequest.TotalDays = totalDays
	}

	if req.LeaveTracker != nil {
		leaveRequest.PreviousDaysUsed = req.LeaveTracker.PreviousDaysUsed
		leaveRequest.DaysRemainingAfterRequest = req.LeaveTracker.DaysRemainingAfterRequest
		leaveRequest.LeaveTakenFromPreviousYear = req.LeaveTracker.LeaveTakenFromPreviousYear
		leaveRequest.LeaveBalanceFromPreviousYear = req.LeaveTracker.LeaveBalanceFromPreviousYear
	}

	leaveRequest.UpdatedBy = updatedBy

	if err := s.repo.Update(leaveRequest); err != nil {
		return nil, fmt.Errorf("failed to update leave request: %w", err)
	}

	return s.repo.FindByID(id)
}

// ApproveLeaveRequest approves a leave request
func (s *LeaveRequestService) ApproveLeaveRequest(id uint, req *models.ApproveLeaveRequestRequest, approverID string, tenantID *uint) (*models.LeaveRequest, error) {
	leaveRequest, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("leave request not found")
	}

	if leaveRequest.Status != "pending" && leaveRequest.Status != "returned_for_info" {
		return nil, errors.New("only pending or returned_for_info requests can be approved")
	}

	// Get pending approvals
	approvals, err := s.approvalRepo.FindPendingByLeaveRequestID(id)
	if err != nil || len(approvals) == 0 {
		return nil, errors.New("no pending approvals found")
	}

	// Find the current approval level
	currentApproval := &approvals[0]
	if currentApproval.ApproverType != req.ApproverType {
		return nil, errors.New("approver type mismatch")
	}

	// Update approval
	now := time.Now()
	currentApproval.ApproverID = &approverID
	currentApproval.ApproverName = &req.ApproverName
	currentApproval.ApproverSignature = &req.ApproverSignature
	currentApproval.Status = "approved"
	currentApproval.ApprovedAt = &now
	currentApproval.Remarks = req.Remarks

	if err := s.approvalRepo.Update(currentApproval); err != nil {
		return nil, fmt.Errorf("failed to update approval: %w", err)
	}

	// Check if there are more approval levels
	remainingApprovals, _ := s.approvalRepo.FindPendingByLeaveRequestID(id)
	if len(remainingApprovals) == 0 {
		// All approvals done - mark as approved
		leaveRequest.Status = "approved"
		leaveRequest.BalanceAfterApproval = &leaveRequest.DaysRemainingAfterRequest

		// Update balance
		if err := s.balanceService.UpdateBalanceAfterApproval(
			leaveRequest.EmployeeID,
			leaveRequest.LeaveTypeCode,
			leaveRequest.LeavePeriodYear,
			leaveRequest.TotalDays,
		); err != nil {
			// Log error but continue
			_ = err
		}

		// TODO: Update employee status to "on_leave" if leave is active
	}

	leaveRequest.UpdatedBy = nil // Will be set by handler

	if err := s.repo.Update(leaveRequest); err != nil {
		return nil, fmt.Errorf("failed to update leave request: %w", err)
	}

	return s.repo.FindByID(id)
}

// RejectLeaveRequest rejects a leave request
func (s *LeaveRequestService) RejectLeaveRequest(id uint, req *models.RejectLeaveRequestRequest, approverID string) (*models.LeaveRequest, error) {
	leaveRequest, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("leave request not found")
	}

	if leaveRequest.Status != "pending" && leaveRequest.Status != "returned_for_info" {
		return nil, errors.New("only pending or returned_for_info requests can be rejected")
	}

	// Get pending approvals
	approvals, err := s.approvalRepo.FindPendingByLeaveRequestID(id)
	if err != nil || len(approvals) == 0 {
		return nil, errors.New("no pending approvals found")
	}

	// Update first pending approval
	currentApproval := &approvals[0]
	now := time.Now()
	currentApproval.ApproverID = &approverID
	currentApproval.ApproverName = &req.ApproverName
	currentApproval.ApproverSignature = &req.ApproverSignature
	currentApproval.Status = "rejected"
	currentApproval.ApprovedAt = &now
	currentApproval.Remarks = &req.RejectionReason

	if err := s.approvalRepo.Update(currentApproval); err != nil {
		return nil, fmt.Errorf("failed to update approval: %w", err)
	}

	// Update request status
	leaveRequest.Status = "rejected"

	// Remove from pending balance
	if err := s.balanceService.RemovePendingBalance(
		leaveRequest.EmployeeID,
		leaveRequest.LeaveTypeCode,
		leaveRequest.LeavePeriodYear,
		leaveRequest.TotalDays,
	); err != nil {
		// Log error but continue
		_ = err
	}

	if err := s.repo.Update(leaveRequest); err != nil {
		return nil, fmt.Errorf("failed to update leave request: %w", err)
	}

	return s.repo.FindByID(id)
}

// CancelLeaveRequest cancels a leave request
func (s *LeaveRequestService) CancelLeaveRequest(id uint, req *models.CancelLeaveRequestRequest) (*models.LeaveRequest, error) {
	leaveRequest, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("leave request not found")
	}

	if leaveRequest.Status != "pending" && leaveRequest.Status != "approved" {
		return nil, errors.New("only pending or approved requests can be cancelled")
	}

	leaveRequest.Status = "cancelled"
	leaveRequest.CancellationReason = &req.CancellationReason

	// Remove from pending balance if pending
	if leaveRequest.Status == "pending" {
		if err := s.balanceService.RemovePendingBalance(
			leaveRequest.EmployeeID,
			leaveRequest.LeaveTypeCode,
			leaveRequest.LeavePeriodYear,
			leaveRequest.TotalDays,
		); err != nil {
			// Log error but continue
			_ = err
		}
	}

	if err := s.repo.Update(leaveRequest); err != nil {
		return nil, fmt.Errorf("failed to update leave request: %w", err)
	}

	return s.repo.FindByID(id)
}

// DeleteLeaveRequest deletes a draft leave request
func (s *LeaveRequestService) DeleteLeaveRequest(id uint) error {
	leaveRequest, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("leave request not found")
	}

	if leaveRequest.Status != "draft" {
		return errors.New("only draft requests can be deleted")
	}

	return s.repo.Delete(id)
}

// ListLeaveRequests lists leave requests with pagination and filters
func (s *LeaveRequestService) ListLeaveRequests(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.LeaveRequest, int64, error) {
	return s.repo.List(tenantID, page, pageSize, filters)
}

// ReturnForInfo returns a leave request for additional information
func (s *LeaveRequestService) ReturnForInfo(id uint, req *models.ReturnForInfoRequest) (*models.LeaveRequest, error) {
	leaveRequest, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("leave request not found")
	}

	if leaveRequest.Status != "pending" {
		return nil, errors.New("only pending requests can be returned for information")
	}

	leaveRequest.Status = "returned_for_info"
	leaveRequest.InformationRequired = &req.InformationRequired

	if err := s.repo.Update(leaveRequest); err != nil {
		return nil, fmt.Errorf("failed to update leave request: %w", err)
	}

	return s.repo.FindByID(id)
}
