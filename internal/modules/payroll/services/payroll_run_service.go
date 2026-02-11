package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/repositories"
	"gorm.io/gorm"
)

// PayrollRunService handles business logic for payroll runs
type PayrollRunService struct {
	repo                *repositories.PayrollRunRepository
	empRepo             *repositories.PayrollRunEmployeeRepository
	payslipRepo         *repositories.PayslipRepository
	loanRepo            *repositories.LoanRepository
	salaryStructureRepo *repositories.SalaryStructureRepository
	taxCalculator       *TaxCalculator
}

// NewPayrollRunService creates a new service instance
func NewPayrollRunService() *PayrollRunService {
	return &PayrollRunService{
		repo:                repositories.NewPayrollRunRepository(),
		empRepo:             repositories.NewPayrollRunEmployeeRepository(),
		payslipRepo:         repositories.NewPayslipRepository(),
		loanRepo:            repositories.NewLoanRepository(),
		salaryStructureRepo: repositories.NewSalaryStructureRepository(),
		taxCalculator:       NewTaxCalculator(),
	}
}

// CreatePayrollRun creates a new payroll run
func (s *PayrollRunService) CreatePayrollRun(run *models.PayrollRun) error {
	// Check if a run already exists for this period
	existing, err := s.repo.GetByPeriod(nil, run.PayMonth, run.PayYear)
	if err == nil && existing != nil {
		return errors.New("payroll run already exists for this period")
	}

	// Set default values
	run.Status = models.PayrollRunStatusDraft
	run.CurrentStep = 0
	run.IsLocked = false

	// Generate run name if not provided
	if run.RunName == "" {
		run.RunName = fmt.Sprintf("Payroll Run %s %d", time.Month(run.PayMonth).String(), run.PayYear)
	}

	// Set pay period string
	run.PayPeriod = fmt.Sprintf("%s %d", time.Month(run.PayMonth).String()[:3], run.PayYear)

	return s.repo.Create(run)
}

// GetPayrollRun retrieves a payroll run by ID
func (s *PayrollRunService) GetPayrollRun(id uint, tenantID *uint) (*models.PayrollRun, error) {
	return s.repo.GetByID(id, tenantID)
}

// UpdatePayrollRun updates a payroll run
func (s *PayrollRunService) UpdatePayrollRun(run *models.PayrollRun) error {
	if run.IsLocked && run.Status != models.PayrollRunStatusDraft {
		return errors.New("cannot update locked payroll run")
	}
	return s.repo.Update(run)
}

// DeletePayrollRun deletes a payroll run (only if in Draft status)
func (s *PayrollRunService) DeletePayrollRun(id uint, tenantID *uint) error {
	run, err := s.repo.GetByID(id, tenantID)
	if err != nil {
		return err
	}
	if run.Status != models.PayrollRunStatusDraft {
		return errors.New("can only delete payroll runs in Draft status")
	}
	return s.repo.Delete(id, tenantID)
}

// ListPayrollRuns lists payroll runs with filters
func (s *PayrollRunService) ListPayrollRuns(tenantID *uint, status string, payYear, payMonth, page, pageSize int) ([]models.PayrollRun, int64, error) {
	return s.repo.List(tenantID, status, payYear, payMonth, page, pageSize)
}

// GetCurrentPayrollRun retrieves the current active payroll run
func (s *PayrollRunService) GetCurrentPayrollRun(tenantID *uint) (*models.PayrollRun, error) {
	return s.repo.GetCurrentRun(tenantID)
}

