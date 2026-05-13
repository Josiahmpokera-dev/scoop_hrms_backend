package handlers

import (
	"errors"
	"strconv"
	"strings"

	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/realtime"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type RecruitmentHandler struct {
	service        *services.RecruitmentService
	dashboardSvc   *services.RecruitmentDashboardService
	employeeRepo   *employeeRepos.EmployeeRepository
	hub            *realtime.Hub
}

func NewRecruitmentHandler() *RecruitmentHandler {
	svc := services.NewRecruitmentService()
	svc.EnsureApplicationSubmissionWorker()
	return &RecruitmentHandler{
		service:      svc,
		dashboardSvc: services.NewRecruitmentDashboardService(),
		employeeRepo: employeeRepos.NewEmployeeRepository(),
		hub:          realtime.GetHub(),
	}
}

// --- Requisitions ---

func (h *RecruitmentHandler) CreateRequisition(c *gin.Context) {
	var req models.CreateRequisitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := recruitmentBindErrMessage(err)
		response.BadRequest(c, msg, nil)
		return
	}

	requisition, err := h.service.CreateRequisition(req)
	if err != nil {
		response.InternalServerError(c, "Failed to create requisition", err)
		return
	}

	response.Success(c, "Requisition created successfully", map[string]string{"id": requisition.ID, "message": "Requisition created successfully"})
}

func (h *RecruitmentHandler) ListRequisitions(c *gin.Context) {
	status := c.Query("status")
	department := c.Query("department")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.service.ListRequisitions(status, department, page, limit)
	if err != nil {
		response.InternalServerError(c, "Failed to list requisitions", err)
		return
	}

	// Unpack result for cleaner response if needed, or just pass it
	response.Success(c, "Requisitions retrieved successfully", result)
}

func (h *RecruitmentHandler) GetRequisition(c *gin.Context) {
	id := c.Param("id")
	req, err := h.service.GetRequisition(id)
	if err != nil {
		response.NotFound(c, "Requisition not found")
		return
	}
	response.Success(c, "Requisition retrieved successfully", req)
}

func (h *RecruitmentHandler) UpdateRequisition(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateRequisitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	if err := h.service.UpdateRequisition(id, req); err != nil {
		response.InternalServerError(c, "Failed to update requisition", err)
		return
	}

	response.Success(c, "Requisition updated successfully", nil)
}

func (h *RecruitmentHandler) ApproveRequisition(c *gin.Context) {
	id := c.Param("id")
	var approval models.ApprovalRequest
	if err := c.ShouldBindJSON(&approval); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	if err := h.service.ApproveRequisition(id, approval); err != nil {
		response.InternalServerError(c, "Failed to update requisition status", err)
		return
	}

	response.Success(c, "Requisition status updated successfully", nil)
}

// --- Job Openings ---

func (h *RecruitmentHandler) CreateJobOpening(c *gin.Context) {
	var req models.CreateJobOpeningRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}

	job, err := h.service.CreateJobOpening(req)
	if err != nil {
		if recruitmentClientError(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to create job opening", err.Error())
		return
	}

	payload := interface{}(job)
	if strings.EqualFold(c.Query("include_share_link"), "true") {
		scheme, host := requestSchemeAndHost(c)
		if link, err := h.service.BuildJobOpeningShareLink(job.ID, scheme, host); err == nil {
			payload = gin.H{
				"opening":    job,
				"share_link": link,
			}
		}
	}

	response.Created(c, "Job opening created successfully", payload)
}

func (h *RecruitmentHandler) PublishJobOpening(c *gin.Context) {
	id := c.Param("id")
	var req models.PublishJobOpeningRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}

	if err := h.service.PublishJobOpening(id, req); err != nil {
		if recruitmentClientError(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to activate job opening", err.Error())
		return
	}

	response.Success(c, "Job opening activated successfully", gin.H{"status": models.JobOpeningStatusActive})
}

// ActivateJobOpening is an alias for PublishJobOpening — sets status to Active.
func (h *RecruitmentHandler) ActivateJobOpening(c *gin.Context) {
	h.PublishJobOpening(c)
}

