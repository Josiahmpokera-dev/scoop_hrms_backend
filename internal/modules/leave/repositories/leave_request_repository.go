package repositories

import (
	"fmt"
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"gorm.io/gorm"
)

type LeaveRequestRepository struct {
	db *gorm.DB
}

func NewLeaveRequestRepository() *LeaveRequestRepository {
	return &LeaveRequestRepository{
		db: database.GetDB(),
	}
}

// Create creates a new leave request
func (r *LeaveRequestRepository) Create(request *models.LeaveRequest) error {
	return r.db.Create(request).Error
}

// FindByID finds a leave request by ID
func (r *LeaveRequestRepository) FindByID(id uint) (*models.LeaveRequest, error) {
	var request models.LeaveRequest
	err := r.db.Preload("LeaveType").
		Preload("Documents").
		Preload("Approvals").
		First(&request, id).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// FindByApplicationNumber finds a leave request by application number
func (r *LeaveRequestRepository) FindByApplicationNumber(appNumber string) (*models.LeaveRequest, error) {
	var request models.LeaveRequest
	err := r.db.Preload("LeaveType").
		Preload("Documents").
		Preload("Approvals").
		Where("application_number = ?", appNumber).
		First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// Update updates a leave request
func (r *LeaveRequestRepository) Update(request *models.LeaveRequest) error {
	return r.db.Save(request).Error
}

// Delete soft deletes a leave request
func (r *LeaveRequestRepository) Delete(id uint) error {
	return r.db.Delete(&models.LeaveRequest{}, id).Error
}

// List returns leave requests with pagination and filters
func (r *LeaveRequestRepository) List(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.LeaveRequest, int64, error) {
	var requests []models.LeaveRequest
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.LeaveRequest{}).
		Preload("LeaveType").
		Preload("Documents")

	// Apply tenant filter
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Apply filters
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if employeeID, ok := filters["employee_id"].(string); ok && employeeID != "" {
		query = query.Where("employee_id = ?", employeeID)
	}

	if leaveTypeCode, ok := filters["leave_type_code"].(string); ok && leaveTypeCode != "" {
		query = query.Where("leave_type_code = ?", leaveTypeCode)
	}

	if fromDate, ok := filters["from_date"].(string); ok && fromDate != "" {
		query = query.Where("from_date >= ?", fromDate)
	}

	if toDate, ok := filters["to_date"].(string); ok && toDate != "" {
		query = query.Where("to_date <= ?", toDate)
	}

	if search, ok := filters["search"].(string); ok && strings.TrimSpace(search) != "" {
		term := "%" + strings.ToLower(strings.TrimSpace(search)) + "%"
		query = query.Joins("LEFT JOIN employees e ON e.employee_id = leave_requests.employee_id").
			Where(`(
				LOWER(leave_requests.application_number) LIKE ? OR
				LOWER(leave_requests.employee_id) LIKE ? OR
				LOWER(leave_requests.leave_type_code) LIKE ? OR
				LOWER(leave_requests.reason) LIKE ? OR
				LOWER(e.first_name) LIKE ? OR
				LOWER(e.last_name) LIKE ?
			)`, term, term, term, term, term, term)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("application_date DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&requests).Error
	return requests, total, err
}

// GenerateApplicationNumber generates a unique application number
func (r *LeaveRequestRepository) GenerateApplicationNumber(year int) (string, error) {
	var count int64
	prefix := "LV"
	yearStr := fmt.Sprintf("%d", year)

	// Get count of requests for this year
	pattern := fmt.Sprintf("%s-%s-%%", prefix, yearStr)
	r.db.Model(&models.LeaveRequest{}).
		Where("application_number LIKE ?", pattern).
		Count(&count)

	// Format: LV-2026-00123
	sequence := count + 1
	appNumber := fmt.Sprintf("%s-%s-%05d", prefix, yearStr, sequence)

	return appNumber, nil
}

// FindByDateRange finds leave requests within a date range
func (r *LeaveRequestRepository) FindByDateRange(startDate, endDate string, tenantID *uint) ([]models.LeaveRequest, error) {
	var requests []models.LeaveRequest
	query := r.db.Preload("LeaveType").
		Where("(from_date <= ? AND to_date >= ?) OR (from_date BETWEEN ? AND ?) OR (to_date BETWEEN ? AND ?)",
			endDate, startDate, startDate, endDate, startDate, endDate)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	err := query.Find(&requests).Error
	return requests, err
}

// SumConsumedDaysForYear sums approved / partially approved working days for balance used totals.
func (r *LeaveRequestRepository) SumConsumedDaysForYear(employeeID, leaveTypeCode string, year int, tenantID *uint) (float64, error) {
	var sum float64
	query := r.db.Model(&models.LeaveRequest{}).
		Select(`COALESCE(SUM(
			CASE
				WHEN status = 'partially_approved' AND approved_days IS NOT NULL THEN approved_days
				WHEN status IN ('approved', 'partially_approved') THEN total_days
				ELSE 0
			END
		), 0)`).
		Where("employee_id = ? AND leave_type_code = ? AND leave_period_year = ?", employeeID, leaveTypeCode, year)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Scan(&sum).Error; err != nil {
		return 0, err
	}
	return sum, nil
}

// SumPendingDaysForYear sums days still awaiting approval (pending / returned for info).
func (r *LeaveRequestRepository) SumPendingDaysForYear(employeeID, leaveTypeCode string, year int, tenantID *uint) (float64, error) {
	var sum float64
	query := r.db.Model(&models.LeaveRequest{}).
		Select("COALESCE(SUM(total_days), 0)").
		Where("employee_id = ? AND leave_type_code = ? AND leave_period_year = ?", employeeID, leaveTypeCode, year).
		Where("status IN ?", []string{"pending", "returned_for_info"})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Scan(&sum).Error; err != nil {
		return 0, err
	}
	return sum, nil
}
