package services

import (
	"errors"
	"math"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/repositories"
)

// LoanService handles business logic for loans
type LoanService struct {
	repo *repositories.LoanRepository
}

// NewLoanService creates a new service instance
func NewLoanService() *LoanService {
	return &LoanService{
		repo: repositories.NewLoanRepository(),
	}
}

// CreateLoan creates a new loan application
func (s *LoanService) CreateLoan(loan *models.Loan) error {
	// Calculate EMI if not provided
	if loan.EMIAmount == 0 {
		loan.EMIAmount = s.CalculateEMI(loan.Amount, loan.InterestRate, loan.Tenure)
	}

	// Set initial values
	loan.Status = models.LoanStatusPendingApproval
	loan.OutstandingBalance = loan.Amount + (loan.Amount * loan.InterestRate / 100)
	loan.PaymentsMade = 0
	loan.PaymentsRemaining = loan.Tenure
	loan.TotalPaid = 0

	return s.repo.Create(loan)
}

// GetLoan retrieves a loan by ID
func (s *LoanService) GetLoan(id uint, tenantID *uint) (*models.Loan, error) {
	return s.repo.GetByID(id, tenantID)
}

// UpdateLoan updates a loan
func (s *LoanService) UpdateLoan(loan *models.Loan) error {
	return s.repo.Update(loan)
}

// ListLoans lists loans with filters
func (s *LoanService) ListLoans(tenantID *uint, employeeID *uint, status, loanType string, page, pageSize int) ([]models.Loan, int64, error) {
	return s.repo.List(tenantID, employeeID, status, loanType, page, pageSize)
}

