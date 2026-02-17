package handlers

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/gin-gonic/gin"
)

type GoalHandler struct {
	goalService *services.GoalService
}

func NewGoalHandler() *GoalHandler {
	return &GoalHandler{goalService: services.NewGoalService()}
}

func (h *GoalHandler) GetDashboardStats(c *gin.Context) {
	userID, _ := c.Get("user_id")
	callerID := userID.(uint)
	data, err := h.goalService.GetDashboardStats(callerID)
	if err != nil {
		response.InternalServerError(c, "Failed to get dashboard stats", err.Error())
		return
	}
	response.Success(c, "Dashboard stats retrieved successfully", data)
}

func (h *GoalHandler) GetUpcomingActions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	callerID := userID.(uint)
	data, err := h.goalService.GetUpcomingActions(callerID)
	if err != nil {
		response.InternalServerError(c, "Failed to get upcoming actions", err.Error())
		return
	}
	response.Success(c, "Upcoming actions retrieved successfully", data)
}

func (h *GoalHandler) ListGoals(c *gin.Context) {
	level := c.Query("level")
	status := c.Query("status")
	ownerIDStr := c.Query("owner_id")
	department := c.Query("department")
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	userID, _ := c.Get("user_id")
	callerID := userID.(uint)
	isHROrAdmin := false
	if v, exists := c.Get("is_hr_or_admin"); exists {
		if b, ok := v.(bool); ok {
			isHROrAdmin = b
		}
	}

	var ownerID *uint
	if ownerIDStr != "" {
		if oid, err := strconv.ParseUint(ownerIDStr, 10, 32); err == nil {
			oidU := uint(oid)
			ownerID = &oidU
		}
	}
	restrictToOwner := !isHROrAdmin
	if restrictToOwner {
		ownerID = &callerID
	}

	goals, total, err := h.goalService.ListGoals(ownerID, level, status, department, search, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to list goals", err.Error())
		return
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &response.Meta{Page: page, PerPage: pageSize, Total: total, TotalPages: totalPages}
	response.SuccessWithMeta(c, "Goals retrieved successfully", goals, meta)
}

func (h *GoalHandler) GetGoalStats(c *gin.Context) {
	userID, _ := c.Get("user_id")
	callerID := userID.(uint)
	isHROrAdmin := false
	if v, exists := c.Get("is_hr_or_admin"); exists {
		if b, ok := v.(bool); ok {
			isHROrAdmin = b
		}
	}
	var ownerID *uint
	if !isHROrAdmin {
		ownerID = &callerID
	}
	data, err := h.goalService.GetGoalStats(ownerID)
	if err != nil {
		response.InternalServerError(c, "Failed to get goal stats", err.Error())
		return
	}
	response.Success(c, "Goal stats retrieved successfully", data)
}

func (h *GoalHandler) GetGoal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid goal ID", nil)
		return
	}
	goal, err := h.goalService.GetGoal(uint(id))
	if err != nil {
		response.NotFound(c, "Goal not found")
		return
	}
	response.Success(c, "Goal retrieved successfully", goal)
}

func (h *GoalHandler) CreateGoal(c *gin.Context) {
	var req models.CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	userName := ""
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(*userModels.User); ok {
			userName = u.FirstName + " " + u.LastName
		}
	}
	goal, err := h.goalService.CreateGoal(req.OwnerID, userName, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "Goal created successfully", goal)
}

func (h *GoalHandler) UpdateGoal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid goal ID", nil)
		return
	}
	var req models.UpdateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	goal, err := h.goalService.UpdateGoal(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Goal updated successfully", goal)
}

func (h *GoalHandler) DeleteGoal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid goal ID", nil)
		return
	}
	if err := h.goalService.DeleteGoal(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Goal deleted successfully", nil)
}

func (h *GoalHandler) CreateCheckIn(c *gin.Context) {
	goalIDStr := c.Param("id")
	goalID, err := strconv.ParseUint(goalIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid goal ID", nil)
		return
	}
	var req models.CreateCheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	callerID := userID.(uint)
	userName := ""
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(*userModels.User); ok {
			userName = u.FirstName + " " + u.LastName
		}
	}
	req.GoalID = uint(goalID)
	req.UpdatedByID = callerID
	req.UpdatedByName = userName
	checkIn, err := h.goalService.CreateCheckIn(uint(goalID), callerID, userName, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "Check-in created successfully", checkIn)
}

