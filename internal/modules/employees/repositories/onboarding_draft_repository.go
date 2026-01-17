package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

type OnboardingDraftRepository struct {
	db *gorm.DB
}

func NewOnboardingDraftRepository() *OnboardingDraftRepository {
	return &OnboardingDraftRepository{
		db: database.GetDB(),
	}
}

// Create creates a new onboarding draft
func (r *OnboardingDraftRepository) Create(draft *models.EmployeeOnboardingDraft) error {
	return r.db.Create(draft).Error
}

// FindByID finds a draft by ID
func (r *OnboardingDraftRepository) FindByID(id uint) (*models.EmployeeOnboardingDraft, error) {
	var draft models.EmployeeOnboardingDraft
	err := r.db.First(&draft, id).Error
	if err != nil {
		return nil, err
	}
	return &draft, nil
}

// FindByTenantID finds all drafts for a tenant
func (r *OnboardingDraftRepository) FindByTenantID(tenantID *uint) ([]models.EmployeeOnboardingDraft, error) {
	var drafts []models.EmployeeOnboardingDraft
	query := r.db.Model(&models.EmployeeOnboardingDraft{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.Order("created_at DESC").Find(&drafts).Error
	return drafts, err
}

// Update updates an onboarding draft
func (r *OnboardingDraftRepository) Update(draft *models.EmployeeOnboardingDraft) error {
	// Use Save to update all fields, ensuring employee_id and all other fields are persisted
	return r.db.Save(draft).Error
}

// Delete soft deletes an onboarding draft
func (r *OnboardingDraftRepository) Delete(id uint) error {
	return r.db.Delete(&models.EmployeeOnboardingDraft{}, id).Error
}

// FindByEmployeeIDFinal finds draft by final employee ID (after completion)
func (r *OnboardingDraftRepository) FindByEmployeeIDFinal(employeeID uint) (*models.EmployeeOnboardingDraft, error) {
	var draft models.EmployeeOnboardingDraft
	err := r.db.Where("employee_id_final = ?", employeeID).First(&draft).Error
	if err != nil {
		return nil, err
	}
	return &draft, nil
}

// FindByEmployeeIDString finds draft by employee ID string (before completion)
// Uses case-insensitive matching and trims whitespace
func (r *OnboardingDraftRepository) FindByEmployeeIDString(employeeID string, tenantID *uint) (*models.EmployeeOnboardingDraft, error) {
	var draft models.EmployeeOnboardingDraft
	// Trim whitespace and use case-insensitive matching
	query := r.db.Where("LOWER(TRIM(employee_id)) = LOWER(TRIM(?))", employeeID)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.First(&draft).Error
	if err != nil {
		return nil, err
	}
	return &draft, nil
}

// FindIncompleteDrafts finds all incomplete drafts for a tenant
func (r *OnboardingDraftRepository) FindIncompleteDrafts(tenantID *uint) ([]models.EmployeeOnboardingDraft, error) {
	var drafts []models.EmployeeOnboardingDraft
	query := r.db.Model(&models.EmployeeOnboardingDraft{}).
		Where("is_completed = ?", false)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.Order("updated_at DESC").Find(&drafts).Error
	return drafts, err
}

// FindIncompleteDraftsWithPagination finds incomplete drafts with pagination
func (r *OnboardingDraftRepository) FindIncompleteDraftsWithPagination(tenantID *uint, page, pageSize int) ([]models.EmployeeOnboardingDraft, int64, error) {
	var drafts []models.EmployeeOnboardingDraft
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.EmployeeOnboardingDraft{}).
		Where("is_completed = ?", false)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	err := query.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&drafts).Error
	return drafts, total, err
}
