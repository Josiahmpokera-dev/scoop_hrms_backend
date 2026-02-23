package services

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	deptModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/models"
	empModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	posModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/models"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ImportResult struct {
	TotalRows int            `json:"total_rows"`
	Imported  int            `json:"imported"`
	Skipped   int            `json:"skipped"`
	Errors    []RowError     `json:"errors,omitempty"`
}

type RowError struct {
	Row     int    `json:"row"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// ── Import Departments ─────────────────────────────────────────────────

func ImportDepartments(file *multipart.FileHeader) (*ImportResult, error) {
	f, sheet, err := openExcel(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("no data rows found (only header)")
	}

	db := database.GetDB()
	result := &ImportResult{TotalRows: len(rows) - 1}

	codeMap := map[string]uint{}
	var existing []deptModels.Department
	db.Select("id, code").Find(&existing)
	for _, d := range existing {
		codeMap[strings.ToUpper(d.Code)] = d.ID
	}

	for i, row := range rows[1:] {
		rowNum := i + 2
		code := cellVal(row, 0)
		name := cellVal(row, 1)

		if code == "" || name == "" {
			result.Errors = append(result.Errors, RowError{Row: rowNum, Message: "Code and Name are required"})
			result.Skipped++
			continue
		}

		if _, exists := codeMap[strings.ToUpper(code)]; exists {
			result.Errors = append(result.Errors, RowError{Row: rowNum, Field: "Code", Message: fmt.Sprintf("Department code '%s' already exists", code)})
			result.Skipped++
			continue
		}

		dept := deptModels.Department{
			Code:     code,
			Name:     name,
			IsActive: true,
		}

		if v := cellVal(row, 2); v != "" {
			dept.Description = &v
		}
		if v := cellVal(row, 3); v != "" {
			dept.Level = &v
		}
		if v := cellVal(row, 4); v != "" {
			dept.DepartmentType = &v
		}
		if parentCode := cellVal(row, 5); parentCode != "" {
			if pid, ok := codeMap[strings.ToUpper(parentCode)]; ok {
				dept.ParentDepartmentID = &pid
			} else {
				result.Errors = append(result.Errors, RowError{Row: rowNum, Field: "Parent Department Code", Message: fmt.Sprintf("Parent '%s' not found", parentCode)})
				result.Skipped++
				continue
			}
		}
		if v := cellVal(row, 6); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				dept.EmployeeCapacity = &n
			}
		}

		if err := db.Create(&dept).Error; err != nil {
			result.Errors = append(result.Errors, RowError{Row: rowNum, Message: err.Error()})
			result.Skipped++
			continue
		}

		codeMap[strings.ToUpper(code)] = dept.ID
		result.Imported++
	}

	return result, nil
}

// ── Import Positions ───────────────────────────────────────────────────

func ImportPositions(file *multipart.FileHeader) (*ImportResult, error) {
	f, sheet, err := openExcel(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("no data rows found (only header)")
	}

	db := database.GetDB()
	result := &ImportResult{TotalRows: len(rows) - 1}

	deptMap := buildCodeMap[deptModels.Department](db, "departments")
	posMap := buildCodeMap[posModels.JobPosition](db, "job_positions")

	for i, row := range rows[1:] {
		rowNum := i + 2
		code := cellVal(row, 0)
		title := cellVal(row, 1)

		if code == "" || title == "" {
			result.Errors = append(result.Errors, RowError{Row: rowNum, Message: "Code and Title are required"})
			result.Skipped++
			continue
		}

		if _, exists := posMap[strings.ToUpper(code)]; exists {
			result.Errors = append(result.Errors, RowError{Row: rowNum, Field: "Code", Message: fmt.Sprintf("Position code '%s' already exists", code)})
			result.Skipped++
			continue
		}

		pos := posModels.JobPosition{
			Code:     code,
			Title:    title,
			IsActive: true,
		}

		if v := cellVal(row, 2); v != "" {
			pos.Grade = &v
		}
		if deptCode := cellVal(row, 3); deptCode != "" {
			if did, ok := deptMap[strings.ToUpper(deptCode)]; ok {
				pos.DepartmentID = &did
			} else {
				result.Errors = append(result.Errors, RowError{Row: rowNum, Field: "Department Code", Message: fmt.Sprintf("Department '%s' not found", deptCode)})
				result.Skipped++
				continue
			}
		}
		if reportsCode := cellVal(row, 4); reportsCode != "" {
			if rid, ok := posMap[strings.ToUpper(reportsCode)]; ok {
				pos.ReportsToPositionID = &rid
			}
		}
		if v := cellVal(row, 5); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				pos.BudgetedHeadcount = &n
			}
		}
		if v := cellVal(row, 6); v != "" {
			pos.KeyCompetencies = &v
		}

		if err := db.Create(&pos).Error; err != nil {
			result.Errors = append(result.Errors, RowError{Row: rowNum, Message: err.Error()})
			result.Skipped++
			continue
		}

		posMap[strings.ToUpper(code)] = pos.ID
		result.Imported++
	}

	return result, nil
}

// ── Import Employees ───────────────────────────────────────────────────

func ImportEmployees(file *multipart.FileHeader) (*ImportResult, error) {
	f, sheet, err := openExcel(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("no data rows found (only header)")
	}

	db := database.GetDB()
	result := &ImportResult{TotalRows: len(rows) - 1}

	deptMap := buildCodeMap[deptModels.Department](db, "departments")
	posMap := buildCodeMap[posModels.JobPosition](db, "job_positions")

	empIDMap := map[string]uint{}
	var existingEmps []empModels.Employee
	db.Select("id, employee_id").Find(&existingEmps)
	for _, e := range existingEmps {
		empIDMap[strings.ToUpper(e.EmployeeID)] = e.ID
	}

	emailMap := map[string]bool{}
	var existingEmails []struct{ WorkEmail *string }
	db.Model(&empModels.Employee{}).Select("work_email").Where("work_email IS NOT NULL").Find(&existingEmails)
	for _, e := range existingEmails {
		if e.WorkEmail != nil {
			emailMap[strings.ToLower(*e.WorkEmail)] = true
		}
	}

	for i, row := range rows[1:] {
		rowNum := i + 2
		employeeID := cellVal(row, 0)
		firstName := cellVal(row, 1)
		lastName := cellVal(row, 3)

		if employeeID == "" || firstName == "" || lastName == "" {
			result.Errors = append(result.Errors, RowError{Row: rowNum, Message: "Employee ID, First Name, and Last Name are required"})
			result.Skipped++
			continue
		}

		if _, exists := empIDMap[strings.ToUpper(employeeID)]; exists {
			result.Errors = append(result.Errors, RowError{Row: rowNum, Field: "Employee ID", Message: fmt.Sprintf("Employee ID '%s' already exists", employeeID)})
			result.Skipped++
			continue
		}

		emp := empModels.Employee{
			EmployeeID: employeeID,
			FirstName:  firstName,
			LastName:   lastName,
			Status:     empModels.StatusActive,
			IsActive:   true,
			Currency:   "TZS",
		}

		if v := cellVal(row, 2); v != "" {
			emp.MiddleName = &v
		}
		if v := cellVal(row, 4); v != "" {
			lower := strings.ToLower(v)
			if emailMap[lower] {
				result.Errors = append(result.Errors, RowError{Row: rowNum, Field: "Work Email", Message: fmt.Sprintf("Email '%s' already exists", v)})
				result.Skipped++
				continue
			}
			emp.WorkEmail = &v
			emailMap[lower] = true
		}
		if v := cellVal(row, 5); v != "" {
			emp.PersonalEmail = &v
		}
		if v := cellVal(row, 6); v != "" {
			emp.PhoneNumber = &v
		}
		if v := cellVal(row, 7); v != "" {
			emp.Gender = &v
		}
		if v := cellVal(row, 8); v != "" {
			if t, err := parseDate(v); err == nil {
				emp.DateOfBirth = &t
			}
		}
		if v := cellVal(row, 9); v != "" {
			emp.Nationality = &v
		}
		if v := cellVal(row, 10); v != "" {
			emp.MaritalStatus = &v
		}
		if deptCode := cellVal(row, 11); deptCode != "" {
			if did, ok := deptMap[strings.ToUpper(deptCode)]; ok {
				emp.DepartmentID = &did
			} else {
				result.Errors = append(result.Errors, RowError{Row: rowNum, Field: "Department Code", Message: fmt.Sprintf("Department '%s' not found", deptCode)})
				result.Skipped++
				continue
			}
		}
		if posCode := cellVal(row, 12); posCode != "" {
			if pid, ok := posMap[strings.ToUpper(posCode)]; ok {
				emp.PositionID = &pid
			} else {
				result.Errors = append(result.Errors, RowError{Row: rowNum, Field: "Position Code", Message: fmt.Sprintf("Position '%s' not found", posCode)})
				result.Skipped++
				continue
			}
		}
		if v := cellVal(row, 13); v != "" {
			emp.EmploymentType = &v
		}
		if v := cellVal(row, 14); v != "" {
			if t, err := parseDate(v); err == nil {
				emp.HireDate = &t
			}
		}
		if v := cellVal(row, 15); v != "" {
			emp.Grade = &v
		}
		if v := cellVal(row, 16); v != "" {
			emp.Shift = &v
		}
		if mgrID := cellVal(row, 17); mgrID != "" {
			if mid, ok := empIDMap[strings.ToUpper(mgrID)]; ok {
				emp.ManagerID = &mid
				emp.ReportsToID = &mid
			}
		}
		if v := cellVal(row, 18); v != "" {
			if sal, err := strconv.ParseFloat(v, 64); err == nil {
				emp.Salary = &sal
			}
		}
		if v := cellVal(row, 19); v != "" {
			emp.Currency = v
		}
		if v := cellVal(row, 20); v != "" {
			emp.EmergencyContactName = v
		}
		if v := cellVal(row, 21); v != "" {
			emp.EmergencyContactPhone = v
		}
		if v := cellVal(row, 22); v != "" {
			emp.EmergencyContactRelation = v
		}
		if v := cellVal(row, 23); v != "" {
			emp.Notes = v
		}

		if err := db.Create(&emp).Error; err != nil {
			result.Errors = append(result.Errors, RowError{Row: rowNum, Message: err.Error()})
			result.Skipped++
			continue
		}

		empIDMap[strings.ToUpper(employeeID)] = emp.ID
		result.Imported++
	}

	return result, nil
}

// ── Helpers ────────────────────────────────────────────────────────────

func openExcel(file *multipart.FileHeader) (*excelize.File, string, error) {
	src, err := file.Open()
	if err != nil {
		return nil, "", fmt.Errorf("failed to open file: %w", err)
	}

	f, err := excelize.OpenReader(src)
	if err != nil {
		src.Close()
		return nil, "", fmt.Errorf("invalid Excel file: %w", err)
	}

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		f.Close()
		return nil, "", fmt.Errorf("Excel file has no sheets")
	}

	// Use first sheet (skip "Instructions" if it's first)
	sheet := sheets[0]
	if strings.EqualFold(sheet, "Instructions") && len(sheets) > 1 {
		sheet = sheets[1]
	}

	return f, sheet, nil
}

func cellVal(row []string, idx int) string {
	if idx < len(row) {
		return strings.TrimSpace(row[idx])
	}
	return ""
}

func parseDate(s string) (time.Time, error) {
	formats := []string{"2006-01-02", "01/02/2006", "02/01/2006", "2006/01/02", "Jan 2, 2006"}
	for _, fmt := range formats {
		if t, err := time.Parse(fmt, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse date: %s", s)
}

type codeEntity interface {
	deptModels.Department | posModels.JobPosition
}

func buildCodeMap[T codeEntity](db *gorm.DB, table string) map[string]uint {
	var results []struct {
		ID   uint
		Code string
	}
	db.Table(table).Select("id, code").Where("deleted_at IS NULL").Find(&results)
	m := make(map[string]uint, len(results))
	for _, r := range results {
		m[strings.ToUpper(r.Code)] = r.ID
	}
	return m
}
