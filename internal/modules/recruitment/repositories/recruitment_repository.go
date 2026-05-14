package repositories

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/models"
	"gorm.io/gorm"
)

type RecruitmentRepository struct {
	db *gorm.DB
}

func NewRecruitmentRepository() *RecruitmentRepository {
	return &RecruitmentRepository{
		db: database.GetDB(),
	}
}

// --- Requisitions ---

func (r *RecruitmentRepository) CreateRequisition(req *models.JobRequisition) error {
	err := r.db.Session(&gorm.Session{FullSaveAssociations: false}).
		Omit("ApprovalFlow").
		Create(req).Error
	if err == nil {
		return nil
	}
	el := strings.ToLower(err.Error())
	if strings.Contains(el, "priority") && (strings.Contains(el, "column") || strings.Contains(el, "does not exist")) {
		// Older DBs created before the priority column; insert without it.
		reqLegacy := *req
		reqLegacy.Priority = ""
		return r.db.Session(&gorm.Session{FullSaveAssociations: false}).
			Omit("ApprovalFlow", "Priority").
			Create(&reqLegacy).Error
	}
	return err
}

func (r *RecruitmentRepository) GetRequisitionByID(id string) (*models.JobRequisition, error) {
	var req models.JobRequisition
	err := r.db.Preload("ApprovalFlow").First(&req, "id = ?", id).Error
	return &req, err
}

