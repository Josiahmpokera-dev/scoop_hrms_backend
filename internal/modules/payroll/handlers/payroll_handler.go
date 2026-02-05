package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/services"
	"github.com/gin-gonic/gin"
)

// PayrollHandler handles payroll-related HTTP requests
type PayrollHandler struct {
	runService        *services.PayrollRunService
	structureService  *services.SalaryStructureService
	payslipService    *services.PayslipService
	loanService       *services.LoanService
	complianceService *services.ComplianceService
}

// NewPayrollHandler creates a new PayrollHandler instance
func NewPayrollHandler() *PayrollHandler {
	return &PayrollHandler{
		runService:        services.NewPayrollRunService(),
		structureService:  services.NewSalaryStructureService(),
		payslipService:    services.NewPayslipService(),
		loanService:       services.NewLoanService(),
		complianceService: services.NewComplianceService(),
	}
}

// Helper functions for extracting request data
func getTenantID(c *gin.Context) *uint {
	if tenantID, exists := c.Get("tenant_id"); exists {
		if tid, ok := tenantID.(uint); ok {
			return &tid
		}
	}
	return nil
}

func getUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uint); ok {
			return uid
		}
	}
	return 0
}

func getUserName(c *gin.Context) string {
	if userName, exists := c.Get("user_name"); exists {
		if name, ok := userName.(string); ok {
			return name
		}
	}
	return ""
}

func getPage(c *gin.Context) int {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	return page
}

func getPageSize(c *gin.Context) int {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return pageSize
}

// ============ Dashboard Handlers ============

// GetDashboard returns payroll dashboard metrics
func (h *PayrollHandler) GetDashboard(c *gin.Context) {
	tenantID := getTenantID(c)

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

// ============ Payroll Run Handlers ============

// ListPayrollRuns lists all payroll runs
func (h *PayrollHandler) ListPayrollRuns(c *gin.Context) {
	tenantID := getTenantID(c)
	status := c.Query("status")
	payYear, _ := strconv.Atoi(c.Query("year"))
	payMonth, _ := strconv.Atoi(c.Query("month"))
	page := getPage(c)
	pageSize := getPageSize(c)

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

	tenantID := getTenantID(c)
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
		PayMonth     int                  `json:"payMonth" binding:"required"`
		PayYear      int                  `json:"payYear" binding:"required"`
		PayFrequency models.PayFrequency  `json:"payFrequency"`
		RunName      string               `json:"runName"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := getTenantID(c)
	userID := getUserID(c)
	userName := getUserName(c)

	run := &models.PayrollRun{
		TenantID:      tenantID,
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

	tenantID := getTenantID(c)
	run, err := h.runService.GetPayrollRun(uint(id), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payroll run not found"})
		return
	}

	var req struct {
		RunName      string     `json:"runName"`
		CutoffDate   *time.Time `json:"cutoffDate"`
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

	tenantID := getTenantID(c)
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
	tenantID := getTenantID(c)
	run, err := h.runService.GetCurrentPayrollRun(tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no active payroll run"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
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

	tenantID := getTenantID(c)
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

	tenantID := getTenantID(c)
	userID := getUserID(c)
	userName := getUserName(c)

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

	tenantID := getTenantID(c)
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

	tenantID := getTenantID(c)
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

	tenantID := getTenantID(c)
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

	tenantID := getTenantID(c)
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
	page := getPage(c)
	pageSize := getPageSize(c)

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
