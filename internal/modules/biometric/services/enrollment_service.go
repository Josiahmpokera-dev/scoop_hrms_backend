package services

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/repositories"
	employeeModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	"gorm.io/gorm"
)

type EnrollmentService struct {
	repo *repositories.EnrollmentRepository
}

func NewEnrollmentService() *EnrollmentService {
	return &EnrollmentService{repo: repositories.NewEnrollmentRepository()}
}

// ---------- Link / Unlink ----------

func (s *EnrollmentService) LinkEmployee(req *models.LinkEmployeeRequest, callerID uint) (*models.BiometricEnrollment, error) {
	db := database.GetDB()

	var emp employeeModels.Employee
	if err := db.First(&emp, req.EmployeeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("employee with ID %d not found", req.EmployeeID)
		}
		return nil, err
	}

	if existing, _ := s.repo.FindByEmployeeID(req.EmployeeID); existing != nil {
		return nil, fmt.Errorf("employee %d is already linked to biometric code '%s'", req.EmployeeID, existing.EmpCode)
	}
	if existing, _ := s.repo.FindByEmpCode(req.EmpCode); existing != nil {
		return nil, fmt.Errorf("biometric code '%s' is already linked to employee %d", req.EmpCode, existing.EmployeeID)
	}

	var deviceName string
	db.Raw(`
		SELECT CONCAT(first_name, ' ', last_name) FROM biotime_transactions
		WHERE emp_code = ? LIMIT 1
	`, req.EmpCode).Scan(&deviceName)

	enrollment := &models.BiometricEnrollment{
		EmployeeID:     req.EmployeeID,
		EmpCode:        req.EmpCode,
		DeviceUserName: deviceName,
		Status:         "linked",
		LinkedAt:       time.Now(),
		LinkedBy:       &callerID,
	}
	if req.Notes != "" {
		enrollment.Notes = &req.Notes
	}

	if err := s.repo.Create(enrollment); err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return nil, fmt.Errorf("this employee or biometric code is already enrolled")
		}
		return nil, err
	}
	return enrollment, nil
}

func (s *EnrollmentService) BulkLink(req *models.BulkLinkRequest, callerID uint) ([]models.BiometricEnrollment, []string, error) {
	var linked []models.BiometricEnrollment
	var errs []string

	for _, m := range req.Mappings {
		e, err := s.LinkEmployee(&m, callerID)
		if err != nil {
			errs = append(errs, fmt.Sprintf("employee_id=%d emp_code=%s: %s", m.EmployeeID, m.EmpCode, err.Error()))
			continue
		}
		linked = append(linked, *e)
	}
	return linked, errs, nil
}

func (s *EnrollmentService) UnlinkEmployee(employeeID uint) error {
	if existing, _ := s.repo.FindByEmployeeID(employeeID); existing == nil {
		return fmt.Errorf("employee %d is not linked to any biometric device user", employeeID)
	}
	return s.repo.Unlink(employeeID)
}

// ---------- Auto-Link ----------

