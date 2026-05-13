package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/models"
	"github.com/google/uuid"
)

// ListApplicationsForJobOpening returns paginated applications with candidates for one opening.
func (s *RecruitmentService) ListApplicationsForJobOpening(jobOpeningID, stage string, page, limit int) (map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	apps, total, err := s.repo.ListApplicationsForJobOpening(jobOpeningID, stage, page, limit)
	if err != nil {
		return nil, err
	}
	tp := int((total + int64(limit) - 1) / int64(limit))
	if tp < 1 {
		tp = 1
	}
	return map[string]interface{}{
		"data": apps,
		"meta": map[string]interface{}{
			"total":          total,
			"page":           page,
			"limit":          limit,
			"total_pages":    tp,
			"job_opening_id": jobOpeningID,
			"stage":          stage,
		},
	}, nil
}

// ListInterviewsForJobOpening lists scheduled interviews for a job opening.
func (s *RecruitmentService) ListInterviewsForJobOpening(jobOpeningID string) ([]models.Interview, error) {
	return s.repo.ListInterviewsForJobOpening(jobOpeningID)
}

// ListInterviewsPaged returns interviews across all stages by default.
// Use round (1/2/3) to filter stage-specific interviews.
func (s *RecruitmentService) ListInterviewsPaged(round int, status, jobID, candidateID string, page, limit int) (map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	rows, total, err := s.repo.ListInterviewsPaginated(round, status, jobID, candidateID, page, limit)
	if err != nil {
		return nil, err
	}
	tp := int((total + int64(limit) - 1) / int64(limit))
	if tp < 1 {
		tp = 1
	}
	return map[string]interface{}{
		"data": rows,
		"meta": map[string]interface{}{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": tp,
			"round":       round,
			"status":      status,
			"job_id":      jobID,
			"candidate_id": candidateID,
		},
	}, nil
}

func (s *RecruitmentService) GetInterviewDetails(id string) (*models.Interview, error) {
	return s.repo.GetInterviewByID(id)
}

// ListOffersForJobOpening lists draft/sent offers tied to a job opening.
func (s *RecruitmentService) ListOffersForJobOpening(jobOpeningID string) ([]models.Offer, error) {
	return s.repo.ListOffersForJobOpening(jobOpeningID)
}