func (h *RecruitmentHandler) CloseJobOpening(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.CloseJobOpening(id); err != nil {
		if recruitmentClientError(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to close job opening", err.Error())
		return
	}
	response.Success(c, "Job opening closed successfully", gin.H{"status": models.JobOpeningStatusClosed})
}

func (h *RecruitmentHandler) ReopenJobOpening(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.ReopenJobOpening(id); err != nil {
		if recruitmentClientError(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to reopen job opening", err.Error())
		return
	}
	response.Success(c, "Job opening moved back to Draft", gin.H{"status": models.JobOpeningStatusDraft})
}

func (h *RecruitmentHandler) UpdateJobOpening(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateJobOpeningRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}
	if err := h.service.UpdateJobOpening(id, req); err != nil {
		if recruitmentClientError(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to update job opening", err.Error())
		return
	}
	response.Success(c, "Job opening updated successfully", nil)
}

func (h *RecruitmentHandler) GetJobOpening(c *gin.Context) {
	id := c.Param("id")
	job, err := h.service.GetJobOpening(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "Job opening not found")
			return
		}
		response.InternalServerError(c, "Failed to retrieve job opening", err.Error())
		return
	}
	response.Success(c, "Job opening retrieved successfully", job)
}

func (h *RecruitmentHandler) ShareJobOpeningLink(c *gin.Context) {
	id := c.Param("id")
	scheme, host := requestSchemeAndHost(c)
	link, err := h.service.BuildJobOpeningShareLink(id, scheme, host)
	if err != nil {
		if recruitmentClientError(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "Job opening not found")
			return
		}
		response.InternalServerError(c, "Failed to build share link", err.Error())
		return
	}
	response.Success(c, "Share link generated", link)
}

// GetRecruitmentDashboardSummary returns consolidated KPIs, pipeline, quick stats, and lists for the recruitment dashboard.
func (h *RecruitmentHandler) GetRecruitmentDashboardSummary(c *gin.Context) {
	deptID := c.Query("department_id")
	if deptID == "" {
		deptID = c.Query("departmentId")
	}
	asOf := c.Query("as_of")
	if asOf == "" {
		asOf = c.Query("asOf")
	}
	out, err := h.dashboardSvc.GetSummary(deptID, asOf)
	if err != nil {
		response.InternalServerError(c, "Failed to build recruitment dashboard summary", err.Error())
		return
	}
	response.Success(c, "Recruitment dashboard summary retrieved successfully", out)
}

func (h *RecruitmentHandler) ListJobOpenings(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		status = c.Query("Status")
	}
	department := c.Query("department")
	requisitionID := c.Query("requisition_id")
	if requisitionID == "" {
		requisitionID = c.Query("requisitionId")
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	out, err := h.service.ListJobOpeningsPaged(status, department, requisitionID, page, limit)
	if err != nil {
		response.InternalServerError(c, "Failed to list job openings", err.Error())
		return
	}
	response.SuccessWithMeta(c, "Job openings retrieved successfully", out["data"], recruitmentMeta(out["meta"]))
}

// ListJobOpeningsGrouped returns openings grouped by status for dashboards (Draft/Active/Closed).
func (h *RecruitmentHandler) ListJobOpeningsGrouped(c *gin.Context) {
	department := c.Query("department")
	requisitionID := c.Query("requisition_id")
	if requisitionID == "" {
		requisitionID = c.Query("requisitionId")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	out, err := h.service.ListJobOpeningsGrouped(department, requisitionID, page, limit)
	if err != nil {
		response.InternalServerError(c, "Failed to list job openings", err.Error())
		return
	}
	response.SuccessWithMeta(c, "Job openings retrieved successfully", out["data"], recruitmentMeta(out["meta"]))
}

// SearchJobOpenings returns filtered openings with the same response format as GET /openings.
func (h *RecruitmentHandler) SearchJobOpenings(c *gin.Context) {
	keyword := c.Query("q")
	status := c.Query("status")
	if status == "" {
		status = c.Query("Status")
	}
	department := c.Query("department")
	requisitionID := c.Query("requisition_id")
	if requisitionID == "" {
		requisitionID = c.Query("requisitionId")
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	out, err := h.service.SearchJobOpeningsPaged(keyword, status, department, requisitionID, page, limit)
	if err != nil {
		response.InternalServerError(c, "Failed to search job openings", err.Error())
		return
	}
	response.SuccessWithMeta(c, "Job openings retrieved successfully", out["data"], recruitmentMeta(out["meta"]))
}

func recruitmentMeta(raw interface{}) *response.Meta {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	return &response.Meta{
		Page:       anyInt(m["page"]),
		PerPage:    anyInt(m["limit"]),
		Total:      anyInt64(m["total"]),
		TotalPages: anyInt(m["total_pages"]),
	}
}

func anyInt(v interface{}) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case int32:
		return int(n)
	case float64:
		return int(n)
	default:
		return 0
	}
}

