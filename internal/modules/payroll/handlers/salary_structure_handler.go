package handlers

import (
	"net/http"
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/services"
	"github.com/gin-gonic/gin"
)

// SalaryStructureHandler handles salary structure-related HTTP requests
type SalaryStructureHandler struct {
	service *services.SalaryStructureService
}

// NewSalaryStructureHandler creates a new handler instance
func NewSalaryStructureHandler() *SalaryStructureHandler {
	return &SalaryStructureHandler{
		service: services.NewSalaryStructureService(),
	}
}

// ============ Salary Structure Handlers ============

// ListSalaryStructures lists salary structure templates
func (h *SalaryStructureHandler) ListSalaryStructures(c *gin.Context) {
	tenantID := getTenantID(c)
	grade := c.Query("grade")
	location := c.Query("location")
	country := c.DefaultQuery("country", "Tanzania")
	page := getPage(c)
	pageSize := getPageSize(c)

	var isActive *bool
	if active := c.Query("isActive"); active != "" {
		val := active == "true"
		isActive = &val
	}

	structures, total, err := h.service.ListSalaryStructures(tenantID, isActive, grade, location, country, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    structures,
		"meta": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetSalaryStructure retrieves a salary structure by ID
func (h *SalaryStructureHandler) GetSalaryStructure(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := getTenantID(c)
	structure, err := h.service.GetSalaryStructure(uint(id), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "salary structure not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    structure,
	})
}

// CreateSalaryStructure creates a new salary structure template
func (h *SalaryStructureHandler) CreateSalaryStructure(c *gin.Context) {
	var req struct {
		TemplateName    string  `json:"templateName" binding:"required"`
		Grade           string  `json:"grade"`
		Location        string  `json:"location"`
		Country         string  `json:"country"`
		CTC             float64 `json:"ctc"`
		Basic           float64 `json:"basic" binding:"required"`
		HRA             float64 `json:"hra"`
		Transport       float64 `json:"transport"`
		Medical         float64 `json:"medical"`
		OtherAllowances float64 `json:"otherAllowances"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := getTenantID(c)
	userID := getUserID(c)

	structure := &models.SalaryStructure{
		TenantID:        tenantID,
		TemplateName:    req.TemplateName,
		Grade:           req.Grade,
		Location:        req.Location,
		Country:         req.Country,
		CTC:             req.CTC,
		Basic:           req.Basic,
		HRA:             req.HRA,
		Transport:       req.Transport,
		Medical:         req.Medical,
		OtherAllowances: req.OtherAllowances,
		IsActive:        true,
		UpdatedByID:     &userID,
	}

	if structure.Country == "" {
		structure.Country = "Tanzania"
	}

	if err := h.service.CreateSalaryStructure(structure); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Salary structure created successfully",
		"data":    structure,
	})
}

// UpdateSalaryStructure updates a salary structure
func (h *SalaryStructureHandler) UpdateSalaryStructure(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := getTenantID(c)
	structure, err := h.service.GetSalaryStructure(uint(id), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "salary structure not found"})
		return
	}

	var req struct {
		TemplateName    *string  `json:"templateName"`
		Grade           *string  `json:"grade"`
		Location        *string  `json:"location"`
		Basic           *float64 `json:"basic"`
		HRA             *float64 `json:"hra"`
		Transport       *float64 `json:"transport"`
		Medical         *float64 `json:"medical"`
		OtherAllowances *float64 `json:"otherAllowances"`
		IsActive        *bool    `json:"isActive"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.TemplateName != nil {
		structure.TemplateName = *req.TemplateName
	}
	if req.Grade != nil {
		structure.Grade = *req.Grade
	}
	if req.Location != nil {
		structure.Location = *req.Location
	}
	if req.Basic != nil {
		structure.Basic = *req.Basic
	}
	if req.HRA != nil {
		structure.HRA = *req.HRA
	}
	if req.Transport != nil {
		structure.Transport = *req.Transport
	}
	if req.Medical != nil {
		structure.Medical = *req.Medical
	}
	if req.OtherAllowances != nil {
		structure.OtherAllowances = *req.OtherAllowances
	}
	if req.IsActive != nil {
		structure.IsActive = *req.IsActive
	}

	userID := getUserID(c)
	structure.UpdatedByID = &userID

	if err := h.service.UpdateSalaryStructure(structure); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Salary structure updated successfully",
		"data":    structure,
	})
}

// DeleteSalaryStructure deletes a salary structure
func (h *SalaryStructureHandler) DeleteSalaryStructure(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := getTenantID(c)
	if err := h.service.DeleteSalaryStructure(uint(id), tenantID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Salary structure deleted successfully",
	})
}

