package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/teams/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/teams/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/types"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	service *services.TeamService
}

func NewTeamHandler() *TeamHandler {
	return &TeamHandler{
		service: services.NewTeamService(),
	}
}

// CreateTeam handles team creation
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req models.CreateTeamRequest
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

	team, err := h.service.CreateTeam(&req, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Team created successfully", team)
}

// GetTeam handles getting a team by ID
func (h *TeamHandler) GetTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid team ID", nil)
		return
	}

	team, err := h.service.GetTeamByID(uint(id))
	if err != nil {
		response.NotFound(c, "Team not found")
		return
	}

	response.Success(c, "Team retrieved successfully", team)
}

// UpdateTeam handles team updates
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid team ID", nil)
		return
	}

	var req models.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	team, err := h.service.UpdateTeam(uint(id), &req, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Team updated successfully", team)
}

// DeleteTeam handles team deletion
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid team ID", nil)
		return
	}

	if err := h.service.DeleteTeam(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Team deleted successfully", nil)
}

// ListTeams handles listing teams
func (h *TeamHandler) ListTeams(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}

	tenantID := middleware.GetTenantID(c)
	var departmentID *uint
	if deptIDStr := c.Query("department_id"); deptIDStr != "" {
		if deptID, err := strconv.ParseUint(deptIDStr, 10, 32); err == nil {
			deptIDUint := uint(deptID)
			departmentID = &deptIDUint
		}
	}
	
	// Build filters from query parameters
	filters := make(map[string]interface{})
	if isActive := c.Query("is_active"); isActive != "" {
		filters["is_active"] = isActive == "true"
	}
	if teamType := c.Query("team_type"); teamType != "" {
		filters["team_type"] = teamType
	}

	teams, total, err := h.service.ListTeams(tenantID, departmentID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to list teams", err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	response.SuccessWithMeta(c, "Teams retrieved successfully", teams, meta)
}

// GetTeamsByDepartment handles getting teams by department
func (h *TeamHandler) GetTeamsByDepartment(c *gin.Context) {
	deptIDStr := c.Param("department_id")
	deptID, err := strconv.ParseUint(deptIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid department ID", nil)
		return
	}

	teams, err := h.service.GetTeamsByDepartment(uint(deptID))
	if err != nil {
		response.InternalServerError(c, "Failed to get teams", err.Error())
		return
	}

	response.Success(c, "Teams retrieved successfully", teams)
}

// HandleAction handles POST-only action-based requests
func (h *TeamHandler) HandleAction(c *gin.Context) {
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

	var departmentID *uint
	if deptID, ok := req.Filters["department_id"]; ok {
		if deptIDFloat, ok := deptID.(float64); ok {
			deptIDUint := uint(deptIDFloat)
			departmentID = &deptIDUint
		}
	}

	switch req.Action {
	case "create":
		var createReq models.CreateTeamRequest
		if err := mapToStruct(req.Data, &createReq); err != nil {
			response.BadRequest(c, "Invalid data format", err.Error())
			return
		}
		team, err := h.service.CreateTeam(&createReq, tenantID, updatedBy)
		if err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Created(c, "Team created successfully", team)

	case "read":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for read action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid team ID", nil)
			return
		}
		team, err := h.service.GetTeamByID(uint(id))
		if err != nil {
			response.NotFound(c, "Team not found")
			return
		}
		response.Success(c, "Team retrieved successfully", team)

	case "update":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for update action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid team ID", nil)
			return
		}
		var updateReq models.UpdateTeamRequest
		if err := mapToStruct(req.Data, &updateReq); err != nil {
			response.BadRequest(c, "Invalid data format", err.Error())
			return
		}
		team, err := h.service.UpdateTeam(uint(id), &updateReq, updatedBy)
		if err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Success(c, "Team updated successfully", team)

	case "delete":
		if req.ID == nil {
			response.BadRequest(c, "ID is required for delete action", nil)
			return
		}
		id, err := strconv.ParseUint(*req.ID, 10, 32)
		if err != nil {
			response.BadRequest(c, "Invalid team ID", nil)
			return
		}
		if err := h.service.DeleteTeam(uint(id)); err != nil {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.Success(c, "Team deleted successfully", nil)

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

		teams, total, err := h.service.ListTeams(tenantID, departmentID, page, pageSize, filters)
		if err != nil {
			response.InternalServerError(c, "Failed to list teams", err.Error())
			return
		}

		meta := &response.Meta{
			Page:       page,
			PerPage:    pageSize,
			Total:      total,
			TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
		}

		response.SuccessWithMeta(c, "Teams retrieved successfully", teams, meta)

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
