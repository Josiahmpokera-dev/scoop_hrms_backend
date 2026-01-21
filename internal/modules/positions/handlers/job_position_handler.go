package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/types"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type JobPositionHandler struct {
	service *services.JobPositionService
}

func NewJobPositionHandler() *JobPositionHandler {
	return &JobPositionHandler{
		service: services.NewJobPositionService(),
	}
}

// CreateJobPosition handles job position creation
func (h *JobPositionHandler) CreateJobPosition(c *gin.Context) {
	var req models.CreateJobPositionRequest
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

	position, err := h.service.CreateJobPosition(&req, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Job position created successfully", position)
}

// GetJobPosition handles getting a job position by ID (POST with ID in body)
func (h *JobPositionHandler) GetJobPosition(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	if req.ID == 0 {
		response.BadRequest(c, "ID is required in request body", nil)
		return
	}

	position, err := h.service.GetJobPositionByID(req.ID)
	if err != nil {
		response.NotFound(c, "Job position not found")
		return
	}

	// Convert to response format with reporting position name
	posMap := map[string]interface{}{
		"id":                  position.ID,
		"code":                position.Code,
		"title":               position.Title,
		"grade":               position.Grade,
		"department_id":       position.DepartmentID,
		"reports_to_position": nil,
		"reports_to_position_id": position.ReportsToPositionID,
		"budgeted_headcount":  position.BudgetedHeadcount,
		"current_headcount":   position.CurrentHeadcount,
		"employment_type":     position.EmploymentType,
		"key_competencies":    position.KeyCompetencies,
		"is_active":           position.IsActive,
		"created_at":          position.CreatedAt,
		"updated_at":          position.UpdatedAt,
	}

	// Load reporting position name if reports_to_position_id exists
	if position.ReportsToPositionID != nil {
		reportingPos, err := h.service.GetJobPositionByID(*position.ReportsToPositionID)
		if err == nil && reportingPos != nil {
			posMap["reports_to_position"] = reportingPos.Title
		}
	}

	response.Success(c, "Job position retrieved successfully", posMap)
}

// UpdateJobPosition handles job position updates (POST with ID in body)
func (h *JobPositionHandler) UpdateJobPosition(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
		models.UpdateJobPositionRequest
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	if req.ID == 0 {
		response.BadRequest(c, "ID is required in request body", nil)
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	position, err := h.service.UpdateJobPosition(req.ID, &req.UpdateJobPositionRequest, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Convert to response format with reporting position name
	posMap := map[string]interface{}{
		"id":                  position.ID,
		"code":                position.Code,
		"title":               position.Title,
		"grade":               position.Grade,
		"department_id":       position.DepartmentID,
		"reports_to_position": nil,
		"reports_to_position_id": position.ReportsToPositionID,
		"budgeted_headcount":  position.BudgetedHeadcount,
		"current_headcount":   position.CurrentHeadcount,
		"employment_type":     position.EmploymentType,
		"key_competencies":    position.KeyCompetencies,
		"is_active":           position.IsActive,
		"created_at":          position.CreatedAt,
		"updated_at":          position.UpdatedAt,
	}

	// Load reporting position name if reports_to_position_id exists
	if position.ReportsToPositionID != nil {
		reportingPos, err := h.service.GetJobPositionByID(*position.ReportsToPositionID)
		if err == nil && reportingPos != nil {
			posMap["reports_to_position"] = reportingPos.Title
		}
	}

	response.Success(c, "Job position updated successfully", posMap)
}

// DeleteJobPosition handles job position deletion (POST with ID in body)
func (h *JobPositionHandler) DeleteJobPosition(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	if req.ID == 0 {
		response.BadRequest(c, "ID is required in request body", nil)
		return
	}

	if err := h.service.DeleteJobPosition(req.ID); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Job position deleted successfully", nil)
}

// ListJobPositions handles listing job positions
func (h *JobPositionHandler) ListJobPositions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}

	tenantID := middleware.GetTenantID(c)
	
	// Build filters from query parameters
	filters := make(map[string]interface{})
	if isActive := c.Query("is_active"); isActive != "" {
		filters["is_active"] = isActive == "true"
	}
	if grade := c.Query("grade"); grade != "" {
		filters["grade"] = grade
	}
	if level := c.Query("level"); level != "" {
		if l, err := strconv.Atoi(level); err == nil {
			filters["level"] = l
		}
	}
	if departmentIDStr := c.Query("department_id"); departmentIDStr != "" {
		if departmentID, err := strconv.ParseUint(departmentIDStr, 10, 32); err == nil {
			filters["department_id"] = uint(departmentID)
		}
	}

	positions, total, err := h.service.ListJobPositions(tenantID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to list job positions", err.Error())
		return
	}

	// Convert to response format with reporting position name
	responseData := make([]map[string]interface{}, len(positions))
	for i, pos := range positions {
		posMap := map[string]interface{}{
			"id":                  pos.ID,
			"code":                pos.Code,
			"title":               pos.Title,
			"grade":               pos.Grade,
			"department_id":       pos.DepartmentID,
			"reports_to_position": nil, // Will be set below
			"reports_to_position_id": pos.ReportsToPositionID, // Keep ID for reference
			"budgeted_headcount":  pos.BudgetedHeadcount,
			"current_headcount":   pos.CurrentHeadcount,
			"employment_type":     pos.EmploymentType,
			"key_competencies":    pos.KeyCompetencies,
			"is_active":           pos.IsActive,
			"created_at":          pos.CreatedAt,
			"updated_at":          pos.UpdatedAt,
		}

		// Load reporting position name if reports_to_position_id exists
		if pos.ReportsToPositionID != nil {
			reportingPos, err := h.service.GetJobPositionByID(*pos.ReportsToPositionID)
			if err == nil && reportingPos != nil {
				posMap["reports_to_position"] = reportingPos.Title
			}
		}

		responseData[i] = posMap
	}

	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	response.SuccessWithMeta(c, "Job positions retrieved successfully", responseData, meta)
}

// GetPositionsByDepartment handles getting positions filtered by department ID
func (h *JobPositionHandler) GetPositionsByDepartment(c *gin.Context) {
	departmentIDStr := c.Param("department_id")
	if departmentIDStr == "" {
		response.BadRequest(c, "Department ID is required", nil)
		return
	}

	departmentID, err := strconv.ParseUint(departmentIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid department ID", nil)
		return
	}

	tenantID := middleware.GetTenantID(c)

	positions, err := h.service.GetPositionsByDepartment(uint(departmentID), tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Convert to response format with reporting position name
	responseData := make([]map[string]interface{}, len(positions))
	for i, pos := range positions {
		posMap := map[string]interface{}{
			"id":                    pos.ID,
			"code":                  pos.Code,
			"title":                 pos.Title,
			"grade":                 pos.Grade,
			"department_id":         pos.DepartmentID,
			"reports_to_position":   nil,
			"reports_to_position_id": pos.ReportsToPositionID,
			"budgeted_headcount":    pos.BudgetedHeadcount,
			"current_headcount":     pos.CurrentHeadcount,
			"employment_type":       pos.EmploymentType,
			"key_competencies":      pos.KeyCompetencies,
			"is_active":             pos.IsActive,
			"created_at":            pos.CreatedAt,
			"updated_at":            pos.UpdatedAt,
		}

		// Load reporting position name if reports_to_position_id exists
		if pos.ReportsToPositionID != nil {
			reportingPos, err := h.service.GetJobPositionByID(*pos.ReportsToPositionID)
			if err == nil && reportingPos != nil {
				posMap["reports_to_position"] = reportingPos.Title
			}
		}

		responseData[i] = posMap
	}

	response.Success(c, "Positions retrieved successfully", responseData)
}

// HandleAction handles POST-only action-based requests
func (h *JobPositionHandler) HandleAction(c *gin.Context) {
	var req types.APIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	switch req.Action {
	case "create":
		var createReq models.CreateJobPositionRequest
		if err := mapToStruct(req.Data, &createReq); err != nil {
			response.BadRequest(c, "Invalid data format", err.Error())
			return
		}
		position, err := h.service.CreateJobPosition(&createReq, tenantID, updatedBy)
		if err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Created(c, "Job position created successfully", position)

	case "read":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for read action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid job position ID", nil)
			return
		}
		position, err := h.service.GetJobPositionByID(uint(id))
		if err != nil {
			response.NotFound(c, "Job position not found")
			return
		}
		response.Success(c, "Job position retrieved successfully", position)

	case "update":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for update action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid job position ID", nil)
			return
		}
		var updateReq models.UpdateJobPositionRequest
		if err := mapToStruct(req.Data, &updateReq); err != nil {
			response.BadRequest(c, "Invalid data format", err.Error())
			return
		}
		position, err := h.service.UpdateJobPosition(uint(id), &updateReq, updatedBy)
		if err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Success(c, "Job position updated successfully", position)

	case "delete":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for delete action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid job position ID", nil)
			return
		}
		if err := h.service.DeleteJobPosition(uint(id)); err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Success(c, "Job position deleted successfully", nil)

	case "list":
		page := 1
		pageSize := 20
		if req.Pagination != nil {
			page = req.Pagination.GetPage()
			pageSize = req.Pagination.GetPageSize()
		}

		filters := req.Filters
		if filters == nil {
			filters = make(map[string]interface{})
		}

		positions, total, err := h.service.ListJobPositions(tenantID, page, pageSize, filters)
		if err != nil {
			response.InternalServerError(c, "Failed to list job positions", err.Error())
			return
		}

		meta := &response.Meta{
			Page:       page,
			PerPage:    pageSize,
			Total:      total,
			TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
		}

		response.SuccessWithMeta(c, "Job positions retrieved successfully", positions, meta)

	default:
		response.BadRequest(c, "Invalid action. Must be one of: create, read, update, delete, list", nil)
	}
}

// Helper function to map map[string]interface{} to struct
func mapToStruct(data map[string]interface{}, target interface{}) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, target)
}
