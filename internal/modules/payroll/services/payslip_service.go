package services

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"

	deptRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
	empModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	empRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/repositories"
	posRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/repositories"
	"github.com/xuri/excelize/v2"
)

// PayslipService handles business logic for payslips
type PayslipService struct {
	repo          *repositories.PayslipRepository
	taxCalculator *TaxCalculator
	empRepo       *empRepos.EmployeeRepository
	empSalaryRepo *empRepos.EmployeeSalaryRepository
	deptRepo      *deptRepos.DepartmentRepository
	posRepo       *posRepos.JobPositionRepository
}

// NewPayslipService creates a new service instance
func NewPayslipService() *PayslipService {
	return &PayslipService{
		repo:          repositories.NewPayslipRepository(),
		taxCalculator: NewTaxCalculator(),
		empRepo:       empRepos.NewEmployeeRepository(),
		empSalaryRepo: empRepos.NewEmployeeSalaryRepository(),
		deptRepo:      deptRepos.NewDepartmentRepository(),
		posRepo:       posRepos.NewJobPositionRepository(),
	}
}

// GetPayslip retrieves a payslip by ID
func (s *PayslipService) GetPayslip(id uint, tenantID *uint) (*models.Payslip, error) {
	return s.repo.GetByID(id, tenantID)
}

// GetPayslipByEmployeeAndPeriod retrieves a payslip for a specific employee and period
func (s *PayslipService) GetPayslipByEmployeeAndPeriod(employeeID uint, payMonth, payYear int, tenantID *uint) (*models.Payslip, error) {
	return s.repo.GetByEmployeeAndPeriod(employeeID, payMonth, payYear, tenantID)
}

// GetPayslipWithItems retrieves a payslip with all its items
func (s *PayslipService) GetPayslipWithItems(id uint, tenantID *uint) (map[string]interface{}, error) {
	payslip, err := s.repo.GetByID(id, tenantID)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.GetItemsByPayslipID(id)
	if err != nil {
		return nil, err
	}

	// Organize items by type
	var earnings []map[string]interface{}
	var deductions []map[string]interface{}
	var employerContributions []map[string]interface{}

	for _, item := range items {
		itemData := map[string]interface{}{
			"code":   item.ComponentCode,
			"name":   item.ComponentName,
			"amount": item.Amount,
		}

		switch item.ItemType {
		case models.ComponentTypeEarning:
			earnings = append(earnings, itemData)
		case models.ComponentTypeDeduction:
			deductions = append(deductions, itemData)
		case models.ComponentTypeEmployerContribution:
			employerContributions = append(employerContributions, itemData)
		}
	}

	return map[string]interface{}{
		"id":                    payslip.ID,
		"employeeId":            payslip.EmployeeID,
		"empId":                 payslip.EmployeeCode,
		"employeeName":          payslip.EmployeeName,
		"employeePhoto":         payslip.EmployeePhoto,
		"department":            payslip.Department,
		"designation":           payslip.Designation,
		"dateOfJoining":         payslip.DateOfJoining,
		"payPeriod":             payslip.PayPeriod,
		"payMonth":              payslip.PayMonth,
		"payYear":               payslip.PayYear,
		"bankName":              payslip.BankName,
		"bankAccount":           payslip.BankAccount,
		"tinNumber":             payslip.TINNumber,
		"nssfNumber":            payslip.NSSFNumber,
		"nhifNumber":            payslip.NHIFNumber,
		"earnings":              earnings,
		"deductions":            deductions,
		"employerContributions": employerContributions,
		"summary": map[string]interface{}{
			"grossSalary":     payslip.GrossSalary,
			"totalDeductions": payslip.TotalDeductions,
			"netPay":          payslip.NetPay,
		},
		"ytd": map[string]interface{}{
			"grossSalary": payslip.YTDGross,
			"tax":         payslip.YTDTax,
			"nssf":        payslip.YTDNSSF,
			"nhif":        payslip.YTDNHIF,
			"netPay":      payslip.YTDNet,
		},
		"attendance": map[string]interface{}{
			"workingDays":   payslip.WorkingDays,
			"daysWorked":    payslip.DaysWorked,
			"lopDays":       payslip.LOPDays,
			"overtimeHours": payslip.OvertimeHours,
		},
		"status":    payslip.Status,
		"createdAt": payslip.CreatedAt,
	}, nil
}