// GetDashboard retrieves dashboard metrics
func (s *PayrollRunService) GetDashboard(tenantID *uint) (map[string]interface{}, error) {
	now := time.Now()
	currentMonth := int(now.Month())
	currentYear := now.Year()

	// Get current run
	currentRun, err := s.repo.GetCurrentRun(tenantID)
	var currentRunData map[string]interface{}
	if err == nil && currentRun != nil {
		currentRunData = map[string]interface{}{
			"id":          currentRun.ID,
			"runName":     currentRun.RunName,
			"status":      currentRun.Status,
			"currentStep": currentRun.CurrentStep,
			"payMonth":    currentRun.PayMonth,
			"payYear":     currentRun.PayYear,
		}
	}

	// Get last finalized run
	lastRun, err := s.repo.GetLastFinalized(tenantID)
	var lastRunData map[string]interface{}
	if err == nil && lastRun != nil {
		lastRunData = map[string]interface{}{
			"id":               lastRun.ID,
			"runName":          lastRun.RunName,
			"payPeriod":        lastRun.PayPeriod,
			"totalEmployees":   lastRun.TotalEmployees,
			"totalNet":         lastRun.TotalNet,
			"disbursementDate": lastRun.DisbursementDate,
		}
	}

	// Get recent runs
	recentRuns, _ := s.repo.GetRecentRuns(tenantID, 5)
	var recentRunsData []map[string]interface{}
	for _, run := range recentRuns {
		recentRunsData = append(recentRunsData, map[string]interface{}{
			"id":        run.ID,
			"runName":   run.RunName,
			"payPeriod": run.PayPeriod,
			"status":    run.Status,
			"totalNet":  run.TotalNet,
		})
	}

	// Loan summary
	loanSummary, _ := s.loanRepo.GetSummary(tenantID, nil)

	return map[string]interface{}{
		"currentRun":   currentRunData,
		"lastRun":      lastRunData,
		"recentRuns":   recentRunsData,
		"currentMonth": currentMonth,
		"currentYear":  currentYear,
		"loanSummary":  loanSummary,
	}, nil
}

// AdvanceStep advances the payroll run to the next step
func (s *PayrollRunService) AdvanceStep(runID uint, tenantID *uint, userID uint, userName string) error {
	run, err := s.repo.GetByID(runID, tenantID)
	if err != nil {
		return err
	}

	if run.IsLocked {
		return errors.New("payroll run is locked")
	}

	switch run.Status {
	case models.PayrollRunStatusDraft:
		run.Status = models.PayrollRunStatusInReview
		run.CurrentStep = 1
	case models.PayrollRunStatusInReview:
		run.Status = models.PayrollRunStatusApproved
		run.CurrentStep = 2
		run.ApprovedByID = &userID
		run.ApprovedByName = &userName
		now := time.Now()
		run.ApprovedAt = &now
	case models.PayrollRunStatusApproved:
		run.Status = models.PayrollRunStatusFinalized
		run.CurrentStep = 3
		run.FinalizedByID = &userID
		run.FinalizedByName = &userName
		now := time.Now()
		run.FinalizedAt = &now
		run.IsLocked = true
	case models.PayrollRunStatusFinalized:
		run.Status = models.PayrollRunStatusDisbursed
		run.CurrentStep = 4
		now := time.Now()
		run.DisbursementDate = &now
	case models.PayrollRunStatusDisbursed:
		run.Status = models.PayrollRunStatusClosed
		run.CurrentStep = 5
	default:
		return errors.New("cannot advance from current status")
	}

	return s.repo.Update(run)
}

// RevertStep reverts the payroll run to the previous step
func (s *PayrollRunService) RevertStep(runID uint, tenantID *uint) error {
	run, err := s.repo.GetByID(runID, tenantID)
	if err != nil {
		return err
	}

	switch run.Status {
	case models.PayrollRunStatusInReview:
		run.Status = models.PayrollRunStatusDraft
		run.CurrentStep = 0
	case models.PayrollRunStatusApproved:
		run.Status = models.PayrollRunStatusInReview
		run.CurrentStep = 1
		run.ApprovedByID = nil
		run.ApprovedByName = nil
		run.ApprovedAt = nil
	default:
		return errors.New("cannot revert from current status")
	}

	run.IsLocked = false
	return s.repo.Update(run)
}

// GetEmployees retrieves employees in a payroll run
func (s *PayrollRunService) GetEmployees(runID uint, search, departmentID string, hasChanges *bool, page, pageSize int) ([]models.PayrollRunEmployee, int64, error) {
	return s.empRepo.GetByRunID(runID, search, departmentID, hasChanges, page, pageSize)
}

