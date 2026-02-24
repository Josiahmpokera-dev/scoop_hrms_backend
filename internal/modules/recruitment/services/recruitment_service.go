package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/repositories"
	"github.com/google/uuid"
)

type RecruitmentService struct {
	repo *repositories.RecruitmentRepository
}

func NewRecruitmentService() *RecruitmentService {
	return &RecruitmentService{
		repo: repositories.NewRecruitmentRepository(),
	}
}

// --- Requisitions ---

func (s *RecruitmentService) CreateRequisition(reqData models.CreateRequisitionRequest) (*models.JobRequisition, error) {
	targetDate, err := time.Parse("2006-01-02", reqData.TargetStartDate)
	if err != nil {
		return nil, errors.New("invalid target start date format (YYYY-MM-DD)")
	}

	// Generate Requisition Number (Simplified)
	reqNum := fmt.Sprintf("REQ-%s-%d", time.Now().Format("2006"), time.Now().Unix()%1000)

	requisition := &models.JobRequisition{
		ID:              uuid.New().String(),
		RequisitionNo:   reqNum,
		JobTitle:        reqData.JobTitle,
		Department:      reqData.Department,
		Location:        reqData.Location,
		Grade:           reqData.Grade,
		Headcount:       reqData.Headcount,
		EmploymentType:  reqData.EmploymentType,
		HiringManager:   reqData.HiringManager,
		RequestedBy:     "System User", // Should come from context/auth
		TargetStartDate: targetDate,
		EstimatedBudget: reqData.EstimatedBudget,
		Currency:        reqData.Currency,
		Justification:   reqData.Justification,
		Status:          "Pending Approval",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.CreateRequisition(requisition); err != nil {
		return nil, err
	}

	return requisition, nil
}

func (s *RecruitmentService) ListRequisitions(status, department string, page, limit int) (map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	reqs, total, err := s.repo.ListRequisitions(status, department, page, limit)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"data": reqs,
		"meta": map[string]interface{}{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	}, nil
}

func (s *RecruitmentService) GetRequisition(id string) (*models.JobRequisition, error) {
	return s.repo.GetRequisitionByID(id)
}

func (s *RecruitmentService) UpdateRequisition(id string, reqData models.UpdateRequisitionRequest) error {
	req, err := s.repo.GetRequisitionByID(id)
	if err != nil {
		return err
	}

	if reqData.Headcount != nil {
		req.Headcount = *reqData.Headcount
	}
	if reqData.Justification != nil {
		req.Justification = *reqData.Justification
	}
	if reqData.Status != nil {
		req.Status = *reqData.Status
	}
	req.UpdatedAt = time.Now()

	return s.repo.UpdateRequisition(req)
}

func (s *RecruitmentService) ApproveRequisition(id string, approval models.ApprovalRequest) error {
	req, err := s.repo.GetRequisitionByID(id)
	if err != nil {
		return err
	}

	// Logic for approval flow
	newStatus := req.Status
	if approval.Action == "Approve" {
		newStatus = "Approved"
	} else if approval.Action == "Reject" {
		newStatus = "Rejected"
	} else {
		return errors.New("invalid action")
	}

	req.Status = newStatus
	// Add approval step record logic here if needed (appending to ApprovalFlow)
	req.ApprovalFlow = append(req.ApprovalFlow, models.ApprovalStep{
		RequisitionID: req.ID,
		Approver:      "Current User", // Should be from context
		Role:          "Approver",
		Status:        newStatus,
		Comments:      approval.Comments,
		Timestamp:     time.Now(),
	})

	return s.repo.UpdateRequisition(req)
}

// --- Job Openings ---

func (s *RecruitmentService) CreateJobOpening(reqData models.CreateJobOpeningRequest) (*models.JobOpening, error) {
	// Verify Requisition exists and is approved
	req, err := s.repo.GetRequisitionByID(reqData.RequisitionID)
	if err != nil {
		return nil, errors.New("requisition not found")
	}
	if req.Status != "Approved" {
		return nil, errors.New("requisition must be approved before creating a job opening")
	}

	expiryDate, _ := time.Parse(time.RFC3339, reqData.ExpiryDate) // Or specific format

	opening := &models.JobOpening{
		ID:               uuid.New().String(),
		RequisitionID:    reqData.RequisitionID,
		JobTitle:         reqData.JobTitle,
		JobDescription:   reqData.JobDescription,
		Responsibilities: reqData.Responsibilities,
		MustHaveSkills:   reqData.MustHaveSkills,
		NiceToHaveSkills: reqData.NiceToHaveSkills,
		ExperienceMin:    reqData.Experience.Min,
		ExperienceMax:    reqData.Experience.Max,
		Education:        reqData.Education,
		SalaryMin:        reqData.SalaryRange.Min,
		SalaryMax:        reqData.SalaryRange.Max,
		SalaryCurrency:   reqData.SalaryRange.Currency,
		ExpiryDate:       expiryDate,
		Status:           "Draft",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.repo.CreateJobOpening(opening); err != nil {
		return nil, err
	}
	return opening, nil
}

func (s *RecruitmentService) PublishJobOpening(id string, reqData models.PublishJobOpeningRequest) error {
	job, err := s.repo.GetJobOpeningByID(id)
	if err != nil {
		return err
	}

	job.Status = "Published"
	job.Platforms = reqData.Platforms
	if reqData.ExpiryDate != "" {
		expiry, _ := time.Parse(time.RFC3339, reqData.ExpiryDate)
		job.ExpiryDate = expiry
	}
	job.UpdatedAt = time.Now()

	return s.repo.UpdateJobOpening(job)
}

func (s *RecruitmentService) ListJobOpenings(status, department string) ([]models.JobOpening, error) {
	return s.repo.ListJobOpenings(status, department)
}

// --- Applications ---

func (s *RecruitmentService) SubmitApplication(reqData models.ApplyJobRequest) (*models.JobApplication, error) {
	// Check if candidate exists by email
	candidate, _ := s.repo.GetCandidateByEmail(reqData.Email)

	if candidate == nil {
		// Create new candidate
		candidate = &models.Candidate{
			ID:           uuid.New().String(),
			FirstName:    reqData.FirstName,
			LastName:     reqData.LastName,
			Email:        reqData.Email,
			Phone:        reqData.Phone,
			ResumeURL:    reqData.Resume, // Assume URL or handle upload separately
			LinkedInURL:  reqData.LinkedInURL,
			PortfolioURL: reqData.PortfolioURL,
			CoverLetter:  reqData.CoverLetter,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := s.repo.CreateCandidate(candidate); err != nil {
			return nil, err
		}
	}

	// Create Application
	application := &models.JobApplication{
		ID:           uuid.New().String(),
		CandidateID:  candidate.ID,
		JobOpeningID: reqData.JobID,
		Stage:        "Applied",
		AppliedDate:  time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateApplication(application); err != nil {
		return nil, err
	}

	return application, nil
}

func (s *RecruitmentService) ListCandidates(jobID, stage string) ([]models.Candidate, error) {
	return s.repo.ListCandidates(jobID, stage)
}

func (s *RecruitmentService) UpdateApplicationStage(id string, reqData models.UpdateStageRequest) error {
	// ID here is Application ID? Or Candidate ID? Doc says /candidates/{id}/stage
	// Assuming Candidate ID if following REST strictly on /candidates resource, but Application ID is more precise.
	// Let's assume the ID in URL is Candidate ID (from doc), but that's ambiguous if multiple applications.
	// Implementation: Assuming ID is ApplicationID for correctness, or CandidateID and we find the active application.
	// Given the route is /candidates/{id}/stage, it likely means ApplicationID or CandidateID.
	// If it's CandidateID, we need JobID to identify the application.

	// Let's assume the route parameter is actually the Application ID for simplicity in backend,
	// or we look up the application by CandidateID + JobID (if provided in body).

	// If the ID passed is a Candidate ID, we need to find the application for the job.
	// But `UpdateStageRequest` has `jobId`.

	// Let's implement fetching Application by ID directly if the ID is ApplicationID.
	// Or search by CandidateID + JobID.

	// Based on typical patterns, let's treat the URL param as ApplicationID.
	app, err := s.repo.GetApplicationByID(id)
	if err != nil {
		// Try treating as CandidateID if JobID is present
		if reqData.JobID != "" {
			// Find application by candidate and job
			// This method is missing in repo, let's skip and assume ID is ApplicationID
			return errors.New("application not found")
		}
		return err
	}

	app.Stage = reqData.Stage
	if reqData.Notes != "" {
		app.Notes = reqData.Notes
	}
	app.LastStatusDate = time.Now()
	app.UpdatedAt = time.Now()

	return s.repo.UpdateApplication(app)
}

// --- Interviews ---

func (s *RecruitmentService) ScheduleInterview(reqData models.ScheduleInterviewRequest) (*models.Interview, error) {
	// Find application using CandidateID and JobID
	app, err := s.repo.GetApplicationByCandidateAndJob(reqData.CandidateID, reqData.JobID)
	if err != nil {
		return nil, errors.New("active application not found for this candidate and job")
	}

	interview := &models.Interview{
		ID:            uuid.New().String(),
		CandidateID:   reqData.CandidateID,
		JobOpeningID:  reqData.JobID,
		ApplicationID: app.ID,
		InterviewType: reqData.InterviewType,
		Round:         reqData.Round,
		ScheduledTime: reqData.ScheduledTime,
		DurationMin:   reqData.Duration,
		Mode:          reqData.Mode,
		Interviewers:  reqData.Interviewers,
		Status:        "Scheduled",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Parse date
	if t, err := time.Parse("2006-01-02", reqData.ScheduledDate); err == nil {
		interview.ScheduledDate = t
	}

	return interview, s.repo.CreateInterview(interview)
}

func (s *RecruitmentService) SubmitFeedback(interviewID string, reqData models.SubmitFeedbackRequest) error {
	interview, err := s.repo.GetInterviewByID(interviewID)
	if err != nil {
		return err
	}

	feedback := &models.Feedback{
		InterviewID:      interviewID,
		InterviewerEmail: reqData.InterviewerEmail,
		Rating:           reqData.Rating,
		Strengths:        reqData.Strengths,
		Concerns:         reqData.Concerns,
		Recommendation:   reqData.Recommendation,
		Comments:         reqData.Comments,
		CreatedAt:        time.Now(),
	}

	interview.Feedback = feedback
	interview.Status = "Completed"

	return s.repo.UpdateInterview(interview)
}

// --- Offers ---

func (s *RecruitmentService) CreateOffer(reqData models.CreateOfferRequest) (*models.Offer, error) {
	joiningDate, _ := time.Parse("2006-01-02", reqData.JoiningDate)
	expiryDate, _ := time.Parse(time.RFC3339, reqData.ExpiryDate)

	app, err := s.repo.GetApplicationByCandidateAndJob(reqData.CandidateID, reqData.JobID)
	if err != nil {
		return nil, errors.New("active application not found for this candidate and job")
	}

	offer := &models.Offer{
		ID:              uuid.New().String(),
		CandidateID:     reqData.CandidateID,
		JobOpeningID:    reqData.JobID,
		ApplicationID:   app.ID,
		JoiningDate:     joiningDate,
		ProbationPeriod: reqData.ProbationPeriod,
		AnnualCTC:       reqData.Compensation.AnnualCTC,
		Currency:        reqData.Compensation.Currency,
		Components:      reqData.Compensation.Components,
		Benefits:        reqData.Benefits,
		ExpiryDate:      expiryDate,
		Status:          "Draft",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return offer, s.repo.CreateOffer(offer)
}

func (s *RecruitmentService) ApproveOffer(id string, reqData models.OfferApprovalRequest) error {
	offer, err := s.repo.GetOfferByID(id)
	if err != nil {
		return err
	}

	if reqData.Action == "Approve" {
		offer.Status = "Approved"
	} else {
		offer.Status = "Rejected"
	}

	return s.repo.UpdateOffer(offer)
}

func (s *RecruitmentService) SendOffer(id string) error {
	offer, err := s.repo.GetOfferByID(id)
	if err != nil {
		return err
	}

	if offer.Status != "Approved" {
		return errors.New("offer must be approved before sending")
	}

	offer.Status = "Sent"
	now := time.Now()
	offer.SentAt = &now

	// Mock Email Sending
	fmt.Printf("Sending offer to candidate %s\n", offer.CandidateID)

	return s.repo.UpdateOffer(offer)
}

// --- Talent Pool ---

func (s *RecruitmentService) AddToTalentPool(reqData models.AddToTalentPoolRequest) error {
	candidate := &models.TalentPoolCandidate{
		ID:            uuid.New().String(),
		Name:          reqData.Name,
		Email:         reqData.Email,
		Skills:        reqData.Skills,
		PreferredRole: reqData.PreferredRole,
		ResumeURL:     reqData.ResumeURL,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	return s.repo.AddToTalentPool(candidate)
}

func (s *RecruitmentService) SearchTalentPool(skills []string, location string, expMin int) ([]models.TalentPoolCandidate, error) {
	return s.repo.SearchTalentPool(skills, location, expMin)
}
