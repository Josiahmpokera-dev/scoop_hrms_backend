package services

import (
	employeeservices "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/repositories"
)

// PrePayrollService handles business logic for pre-payroll employee listing and confirmation
// This service is used BEFORE creating any payroll run to review employee salary configurations
// and handle HR attendance confirmation and Finance payroll confirmation

type PrePayrollService struct {
	employeeService     *employeeservices.EmployeeService
	salaryStructureRepo *repositories.SalaryStructureRepository
}

// NewPrePayrollService creates a new pre-payroll service instance
func NewPrePayrollService() *PrePayrollService {
	return &PrePayrollService{
		employeeService:     employeeservices.NewEmployeeService(),
		salaryStructureRepo: repositories.NewSalaryStructureRepository(),
	}
}

// EmployeeSalaryInfo represents employee salary information for pre-payroll review
type EmployeeSalaryInfo struct {
	EmployeeID       uint   `json:"employeeId"`
	EmployeeCode     string `json:"employeeCode"`
	EmployeeName     string `json:"employeeName"`
	Department       string `json:"department"`
	Designation      string `json:"designation"`
	EmploymentStatus string `json:"employmentStatus"`

	// Basic Salary Information
	BasicSalary         float64 `json:"basicSalary"`
	GrossSalary         float64 `json:"grossSalary"`
	SalaryCurrency      string  `json:"salaryCurrency"`
	SalaryGrade         string  `json:"salaryGrade"`
	SalaryEffectiveDate string  `json:"salaryEffectiveDate"`

	// Bank Information for Payment
	BankName    string `json:"bankName"`
	BankAccount string `json:"bankAccount"`

	// Statutory Information
	TINNumber    string `json:"tinNumber"`
	NSSFNumber   string `json:"nssfNumber"`
	HasHESLBLoan bool   `json:"hasHeslbLoan"`

	// Status Tracking
	HasSalaryConfigured bool `json:"hasSalaryConfigured"`
	HasBankInfo         bool `json:"hasBankInfo"`
	HasStatutoryInfo    bool `json:"hasStatutoryInfo"`
	ReadyForPayroll     bool `json:"readyForPayroll"`

	// HR Attendance Status (To be confirmed by HR)
	AttendanceStatus      string  `json:"attendanceStatus"` // pending, confirmed, rejected, blocked
	AttendanceConfirmedBy *uint   `json:"attendanceConfirmedBy"`
	AttendanceConfirmedAt *string `json:"attendanceConfirmedAt"`

	// Finance Payroll Status (To be confirmed by Finance after HR)
	PayrollStatus      string  `json:"payrollStatus"` // pending, confirmed, rejected, blocked
	PayrollConfirmedBy *uint   `json:"payrollConfirmedBy"`
	PayrollConfirmedAt *string `json:"payrollConfirmedAt"`
}

// GetEmployeesWithSalaryInfo retrieves employees with their salary information for pre-payroll review
// This is used BEFORE creating any payroll run to review employee salary configurations
func (s *PrePayrollService) GetEmployeesWithSalaryInfo(tenantID *uint, status string, page, pageSize int) ([]EmployeeSalaryInfo, map[string]interface{}, error) {
	// Get active employees with full details including department and designation
	employees, _, totalPages, err := s.employeeService.ListEmployeesWithDetails(page, pageSize)
	if err != nil {
		return nil, nil, err
	}

	var results []EmployeeSalaryInfo
	var readyCount, missingSalaryCount, missingBankCount, missingStatutoryCount int

	for _, emp := range employees {
		if status != "" && emp.Status != status {
			continue
		}

		var hasSalary, hasBankInfo, hasStatutoryInfo, readyForPayroll bool
		var basicSalary, grossSalary float64
		var bankName, bankAccount, tinNumber, nssfNumber string
		var salaryGrade, effectiveDate string

		if emp.PositionID != nil {
			if structure, err := s.salaryStructureRepo.GetByJobPositionID(tenantID, *emp.PositionID); err == nil && structure != nil {
				hasSalary = true
				basicSalary = structure.Basic
				if structure.GrossSalary > 0 {
					grossSalary = structure.GrossSalary
				} else {
					grossSalary = structure.Basic + structure.HRA + structure.Transport + structure.Medical + structure.OtherAllowances
				}
				salaryGrade = structure.Grade
			}
		}

		hasBankInfo = bankName != "" && bankAccount != ""
		hasStatutoryInfo = tinNumber != "" && nssfNumber != ""
		readyForPayroll = hasSalary

		// Update counters
		if readyForPayroll {
			readyCount++
		}
		if !hasSalary {
			missingSalaryCount++
		}
		if !hasBankInfo {
			missingBankCount++
		}
		if !hasStatutoryInfo {
			missingStatutoryCount++
		}

		// Handle nil department and designation pointers
		department := ""
		if emp.Department != nil {
			department = *emp.Department
		}

		designation := ""
		if emp.Designation != nil {
			designation = *emp.Designation
		}

		result := EmployeeSalaryInfo{
			EmployeeID:       emp.ID,
			EmployeeCode:     emp.EmployeeID,
			EmployeeName:     emp.FullName,
			Department:       department,
			Designation:      designation,
			EmploymentStatus: emp.Status,

			// Basic Salary Information
			BasicSalary:         basicSalary,
			GrossSalary:         grossSalary,
			SalaryCurrency:      "TZS",
			SalaryGrade:         salaryGrade,
			SalaryEffectiveDate: effectiveDate,

			// Bank Information for Payment
			BankName:    bankName,
			BankAccount: bankAccount,

			// Statutory Information
			TINNumber:    tinNumber,
			NSSFNumber:   nssfNumber,
			HasHESLBLoan: false,

			// Status Tracking
			HasSalaryConfigured: hasSalary,
			HasBankInfo:         hasBankInfo,
			HasStatutoryInfo:    hasStatutoryInfo,
			ReadyForPayroll:     readyForPayroll,

			// HR Attendance Status (To be confirmed by HR)
			AttendanceStatus:      "pending",
			AttendanceConfirmedBy: nil,
			AttendanceConfirmedAt: nil,

			// Finance Payroll Status (To be confirmed by Finance after HR)
			PayrollStatus:      "pending",
			PayrollConfirmedBy: nil,
			PayrollConfirmedAt: nil,
		}

		results = append(results, result)
	}

	// Prepare metadata
	metadata := map[string]interface{}{
		"totalEmployees":       len(employees),
		"activeEmployees":      len(results), // Only active employees in results
		"readyForPayroll":      readyCount,
		"missingSalaryInfo":    missingSalaryCount,
		"missingBankInfo":      missingBankCount,
		"missingStatutoryInfo": missingStatutoryCount,
		"page":                 page,
		"pageSize":             pageSize,
		"totalPages":           totalPages,
	}

	return results, metadata, nil
}

// ConfirmAttendance confirms employee attendance for payroll (HR role only)
func (s *PrePayrollService) ConfirmAttendance(employeeID uint, confirmedBy uint, comments string) error {
	// This would update the attendance status in the database
	// For now, return success as placeholder
	return nil
}

// ConfirmPayroll confirms payroll calculations for employee (Finance role only)
func (s *PrePayrollService) ConfirmPayroll(employeeID uint, confirmedBy uint, comments string) error {
	// This would update the payroll status in the database
	// For now, return success as placeholder
	return nil
}