// ListPayslips lists payslips with filters
func (s *PayslipService) ListPayslips(tenantID *uint, payMonth, payYear int, employeeID, departmentID *uint, status string, page, pageSize int) ([]models.Payslip, int64, error) {
	return s.repo.List(tenantID, payMonth, payYear, employeeID, departmentID, status, page, pageSize)
}

// GetEmployeePayslipHistory retrieves payslip history for an employee
func (s *PayslipService) GetEmployeePayslipHistory(employeeID uint, tenantID *uint, page, pageSize int) ([]models.Payslip, int64, error) {
	return s.repo.List(tenantID, 0, 0, &employeeID, nil, "", page, pageSize)
}

// GetSummary retrieves payslip summary for a period
func (s *PayslipService) GetSummary(tenantID *uint, employeeID *uint, payMonth, payYear int) (map[string]interface{}, error) {
	return s.repo.GetSummary(tenantID, employeeID, payMonth, payYear)
}

// ReleasePayslips releases payslips for employees to view
func (s *PayslipService) ReleasePayslips(runID uint, tenantID *uint) error {
	payslips, err := s.repo.GetByPayrollRunID(runID)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, payslip := range payslips {
		payslip.Status = models.PayslipStatusReleased
		payslip.ReleasedAt = &now
		if err := s.repo.Update(&payslip); err != nil {
			return err
		}
	}

	return nil
}

// GeneratePDF generates a PDF payslip (returns bytes)
func (s *PayslipService) GeneratePDF(payslipID uint, tenantID *uint) ([]byte, string, error) {
	payslipData, err := s.GetPayslipWithItems(payslipID, tenantID)
	if err != nil {
		return nil, "", err
	}

	// For now, return an error indicating PDF generation is not implemented
	// In production, you would use a PDF library like github.com/jung-kurt/gofpdf
	_ = payslipData // Use the data in actual implementation

	return nil, "", errors.New("PDF generation requires additional library integration")
}

// GenerateExcel generates an Excel payslip
func (s *PayslipService) GenerateExcel(payslipID uint, tenantID *uint) ([]byte, string, error) {
	payslipData, err := s.GetPayslipWithItems(payslipID, tenantID)
	if err != nil {
		return nil, "", err
	}

	f := excelize.NewFile()
	sheet := "Payslip"
	f.SetSheetName("Sheet1", sheet)

	// Set column widths
	f.SetColWidth(sheet, "A", "A", 20)
	f.SetColWidth(sheet, "B", "B", 20)
	f.SetColWidth(sheet, "C", "C", 15)
	f.SetColWidth(sheet, "D", "D", 15)

	// Header
	f.SetCellValue(sheet, "A1", "PAYSLIP")
	f.MergeCell(sheet, "A1", "D1")
	f.SetCellValue(sheet, "A2", fmt.Sprintf("Pay Period: %s", payslipData["payPeriod"]))
	f.MergeCell(sheet, "A2", "D2")

	// Employee Info
	f.SetCellValue(sheet, "A4", "Employee Information")
	f.MergeCell(sheet, "A4", "D4")
	f.SetCellValue(sheet, "A5", "Name:")
	f.SetCellValue(sheet, "B5", payslipData["employeeName"])
	f.SetCellValue(sheet, "C5", "Employee ID:")
	f.SetCellValue(sheet, "D5", payslipData["empId"])
	f.SetCellValue(sheet, "A6", "Department:")
	f.SetCellValue(sheet, "B6", payslipData["department"])
	f.SetCellValue(sheet, "C6", "Designation:")
	f.SetCellValue(sheet, "D6", payslipData["designation"])

	row := 8

	// Earnings
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Earnings")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "Amount")
	row++
	if earnings, ok := payslipData["earnings"].([]map[string]interface{}); ok {
		for _, earning := range earnings {
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), earning["name"])
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), earning["amount"])
			row++
		}
	}

	row++
	// Deductions
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Deductions")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "Amount")
	row++
	if deductions, ok := payslipData["deductions"].([]map[string]interface{}); ok {
		for _, deduction := range deductions {
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), deduction["name"])
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), deduction["amount"])
			row++
		}
	}

	row++
	// Summary
	if summary, ok := payslipData["summary"].(map[string]interface{}); ok {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Gross Salary:")
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), summary["grossSalary"])
		row++
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Total Deductions:")
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), summary["totalDeductions"])
		row++
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Net Pay:")
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), summary["netPay"])
	}

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("payslip_%s_%d_%d.xlsx",
		payslipData["empId"],
		payslipData["payMonth"],
		payslipData["payYear"])

	return buf.Bytes(), filename, nil
}

