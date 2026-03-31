package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

// PayrollHandler handles payroll-related HTTP requests
type PayrollHandler struct {
	runService        *services.PayrollRunService
	structureService  *services.SalaryStructureService
	payslipService    *services.PayslipService
	loanService       *services.LoanService
	complianceService *services.ComplianceService
	prePayrollService *services.PrePayrollService
}

// NewPayrollHandler creates a new PayrollHandler instance
func NewPayrollHandler() *PayrollHandler {
	return &PayrollHandler{
		runService:        services.NewPayrollRunService(),
		structureService:  services.NewSalaryStructureService(),
		payslipService:    services.NewPayslipService(),
		loanService:       services.NewLoanService(),
		complianceService: services.NewComplianceService(),
		prePayrollService: services.NewPrePayrollService(),
	}
}

// ============ Dashboard Handlers ============

// GetDashboard returns payroll dashboard metrics
func (h *PayrollHandler) GetDashboard(c *gin.Context) {
	tenantID := utils.GetTenantID(c)

	dashboard, err := h.runService.GetDashboard(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dashboard,
	})
}

func (h *PayrollHandler) ListEmployeesWithSalaryInfo(c *gin.Context) {
	tenantID := utils.GetTenantID(c)
	status := c.Query("status")
	page := utils.GetPage(c)

	pageSize := utils.GetPageSize(c)
	if ps := c.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	employees, meta, err := h.prePayrollService.GetEmployeesWithSalaryInfo(tenantID, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    employees,
		"meta":    meta,
	})
}

// ============ Payroll Run Handlers ============

// ListPayrollRuns lists all payroll runs
func (h *PayrollHandler) ListPayrollRuns(c *gin.Context) {
	tenantID := utils.GetTenantID(c)
	status := c.Query("status")
	payYear, _ := strconv.Atoi(c.Query("year"))
	payMonth, _ := strconv.Atoi(c.Query("month"))
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	runs, total, err := h.runService.ListPayrollRuns(tenantID, status, payYear, payMonth, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    runs,
		"meta": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
			"pages":    (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetPayrollRun retrieves a specific payroll run
func (h *PayrollHandler) GetPayrollRun(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := utils.GetTenantID(c)
	run, err := h.runService.GetPayrollRun(uint(id), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payroll run not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    run,
	})
}

// CreatePayrollRun creates a new payroll run
func (h *PayrollHandler) CreatePayrollRun(c *gin.Context) {
	var req struct {
		PayMonth     int                 `json:"payMonth" binding:"required"`
		PayYear      int                 `json:"payYear" binding:"required"`
		PayFrequency models.PayFrequency `json:"payFrequency"`
		RunName      string              `json:"runName"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := utils.GetUserID(c)
	userName := utils.GetUserName(c)

	run := &models.PayrollRun{
		PayMonth:      req.PayMonth,
		PayYear:       req.PayYear,
		PayFrequency:  req.PayFrequency,
		RunName:       req.RunName,
		CreatedByID:   &userID,
		CreatedByName: userName,
	}

	if run.PayFrequency == "" {
		run.PayFrequency = models.PayFrequencyMonthly
	}

	if err := h.runService.CreatePayrollRun(run); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Payroll run created successfully",
		"data":    run,
	})
}

// UpdatePayrollRun updates a payroll run
func (h *PayrollHandler) UpdatePayrollRun(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := utils.GetTenantID(c)
	run, err := h.runService.GetPayrollRun(uint(id), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payroll run not found"})
		return
	}

	var req struct {
		RunName    string     `json:"runName"`
		CutoffDate *time.Time `json:"cutoffDate"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.RunName != "" {
		run.RunName = req.RunName
	}
	if req.CutoffDate != nil {
		run.CutoffDate = req.CutoffDate
	}

	if err := h.runService.UpdatePayrollRun(run); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payroll run updated successfully",
		"data":    run,
	})
}

// DeletePayrollRun deletes a payroll run
func (h *PayrollHandler) DeletePayrollRun(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := utils.GetTenantID(c)
	if err := h.runService.DeletePayrollRun(uint(id), tenantID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payroll run deleted successfully",
	})
}

// GetCurrentPayrollRun retrieves the current active payroll run
func (h *PayrollHandler) GetCurrentPayrollRun(c *gin.Context) {
	tenantID := utils.GetTenantID(c)
	run, err := h.runService.GetCurrentPayrollRun(tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No active payroll run found",
			"error":   "no_active_payroll_run",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Current payroll run retrieved successfully",
		"data":    run,
	})
}