// ScheduleInterviewManual creates candidate + application + interview when the person was not in the pipeline.
func (s *RecruitmentService) ScheduleInterviewManual(req models.ManualScheduleInterviewRequest) (*models.Interview, error) {
	if _, err := s.repo.GetJobOpeningByID(req.JobOpeningID); err != nil {
		return nil, errors.New("job opening not found")
	}

	email := strings.TrimSpace(req.Email)
	existingCand, _ := s.repo.GetCandidateByEmail(email)
	if existingCand != nil {
		app, _ := s.repo.GetApplicationByCandidateAndJob(existingCand.ID, req.JobOpeningID)
		if app != nil {
			return nil, errors.New("this email already has an application for this job opening")
		}
	}

	var cand *models.Candidate
	if existingCand != nil {
		cand = existingCand
	} else {
		cand = &models.Candidate{
			ID:        uuid.New().String(),
			FirstName: strings.TrimSpace(req.FirstName),
			LastName:  strings.TrimSpace(req.LastName),
			Email:     strings.TrimSpace(req.Email),
			Phone:     strings.TrimSpace(req.Phone),
			ResumeURL: strings.TrimSpace(req.ResumeURL),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.repo.CreateCandidate(cand); err != nil {
			return nil, err
		}
	}

	app := &models.JobApplication{
		ID:             uuid.New().String(),
		CandidateID:    cand.ID,
		JobOpeningID:   req.JobOpeningID,
		Stage:          models.StageInterview,
		Notes:          strings.TrimSpace(req.Notes),
		AppliedDate:    time.Now(),
		LastStatusDate: time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.repo.CreateApplication(app); err != nil {
		return nil, err
	}

	return s.ScheduleInterview(models.ScheduleInterviewRequest{
		CandidateID:            cand.ID,
		JobID:                  req.JobOpeningID,
		InterviewType:          req.InterviewType,
		Round:                  req.Round,
		ScheduledDate:          req.ScheduledDate,
		ScheduledTime:          req.ScheduledTime,
		Duration:               req.Duration,
		Mode:                   req.Mode,
		Interviewers:           req.Interviewers,
		InterviewerEmployeeIDs: req.InterviewerEmployeeIDs,
		IsReinterview:          req.IsReinterview,
	})
}

// ScheduleInterviewsBatch schedules the same slot for multiple existing applications (pipeline).
func (s *RecruitmentService) ScheduleInterviewsBatch(req models.BatchScheduleInterviewRequest) ([]models.Interview, error) {
	updateStage := true
	if req.UpdateStage != nil {
		updateStage = *req.UpdateStage
	}
	var out []models.Interview
	for _, appID := range req.ApplicationIDs {
		app, err := s.repo.GetApplicationByID(strings.TrimSpace(appID))
		if err != nil {
			return nil, fmt.Errorf("application %s: %w", appID, err)
		}
		if app.JobOpeningID != req.JobOpeningID {
			return nil, fmt.Errorf("application %s does not belong to job opening %s", appID, req.JobOpeningID)
		}
		inv, err := s.ScheduleInterview(models.ScheduleInterviewRequest{
			CandidateID:            app.CandidateID,
			JobID:                  req.JobOpeningID,
			InterviewType:          req.InterviewType,
			Round:                  req.Round,
			ScheduledDate:          req.ScheduledDate,
			ScheduledTime:          req.ScheduledTime,
			Duration:               req.Duration,
			Mode:                   req.Mode,
			Interviewers:           req.Interviewers,
			InterviewerEmployeeIDs: req.InterviewerEmployeeIDs,
			IsReinterview:          req.IsReinterview,
		})
		if err != nil {
			return nil, fmt.Errorf("application %s: %w", appID, err)
		}
		if updateStage {
			app.Stage = models.StageInterview
			app.LastStatusDate = time.Now()
			app.UpdatedAt = time.Now()
			if err := s.repo.UpdateApplication(app); err != nil {
				return nil, err
			}
		}
		out = append(out, *inv)
	}
	return out, nil
}

// MoveApplicationToTalentPool moves an application to the TalentPool stage and optionally syncs the talent_pool_candidates table.
func (s *RecruitmentService) MoveApplicationToTalentPool(applicationID string, req models.MoveApplicationToTalentPoolRequest) error {
	app, err := s.repo.GetApplicationByID(applicationID)
	if err != nil {
		return err
	}
	if app.Stage == models.StageTalentPool {
		return errors.New("application is already in talent pool stage")
	}
	if strings.TrimSpace(req.Notes) != "" {
		if app.Notes != "" {
			app.Notes += "\n---\n"
		}
		app.Notes += strings.TrimSpace(req.Notes)
	}
	app.Stage = models.StageTalentPool
	app.LastStatusDate = time.Now()
	app.UpdatedAt = time.Now()
	if err := s.repo.UpdateApplication(app); err != nil {
		return err
	}
	sync := true
	if req.SyncTalentPoolRow != nil {
		sync = *req.SyncTalentPoolRow
	}
	if !sync {
		return nil
	}
	return s.upsertTalentPoolFromApplication(app, strings.TrimSpace(req.InternalPoolNotes))
}

func (s *RecruitmentService) upsertTalentPoolFromApplication(app *models.JobApplication, poolNotes string) error {
	cand, err := s.repo.GetCandidateByID(app.CandidateID)
	if err != nil || cand == nil {
		return fmt.Errorf("candidate not found for application")
	}
	name := strings.TrimSpace(cand.FirstName + " " + cand.LastName)
	email := strings.TrimSpace(cand.Email)

	existing, err := s.repo.GetTalentPoolByEmail(email)
	if err != nil {
		return err
	}
	now := time.Now()
	if existing != nil {
		existing.Name = name
		existing.ResumeURL = cand.ResumeURL
		existing.Skills = cand.Skills
		existing.SourceJobOpeningID = app.JobOpeningID
		existing.SourceApplicationID = app.ID
		if poolNotes != "" {
			if existing.InternalNotes != "" {
				existing.InternalNotes += "\n---\n"
			}
			existing.InternalNotes += poolNotes
		}
		existing.UpdatedAt = now
		return s.repo.SaveTalentPoolCandidate(existing)
	}

	row := &models.TalentPoolCandidate{
		ID:                  uuid.New().String(),
		Name:                name,
		Email:               cand.Email,
		Skills:              cand.Skills,
		ResumeURL:           cand.ResumeURL,
		SourceJobOpeningID:  app.JobOpeningID,
		SourceApplicationID: app.ID,
		InternalNotes:       poolNotes,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	return s.repo.AddToTalentPool(row)
}

// ListTalentPoolPaged lists talent pool candidates for HR outreach.
func (s *RecruitmentService) ListTalentPoolPaged(page, limit int) (map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	rows, total, err := s.repo.ListTalentPoolPaged(page, limit)
	if err != nil {
		return nil, err
	}
	tp := int((total + int64(limit) - 1) / int64(limit))
	if tp < 1 {
		tp = 1
	}
	return map[string]interface{}{
		"data": rows,
		"meta": map[string]interface{}{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": tp,
		},
	}, nil
}

// ContactTalentPoolCandidate sends an email and records last contacted time.
func (s *RecruitmentService) ContactTalentPoolCandidate(id string, req models.ContactTalentPoolRequest) error {
	row, err := s.repo.GetTalentPoolCandidateByID(id)
	if err != nil {
		return err
	}
	if err := s.email.SendEmail(row.Email, req.Subject, req.Message); err != nil {
		return err
	}
	now := time.Now()
	row.LastContactedAt = &now
	row.UpdatedAt = now
	return s.repo.SaveTalentPoolCandidate(row)
}
