package handlers

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/services"
	assetRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// OffboardingHandler handles offboarding HTTP requests
type OffboardingHandler struct {
	offboardingService *services.OffboardingService
	employeeRepo      *employeeRepos.EmployeeRepository
	assetRepo         *assetRepos.AssetRepository
}

// NewOffboardingHandler creates a new offboarding handler
func NewOffboardingHandler() *OffboardingHandler {
	return &OffboardingHandler{
		offboardingService: services.NewOffboardingService(),
		employeeRepo:       employeeRepos.NewEmployeeRepository(),
		assetRepo:          assetRepos.NewAssetRepository(),
	}
}

// InitiateSeparation initiates a new offboarding workflow
// @Summary Initiate employee separation
// @Description Creates a new offboarding workflow for an employee
// @Tags Employee Offboarding
// @Accept json
// @Produce json
// @Param request body models.InitiateSeparationRequest true "Separation request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/initiate [post]
func (h *OffboardingHandler) InitiateSeparation(c *gin.Context) {
	var req models.InitiateSeparationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var createdBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		createdBy = &userObj.ID
	}

	workflow, err := h.offboardingService.InitiateSeparation(&req, tenantID, createdBy)
	if err != nil {
		if err.Error() == "employee not found" {
			response.NotFound(c, err.Error())
			return
		}
		if err.Error() == "an active offboarding workflow already exists for this employee" {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Build response with employee details
	employee, _ := h.employeeRepo.FindByEmployeeID(req.EmployeeID)
	responseData := h.buildWorkflowResponse(workflow, employee, nil, nil, nil)

	response.Success(c, "Offboarding workflow initiated successfully", responseData)
}

// ListWorkflows lists offboarding workflows with pagination and filters
// @Summary List offboarding workflows
// @Description Retrieves a paginated list of offboarding workflows
// @Tags Employee Offboarding
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Status filter"
// @Param separation_type query string false "Separation type filter"
// @Param clearance_status query string false "Clearance status filter"
// @Param final_settlement query string false "Final settlement status filter"
// @Param search query string false "Search term"
// @Param department query string false "Department filter"
// @Param date_from query string false "Date from"
// @Param date_to query string false "Date to"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/workflows [get]
func (h *OffboardingHandler) ListWorkflows(c *gin.Context) {
	var req models.ListOffboardingWorkflowsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	workflows, total, err := h.offboardingService.ListWorkflows(&req, tenantID)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve workflows", err.Error())
		return
	}

	// Build response list
	responseList := make([]models.OffboardingWorkflowResponse, 0, len(workflows))
	for _, workflow := range workflows {
		employee, _ := h.employeeRepo.FindByEmployeeID(workflow.EmployeeID)
		responseList = append(responseList, *h.buildWorkflowResponse(&workflow, employee, nil, nil, nil))
	}

	// Pagination metadata
	page := req.Page
	if page == 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 20
	}
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	response.Success(c, "Offboarding workflows retrieved successfully", gin.H{
		"data": responseList,
		"meta": gin.H{
			"page":        page,
			"per_page":    pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetWorkflowDetails retrieves detailed information about a specific workflow
// @Summary Get offboarding workflow details
// @Description Retrieves complete details of a specific offboarding workflow
// @Tags Employee Offboarding
// @Produce json
// @Param offboarding_id path string true "Offboarding ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/workflows/{offboarding_id} [get]
func (h *OffboardingHandler) GetWorkflowDetails(c *gin.Context) {
	offboardingID := c.Param("offboarding_id")
	tenantID := middleware.GetTenantID(c)

	workflow, err := h.offboardingService.GetWorkflowByID(offboardingID, tenantID)
	if err != nil {
		response.NotFound(c, "Offboarding workflow not found")
		return
	}

	// Get employee details
	employee, _ := h.employeeRepo.FindByEmployeeID(workflow.EmployeeID)

	// Get clearances
	clearances, _ := h.offboardingService.GetClearances(offboardingID, tenantID)

	// Get asset returns
	assetReturns, _ := h.offboardingService.GetAssetReturns(offboardingID, tenantID)

	// Get settlement
	settlement, _ := h.offboardingService.GetSettlement(offboardingID, tenantID)

	responseData := h.buildWorkflowResponse(workflow, employee, clearances, assetReturns, settlement)
	response.Success(c, "Offboarding workflow retrieved successfully", responseData)
}

// UpdateWorkflow updates an offboarding workflow
// @Summary Update offboarding workflow
// @Description Updates basic information of an offboarding workflow
// @Tags Employee Offboarding
// @Accept json
// @Produce json
// @Param offboarding_id path string true "Offboarding ID"
// @Param request body models.UpdateOffboardingWorkflowRequest true "Update request"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/workflows/{offboarding_id} [patch]
func (h *OffboardingHandler) UpdateWorkflow(c *gin.Context) {
	offboardingID := c.Param("offboarding_id")
	var req models.UpdateOffboardingWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	workflow, err := h.offboardingService.UpdateWorkflow(offboardingID, &req, tenantID, updatedBy)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	employee, _ := h.employeeRepo.FindByEmployeeID(workflow.EmployeeID)
	responseData := h.buildWorkflowResponse(workflow, employee, nil, nil, nil)
	response.Success(c, "Offboarding workflow updated successfully", responseData)
}

// GetClearances retrieves all clearances for a workflow
// @Summary Get clearances
// @Description Retrieves all clearance tasks for an offboarding workflow
// @Tags Employee Offboarding
// @Produce json
// @Param offboarding_id path string true "Offboarding ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/workflows/{offboarding_id}/clearances [get]
func (h *OffboardingHandler) GetClearances(c *gin.Context) {
	offboardingID := c.Param("offboarding_id")
	tenantID := middleware.GetTenantID(c)

	clearances, err := h.offboardingService.GetClearances(offboardingID, tenantID)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve clearances", err.Error())
		return
	}

	// Build response
	responseList := make([]models.OffboardingClearanceResponse, 0, len(clearances))
	for _, clearance := range clearances {
		var clearedByName *string
		if clearance.ClearedByID != nil {
			employee, _ := h.employeeRepo.FindByEmployeeID(*clearance.ClearedByID)
			if employee != nil {
				name := employee.FullName()
				clearedByName = &name
			}
		}

		responseList = append(responseList, models.OffboardingClearanceResponse{
			ID:            clearance.ClearanceID,
			OffboardingID: clearance.OffboardingID,
			Department:    clearance.Department,
			ClearedBy:     clearedByName,
			ClearedByID:   clearance.ClearedByID,
			ClearanceDate: clearance.ClearanceDate,
			Status:        clearance.Status,
			Notes:         clearance.Notes,
			Issues:        clearance.Issues,
			CreatedAt:     clearance.CreatedAt,
			UpdatedAt:     clearance.UpdatedAt,
		})
	}

	response.Success(c, "Clearances retrieved successfully", responseList)
}

// UpdateClearance updates a clearance status
// @Summary Update clearance
// @Description Updates the status of a specific clearance
// @Tags Employee Offboarding
// @Accept json
// @Produce json
// @Param offboarding_id path string true "Offboarding ID"
// @Param clearance_id path string true "Clearance ID"
// @Param request body models.UpdateClearanceRequest true "Update request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/workflows/{offboarding_id}/clearances/{clearance_id} [patch]
func (h *OffboardingHandler) UpdateClearance(c *gin.Context) {
	offboardingID := c.Param("offboarding_id")
	clearanceID := c.Param("clearance_id")
	var req models.UpdateClearanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	clearance, err := h.offboardingService.UpdateClearance(offboardingID, clearanceID, &req, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	var clearedByName *string
	if clearance.ClearedByID != nil {
		employee, _ := h.employeeRepo.FindByEmployeeID(*clearance.ClearedByID)
		if employee != nil {
			name := employee.FullName()
			clearedByName = &name
		}
	}

	responseData := models.OffboardingClearanceResponse{
		ID:            clearance.ClearanceID,
		Status:        clearance.Status,
		ClearedBy:     clearedByName,
		ClearedByID:   clearance.ClearedByID,
		ClearanceDate: clearance.ClearanceDate,
		Notes:         clearance.Notes,
		UpdatedAt:     clearance.UpdatedAt,
	}

	response.Success(c, "Clearance updated successfully", responseData)
}

// buildWorkflowResponse builds a workflow response with employee details
func (h *OffboardingHandler) buildWorkflowResponse(
	workflow *models.OffboardingWorkflow,
	employee *models.Employee,
	clearances []models.OffboardingClearance,
	assetReturns []models.OffboardingAssetReturn,
	settlement *models.FinalSettlement,
) *models.OffboardingWorkflowResponse {
	resp := &models.OffboardingWorkflowResponse{
		OffboardingID:       workflow.OffboardingID,
		EmployeeID:          workflow.EmployeeID,
		EmployeeDBID:        workflow.EmployeeDBID,
		EmpID:               workflow.EmployeeID,
		SeparationType:      workflow.SeparationType,
		ResignationDate:     workflow.ResignationDate,
		LastWorkingDate:     workflow.LastWorkingDate,
		NoticePeriod:        workflow.NoticePeriodDays,
		ServedNoticePeriod:  workflow.ServedNoticePeriod,
		BuyoutAmount:        workflow.BuyoutAmount,
		Reason:              workflow.Reason,
		ReasonCode:          workflow.ReasonCode,
		ExitInterviewStatus: workflow.ExitInterviewStatus,
		ExitInterviewDate:   workflow.ExitInterviewDate,
		ExitInterviewConductedBy: workflow.ExitInterviewConductedByID,
		ExitInterviewNotes:  workflow.ExitInterviewNotes,
		ClearanceStatus:     workflow.ClearanceStatus,
		FinalSettlement:     workflow.FinalSettlementStatus,
		Status:              workflow.Status,
		ClearanceProgress:   workflow.ClearanceProgress,
		CompletedClearances: workflow.CompletedClearances,
		TotalClearances:     workflow.TotalClearances,
		AccessRevokedDate:   workflow.AccessRevokedDate,
		CreatedAt:           workflow.CreatedAt,
		UpdatedAt:           workflow.UpdatedAt,
	}

	if employee != nil {
		name := employee.FullName()
		resp.EmployeeName = &name
		resp.PhotoURL = employee.PhotoURL
		if employee.PositionID != nil {
			// Get position name if needed
		}
		if employee.DepartmentID != nil {
			// Get department name if needed
		}
	}

	// Build clearances response
	if clearances != nil {
		clearanceResponses := make([]models.OffboardingClearanceResponse, 0, len(clearances))
		for _, clearance := range clearances {
			var clearedByName *string
			if clearance.ClearedByID != nil {
				emp, _ := h.employeeRepo.FindByEmployeeID(*clearance.ClearedByID)
				if emp != nil {
					name := emp.FullName()
					clearedByName = &name
				}
			}
			clearanceResponses = append(clearanceResponses, models.OffboardingClearanceResponse{
				ID:            clearance.ClearanceID,
				OffboardingID: clearance.OffboardingID,
				Department:    clearance.Department,
				ClearedBy:     clearedByName,
				ClearedByID:   clearance.ClearedByID,
				ClearanceDate: clearance.ClearanceDate,
				Status:        clearance.Status,
				Notes:         clearance.Notes,
				Issues:        clearance.Issues,
				CreatedAt:     clearance.CreatedAt,
				UpdatedAt:     clearance.UpdatedAt,
			})
		}
		resp.Clearances = clearanceResponses
	}

	// Build assets response
	if assetReturns != nil {
		assetResponses := make([]models.OffboardingAssetResponse, 0, len(assetReturns))
		for _, assetReturn := range assetReturns {
			asset, _ := h.assetRepo.FindByID(assetReturn.AssetID)
			if asset != nil {
				assetName := asset.Brand + " " + asset.Model
				assetResponses = append(assetResponses, models.OffboardingAssetResponse{
					AssetID:            assetReturn.AssetID,
					AssetCode:          asset.AssetCode,
					AssetType:          asset.AssetType,
					AssetName:          &assetName,
					Brand:              &asset.Brand,
					Model:              &asset.Model,
					SerialNumber:       &asset.SerialNumber,
					AssetTag:           &asset.AssetCode, // Use asset code as tag
					AssignedDate:       asset.AssignedDate,
					ReturnStatus:       assetReturn.ReturnStatus,
					ReturnDate:         assetReturn.ReturnDate,
					ReturnCondition:    assetReturn.ReturnCondition,
					ReturnNotes:        assetReturn.ReturnNotes,
					ExpectedReturnDate: assetReturn.ExpectedReturnDate,
					IssueType:          assetReturn.IssueType,
					IssueDescription:   assetReturn.IssueDescription,
				})
			}
		}
		resp.Assets = assetResponses
	}

	// Build settlement response
	if settlement != nil {
		var breakdown *models.SettlementBreakdown
		if settlement.SettlementBreakdown != "" {
			json.Unmarshal([]byte(settlement.SettlementBreakdown), &breakdown)
		}

		var calculatedByName, approvedByName, paidByName *string
		if settlement.CalculatedByID != nil {
			emp, _ := h.employeeRepo.FindByEmployeeID(*settlement.CalculatedByID)
			if emp != nil {
				name := emp.FullName()
				calculatedByName = &name
			}
		}
		if settlement.ApprovedByID != nil {
			emp, _ := h.employeeRepo.FindByEmployeeID(*settlement.ApprovedByID)
			if emp != nil {
				name := emp.FullName()
				approvedByName = &name
			}
		}
		if settlement.PaidByID != nil {
			emp, _ := h.employeeRepo.FindByEmployeeID(*settlement.PaidByID)
			if emp != nil {
				name := emp.FullName()
				paidByName = &name
			}
		}

		var totalAmount *float64
		if breakdown != nil {
			totalAmount = &breakdown.NetSettlement
		}

		resp.FinalSettlementData = &models.FinalSettlementResponse{
			SettlementID:      settlement.SettlementID,
			OffboardingID:     settlement.OffboardingID,
			EmployeeID:        settlement.EmployeeID,
			Status:            settlement.Status,
			CalculatedDate:    settlement.CalculatedDate,
			CalculatedBy:      calculatedByName,
			ApprovedDate:      settlement.ApprovedDate,
			ApprovedBy:        approvedByName,
			ApprovedByID:      settlement.ApprovedByID,
			PaidDate:          settlement.PaidDate,
			PaidBy:            paidByName,
			PaidByID:          settlement.PaidByID,
			TotalAmount:       totalAmount,
			Breakdown:         breakdown,
			PaymentMethod:     settlement.PaymentMethod,
			PaymentReference:  settlement.PaymentReference,
			Notes:             settlement.Notes,
			CreatedAt:         settlement.CreatedAt,
			UpdatedAt:         settlement.UpdatedAt,
		}
	}

	return resp
}

// GetAssetReturns retrieves all asset returns for an offboarding workflow
// @Summary Get asset returns
// @Description Retrieves all assets assigned to the employee that need to be returned
// @Tags Employee Offboarding
// @Produce json
// @Param offboarding_id path string true "Offboarding ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/workflows/{offboarding_id}/assets [get]
func (h *OffboardingHandler) GetAssetReturns(c *gin.Context) {
	offboardingID := c.Param("offboarding_id")
	tenantID := middleware.GetTenantID(c)

	assetReturns, err := h.offboardingService.GetAssetReturns(offboardingID, tenantID)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve assets", err.Error())
		return
	}

	// Build response
	responseList := make([]models.OffboardingAssetResponse, 0, len(assetReturns))
	for _, assetReturn := range assetReturns {
		asset, _ := h.assetRepo.FindByID(assetReturn.AssetID)
		if asset != nil {
			assetName := asset.Brand + " " + asset.Model
			responseList = append(responseList, models.OffboardingAssetResponse{
				AssetID:            assetReturn.AssetID,
				AssetCode:          asset.AssetCode,
				AssetType:          asset.AssetType,
				AssetName:          &assetName,
				Brand:              &asset.Brand,
				Model:              &asset.Model,
				SerialNumber:       &asset.SerialNumber,
				AssetTag:           &asset.AssetCode,
				AssignedDate:       asset.AssignedDate,
				ReturnStatus:       assetReturn.ReturnStatus,
				ReturnDate:         assetReturn.ReturnDate,
				ReturnCondition:    assetReturn.ReturnCondition,
				ReturnNotes:        assetReturn.ReturnNotes,
				ExpectedReturnDate: assetReturn.ExpectedReturnDate,
				IssueType:          assetReturn.IssueType,
				IssueDescription:   assetReturn.IssueDescription,
			})
		}
	}

	response.Success(c, "Assets retrieved successfully", responseList)
}

// RecordAssetReturn records the return of an asset
// @Summary Record asset return
// @Description Records the return of an asset
// @Tags Employee Offboarding
// @Accept json
// @Produce json
// @Param offboarding_id path string true "Offboarding ID"
// @Param asset_id path int true "Asset ID"
// @Param request body models.RecordAssetReturnRequest true "Asset return request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/workflows/{offboarding_id}/assets/{asset_id}/return [post]
func (h *OffboardingHandler) RecordAssetReturn(c *gin.Context) {
	offboardingID := c.Param("offboarding_id")
	assetIDStr := c.Param("asset_id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		response.ValidationError(c, "Validation failed", "invalid asset_id")
		return
	}

	var req models.RecordAssetReturnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	assetReturn, err := h.offboardingService.RecordAssetReturn(offboardingID, uint(assetID), &req, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	asset, _ := h.assetRepo.FindByID(uint(assetID))
	var returnedByName *string
	if assetReturn.ReturnedByID != nil {
		employee, _ := h.employeeRepo.FindByEmployeeID(*assetReturn.ReturnedByID)
		if employee != nil {
			name := employee.FullName()
			returnedByName = &name
		}
	}

	responseData := gin.H{
		"asset_id":         assetReturn.AssetID,
		"asset_code":       asset.AssetCode,
		"return_status":    assetReturn.ReturnStatus,
		"return_date":      assetReturn.ReturnDate,
		"return_condition": assetReturn.ReturnCondition,
		"return_notes":     assetReturn.ReturnNotes,
		"returned_by":      returnedByName,
		"returned_by_id":   assetReturn.ReturnedByID,
		"updated_at":       assetReturn.UpdatedAt,
	}

	response.Success(c, "Asset return recorded successfully", responseData)
}

// RecordAssetIssue records an issue with an asset
// @Summary Record asset issue
// @Description Marks an asset as having issues
// @Tags Employee Offboarding
// @Accept json
// @Produce json
// @Param offboarding_id path string true "Offboarding ID"
// @Param asset_id path int true "Asset ID"
// @Param request body models.RecordAssetIssueRequest true "Asset issue request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/workflows/{offboarding_id}/assets/{asset_id}/issue [post]
func (h *OffboardingHandler) RecordAssetIssue(c *gin.Context) {
	offboardingID := c.Param("offboarding_id")
	assetIDStr := c.Param("asset_id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		response.ValidationError(c, "Validation failed", "invalid asset_id")
		return
	}

	var req models.RecordAssetIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	assetReturn, err := h.offboardingService.RecordAssetIssue(offboardingID, uint(assetID), &req, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	asset, _ := h.assetRepo.FindByID(uint(assetID))
	var reportedByName *string
	if assetReturn.ReportedByID != nil {
		employee, _ := h.employeeRepo.FindByEmployeeID(*assetReturn.ReportedByID)
		if employee != nil {
			name := employee.FullName()
			reportedByName = &name
		}
	}

	responseData := gin.H{
		"asset_id":          assetReturn.AssetID,
		"asset_code":        asset.AssetCode,
		"return_status":    assetReturn.ReturnStatus,
		"issue_type":       assetReturn.IssueType,
		"issue_description": assetReturn.IssueDescription,
		"reported_by":       reportedByName,
		"reported_by_id":    assetReturn.ReportedByID,
		"updated_at":        assetReturn.UpdatedAt,
	}

	response.Success(c, "Asset issue recorded successfully", responseData)
}

// ScheduleExitInterview schedules an exit interview
// @Summary Schedule exit interview
// @Description Schedules an exit interview for the offboarding employee
// @Tags Employee Offboarding
// @Accept json
// @Produce json
// @Param offboarding_id path string true "Offboarding ID"
// @Param request body models.ScheduleExitInterviewRequest true "Schedule request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/workflows/{offboarding_id}/exit-interview/schedule [post]
func (h *OffboardingHandler) ScheduleExitInterview(c *gin.Context) {
	offboardingID := c.Param("offboarding_id")
	var req models.ScheduleExitInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	workflow, err := h.offboardingService.ScheduleExitInterview(offboardingID, &req, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	var interviewerName *string
	if workflow.ExitInterviewerID != nil {
		employee, _ := h.employeeRepo.FindByEmployeeID(*workflow.ExitInterviewerID)
		if employee != nil {
			name := employee.FullName()
			interviewerName = &name
		}
	}

	responseData := gin.H{
		"offboarding_id":          workflow.OffboardingID,
		"exit_interview_status":   workflow.ExitInterviewStatus,
		"exit_interview_date":     workflow.ExitInterviewDate,
		"interviewer":             interviewerName,
		"interviewer_id":          workflow.ExitInterviewerID,
		"location":                workflow.ExitInterviewLocation,
		"notes":                   workflow.ExitInterviewNotes,
		"updated_at":              workflow.UpdatedAt,
	}

	response.Success(c, "Exit interview scheduled successfully", responseData)
}

// CompleteExitInterview completes an exit interview
// @Summary Complete exit interview
// @Description Marks the exit interview as completed and records feedback
// @Tags Employee Offboarding
// @Accept json
// @Produce json
// @Param offboarding_id path string true "Offboarding ID"
// @Param request body models.CompleteExitInterviewRequest true "Complete request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/workflows/{offboarding_id}/exit-interview/complete [post]
func (h *OffboardingHandler) CompleteExitInterview(c *gin.Context) {
	offboardingID := c.Param("offboarding_id")
	var req models.CompleteExitInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	workflow, err := h.offboardingService.CompleteExitInterview(offboardingID, &req, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	var conductedByName *string
	if workflow.ExitInterviewConductedByID != nil {
		employee, _ := h.employeeRepo.FindByEmployeeID(*workflow.ExitInterviewConductedByID)
		if employee != nil {
			name := employee.FullName()
			conductedByName = &name
		}
	}

	responseData := gin.H{
		"offboarding_id":          workflow.OffboardingID,
		"exit_interview_status":   workflow.ExitInterviewStatus,
		"exit_interview_date":     workflow.ExitInterviewDate,
		"interview_notes":        workflow.ExitInterviewNotes,
		"feedback_rating":         workflow.ExitInterviewRating,
		"would_recommend":         workflow.WouldRecommend,
		"conducted_by":            conductedByName,
		"conducted_by_id":         workflow.ExitInterviewConductedByID,
		"completed_at":            workflow.ExitInterviewCompletedAt,
		"updated_at":              workflow.UpdatedAt,
	}

	response.Success(c, "Exit interview completed successfully", responseData)
}

// CancelExitInterview cancels an exit interview
// @Summary Cancel exit interview
// @Description Cancels a scheduled exit interview
// @Tags Employee Offboarding
// @Accept json
// @Produce json
// @Param offboarding_id path string true "Offboarding ID"
// @Param request body models.CancelExitInterviewRequest true "Cancel request"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/workflows/{offboarding_id}/exit-interview/cancel [post]
func (h *OffboardingHandler) CancelExitInterview(c *gin.Context) {
	offboardingID := c.Param("offboarding_id")
	var req models.CancelExitInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	workflow, err := h.offboardingService.CancelExitInterview(offboardingID, &req, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	responseData := gin.H{
		"offboarding_id":        workflow.OffboardingID,
		"exit_interview_status": workflow.ExitInterviewStatus,
		"exit_interview_date":   workflow.ExitInterviewDate,
		"updated_at":            workflow.UpdatedAt,
	}

	response.Success(c, "Exit interview cancelled successfully", responseData)
}

// CalculateSettlement calculates the final settlement
// @Summary Calculate final settlement
// @Description Calculates the final settlement amount for the employee
// @Tags Employee Offboarding
// @Accept json
// @Produce json
// @Param offboarding_id path string true "Offboarding ID"
// @Param request body models.CalculateSettlementRequest true "Calculate request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/workflows/{offboarding_id}/settlement/calculate [post]
func (h *OffboardingHandler) CalculateSettlement(c *gin.Context) {
	offboardingID := c.Param("offboarding_id")
	var req models.CalculateSettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	// Get calculated_by_id from request if provided, otherwise use current user's employee ID
	var calculatedByID *string
	// For now, we'll use the user's employee ID if available
	// This can be enhanced to get from user context

	settlement, err := h.offboardingService.CalculateSettlement(offboardingID, &req, tenantID, calculatedByID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Parse breakdown
	var breakdown *models.SettlementBreakdown
	if settlement.SettlementBreakdown != "" {
		json.Unmarshal([]byte(settlement.SettlementBreakdown), &breakdown)
	}

	responseData := gin.H{
		"settlement_id":   settlement.SettlementID,
		"offboarding_id":  settlement.OffboardingID,
		"employee_id":     settlement.EmployeeID,
		"status":          settlement.Status,
		"calculated_date": settlement.CalculatedDate,
		"calculated_by":   settlement.CalculatedByID,
		"breakdown":       breakdown,
		"notes":           settlement.Notes,
		"created_at":      settlement.CreatedAt,
		"updated_at":      settlement.UpdatedAt,
	}

	response.Success(c, "Final settlement calculated successfully", responseData)
}

// GetSettlement retrieves settlement details
// @Summary Get settlement details
// @Description Retrieves the final settlement details for an offboarding workflow
// @Tags Employee Offboarding
// @Produce json
// @Param offboarding_id path string true "Offboarding ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/workflows/{offboarding_id}/settlement [get]
func (h *OffboardingHandler) GetSettlement(c *gin.Context) {
	offboardingID := c.Param("offboarding_id")
	tenantID := middleware.GetTenantID(c)

	settlement, err := h.offboardingService.GetSettlement(offboardingID, tenantID)
	if err != nil || settlement == nil {
		response.NotFound(c, "Settlement not found")
		return
	}

	// Parse breakdown
	var breakdown *models.SettlementBreakdown
	if settlement.SettlementBreakdown != "" {
		json.Unmarshal([]byte(settlement.SettlementBreakdown), &breakdown)
	}

	var calculatedByName, approvedByName, paidByName *string
	if settlement.CalculatedByID != nil {
		emp, _ := h.employeeRepo.FindByEmployeeID(*settlement.CalculatedByID)
		if emp != nil {
			name := emp.FullName()
			calculatedByName = &name
		}
	}
	if settlement.ApprovedByID != nil {
		emp, _ := h.employeeRepo.FindByEmployeeID(*settlement.ApprovedByID)
		if emp != nil {
			name := emp.FullName()
			approvedByName = &name
		}
	}
	if settlement.PaidByID != nil {
		emp, _ := h.employeeRepo.FindByEmployeeID(*settlement.PaidByID)
		if emp != nil {
			name := emp.FullName()
			paidByName = &name
		}
	}

	var totalAmount *float64
	if breakdown != nil {
		totalAmount = &breakdown.NetSettlement
	}

	responseData := gin.H{
		"settlement_id":     settlement.SettlementID,
		"offboarding_id":    settlement.OffboardingID,
		"employee_id":       settlement.EmployeeID,
		"status":            settlement.Status,
		"calculated_date":   settlement.CalculatedDate,
		"calculated_by":    calculatedByName,
		"approved_date":     settlement.ApprovedDate,
		"approved_by":       approvedByName,
		"approved_by_id":    settlement.ApprovedByID,
		"paid_date":         settlement.PaidDate,
		"paid_by":           paidByName,
		"paid_by_id":        settlement.PaidByID,
		"total_amount":      totalAmount,
		"breakdown":         breakdown,
		"payment_method":    settlement.PaymentMethod,
		"payment_reference": settlement.PaymentReference,
		"notes":             settlement.Notes,
		"created_at":        settlement.CreatedAt,
		"updated_at":        settlement.UpdatedAt,
	}

	response.Success(c, "Settlement details retrieved successfully", responseData)
}

// ApproveSettlement approves a final settlement
// @Summary Approve final settlement
// @Description Approves the calculated final settlement
// @Tags Employee Offboarding
// @Accept json
// @Produce json
// @Param offboarding_id path string true "Offboarding ID"
// @Param request body models.ApproveSettlementRequest true "Approve request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/workflows/{offboarding_id}/settlement/approve [post]
func (h *OffboardingHandler) ApproveSettlement(c *gin.Context) {
	offboardingID := c.Param("offboarding_id")
	var req models.ApproveSettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	settlement, err := h.offboardingService.ApproveSettlement(offboardingID, &req, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	var approvedByName *string
	if settlement.ApprovedByID != nil {
		employee, _ := h.employeeRepo.FindByEmployeeID(*settlement.ApprovedByID)
		if employee != nil {
			name := employee.FullName()
			approvedByName = &name
		}
	}

	responseData := gin.H{
		"settlement_id":   settlement.SettlementID,
		"status":          settlement.Status,
		"approved_date":   settlement.ApprovedDate,
		"approved_by":     approvedByName,
		"approved_by_id":  settlement.ApprovedByID,
		"approval_notes":  settlement.ApprovalNotes,
		"updated_at":      settlement.UpdatedAt,
	}

	response.Success(c, "Final settlement approved successfully", responseData)
}

// PaySettlement marks settlement as paid
// @Summary Mark settlement as paid
// @Description Marks the final settlement as paid
// @Tags Employee Offboarding
// @Accept json
// @Produce json
// @Param offboarding_id path string true "Offboarding ID"
// @Param request body models.PaySettlementRequest true "Pay request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/workflows/{offboarding_id}/settlement/pay [post]
func (h *OffboardingHandler) PaySettlement(c *gin.Context) {
	offboardingID := c.Param("offboarding_id")
	var req models.PaySettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	settlement, err := h.offboardingService.PaySettlement(offboardingID, &req, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	var paidByName *string
	if settlement.PaidByID != nil {
		employee, _ := h.employeeRepo.FindByEmployeeID(*settlement.PaidByID)
		if employee != nil {
			name := employee.FullName()
			paidByName = &name
		}
	}

	responseData := gin.H{
		"settlement_id":      settlement.SettlementID,
		"status":             settlement.Status,
		"paid_date":          settlement.PaidDate,
		"payment_method":    settlement.PaymentMethod,
		"payment_reference": settlement.PaymentReference,
		"paid_by":            paidByName,
		"paid_by_id":         settlement.PaidByID,
		"payment_notes":      settlement.PaymentNotes,
		"updated_at":         settlement.UpdatedAt,
	}

	response.Success(c, "Settlement marked as paid successfully", responseData)
}

// CompleteOffboarding completes the offboarding workflow
// @Summary Complete offboarding
// @Description Completes the offboarding workflow
// @Tags Employee Offboarding
// @Accept json
// @Produce json
// @Param offboarding_id path string true "Offboarding ID"
// @Param request body models.CompleteOffboardingRequest true "Complete request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/workflows/{offboarding_id}/complete [post]
func (h *OffboardingHandler) CompleteOffboarding(c *gin.Context) {
	offboardingID := c.Param("offboarding_id")
	var req models.CompleteOffboardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	workflow, err := h.offboardingService.CompleteOffboarding(offboardingID, &req, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	responseData := gin.H{
		"offboarding_id":      workflow.OffboardingID,
		"status":              workflow.Status,
		"completed_at":        workflow.CompletedAt,
		"completed_by_id":     workflow.CompletedByID,
		"employee_status":     "Separated",
		"access_revoked_date": workflow.AccessRevokedDate,
		"updated_at":          workflow.UpdatedAt,
	}

	response.Success(c, "Offboarding workflow completed successfully", responseData)
}

// GetStatistics retrieves offboarding statistics
// @Summary Get offboarding statistics
// @Description Retrieves dashboard statistics for the offboarding page
// @Tags Employee Offboarding
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /api/v1/employees/offboarding/statistics [get]
func (h *OffboardingHandler) GetStatistics(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	stats, err := h.offboardingService.GetStatistics(tenantID)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve statistics", err.Error())
		return
	}

	response.Success(c, "Statistics retrieved successfully", stats)
}