// SimulateSalary simulates salary calculation from CTC
func (h *SalaryStructureHandler) SimulateSalary(c *gin.Context) {
	var req struct {
		CTC float64 `json:"ctc" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := getTenantID(c)
	result := h.service.SimulateSalary(req.CTC, tenantID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ============ Salary Component Handlers ============

// ListSalaryComponents lists salary components
func (h *SalaryStructureHandler) ListSalaryComponents(c *gin.Context) {
	tenantID := getTenantID(c)
	componentType := c.Query("type")
	country := c.DefaultQuery("country", "Tanzania")
	page := getPage(c)
	pageSize := getPageSize(c)

	var isActive, isStatutory *bool
	if active := c.Query("isActive"); active != "" {
		val := active == "true"
		isActive = &val
	}
	if statutory := c.Query("isStatutory"); statutory != "" {
		val := statutory == "true"
		isStatutory = &val
	}

	components, total, err := h.service.ListSalaryComponents(tenantID, componentType, isActive, isStatutory, country, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    components,
		"meta": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetSalaryComponent retrieves a salary component by ID
func (h *SalaryStructureHandler) GetSalaryComponent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := getTenantID(c)
	component, err := h.service.GetSalaryComponent(uint(id), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "salary component not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    component,
	})
}

// CreateSalaryComponent creates a new salary component
func (h *SalaryStructureHandler) CreateSalaryComponent(c *gin.Context) {
	var req struct {
		ComponentName     string                   `json:"componentName" binding:"required"`
		ComponentCode     string                   `json:"componentCode" binding:"required"`
		ComponentType     models.ComponentType     `json:"componentType" binding:"required"`
		CalculationType   models.CalculationType   `json:"calculationType" binding:"required"`
		DefaultAmount     *float64                 `json:"defaultAmount"`
		DefaultPercentage *float64                 `json:"defaultPercentage"`
		Formula           *string                  `json:"formula"`
		IsTaxable         bool                     `json:"isTaxable"`
		IsStatutory       bool                     `json:"isStatutory"`
		IsRecurring       bool                     `json:"isRecurring"`
		ApplicableFor     string                   `json:"applicableFor"`
		DisplayInPayslip  bool                     `json:"displayInPayslip"`
		Country           string                   `json:"country"`
		SortOrder         int                      `json:"sortOrder"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := getTenantID(c)

	component := &models.SalaryComponent{
		TenantID:          tenantID,
		ComponentName:     req.ComponentName,
		ComponentCode:     req.ComponentCode,
		ComponentType:     req.ComponentType,
		CalculationType:   req.CalculationType,
		DefaultAmount:     req.DefaultAmount,
		DefaultPercentage: req.DefaultPercentage,
		Formula:           req.Formula,
		IsTaxable:         req.IsTaxable,
		IsStatutory:       req.IsStatutory,
		IsRecurring:       req.IsRecurring,
		ApplicableFor:     req.ApplicableFor,
		DisplayInPayslip:  req.DisplayInPayslip,
		Country:           req.Country,
		SortOrder:         req.SortOrder,
		IsActive:          true,
	}

	if component.Country == "" {
		component.Country = "Tanzania"
	}
	if component.ApplicableFor == "" {
		component.ApplicableFor = `["All"]`
	}

	if err := h.service.CreateSalaryComponent(component); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Salary component created successfully",
		"data":    component,
	})
}

// UpdateSalaryComponent updates a salary component
func (h *SalaryStructureHandler) UpdateSalaryComponent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := getTenantID(c)
	component, err := h.service.GetSalaryComponent(uint(id), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "salary component not found"})
		return
	}

	var req struct {
		ComponentName     *string  `json:"componentName"`
		DefaultAmount     *float64 `json:"defaultAmount"`
		DefaultPercentage *float64 `json:"defaultPercentage"`
		IsTaxable         *bool    `json:"isTaxable"`
		DisplayInPayslip  *bool    `json:"displayInPayslip"`
		IsActive          *bool    `json:"isActive"`
		SortOrder         *int     `json:"sortOrder"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ComponentName != nil {
		component.ComponentName = *req.ComponentName
	}
	if req.DefaultAmount != nil {
		component.DefaultAmount = req.DefaultAmount
	}
	if req.DefaultPercentage != nil {
		component.DefaultPercentage = req.DefaultPercentage
	}
	if req.IsTaxable != nil {
		component.IsTaxable = *req.IsTaxable
	}
	if req.DisplayInPayslip != nil {
		component.DisplayInPayslip = *req.DisplayInPayslip
	}
	if req.IsActive != nil {
		component.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		component.SortOrder = *req.SortOrder
	}

	if err := h.service.UpdateSalaryComponent(component); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Salary component updated successfully",
		"data":    component,
	})
}

// DeleteSalaryComponent deletes a salary component
func (h *SalaryStructureHandler) DeleteSalaryComponent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tenantID := getTenantID(c)
	if err := h.service.DeleteSalaryComponent(uint(id), tenantID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Salary component deleted successfully",
	})
}