func anyInt64(v interface{}) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case int32:
		return int64(n)
	case float64:
		return int64(n)
	default:
		return 0
	}
}

func (h *RecruitmentHandler) ListPublicJobOpenings(c *gin.Context) {
	department := c.Query("department")

	jobs, err := h.service.ListPublicActiveOpenings(department)
	if err != nil {
		response.InternalServerError(c, "Failed to list job openings", err.Error())
		return
	}

	response.Success(c, "Job openings retrieved successfully", jobs)
}

func (h *RecruitmentHandler) GetPublicJobOpeningByToken(c *gin.Context) {
	token := strings.TrimSpace(c.Param("token"))
	job, err := h.service.GetPublicOpeningByApplyToken(token)
	if err != nil || job == nil {
		response.NotFound(c, "Job opening not found or not accepting applications")
		return
	}
	response.Success(c, "Job opening retrieved successfully", job)
}

// CandidatesWS opens a websocket stream for recruitment candidate/interview updates.
func (h *RecruitmentHandler) CandidatesWS(c *gin.Context) {
	h.hub.HandleWS(c)
}

// ListJobOpeningApplications returns applicants for one job opening (includes candidate profile).
func (h *RecruitmentHandler) ListJobOpeningApplications(c *gin.Context) {
	jobID := c.Param("id")
	stage := c.Query("stage")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	out, err := h.service.ListApplicationsForJobOpening(jobID, stage, page, limit)
	if err != nil {
		response.InternalServerError(c, "Failed to list applications", err.Error())
		return
	}
	response.Success(c, "Applications retrieved successfully", out)
}

// ListJobOpeningInterviews lists interviews scheduled for this opening.
func (h *RecruitmentHandler) ListJobOpeningInterviews(c *gin.Context) {
	jobID := c.Param("id")
	rows, err := h.service.ListInterviewsForJobOpening(jobID)
	if err != nil {
		response.InternalServerError(c, "Failed to list interviews", err.Error())
		return
	}
	response.Success(c, "Interviews retrieved successfully", rows)
}