// BulkDownloadPayslips generates a ZIP file with multiple payslips
func (s *PayslipService) BulkDownloadPayslips(payslipIDs []uint, tenantID *uint) ([]byte, string, error) {
	// For bulk download, we would create a ZIP file with individual payslips
	// This is a placeholder - actual implementation would use archive/zip
	return nil, "", errors.New("bulk download requires additional implementation")
}

// SendPayslipEmail sends payslip to employee email
func (s *PayslipService) SendPayslipEmail(payslipID uint, tenantID *uint) error {
	payslip, err := s.repo.GetByID(payslipID, tenantID)
	if err != nil {
		return err
	}

	// This is a placeholder - actual implementation would use email service
	_ = payslip
	return errors.New("email sending requires email service configuration")
}

// BulkEmailPayslips sends payslips to multiple employees
func (s *PayslipService) BulkEmailPayslips(runID uint, tenantID *uint) (map[string]interface{}, error) {
	payslips, err := s.repo.GetByPayrollRunID(runID)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"total":   len(payslips),
		"sent":    0,
		"failed":  len(payslips), // All fail for now since email not configured
		"message": "Email service not configured",
	}

	return result, nil
}

// GetPayslipSummaryByRun returns summary for a payroll run's payslips
func (s *PayslipService) GetPayslipSummaryByRun(runID uint, tenantID *uint) (map[string]interface{}, error) {
	payslips, err := s.repo.GetByPayrollRunID(runID)
	if err != nil {
		return nil, err
	}

	var totalGross, totalDeductions, totalNet float64
	var releaseCount, draftCount int

	for _, p := range payslips {
		totalGross += p.GrossSalary
		totalDeductions += p.TotalDeductions
		totalNet += p.NetPay

		if p.Status == models.PayslipStatusReleased {
			releaseCount++
		} else {
			draftCount++
		}
	}

	return map[string]interface{}{
		"payrollRunId":    runID,
		"totalPayslips":   len(payslips),
		"releasedCount":   releaseCount,
		"draftCount":      draftCount,
		"totalGross":      totalGross,
		"totalDeductions": totalDeductions,
		"totalNet":        totalNet,
	}, nil
}

// GetEmployeeSalaryStructures retrieves the salary structures for all employees
func (s *PayslipService) GetEmployeeSalaryStructures(tenantID *uint) ([]map[string]interface{}, error) {
	employees, _, err := s.empRepo.List(1000, 0) // Fetch up to 1000 employees
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	taxYear := time.Now().Year()

	for _, emp := range employees {
		// Only include active employees
		if emp.Status != empModels.StatusActive {
			continue
		}

		// Fetch salary components
		salary, _ := s.empSalaryRepo.FindByEmployeeID(emp.ID)

		var basicSalary, housingAllowance, transportAllowance, otherAllowances float64

		if salary != nil {
			if salary.BasicSalary != nil {
				basicSalary = *salary.BasicSalary
			}
			if salary.HouseRentAllowance != nil {
				housingAllowance = *salary.HouseRentAllowance
			}
			if salary.TransportAllowance != nil {
				transportAllowance = *salary.TransportAllowance
			}
			if salary.OtherAllowances != nil {
				otherAllowances = *salary.OtherAllowances
			}
			if salary.SpecialAllowance != nil {
				otherAllowances += *salary.SpecialAllowance
			}
		} else if emp.Salary != nil {
			// Fallback to employee base salary if no component record found
			basicSalary = *emp.Salary
		}

		// Skip employees with no salary info
		if basicSalary <= 0 {
			continue
		}

		grossSalary := basicSalary + housingAllowance + transportAllowance + otherAllowances

		// Calculate statutory deductions
		taxResult := s.taxCalculator.CalculateAllDeductions(grossSalary, taxYear, tenantID)

		// Get department name
		departmentName := "Unassigned"
		if emp.DepartmentID != nil {
			dept, err := s.deptRepo.FindByID(*emp.DepartmentID)
			if err == nil && dept != nil {
				departmentName = dept.Name
			}
		}

		// Get designation
		designation := "Staff"
		if emp.PositionID != nil {
			pos, err := s.posRepo.FindByID(*emp.PositionID)
			if err == nil && pos != nil {
				designation = pos.Title
			}
		}

		results = append(results, map[string]interface{}{
			"employeeId":      emp.ID,
			"employeeName":    emp.FullName(),
			"department":      departmentName,
			"designation":     designation,
			"grossSalary":     grossSalary,
			"totalDeductions": taxResult.TotalEmployeeDeductions,
			"netPay":          taxResult.NetPay,
			"earnings": []map[string]interface{}{
				{"name": "Basic Pay", "amount": basicSalary},
				{"name": "Housing Allowance", "amount": housingAllowance},
				{"name": "Transport Allowance", "amount": transportAllowance},
				{"name": "Other Allowances", "amount": otherAllowances},
			},
			"deductions": []map[string]interface{}{
				{"name": "PAYE Tax", "amount": taxResult.PAYE},
				{"name": "NSSF (10%)", "amount": taxResult.NSSFEmployee},
				{"name": "NHIF", "amount": taxResult.NHIFEmployee},
			},
		})
	}

	return results, nil
}