// RunPreCheck performs pre-checks before processing payroll
func (s *PayrollRunService) RunPreCheck(runID uint, tenantID *uint) (*models.PreCheckResult, error) {
	run, err := s.repo.GetByID(runID, tenantID)
	if err != nil {
		return nil, err
	}

	result := &models.PreCheckResult{
		PayrollRunID:   runID,
		RunName:        run.RunName,
		TotalEmployees: run.TotalEmployees,
		ReadyCount:     0,
		IssuesCount:    0,
		Issues:         []models.PreCheckIssue{},
	}

	var criticalIssues int
	var highIssues int

	// Check 1: Verify salary structures/grades exist
	salaryStructureCount, err := s.salaryStructureRepo.CountActiveSalaryStructures(tenantID)
	if err != nil || salaryStructureCount == 0 {
		actionRequired := "Navigate to Payroll > Salary Structures and create at least one salary grade/structure before running payroll."
		resolutionURL := "/payroll/salary-structures"
		result.Issues = append(result.Issues, models.PreCheckIssue{
			ID:             "no_salary_structures",
			Type:           "missing_configuration",
			Category:       "salary_structure",
			Count:          1,
			Severity:       "High",
			EmployeeIDs:    []string{},
			Employees:      []models.PreCheckEmployee{},
			Description:    "No salary structures/grades are configured. Payroll cannot be calculated without salary grades.",
			ActionRequired: &actionRequired,
			ResolutionURL:  &resolutionURL,
		})
		criticalIssues++
	}

	// Check 2: Get available grades
	availableGrades, _ := s.salaryStructureRepo.GetDistinctGrades(tenantID)
	if len(availableGrades) == 0 && salaryStructureCount > 0 {
		actionRequired := "Ensure salary structures have grades assigned (e.g., G1, G2, G3)."
		resolutionURL := "/payroll/salary-structures"
		result.Issues = append(result.Issues, models.PreCheckIssue{
			ID:             "no_salary_grades",
			Type:           "missing_configuration",
			Category:       "salary_structure",
			Count:          1,
			Severity:       "High",
			EmployeeIDs:    []string{},
			Employees:      []models.PreCheckEmployee{},
			Description:    "Salary structures exist but no grades are defined. Please assign grades to salary structures.",
			ActionRequired: &actionRequired,
			ResolutionURL:  &resolutionURL,
		})
		criticalIssues++
	}

	// Check 3: Verify employees in the payroll run have salary data
	employees, _, err := s.empRepo.GetByRunID(runID, "", "", nil, 1, 10000)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve employees: %w", err)
	}

	if len(employees) == 0 {
		actionRequired := "Add employees to this payroll run before proceeding."
		resolutionURL := fmt.Sprintf("/payroll/runs/%d/employees", runID)
		result.Issues = append(result.Issues, models.PreCheckIssue{
			ID:             "no_employees",
			Type:           "missing_data",
			Category:       "employee",
			Count:          1,
			Severity:       "High",
			EmployeeIDs:    []string{},
			Employees:      []models.PreCheckEmployee{},
			Description:    "No employees are assigned to this payroll run.",
			ActionRequired: &actionRequired,
			ResolutionURL:  &resolutionURL,
		})
		criticalIssues++
	}

	// Check 4: Verify employees have gross salary set
	var employeesWithoutSalary []models.PreCheckEmployee
	var employeeIDsWithoutSalary []string
	var employeesWithZeroSalary []models.PreCheckEmployee
	var employeeIDsWithZeroSalary []string

	for _, emp := range employees {
		if emp.GrossSalary <= 0 {
			if emp.GrossSalary == 0 {
				employeesWithZeroSalary = append(employeesWithZeroSalary, models.PreCheckEmployee{
					ID:         emp.EmployeeCode,
					Name:       emp.EmployeeName,
					Department: emp.Department,
				})
				employeeIDsWithZeroSalary = append(employeeIDsWithZeroSalary, emp.EmployeeCode)
			} else {
				employeesWithoutSalary = append(employeesWithoutSalary, models.PreCheckEmployee{
					ID:         emp.EmployeeCode,
					Name:       emp.EmployeeName,
					Department: emp.Department,
				})
				employeeIDsWithoutSalary = append(employeeIDsWithoutSalary, emp.EmployeeCode)
			}
		}
	}

	if len(employeesWithoutSalary) > 0 {
		actionRequired := "Update employee salary information in their profiles or assign them to a salary structure."
		resolutionURL := "/employees"
		result.Issues = append(result.Issues, models.PreCheckIssue{
			ID:             "employees_without_salary",
			Type:           "missing_data",
			Category:       "employee",
			Count:          len(employeesWithoutSalary),
			Severity:       "High",
			EmployeeIDs:    employeeIDsWithoutSalary,
			Employees:      employeesWithoutSalary,
			Description:    fmt.Sprintf("%d employee(s) do not have salary information configured.", len(employeesWithoutSalary)),
			ActionRequired: &actionRequired,
			ResolutionURL:  &resolutionURL,
		})
		highIssues += len(employeesWithoutSalary)
	}

	if len(employeesWithZeroSalary) > 0 {
		actionRequired := "Verify and update salary amounts for these employees. Zero salary may be intentional for unpaid leave, but review is recommended."
		resolutionURL := "/employees"
		result.Issues = append(result.Issues, models.PreCheckIssue{
			ID:             "employees_zero_salary",
			Type:           "warning",
			Category:       "employee",
			Count:          len(employeesWithZeroSalary),
			Severity:       "Medium",
			EmployeeIDs:    employeeIDsWithZeroSalary,
			Employees:      employeesWithZeroSalary,
			Description:    fmt.Sprintf("%d employee(s) have zero gross salary. Please verify this is intentional.", len(employeesWithZeroSalary)),
			ActionRequired: &actionRequired,
			ResolutionURL:  &resolutionURL,
		})
	}

	// Check 5: Verify basic salary is set
	var employeesWithoutBasic []models.PreCheckEmployee
	var employeeIDsWithoutBasic []string

	for _, emp := range employees {
		if emp.BasicSalary <= 0 && emp.GrossSalary > 0 {
			employeesWithoutBasic = append(employeesWithoutBasic, models.PreCheckEmployee{
				ID:         emp.EmployeeCode,
				Name:       emp.EmployeeName,
				Department: emp.Department,
			})
			employeeIDsWithoutBasic = append(employeeIDsWithoutBasic, emp.EmployeeCode)
		}
	}

	if len(employeesWithoutBasic) > 0 {
		actionRequired := "Set basic salary for these employees. Basic salary is required for accurate tax calculations."
		resolutionURL := "/employees"
		result.Issues = append(result.Issues, models.PreCheckIssue{
			ID:             "employees_without_basic",
			Type:           "missing_data",
			Category:       "employee",
			Count:          len(employeesWithoutBasic),
			Severity:       "Medium",
			EmployeeIDs:    employeeIDsWithoutBasic,
			Employees:      employeesWithoutBasic,
			Description:    fmt.Sprintf("%d employee(s) have gross salary but no basic salary set.", len(employeesWithoutBasic)),
			ActionRequired: &actionRequired,
			ResolutionURL:  &resolutionURL,
		})
	}

	// Calculate summary
	result.IssuesCount = len(result.Issues)
	readyEmployees := len(employees) - len(employeesWithoutSalary)
	if readyEmployees < 0 {
		readyEmployees = 0
	}
	result.ReadyCount = readyEmployees

	// Determine if payroll can proceed
	// Cannot proceed if: no salary structures, no grades, no employees, or employees without salary
	if criticalIssues > 0 || len(employeesWithoutSalary) > 0 {
		result.CanProceed = false
		if criticalIssues > 0 {
			result.Summary = "Payroll cannot proceed. Critical configuration issues must be resolved first."
		} else {
			result.Summary = fmt.Sprintf("Payroll cannot proceed. %d employee(s) do not have salary information.", len(employeesWithoutSalary))
		}
	} else if len(employeesWithZeroSalary) > 0 || len(employeesWithoutBasic) > 0 {
		result.CanProceed = true
		result.Summary = fmt.Sprintf("%d of %d employees are ready. Review warnings before proceeding.", readyEmployees, len(employees))
	} else {
		result.CanProceed = true
		result.Summary = fmt.Sprintf("All %d employees are ready for payroll processing.", len(employees))
	}

	return result, nil
}

