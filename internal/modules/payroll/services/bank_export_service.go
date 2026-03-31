package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
)

// BankExportService handles bank export functionality
type BankExportService struct {
	db *gorm.DB
}

// NewBankExportService creates a new bank export service
func NewBankExportService(db *gorm.DB) *BankExportService {
	return &BankExportService{db: db}
}

// BankExportResult represents the result of a bank export operation
type BankExportResult struct {
	CRDBEmployees      int
	CRDBTotalAmount    float64
	DTBEmployees       int
	DTBTotalAmount     float64
	ValidationErrors   []string
	PaymentDate        time.Time
	CompanyCRDBAccount string
	CompanyDTBAccount  string
}

// ValidateExportRequirements validates if a payroll run can be exported
func (s *BankExportService) ValidateExportRequirements(payrollRunID uint) (bool, []string, error) {
	var errors []string

	// Check if payroll run exists and is LOCKED
	var payrollRun models.PayrollRun
	if err := s.db.First(&payrollRun, payrollRunID).Error; err != nil {
		return false, nil, fmt.Errorf("payroll run not found: %v", err)
	}

	if payrollRun.Status != models.PayrollRunStatusClosed {
		errors = append(errors, "Payroll run status must be CLOSED for export")
	}

	// Check if all employees have bank account numbers and net pay > 0
	var employees []models.PayrollRunEmployee
	if err := s.db.Where("payroll_run_id = ?", payrollRunID).Find(&employees).Error; err != nil {
		return false, nil, fmt.Errorf("failed to fetch employees: %v", err)
	}

	for _, emp := range employees {
		// Check if employee has bank configuration
		var config models.EmployeePayrollConfiguration
		if err := s.db.Where("employee_id = ? AND is_active = true", emp.EmployeeID).First(&config).Error; err != nil {
			errors = append(errors, fmt.Sprintf("Employee %s (%s) has no active payroll configuration", emp.EmployeeName, emp.EmployeeCode))
			continue
		}

		// Check bank account number
		if config.BankAccountNumber == "" {
			errors = append(errors, fmt.Sprintf("Employee %s (%s) has no bank account number", emp.EmployeeName, emp.EmployeeCode))
		}

		// Check net pay
		if emp.NetPay <= 0 {
			errors = append(errors, fmt.Sprintf("Employee %s (%s) has zero or negative net pay: %.2f TZS", emp.EmployeeName, emp.EmployeeCode, emp.NetPay))
		}
	}

	// Check if company bank accounts are configured
	var bankSettings models.BankSettings
	if err := s.db.First(&bankSettings).Error; err != nil {
		errors = append(errors, "Company bank account settings are not configured")
	} else {
		if bankSettings.CompanyCRDBAccount == "" {
			errors = append(errors, "Company CRDB account number is not configured")
		}
		if bankSettings.CompanyDTBAccount == "" {
			errors = append(errors, "Company DTB account number is not configured")
		}
		if bankSettings.DefaultPaymentDate.IsZero() {
			errors = append(errors, "Default payment date is not configured")
		}
	}

	return len(errors) == 0, errors, nil
}

// GenerateCRDBExport generates CRDB bank export file
func (s *BankExportService) GenerateCRDBExport(payrollRunID uint, paymentDate time.Time) (*excelize.File, *BankExportResult, error) {
	result, err := s.prepareExportData(payrollRunID, paymentDate, "CRDB")
	if err != nil {
		return nil, nil, err
	}

	if result.CRDBEmployees == 0 {
		return nil, result, fmt.Errorf("no CRDB employees found for export")
	}

	f := excelize.NewFile()
	sheet := "Payments"
	f.SetSheetName("Sheet1", sheet)

	// Set column headers
	headers := []string{
		"Beneficiary Account",
		"Beneficiary Name",
		"Amount",
		"Currency",
		"Narration",
		"Reference",
		"Debit Account",
		"Transaction Date",
	}

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, s.getHeaderStyle(f))
	}

	// Add employee rows
	row := 2
	var employees []models.PayrollRunEmployee
	s.db.Where("payroll_run_id = ?", payrollRunID).Order("employee_code asc").Find(&employees)

	for _, emp := range employees {
		var config models.EmployeePayrollConfiguration
		if err := s.db.Where("employee_id = ? AND is_active = true", emp.EmployeeID).First(&config).Error; err != nil {
			continue
		}

		if config.BankName != "CRDB" || emp.NetPay <= 0 {
			continue
		}

		// CRDB format row
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), config.BankAccountNumber)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), strings.ToUpper(emp.EmployeeName))
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), int(emp.NetPay)) // Integer amount
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "TZS")
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("Salary %s %d", time.Month(result.PaymentDate.Month()).String(), result.PaymentDate.Year()))
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("GTL-%s-%d-%s",
			strings.ToUpper(time.Month(result.PaymentDate.Month()).String()[:3]),
			result.PaymentDate.Year(),
			emp.EmployeeCode))
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), result.CompanyCRDBAccount)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), result.PaymentDate.Format("02/01/2006"))

		// Set amount column as number format
		f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), s.getNumberStyle(f))

		row++
	}

	// Add total row
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "TOTAL")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), int(result.CRDBTotalAmount))
	f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), s.getNumberStyle(f))
	f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), s.getTotalStyle(f))

	// Set column widths
	f.SetColWidth(sheet, "A", "A", 20)
	f.SetColWidth(sheet, "B", "B", 25)
	f.SetColWidth(sheet, "C", "C", 15)
	f.SetColWidth(sheet, "D", "D", 10)
	f.SetColWidth(sheet, "E", "E", 20)
	f.SetColWidth(sheet, "F", "F", 25)
	f.SetColWidth(sheet, "G", "G", 20)
	f.SetColWidth(sheet, "H", "H", 15)

	return f, result, nil
}