// GetPayrollRunSummary retrieves summary for a payroll run
func (h *PayrollHandler) GetPayrollRunSummary(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := utils.GetTenantID(c)
	summary, err := h.runService.GetPayrollRunSummary(uint(id), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

// AdvancePayrollStep advances payroll run to next step
func (h *PayrollHandler) AdvancePayrollStep(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := utils.GetTenantID(c)
	userID := utils.GetUserID(c)
	userName := utils.GetUserName(c)

	if err := h.runService.AdvanceStep(uint(id), tenantID, userID, userName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get updated run
	run, _ := h.runService.GetPayrollRun(uint(id), tenantID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payroll advanced to next step",
		"data":    run,
	})
}

// RevertPayrollStep reverts payroll run to previous step
func (h *PayrollHandler) RevertPayrollStep(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := utils.GetTenantID(c)
	if err := h.runService.RevertStep(uint(id), tenantID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	run, _ := h.runService.GetPayrollRun(uint(id), tenantID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payroll reverted to previous step",
		"data":    run,
	})
}

// RunPreCheck runs pre-checks for payroll processing
func (h *PayrollHandler) RunPreCheck(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := utils.GetTenantID(c)
	result, err := h.runService.RunPreCheck(uint(id), tenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// CalculatePayroll calculates payroll for all employees
func (h *PayrollHandler) CalculatePayroll(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := utils.GetTenantID(c)
	taxYear := time.Now().Year()

	if err := h.runService.CalculatePayroll(uint(id), tenantID, taxYear); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	run, _ := h.runService.GetPayrollRun(uint(id), tenantID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payroll calculated successfully",
		"data":    run,
	})
}

// GeneratePayslips generates payslips for finalized payroll
func (h *PayrollHandler) GeneratePayslips(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := utils.GetTenantID(c)
	taxYear := time.Now().Year()

	if err := h.runService.GeneratePayslips(uint(id), tenantID, taxYear); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payslips generated successfully",
	})
}

// GetPayrollEmployees retrieves employees in a payroll run
func (h *PayrollHandler) GetPayrollEmployees(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	search := c.Query("search")
	departmentID := c.Query("departmentId")
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	var hasChanges *bool
	if hc := c.Query("hasChanges"); hc != "" {
		val := hc == "true"
		hasChanges = &val
	}

	employees, total, err := h.runService.GetEmployees(uint(id), search, departmentID, hasChanges, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    employees,
		"meta": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// ============ New Workflow Handlers ============

// SubmitForHRReview submits payroll run for HR review
func (h *PayrollHandler) SubmitForHRReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := utils.GetTenantID(c)
	userID := utils.GetUserID(c)
	userName := utils.GetUserName(c)

	if err := h.runService.SubmitForHRReview(uint(id), tenantID, userID, userName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	run, _ := h.runService.GetPayrollRun(uint(id), tenantID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payroll run submitted for HR review",
		"data":    run,
	})
}

// HRReview handles HR review of payroll run
func (h *PayrollHandler) HRReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Status   string `json:"status" binding:"required"`
		Comments string `json:"comments"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	validStatuses := []string{"Approved", "Rejected", "Needs Correction"}
	isValid := false
	for _, valid := range validStatuses {
		if req.Status == valid {
			isValid = true
			break
		}
	}
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status. Must be: Approved, Rejected, or Needs Correction"})
		return
	}

	tenantID := utils.GetTenantID(c)
	userID := utils.GetUserID(c)
	userName := utils.GetUserName(c)

	status := models.WorkflowReviewStatus(req.Status)
	if err := h.runService.HRReview(uint(id), tenantID, userID, userName, status, req.Comments); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	run, _ := h.runService.GetPayrollRun(uint(id), tenantID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "HR review completed",
		"data":    run,
	})
}

// FinanceReview handles Finance review of payroll run
func (h *PayrollHandler) FinanceReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Status         string `json:"status" binding:"required"`
		Comments       string `json:"comments"`
		BudgetVerified bool   `json:"budgetVerified"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	validStatuses := []string{"Approved", "Rejected", "Needs Correction"}
	isValid := false
	for _, valid := range validStatuses {
		if req.Status == valid {
			isValid = true
			break
		}
	}
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status. Must be: Approved, Rejected, or Needs Correction"})
		return
	}

	tenantID := utils.GetTenantID(c)
	userID := utils.GetUserID(c)
	userName := utils.GetUserName(c)

	status := models.WorkflowReviewStatus(req.Status)
	if err := h.runService.FinanceReview(uint(id), tenantID, userID, userName, status, req.Comments, req.BudgetVerified); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	run, _ := h.runService.GetPayrollRun(uint(id), tenantID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Finance review completed",
		"data":    run,
	})
}

// ManagementApproval handles final management approval
func (h *PayrollHandler) ManagementApproval(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Approved bool   `json:"approved" binding:"required"`
		Comments string `json:"comments"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := utils.GetTenantID(c)
	userID := utils.GetUserID(c)
	userName := utils.GetUserName(c)

	if err := h.runService.ManagementApproval(uint(id), tenantID, userID, userName, req.Approved, req.Comments); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	run, _ := h.runService.GetPayrollRun(uint(id), tenantID)

	message := "Management approval completed"
	if !req.Approved {
		message = "Management approval rejected"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    run,
	})
}