// GetEmployeeSalaryStructure retrieves the salary structure for a single employee
func (s *PayslipService) GetEmployeeSalaryStructure(employeeID uint, tenantID *uint) (map[string]interface{}, error) {
	emp, err := s.empRepo.FindByID(employeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Fetch salary components
	salary, _ := s.empSalaryRepo.FindByEmployeeID(emp.ID)

	var basicSalary, housingAllowance, transportAllowance, otherAllowances float64

	if salary != nil {
		if salary.BasicSalary != nil {
			basicSalary = *salary.BasicSalary
		}
		if salary.HouseRentAllowance != nil {
			housingAllowance = *salary.HouseRentAllowance
		}
		if salary.TransportAllowance != nil {
			transportAllowance = *salary.TransportAllowance
		}
		if salary.OtherAllowances != nil {
			otherAllowances = *salary.OtherAllowances
		}
		if salary.SpecialAllowance != nil {
			otherAllowances += *salary.SpecialAllowance
		}
	} else if emp.Salary != nil {
		basicSalary = *emp.Salary
	}

	grossSalary := basicSalary + housingAllowance + transportAllowance + otherAllowances
	taxYear := time.Now().Year()

	// Calculate statutory deductions
	taxResult := s.taxCalculator.CalculateAllDeductions(grossSalary, taxYear, tenantID)

	// Get department name
	departmentName := "Unassigned"
	if emp.DepartmentID != nil {
		dept, err := s.deptRepo.FindByID(*emp.DepartmentID)
		if err == nil && dept != nil {
			departmentName = dept.Name
		}
	}

	// Get designation
	designation := "Staff"
	if emp.PositionID != nil {
		pos, err := s.posRepo.FindByID(*emp.PositionID)
		if err == nil && pos != nil {
			designation = pos.Title
		}
	}

	return map[string]interface{}{
		"employeeId":      emp.ID,
		"employeeName":    emp.FullName(),
		"department":      departmentName,
		"designation":     designation,
		"grossSalary":     grossSalary,
		"totalDeductions": taxResult.TotalEmployeeDeductions,
		"netPay":          taxResult.NetPay,
		"earnings": []map[string]interface{}{
			{"name": "Basic Pay", "amount": basicSalary},
			{"name": "Housing Allowance", "amount": housingAllowance},
			{"name": "Transport Allowance", "amount": transportAllowance},
			{"name": "Other Allowances", "amount": otherAllowances},
		},
		"deductions": []map[string]interface{}{
			{"name": "PAYE Tax", "amount": taxResult.PAYE},
			{"name": "NSSF (10%)", "amount": taxResult.NSSFEmployee},
			{"name": "NHIF", "amount": taxResult.NHIFEmployee},
		},
	}, nil
}

// FormatPayslipPeriod formats a pay period string
func FormatPayslipPeriod(month, year int) string {
	return time.Month(month).String()[:3] + " " + strconv.Itoa(year)
}
