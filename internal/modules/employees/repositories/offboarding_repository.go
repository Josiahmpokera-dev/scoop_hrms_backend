package repositories

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

// OffboardingWorkflowRepository handles database operations for offboarding workflows
type OffboardingWorkflowRepository struct {
	db *gorm.DB
}

// NewOffboardingWorkflowRepository creates a new offboarding workflow repository
func NewOffboardingWorkflowRepository() *OffboardingWorkflowRepository {
	return &OffboardingWorkflowRepository{
		db: database.GetDB(),
	}
}

// Create creates a new offboarding workflow
func (r *OffboardingWorkflowRepository) Create(workflow *models.OffboardingWorkflow) error {
	return r.db.Create(workflow).Error
}

// FindByID finds an offboarding workflow by ID
func (r *OffboardingWorkflowRepository) FindByID(id uint) (*models.OffboardingWorkflow, error) {
	var workflow models.OffboardingWorkflow
	err := r.db.First(&workflow, id).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

// FindByOffboardingID finds an offboarding workflow by offboarding_id (string)
func (r *OffboardingWorkflowRepository) FindByOffboardingID(offboardingID string, tenantID *uint) (*models.OffboardingWorkflow, error) {
	var workflow models.OffboardingWorkflow
	query := r.db.Where("offboarding_id = ?", offboardingID)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.First(&workflow).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

// FindByEmployeeID finds an active offboarding workflow by employee ID
func (r *OffboardingWorkflowRepository) FindByEmployeeID(employeeID string, tenantID *uint) (*models.OffboardingWorkflow, error) {
	var workflow models.OffboardingWorkflow
	query := r.db.Where("employee_id = ? AND status != ?", employeeID, "Completed")
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.Order("created_at DESC").First(&workflow).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

// Update updates an offboarding workflow
func (r *OffboardingWorkflowRepository) Update(workflow *models.OffboardingWorkflow) error {
	return r.db.Save(workflow).Error
}

// List retrieves a paginated list of offboarding workflows with filters
func (r *OffboardingWorkflowRepository) List(req *models.ListOffboardingWorkflowsRequest, tenantID *uint) ([]models.OffboardingWorkflow, int64, error) {
	var workflows []models.OffboardingWorkflow
	var total int64

	query := r.db.Model(&models.OffboardingWorkflow{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply filters
	if req.Status != nil && *req.Status != "" {
		query = query.Where("status = ?", *req.Status)
	}
	if req.SeparationType != nil && *req.SeparationType != "" {
		query = query.Where("separation_type = ?", *req.SeparationType)
	}
	if req.ClearanceStatus != nil && *req.ClearanceStatus != "" {
		query = query.Where("clearance_status = ?", *req.ClearanceStatus)
	}
	if req.FinalSettlement != nil && *req.FinalSettlement != "" {
		query = query.Where("final_settlement_status = ?", *req.FinalSettlement)
	}
	if req.DateFrom != nil && *req.DateFrom != "" {
		query = query.Where("last_working_date >= ?", *req.DateFrom)
	}
	if req.DateTo != nil && *req.DateTo != "" {
		query = query.Where("last_working_date <= ?", *req.DateTo)
	}

	// Search by employee name, employee ID, or designation (requires join with employees table)
	if req.Search != nil && *req.Search != "" {
		searchTerm := "%" + strings.ToLower(*req.Search) + "%"
		query = query.Where("LOWER(employee_id) LIKE ?", searchTerm)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	page := 1
	if req.Page > 0 {
		page = req.Page
	}
	pageSize := 20
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}
	offset := (page - 1) * pageSize

	// Fetch workflows
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&workflows).Error
	if err != nil {
		return nil, 0, err
	}

	return workflows, total, nil
}

// GenerateOffboardingID generates a unique offboarding ID
func (r *OffboardingWorkflowRepository) GenerateOffboardingID() (string, error) {
	var count int64
	r.db.Model(&models.OffboardingWorkflow{}).Count(&count)
	return fmt.Sprintf("off-%03d", count+1), nil
}

// OffboardingClearanceRepository handles database operations for offboarding clearances
type OffboardingClearanceRepository struct {
	db *gorm.DB
}

// NewOffboardingClearanceRepository creates a new offboarding clearance repository
func NewOffboardingClearanceRepository() *OffboardingClearanceRepository {
	return &OffboardingClearanceRepository{
		db: database.GetDB(),
	}
}

// Create creates a new clearance
func (r *OffboardingClearanceRepository) Create(clearance *models.OffboardingClearance) error {
	return r.db.Create(clearance).Error
}

// CreateBatch creates multiple clearances at once
func (r *OffboardingClearanceRepository) CreateBatch(clearances []models.OffboardingClearance) error {
	if len(clearances) == 0 {
		return nil
	}
	return r.db.Create(&clearances).Error
}

// FindByID finds a clearance by ID
func (r *OffboardingClearanceRepository) FindByID(id uint) (*models.OffboardingClearance, error) {
	var clearance models.OffboardingClearance
	err := r.db.First(&clearance, id).Error
	if err != nil {
		return nil, err
	}
	return &clearance, nil
}

// FindByClearanceID finds a clearance by clearance_id (string)
func (r *OffboardingClearanceRepository) FindByClearanceID(clearanceID string, tenantID *uint) (*models.OffboardingClearance, error) {
	var clearance models.OffboardingClearance
	query := r.db.Where("clearance_id = ?", clearanceID)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.First(&clearance).Error
	if err != nil {
		return nil, err
	}
	return &clearance, nil
}

// FindByOffboardingID finds all clearances for an offboarding workflow
func (r *OffboardingClearanceRepository) FindByOffboardingID(offboardingID string, tenantID *uint) ([]models.OffboardingClearance, error) {
	var clearances []models.OffboardingClearance
	query := r.db.Where("offboarding_id = ?", offboardingID)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.Order("created_at ASC").Find(&clearances).Error
	if err != nil {
		return nil, err
	}
	return clearances, nil
}

// Update updates a clearance
func (r *OffboardingClearanceRepository) Update(clearance *models.OffboardingClearance) error {
	return r.db.Save(clearance).Error
}

// GenerateClearanceID generates a unique clearance ID
func (r *OffboardingClearanceRepository) GenerateClearanceID() (string, error) {
	var count int64
	r.db.Model(&models.OffboardingClearance{}).Count(&count)
	return fmt.Sprintf("clear-%03d", count+1), nil
}

// OffboardingAssetReturnRepository handles database operations for asset returns
type OffboardingAssetReturnRepository struct {
	db *gorm.DB
}

// NewOffboardingAssetReturnRepository creates a new asset return repository
func NewOffboardingAssetReturnRepository() *OffboardingAssetReturnRepository {
	return &OffboardingAssetReturnRepository{
		db: database.GetDB(),
	}
}

// Create creates a new asset return record
func (r *OffboardingAssetReturnRepository) Create(assetReturn *models.OffboardingAssetReturn) error {
	return r.db.Create(assetReturn).Error
}

// FindByID finds an asset return by ID
func (r *OffboardingAssetReturnRepository) FindByID(id uint) (*models.OffboardingAssetReturn, error) {
	var assetReturn models.OffboardingAssetReturn
	err := r.db.First(&assetReturn, id).Error
	if err != nil {
		return nil, err
	}
	return &assetReturn, nil
}

// FindByOffboardingID finds all asset returns for an offboarding workflow
func (r *OffboardingAssetReturnRepository) FindByOffboardingID(offboardingID string, tenantID *uint) ([]models.OffboardingAssetReturn, error) {
	var assetReturns []models.OffboardingAssetReturn
	query := r.db.Where("offboarding_id = ?", offboardingID)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.Order("created_at ASC").Find(&assetReturns).Error
	if err != nil {
		return nil, err
	}
	return assetReturns, nil
}

// FindByAssetIDAndOffboardingID finds an asset return by asset ID and offboarding ID
func (r *OffboardingAssetReturnRepository) FindByAssetIDAndOffboardingID(assetID uint, offboardingID string, tenantID *uint) (*models.OffboardingAssetReturn, error) {
	var assetReturn models.OffboardingAssetReturn
	query := r.db.Where("asset_id = ? AND offboarding_id = ?", assetID, offboardingID)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.First(&assetReturn).Error
	if err != nil {
		return nil, err
	}
	return &assetReturn, nil
}

// Update updates an asset return
func (r *OffboardingAssetReturnRepository) Update(assetReturn *models.OffboardingAssetReturn) error {
	return r.db.Save(assetReturn).Error
}

// FinalSettlementRepository handles database operations for final settlements
type FinalSettlementRepository struct {
	db *gorm.DB
}

// NewFinalSettlementRepository creates a new final settlement repository
func NewFinalSettlementRepository() *FinalSettlementRepository {
	return &FinalSettlementRepository{
		db: database.GetDB(),
	}
}

// Create creates a new final settlement
func (r *FinalSettlementRepository) Create(settlement *models.FinalSettlement) error {
	return r.db.Create(settlement).Error
}

// FindByID finds a settlement by ID
func (r *FinalSettlementRepository) FindByID(id uint) (*models.FinalSettlement, error) {
	var settlement models.FinalSettlement
	err := r.db.First(&settlement, id).Error
	if err != nil {
		return nil, err
	}
	return &settlement, nil
}

// FindBySettlementID finds a settlement by settlement_id (string)
func (r *FinalSettlementRepository) FindBySettlementID(settlementID string, tenantID *uint) (*models.FinalSettlement, error) {
	var settlement models.FinalSettlement
	query := r.db.Where("settlement_id = ?", settlementID)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.First(&settlement).Error
	if err != nil {
		return nil, err
	}
	return &settlement, nil
}

// FindByOffboardingID finds a settlement by offboarding ID
func (r *FinalSettlementRepository) FindByOffboardingID(offboardingID string, tenantID *uint) (*models.FinalSettlement, error) {
	var settlement models.FinalSettlement
	query := r.db.Where("offboarding_id = ?", offboardingID)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.First(&settlement).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Not found is not an error, return nil
	}
	if err != nil {
		return nil, err
	}
	return &settlement, nil
}

// Update updates a settlement
func (r *FinalSettlementRepository) Update(settlement *models.FinalSettlement) error {
	return r.db.Save(settlement).Error
}

// GenerateSettlementID generates a unique settlement ID
func (r *FinalSettlementRepository) GenerateSettlementID() (string, error) {
	var count int64
	r.db.Model(&models.FinalSettlement{}).Count(&count)
	return fmt.Sprintf("sett-%03d", count+1), nil
}
