package handlers

import (
	"strings"

	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// InterviewWorkflowHandler serves the structured stage-driven interview API.
type InterviewWorkflowHandler struct {
	svc      *services.InterviewWorkflowService
	empRepo  *employeeRepos.EmployeeRepository
}

func NewInterviewWorkflowHandler() *InterviewWorkflowHandler {
	return &InterviewWorkflowHandler{
		svc:     services.NewInterviewWorkflowService(),
		empRepo: employeeRepos.NewEmployeeRepository(),
	}
}

func (h *InterviewWorkflowHandler) resolveActorEmployeeID(c *gin.Context) string {
	user, _ := c.Get("user")
	if userObj, ok := user.(*userModels.User); ok {
		if emp, err := h.empRepo.FindByUserID(userObj.ID); err == nil && emp != nil {
			return strings.TrimSpace(emp.EmployeeID)
		}
	}
	return ""
}

func interviewWorkflowClientErr(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	switch {
	case strings.Contains(s, "not found"),
		strings.Contains(s, "already exists"),
		strings.Contains(s, "only the assigned"),
		strings.Contains(s, "invalid "),
		strings.Contains(s, "cannot "),
		strings.Contains(s, "must be"),
		strings.Contains(s, "requires"),
		strings.Contains(s, "mandatory"),
		strings.Contains(s, "previous "),
		strings.Contains(s, "open attempt"),
		strings.Contains(s, "not accepting"),
		strings.Contains(s, "not in final"),
		strings.Contains(s, "no longer pending"),
		strings.Contains(s, "unknown score"),
		strings.Contains(s, "missing score"),
		strings.Contains(s, "out of range"),
		strings.Contains(s, "attendance was already"),
		strings.Contains(s, "evaluation already"),
		strings.Contains(s, "employee context required"):
		return true
	default:
		return false
	}
}

// --- Definitions ---

func (h *InterviewWorkflowHandler) CreateInterviewWorkflowDefinition(c *gin.Context) {
	var req models.CreateInterviewWorkflowDefinitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}
	def, err := h.svc.CreateDefinition(req)
	if err != nil {
		if interviewWorkflowClientErr(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to create definition", err.Error())
		return
	}
	response.Created(c, "Interview workflow definition created", def)
}

func (h *InterviewWorkflowHandler) UpdateInterviewWorkflowDefinition(c *gin.Context) {
	id := c.Param("id")
	var req models.CreateInterviewWorkflowDefinitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}
	def, err := h.svc.UpdateDefinition(id, req)
	if err != nil {
		if interviewWorkflowClientErr(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to update definition", err.Error())
		return
	}
	response.Success(c, "Definition updated", def)
}

func (h *InterviewWorkflowHandler) ListInterviewWorkflowDefinitions(c *gin.Context) {
	jobID := c.Query("jobOpeningId")
	if jobID == "" {
		jobID = c.Query("jobId")
	}
	if strings.TrimSpace(jobID) == "" {
		response.BadRequest(c, "jobOpeningId is required", nil)
		return
	}
	rows, err := h.svc.ListDefinitionsByJob(strings.TrimSpace(jobID))
	if err != nil {
		response.InternalServerError(c, "Failed to list definitions", err.Error())
		return
	}
	response.Success(c, "Definitions retrieved", rows)
}

func (h *InterviewWorkflowHandler) GetInterviewWorkflowDefinition(c *gin.Context) {
	def, err := h.svc.GetDefinition(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Definition not found")
		return
	}
	response.Success(c, "Definition retrieved", def)
}

// --- Process ---

func (h *InterviewWorkflowHandler) StartInterviewWorkflowProcess(c *gin.Context) {
	var req models.StartInterviewWorkflowProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}
	actor := h.resolveActorEmployeeID(c)
	proc, err := h.svc.StartProcess(actor, req)
	if err != nil {
		if interviewWorkflowClientErr(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to start process", err.Error())
		return
	}
	response.Created(c, "Interview workflow process started", proc)
}

func (h *InterviewWorkflowHandler) GetInterviewWorkflowProcess(c *gin.Context) {
	proc, err := h.svc.GetProcess(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Process not found")
		return
	}
	response.Success(c, "Process retrieved", proc)
}

func (h *InterviewWorkflowHandler) GetInterviewWorkflowProcessByApplication(c *gin.Context) {
	proc, err := h.svc.GetProcessByApplication(c.Param("applicationId"))
	if err != nil {
		response.NotFound(c, "Process not found")
		return
	}
	response.Success(c, "Process retrieved", proc)
}

func (h *InterviewWorkflowHandler) GetInterviewWorkflowTimeline(c *gin.Context) {
	out, err := h.svc.GetTimeline(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Process not found")
		return
	}
	response.Success(c, "Timeline retrieved", out)
}

func (h *InterviewWorkflowHandler) ScheduleInterviewWorkflowStage(c *gin.Context) {
	var req models.ScheduleInterviewWorkflowStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}
	actor := h.resolveActorEmployeeID(c)
	att, err := h.svc.ScheduleStage(c.Param("id"), actor, req)
	if err != nil {
		if interviewWorkflowClientErr(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to schedule stage", err.Error())
		return
	}
	response.Created(c, "Stage scheduled", att)
}

// --- Assignee actions ---

func (h *InterviewWorkflowHandler) RecordInterviewWorkflowAttendance(c *gin.Context) {
	var req models.RecordInterviewWorkflowAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}
	actor := h.resolveActorEmployeeID(c)
	if actor == "" {
		response.BadRequest(c, "Authenticated user must be linked to an employee record", nil)
		return
	}
	if err := h.svc.RecordAttendance(c.Param("id"), actor, req); err != nil {
		if interviewWorkflowClientErr(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to record attendance", err.Error())
		return
	}
	response.Success(c, "Attendance recorded", nil)
}

func (h *InterviewWorkflowHandler) StartInterviewWorkflowStage(c *gin.Context) {
	actor := h.resolveActorEmployeeID(c)
	if actor == "" {
		response.BadRequest(c, "Authenticated user must be linked to an employee record", nil)
		return
	}
	if err := h.svc.StartInProgress(c.Param("id"), actor); err != nil {
		if interviewWorkflowClientErr(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to start stage", err.Error())
		return
	}
	response.Success(c, "Stage marked in progress", nil)
}

func (h *InterviewWorkflowHandler) SubmitInterviewWorkflowEvaluation(c *gin.Context) {
	var req models.SubmitInterviewWorkflowEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}
	actor := h.resolveActorEmployeeID(c)
	if actor == "" {
		response.BadRequest(c, "Authenticated user must be linked to an employee record", nil)
		return
	}
	if err := h.svc.SubmitEvaluation(c.Param("id"), actor, req); err != nil {
		if interviewWorkflowClientErr(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to submit evaluation", err.Error())
		return
	}
	response.Success(c, "Evaluation submitted", nil)
}

func (h *InterviewWorkflowHandler) InterviewWorkflowAbsenceAction(c *gin.Context) {
	var req models.InterviewWorkflowAbsenceActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}
	actor := h.resolveActorEmployeeID(c)
	if actor == "" {
		response.BadRequest(c, "Authenticated user must be linked to an employee record", nil)
		return
	}
	if err := h.svc.AbsenceAction(c.Param("id"), actor, req); err != nil {
		if interviewWorkflowClientErr(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to process absence action", err.Error())
		return
	}
	response.Success(c, "Absence action recorded", nil)
}

func (h *InterviewWorkflowHandler) DecideInterviewWorkflowFinalApproval(c *gin.Context) {
	var req models.InterviewWorkflowFinalApprovalDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}
	actor := h.resolveActorEmployeeID(c)
	if actor == "" {
		response.BadRequest(c, "Authenticated user must be linked to an employee record", nil)
		return
	}
	if err := h.svc.DecideFinalApproval(c.Param("id"), actor, req); err != nil {
		if interviewWorkflowClientErr(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to record decision", err.Error())
		return
	}
	response.Success(c, "Decision recorded", nil)
}