// ApproveLoan approves a loan application
func (s *LoanService) ApproveLoan(loanID uint, tenantID *uint, approverID uint, approverName string, disbursementDate *time.Time, remarks string) error {
	loan, err := s.repo.GetByID(loanID, tenantID)
	if err != nil {
		return err
	}

	if loan.Status != models.LoanStatusPendingApproval {
		return errors.New("loan is not pending approval")
	}

	now := time.Now()
	loan.Status = models.LoanStatusActive
	loan.ApprovedByID = &approverID
	loan.ApprovedByName = &approverName
	loan.ApprovedAt = &now

	if disbursementDate != nil {
		loan.DisbursedDate = disbursementDate
	} else {
		loan.DisbursedDate = &now
	}

	// Set start and end dates
	startDate := loan.DisbursedDate
	loan.StartDate = startDate
	endDate := startDate.AddDate(0, loan.Tenure, 0)
	loan.EndDate = &endDate

	// Set next EMI date (first day of next month)
	nextMonth := time.Date(startDate.Year(), startDate.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	loan.NextEMIDate = &nextMonth

	if remarks != "" {
		loan.Remarks = &remarks
	}

	if err := s.repo.Update(loan); err != nil {
		return err
	}

	// Generate repayment schedule
	return s.GenerateRepaymentSchedule(loan)
}

// RejectLoan rejects a loan application
func (s *LoanService) RejectLoan(loanID uint, tenantID *uint, rejecterID uint, rejecterName, reason string) error {
	loan, err := s.repo.GetByID(loanID, tenantID)
	if err != nil {
		return err
	}

	if loan.Status != models.LoanStatusPendingApproval {
		return errors.New("loan is not pending approval")
	}

	now := time.Now()
	loan.Status = models.LoanStatusRejected
	loan.RejectedByID = &rejecterID
	loan.RejectedByName = &rejecterName
	loan.RejectedAt = &now
	loan.RejectionReason = &reason

	return s.repo.Update(loan)
}

// CalculateEMI calculates Equated Monthly Installment
func (s *LoanService) CalculateEMI(principal, annualInterestRate float64, tenure int) float64 {
	if annualInterestRate == 0 {
		// Simple division for zero interest
		return math.Round(principal/float64(tenure)*100) / 100
	}

	// Monthly interest rate
	monthlyRate := annualInterestRate / 12 / 100

	// EMI formula: P * r * (1+r)^n / ((1+r)^n - 1)
	emi := principal * monthlyRate * math.Pow(1+monthlyRate, float64(tenure)) / (math.Pow(1+monthlyRate, float64(tenure)) - 1)

	return math.Round(emi*100) / 100
}

// GenerateRepaymentSchedule generates the repayment schedule for a loan
func (s *LoanService) GenerateRepaymentSchedule(loan *models.Loan) error {
	if loan.StartDate == nil {
		return errors.New("loan start date is required")
	}

	var repayments []models.LoanRepayment
	balance := loan.OutstandingBalance
	monthlyRate := loan.InterestRate / 12 / 100

	dueDate := *loan.StartDate
	for i := 1; i <= loan.Tenure; i++ {
		// Move to next month
		dueDate = dueDate.AddDate(0, 1, 0)

		var interest, principal float64
		if loan.InterestRate > 0 {
			interest = balance * monthlyRate
			principal = loan.EMIAmount - interest
		} else {
			interest = 0
			principal = loan.EMIAmount
		}

		repayment := models.LoanRepayment{
			LoanID:      loan.ID,
			Installment: i,
			DueDate:     dueDate,
			EMIAmount:   loan.EMIAmount,
			Principal:   math.Round(principal*100) / 100,
			Interest:    math.Round(interest*100) / 100,
			Status:      models.LoanRepaymentStatusPending,
		}

		repayments = append(repayments, repayment)
		balance -= principal
	}

	return s.repo.CreateRepaymentsBatch(repayments)
}

// GetRepaymentSchedule retrieves the repayment schedule for a loan
func (s *LoanService) GetRepaymentSchedule(loanID uint) ([]models.LoanRepayment, error) {
	return s.repo.GetRepaymentsByLoanID(loanID)
}

// RecordRepayment records a loan repayment
func (s *LoanService) RecordRepayment(loanID uint, installment int, paidAmount float64, payslipID *uint) error {
	repayments, err := s.repo.GetRepaymentsByLoanID(loanID)
	if err != nil {
		return err
	}

	var repayment *models.LoanRepayment
	for i := range repayments {
		if repayments[i].Installment == installment {
			repayment = &repayments[i]
			break
		}
	}

	if repayment == nil {
		return errors.New("installment not found")
	}

	if repayment.Status == models.LoanRepaymentStatusPaid {
		return errors.New("installment already paid")
	}

	now := time.Now()
	repayment.Status = models.LoanRepaymentStatusPaid
	repayment.PaidDate = &now
	repayment.PaidAmount = &paidAmount
	repayment.PayslipID = payslipID

	if err := s.repo.UpdateRepayment(repayment); err != nil {
		return err
	}

	// Update loan totals
	loan, err := s.repo.GetByID(loanID, nil)
	if err != nil {
		return err
	}

	loan.TotalPaid += paidAmount
	loan.OutstandingBalance -= paidAmount
	loan.PaymentsMade++
	loan.PaymentsRemaining--

	// Update next EMI date
	if loan.PaymentsRemaining > 0 {
		for _, rep := range repayments {
			if rep.Status == models.LoanRepaymentStatusPending {
				loan.NextEMIDate = &rep.DueDate
				break
			}
		}
	} else {
		loan.NextEMIDate = nil
		loan.Status = models.LoanStatusClosed
	}

	return s.repo.Update(loan)
}

// GetLoanSummary retrieves loan summary
func (s *LoanService) GetLoanSummary(tenantID *uint, employeeID *uint) (map[string]interface{}, error) {
	return s.repo.GetSummary(tenantID, employeeID)
}

// GetActiveLoansForEmployee retrieves active loans for an employee
func (s *LoanService) GetActiveLoansForEmployee(employeeID uint) ([]models.Loan, error) {
	return s.repo.GetActiveLoansForEmployee(employeeID)
}

// GetPendingRepayments retrieves pending repayments due by a date
func (s *LoanService) GetPendingRepayments(employeeID uint, dueDate time.Time) ([]models.LoanRepayment, error) {
	return s.repo.GetPendingRepayments(employeeID, dueDate)
}

// GetLoanWithRepayments retrieves a loan with its repayment schedule
func (s *LoanService) GetLoanWithRepayments(loanID uint, tenantID *uint) (map[string]interface{}, error) {
	loan, err := s.repo.GetByID(loanID, tenantID)
	if err != nil {
		return nil, err
	}

	repayments, err := s.repo.GetRepaymentsByLoanID(loanID)
	if err != nil {
		return nil, err
	}

	var schedule []map[string]interface{}
	for _, rep := range repayments {
		schedule = append(schedule, map[string]interface{}{
			"installment": rep.Installment,
			"dueDate":     rep.DueDate.Format("2006-01-02"),
			"emiAmount":   rep.EMIAmount,
			"principal":   rep.Principal,
			"interest":    rep.Interest,
			"status":      rep.Status,
			"paidDate":    rep.PaidDate,
			"paidAmount":  rep.PaidAmount,
		})
	}

	return map[string]interface{}{
		"id":                 loan.ID,
		"employeeId":         loan.EmployeeID,
		"empId":              loan.EmployeeCode,
		"employeeName":       loan.EmployeeName,
		"department":         loan.Department,
		"loanType":           loan.LoanType,
		"amount":             loan.Amount,
		"interestRate":       loan.InterestRate,
		"tenure":             loan.Tenure,
		"emiAmount":          loan.EMIAmount,
		"status":             loan.Status,
		"totalPaid":          loan.TotalPaid,
		"outstandingBalance": loan.OutstandingBalance,
		"paymentsMade":       loan.PaymentsMade,
		"paymentsRemaining":  loan.PaymentsRemaining,
		"nextEmiDate":        loan.NextEMIDate,
		"disbursedDate":      loan.DisbursedDate,
		"startDate":          loan.StartDate,
		"endDate":            loan.EndDate,
		"purpose":            loan.Purpose,
		"approvedBy":         loan.ApprovedByName,
		"approvedAt":         loan.ApprovedAt,
		"repaymentSchedule":  schedule,
	}, nil
}