func (s *EnrollmentService) AutoLink(callerID uint) (*models.AutoLinkResult, error) {
	db := database.GetDB()

	var employees []employeeModels.Employee
	if err := db.Where("employee_id IS NOT NULL AND employee_id != '' AND status = 'active'").Find(&employees).Error; err != nil {
		return nil, err
	}

	deviceUsers, err := s.repo.GetDistinctDeviceUsers()
	if err != nil {
		return nil, err
	}
	deviceMap := make(map[string]models.DeviceUser, len(deviceUsers))
	for _, du := range deviceUsers {
		deviceMap[strings.TrimSpace(du.EmpCode)] = du
	}

	linkedEmpIDs, _ := s.repo.GetLinkedEmployeeIDs()
	linkedCodes, _ := s.repo.GetLinkedEmpCodes()

	result := &models.AutoLinkResult{
		TotalEmployees:   len(employees),
		TotalDeviceUsers: len(deviceUsers),
	}

	for _, emp := range employees {
		detail := models.AutoLinkDetail{
			EmployeeID:   emp.ID,
			EmployeeCode: emp.EmployeeID,
			EmployeeName: emp.FirstName + " " + emp.LastName,
		}

		if _, ok := linkedEmpIDs[emp.ID]; ok {
			detail.Action = "already_linked"
			detail.EmpCode = linkedEmpIDs[emp.ID]
			result.AlreadyLinked++
			result.Details = append(result.Details, detail)
			continue
		}

		code := strings.TrimSpace(emp.EmployeeID)
		du, found := deviceMap[code]
		if !found {
			detail.Action = "no_match"
			result.NoMatch++
			result.Details = append(result.Details, detail)
			continue
		}

		if _, taken := linkedCodes[code]; taken {
			detail.Action = "already_linked"
			detail.EmpCode = code
			result.AlreadyLinked++
			result.Details = append(result.Details, detail)
			continue
		}

		enrollment := &models.BiometricEnrollment{
			EmployeeID:     emp.ID,
			EmpCode:        du.EmpCode,
			DeviceUserName: du.FirstName + " " + du.LastName,
			Status:         "linked",
			LinkedAt:       time.Now(),
			LinkedBy:       &callerID,
		}
		if err := s.repo.Create(enrollment); err == nil {
			detail.Action = "linked"
			detail.EmpCode = du.EmpCode
			result.Linked++
			linkedCodes[du.EmpCode] = emp.ID
			linkedEmpIDs[emp.ID] = du.EmpCode
		} else {
			detail.Action = "no_match"
			result.NoMatch++
		}
		result.Details = append(result.Details, detail)
	}

	return result, nil
}

// ---------- Listing ----------

