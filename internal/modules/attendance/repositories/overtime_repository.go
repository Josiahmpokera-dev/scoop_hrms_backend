package repositories

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	"gorm.io/gorm"
)

// OvertimeRepository handles overtime database operations
type OvertimeRepository struct {
	db *gorm.DB
}

// NewOvertimeRepository creates a new overtime repository
func NewOvertimeRepository() *OvertimeRepository {
	return &OvertimeRepository{
		db: database.GetDB(),
	}
}

// --- OvertimePolicy ---

// CreatePolicy creates a new OT policy
func (r *OvertimeRepository) CreatePolicy(policy *models.OvertimePolicy) error {
	return r.db.Create(policy).Error
}

// FindPolicyByID finds a policy by ID
func (r *OvertimeRepository) FindPolicyByID(id uint) (*models.OvertimePolicy, error) {
	var policy models.OvertimePolicy
	err := r.db.Where("id = ?", id).First(&policy).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// UpdatePolicy updates a policy
func (r *OvertimeRepository) UpdatePolicy(policy *models.OvertimePolicy) error {
	return r.db.Save(policy).Error
}

// DeletePolicy soft deletes a policy
func (r *OvertimeRepository) DeletePolicy(id uint) error {
	return r.db.Delete(&models.OvertimePolicy{}, id).Error
}

// ListPolicies returns all active policies
func (r *OvertimeRepository) ListPolicies(page, pageSize int) ([]models.OvertimePolicy, int64, error) {
	var policies []models.OvertimePolicy
	var total int64

	query := r.db.Model(&models.OvertimePolicy{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&policies).Error
	return policies, total, err
}

// FindPolicyByDepartment finds the active policy for a department (or the default one)
func (r *OvertimeRepository) FindPolicyByDepartment(departmentID *uint) (*models.OvertimePolicy, error) {
	var policy models.OvertimePolicy
	// Try department-specific first
	if departmentID != nil {
		err := r.db.Where("department_id = ? AND is_active = ?", *departmentID, true).First(&policy).Error
		if err == nil {
			return &policy, nil
		}
	}
	// Fall back to default (no department)
	err := r.db.Where("department_id IS NULL AND is_active = ?", true).First(&policy).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// --- OvertimeRequest ---

// CreateRequest creates a new OT request
func (r *OvertimeRepository) CreateRequest(request *models.OvertimeRequest) error {
	return r.db.Create(request).Error
}

// FindRequestByID finds an OT request by ID
func (r *OvertimeRepository) FindRequestByID(id uint) (*models.OvertimeRequest, error) {
	var request models.OvertimeRequest
	err := r.db.Preload("Approvals").Where("id = ?", id).First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// UpdateRequest updates an OT request
func (r *OvertimeRepository) UpdateRequest(request *models.OvertimeRequest) error {
	return r.db.Save(request).Error
}

// ListRequestsByEmployee returns OT requests for an employee with pagination
func (r *OvertimeRepository) ListRequestsByEmployee(employeeID uint, status string, page, pageSize int) ([]models.OvertimeRequest, int64, error) {
	var requests []models.OvertimeRequest
	var total int64

	query := r.db.Model(&models.OvertimeRequest{}).Where("employee_id = ?", employeeID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Approvals").Order("date DESC").Offset(offset).Limit(pageSize).Find(&requests).Error
	return requests, total, err
}

// ListPendingRequests returns OT requests pending approval
func (r *OvertimeRepository) ListPendingRequests(page, pageSize int) ([]models.OvertimeRequest, int64, error) {
	var requests []models.OvertimeRequest
	var total int64

	query := r.db.Model(&models.OvertimeRequest{}).Where("status = ?", models.OTStatusPending)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Approvals").Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&requests).Error
	return requests, total, err
}

// ListAllRequests returns all OT requests with filters
func (r *OvertimeRepository) ListAllRequests(status string, startDate, endDate *time.Time, page, pageSize int) ([]models.OvertimeRequest, int64, error) {
	var requests []models.OvertimeRequest
	var total int64

	query := r.db.Model(&models.OvertimeRequest{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != nil {
		query = query.Where("date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("date <= ?", *endDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Approvals").Order("date DESC").Offset(offset).Limit(pageSize).Find(&requests).Error
	return requests, total, err
}

// --- Aggregation / Reports ---

// GetOvertimeStats returns OT stats for a date range
func (r *OvertimeRepository) GetOvertimeStats(startDate, endDate time.Time) (map[string]interface{}, error) {
	var result struct {
		TotalRequests  int64
		PendingCount   int64
		ApprovedHours  float64
		TotalPayout    float64
	}

	r.db.Model(&models.OvertimeRequest{}).
		Where("date >= ? AND date <= ?", startDate, endDate).
		Select("COUNT(*) as total_requests").
		Scan(&result)

	r.db.Model(&models.OvertimeRequest{}).
		Where("date >= ? AND date <= ? AND status = ?", startDate, endDate, models.OTStatusPending).
		Select("COUNT(*) as pending_count").
		Scan(&result)

	r.db.Model(&models.OvertimeRequest{}).
		Where("date >= ? AND date <= ? AND status = ?", startDate, endDate, models.OTStatusApproved).
		Select("COALESCE(SUM(hours), 0) as approved_hours").
		Scan(&result)

	r.db.Model(&models.OvertimeRequest{}).
		Where("date >= ? AND date <= ? AND status IN ?", startDate, endDate, []string{string(models.OTStatusApproved), string(models.OTStatusProcessed)}).
		Select("COALESCE(SUM(payout_amount), 0) as total_payout").
		Scan(&result)

	return map[string]interface{}{
		"total_requests": result.TotalRequests,
		"pending_count":  result.PendingCount,
		"approved_hours": result.ApprovedHours,
		"total_payout":   result.TotalPayout,
	}, nil
}

// GetEmployeeOTHoursInMonth returns total approved/pending OT hours for an employee in a month
func (r *OvertimeRepository) GetEmployeeOTHoursInMonth(employeeID uint, year int, month time.Month) (float64, error) {
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	var total float64
	err := r.db.Model(&models.OvertimeRequest{}).
		Where("employee_id = ? AND date >= ? AND date <= ? AND status IN ?",
			employeeID, startDate, endDate, []string{string(models.OTStatusPending), string(models.OTStatusApproved)}).
		Select("COALESCE(SUM(hours), 0)").
		Scan(&total).Error
	return total, err
}

// GetEmployeeOTHoursInWeek returns total approved/pending OT hours for an employee in a week
func (r *OvertimeRepository) GetEmployeeOTHoursInWeek(employeeID uint, weekStart, weekEnd time.Time) (float64, error) {
	var total float64
	err := r.db.Model(&models.OvertimeRequest{}).
		Where("employee_id = ? AND date >= ? AND date <= ? AND status IN ?",
			employeeID, weekStart, weekEnd, []string{string(models.OTStatusPending), string(models.OTStatusApproved)}).
		Select("COALESCE(SUM(hours), 0)").
		Scan(&total).Error
	return total, err
}

// --- OvertimeApproval ---

// CreateApproval creates a new OT approval record
func (r *OvertimeRepository) CreateApproval(approval *models.OvertimeApproval) error {
	return r.db.Create(approval).Error
}
