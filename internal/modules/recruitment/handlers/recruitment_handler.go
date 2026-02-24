package handlers

import (
	"strconv"
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type RecruitmentHandler struct {
	service *services.RecruitmentService
}

func NewRecruitmentHandler() *RecruitmentHandler {
	return &RecruitmentHandler{
		service: services.NewRecruitmentService(),
	}
}

// --- Requisitions ---

func (h *RecruitmentHandler) CreateRequisition(c *gin.Context) {
	var req models.CreateRequisitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
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
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	job, err := h.service.CreateJobOpening(req)
	if err != nil {
		response.InternalServerError(c, "Failed to create job opening", err)
		return
	}

	response.Success(c, "Job opening created successfully", job)
}

func (h *RecruitmentHandler) PublishJobOpening(c *gin.Context) {
	id := c.Param("id")
	var req models.PublishJobOpeningRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	if err := h.service.PublishJobOpening(id, req); err != nil {
		response.InternalServerError(c, "Failed to publish job opening", err)
		return
	}

	response.Success(c, "Job opening published successfully", nil)
}

func (h *RecruitmentHandler) ListJobOpenings(c *gin.Context) {
	status := c.Query("status")
	department := c.Query("department")

	jobs, err := h.service.ListJobOpenings(status, department)
	if err != nil {
		response.InternalServerError(c, "Failed to list job openings", err)
		return
	}

	response.Success(c, "Job openings retrieved successfully", jobs)
}

// --- Applications ---

func (h *RecruitmentHandler) SubmitApplication(c *gin.Context) {
	var req models.ApplyJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	app, err := h.service.SubmitApplication(req)
	if err != nil {
		response.InternalServerError(c, "Failed to submit application", err)
		return
	}

	response.Success(c, "Application submitted successfully", app)
}

func (h *RecruitmentHandler) ListCandidates(c *gin.Context) {
	jobID := c.Query("jobId")
	stage := c.Query("stage")

	candidates, err := h.service.ListCandidates(jobID, stage)
	if err != nil {
		response.InternalServerError(c, "Failed to list candidates", err)
		return
	}

	response.Success(c, "Candidates retrieved successfully", candidates)
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

// --- Interviews ---

func (h *RecruitmentHandler) ScheduleInterview(c *gin.Context) {
	var req models.ScheduleInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	interview, err := h.service.ScheduleInterview(req)
	if err != nil {
		response.InternalServerError(c, "Failed to schedule interview", err)
		return
	}

	response.Success(c, "Interview scheduled successfully", interview)
}

func (h *RecruitmentHandler) SubmitFeedback(c *gin.Context) {
	id := c.Param("id")
	var req models.SubmitFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	if err := h.service.SubmitFeedback(id, req); err != nil {
		response.InternalServerError(c, "Failed to submit feedback", err)
		return
	}

	response.Success(c, "Feedback submitted successfully", nil)
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