// CalculatePayroll calculates payroll for all employees in a run
func (s *PayrollRunService) CalculatePayroll(runID uint, tenantID *uint, taxYear int) error {
	run, err := s.repo.GetByID(runID, tenantID)
	if err != nil {
		return err
	}

	if run.Status != models.PayrollRunStatusDraft && run.Status != models.PayrollRunStatusInReview {
		return errors.New("can only calculate payroll for Draft or In Review runs")
	}

	// Run pre-checks before calculating payroll
	preCheckResult, err := s.RunPreCheck(runID, tenantID)
	if err != nil {
		return fmt.Errorf("pre-check failed: %w", err)
	}

	if !preCheckResult.CanProceed {
		// Build error message with details about issues
		errMsg := "Cannot calculate payroll: " + preCheckResult.Summary
		if len(preCheckResult.Issues) > 0 {
			for _, issue := range preCheckResult.Issues {
				if issue.Severity == "High" {
					errMsg += fmt.Sprintf("\n- %s (%d affected)", issue.Description, issue.Count)
				}
			}
		}
		return errors.New(errMsg)
	}

	// Get employees in the run
	employees, _, err := s.empRepo.GetByRunID(runID, "", "", nil, 1, 10000)
	if err != nil {
		return err
	}

	var totalGross, totalDeductions, totalNet, totalEmployerContributions float64
	var totalBasic, totalAllowances float64
	var totalPAYE, totalNSSFEmployee, totalNHIFEmployee, totalLoanDeductions float64
	var totalNSSFEmployer, totalNHIFEmployer, totalSDL, totalWCF float64

	for i := range employees {
		emp := &employees[i]

		// Calculate taxes
		taxResult := s.taxCalculator.CalculateAllDeductions(emp.GrossSalary, taxYear, tenantID)

		// Get loan deductions for this employee
		pendingRepayments, _ := s.loanRepo.GetPendingRepayments(emp.EmployeeID, time.Date(run.PayYear, time.Month(run.PayMonth), 28, 0, 0, 0, 0, time.UTC))
		var loanDeduction float64
		for _, rep := range pendingRepayments {
			loanDeduction += rep.EMIAmount
		}

		// Update employee record
		emp.TotalDeductions = taxResult.TotalEmployeeDeductions + loanDeduction
		emp.NetPay = emp.GrossSalary - emp.TotalDeductions

		// Accumulate totals
		totalGross += emp.GrossSalary
		totalBasic += emp.BasicSalary
		totalAllowances += (emp.GrossSalary - emp.BasicSalary)
		totalDeductions += emp.TotalDeductions
		totalNet += emp.NetPay
		totalPAYE += taxResult.PAYE
		totalNSSFEmployee += taxResult.NSSFEmployee
		totalNHIFEmployee += taxResult.NHIFEmployee
		totalLoanDeductions += loanDeduction
		totalNSSFEmployer += taxResult.NSSFEmployer
		totalNHIFEmployer += taxResult.NHIFEmployer
		totalSDL += taxResult.SDL
		totalWCF += taxResult.WCF
		totalEmployerContributions += taxResult.TotalEmployerContributions
	}

	// Update run totals
	run.TotalGross = totalGross
	run.TotalBasic = totalBasic
	run.TotalAllowances = totalAllowances
	run.TotalDeductions = totalDeductions
	run.TotalNet = totalNet
	run.TotalPAYE = totalPAYE
	run.TotalNSSFEmployee = totalNSSFEmployee
	run.TotalNHIFEmployee = totalNHIFEmployee
	run.TotalLoanDeductions = totalLoanDeductions
	run.TotalNSSFEmployer = totalNSSFEmployer
	run.TotalNHIFEmployer = totalNHIFEmployer
	run.TotalSDL = totalSDL
	run.TotalWCF = totalWCF
	run.TotalEmployerContributions = totalEmployerContributions

	return s.repo.Update(run)
}