// GenerateDTBExport generates DTB bank export file
func (s *BankExportService) GenerateDTBExport(payrollRunID uint, paymentDate time.Time) (*excelize.File, *BankExportResult, error) {
	result, err := s.prepareExportData(payrollRunID, paymentDate, "DTB")
	if err != nil {
		return nil, nil, err
	}

	if result.DTBEmployees == 0 {
		return nil, result, fmt.Errorf("no DTB employees found for export")
	}

	f := excelize.NewFile()
	sheet := "BulkPayment"
	f.SetSheetName("Sheet1", sheet)

	// Set column headers
	headers := []string{
		"Account Number",
		"Beneficiary Name",
		"Amount",
		"Currency",
		"Payment Details",
		"Reference Number",
		"Value Date",
		"Bank Code",
	}

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, s.getHeaderStyle(f))
	}

	// Add employee rows
	row := 2
	var employees []models.PayrollRunEmployee
	s.db.Where("payroll_run_id = ?", payrollRunID).Order("employee_code asc").Find(&employees)

	for _, emp := range employees {
		var config models.EmployeePayrollConfiguration
		if err := s.db.Where("employee_id = ? AND is_active = true", emp.EmployeeID).First(&config).Error; err != nil {
			continue
		}

		if config.BankName != "DTB" || emp.NetPay <= 0 {
			continue
		}

		// DTB format row
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), config.BankAccountNumber)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), strings.ToUpper(emp.EmployeeName))
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), int(emp.NetPay)) // Integer amount
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "TZS")
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("Salary %s %d", time.Month(result.PaymentDate.Month()).String(), result.PaymentDate.Year()))
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("GTL-%s-%d-%s",
			strings.ToUpper(time.Month(result.PaymentDate.Month()).String()[:3]),
			result.PaymentDate.Year(),
			emp.EmployeeCode))
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), result.PaymentDate.Format("02/01/2006"))
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), "085") // DTB bank code

		// Set amount column as number format
		f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), s.getNumberStyle(f))

		row++
	}

	// Add total row
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "TOTAL")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), int(result.DTBTotalAmount))
	f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), s.getNumberStyle(f))
	f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), s.getTotalStyle(f))

	// Set column widths
	f.SetColWidth(sheet, "A", "A", 20)
	f.SetColWidth(sheet, "B", "B", 25)
	f.SetColWidth(sheet, "C", "C", 15)
	f.SetColWidth(sheet, "D", "D", 10)
	f.SetColWidth(sheet, "E", "E", 20)
	f.SetColWidth(sheet, "F", "F", 25)
	f.SetColWidth(sheet, "G", "G", 15)
	f.SetColWidth(sheet, "H", "H", 10)

	return f, result, nil
}

