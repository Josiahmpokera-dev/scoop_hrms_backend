package models

import "time"

// Step 1: Personal Information Request
type Step1PersonalInfoRequest struct {
	// Profile Picture
	PhotoURL *string `json:"photo_url,omitempty"`

	// Basic Information
	FirstName    string     `json:"first_name" binding:"required,min=2"`
	MiddleName   *string    `json:"middle_name,omitempty"`
	LastName     string     `json:"last_name" binding:"required,min=2"`
	DateOfBirth  *time.Time `json:"date_of_birth,omitempty"`
	Gender       *string    `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	MaritalStatus *string   `json:"marital_status,omitempty" binding:"omitempty,oneof=single married divorced widowed"`
	BloodGroup   *string    `json:"blood_group,omitempty" binding:"omitempty,oneof=A+ A- B+ B- AB+ AB- O+ O-"`
	Nationality  *string    `json:"nationality,omitempty"`

	// Contact Information
	PersonalEmail    *string `json:"personal_email,omitempty" binding:"omitempty,email"`
	MobileNumber     *string `json:"mobile_number,omitempty"`
	AlternateNumber  *string `json:"alternate_number,omitempty"`
	CurrentAddress   *string `json:"current_address,omitempty"`
	PermanentAddress *string `json:"permanent_address,omitempty"`
	City             *string `json:"city,omitempty"`
	State            *string `json:"state,omitempty"`
	PostalCode       *string `json:"postal_code,omitempty"`
	Country          *string `json:"country,omitempty"`
}

// Step 2: Employment Details Request
type Step2EmploymentRequest struct {
	// Employment ID & Email
	EmployeeID  *string    `json:"employee_id,omitempty"` // Can be auto-generated if not provided
	OfficialEmail *string  `json:"official_email,omitempty" binding:"omitempty,email"`
	DateOfJoining *time.Time `json:"date_of_joining,omitempty"`

	// Organization Structure
	DepartmentID  *uint   `json:"department_id,omitempty"`
	PositionID    *uint   `json:"position_id,omitempty"` // Designation
	Grade         *string `json:"grade,omitempty"`
	ReportingManagerID *uint `json:"reporting_manager_id,omitempty"` // Manager ID
	EmploymentType *string `json:"employment_type,omitempty" binding:"omitempty,oneof=full_time part_time contract intern"`

	// Work Location & Shift
	LocationID *uint   `json:"location_id,omitempty"`
	Shift      *string `json:"shift,omitempty"` // Day, Night, Rotating, etc.
	WorkPhone  *string `json:"work_phone,omitempty"`

	// Probation Period
	ProbationPeriodDays *int       `json:"probation_period_days,omitempty" binding:"omitempty,min=90"` // Minimum 90 days
	ExpectedConfirmationDate *time.Time `json:"expected_confirmation_date,omitempty"`
}

// Step 3: Salary & CTC Request
type Step3SalaryRequest struct {
	// Cost to Company (CTC)
	AnnualCTC        *float64   `json:"annual_ctc,omitempty"`
	CTCEffectiveDate *time.Time `json:"ctc_effective_date,omitempty"`
	Currency         string     `json:"currency" binding:"omitempty,len=3"` // Default: TZS

	// Salary Components
	BasicSalary        *float64 `json:"basic_salary,omitempty"`
	HouseRentAllowance *float64 `json:"house_rent_allowance,omitempty"`
	TransportAllowance *float64 `json:"transport_allowance,omitempty"`
	SpecialAllowance   *float64 `json:"special_allowance,omitempty"`
	OtherAllowances    *float64 `json:"other_allowances,omitempty"`
	// Gross Salary will be calculated automatically

	// Deductions
	IncomeTax       *float64 `json:"income_tax,omitempty"`
	ProvidentFund   *float64 `json:"provident_fund,omitempty"`
	ProfessionalTax *float64 `json:"professional_tax,omitempty"`
	OtherDeductions *float64 `json:"other_deductions,omitempty"`
	// Net Monthly Salary will be calculated automatically
}

// Step 4: Bank Account Details Request
type Step4BankRequest struct {
	BankName         *string `json:"bank_name,omitempty"`
	AccountHolderName *string `json:"account_holder_name,omitempty"`
	AccountNumber    *string `json:"account_number,omitempty"`
	AccountType      *string `json:"account_type,omitempty" binding:"omitempty,oneof=savings current"`
	BranchName       *string `json:"branch_name,omitempty"`
	SWIFTCode        *string `json:"swift_code,omitempty"`
}

// Step 5: Statutory Requirements Request
type Step5StatutoryRequest struct {
	// Tanzania Statutory Requirements
	TINNumber  *string `json:"tin_number,omitempty"`  // Tax Identification Number
	NSSFNumber *string `json:"nssf_number,omitempty"` // National Social Security Fund
	NHIFNumber *string `json:"nhif_number,omitempty"` // National Health Insurance Fund
	WCFNumber  *string `json:"wcf_number,omitempty"`  // Workers Compensation Fund
	SDLNumber  *string `json:"sdl_number,omitempty"`  // Skills Development Levy

	// Work Authorization (For Expatriates)
	PassportNumber       *string    `json:"passport_number,omitempty"`
	PassportExpiryDate   *time.Time `json:"passport_expiry_date,omitempty"`
	WorkPermitNumber     *string    `json:"work_permit_number,omitempty"`
	WorkPermitExpiryDate *time.Time `json:"work_permit_expiry_date,omitempty"`
}

// Step 6: Documents Request (for file uploads, URLs will be provided)
type Step6DocumentsRequest struct {
	Documents []DocumentUpload `json:"documents,omitempty"`
}

// DocumentUpload represents a document to be uploaded
type DocumentUpload struct {
	DocumentType DocumentType `json:"document_type" binding:"required,oneof=identity work_permit education contract tax_statutory other"`
	FileName     string       `json:"file_name" binding:"required"`
	FileURL      string       `json:"file_url" binding:"required"` // URL to uploaded file
	FileSize     *int64       `json:"file_size,omitempty"`
	MimeType     *string      `json:"mime_type,omitempty"`
	Description  *string      `json:"description,omitempty"`
}

// Step 7: Company Assets Request (placeholder - will be handled by separate API)
type Step7AssetsRequest struct {
	AssetIDs []uint `json:"asset_ids,omitempty"` // IDs of assigned assets
	Note     *string `json:"note,omitempty"`    // Notes about asset assignment
}

// Step 8: Policies Request
type Step8PoliciesRequest struct {
	LeavePolicyID     *uint      `json:"leave_policy_id,omitempty"`
	AttendancePolicyID *uint     `json:"attendance_policy_id,omitempty"`
	WeeklyOffDays     *string    `json:"weekly_off_days,omitempty"` // e.g., "Saturday,Sunday"
	EffectiveDate     *time.Time `json:"effective_date,omitempty"`
}

// Step 9: Emergency Contacts Request
type Step9EmergencyContactsRequest struct {
	Contacts []EmergencyContactInput `json:"contacts,omitempty"`
}

// EmergencyContactInput represents an emergency contact input
type EmergencyContactInput struct {
	ContactName    string  `json:"contact_name" binding:"required"`
	Relationship   *string `json:"relationship,omitempty"`
	PhoneNumber    string  `json:"phone_number" binding:"required"`
	AlternatePhone *string `json:"alternate_phone,omitempty"`
	Email          *string `json:"email,omitempty" binding:"omitempty,email"`
	Address        *string `json:"address,omitempty"`
	IsPrimary      bool    `json:"is_primary"` // Primary emergency contact
}

// Step 10: Internal Notes Request
type Step10NotesRequest struct {
	Notes *string `json:"notes,omitempty"`
}

// CreateDraftRequest represents a request to create a new onboarding draft
// Optionally accepts step data to create draft and save a step in one call
type CreateDraftRequest struct {
	Step1Data map[string]interface{} `json:"step1_data,omitempty"` // Optional: Step 1 data to save immediately (legacy)
	Step      *int                   `json:"step,omitempty"`      // Optional: Step number (1-10) if providing step data
	Data      map[string]interface{} `json:"data,omitempty"`      // Optional: Step data (requires step number)
}

// SaveDraftRequest represents a request to save a draft for a specific step
type SaveDraftRequest struct {
	Step  int                    `json:"step" binding:"required,min=1,max=10"` // Step number (1-10)
	Data  map[string]interface{} `json:"data" binding:"required"`              // Step-specific data
}

// GetDraftResponse represents the draft data with progress
type GetDraftResponse struct {
	DraftID         uint                   `json:"draft_id"`
	EmployeeID      *string                `json:"employee_id,omitempty"`
	Progress        float64                `json:"progress"` // Completion percentage (0-100)
	CompletedSteps  []int                  `json:"completed_steps"`
	FinishedSteps   []int                  `json:"finished_steps"`   // Steps that have data saved
	UnfinishedSteps []int                  `json:"unfinished_steps"` // Steps that don't have data yet
	IsCompleted     bool                   `json:"is_completed"`
	Steps           map[string]interface{} `json:"steps"` // Step data keyed by step number
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// CompleteOnboardingRequest represents the final request to complete onboarding
type CompleteOnboardingRequest struct {
	DraftID uint `json:"draft_id" binding:"required"` // Draft ID to complete
}

// GetDraftEmployeeRequest represents the request to get a draft employee by employee ID
type GetDraftEmployeeRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"` // Employee ID (e.g., "EMP001")
}

