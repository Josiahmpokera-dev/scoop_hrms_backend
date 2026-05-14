package services

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	empModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/realtime"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/email"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/storage"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RecruitmentService struct {
	repo         *repositories.RecruitmentRepository
	employeeRepo *employeeRepos.EmployeeRepository
	storage      *storage.StorageService
	email        *email.EmailService
	hub          *realtime.Hub
}

var applicationQueueWorkerOnce sync.Once

func NewRecruitmentService() *RecruitmentService {
	return &RecruitmentService{
		repo:         repositories.NewRecruitmentRepository(),
		employeeRepo: employeeRepos.NewEmployeeRepository(),
		storage:      storage.NewStorageService(),
		email:        email.NewEmailService(),
		hub:          realtime.GetHub(),
	}
}

func (s *RecruitmentService) EnsureApplicationSubmissionWorker() {
	applicationQueueWorkerOnce.Do(func() {
		go s.runApplicationSubmissionQueueWorker()
	})
}

// --- Requisitions ---

func (s *RecruitmentService) CreateRequisition(reqData models.CreateRequisitionRequest) (*models.JobRequisition, error) {
	targetDate, err := time.Parse("2006-01-02", reqData.TargetStartDate)
	if err != nil {
		return nil, errors.New("invalid target start date format (YYYY-MM-DD)")
	}

	reqNum := strings.TrimSpace(reqData.RequisitionNo)
	if reqNum == "" {
		reqNum = fmt.Sprintf("REQ-%s-%d", time.Now().Format("2006"), time.Now().Unix()%1000)
	}

	employmentType := strings.TrimSpace(reqData.EmploymentType)
	if employmentType == "" {
		employmentType = "full_time"
	}

	currency := strings.TrimSpace(reqData.Currency)
	if currency == "" {
		currency = "TZS"
	}

	status := strings.TrimSpace(reqData.Status)
	if status == "" {
		status = "Pending Approval"
	}

	requisition := &models.JobRequisition{
		ID:              uuid.New().String(),
		RequisitionNo:   reqNum,
		JobTitle:        reqData.JobTitle,
		Department:      reqData.Department,
		Location:        reqData.Location,
		Grade:           reqData.Grade,
		Headcount:       reqData.Headcount,
		EmploymentType:  employmentType,
		HiringManager: func() string {
			hm := strings.TrimSpace(reqData.HiringManager)
			if hm == "" {
				return "HR"
			}
			return hm
		}(),
		RequestedBy:     "System User", // Should come from context/auth
		TargetStartDate: targetDate,
		EstimatedBudget: reqData.EstimatedBudget,
		Currency:        currency,
		Justification:   reqData.Justification,
		Priority:        strings.TrimSpace(reqData.Priority),
		Status:          status,
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

func parseFlexibleDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, nil
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05Z07:00", "2006-01-02"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s (use YYYY-MM-DD or RFC3339)", s)
}

func (s *RecruitmentService) generateApplyToken() (string, error) {
	for range 10 {
		b := make([]byte, 9)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		tok := base64.RawURLEncoding.EncodeToString(b)
		taken, err := s.repo.IsApplyTokenTaken(tok)
		if err != nil {
			return "", err
		}
		if !taken {
			return tok, nil
		}
	}
	return "", errors.New("could not generate unique apply token")
}

func (s *RecruitmentService) jobIsPubliclyVisible(j *models.JobOpening) bool {
	if j == nil {
		return false
	}
	if j.Status != models.JobOpeningStatusActive && j.Status != models.LegacyStatusPublished {
		return false
	}
	if !j.ExpiryDate.IsZero() && j.ExpiryDate.Before(time.Now()) {
		return false
	}
	return true
}

