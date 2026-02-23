package repositories

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/models"
	"gorm.io/gorm"
)

type EnrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository() *EnrollmentRepository {
	return &EnrollmentRepository{db: database.GetDB()}
}

func (r *EnrollmentRepository) Create(e *models.BiometricEnrollment) error {
	return r.db.Create(e).Error
}

func (r *EnrollmentRepository) FindByEmployeeID(employeeID uint) (*models.BiometricEnrollment, error) {
	var e models.BiometricEnrollment
	if err := r.db.Where("employee_id = ? AND status = 'linked'", employeeID).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *EnrollmentRepository) FindByEmpCode(empCode string) (*models.BiometricEnrollment, error) {
	var e models.BiometricEnrollment
	if err := r.db.Where("emp_code = ? AND status = 'linked'", empCode).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *EnrollmentRepository) Unlink(employeeID uint) error {
	return r.db.Model(&models.BiometricEnrollment{}).
		Where("employee_id = ? AND status = 'linked'", employeeID).
		Updates(map[string]interface{}{
			"status":     "unlinked",
			"updated_at": time.Now(),
		}).Error
}

func (r *EnrollmentRepository) ListLinked(page, pageSize int) ([]models.BiometricEnrollment, int64, error) {
	var items []models.BiometricEnrollment
	var total int64

	r.db.Model(&models.BiometricEnrollment{}).Where("status = 'linked'").Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Where("status = 'linked'").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&items).Error
	return items, total, err
}

// GetLinkedEmpCodes returns a set of all currently linked emp_codes.
func (r *EnrollmentRepository) GetLinkedEmpCodes() (map[string]uint, error) {
	var enrollments []models.BiometricEnrollment
	if err := r.db.Where("status = 'linked'").Find(&enrollments).Error; err != nil {
		return nil, err
	}
	m := make(map[string]uint, len(enrollments))
	for _, e := range enrollments {
		m[e.EmpCode] = e.EmployeeID
	}
	return m, nil
}

// GetLinkedEmployeeIDs returns a set of all currently linked employee IDs.
func (r *EnrollmentRepository) GetLinkedEmployeeIDs() (map[uint]string, error) {
	var enrollments []models.BiometricEnrollment
	if err := r.db.Where("status = 'linked'").Find(&enrollments).Error; err != nil {
		return nil, err
	}
	m := make(map[uint]string, len(enrollments))
	for _, e := range enrollments {
		m[e.EmployeeID] = e.EmpCode
	}
	return m, nil
}

// GetDistinctDeviceUsers returns distinct emp_code/name from biometric transactions.
func (r *EnrollmentRepository) GetDistinctDeviceUsers() ([]models.DeviceUser, error) {
	var results []struct {
		EmpCode    string `json:"emp_code"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Department string `json:"department"`
		Position   string `json:"position"`
	}

	err := r.db.Raw(`
		SELECT DISTINCT ON (emp_code)
			emp_code, first_name, last_name, department, position
		FROM biotime_transactions
		WHERE emp_code IS NOT NULL AND emp_code != ''
		ORDER BY emp_code, punch_time DESC
	`).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	linkedMap, _ := r.GetLinkedEmpCodes()

	var users []models.DeviceUser
	for _, row := range results {
		u := models.DeviceUser{
			EmpCode:    row.EmpCode,
			FirstName:  row.FirstName,
			LastName:   row.LastName,
			Department: row.Department,
			Position:   row.Position,
		}
		if empID, ok := linkedMap[row.EmpCode]; ok {
			u.Linked = true
			u.EmployeeID = &empID
		}
		users = append(users, u)
	}
	return users, nil
}

// GetMergedDailyAttendance joins biometric daily attendance with employee data
// for all linked employees in a date range.
func (r *EnrollmentRepository) GetMergedDailyAttendance(startDate, endDate string, employeeID *uint, page, pageSize int) ([]models.MergedAttendance, int64, error) {
	var results []models.MergedAttendance
	var total int64

	baseQuery := `
		FROM biometric_enrollments be
		JOIN employees e ON e.id = be.employee_id AND e.deleted_at IS NULL
		JOIN (
			SELECT
				emp_code,
				DATE(punch_time) AS punch_date,
				MIN(punch_time) AS check_in,
				MAX(punch_time) AS check_out,
				COUNT(*) AS punch_count
			FROM biotime_transactions
			WHERE DATE(punch_time) >= ? AND DATE(punch_time) <= ?
			GROUP BY emp_code, DATE(punch_time)
		) att ON att.emp_code = be.emp_code
		LEFT JOIN departments d ON d.id = e.department_id AND d.deleted_at IS NULL
		LEFT JOIN job_positions p ON p.id = e.position_id AND p.deleted_at IS NULL
		WHERE be.status = 'linked' AND be.deleted_at IS NULL
	`
	args := []interface{}{startDate, endDate}

	if employeeID != nil {
		baseQuery += " AND be.employee_id = ?"
		args = append(args, *employeeID)
	}

	countSQL := "SELECT COUNT(*) " + baseQuery
	r.db.Raw(countSQL, args...).Scan(&total)

	selectSQL := `
		SELECT
			e.id AS employee_id,
			e.employee_id AS employee_code,
			CONCAT(e.first_name, ' ', e.last_name) AS employee_name,
			d.name AS department,
			p.title AS position,
			e.photo_url,
			be.emp_code,
			TO_CHAR(att.punch_date, 'YYYY-MM-DD') AS date,
			TO_CHAR(att.check_in, 'HH24:MI:SS') AS check_in,
			TO_CHAR(att.check_out, 'HH24:MI:SS') AS check_out,
			att.punch_count
	` + baseQuery + " ORDER BY att.punch_date DESC, e.first_name ASC"

	offset := (page - 1) * pageSize
	selectSQL += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	err := r.db.Raw(selectSQL, args...).Scan(&results).Error
	return results, total, err
}
