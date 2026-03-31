package handlers

import (
	"net/http"
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/services"
	"github.com/gin-gonic/gin"
)

// PayrollConfigHandler handles payroll configuration HTTP requests
type PayrollConfigHandler struct {
	service *services.PayrollConfigService
}

// NewPayrollConfigHandler creates a new PayrollConfigHandler instance
func NewPayrollConfigHandler() *PayrollConfigHandler {
	return &PayrollConfigHandler{
		service: services.NewPayrollConfigService(),
	}
}

// GetAllConfigurations returns all payroll configurations
func (h *PayrollConfigHandler) GetAllConfigurations(c *gin.Context) {
	configs, err := h.service.GetAllConfigurations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to load payroll configurations",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    configs,
	})
}

// GetConfiguration returns a single configuration by key
func (h *PayrollConfigHandler) GetConfiguration(c *gin.Context) {
	configKey := c.Param("key")
	if configKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Configuration key is required",
		})
		return
	}

	config, err := h.service.GetConfiguration(configKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Configuration not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    config,
	})
}

// CalculatePAYE calculates PAYE tax for given taxable pay
func (h *PayrollConfigHandler) CalculatePAYE(c *gin.Context) {
	taxablePayStr := c.Query("taxable_pay")
	if taxablePayStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "taxable_pay parameter is required",
		})
		return
	}

	taxablePay, err := strconv.ParseFloat(taxablePayStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid taxable_pay value",
			"details": "Must be a valid number",
		})
		return
	}

	result, err := h.service.CalculatePAYE(taxablePay)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to calculate PAYE",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// CalculateNSSF calculates NSSF contributions
func (h *PayrollConfigHandler) CalculateNSSF(c *gin.Context) {
	grossSalaryStr := c.Query("gross")
	if grossSalaryStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "gross parameter is required",
		})
		return
	}

	grossSalary, err := strconv.ParseFloat(grossSalaryStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid gross value",
			"details": "Must be a valid number",
		})
		return
	}

	result, err := h.service.CalculateNSSF(grossSalary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to calculate NSSF",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// CalculateCTC calculates Cost to Company
func (h *PayrollConfigHandler) CalculateCTC(c *gin.Context) {
	grossSalaryStr := c.Query("gross")
	if grossSalaryStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "gross parameter is required",
		})
		return
	}

	grossSalary, err := strconv.ParseFloat(grossSalaryStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid gross value",
			"details": "Must be a valid number",
		})
		return
	}

	result, err := h.service.CalculateCTC(grossSalary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to calculate CTC",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetPAYEBands returns all PAYE tax bands
func (h *PayrollConfigHandler) GetPAYEBands(c *gin.Context) {
	bands, err := h.service.GetPAYEBands()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to load PAYE bands",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bands,
	})
}

// GetNSSFRates returns NSSF rates
func (h *PayrollConfigHandler) GetNSSFRates(c *gin.Context) {
	rates, err := h.service.GetNSSFRates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to load NSSF rates",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rates,
	})
}

// GetEmployerContributions returns employer contribution rates
func (h *PayrollConfigHandler) GetEmployerContributions(c *gin.Context) {
	contributions, err := h.service.GetEmployerContributions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to load employer contributions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    contributions,
	})
}

// GetFormulaConfigurations returns formula configurations
func (h *PayrollConfigHandler) GetFormulaConfigurations(c *gin.Context) {
	formulas, err := h.service.GetFormulaConfigurations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to load formula configurations",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    formulas,
	})
}

// BlockModificationEndpoints blocks any POST, PUT, DELETE operations
func (h *PayrollConfigHandler) BlockModificationEndpoints(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{
		"error":   "Configuration modifications are not allowed",
		"details": "These rates are set by Tanzanian law and can only be changed by a system administrator via database migration",
	})
}

// CreateConfiguration blocks creation of configurations
func (h *PayrollConfigHandler) CreateConfiguration(c *gin.Context) {
	h.BlockModificationEndpoints(c)
}

// UpdateConfiguration blocks updates to configurations
func (h *PayrollConfigHandler) UpdateConfiguration(c *gin.Context) {
	h.BlockModificationEndpoints(c)
}

// DeleteConfiguration blocks deletion of configurations
func (h *PayrollConfigHandler) DeleteConfiguration(c *gin.Context) {
	h.BlockModificationEndpoints(c)
}
