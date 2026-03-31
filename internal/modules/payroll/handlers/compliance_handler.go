package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

// ComplianceHandler handles compliance-related HTTP requests
type ComplianceHandler struct {
	service *services.ComplianceService
}

// NewComplianceHandler creates a new handler instance
func NewComplianceHandler() *ComplianceHandler {
	return &ComplianceHandler{
		service: services.NewComplianceService(),
	}
}

// GetTaxSlabs retrieves tax slabs
func (h *ComplianceHandler) GetTaxSlabs(c *gin.Context) {
	tenantID := utils.GetTenantID(c)
	country := c.DefaultQuery("country", "Tanzania")
	taxYear, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))

	slabs, err := h.service.GetTaxSlabs(country, taxYear, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    slabs,
	})
}

// GetStatutoryRules retrieves statutory rules
func (h *ComplianceHandler) GetStatutoryRules(c *gin.Context) {
	tenantID := utils.GetTenantID(c)
	country := c.DefaultQuery("country", "Tanzania")

	rules, err := h.service.GetStatutoryRules(country, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rules,
	})
}

// GetNHIFSchedule retrieves NHIF schedule
func (h *ComplianceHandler) GetNHIFSchedule(c *gin.Context) {
	tenantID := utils.GetTenantID(c)
	country := c.DefaultQuery("country", "Tanzania")

	schedule, err := h.service.GetNHIFSchedule(country, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    schedule,
	})
}

// GetComplianceSummary retrieves compliance summary for a period
func (h *ComplianceHandler) GetComplianceSummary(c *gin.Context) {
	tenantID := utils.GetTenantID(c)
	payMonth, _ := strconv.Atoi(c.Query("month"))
	payYear, _ := strconv.Atoi(c.Query("year"))

	if payMonth == 0 {
		payMonth = int(time.Now().Month())
	}
	if payYear == 0 {
		payYear = time.Now().Year()
	}

	summary, err := h.service.GetComplianceSummary(tenantID, payMonth, payYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

// SimulateTax performs tax simulation
func (h *ComplianceHandler) SimulateTax(c *gin.Context) {
	var req struct {
		GrossSalary float64 `json:"grossSalary" binding:"required"`
		TaxYear     int     `json:"taxYear"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.TaxYear == 0 {
		req.TaxYear = time.Now().Year()
	}

	tenantID := utils.GetTenantID(c)
	result := h.service.SimulateTax(req.GrossSalary, req.TaxYear, tenantID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetAllComplianceRules retrieves all compliance rules
func (h *ComplianceHandler) GetAllComplianceRules(c *gin.Context) {
	tenantID := utils.GetTenantID(c)
	country := c.DefaultQuery("country", "Tanzania")

	rules, err := h.service.GetAllComplianceRules(country, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rules,
	})
}

// GenerateComplianceReturn generates compliance return for filing
func (h *ComplianceHandler) GenerateComplianceReturn(c *gin.Context) {
	tenantID := utils.GetTenantID(c)
	payMonth, _ := strconv.Atoi(c.Query("month"))
	payYear, _ := strconv.Atoi(c.Query("year"))
	returnType := c.Query("type") // PAYE, NSSF, NHIF, SDL, WCF

	if payMonth == 0 {
		payMonth = int(time.Now().Month())
	}
	if payYear == 0 {
		payYear = time.Now().Year()
	}

	result, err := h.service.GenerateComplianceReturn(tenantID, payMonth, payYear, returnType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// CalculateMonthlyStatutory calculates statutory contributions for a month
func (h *ComplianceHandler) CalculateMonthlyStatutory(c *gin.Context) {
	var req struct {
		TotalGross    float64 `json:"totalGross" binding:"required"`
		EmployeeCount int     `json:"employeeCount" binding:"required"`
		PayMonth      int     `json:"payMonth"`
		PayYear      int     `json:"payYear"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.PayMonth == 0 {
		req.PayMonth = int(time.Now().Month())
	}
	if req.PayYear == 0 {
		req.PayYear = time.Now().Year()
	}

	tenantID := utils.GetTenantID(c)
	result, err := h.service.CalculateMonthlyStatutory(tenantID, req.PayMonth, req.PayYear, req.TotalGross, req.EmployeeCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// SeedComplianceData seeds Tanzania compliance data
func (h *ComplianceHandler) SeedComplianceData(c *gin.Context) {
	if err := h.service.SeedTanzaniaCompliance(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tanzania compliance data seeded successfully",
	})
}

// GetCompliancePayments retrieves compliance payments for a period
func (h *ComplianceHandler) GetCompliancePayments(c *gin.Context) {
	tenantID := utils.GetTenantID(c)
	payMonth, _ := strconv.Atoi(c.Query("month"))
	payYear, _ := strconv.Atoi(c.Query("year"))

	if payMonth == 0 {
		payMonth = int(time.Now().Month())
	}
	if payYear == 0 {
		payYear = time.Now().Year()
	}

	payments, err := h.service.GetCompliancePayments(tenantID, payMonth, payYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    payments,
	})
}