// GeneratePayslips generates payslips for a finalized payroll run
func (s *PayrollRunService) GeneratePayslips(runID uint, tenantID *uint, taxYear int) error {
	run, err := s.repo.GetByID(runID, tenantID)
	if err != nil {
		return err
	}

	if run.Status != models.PayrollRunStatusFinalized && run.Status != models.PayrollRunStatusApproved {
		return errors.New("can only generate payslips for Approved or Finalized runs")
	}

	employees, _, err := s.empRepo.GetByRunID(runID, "", "", nil, 1, 10000)
	if err != nil {
		return err
	}

	var payslips []models.Payslip
	var allItems []models.PayslipItem

	for _, emp := range employees {
		// Calculate taxes
		taxResult := s.taxCalculator.CalculateAllDeductions(emp.GrossSalary, taxYear, tenantID)

		// Get YTD totals
		ytdGross, ytdTax, ytdNSSF, ytdNHIF, ytdNet, _ := s.payslipRepo.GetYTDTotals(emp.EmployeeID, run.PayYear, run.PayMonth-1)

		payslip := models.Payslip{
			PayrollRunID:               runID,
			EmployeeID:                 emp.EmployeeID,
			EmployeeCode:               emp.EmployeeCode,
			EmployeeName:               emp.EmployeeName,
			EmployeePhoto:              emp.EmployeePhoto,
			Department:                 emp.Department,
			DepartmentID:               emp.DepartmentID,
			Designation:                emp.Designation,
			PayPeriod:                  run.PayPeriod,
			PayMonth:                   run.PayMonth,
			PayYear:                    run.PayYear,
			EarningsTotal:              emp.GrossSalary,
			DeductionsTotal:            taxResult.TotalEmployeeDeductions,
			EmployerContributionsTotal: taxResult.TotalEmployerContributions,
			GrossSalary:                emp.GrossSalary,
			TotalDeductions:            emp.TotalDeductions,
			NetPay:                     emp.NetPay,
			YTDGross:                   ytdGross + emp.GrossSalary,
			YTDTax:                     ytdTax + taxResult.PAYE,
			YTDNSSF:                    ytdNSSF + taxResult.NSSFEmployee,
			YTDNHIF:                    ytdNHIF + taxResult.NHIFEmployee,
			YTDNet:                     ytdNet + emp.NetPay,
			WorkingDays:                22, // Default working days
			DaysWorked:                 emp.DaysWorked,
			LOPDays:                    emp.LOPDays,
			OvertimeHours:              emp.OvertimeHours,
			Status:                     models.PayslipStatusDraft,
		}
		payslips = append(payslips, payslip)
	}

	// Create payslips
	if err := s.payslipRepo.CreateBatch(payslips); err != nil {
		return err
	}

	// Create payslip items
	if len(allItems) > 0 {
		if err := s.payslipRepo.CreateItemsBatch(allItems); err != nil {
			return err
		}
	}

	return nil
}

