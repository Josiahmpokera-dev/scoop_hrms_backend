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

// LoanHandler handles loan-related HTTP requests
type LoanHandler struct {
	service *services.LoanService
}

// NewLoanHandler creates a new handler instance
func NewLoanHandler() *LoanHandler {
	return &LoanHandler{
		service: services.NewLoanService(),
	}
}

// ListLoans lists loans with filters
func (h *LoanHandler) ListLoans(c *gin.Context) {
	tenantID := utils.GetTenantID(c)
	status := c.Query("status")
	loanType := c.Query("type")
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	var employeeID *uint
	if empID := c.Query("employeeId"); empID != "" {
		id, _ := strconv.ParseUint(empID, 10, 32)
		eid := uint(id)
		employeeID = &eid
	}

	loans, total, err := h.service.ListLoans(tenantID, employeeID, status, loanType, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    loans,
		"meta": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetLoan retrieves a loan by ID
func (h *LoanHandler) GetLoan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := utils.GetTenantID(c)
	loanData, err := h.service.GetLoanWithRepayments(uint(id), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    loanData,
	})
}

// CreateLoan creates a new loan application
func (h *LoanHandler) CreateLoan(c *gin.Context) {
	var req struct {
		EmployeeID   uint            `json:"employeeId" binding:"required"`
		EmployeeCode string          `json:"empId"`
		EmployeeName string          `json:"employeeName"`
		Department   string          `json:"department"`
		DepartmentID *uint           `json:"departmentId"`
		Designation  *string         `json:"designation"`
		LoanType     models.LoanType `json:"loanType" binding:"required"`
		Amount       float64         `json:"amount" binding:"required"`
		InterestRate float64         `json:"interestRate"`
		Tenure       int             `json:"tenure" binding:"required"`
		Purpose      *string         `json:"purpose"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan := &models.Loan{
		EmployeeID:   req.EmployeeID,
		EmployeeCode: req.EmployeeCode,
		EmployeeName: req.EmployeeName,
		Department:   req.Department,
		DepartmentID: req.DepartmentID,
		Designation:  req.Designation,
		LoanType:     req.LoanType,
		Amount:       req.Amount,
		InterestRate: req.InterestRate,
		Tenure:       req.Tenure,
		Purpose:      req.Purpose,
	}

	if err := h.service.CreateLoan(loan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Loan application created successfully",
		"data":    loan,
	})
}

// ApproveLoan approves a loan application
func (h *LoanHandler) ApproveLoan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		DisbursementDate *time.Time `json:"disbursementDate"`
		Remarks          string     `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := utils.GetTenantID(c)
	approverID := utils.GetUserID(c)
	approverName := utils.GetUserName(c)

	if err := h.service.ApproveLoan(uint(id), tenantID, approverID, approverName, req.DisbursementDate, req.Remarks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan, _ := h.service.GetLoan(uint(id), tenantID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Loan approved successfully",
		"data":    loan,
	})
}

// RejectLoan rejects a loan application
func (h *LoanHandler) RejectLoan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := utils.GetTenantID(c)
	rejecterID := utils.GetUserID(c)
	rejecterName := utils.GetUserName(c)

	if err := h.service.RejectLoan(uint(id), tenantID, rejecterID, rejecterName, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Loan rejected",
	})
}

// GetRepaymentSchedule retrieves repayment schedule for a loan
func (h *LoanHandler) GetRepaymentSchedule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	schedule, err := h.service.GetRepaymentSchedule(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    schedule,
	})
}

// RecordRepayment records a loan repayment
func (h *LoanHandler) RecordRepayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Installment int     `json:"installment" binding:"required"`
		PaidAmount  float64 `json:"paidAmount" binding:"required"`
		PayslipID   *uint   `json:"payslipId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RecordRepayment(uint(id), req.Installment, req.PaidAmount, req.PayslipID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Repayment recorded successfully",
	})
}

// GetLoanSummary retrieves loan summary
func (h *LoanHandler) GetLoanSummary(c *gin.Context) {
	tenantID := utils.GetTenantID(c)

	var employeeID *uint
	if empID := c.Query("employeeId"); empID != "" {
		id, _ := strconv.ParseUint(empID, 10, 32)
		eid := uint(id)
		employeeID = &eid
	}

	summary, err := h.service.GetLoanSummary(tenantID, employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

// CalculateEMI calculates EMI for given parameters
func (h *LoanHandler) CalculateEMI(c *gin.Context) {
	var req struct {
		Principal    float64 `json:"principal" binding:"required"`
		InterestRate float64 `json:"interestRate"`
		Tenure       int     `json:"tenure" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emi := h.service.CalculateEMI(req.Principal, req.InterestRate, req.Tenure)
	totalPayable := emi * float64(req.Tenure)
	totalInterest := totalPayable - req.Principal

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"principal":     req.Principal,
			"interestRate":  req.InterestRate,
			"tenure":        req.Tenure,
			"emi":           emi,
			"totalPayable":  totalPayable,
			"totalInterest": totalInterest,
		},
	})
}

// ============ Employee Self-Service Handlers ============

// GetMyLoans retrieves current user's loans
func (h *LoanHandler) GetMyLoans(c *gin.Context) {
	employeeID := utils.GetUserID(c)
	tenantID := utils.GetTenantID(c)
	status := c.Query("status")
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	empID := employeeID
	loans, total, err := h.service.ListLoans(tenantID, &empID, status, "", page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    loans,
		"meta": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// ApplyForLoan allows employee to apply for a loan
func (h *LoanHandler) ApplyForLoan(c *gin.Context) {
	employeeID := utils.GetUserID(c)

	var req struct {
		LoanType models.LoanType `json:"loanType" binding:"required"`
		Amount   float64         `json:"amount" binding:"required"`
		Tenure   int             `json:"tenure" binding:"required"`
		Purpose  *string         `json:"purpose"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan := &models.Loan{
		EmployeeID: employeeID,
		LoanType:   req.LoanType,
		Amount:     req.Amount,
		Tenure:     req.Tenure,
		Purpose:    req.Purpose,
	}

	if err := h.service.CreateLoan(loan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Loan application submitted successfully",
		"data":    loan,
	})
}