func (s *EnrollmentService) ListEnrollments(page, pageSize int) ([]models.EnrollmentDetail, int64, error) {
	enrollments, total, err := s.repo.ListLinked(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	db := database.GetDB()
	var details []models.EnrollmentDetail

	for _, e := range enrollments {
		var emp employeeModels.Employee
		if err := db.First(&emp, e.EmployeeID).Error; err != nil {
			continue
		}

		var deptName, posTitle *string
		if emp.DepartmentID != nil {
			var name string
			db.Raw("SELECT name FROM departments WHERE id = ?", *emp.DepartmentID).Scan(&name)
			if name != "" {
				deptName = &name
			}
		}
		if emp.PositionID != nil {
			var title string
			db.Raw("SELECT title FROM job_positions WHERE id = ?", *emp.PositionID).Scan(&title)
			if title != "" {
				posTitle = &title
			}
		}

		details = append(details, models.EnrollmentDetail{
			ID:             e.ID,
			EmployeeID:     e.EmployeeID,
			EmployeeCode:   emp.EmployeeID,
			EmployeeName:   emp.FirstName + " " + emp.LastName,
			Department:     deptName,
			Position:       posTitle,
			EmpCode:        e.EmpCode,
			DeviceUserName: e.DeviceUserName,
			Status:         e.Status,
			LinkedAt:       e.LinkedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return details, total, nil
}

func (s *EnrollmentService) GetUnlinkedEmployees() ([]models.UnlinkedEmployee, error) {
	db := database.GetDB()

	linkedIDs, err := s.repo.GetLinkedEmployeeIDs()
	if err != nil {
		return nil, err
	}

	var employees []employeeModels.Employee
	query := db.Where("status = 'active'")
	if len(linkedIDs) > 0 {
		ids := make([]uint, 0, len(linkedIDs))
		for id := range linkedIDs {
			ids = append(ids, id)
		}
		query = query.Where("id NOT IN ?", ids)
	}
	if err := query.Find(&employees).Error; err != nil {
		return nil, err
	}

	deviceUsers, _ := s.repo.GetDistinctDeviceUsers()
	deviceMap := make(map[string]bool, len(deviceUsers))
	for _, du := range deviceUsers {
		deviceMap[strings.TrimSpace(du.EmpCode)] = true
	}

	var result []models.UnlinkedEmployee
	for _, emp := range employees {
		u := models.UnlinkedEmployee{
			ID:           emp.ID,
			EmployeeCode: emp.EmployeeID,
			FirstName:    emp.FirstName,
			LastName:     emp.LastName,
		}
		if emp.DepartmentID != nil {
			var name string
			db.Raw("SELECT name FROM departments WHERE id = ?", *emp.DepartmentID).Scan(&name)
			if name != "" {
				u.Department = &name
			}
		}
		if emp.PositionID != nil {
			var title string
			db.Raw("SELECT title FROM job_positions WHERE id = ?", *emp.PositionID).Scan(&title)
			if title != "" {
				u.Position = &title
			}
		}
		code := strings.TrimSpace(emp.EmployeeID)
		if deviceMap[code] {
			u.SuggestedEmpCode = &code
		}

		result = append(result, u)
	}
	return result, nil
}

func (s *EnrollmentService) GetDeviceUsers() ([]models.DeviceUser, error) {
	return s.repo.GetDistinctDeviceUsers()
}

// ---------- Merged Attendance ----------

func (s *EnrollmentService) GetMergedAttendance(startDate, endDate string, employeeID *uint, page, pageSize int) ([]models.MergedAttendance, int64, int, error) {
	records, total, err := s.repo.GetMergedDailyAttendance(startDate, endDate, employeeID, page, pageSize)
	if err != nil {
		return nil, 0, 0, err
	}

	for i := range records {
		if records[i].CheckIn != nil && records[i].CheckOut != nil {
			cin, _ := time.Parse("15:04:05", *records[i].CheckIn)
			cout, _ := time.Parse("15:04:05", *records[i].CheckOut)
			if !cin.IsZero() && !cout.IsZero() {
				diff := cout.Sub(cin)
				hours := math.Floor(diff.Hours())
				mins := math.Floor(diff.Minutes()) - hours*60
				wh := fmt.Sprintf("%dh %dm", int(hours), int(mins))
				records[i].WorkingHours = &wh
			}
		}

		records[i].Status = s.computeStatus(records[i])
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	return records, total, totalPages, nil
}

func (s *EnrollmentService) computeStatus(r models.MergedAttendance) string {
	if r.CheckIn == nil && r.CheckOut == nil {
		return "Absent"
	}

	if r.CheckIn != nil {
		cin, _ := time.Parse("15:04:05", *r.CheckIn)
		threshold, _ := time.Parse("15:04:05", "08:30:00")
		if cin.After(threshold) {
			return "Late"
		}
	}

	if r.CheckIn != nil && r.CheckOut != nil {
		cin, _ := time.Parse("15:04:05", *r.CheckIn)
		cout, _ := time.Parse("15:04:05", *r.CheckOut)
		if cout.Sub(cin).Hours() < 4.5 {
			return "Half Day"
		}
	}

	return "Present"
}

// ---------- Enrollment Summary / Statistics ----------

func (s *EnrollmentService) GetStatistics() (map[string]interface{}, error) {
	db := database.GetDB()

	var totalEmployees int64
	db.Model(&employeeModels.Employee{}).Where("status = 'active'").Count(&totalEmployees)

	var totalLinked int64
	db.Model(&models.BiometricEnrollment{}).Where("status = 'linked'").Count(&totalLinked)

	deviceUsers, _ := s.repo.GetDistinctDeviceUsers()
	totalDevice := int64(len(deviceUsers))

	var linkedDeviceUsers int64
	for _, du := range deviceUsers {
		if du.Linked {
			linkedDeviceUsers++
		}
	}

	return map[string]interface{}{
		"total_employees":        totalEmployees,
		"linked_employees":       totalLinked,
		"unlinked_employees":     totalEmployees - totalLinked,
		"total_device_users":     totalDevice,
		"linked_device_users":    linkedDeviceUsers,
		"unlinked_device_users":  totalDevice - linkedDeviceUsers,
		"coverage_percentage":    math.Round(float64(totalLinked)/float64(max(totalEmployees, 1))*10000) / 100,
	}, nil
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