// ListInterviews lists interviews globally (all stages by default) with optional filters.
func (h *RecruitmentHandler) ListInterviews(c *gin.Context) {
	roundStr := c.Query("stage")
	if roundStr == "" {
		roundStr = c.Query("round")
	}
	round := 0
	if roundStr != "" {
		r, err := strconv.Atoi(roundStr)
		if err != nil || r < 1 || r > 3 {
			response.BadRequest(c, "stage (or round) must be 1, 2, or 3", nil)
			return
		}
		round = r
	}
	status := c.Query("status")
	jobID := c.Query("jobId")
	candidateID := c.Query("candidateId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	out, err := h.service.ListInterviewsPaged(round, status, jobID, candidateID, page, limit)
	if err != nil {
		response.InternalServerError(c, "Failed to list interviews", err.Error())
		return
	}
	response.SuccessWithMeta(c, "Interviews retrieved successfully", out["data"], recruitmentMeta(out["meta"]))
}

func (h *RecruitmentHandler) GetInterviewDetails(c *gin.Context) {
	id := c.Param("id")
	row, err := h.service.GetInterviewDetails(id)
	if err != nil || row == nil {
		response.NotFound(c, "Interview not found")
		return
	}
	response.Success(c, "Interview retrieved successfully", row)
}

// ListJobOpeningOffers lists offers for this opening.
func (h *RecruitmentHandler) ListJobOpeningOffers(c *gin.Context) {
	jobID := c.Param("id")
	rows, err := h.service.ListOffersForJobOpening(jobID)
	if err != nil {
		response.InternalServerError(c, "Failed to list offers", err.Error())
		return
	}
	response.Success(c, "Offers retrieved successfully", rows)
}

// --- Applications ---

func (h *RecruitmentHandler) SubmitApplication(c *gin.Context) {
	var req models.ApplyJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}

	app, err := h.service.SubmitApplication(req)
	if err != nil {
		if recruitmentClientError(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to submit application", err.Error())
		return
	}

	response.Created(c, "Application submitted successfully", app)
}

func (h *RecruitmentHandler) ListCandidates(c *gin.Context) {
	jobID := c.Query("jobId")
	stage := c.Query("stage")
	q := c.Query("q")
	if q == "" {
		q = c.Query("search")
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	out, err := h.service.ListCandidatesPaged(jobID, stage, q, page, limit)
	if err != nil {
		response.InternalServerError(c, "Failed to list candidates", err)
		return
	}

	response.SuccessWithMeta(c, "Candidates retrieved successfully", out["data"], recruitmentMeta(out["meta"]))
}

func (h *RecruitmentHandler) GetCandidateDetails(c *gin.Context) {
	id := c.Param("id")
	row, err := h.service.GetCandidateDetails(id)
	if err != nil || row == nil {
		response.NotFound(c, "Candidate not found")
		return
	}
	response.Success(c, "Candidate retrieved successfully", row)
}

// ListCandidatesByInterviewStage lists candidates selected for interview stage (round) 1/2/3.
func (h *RecruitmentHandler) ListCandidatesByInterviewStage(c *gin.Context) {
	stageStr := c.Query("stage")
	if stageStr == "" {
		stageStr = c.Query("round")
	}
	stage, err := strconv.Atoi(stageStr)
	if err != nil || stage < 1 || stage > 3 {
		response.BadRequest(c, "stage (or round) must be 1, 2, or 3", nil)
		return
	}
	jobID := c.Query("jobId")
	q := c.Query("q")
	if q == "" {
		q = c.Query("search")
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	out, err := h.service.ListCandidatesByInterviewStage(stage, jobID, q, page, limit)
	if err != nil {
		if recruitmentClientError(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to list candidates", err.Error())
		return
	}
	response.SuccessWithMeta(c, "Candidates retrieved successfully", out["data"], recruitmentMeta(out["meta"]))
}

func (h *RecruitmentHandler) UpdateCandidateStage(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	if err := h.service.UpdateApplicationStage(id, req); err != nil {
		response.InternalServerError(c, "Failed to update candidate stage", err)
		return
	}

	response.Success(c, "Candidate stage updated successfully", nil)
}

// UpdateApplicationStage updates pipeline stage by application ID (same payload as PATCH /candidates/:id/stage).
func (h *RecruitmentHandler) UpdateApplicationStage(c *gin.Context) {
	h.UpdateCandidateStage(c)
}

func (h *RecruitmentHandler) MoveApplicationToTalentPool(c *gin.Context) {
	id := c.Param("id")
	var req models.MoveApplicationToTalentPoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}
	if err := h.service.MoveApplicationToTalentPool(id, req); err != nil {
		if recruitmentClientError(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to update application", err.Error())
		return
	}
	response.Success(c, "Application moved to talent pool", gin.H{"stage": models.StageTalentPool})
}

// ActionApplication performs accept/reject/interview-stage-one actions on one application.
func (h *RecruitmentHandler) ActionApplication(c *gin.Context) {
	id := c.Param("id")
	var req models.ApplicationActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}
	if err := h.service.ActionApplication(id, req); err != nil {
		if recruitmentClientError(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to process application action", err.Error())
		return
	}
	response.Success(c, "Application action processed successfully", nil)
}

// --- Interviews ---

func (h *RecruitmentHandler) ScheduleInterview(c *gin.Context) {
	var req models.ScheduleInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	interview, err := h.service.ScheduleInterview(req)
	if err != nil {
		if recruitmentClientError(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to schedule interview", err)
		return
	}

	response.Success(c, "Interview scheduled successfully", interview)
}

// ScheduleInterviewManual adds a candidate manually and schedules their first interview (walk-in / referral).
func (h *RecruitmentHandler) ScheduleInterviewManual(c *gin.Context) {
	var req models.ManualScheduleInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}
	inv, err := h.service.ScheduleInterviewManual(req)
	if err != nil {
		if recruitmentClientError(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to schedule interview", err.Error())
		return
	}
	response.Created(c, "Interview scheduled successfully", inv)
}

// ScheduleInterviewBatch schedules the same slot for multiple applications (pipeline).
func (h *RecruitmentHandler) ScheduleInterviewBatch(c *gin.Context) {
	var req models.BatchScheduleInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}
	list, err := h.service.ScheduleInterviewsBatch(req)
	if err != nil {
		if recruitmentClientError(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to schedule interviews", err.Error())
		return
	}
	response.Created(c, "Interviews scheduled successfully", gin.H{"interviews": list, "count": len(list)})
}

func (h *RecruitmentHandler) SubmitFeedback(c *gin.Context) {
	id := c.Param("id")
	var req models.SubmitFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	if err := h.service.SubmitFeedback(id, req); err != nil {
		if recruitmentClientError(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to submit feedback", err)
		return
	}

	response.Success(c, "Feedback submitted successfully", nil)
}

// SubmitHrInterviewDecision applies HR remarks and advances the process (next round, re-interview, or reject).
func (h *RecruitmentHandler) SubmitHrInterviewDecision(c *gin.Context) {
	id := c.Param("id")
	var req models.HRInterviewDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	deciderID := strings.TrimSpace(req.DecidedByEmployeeID)
	if deciderID == "" {
		user, _ := c.Get("user")
		if userObj, ok := user.(*userModels.User); ok {
			if emp, err := h.employeeRepo.FindByUserID(userObj.ID); err == nil && emp != nil {
				deciderID = emp.EmployeeID
			}
		}
	}
	if deciderID == "" {
		response.BadRequest(c, "Could not resolve HR approver; pass decidedByEmployeeId in the body", nil)
		return
	}

	if err := h.service.SubmitHrInterviewDecision(id, deciderID, req); err != nil {
		if recruitmentClientError(err) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to record HR decision", err)
		return
	}

	response.Success(c, "HR decision recorded successfully", nil)
}

// --- Offers ---

func (h *RecruitmentHandler) CreateOffer(c *gin.Context) {
	var req models.CreateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	offer, err := h.service.CreateOffer(req)
	if err != nil {
		response.InternalServerError(c, "Failed to create offer", err)
		return
	}

	response.Success(c, "Offer created successfully", offer)
}

func (h *RecruitmentHandler) ApproveOffer(c *gin.Context) {
	id := c.Param("id")
	var req models.OfferApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	if err := h.service.ApproveOffer(id, req); err != nil {
		response.InternalServerError(c, "Failed to approve offer", err)
		return
	}

	response.Success(c, "Offer status updated successfully", nil)
}

func (h *RecruitmentHandler) SendOffer(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.SendOffer(id); err != nil {
		response.InternalServerError(c, "Failed to send offer", err)
		return
	}

	response.Success(c, "Offer sent successfully", map[string]string{"message": "Offer letter sent to candidate email."})
}

// --- Talent Pool ---

func (h *RecruitmentHandler) AddToTalentPool(c *gin.Context) {
	var req models.AddToTalentPoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	if err := h.service.AddToTalentPool(req); err != nil {
		response.InternalServerError(c, "Failed to add to talent pool", err)
		return
	}

	response.Success(c, "Added to talent pool successfully", nil)
}

func (h *RecruitmentHandler) SearchTalentPool(c *gin.Context) {
	skillsStr := c.Query("skills")
	var skills []string
	if skillsStr != "" {
		skills = strings.Split(skillsStr, ",")
	}

	location := c.Query("location")
	expMin, _ := strconv.Atoi(c.DefaultQuery("experienceMin", "0"))

	candidates, err := h.service.SearchTalentPool(skills, location, expMin)
	if err != nil {
		response.InternalServerError(c, "Failed to search talent pool", err)
		return
	}

	response.Success(c, "Talent pool search results", candidates)
}

// ListTalentPool lists stored talent pool profiles for outreach (paginated).
func (h *RecruitmentHandler) ListTalentPool(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	out, err := h.service.ListTalentPoolPaged(page, limit)
	if err != nil {
		response.InternalServerError(c, "Failed to list talent pool", err.Error())
		return
	}
	response.Success(c, "Talent pool retrieved successfully", out)
}

// ContactTalentPoolCandidate sends an email to a talent pool entry (SMTP required).
func (h *RecruitmentHandler) ContactTalentPoolCandidate(c *gin.Context) {
	id := c.Param("id")
	var req models.ContactTalentPoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}
	if err := h.service.ContactTalentPoolCandidate(id, req); err != nil {
		response.InternalServerError(c, "Failed to send message", err.Error())
		return
	}
	response.Success(c, "Message sent successfully", nil)
}

func recruitmentClientError(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	switch {
	case strings.Contains(s, "requisition must be approved"):
		return true
	case strings.Contains(s, "only Draft openings"):
		return true
	case strings.Contains(s, "only Draft"):
		return true
	case strings.Contains(s, "only active openings"):
		return true
	case strings.Contains(s, "only closed openings"):
		return true
	case strings.Contains(s, "must be activated"):
		return true
	case strings.Contains(s, "expiryDate is required"):
		return true
	case strings.Contains(s, "invalid date"):
		return true
	case strings.Contains(s, "job opening not found"):
		return true
	case strings.Contains(s, "this job is not accepting"):
		return true
	case strings.Contains(s, "jobId or applyToken"):
		return true
	case strings.Contains(s, "opening has no apply token"):
		return true
	case strings.Contains(s, "already has an application"):
		return true
	case strings.Contains(s, "does not belong to job opening"):
		return true
	case strings.Contains(s, "already in talent pool"):
		return true
	case strings.Contains(s, "another interview is still pending"):
		return true
	case strings.Contains(s, "employee not found:"):
		return true
	case strings.Contains(s, "provide interviewerEmployeeIds"):
		return true
	case strings.Contains(s, "remarks are required"):
		return true
	case strings.Contains(s, "HR decision is only allowed"):
		return true
	case strings.Contains(s, "invalid decision"):
		return true
	case strings.Contains(s, "cannot schedule re-interview"):
		return true
	case strings.Contains(s, "cannot schedule round"):
		return true
	case strings.Contains(s, "previous round must be HR-approved"):
		return true
	case strings.Contains(s, "cannot schedule before"):
		return true
	case strings.Contains(s, "HR must request a re-interview"):
		return true
	case strings.Contains(s, "use isReinterview=true"):
		return true
	case strings.Contains(s, "feedback was already submitted"):
		return true
	case strings.Contains(s, "feedback can only be submitted"):
		return true
	case strings.Contains(s, "invalid attempt state"):
		return true
	case strings.Contains(s, "round must be between"):
		return true
	case strings.Contains(s, "invalid action: use accept"):
		return true
	case strings.Contains(s, "scheduledDate and scheduledTime are required"):
		return true
	case strings.Contains(s, "stage must be 1, 2, or 3"):
		return true
	case strings.Contains(s, "round) must be 1, 2, or 3"):
		return true
	default:
		return false
	}
}

func requestSchemeAndHost(c *gin.Context) (scheme, host string) {
	scheme = "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	if p := strings.ToLower(strings.TrimSpace(c.GetHeader("X-Forwarded-Proto"))); p == "https" || p == "http" {
		scheme = p
	}
	host = c.Request.Host
	if xh := strings.TrimSpace(c.GetHeader("X-Forwarded-Host")); xh != "" {
		host = xh
	}
	return scheme, host
}

func recruitmentBindErrMessage(err error) string {
	if err == nil {
		return ""
	}
	var ves validator.ValidationErrors
	if errors.As(err, &ves) {
		var b strings.Builder
		b.WriteString("Invalid request body: ")
		for i, fe := range ves {
			if i > 0 {
				b.WriteString("; ")
			}
			switch fe.Tag() {
			case "required":
				b.WriteString(fe.Field() + " is required")
			case "gt":
				b.WriteString(fe.Field() + " must be greater than " + fe.Param())
			default:
				b.WriteString(fe.Field() + ": " + fe.Tag())
			}
		}
		return b.String()
	}
	return "Invalid request body: " + err.Error()
}