func (r *RecruitmentRepository) ListRequisitions(status, department string, page, limit int) ([]models.JobRequisition, int64, error) {
	var reqs []models.JobRequisition
	var total int64
	query := r.db.Model(&models.JobRequisition{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if department != "" {
		query = query.Where("department = ?", department)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Order("created_at desc").Find(&reqs).Error
	return reqs, total, err
}

func (r *RecruitmentRepository) UpdateRequisition(req *models.JobRequisition) error {
	// Persist nested ApprovalFlow rows when approvers append steps (ApproveRequisition).
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(req).Error
}

// --- Job Openings ---

func (r *RecruitmentRepository) CreateJobOpening(job *models.JobOpening) error {
	return r.db.Create(job).Error
}

func (r *RecruitmentRepository) GetJobOpeningByID(id string) (*models.JobOpening, error) {
	var job models.JobOpening
	err := r.db.Preload("Requisition").First(&job, "id = ?", id).Error
	return &job, err
}

// GetJobOpeningByApplyToken loads an opening by its public apply_token.
func (r *RecruitmentRepository) GetJobOpeningByApplyToken(token string) (*models.JobOpening, error) {
	var job models.JobOpening
	err := r.db.Preload("Requisition").Where("apply_token = ?", token).First(&job).Error
	return &job, err
}

func (r *RecruitmentRepository) IsApplyTokenTaken(token string) (bool, error) {
	var count int64
	err := r.db.Model(&models.JobOpening{}).Where("apply_token = ?", token).Count(&count).Error
	return count > 0, err
}

func (r *RecruitmentRepository) ListJobOpenings(status, department string) ([]models.JobOpening, error) {
	jobs, _, err := r.ListJobOpeningsPaginated(status, department, "", 0, 0)
	return jobs, err
}

// ListJobOpeningsPaginated returns job openings with optional filters.
// page/limit 0 = no pagination (all rows).
func (r *RecruitmentRepository) ListJobOpeningsPaginated(status, department, requisitionID string, page, limit int) ([]models.JobOpening, int64, error) {
	var jobs []models.JobOpening
	var total int64
	query := r.db.Model(&models.JobOpening{})

	if status != "" {
		switch status {
		case models.JobOpeningStatusActive:
			query = query.Where("status IN ?", []string{models.JobOpeningStatusActive, models.LegacyStatusPublished})
		default:
			query = query.Where("status = ?", status)
		}
	}
	if requisitionID != "" {
		query = query.Where("requisition_id = ?", requisitionID)
	}
	if department != "" {
		query = query.Joins("JOIN job_requisitions ON job_requisitions.id = job_openings.requisition_id").
			Where("job_requisitions.department = ?", department)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q2 := query.Session(&gorm.Session{})
	if limit > 0 {
		offset := (page - 1) * limit
		if offset < 0 {
			offset = 0
		}
		q2 = q2.Offset(offset).Limit(limit)
	}
	err := q2.Preload("Requisition").Order("created_at DESC").Find(&jobs).Error
	return jobs, total, err
}

// SearchJobOpeningsPaginated returns openings filtered by keyword + existing filters.
// Keyword currently matches job_title and job_description.
func (r *RecruitmentRepository) SearchJobOpeningsPaginated(keyword, status, department, requisitionID string, page, limit int) ([]models.JobOpening, int64, error) {
	var jobs []models.JobOpening
	var total int64
	query := r.db.Model(&models.JobOpening{})

	kw := strings.TrimSpace(keyword)
	if kw != "" {
		like := "%" + strings.ToLower(kw) + "%"
		query = query.Where("(LOWER(job_title) LIKE ? OR LOWER(job_description) LIKE ?)", like, like)
	}

	if status != "" {
		switch status {
		case models.JobOpeningStatusActive:
			query = query.Where("status IN ?", []string{models.JobOpeningStatusActive, models.LegacyStatusPublished})
		default:
			query = query.Where("status = ?", status)
		}
	}
	if requisitionID != "" {
		query = query.Where("requisition_id = ?", requisitionID)
	}
	if department != "" {
		query = query.Joins("JOIN job_requisitions ON job_requisitions.id = job_openings.requisition_id").
			Where("job_requisitions.department = ?", department)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q2 := query.Session(&gorm.Session{})
	if limit > 0 {
		offset := (page - 1) * limit
		if offset < 0 {
			offset = 0
		}
		q2 = q2.Offset(offset).Limit(limit)
	}
	err := q2.Preload("Requisition").Order("created_at DESC").Find(&jobs).Error
	return jobs, total, err
}

// ListPublicActiveOpenings returns openings candidates can apply to (Active or legacy Published, not expired).
func (r *RecruitmentRepository) ListPublicActiveOpenings(department string, now time.Time) ([]models.JobOpening, error) {
	var jobs []models.JobOpening
	query := r.db.Model(&models.JobOpening{}).
		Where("status IN ?", []string{models.JobOpeningStatusActive, models.LegacyStatusPublished}).
		Where("(expiry_date IS NULL OR expiry_date = ? OR expiry_date > ?)", time.Time{}, now)
	if department != "" {
		query = query.Joins("JOIN job_requisitions ON job_requisitions.id = job_openings.requisition_id").
			Where("job_requisitions.department = ?", department)
	}
	err := query.Preload("Requisition").Order("created_at DESC").Find(&jobs).Error
	return jobs, err
}

func (r *RecruitmentRepository) UpdateJobOpening(job *models.JobOpening) error {
	return r.db.Save(job).Error
}

// --- Candidates & Applications ---

func (r *RecruitmentRepository) CreateCandidate(candidate *models.Candidate) error {
	return r.db.Create(candidate).Error
}

func (r *RecruitmentRepository) UpdateCandidate(candidate *models.Candidate) error {
	return r.db.Save(candidate).Error
}

func (r *RecruitmentRepository) GetCandidateByEmail(email string) (*models.Candidate, error) {
	var candidate models.Candidate
	err := r.db.Where("email = ?", email).First(&candidate).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &candidate, err
}

func (r *RecruitmentRepository) GetCandidateByID(id string) (*models.Candidate, error) {
	var candidate models.Candidate
	err := r.db.First(&candidate, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &candidate, err
}

func (r *RecruitmentRepository) CreateApplication(app *models.JobApplication) error {
	return r.db.Create(app).Error
}

func (r *RecruitmentRepository) GetApplicationByID(id string) (*models.JobApplication, error) {
	var app models.JobApplication
	err := r.db.Preload("JobOpening").Preload("Candidate").Preload("Interviews").Preload("Offer").First(&app, "id = ?", id).Error
	return &app, err
}

func (r *RecruitmentRepository) GetApplicationByCandidateAndJob(candidateID, jobID string) (*models.JobApplication, error) {
	var app models.JobApplication
	err := r.db.Where("candidate_id = ? AND job_opening_id = ?", candidateID, jobID).First(&app).Error
	return &app, err
}

func (r *RecruitmentRepository) ListCandidates(jobID, stage string) ([]models.Candidate, error) {
	var candidates []models.Candidate

	// Start with the base query for Candidates
	query := r.db.Model(&models.Candidate{}).Preload("Applications")

	// Join with JobApplication to filter by JobID or Stage
	if jobID != "" || stage != "" {
		query = query.Joins("JOIN job_applications ON job_applications.candidate_id = candidates.id")

		if jobID != "" {
			query = query.Where("job_applications.job_opening_id = ?", jobID)
		}
		if stage != "" {
			query = query.Where("job_applications.stage = ?", stage)
		}
	}

	err := query.Distinct().Find(&candidates).Error
	return candidates, err
}

// ListCandidatesFiltered returns candidates globally, with optional filters and keyword search.
func (r *RecruitmentRepository) ListCandidatesFiltered(jobID, stage, keyword string, page, limit int) ([]models.Candidate, int64, error) {
	var rows []models.Candidate
	var total int64

	base := r.db.Model(&models.Candidate{})
	needsJoinApplications := strings.TrimSpace(jobID) != "" || strings.TrimSpace(stage) != ""
	if needsJoinApplications {
		base = base.Joins("JOIN job_applications ON job_applications.candidate_id = candidates.id")
	}
	if strings.TrimSpace(jobID) != "" {
		base = base.Where("job_applications.job_opening_id = ?", strings.TrimSpace(jobID))
	}
	if strings.TrimSpace(stage) != "" {
		base = base.Where("job_applications.stage = ?", strings.TrimSpace(stage))
	}
	if kw := strings.TrimSpace(strings.ToLower(keyword)); kw != "" {
		like := "%" + kw + "%"
		base = base.Where(
			"(LOWER(candidates.first_name) LIKE ? OR LOWER(candidates.last_name) LIKE ? OR LOWER(candidates.email) LIKE ? OR LOWER(candidates.phone) LIKE ?)",
			like, like, like, like,
		)
	}

	countQ := base.Session(&gorm.Session{})
	if err := countQ.Distinct("candidates.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q := base.Session(&gorm.Session{}).
		Select("candidates.*").
		Distinct().
		Preload("Applications").
		Order("candidates.created_at DESC")
	if limit > 0 {
		offset := (page - 1) * limit
		if offset < 0 {
			offset = 0
		}
		q = q.Offset(offset).Limit(limit)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// GetCandidateWithDetails returns candidate profile plus recruitment details.
func (r *RecruitmentRepository) GetCandidateWithDetails(candidateID string) (*models.Candidate, error) {
	var row models.Candidate
	err := r.db.
		Preload("Applications").
		Preload("Applications.JobOpening").
		Preload("Applications.Interviews").
		Preload("Applications.Offer").
		First(&row, "id = ?", candidateID).Error
	return &row, err
}

// ListCandidatesByInterviewRound returns candidates that have interviews in the requested round (1..3).
func (r *RecruitmentRepository) ListCandidatesByInterviewRound(round int, jobID, keyword string, page, limit int) ([]models.Candidate, int64, error) {
	var rows []models.Candidate
	var total int64

	base := r.db.Model(&models.Candidate{}).
		Joins("JOIN job_applications ON job_applications.candidate_id = candidates.id").
		Joins("JOIN interviews ON interviews.application_id = job_applications.id").
		Where("interviews.round = ?", round)

	if strings.TrimSpace(jobID) != "" {
		base = base.Where("job_applications.job_opening_id = ?", strings.TrimSpace(jobID))
	}
	if kw := strings.TrimSpace(strings.ToLower(keyword)); kw != "" {
		like := "%" + kw + "%"
		base = base.Where(
			"(LOWER(candidates.first_name) LIKE ? OR LOWER(candidates.last_name) LIKE ? OR LOWER(candidates.email) LIKE ? OR LOWER(candidates.phone) LIKE ?)",
			like, like, like, like,
		)
	}

	countQ := base.Session(&gorm.Session{})
	if err := countQ.Distinct("candidates.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q := base.Session(&gorm.Session{}).
		Select("candidates.*").
		Distinct().
		Preload("Applications", func(db *gorm.DB) *gorm.DB {
			if strings.TrimSpace(jobID) != "" {
				db = db.Where("job_opening_id = ?", strings.TrimSpace(jobID))
			}
			return db
		}).
		Preload("Applications.JobOpening").
		Preload("Applications.Interviews", "round = ?", round).
		Order("candidates.created_at DESC")

	if limit > 0 {
		offset := (page - 1) * limit
		if offset < 0 {
			offset = 0
		}
		q = q.Offset(offset).Limit(limit)
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *RecruitmentRepository) UpdateApplication(app *models.JobApplication) error {
	return r.db.Save(app).Error
}

// ListRecentJobApplications returns the most recent applications by applied_date (then created_at).
func (r *RecruitmentRepository) ListRecentJobApplications(limit int) ([]models.JobApplication, error) {
	if limit < 1 {
		limit = 3
	}
	var apps []models.JobApplication
	err := r.db.Model(&models.JobApplication{}).
		Preload("Candidate").
		Preload("JobOpening.Requisition").
		Where("job_applications.deleted_at IS NULL").
		Order("job_applications.applied_date DESC, job_applications.created_at DESC").
		Limit(limit).
		Find(&apps).Error
	return apps, err
}

// CountApplicationsByJobOpeningIDs returns application counts grouped by job opening id.
func (r *RecruitmentRepository) CountApplicationsByJobOpeningIDs(jobOpeningIDs []string) (map[string]int64, error) {
	out := make(map[string]int64)
	if len(jobOpeningIDs) == 0 {
		return out, nil
	}
	type row struct {
		JobOpeningID string
		Count        int64
	}
	var rows []row
	err := r.db.Model(&models.JobApplication{}).
		Select("job_opening_id, COUNT(*) as count").
		Where("job_opening_id IN ?", jobOpeningIDs).
		Group("job_opening_id").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.JobOpeningID] = r.Count
	}
	return out, nil
}

// --- Interviews ---

func (r *RecruitmentRepository) CreateInterview(interview *models.Interview) error {
	return r.db.Create(interview).Error
}

func (r *RecruitmentRepository) GetInterviewByID(id string) (*models.Interview, error) {
	var interview models.Interview
	err := r.db.Preload("Feedback").First(&interview, "id = ?", id).Error
	return &interview, err
}

func (r *RecruitmentRepository) UpdateInterview(interview *models.Interview) error {
	return r.db.Save(interview).Error
}

// ListInterviewsPaginated lists interviews globally with optional filters.
func (r *RecruitmentRepository) ListInterviewsPaginated(round int, status, jobOpeningID, candidateID string, page, limit int) ([]models.Interview, int64, error) {
	var rows []models.Interview
	var total int64

	q := r.db.Model(&models.Interview{})
	if round > 0 {
		q = q.Where("round = ?", round)
	}
	if strings.TrimSpace(status) != "" {
		q = q.Where("status = ?", strings.TrimSpace(status))
	}
	if strings.TrimSpace(jobOpeningID) != "" {
		q = q.Where("job_opening_id = ?", strings.TrimSpace(jobOpeningID))
	}
	if strings.TrimSpace(candidateID) != "" {
		q = q.Where("candidate_id = ?", strings.TrimSpace(candidateID))
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q2 := q.Session(&gorm.Session{}).Preload("Feedback").Order("scheduled_date DESC, created_at DESC")
	if limit > 0 {
		offset := (page - 1) * limit
		if offset < 0 {
			offset = 0
		}
		q2 = q2.Offset(offset).Limit(limit)
	}

	if err := q2.Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// MaxAttemptForRound returns the highest attempt number recorded for this application and round.
func (r *RecruitmentRepository) MaxAttemptForRound(applicationID string, round int) (int, error) {
	var max sql.NullInt64
	err := r.db.Model(&models.Interview{}).
		Where("application_id = ? AND round = ?", applicationID, round).
		Select("MAX(attempt)").
		Scan(&max).Error
	if err != nil {
		return 0, err
	}
	if !max.Valid {
		return 0, nil
	}
	return int(max.Int64), nil
}

// GetLatestInterviewForRound returns the interview row with the highest attempt for that round.
func (r *RecruitmentRepository) GetLatestInterviewForRound(applicationID string, round int) (*models.Interview, error) {
	var inv models.Interview
	err := r.db.Where("application_id = ? AND round = ?", applicationID, round).
		Order("attempt DESC, created_at DESC").
		First(&inv).Error
	return &inv, err
}

// CountBlockingInterviews counts sessions that still need panel completion or HR action before scheduling another interview.
func (r *RecruitmentRepository) CountBlockingInterviews(applicationID string) (int64, error) {
	var n int64
	err := r.db.Model(&models.Interview{}).
		Where("application_id = ? AND status IN ?", applicationID, blockingInterviewStatuses()).
		Count(&n).Error
	return n, err
}

func blockingInterviewStatuses() []string {
	return []string{
		models.InterviewStatusScheduled,
		models.InterviewStatusPanelCompleted,
		models.InterviewStatusLegacyCompleted,
	}
}

// --- Offers ---

func (r *RecruitmentRepository) CreateOffer(offer *models.Offer) error {
	return r.db.Create(offer).Error
}

func (r *RecruitmentRepository) GetOfferByID(id string) (*models.Offer, error) {
	var offer models.Offer
	err := r.db.First(&offer, "id = ?", id).Error
	return &offer, err
}

func (r *RecruitmentRepository) UpdateOffer(offer *models.Offer) error {
	return r.db.Save(offer).Error
}

// --- Talent Pool ---

func (r *RecruitmentRepository) AddToTalentPool(candidate *models.TalentPoolCandidate) error {
	return r.db.Create(candidate).Error
}

func (r *RecruitmentRepository) SearchTalentPool(skills []string, location string, expMin int) ([]models.TalentPoolCandidate, error) {
	var candidates []models.TalentPoolCandidate
	query := r.db.Model(&models.TalentPoolCandidate{})

	if len(skills) > 0 {
		// Postgres specific array overlap or partial match.
		// For simplicity/portability (assuming JSON storage), we might use LIKE or custom query.
		// Since we use GORM serializer:json, it's stored as text/json.
		// A robust implementation would use Postgres GIN index.
		// Here, simpler approach: OR-based LIKE for each skill (not efficient but functional for small scale)
		for _, skill := range skills {
			query = query.Where("skills LIKE ?", "%"+strings.TrimSpace(skill)+"%")
		}
	}
	if location != "" {
		query = query.Where("location LIKE ?", "%"+location+"%")
	}
	if expMin > 0 {
		query = query.Where("experience_min >= ?", expMin)
	}

	err := query.Find(&candidates).Error
	return candidates, err
}

// ListApplicationsForJobOpening returns applications for one job opening with candidate preloaded.
func (r *RecruitmentRepository) ListApplicationsForJobOpening(jobOpeningID, stage string, page, limit int) ([]models.JobApplication, int64, error) {
	var apps []models.JobApplication
	var total int64
	base := r.db.Model(&models.JobApplication{}).Where("job_opening_id = ?", jobOpeningID)
	if stage != "" {
		base = base.Where("stage = ?", stage)
	}
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	q := r.db.Preload("Candidate").Where("job_opening_id = ?", jobOpeningID)
	if stage != "" {
		q = q.Where("stage = ?", stage)
	}
	if limit > 0 {
		offset := (page - 1) * limit
		if offset < 0 {
			offset = 0
		}
		q = q.Offset(offset).Limit(limit)
	}
	err := q.Order("applied_date DESC").Find(&apps).Error
	return apps, total, err
}

// ListInterviewsForJobOpening lists interviews linked to a job opening.
func (r *RecruitmentRepository) ListInterviewsForJobOpening(jobOpeningID string) ([]models.Interview, error) {
	var rows []models.Interview
	err := r.db.Preload("Feedback").Where("job_opening_id = ?", jobOpeningID).Order("scheduled_date ASC").Find(&rows).Error
	return rows, err
}

// ListOffersForJobOpening lists offers for a job opening.
func (r *RecruitmentRepository) ListOffersForJobOpening(jobOpeningID string) ([]models.Offer, error) {
	var rows []models.Offer
	err := r.db.Where("job_opening_id = ?", jobOpeningID).Order("created_at DESC").Find(&rows).Error
	return rows, err
}

// ListTalentPoolPaged lists talent pool rows with pagination.
func (r *RecruitmentRepository) ListTalentPoolPaged(page, limit int) ([]models.TalentPoolCandidate, int64, error) {
	var rows []models.TalentPoolCandidate
	var total int64
	q := r.db.Model(&models.TalentPoolCandidate{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	err := r.db.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&rows).Error
	return rows, total, err
}

// GetTalentPoolCandidateByID returns one talent pool row.
func (r *RecruitmentRepository) GetTalentPoolCandidateByID(id string) (*models.TalentPoolCandidate, error) {
	var row models.TalentPoolCandidate
	err := r.db.First(&row, "id = ?", id).Error
	return &row, err
}

// GetTalentPoolByEmail finds an existing talent pool row by email.
func (r *RecruitmentRepository) GetTalentPoolByEmail(email string) (*models.TalentPoolCandidate, error) {
	var row models.TalentPoolCandidate
	err := r.db.Where("email = ?", email).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &row, err
}

// SaveTalentPoolCandidate saves updates to an existing talent pool row.
func (r *RecruitmentRepository) SaveTalentPoolCandidate(row *models.TalentPoolCandidate) error {
	return r.db.Save(row).Error
}

// --- Talent pool pipeline (candidates with applications in TalentPool stage) ---

func (r *RecruitmentRepository) talentPoolPipelineSubquery(keyword, jobOpeningID, location, skill string) *gorm.DB {
	sq := r.db.Table("job_applications ja").
		Select("ja.candidate_id, MAX(ja.updated_at) AS mx").
		Joins("JOIN candidates c ON c.id = ja.candidate_id AND c.deleted_at IS NULL").
		Where("ja.deleted_at IS NULL AND ja.stage = ?", models.StageTalentPool)

	if kw := strings.TrimSpace(strings.ToLower(keyword)); kw != "" {
		like := "%" + kw + "%"
		sq = sq.Where(
			"(LOWER(c.first_name) LIKE ? OR LOWER(c.last_name) LIKE ? OR LOWER(c.email) LIKE ? OR LOWER(COALESCE(c.phone, '')) LIKE ?)",
			like, like, like, like,
		)
	}
	if jid := strings.TrimSpace(jobOpeningID); jid != "" {
		sq = sq.Where("ja.job_opening_id = ?", jid)
	}
	if loc := strings.TrimSpace(location); loc != "" {
		like := "%" + strings.ToLower(loc) + "%"
		sq = sq.Joins("JOIN job_openings jo ON jo.id = ja.job_opening_id AND jo.deleted_at IS NULL").
			Joins("JOIN job_requisitions jr ON jr.id = jo.requisition_id AND jr.deleted_at IS NULL").
			Where("LOWER(COALESCE(jr.location, '')) LIKE ?", like)
	}
	if sk := strings.TrimSpace(skill); sk != "" {
		sq = sq.Where("CAST(c.skills AS TEXT) ILIKE ?", "%"+sk+"%")
	}
	return sq.Group("ja.candidate_id")
}

// CountTalentPoolPipelineCandidates counts distinct candidates with at least one TalentPool-stage application.
func (r *RecruitmentRepository) CountTalentPoolPipelineCandidates(keyword, jobOpeningID, location, skill string) (int64, error) {
	sq := r.talentPoolPipelineSubquery(keyword, jobOpeningID, location, skill)
	var total int64
	err := r.db.Table("(?) AS t", sq).Count(&total).Error
	return total, err
}

// ListTalentPoolPipelineCandidateIDs returns candidate IDs ordered by most recent TalentPool application activity.
func (r *RecruitmentRepository) ListTalentPoolPipelineCandidateIDs(keyword, jobOpeningID, location, skill string, page, limit int) ([]string, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit
	sq := r.talentPoolPipelineSubquery(keyword, jobOpeningID, location, skill)
	var ids []string
	err := r.db.Table("(?) AS ranked", sq).
		Order("mx DESC").
		Offset(offset).Limit(limit).
		Pluck("candidate_id", &ids).Error
	return ids, err
}

// ListJobApplicationsTalentPoolByCandidateIDs loads TalentPool-stage applications for the given candidates.
func (r *RecruitmentRepository) ListJobApplicationsTalentPoolByCandidateIDs(candidateIDs []string) ([]models.JobApplication, error) {
	var apps []models.JobApplication
	if len(candidateIDs) == 0 {
		return apps, nil
	}
	err := r.db.Where("candidate_id IN ? AND deleted_at IS NULL AND stage = ?", candidateIDs, models.StageTalentPool).
		Preload("JobOpening.Requisition").
		Order("updated_at DESC").
		Find(&apps).Error
	return apps, err
}

// FindCandidatesByIDs loads candidates by primary keys (order not preserved).
func (r *RecruitmentRepository) FindCandidatesByIDs(ids []string) ([]models.Candidate, error) {
	var rows []models.Candidate
	if len(ids) == 0 {
		return rows, nil
	}
	err := r.db.Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

// --- Application submission queue ---

func (r *RecruitmentRepository) CreateApplicationSubmissionQueue(row *models.ApplicationSubmissionQueue) error {
	return r.db.Create(row).Error
}

func (r *RecruitmentRepository) GetApplicationSubmissionQueueByID(id string) (*models.ApplicationSubmissionQueue, error) {
	var row models.ApplicationSubmissionQueue
	err := r.db.First(&row, "id = ?", id).Error
	return &row, err
}

func (r *RecruitmentRepository) SaveApplicationSubmissionQueue(row *models.ApplicationSubmissionQueue) error {
	return r.db.Save(row).Error
}

func (r *RecruitmentRepository) ListPendingApplicationSubmissionQueues(limit int) ([]models.ApplicationSubmissionQueue, error) {
	var rows []models.ApplicationSubmissionQueue
	if limit <= 0 {
		limit = 20
	}
	err := r.db.Where("status IN ?", []string{
		models.ApplicationQueueStatusPending,
		models.ApplicationQueueStatusFailed,
	}).
		Where("retry_count < ?", 5).
		Order("created_at ASC").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}