func (s *RecruitmentService) CreateJobOpening(reqData models.CreateJobOpeningRequest) (*models.JobOpening, error) {
	req, err := s.repo.GetRequisitionByID(reqData.RequisitionID)
	if err != nil {
		return nil, errors.New("requisition not found")
	}
	if req.Status != "Approved" {
		return nil, errors.New("requisition must be approved before creating a job opening")
	}

	expiryDate := time.Time{}
	if strings.TrimSpace(reqData.ExpiryDate) != "" {
		expiryDate, err = parseFlexibleDate(reqData.ExpiryDate)
		if err != nil {
			return nil, err
		}
	}

	token, err := s.generateApplyToken()
	if err != nil {
		return nil, err
	}

	opening := &models.JobOpening{
		ID:               uuid.New().String(),
		RequisitionID:    reqData.RequisitionID,
		ApplyToken:       token,
		JobTitle:         reqData.JobTitle,
		JobDescription:   reqData.JobDescription,
		Responsibilities: reqData.Responsibilities,
		MustHaveSkills:   reqData.MustHaveSkills,
		NiceToHaveSkills: reqData.NiceToHaveSkills,
		ExperienceMin:    reqData.Experience.Min,
		ExperienceMax:    reqData.Experience.Max,
		Education:        reqData.Education,
		SalaryMin:        0,
		SalaryMax:        0,
		SalaryCurrency:   "TZS",
		ExpiryDate:       expiryDate,
		Status:           models.JobOpeningStatusDraft,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.repo.CreateJobOpening(opening); err != nil {
		return nil, err
	}
	return opening, nil
}

// ActivateJobOpening marks a Draft opening as Active (published for candidates).
func (s *RecruitmentService) ActivateJobOpening(id string, reqData models.PublishJobOpeningRequest) error {
	job, err := s.repo.GetJobOpeningByID(id)
	if err != nil {
		return err
	}
	if job.Status != models.JobOpeningStatusDraft {
		return fmt.Errorf("only %s openings can be activated", models.JobOpeningStatusDraft)
	}
	if strings.TrimSpace(job.ApplyToken) == "" {
		tok, genErr := s.generateApplyToken()
		if genErr != nil {
			return genErr
		}
		job.ApplyToken = tok
	}
	if strings.TrimSpace(reqData.ExpiryDate) != "" {
		exp, perr := parseFlexibleDate(reqData.ExpiryDate)
		if perr != nil {
			return perr
		}
		job.ExpiryDate = exp
	}
	if job.ExpiryDate.IsZero() {
		return errors.New("expiryDate is required to activate — set it on activate or when creating the opening")
	}
	if reqData.Platforms != nil {
		job.Platforms = reqData.Platforms
	}
	job.Status = models.JobOpeningStatusActive
	job.UpdatedAt = time.Now()

	return s.repo.UpdateJobOpening(job)
}

// PublishJobOpening kept for backwards compatibility — same as ActivateJobOpening.
func (s *RecruitmentService) PublishJobOpening(id string, reqData models.PublishJobOpeningRequest) error {
	return s.ActivateJobOpening(id, reqData)
}

func (s *RecruitmentService) CloseJobOpening(id string) error {
	job, err := s.repo.GetJobOpeningByID(id)
	if err != nil {
		return err
	}
	if job.Status != models.JobOpeningStatusActive && job.Status != models.LegacyStatusPublished {
		return errors.New("only active openings can be closed")
	}
	job.Status = models.JobOpeningStatusClosed
	job.UpdatedAt = time.Now()
	return s.repo.UpdateJobOpening(job)
}

func (s *RecruitmentService) ReopenJobOpening(id string) error {
	job, err := s.repo.GetJobOpeningByID(id)
	if err != nil {
		return err
	}
	if job.Status != models.JobOpeningStatusClosed {
		return errors.New("only closed openings can be reopened to draft")
	}
	job.Status = models.JobOpeningStatusDraft
	job.UpdatedAt = time.Now()
	return s.repo.UpdateJobOpening(job)
}

func (s *RecruitmentService) UpdateJobOpening(id string, patch models.UpdateJobOpeningRequest) error {
	job, err := s.repo.GetJobOpeningByID(id)
	if err != nil {
		return err
	}
	if job.Status != models.JobOpeningStatusDraft {
		return errors.New("only draft openings can be updated")
	}
	if patch.JobTitle != nil {
		job.JobTitle = *patch.JobTitle
	}
	if patch.JobDescription != nil {
		job.JobDescription = *patch.JobDescription
	}
	if patch.Responsibilities != nil {
		job.Responsibilities = patch.Responsibilities
	}
	if patch.MustHaveSkills != nil {
		job.MustHaveSkills = patch.MustHaveSkills
	}
	if patch.NiceToHaveSkills != nil {
		job.NiceToHaveSkills = patch.NiceToHaveSkills
	}
	if patch.Education != nil {
		job.Education = patch.Education
	}
	if patch.ExperienceMin != nil {
		job.ExperienceMin = *patch.ExperienceMin
	}
	if patch.ExperienceMax != nil {
		job.ExperienceMax = *patch.ExperienceMax
	}
	if patch.SalaryMin != nil {
		job.SalaryMin = *patch.SalaryMin
	}
	if patch.SalaryMax != nil {
		job.SalaryMax = *patch.SalaryMax
	}
	if patch.SalaryCurrency != nil {
		job.SalaryCurrency = strings.TrimSpace(*patch.SalaryCurrency)
	}
	if patch.ExpiryDate != nil && strings.TrimSpace(*patch.ExpiryDate) != "" {
		exp, perr := parseFlexibleDate(*patch.ExpiryDate)
		if perr != nil {
			return perr
		}
		job.ExpiryDate = exp
	}
	if patch.Platforms != nil {
		job.Platforms = patch.Platforms
	}
	job.UpdatedAt = time.Now()
	return s.repo.UpdateJobOpening(job)
}

func (s *RecruitmentService) GetJobOpening(id string) (*models.JobOpening, error) {
	return s.repo.GetJobOpeningByID(id)
}

func (s *RecruitmentService) GetPublicOpeningByApplyToken(token string) (*models.JobOpening, error) {
	job, err := s.repo.GetJobOpeningByApplyToken(token)
	if err != nil {
		return nil, err
	}
	if !s.jobIsPubliclyVisible(job) {
		return nil, errors.New("job opening not found or not accepting applications")
	}
	return job, nil
}

// ListPublicActiveOpenings returns Active (and legacy Published) openings accepting applications.
func (s *RecruitmentService) ListPublicActiveOpenings(department string) ([]models.JobOpening, error) {
	return s.repo.ListPublicActiveOpenings(department, time.Now())
}

func (s *RecruitmentService) ListJobOpenings(status, department string) ([]models.JobOpening, error) {
	jobs, err := s.repo.ListJobOpenings(status, department)
	if err != nil {
		return nil, err
	}
	return s.attachApplicantCounts(jobs)
}

func (s *RecruitmentService) ListJobOpeningsPaged(status, department, requisitionID string, page, limit int) (map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	jobs, total, err := s.repo.ListJobOpeningsPaginated(status, department, requisitionID, page, limit)
	if err != nil {
		return nil, err
	}
	jobs, err = s.attachApplicantCounts(jobs)
	if err != nil {
		return nil, err
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 1 {
		totalPages = 1
	}
	return map[string]interface{}{
		"data": jobs,
		"meta": map[string]interface{}{
			"total":          total,
			"page":           page,
			"limit":          limit,
			"total_pages":    totalPages,
			"requisition_id": requisitionID,
			"status":         status,
			"department":     department,
		},
	}, nil
}

// ListJobOpeningsGrouped kept for backward route compatibility.
// Despite route name (/openings/grouped), response is flat and consistent with ListJobOpeningsPaged.
func (s *RecruitmentService) ListJobOpeningsGrouped(department, requisitionID string, page, limit int) (map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	rows, total, err := s.repo.ListJobOpeningsPaginated("", department, requisitionID, page, limit)
	if err != nil {
		return nil, err
	}
	rows, err = s.attachApplicantCounts(rows)
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 1 {
		totalPages = 1
	}

	return map[string]interface{}{
		"data": rows,
		"meta": map[string]interface{}{
			"total":          total,
			"page":           page,
			"limit":          limit,
			"total_pages":    totalPages,
			"department":     department,
			"requisition_id": requisitionID,
			"status":         "",
		},
	}, nil
}

// SearchJobOpeningsPaged returns job openings filtered by keyword while preserving
// the same response structure as ListJobOpeningsPaged.
func (s *RecruitmentService) SearchJobOpeningsPaged(keyword, status, department, requisitionID string, page, limit int) (map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	jobs, total, err := s.repo.SearchJobOpeningsPaginated(keyword, status, department, requisitionID, page, limit)
	if err != nil {
		return nil, err
	}
	jobs, err = s.attachApplicantCounts(jobs)
	if err != nil {
		return nil, err
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 1 {
		totalPages = 1
	}
	return map[string]interface{}{
		"data": jobs,
		"meta": map[string]interface{}{
			"total":          total,
			"page":           page,
			"limit":          limit,
			"total_pages":    totalPages,
			"requisition_id": requisitionID,
			"status":         status,
			"department":     department,
			"query":          keyword,
		},
	}, nil
}

// BuildJobOpeningShareLink builds URLs recruiters can paste into email or the careers site.
func (s *RecruitmentService) BuildJobOpeningShareLink(openingID, requestScheme, requestHost string) (*models.JobOpeningShareLinkResponse, error) {
	job, err := s.repo.GetJobOpeningByID(openingID)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(job.ApplyToken) == "" {
		tok, genErr := s.generateApplyToken()
		if genErr != nil {
			return nil, genErr
		}
		job.ApplyToken = tok
		job.UpdatedAt = time.Now()
		if updErr := s.repo.UpdateJobOpening(job); updErr != nil {
			return nil, updErr
		}
	}
	token := strings.TrimSpace(job.ApplyToken)
	baseAPI := strings.TrimSuffix(fmt.Sprintf("%s://%s", requestScheme, requestHost), "/")
	pubJob := fmt.Sprintf("%s/api/v1/recruitment/public/openings/by-token/%s", baseAPI, token)
	pubApply := fmt.Sprintf("%s/api/v1/recruitment/public/apply", baseAPI)

	out := &models.JobOpeningShareLinkResponse{
		OpeningID:       job.ID,
		ApplyToken:      token,
		PublicJobAPIURL: pubJob,
		PublicApplyURL:  pubApply,
		Instructions:    "Share public_job_api_url or build a careers page that reads this JSON. Candidates submit POST public_apply_url with JSON body containing applyToken plus their details.",
	}
	if cfg := config.AppConfig; cfg != nil && cfg.Server.PublicRecruitmentCareersURL != "" {
		out.CareersApplyURL = fmt.Sprintf("%s/apply?token=%s", cfg.Server.PublicRecruitmentCareersURL, token)
	}
	return out, nil
}

// --- Applications ---

func (s *RecruitmentService) SubmitApplication(reqData models.ApplyJobRequest) (*models.JobApplication, error) {
	queueID, queueErr := s.enqueueApplicationSubmission(reqData)
	if queueErr != nil {
		return nil, queueErr
	}

	app, err := s.submitApplicationCore(reqData)
	if err == nil {
		_ = s.markApplicationSubmissionQueueCompleted(queueID, app.ID)
		return app, nil
	}
	_ = s.markApplicationSubmissionQueueFailed(queueID, err)
	return nil, err
}

func (s *RecruitmentService) submitApplicationCore(reqData models.ApplyJobRequest) (*models.JobApplication, error) {
	jobID := strings.TrimSpace(reqData.JobID)
	var job *models.JobOpening
	var err error
	switch {
	case jobID != "":
		job, err = s.repo.GetJobOpeningByID(jobID)
	case strings.TrimSpace(reqData.ApplyToken) != "":
		job, err = s.repo.GetJobOpeningByApplyToken(strings.TrimSpace(reqData.ApplyToken))
		if err == nil && job != nil {
			jobID = job.ID
		}
	default:
		return nil, errors.New("jobId or applyToken is required")
	}
	if err != nil || job == nil || jobID == "" {
		return nil, errors.New("job opening not found")
	}
	if !s.jobIsPubliclyVisible(job) {
		return nil, errors.New("this job is not accepting applications")
	}

	coverLetter := strings.TrimSpace(reqData.CoverLetter)
	if coverLetter == "" {
		coverLetter = strings.TrimSpace(reqData.ApplicationLetter)
	}

	// Handle Resume Upload
	resumeURL := reqData.Resume
	if reqData.Resume != "" && !strings.HasPrefix(reqData.Resume, "http") {
		// Assuming Base64
		b64data := reqData.Resume

		// Extract extension from Base64 header if present (e.g., data:application/pdf;base64,)
		ext := ".pdf" // Default
		if strings.HasPrefix(b64data, "data:") {
			parts := strings.Split(b64data, ";")
			if len(parts) > 0 {
				mime := strings.TrimPrefix(parts[0], "data:")
				switch mime {
				case "application/pdf":
					ext = ".pdf"
				case "application/msword":
					ext = ".doc"
				case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
					ext = ".docx"
				case "image/jpeg":
					ext = ".jpg"
				case "image/png":
					ext = ".png"
				}
			}
		}

		if idx := strings.Index(b64data, ","); idx != -1 {
			b64data = b64data[idx+1:]
		}

		decoded, err := base64.StdEncoding.DecodeString(b64data)
		if err == nil {
			// Upload
			var filename string
			if reqData.ResumeName != "" {
				filename = reqData.ResumeName
			} else {
				filename = fmt.Sprintf("resume_%s_%s_%d%s", reqData.FirstName, reqData.LastName, time.Now().Unix(), ext)
			}

			url, err := s.storage.UploadReader(bytes.NewReader(decoded), filename, "recruitment/resumes", uuid.New().String())
			if err == nil {
				resumeURL = url
			} else {
				return nil, fmt.Errorf("failed to upload resume: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to decode resume base64: %w", err)
		}
	}

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
			ResumeURL:    resumeURL, // Use the uploaded URL
			LinkedInURL:  reqData.LinkedInURL,
			PortfolioURL: reqData.PortfolioURL,
			CoverLetter:  coverLetter,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := s.repo.CreateCandidate(candidate); err != nil {
			return nil, err
		}
	} else {
		// Update candidate resume if provided
		if resumeURL != "" && resumeURL != candidate.ResumeURL {
			candidate.ResumeURL = resumeURL
			candidate.UpdatedAt = time.Now()
			s.repo.UpdateCandidate(candidate)
		}
		if coverLetter != "" {
			candidate.CoverLetter = coverLetter
			candidate.UpdatedAt = time.Now()
			s.repo.UpdateCandidate(candidate)
		}
	}

	// Idempotency guard: if already created for this candidate+job, return it.
	existingApp, appErr := s.repo.GetApplicationByCandidateAndJob(candidate.ID, jobID)
	if appErr == nil && existingApp != nil {
		return existingApp, nil
	}
	if appErr != nil && !errors.Is(appErr, gorm.ErrRecordNotFound) {
		return nil, appErr
	}

	// Create Application
	application := &models.JobApplication{
		ID:           uuid.New().String(),
		CandidateID:  candidate.ID,
		JobOpeningID: jobID,
		Stage:        "Applied",
		AppliedDate:  time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateApplication(application); err != nil {
		return nil, err
	}
	s.hub.Broadcast("candidate.application.created", map[string]interface{}{
		"applicationId": application.ID,
		"candidateId":   application.CandidateID,
		"jobId":         application.JobOpeningID,
		"stage":         application.Stage,
	})

	return application, nil
}

func (s *RecruitmentService) attachApplicantCounts(jobs []models.JobOpening) ([]models.JobOpening, error) {
	if len(jobs) == 0 {
		return jobs, nil
	}
	ids := make([]string, 0, len(jobs))
	for _, j := range jobs {
		if strings.TrimSpace(j.ID) != "" {
			ids = append(ids, j.ID)
		}
	}
	counts, err := s.repo.CountApplicationsByJobOpeningIDs(ids)
	if err != nil {
		return nil, err
	}
	for i := range jobs {
		jobs[i].ApplicantCount = counts[jobs[i].ID]
	}
	return jobs, nil
}

func (s *RecruitmentService) enqueueApplicationSubmission(reqData models.ApplyJobRequest) (string, error) {
	payload, err := json.Marshal(reqData)
	if err != nil {
		return "", err
	}
	row := &models.ApplicationSubmissionQueue{
		ID:             uuid.New().String(),
		JobOpeningID:   strings.TrimSpace(reqData.JobID),
		CandidateEmail: strings.TrimSpace(reqData.Email),
		Payload:        string(payload),
		Status:         models.ApplicationQueueStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.repo.CreateApplicationSubmissionQueue(row); err != nil {
		return "", err
	}
	return row.ID, nil
}

func (s *RecruitmentService) markApplicationSubmissionQueueCompleted(queueID, applicationID string) error {
	row, err := s.repo.GetApplicationSubmissionQueueByID(queueID)
	if err != nil {
		return err
	}
	now := time.Now()
	row.Status = models.ApplicationQueueStatusCompleted
	row.ApplicationID = applicationID
	row.ProcessedAt = &now
	row.LastError = ""
	row.UpdatedAt = now
	return s.repo.SaveApplicationSubmissionQueue(row)
}

func (s *RecruitmentService) markApplicationSubmissionQueueFailed(queueID string, inErr error) error {
	row, err := s.repo.GetApplicationSubmissionQueueByID(queueID)
	if err != nil {
		return err
	}
	row.Status = models.ApplicationQueueStatusFailed
	row.RetryCount++
	row.LastError = inErr.Error()
	row.UpdatedAt = time.Now()
	return s.repo.SaveApplicationSubmissionQueue(row)
}

func (s *RecruitmentService) runApplicationSubmissionQueueWorker() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		rows, err := s.repo.ListPendingApplicationSubmissionQueues(20)
		if err != nil || len(rows) == 0 {
			continue
		}
		for _, row := range rows {
			q := row
			q.Status = models.ApplicationQueueStatusProcessing
			q.UpdatedAt = time.Now()
			if err := s.repo.SaveApplicationSubmissionQueue(&q); err != nil {
				continue
			}
			var req models.ApplyJobRequest
			if err := json.Unmarshal([]byte(q.Payload), &req); err != nil {
				q.Status = models.ApplicationQueueStatusFailed
				q.RetryCount++
				q.LastError = "invalid queue payload: " + err.Error()
				q.UpdatedAt = time.Now()
				_ = s.repo.SaveApplicationSubmissionQueue(&q)
				continue
			}
			app, err := s.submitApplicationCore(req)
			if err != nil {
				q.Status = models.ApplicationQueueStatusFailed
				q.RetryCount++
				q.LastError = err.Error()
				q.UpdatedAt = time.Now()
				_ = s.repo.SaveApplicationSubmissionQueue(&q)
				continue
			}
			now := time.Now()
			q.Status = models.ApplicationQueueStatusCompleted
			q.ApplicationID = app.ID
			q.ProcessedAt = &now
			q.LastError = ""
			q.UpdatedAt = now
			_ = s.repo.SaveApplicationSubmissionQueue(&q)
		}
	}
}

func (s *RecruitmentService) ListCandidates(jobID, stage string) ([]models.Candidate, error) {
	return s.repo.ListCandidates(jobID, stage)
}

// ListCandidatesPaged provides global list + filter + search in one endpoint shape.
func (s *RecruitmentService) ListCandidatesPaged(jobID, stage, keyword string, page, limit int) (map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	rows, total, err := s.repo.ListCandidatesFiltered(jobID, stage, keyword, page, limit)
	if err != nil {
		return nil, err
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 1 {
		totalPages = 1
	}
	return map[string]interface{}{
		"data": rows,
		"meta": map[string]interface{}{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
			"job_id":      jobID,
			"stage":       stage,
			"query":       keyword,
		},
	}, nil
}

// ListRecentApplicants returns the latest job applications (who applied most recently), not candidate account creation order.
func (s *RecruitmentService) ListRecentApplicants(limit int) ([]models.RecentApplicantView, error) {
	if limit < 1 {
		limit = 3
	}
	if limit > 50 {
		limit = 50
	}
	apps, err := s.repo.ListRecentJobApplications(limit)
	if err != nil {
		return nil, err
	}
	out := make([]models.RecentApplicantView, 0, len(apps))
	for _, a := range apps {
		v := models.RecentApplicantView{
			ApplicationID: a.ID,
			CandidateID:   a.CandidateID,
			JobOpeningID:  a.JobOpeningID,
			AppliedAt:     a.AppliedDate,
			Stage:         a.Stage,
			Score:         a.Score,
		}
		if a.Candidate != nil {
			v.FirstName = a.Candidate.FirstName
			v.LastName = a.Candidate.LastName
			v.FullName = strings.TrimSpace(a.Candidate.FirstName + " " + a.Candidate.LastName)
			v.Email = a.Candidate.Email
			v.Phone = a.Candidate.Phone
		}
		if a.JobOpening.ID != "" {
			v.JobTitle = a.JobOpening.JobTitle
			if a.JobOpening.Requisition != nil {
				v.Department = a.JobOpening.Requisition.Department
				v.Location = a.JobOpening.Requisition.Location
			}
		}
		out = append(out, v)
	}
	return out, nil
}

func (s *RecruitmentService) GetCandidateDetails(id string) (*models.Candidate, error) {
	return s.repo.GetCandidateWithDetails(id)
}

// ListCandidatesByInterviewStage returns candidates shortlisted/scheduled for interview round 1/2/3.
func (s *RecruitmentService) ListCandidatesByInterviewStage(stage int, jobID, keyword string, page, limit int) (map[string]interface{}, error) {
	if stage < 1 || stage > 3 {
		return nil, errors.New("stage must be 1, 2, or 3")
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	rows, total, err := s.repo.ListCandidatesByInterviewRound(stage, jobID, keyword, page, limit)
	if err != nil {
		return nil, err
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 1 {
		totalPages = 1
	}
	return map[string]interface{}{
		"data": rows,
		"meta": map[string]interface{}{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
			"stage":       stage,
			"job_id":      jobID,
			"query":       keyword,
		},
	}, nil
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
	if err := s.repo.UpdateApplication(app); err != nil {
		return err
	}
	s.hub.Broadcast("candidate.application.stage_updated", map[string]interface{}{
		"applicationId": app.ID,
		"candidateId":   app.CandidateID,
		"jobId":         app.JobOpeningID,
		"stage":         app.Stage,
	})
	return nil
}

// ActionApplication performs common decisions without requiring clients to map stage names.
func (s *RecruitmentService) ActionApplication(applicationID string, reqData models.ApplicationActionRequest) error {
	app, err := s.repo.GetApplicationByID(applicationID)
	if err != nil {
		return err
	}

	notes := strings.TrimSpace(reqData.Notes)
	appendNotes := func(txt string) {
		if txt == "" {
			return
		}
		if app.Notes != "" {
			app.Notes += "\n---\n"
		}
		app.Notes += txt
	}

	switch strings.TrimSpace(reqData.Action) {
	case models.ApplicationActionAccept:
		app.Stage = models.StageScreening
		appendNotes(notes)
		app.LastStatusDate = time.Now()
		app.UpdatedAt = time.Now()
		if err := s.repo.UpdateApplication(app); err != nil {
			return err
		}
		s.hub.Broadcast("candidate.application.actioned", map[string]interface{}{
			"applicationId": app.ID,
			"candidateId":   app.CandidateID,
			"jobId":         app.JobOpeningID,
			"action":        models.ApplicationActionAccept,
			"stage":         app.Stage,
		})
		return nil

	case models.ApplicationActionReject:
		app.Stage = models.StageRejected
		appendNotes(notes)
		app.LastStatusDate = time.Now()
		app.UpdatedAt = time.Now()
		if err := s.repo.UpdateApplication(app); err != nil {
			return err
		}
		s.hub.Broadcast("candidate.application.actioned", map[string]interface{}{
			"applicationId": app.ID,
			"candidateId":   app.CandidateID,
			"jobId":         app.JobOpeningID,
			"action":        models.ApplicationActionReject,
			"stage":         app.Stage,
		})
		return nil

	case models.ApplicationActionMoveToInterviewStageOne:
		if strings.TrimSpace(reqData.ScheduledDate) == "" || strings.TrimSpace(reqData.ScheduledTime) == "" {
			return errors.New("scheduledDate and scheduledTime are required for move_to_interview_stage_one")
		}
		itype := strings.TrimSpace(reqData.InterviewType)
		if itype == "" {
			itype = "Stage 1"
		}
		_, err := s.ScheduleInterview(models.ScheduleInterviewRequest{
			CandidateID:            app.CandidateID,
			JobID:                  app.JobOpeningID,
			InterviewType:          itype,
			Round:                  1,
			ScheduledDate:          reqData.ScheduledDate,
			ScheduledTime:          reqData.ScheduledTime,
			Duration:               reqData.Duration,
			Mode:                   reqData.Mode,
			InterviewerEmployeeIDs: reqData.InterviewerEmployeeIDs,
			Interviewers:           reqData.Interviewers,
			IsReinterview:          false,
		})
		if err != nil {
			return err
		}
		app.Stage = models.StageInterview
		appendNotes(notes)
		app.LastStatusDate = time.Now()
		app.UpdatedAt = time.Now()
		if err := s.repo.UpdateApplication(app); err != nil {
			return err
		}
		s.hub.Broadcast("candidate.application.actioned", map[string]interface{}{
			"applicationId": app.ID,
			"candidateId":   app.CandidateID,
			"jobId":         app.JobOpeningID,
			"action":        models.ApplicationActionMoveToInterviewStageOne,
			"stage":         app.Stage,
		})
		return nil

	default:
		return errors.New("invalid action: use accept, reject, or move_to_interview_stage_one")
	}
}

// --- Interviews ---

func workEmailFromEmployee(e *empModels.Employee) string {
	if e.WorkEmail != nil && strings.TrimSpace(*e.WorkEmail) != "" {
		return strings.TrimSpace(*e.WorkEmail)
	}
	if e.PersonalEmail != nil && strings.TrimSpace(*e.PersonalEmail) != "" {
		return strings.TrimSpace(*e.PersonalEmail)
	}
	return ""
}

func (s *RecruitmentService) buildInterviewPanelFromEmployeeIDs(ids []string) ([]models.InterviewerAssignment, []string, error) {
	seen := make(map[string]bool)
	var panel []models.InterviewerAssignment
	var emails []string
	for _, raw := range ids {
		id := strings.TrimSpace(raw)
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		emp, err := s.employeeRepo.FindByEmployeeID(id)
		if err != nil || emp == nil {
			return nil, nil, fmt.Errorf("employee not found: %s", id)
		}
		name := strings.TrimSpace(strings.TrimSpace(emp.FirstName) + " " + strings.TrimSpace(emp.LastName))
		panel = append(panel, models.InterviewerAssignment{
			EmployeeID:   emp.EmployeeID,
			DisplayName:  name,
			WorkEmail:    workEmailFromEmployee(emp),
			DepartmentID: emp.DepartmentID,
			PositionID:   emp.PositionID,
		})
		if em := workEmailFromEmployee(emp); em != "" {
			emails = append(emails, em)
		}
	}
	return panel, emails, nil
}

func mergeUniqueEmails(fromPanel []string, legacy []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, e := range append(fromPanel, legacy...) {
		e = strings.TrimSpace(e)
		if e == "" || seen[e] {
			continue
		}
		seen[e] = true
		out = append(out, e)
	}
	return out
}

func (s *RecruitmentService) validateInterviewSchedule(app *models.JobApplication, round int, isReinterview bool) error {
	if round < 1 || round > models.InterviewMaxRound {
		return fmt.Errorf("round must be between 1 and %d", models.InterviewMaxRound)
	}
	blocking, err := s.repo.CountBlockingInterviews(app.ID)
	if err != nil {
		return err
	}
	if blocking > 0 {
		return errors.New("another interview is still pending panel feedback or HR decision; resolve it before scheduling again")
	}
	if round > 1 {
		prev, err := s.repo.GetLatestInterviewForRound(app.ID, round-1)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("cannot schedule round %d before round %d is HR-approved", round, round-1)
			}
			return err
		}
		if prev.Status != models.InterviewStatusHrApprovedNextRound {
			return fmt.Errorf("previous round must be HR-approved (%s) before scheduling round %d", models.InterviewStatusHrApprovedNextRound, round)
		}
	}
	if isReinterview {
		last, err := s.repo.GetLatestInterviewForRound(app.ID, round)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("cannot schedule re-interview: no prior interview for this round")
			}
			return err
		}
		if last.Status != models.InterviewStatusHrRequestedReinterview {
			return errors.New("HR must request a re-interview for this round before rescheduling")
		}
		return nil
	}
	maxA, err := s.repo.MaxAttemptForRound(app.ID, round)
	if err != nil {
		return err
	}
	if maxA > 0 {
		return errors.New("this round already has interviews recorded; use isReinterview=true after HR requests a re-interview")
	}
	return nil
}

func (s *RecruitmentService) ScheduleInterview(reqData models.ScheduleInterviewRequest) (*models.Interview, error) {
	app, err := s.repo.GetApplicationByCandidateAndJob(reqData.CandidateID, reqData.JobID)
	if err != nil {
		return nil, errors.New("active application not found for this candidate and job")
	}

	if err := s.validateInterviewSchedule(app, reqData.Round, reqData.IsReinterview); err != nil {
		return nil, err
	}

	panel, panelEmails, err := s.buildInterviewPanelFromEmployeeIDs(reqData.InterviewerEmployeeIDs)
	if err != nil {
		return nil, err
	}
	legacyEmails := mergeUniqueEmails(panelEmails, reqData.Interviewers)
	if len(panel) == 0 && len(legacyEmails) == 0 {
		return nil, errors.New("provide interviewerEmployeeIds (employee IDs from HR) and/or interviewers (email list)")
	}

	maxA, _ := s.repo.MaxAttemptForRound(app.ID, reqData.Round)
	attempt := 1
	if reqData.IsReinterview {
		attempt = maxA + 1
		if attempt < 2 {
			attempt = 2
		}
	}

	interview := &models.Interview{
		ID:               uuid.New().String(),
		CandidateID:      reqData.CandidateID,
		JobOpeningID:     reqData.JobID,
		ApplicationID:    app.ID,
		InterviewType:    reqData.InterviewType,
		Round:            reqData.Round,
		Attempt:          attempt,
		ScheduledTime:    reqData.ScheduledTime,
		DurationMin:      reqData.Duration,
		Mode:             reqData.Mode,
		Interviewers:     legacyEmails,
		InterviewerPanel: panel,
		Status:           models.InterviewStatusScheduled,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if t, err := parseFlexibleDate(reqData.ScheduledDate); err == nil && !t.IsZero() {
		interview.ScheduledDate = t
	}

	if err := s.repo.CreateInterview(interview); err != nil {
		return nil, err
	}
	s.hub.Broadcast("interview.scheduled", map[string]interface{}{
		"interviewId":   interview.ID,
		"applicationId": interview.ApplicationID,
		"candidateId":   interview.CandidateID,
		"jobId":         interview.JobOpeningID,
		"round":         interview.Round,
		"status":        interview.Status,
	})

	candidate, _ := s.repo.GetCandidateByID(reqData.CandidateID)
	if candidate != nil {
		subject := fmt.Sprintf("Interview Scheduled: Round %d", reqData.Round)
		if attempt > 1 {
			subject = fmt.Sprintf("Interview Scheduled: Round %d (attempt %d)", reqData.Round, attempt)
		}
		body := fmt.Sprintf("Dear %s,\n\nYour interview for %s is scheduled on %s at %s.\nMode: %s\n\nBest regards,\nRecruitment Team",
			candidate.FirstName, fmt.Sprintf("round %d", reqData.Round), reqData.ScheduledDate, reqData.ScheduledTime, reqData.Mode)
		s.email.SendEmail(candidate.Email, subject, body)
	}

	return interview, nil
}

func (s *RecruitmentService) SubmitFeedback(interviewID string, reqData models.SubmitFeedbackRequest) error {
	interview, err := s.repo.GetInterviewByID(interviewID)
	if err != nil {
		return err
	}

	if interview.Status != models.InterviewStatusScheduled {
		if interview.Status == models.InterviewStatusPanelCompleted || interview.Status == models.InterviewStatusLegacyCompleted {
			return errors.New("feedback was already submitted for this interview")
		}
		return errors.New("feedback can only be submitted while the interview is scheduled")
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

	now := time.Now()
	interview.Feedback = feedback
	interview.Status = models.InterviewStatusPanelCompleted
	interview.PanelCompletedAt = &now
	if err := s.repo.UpdateInterview(interview); err != nil {
		return err
	}
	s.hub.Broadcast("interview.feedback_submitted", map[string]interface{}{
		"interviewId":   interview.ID,
		"applicationId": interview.ApplicationID,
		"candidateId":   interview.CandidateID,
		"jobId":         interview.JobOpeningID,
		"round":         interview.Round,
		"status":        interview.Status,
		"rating":        reqData.Rating,
		"remarks":       reqData.Comments,
	})
	return nil
}

// SubmitHrInterviewDecision records HR remarks and unlocks the next round, a same-round re-interview, or rejects the candidate.
func (s *RecruitmentService) SubmitHrInterviewDecision(interviewID string, hrDeciderEmployeeID string, req models.HRInterviewDecisionRequest) error {
	remarks := strings.TrimSpace(req.Remarks)
	if remarks == "" {
		return errors.New("remarks are required")
	}
	decision := strings.TrimSpace(req.Decision)
	switch decision {
	case models.HrInterviewDecisionProceedNextRound, models.HrInterviewDecisionScheduleReinterview, models.HrInterviewDecisionRejectCandidate:
	default:
		return errors.New("invalid decision: use proceed_next_round, schedule_reinterview, or reject_candidate")
	}

	inv, err := s.repo.GetInterviewByID(interviewID)
	if err != nil {
		return err
	}
	if inv.Status != models.InterviewStatusPanelCompleted && inv.Status != models.InterviewStatusLegacyCompleted {
		return errors.New("HR decision is only allowed after panel feedback is submitted")
	}

	now := time.Now()
	inv.HrRemarks = remarks
	inv.HrDecision = decision
	inv.HrDecidedAt = &now
	inv.HrDeciderEmployeeID = hrDeciderEmployeeID

	app, err := s.repo.GetApplicationByID(inv.ApplicationID)
	if err != nil {
		return err
	}

	switch decision {
	case models.HrInterviewDecisionProceedNextRound:
		inv.Status = models.InterviewStatusHrApprovedNextRound
		if inv.Round >= models.InterviewMaxRound {
			app.Stage = models.StageInterviewCompleted
		}
	case models.HrInterviewDecisionScheduleReinterview:
		inv.Status = models.InterviewStatusHrRequestedReinterview
	case models.HrInterviewDecisionRejectCandidate:
		inv.Status = models.InterviewStatusHrRejected
		app.Stage = models.StageRejected
	}

	app.LastStatusDate = now
	app.UpdatedAt = now

	if err := s.repo.UpdateInterview(inv); err != nil {
		return err
	}
	if err := s.repo.UpdateApplication(app); err != nil {
		return err
	}
	s.hub.Broadcast("interview.hr_decision", map[string]interface{}{
		"interviewId":   inv.ID,
		"applicationId": inv.ApplicationID,
		"candidateId":   inv.CandidateID,
		"jobId":         inv.JobOpeningID,
		"round":         inv.Round,
		"decision":      inv.HrDecision,
		"status":        inv.Status,
		"remarks":       inv.HrRemarks,
		"applicationStage": app.Stage,
	})
	return nil
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

	// Fetch candidate to get email
	candidate, err := s.repo.GetCandidateByID(offer.CandidateID)
	if err != nil {
		return fmt.Errorf("candidate not found: %w", err)
	}

	offer.Status = "Sent"
	now := time.Now()
	offer.SentAt = &now

	// Send Email
	subject := "Job Offer from Our Company"
	body := fmt.Sprintf("Dear %s,\n\nWe are pleased to offer you the position. Please find the details attached.\n\nBest regards,\nRecruitment Team", candidate.FirstName)
	s.email.SendEmail(candidate.Email, subject, body)

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
