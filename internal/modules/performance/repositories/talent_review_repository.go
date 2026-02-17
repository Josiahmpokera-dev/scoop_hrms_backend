package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"gorm.io/gorm"
)

// TalentReviewRepository handles database operations for talent reviews, calibration sessions, and succession plans.
type TalentReviewRepository struct {
	db *gorm.DB
}

// NewTalentReviewRepository creates a new TalentReviewRepository.
func NewTalentReviewRepository() *TalentReviewRepository {
	return &TalentReviewRepository{db: database.GetDB()}
}

// FindByEmployeeID finds a talent review by employee ID.
func (r *TalentReviewRepository) FindByEmployeeID(employeeID uint) (*models.TalentReview, error) {
	var review models.TalentReview
	err := r.db.Where("employee_id = ?", employeeID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// Upsert creates or updates a talent review by employee_id (unique).
func (r *TalentReviewRepository) Upsert(review *models.TalentReview) error {
	existing, err := r.FindByEmployeeID(review.EmployeeID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if existing != nil {
		review.ID = existing.ID
		return r.db.Save(review).Error
	}
	return r.db.Create(review).Error
}

// Update updates a talent review.
func (r *TalentReviewRepository) Update(review *models.TalentReview) error {
	return r.db.Save(review).Error
}

// List returns talent reviews with pagination and optional filters.
func (r *TalentReviewRepository) List(department string, box *int, criticalRole *bool, riskOfLoss string, page, pageSize int) ([]models.TalentReview, int64, error) {
	var reviews []models.TalentReview
	var total int64

	query := r.db.Model(&models.TalentReview{})

	if department != "" {
		query = query.Where("department = ?", department)
	}
	if box != nil {
		query = query.Where("box = ?", *box)
	}
	if criticalRole != nil {
		query = query.Where("critical_role = ?", *criticalRole)
	}
	if riskOfLoss != "" {
		query = query.Where("risk_of_loss = ?", riskOfLoss)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&reviews).Error

	return reviews, total, err
}

// CreateCalibrationSession creates a new calibration session.
func (r *TalentReviewRepository) CreateCalibrationSession(session *models.CalibrationSession) error {
	return r.db.Create(session).Error
}

// CreateSuccessionPlan creates a new succession plan.
func (r *TalentReviewRepository) CreateSuccessionPlan(plan *models.SuccessionPlan) error {
	return r.db.Create(plan).Error
}

// ListSuccessionPlans returns all succession plans for an employee.
func (r *TalentReviewRepository) ListSuccessionPlans(employeeID uint) ([]models.SuccessionPlan, error) {
	var plans []models.SuccessionPlan
	err := r.db.Where("employee_id = ?", employeeID).Order("created_at DESC").Find(&plans).Error
	return plans, err
}
