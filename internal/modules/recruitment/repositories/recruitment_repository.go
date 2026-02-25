package repositories

import (
	"errors"
	"strings"

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
	return r.db.Create(req).Error
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
	return r.db.Save(req).Error
}

// --- Job Openings ---

func (r *RecruitmentRepository) CreateJobOpening(job *models.JobOpening) error {
	return r.db.Create(job).Error
}

func (r *RecruitmentRepository) GetJobOpeningByID(id string) (*models.JobOpening, error) {
	var job models.JobOpening
	err := r.db.First(&job, "id = ?", id).Error
	return &job, err
}

func (r *RecruitmentRepository) ListJobOpenings(status, department string) ([]models.JobOpening, error) {
	var jobs []models.JobOpening
	query := r.db.Model(&models.JobOpening{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	// Note: JobOpening doesn't have Department directly, it's in Requisition.
	// We might need a join if filtering by department is strictly required on this table.
	// For now, let's assume if department is passed, we join.
	if department != "" {
		query = query.Joins("JOIN job_requisitions ON job_requisitions.id = job_openings.requisition_id").
			Where("job_requisitions.department = ?", department)
	}

	err := query.Find(&jobs).Error
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
	err := r.db.Preload("JobOpening").Preload("Interviews").Preload("Offer").First(&app, "id = ?", id).Error
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

func (r *RecruitmentRepository) UpdateApplication(app *models.JobApplication) error {
	return r.db.Save(app).Error
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