func (h *GoalHandler) CreateKeyResult(c *gin.Context) {
	goalIDStr := c.Param("id")
	goalID, err := strconv.ParseUint(goalIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid goal ID", nil)
		return
	}
	var req models.CreateKeyResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	req.GoalID = uint(goalID)
	kr, err := h.goalService.CreateKeyResult(uint(goalID), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "Key result created successfully", kr)
}

func (h *GoalHandler) UpdateKeyResult(c *gin.Context) {
	goalIDStr := c.Param("id")
	krIDStr := c.Param("krId")
	_, err1 := strconv.ParseUint(goalIDStr, 10, 32)
	krID, err2 := strconv.ParseUint(krIDStr, 10, 32)
	if err1 != nil || err2 != nil {
		response.BadRequest(c, "Invalid goal or key result ID", nil)
		return
	}
	var req models.UpdateKeyResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	kr, err := h.goalService.UpdateKeyResult(uint(krID), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Key result updated successfully", kr)
}

func (h *GoalHandler) DeleteKeyResult(c *gin.Context) {
	goalIDStr := c.Param("id")
	krIDStr := c.Param("krId")
	_, err1 := strconv.ParseUint(goalIDStr, 10, 32)
	krID, err2 := strconv.ParseUint(krIDStr, 10, 32)
	if err1 != nil || err2 != nil {
		response.BadRequest(c, "Invalid goal or key result ID", nil)
		return
	}
	if err := h.goalService.DeleteKeyResult(uint(krID)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Key result deleted successfully", nil)
}

func (h *GoalHandler) GetAlignmentMap(c *gin.Context) {
	data, err := h.goalService.GetAlignmentMap()
	if err != nil {
		response.InternalServerError(c, "Failed to get alignment map", err.Error())
		return
	}
	response.Success(c, "Alignment map retrieved successfully", data)
}

func (h *GoalHandler) SubmitForApproval(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid goal ID", nil)
		return
	}
	var body struct {
		Comments string `json:"comments"`
	}
	_ = c.ShouldBindJSON(&body)
	goal, err := h.goalService.SubmitForApproval(uint(id), body.Comments)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Goal submitted for approval successfully", goal)
}

func (h *GoalHandler) ApproveGoal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid goal ID", nil)
		return
	}
	var req models.ApproveGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	req.GoalID = uint(id)
	userID, _ := c.Get("user_id")
	callerID := userID.(uint)
	userName := ""
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(*userModels.User); ok {
			userName = u.FirstName + " " + u.LastName
		}
	}
	action := "reject"
	if req.Approved {
		action = "approve"
	}
	comments := ""
	if req.Comments != nil {
		comments = *req.Comments
	}
	goal, err := h.goalService.ApproveGoal(uint(id), callerID, userName, action, comments)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Goal approval updated successfully", goal)
}

func (h *GoalHandler) RequestCompletion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid goal ID", nil)
		return
	}
	var req models.RequestCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	req.GoalID = uint(id)
	evidence := ""
	if req.CompletionEvidence != nil {
		evidence = *req.CompletionEvidence
	}
	evidenceURL := ""
	if req.CompletionEvidenceURL != nil {
		evidenceURL = *req.CompletionEvidenceURL
	}
	goal, err := h.goalService.RequestCompletion(uint(id), evidence, evidenceURL, 100)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Completion requested successfully", goal)
}

func (h *GoalHandler) VerifyCompletion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid goal ID", nil)
		return
	}
	var req models.VerifyCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	req.GoalID = uint(id)
	userID, _ := c.Get("user_id")
	callerID := userID.(uint)
	userName := ""
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(*userModels.User); ok {
			userName = u.FirstName + " " + u.LastName
		}
	}
	action := "reject"
	if req.Verified {
		action = "verify"
	}
	comments := ""
	if req.Comments != nil {
		comments = *req.Comments
	}
	goal, err := h.goalService.VerifyCompletion(uint(id), callerID, userName, action, req.FinalRating, comments)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Completion verification updated successfully", goal)
}

func (h *GoalHandler) AssignGoal(c *gin.Context) {
	var req models.AssignGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	callerID := userID.(uint)
	userName := ""
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(*userModels.User); ok {
			userName = u.FirstName + " " + u.LastName
		}
	}
	goal, err := h.goalService.AssignGoal(callerID, userName, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Goal assigned successfully", goal)
}

func (h *GoalHandler) LinkProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid goal ID", nil)
		return
	}
	var req models.LinkProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}
	req.GoalID = uint(id)
	goal, err := h.goalService.LinkProject(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Project linked successfully", goal)
}

func (h *GoalHandler) UnlinkProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid goal ID", nil)
		return
	}
	goal, err := h.goalService.UnlinkProject(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, "Project unlinked successfully", goal)
}