// GenerateSummaryExport generates internal payment summary file
func (s *BankExportService) GenerateSummaryExport(payrollRunID uint) (*excelize.File, error) {
	var payrollRun models.PayrollRun
	if err := s.db.First(&payrollRun, payrollRunID).Error; err != nil {
		return nil, fmt.Errorf("payroll run not found: %v", err)
	}

	var employees []models.PayrollRunEmployee
	if err := s.db.Where("payroll_run_id = ?", payrollRunID).Order("employee_code asc").Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch employees: %v", err)
	}

	f := excelize.NewFile()
	sheet := "Summary"
	f.SetSheetName("Sheet1", sheet)

	// Set column headers
	headers := []string{
		"SR.NO",
		"Employee Name",
		"Department",
		"Bank",
		"Account Number",
		"Gross Salary",
		"Total Deductions",
		"Net Pay",
		"Period",
		"Status",
	}

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, s.getHeaderStyle(f))
	}

	// Add employee rows
	row := 2
	var totalGross, totalDeductions, totalNetPay float64

	for _, emp := range employees {
		var config models.EmployeePayrollConfiguration
		if err := s.db.Where("employee_id = ?", emp.EmployeeID).First(&config).Error; err != nil {
			continue
		}

		status := "Included"
		if emp.NetPay <= 0 {
			status = "Zero net pay - excluded"
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), emp.EmployeeCode)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), emp.EmployeeName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), emp.Department)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), config.BankName)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), config.BankAccountNumber)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), emp.GrossSalary)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), emp.TotalDeductions)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), emp.NetPay)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), payrollRun.RunName)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), status)

		// Set number formats
		f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("H%d", row), s.getCurrencyStyle(f))

		totalGross += emp.GrossSalary
		totalDeductions += emp.TotalDeductions
		totalNetPay += emp.NetPay
		row++
	}

	// Add total row
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "TOTAL")
	f.SetCellValue(sheet, fmt.Sprintf("F%d", row), totalGross)
	f.SetCellValue(sheet, fmt.Sprintf("G%d", row), totalDeductions)
	f.SetCellValue(sheet, fmt.Sprintf("H%d", row), totalNetPay)
	f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("H%d", row), s.getTotalStyle(f))

	// Set column widths
	f.SetColWidth(sheet, "A", "A", 10)
	f.SetColWidth(sheet, "B", "B", 25)
	f.SetColWidth(sheet, "C", "C", 20)
	f.SetColWidth(sheet, "D", "D", 15)
	f.SetColWidth(sheet, "E", "E", 20)
	f.SetColWidth(sheet, "F", "F", 15)
	f.SetColWidth(sheet, "G", "G", 15)
	f.SetColWidth(sheet, "H", "H", 15)
	f.SetColWidth(sheet, "I", "I", 20)
	f.SetColWidth(sheet, "J", "J", 20)

	return f, nil
}

// prepareExportData prepares data for bank export
func (s *BankExportService) prepareExportData(payrollRunID uint, paymentDate time.Time, bankName string) (*BankExportResult, error) {
	result := &BankExportResult{
		PaymentDate: paymentDate,
	}

	// Get company bank settings
	var bankSettings models.BankSettings
	if err := s.db.First(&bankSettings).Error; err != nil {
		return nil, fmt.Errorf("bank settings not configured: %v", err)
	}
	result.CompanyCRDBAccount = bankSettings.CompanyCRDBAccount
	result.CompanyDTBAccount = bankSettings.CompanyDTBAccount

	// Get employees for this payroll run
	var employees []models.PayrollRunEmployee
	if err := s.db.Where("payroll_run_id = ?", payrollRunID).Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch employees: %v", err)
	}

	// Count employees and calculate totals by bank
	for _, emp := range employees {
		if emp.NetPay <= 0 {
			continue // Skip employees with zero net pay
		}

		var config models.EmployeePayrollConfiguration
		if err := s.db.Where("employee_id = ? AND is_active = true", emp.EmployeeID).First(&config).Error; err != nil {
			continue
		}

		switch config.BankName {
		case "CRDB":
			result.CRDBEmployees++
			result.CRDBTotalAmount += emp.NetPay
		case "DTB":
			result.DTBEmployees++
			result.DTBTotalAmount += emp.NetPay
		}
	}

	return result, nil
}

// Excel styling functions
func (s *BankExportService) getHeaderStyle(f *excelize.File) int {
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"2F5496"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	return style
}

func (s *BankExportService) getNumberStyle(f *excelize.File) int {
	style, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right"},
		NumFmt:    3, // Number format with thousands separator, no decimals
	})
	return style
}

func (s *BankExportService) getCurrencyStyle(f *excelize.File) int {
	style, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right"},
		NumFmt:    4, // Currency format
	})
	return style
}

func (s *BankExportService) getTotalStyle(f *excelize.File) int {
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Alignment: &excelize.Alignment{Horizontal: "right"},
		NumFmt:    4, // Currency format
	})
	return style
}
