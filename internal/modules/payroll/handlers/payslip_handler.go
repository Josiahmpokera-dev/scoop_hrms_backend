package handlers

import (
	"net/http"
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/services"
	"github.com/gin-gonic/gin"
)

// PayslipHandler handles payslip-related HTTP requests
type PayslipHandler struct {
	service *services.PayslipService
}

// NewPayslipHandler creates a new handler instance
func NewPayslipHandler() *PayslipHandler {
	return &PayslipHandler{
		service: services.NewPayslipService(),
	}
}

// ListPayslips lists payslips with filters
func (h *PayslipHandler) ListPayslips(c *gin.Context) {
	tenantID := getTenantID(c)
	payMonth, _ := strconv.Atoi(c.Query("month"))
	payYear, _ := strconv.Atoi(c.Query("year"))
	status := c.Query("status")
	page := getPage(c)
	pageSize := getPageSize(c)

	var employeeID, departmentID *uint
	if empID := c.Query("employeeId"); empID != "" {
		id, _ := strconv.ParseUint(empID, 10, 32)
		eid := uint(id)
		employeeID = &eid
	}
	if deptID := c.Query("departmentId"); deptID != "" {
		id, _ := strconv.ParseUint(deptID, 10, 32)
		did := uint(id)
		departmentID = &did
	}

	payslips, total, err := h.service.ListPayslips(tenantID, payMonth, payYear, employeeID, departmentID, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    payslips,
		"meta": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetPayslip retrieves a payslip by ID
func (h *PayslipHandler) GetPayslip(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := getTenantID(c)
	payslipData, err := h.service.GetPayslipWithItems(uint(id), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payslip not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    payslipData,
	})
}

// GetPayslipSummary retrieves payslip summary for a period
func (h *PayslipHandler) GetPayslipSummary(c *gin.Context) {
	tenantID := getTenantID(c)
	payMonth, _ := strconv.Atoi(c.Query("month"))
	payYear, _ := strconv.Atoi(c.Query("year"))

	var employeeID *uint
	if empID := c.Query("employeeId"); empID != "" {
		id, _ := strconv.ParseUint(empID, 10, 32)
		eid := uint(id)
		employeeID = &eid
	}

	summary, err := h.service.GetSummary(tenantID, employeeID, payMonth, payYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

// DownloadPayslip downloads a payslip in specified format
func (h *PayslipHandler) DownloadPayslip(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	format := c.DefaultQuery("format", "pdf")
	tenantID := getTenantID(c)

	switch format {
	case "pdf":
		data, filename, err := h.service.GeneratePDF(uint(id), tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/pdf", data)

	case "xlsx", "excel":
		data, filename, err := h.service.GenerateExcel(uint(id), tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported format"})
	}
}

// BulkDownloadPayslips downloads multiple payslips as ZIP
func (h *PayslipHandler) BulkDownloadPayslips(c *gin.Context) {
	var req struct {
		PayslipIDs []uint `json:"payslipIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := getTenantID(c)
	data, filename, err := h.service.BulkDownloadPayslips(req.PayslipIDs, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/zip", data)
}

// SendPayslipEmail sends payslip to employee email
func (h *PayslipHandler) SendPayslipEmail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := getTenantID(c)
	if err := h.service.SendPayslipEmail(uint(id), tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payslip email sent successfully",
	})
}

// BulkEmailPayslips sends payslips to all employees for a run
func (h *PayslipHandler) BulkEmailPayslips(c *gin.Context) {
	runID, err := strconv.ParseUint(c.Param("runId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid run id"})
		return
	}

	tenantID := getTenantID(c)
	result, err := h.service.BulkEmailPayslips(uint(runID), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ReleasePayslips releases payslips for a payroll run
func (h *PayslipHandler) ReleasePayslips(c *gin.Context) {
	runID, err := strconv.ParseUint(c.Param("runId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid run id"})
		return
	}

	tenantID := getTenantID(c)
	if err := h.service.ReleasePayslips(uint(runID), tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payslips released successfully",
	})
}

// ============ Employee Self-Service Handlers ============

// GetMyPayslips retrieves current user's payslips
func (h *PayslipHandler) GetMyPayslips(c *gin.Context) {
	// Get employee ID from authenticated user
	employeeID := getUserID(c) // Assuming user ID maps to employee ID
	tenantID := getTenantID(c)
	page := getPage(c)
	pageSize := getPageSize(c)

	payslips, total, err := h.service.GetEmployeePayslipHistory(employeeID, tenantID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    payslips,
		"meta": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetMyLatestPayslip retrieves current user's latest payslip
func (h *PayslipHandler) GetMyLatestPayslip(c *gin.Context) {
	employeeID := getUserID(c)
	tenantID := getTenantID(c)

	payslips, _, err := h.service.GetEmployeePayslipHistory(employeeID, tenantID, 1, 1)
	if err != nil || len(payslips) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no payslips found"})
		return
	}

	payslipData, err := h.service.GetPayslipWithItems(payslips[0].ID, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    payslipData,
	})
}

// GetMySalarySlipSummary retrieves current user's salary summary
func (h *PayslipHandler) GetMySalarySlipSummary(c *gin.Context) {
	employeeID := getUserID(c)
	tenantID := getTenantID(c)
	payMonth, _ := strconv.Atoi(c.Query("month"))
	payYear, _ := strconv.Atoi(c.Query("year"))

	empID := employeeID
	summary, err := h.service.GetSummary(tenantID, &empID, payMonth, payYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

// DownloadMyPayslip downloads current user's payslip
func (h *PayslipHandler) DownloadMyPayslip(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	employeeID := getUserID(c)
	tenantID := getTenantID(c)
	format := c.DefaultQuery("format", "pdf")

	// Verify the payslip belongs to this employee
	payslip, err := h.service.GetPayslip(uint(id), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payslip not found"})
		return
	}

	if payslip.EmployeeID != employeeID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	switch format {
	case "pdf":
		data, filename, err := h.service.GeneratePDF(uint(id), tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/pdf", data)

	case "xlsx", "excel":
		data, filename, err := h.service.GenerateExcel(uint(id), tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported format"})
	}
}
