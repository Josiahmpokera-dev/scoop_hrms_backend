package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"gorm.io/gorm"
)

// SalaryStructureRepository handles database operations for salary structures
type SalaryStructureRepository struct {
	db *gorm.DB
}

// NewSalaryStructureRepository creates a new repository instance
func NewSalaryStructureRepository() *SalaryStructureRepository {
	return &SalaryStructureRepository{
		db: database.DB,
	}
}

// Create creates a new salary structure
func (r *SalaryStructureRepository) Create(structure *models.SalaryStructure) error {
	return r.db.Create(structure).Error
}

// GetByID retrieves a salary structure by ID
func (r *SalaryStructureRepository) GetByID(id uint, tenantID *uint) (*models.SalaryStructure, error) {
	var structure models.SalaryStructure
	query := r.db.Where("id = ?", id)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.First(&structure).Error; err != nil {
		return nil, err
	}
	return &structure, nil
}

// Update updates a salary structure
func (r *SalaryStructureRepository) Update(structure *models.SalaryStructure) error {
	return r.db.Save(structure).Error
}

// Delete soft-deletes a salary structure
func (r *SalaryStructureRepository) Delete(id uint, tenantID *uint) error {
	query := r.db.Where("id = ?", id)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	return query.Delete(&models.SalaryStructure{}).Error
}

// List retrieves salary structures with pagination and filters
func (r *SalaryStructureRepository) List(tenantID *uint, isActive *bool, grade, location, country string, page, pageSize int) ([]models.SalaryStructure, int64, error) {
	var structures []models.SalaryStructure
	var total int64

	query := r.db.Model(&models.SalaryStructure{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	if grade != "" {
		query = query.Where("grade = ?", grade)
	}
	if location != "" {
		query = query.Where("location = ?", location)
	}
	if country != "" {
		query = query.Where("country = ?", country)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("grade ASC, template_name ASC").Offset(offset).Limit(pageSize).Find(&structures).Error; err != nil {
		return nil, 0, err
	}

	return structures, total, nil
}

// SalaryComponentRepository handles database operations for salary components
type SalaryComponentRepository struct {
	db *gorm.DB
}

// NewSalaryComponentRepository creates a new repository instance
func NewSalaryComponentRepository() *SalaryComponentRepository {
	return &SalaryComponentRepository{
		db: database.DB,
	}
}

// Create creates a new salary component
func (r *SalaryComponentRepository) Create(component *models.SalaryComponent) error {
	return r.db.Create(component).Error
}

// GetByID retrieves a salary component by ID
func (r *SalaryComponentRepository) GetByID(id uint, tenantID *uint) (*models.SalaryComponent, error) {
	var component models.SalaryComponent
	query := r.db.Where("id = ?", id)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.First(&component).Error; err != nil {
		return nil, err
	}
	return &component, nil
}

// GetByCode retrieves a salary component by code
func (r *SalaryComponentRepository) GetByCode(code string, tenantID *uint) (*models.SalaryComponent, error) {
	var component models.SalaryComponent
	query := r.db.Where("component_code = ?", code)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.First(&component).Error; err != nil {
		return nil, err
	}
	return &component, nil
}

// Update updates a salary component
func (r *SalaryComponentRepository) Update(component *models.SalaryComponent) error {
	return r.db.Save(component).Error
}

// Delete soft-deletes a salary component
func (r *SalaryComponentRepository) Delete(id uint, tenantID *uint) error {
	query := r.db.Where("id = ?", id)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	return query.Delete(&models.SalaryComponent{}).Error
}

// List retrieves salary components with pagination and filters
func (r *SalaryComponentRepository) List(tenantID *uint, componentType string, isActive, isStatutory *bool, country string, page, pageSize int) ([]models.SalaryComponent, int64, error) {
	var components []models.SalaryComponent
	var total int64

	query := r.db.Model(&models.SalaryComponent{})
	if tenantID != nil {
		query = query.Where("tenant_id = ? OR tenant_id IS NULL", *tenantID)
	}
	if componentType != "" {
		query = query.Where("component_type = ?", componentType)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	if isStatutory != nil {
		query = query.Where("is_statutory = ?", *isStatutory)
	}
	if country != "" {
		query = query.Where("country = ? OR country = 'All'", country)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("sort_order ASC, component_name ASC").Offset(offset).Limit(pageSize).Find(&components).Error; err != nil {
		return nil, 0, err
	}

	return components, total, nil
}

// GetAllActive retrieves all active salary components
func (r *SalaryComponentRepository) GetAllActive(tenantID *uint, country string) ([]models.SalaryComponent, error) {
	var components []models.SalaryComponent
	query := r.db.Where("is_active = ?", true)
	if tenantID != nil {
		query = query.Where("tenant_id = ? OR tenant_id IS NULL", *tenantID)
	}
	if country != "" {
		query = query.Where("country = ? OR country = 'All'", country)
	}
	if err := query.Order("sort_order ASC").Find(&components).Error; err != nil {
		return nil, err
	}
	return components, nil
}

// CountActiveSalaryStructures returns the count of active salary structures
func (r *SalaryStructureRepository) CountActiveSalaryStructures(tenantID *uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.SalaryStructure{}).Where("is_active = ?", true)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetAllActiveStructures returns all active salary structures
func (r *SalaryStructureRepository) GetAllActiveStructures(tenantID *uint) ([]models.SalaryStructure, error) {
	var structures []models.SalaryStructure
	query := r.db.Where("is_active = ?", true)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Order("grade ASC").Find(&structures).Error; err != nil {
		return nil, err
	}
	return structures, nil
}

func (r *SalaryStructureRepository) GetByJobPositionID(tenantID *uint, jobPositionID uint) (*models.SalaryStructure, error) {
	var structure models.SalaryStructure
	query := r.db.Where("is_active = ?", true).Where("job_position_id = ?", jobPositionID)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.First(&structure).Error; err != nil {
		return nil, err
	}
	return &structure, nil
}

// GetDistinctGrades returns all unique grades from active salary structures
func (r *SalaryStructureRepository) GetDistinctGrades(tenantID *uint) ([]string, error) {
	var grades []string
	query := r.db.Model(&models.SalaryStructure{}).
		Select("DISTINCT grade").
		Where("is_active = ?", true).
		Where("grade IS NOT NULL AND grade != ''")
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Pluck("grade", &grades).Error; err != nil {
		return nil, err
	}
	return grades, nil
}