// GetPayrollRunSummary returns a summary of a payroll run
func (s *PayrollRunService) GetPayrollRunSummary(runID uint, tenantID *uint) (map[string]interface{}, error) {
	run, err := s.repo.GetByID(runID, tenantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payroll run not found")
		}
		return nil, err
	}

	summary := map[string]interface{}{
		"id":             run.ID,
		"runName":        run.RunName,
		"payPeriod":      run.PayPeriod,
		"payMonth":       run.PayMonth,
		"payYear":        run.PayYear,
		"payFrequency":   run.PayFrequency,
		"status":         run.Status,
		"currentStep":    run.CurrentStep,
		"totalEmployees": run.TotalEmployees,
		"financialSummary": map[string]interface{}{
			"totalGross":      run.TotalGross,
			"totalBasic":      run.TotalBasic,
			"totalAllowances": run.TotalAllowances,
			"totalDeductions": run.TotalDeductions,
			"totalNet":        run.TotalNet,
		},
		"deductionsBreakdown": map[string]interface{}{
			"paye":           run.TotalPAYE,
			"nssfEmployee":   run.TotalNSSFEmployee,
			"nhifEmployee":   run.TotalNHIFEmployee,
			"loanDeductions": run.TotalLoanDeductions,
		},
		"employerContributions": map[string]interface{}{
			"nssfEmployer": run.TotalNSSFEmployer,
			"nhifEmployer": run.TotalNHIFEmployer,
			"sdl":          run.TotalSDL,
			"wcf":          run.TotalWCF,
			"total":        run.TotalEmployerContributions,
		},
		"createdBy":        run.CreatedByName,
		"approvedBy":       run.ApprovedByName,
		"approvedAt":       run.ApprovedAt,
		"finalizedBy":      run.FinalizedByName,
		"finalizedAt":      run.FinalizedAt,
		"isLocked":         run.IsLocked,
		"cutoffDate":       run.CutoffDate,
		"disbursementDate": run.DisbursementDate,
	}

	return summary, nil
}