// DraftEmployeeListItem represents a simplified draft employee for listing
type DraftEmployeeListItem struct {
	DraftID         uint                   `json:"draft_id"`
	EmployeeID      *string                `json:"employee_id,omitempty"`
	FirstName       *string                `json:"first_name,omitempty"`
	LastName        *string                `json:"last_name,omitempty"`
	FullName        *string                `json:"full_name,omitempty"`
	PersonalEmail   *string                `json:"personal_email,omitempty"`
	MobileNumber    *string                `json:"mobile_number,omitempty"`
	DepartmentName  *string                `json:"department_name,omitempty"`
	PositionName    *string                `json:"position_name,omitempty"`
	Progress        float64                `json:"progress"`
	CompletedSteps  []int                  `json:"completed_steps"`
	FinishedSteps   []int                  `json:"finished_steps"`
	UnfinishedSteps []int                  `json:"unfinished_steps"`
	IsCompleted     bool                   `json:"is_completed"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// CompletedEmployeeOnboardingResponse represents the completed employee onboarding data
type CompletedEmployeeOnboardingResponse struct {
	EmployeeID      string                 `json:"employee_id"`
	EmployeeDBID    uint                   `json:"employee_db_id"` // Database ID
	Progress        float64                `json:"progress"` // Always 100 for completed
	CompletedSteps  []int                  `json:"completed_steps"` // All steps 1-10
	IsCompleted     bool                   `json:"is_completed"` // Always true
	Steps           map[string]interface{} `json:"steps"` // Step data keyed by step number
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}
