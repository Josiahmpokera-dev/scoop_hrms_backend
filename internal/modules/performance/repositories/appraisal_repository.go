package repositories

import (
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"gorm.io/gorm"
)

// AppraisalRepository handles database operations for appraisal cycles, workflow steps, and appraisals.
type AppraisalRepository struct {
	db *gorm.DB
}

// NewAppraisalRepository creates a new AppraisalRepository.
func NewAppraisalRepository() *AppraisalRepository {
	return &AppraisalRepository{db: database.GetDB()}
}

// CreateCycle creates a new appraisal cycle.
func (r *AppraisalRepository) CreateCycle(cycle *models.AppraisalCycle) error {
	return r.db.Create(cycle).Error
}

// FindCycleByID finds an appraisal cycle by ID with WorkflowSteps preloaded.
func (r *AppraisalRepository) FindCycleByID(id uint) (*models.AppraisalCycle, error) {
	var cycle models.AppraisalCycle
	err := r.db.Preload("WorkflowSteps").First(&cycle, id).Error
	if err != nil {
		return nil, err
	}
	return &cycle, nil
}

// UpdateCycle updates an appraisal cycle.
func (r *AppraisalRepository) UpdateCycle(cycle *models.AppraisalCycle) error {
	return r.db.Save(cycle).Error
}

// DeleteCycle soft deletes an appraisal cycle.
func (r *AppraisalRepository) DeleteCycle(id uint) error {
	return r.db.Delete(&models.AppraisalCycle{}, id).Error
}

// ListCycles returns appraisal cycles with pagination and optional filters; preloads WorkflowSteps.
func (r *AppraisalRepository) ListCycles(status, cycleType string, year int, page, pageSize int) ([]models.AppraisalCycle, int64, error) {
	var cycles []models.AppraisalCycle
	var total int64

	query := r.db.Model(&models.AppraisalCycle{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if cycleType != "" {
		query = query.Where("cycle_type = ?", cycleType)
	}
	if year > 0 {
		query = query.Where("year = ?", year)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("WorkflowSteps").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&cycles).Error

	return cycles, total, err
}

// FindActiveCycle returns the active appraisal cycle (status = Active) with WorkflowSteps preloaded.
func (r *AppraisalRepository) FindActiveCycle() (*models.AppraisalCycle, error) {
	var cycle models.AppraisalCycle
	err := r.db.Preload("WorkflowSteps").Where("status = ?", "Active").First(&cycle).Error
	if err != nil {
		return nil, err
	}
	return &cycle, nil
}

// GetNextCycleCode returns the next appraisal cycle code (e.g. AC-001).
func (r *AppraisalRepository) GetNextCycleCode() string {
	var count int64
	r.db.Model(&models.AppraisalCycle{}).Count(&count)
	return fmt.Sprintf("AC-%03d", count+1)
}

// CreateWorkflowStep creates a new appraisal workflow step.
func (r *AppraisalRepository) CreateWorkflowStep(step *models.AppraisalWorkflowStep) error {
	return r.db.Create(step).Error
}

// UpdateWorkflowStep updates an appraisal workflow step.
func (r *AppraisalRepository) UpdateWorkflowStep(step *models.AppraisalWorkflowStep) error {
	return r.db.Save(step).Error
}

// CreateAppraisal creates a new appraisal.
func (r *AppraisalRepository) CreateAppraisal(appraisal *models.Appraisal) error {
	return r.db.Create(appraisal).Error
}

// FindAppraisalByID finds an appraisal by ID.
func (r *AppraisalRepository) FindAppraisalByID(id uint) (*models.Appraisal, error) {
	var appraisal models.Appraisal
	err := r.db.First(&appraisal, id).Error
	if err != nil {
		return nil, err
	}
	return &appraisal, nil
}

// FindAppraisalByCode finds an appraisal by code.
func (r *AppraisalRepository) FindAppraisalByCode(code string) (*models.Appraisal, error) {
	var appraisal models.Appraisal
	err := r.db.Where("code = ?", code).First(&appraisal).Error
	if err != nil {
		return nil, err
	}
	return &appraisal, nil
}

// UpdateAppraisal updates an appraisal.
func (r *AppraisalRepository) UpdateAppraisal(appraisal *models.Appraisal) error {
	return r.db.Save(appraisal).Error
}

// ListAppraisals returns appraisals with pagination and optional filters.
// mine: when true, filter by employeeID (current user's appraisals).
func (r *AppraisalRepository) ListAppraisals(employeeID *uint, cycleID *uint, status, department string, mine bool, page, pageSize int) ([]models.Appraisal, int64, error) {
	var appraisals []models.Appraisal
	var total int64

	query := r.db.Model(&models.Appraisal{})

	if mine && employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	} else if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	if cycleID != nil {
		query = query.Where("appraisal_cycle_id = ?", *cycleID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if department != "" {
		query = query.Where("department = ?", department)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&appraisals).Error

	return appraisals, total, err
}

// GetNextAppraisalCode returns the next appraisal code (e.g. APR-001).
func (r *AppraisalRepository) GetNextAppraisalCode() string {
	var count int64
	r.db.Model(&models.Appraisal{}).Count(&count)
	return fmt.Sprintf("APR-%03d", count+1)
}

// CountAppraisalsByStatus returns counts grouped by status for the given employee (or all if nil).
func (r *AppraisalRepository) CountAppraisalsByStatus(employeeID *uint) (map[string]int64, error) {
	type result struct {
		Status string
		Count  int64
	}
	var results []result
	query := r.db.Model(&models.Appraisal{}).Select("status, count(*) as count").Group("status")
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}
	counts := make(map[string]int64)
	for _, row := range results {
		counts[row.Status] = row.Count
	}
	return counts, nil
}
