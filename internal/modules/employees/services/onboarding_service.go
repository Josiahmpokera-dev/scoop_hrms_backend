package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	departmentRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	locationRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/repositories"
	positionRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/repositories"
	roleRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/repositories"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	userRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/email"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type OnboardingService struct {
	draftRepo             *employeeRepos.OnboardingDraftRepository
	employeeRepo          *employeeRepos.EmployeeRepository
	basicInfoRepo         *employeeRepos.EmployeeBasicInformationRepository
	employmentDetailsRepo *employeeRepos.EmployeeEmploymentDetailsRepository
	addressRepo           *employeeRepos.EmployeeAddressRepository
	salaryRepo            *employeeRepos.EmployeeSalaryRepository
	bankRepo              *employeeRepos.EmployeeBankRepository
	statutoryRepo         *employeeRepos.EmployeeStatutoryRepository
	documentRepo          *employeeRepos.EmployeeDocumentRepository
	contactRepo           *employeeRepos.EmployeeEmergencyContactRepository
	policyRepo            *employeeRepos.EmployeePolicyRepository
	assetRepo             *employeeRepos.EmployeeAssetRepository
	departmentRepo        *departmentRepos.DepartmentRepository
	positionRepo          *positionRepos.JobPositionRepository
	locationRepo          *locationRepos.LocationRepository
	userRepo              *userRepos.UserRepository
	roleRepo              *roleRepos.RoleRepository
}

func NewOnboardingService() *OnboardingService {
	return &OnboardingService{
		draftRepo:             employeeRepos.NewOnboardingDraftRepository(),
		employeeRepo:          employeeRepos.NewEmployeeRepository(),
		basicInfoRepo:         employeeRepos.NewEmployeeBasicInformationRepository(),
		employmentDetailsRepo: employeeRepos.NewEmployeeEmploymentDetailsRepository(),
		addressRepo:           employeeRepos.NewEmployeeAddressRepository(),
		salaryRepo:            employeeRepos.NewEmployeeSalaryRepository(),
		bankRepo:              employeeRepos.NewEmployeeBankRepository(),
		statutoryRepo:         employeeRepos.NewEmployeeStatutoryRepository(),
		documentRepo:          employeeRepos.NewEmployeeDocumentRepository(),
		contactRepo:           employeeRepos.NewEmployeeEmergencyContactRepository(),
		policyRepo:            employeeRepos.NewEmployeePolicyRepository(),
		assetRepo:             employeeRepos.NewEmployeeAssetRepository(),
		departmentRepo:        departmentRepos.NewDepartmentRepository(),
		positionRepo:          positionRepos.NewJobPositionRepository(),
		locationRepo:          locationRepos.NewLocationRepository(),
		userRepo:              userRepos.NewUserRepository(),
		roleRepo:              roleRepos.NewRoleRepository(),
	}
}

// CreateDraftOpts optional params for creating a draft (e.g. for existing user onboarding)
type CreateDraftOpts struct {
	LinkedUserID      *uint   // Onboard this existing user; no credentials at end
	ExistingUserEmail *string // Resolve to user and set LinkedUserID
}

// CreateDraft creates a new onboarding draft. Pass opts with LinkedUserID or ExistingUserEmail to onboard an existing user (no credentials at end).
func (s *OnboardingService) CreateDraft(tenantID *uint, createdBy *uint, opts *CreateDraftOpts) (*models.EmployeeOnboardingDraft, error) {
	draft := &models.EmployeeOnboardingDraft{
		CompletedSteps: "[]",
		Progress:       0,
		IsCompleted:    false,
		CreatedBy:      createdBy,
		UpdatedBy:      createdBy,
	}
	if opts != nil {
		if opts.LinkedUserID != nil {
			existingEmp, _ := s.employeeRepo.FindByUserID(*opts.LinkedUserID)
			if existingEmp != nil {
				return nil, fmt.Errorf("user is already onboarded as an employee (employee_id: %s)", existingEmp.EmployeeID)
			}
			draft.LinkedUserID = opts.LinkedUserID
		} else if opts.ExistingUserEmail != nil && *opts.ExistingUserEmail != "" {
			user, err := s.userRepo.FindByEmail(*opts.ExistingUserEmail)
			if err != nil || user == nil {
				return nil, fmt.Errorf("no user found with email %s", *opts.ExistingUserEmail)
			}
			existingEmp, _ := s.employeeRepo.FindByUserID(user.ID)
			if existingEmp != nil {
				return nil, fmt.Errorf("user is already onboarded as an employee")
			}
			draft.LinkedUserID = &user.ID
		}
	}
	if err := s.draftRepo.Create(draft); err != nil {
		return nil, fmt.Errorf("failed to create draft: %w", err)
	}
	return draft, nil
}

// SaveStep saves data for a specific onboarding step
func (s *OnboardingService) SaveStep(draftID uint, step int, data map[string]interface{}, updatedBy *uint) (*models.EmployeeOnboardingDraft, error) {
	draft, err := s.draftRepo.FindByID(draftID)
	if err != nil {
		return nil, errors.New("draft not found")
	}

	// Update draft
	draft.UpdatedBy = updatedBy

	// Ensure employee_id is set and normalized (should be set in step 1, but handle edge cases)
	if draft.EmployeeID == nil || *draft.EmployeeID == "" {
		if step == 1 {
			employeeID := s.generateEmployeeID()
			draft.EmployeeID = &employeeID
		} else {
			// For steps > 1, employee_id should already be set
			// If not, this is an error condition
			return nil, errors.New("employee_id is not set in draft. Please save step 1 first")
		}
	} else {
		// Normalize employee_id (uppercase, trimmed) to ensure consistency
		normalizedID := strings.ToUpper(strings.TrimSpace(*draft.EmployeeID))
		if *draft.EmployeeID != normalizedID {
			draft.EmployeeID = &normalizedID
			// Save the normalized employee_id immediately
			if err := s.draftRepo.Update(draft); err != nil {
				return nil, fmt.Errorf("failed to normalize employee_id: %w", err)
			}
		}
	}

	// Save step-specific data
	switch step {
	case 1:
		err = s.saveStep1PersonalInfo(draft, data)
	case 2:
		err = s.saveStep2Employment(draft, data)
	case 3:
		err = s.saveStep3Salary(draft, data)
	case 4:
		err = s.saveStep4Bank(draft, data)
	case 5:
		err = s.saveStep5Statutory(draft, data)
	case 6:
		err = s.saveStep6Documents(draft, data)
	case 7:
		err = s.saveStep7Assets(draft, data) // Placeholder - assets API will be separate
	case 8:
		err = s.saveStep8Policies(draft, data)
	case 9:
		err = s.saveStep9EmergencyContacts(draft, data)
	case 10:
		err = s.saveStep10Notes(draft, data)
	default:
		return nil, errors.New("invalid step number")
	}

	if err != nil {
		return nil, err
	}

	// Mark step as completed (this updates CompletedSteps)
	draft.AddCompletedStep(models.OnboardingStep(step))

	// Calculate progress AFTER updating completed steps
	draft.Progress = draft.CalculateProgress()
	draft.UpdatedBy = updatedBy

	// Save to database - ensure all fields are updated including Progress and CompletedSteps
	if err := s.draftRepo.Update(draft); err != nil {
		return nil, fmt.Errorf("failed to update draft: %w", err)
	}

	// Reload draft from database to ensure we have the latest saved data including progress
	var updatedDraft *models.EmployeeOnboardingDraft
	if draft.EmployeeID != nil {
		updatedDraft, err = s.draftRepo.FindByEmployeeIDString(*draft.EmployeeID, nil)
	} else {
		updatedDraft, err = s.draftRepo.FindByID(draft.ID)
	}
	if err != nil {
		// If reload fails, return the draft we have (should still have progress calculated)
		return draft, nil
	}

	return updatedDraft, nil
}

// SaveStepByEmployeeID saves data for a specific onboarding step by employee ID
func (s *OnboardingService) SaveStepByEmployeeID(employeeID string, tenantID *uint, step int, data map[string]interface{}, updatedBy *uint) (*models.EmployeeOnboardingDraft, error) {
	// Trim whitespace from employee ID and normalize (uppercase)
	employeeID = strings.TrimSpace(employeeID)
	employeeID = strings.ToUpper(employeeID)
	if employeeID == "" {
		return nil, errors.New("employee ID cannot be empty")
	}

	// Try to find existing draft by employee ID (case-insensitive search)
	draft, err := s.draftRepo.FindByEmployeeIDString(employeeID, tenantID)
	if err != nil {
		// If draft doesn't exist and this is step 1, create a new draft
		if step == 1 && errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new draft with the provided employee ID
			draft = &models.EmployeeOnboardingDraft{
				EmployeeID:     &employeeID,
				CompletedSteps: "[]",
				Progress:       0,
				IsCompleted:    false,
				CreatedBy:      updatedBy,
				UpdatedBy:      updatedBy,
			}
			if err := s.draftRepo.Create(draft); err != nil {
				return nil, fmt.Errorf("failed to create draft: %w", err)
			}
		} else {
			// For other steps, try to find by case-insensitive search or create if step 1 was saved but employee_id wasn't properly set
			// This handles edge cases where step 1 was saved but employee_id wasn't persisted
			if step > 1 {
				// Try to find any draft for this tenant that might have this employee_id in a different case
				// Or check if there's a draft without employee_id that we can assign
				allDrafts, _ := s.draftRepo.FindIncompleteDrafts(tenantID)
				for _, d := range allDrafts {
					if d.EmployeeID != nil && strings.EqualFold(strings.TrimSpace(*d.EmployeeID), employeeID) {
						// Found a draft with matching employee_id (case-insensitive)
						draft = &d
						// Normalize the employee_id in the draft
						normalizedID := strings.ToUpper(strings.TrimSpace(*draft.EmployeeID))
						if *draft.EmployeeID != normalizedID {
							draft.EmployeeID = &normalizedID
							if err := s.draftRepo.Update(draft); err != nil {
								return nil, fmt.Errorf("failed to normalize employee ID: %w", err)
							}
						}
						break
					}
				}
				// If still not found, return error
				if draft == nil || draft.ID == 0 {
					return nil, fmt.Errorf("draft not found for employee ID '%s'. Please save step 1 first", employeeID)
				}
			} else {
				return nil, fmt.Errorf("draft not found for employee ID '%s'. Please save step 1 first", employeeID)
			}
		}
	}

	// Ensure the draft has the employee ID set and normalized (in case it was created without it or with wrong case)
	normalizedID := strings.ToUpper(strings.TrimSpace(employeeID))
	if draft.EmployeeID == nil || *draft.EmployeeID == "" || !strings.EqualFold(strings.TrimSpace(*draft.EmployeeID), normalizedID) {
		draft.EmployeeID = &normalizedID
		if err := s.draftRepo.Update(draft); err != nil {
			return nil, fmt.Errorf("failed to update draft with employee ID: %w", err)
		}
	}

	return s.SaveStep(draft.ID, step, data, updatedBy)
}

func (s *OnboardingService) EditEmployeeStep(employeeID string, tenantID *uint, step int, data map[string]interface{}, updatedBy *uint) (*models.EmployeeEditResponse, error) {
	_ = tenantID

	employeeID = strings.TrimSpace(employeeID)
	employeeID = strings.ToUpper(employeeID)
	if employeeID == "" {
		return nil, errors.New("employee ID cannot be empty")
	}

	employee, err := s.employeeRepo.FindByEmployeeID(employeeID)
	if err != nil || employee == nil {
		return nil, errors.New("employee not found")
	}

	updatedFields := make([]string, 0)

	switch step {
	case 1:
		err = s.editEmployeeStep1(employee, data, updatedBy, &updatedFields)
	case 2:
		err = s.editEmployeeStep2(employee, data, updatedBy, &updatedFields)
	case 3:
		err = s.editEmployeeStep3(employee, data, updatedBy, &updatedFields)
	case 4:
		err = s.editEmployeeStep4(employee, data, updatedBy, &updatedFields)
	case 5:
		err = s.editEmployeeStep5(employee, data, updatedBy, &updatedFields)
	case 6:
		err = s.editEmployeeStep6(employee, data, updatedBy, &updatedFields)
	case 7:
		err = s.editEmployeeStep7(employee, data, updatedBy, &updatedFields)
	case 8:
		err = s.editEmployeeStep8(employee, data, updatedBy, &updatedFields)
	case 9:
		err = s.editEmployeeStep9(employee, data, updatedBy, &updatedFields)
	case 10:
		err = s.editEmployeeStep10(employee, data, updatedBy, &updatedFields)
	default:
		return nil, errors.New("invalid step number")
	}

	if err != nil {
		return nil, err
	}

	return &models.EmployeeEditResponse{
		EmployeeID:    employee.EmployeeID,
		Step:          step,
		UpdatedFields: updatedFields,
		Message:       "Updated successfully",
		UpdatedAt:     employee.UpdatedAt,
	}, nil
}

func (s *OnboardingService) editEmployeeStep1(employee *models.Employee, data map[string]interface{}, updatedBy *uint, updatedFields *[]string) error {
	basicInfo, err := s.basicInfoRepo.FindByEmployeeID(employee.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basicInfo = &models.EmployeeBasicInformation{
				EmployeeID:       &employee.ID,
				EmployeeIDString: &employee.EmployeeID,
				FirstName:        employee.FirstName,
				MiddleName:       employee.MiddleName,
				LastName:         employee.LastName,
				DateOfBirth:      employee.DateOfBirth,
				Gender:           employee.Gender,
				MaritalStatus:    employee.MaritalStatus,
				BloodGroup:       employee.BloodGroup,
				Nationality:      employee.Nationality,
				PersonalEmail:    employee.PersonalEmail,
				MobileNumber:     employee.PhoneNumber,
				AlternateNumber:  employee.AlternatePhone,
				PhotoURL:         employee.PhotoURL,
			}
			if err := s.basicInfoRepo.Create(basicInfo); err != nil {
				return fmt.Errorf("failed to create basic information: %w", err)
			}
		} else {
			return fmt.Errorf("failed to load basic information: %w", err)
		}
	}

	setEmployeeStringField := func(key string, target **string, recordTarget **string) {
		if v, ok := data[key]; ok {
			if sVal, ok := v.(string); ok {
				*target = &sVal
				*recordTarget = &sVal
				*updatedFields = append(*updatedFields, key)
			}
		}
	}

	if v, ok := data["first_name"]; ok {
		if sVal, ok := v.(string); ok {
			employee.FirstName = sVal
			basicInfo.FirstName = sVal
			*updatedFields = append(*updatedFields, "first_name")
		}
	}
	if v, ok := data["last_name"]; ok {
		if sVal, ok := v.(string); ok {
			employee.LastName = sVal
			basicInfo.LastName = sVal
			*updatedFields = append(*updatedFields, "last_name")
		}
	}
	if v, ok := data["middle_name"]; ok {
		if v == nil {
			employee.MiddleName = nil
			basicInfo.MiddleName = nil
			*updatedFields = append(*updatedFields, "middle_name")
		} else if sVal, ok := v.(string); ok {
			employee.MiddleName = &sVal
			basicInfo.MiddleName = &sVal
			*updatedFields = append(*updatedFields, "middle_name")
		}
	}
	if v, ok := data["photo_url"]; ok {
		if v == nil {
			employee.PhotoURL = nil
			basicInfo.PhotoURL = nil
			*updatedFields = append(*updatedFields, "photo_url")
		} else if sVal, ok := v.(string); ok {
			employee.PhotoURL = &sVal
			basicInfo.PhotoURL = &sVal
			*updatedFields = append(*updatedFields, "photo_url")
		}
	}
	if v, ok := data["date_of_birth"]; ok {
		if v == nil {
			employee.DateOfBirth = nil
			basicInfo.DateOfBirth = nil
			*updatedFields = append(*updatedFields, "date_of_birth")
		} else if sVal, ok := v.(string); ok {
			if t, ok := parseTimeFlexible(sVal); ok {
				employee.DateOfBirth = &t
				basicInfo.DateOfBirth = &t
				*updatedFields = append(*updatedFields, "date_of_birth")
			}
		}
	}

	setEmployeeStringField("gender", &employee.Gender, &basicInfo.Gender)
	setEmployeeStringField("marital_status", &employee.MaritalStatus, &basicInfo.MaritalStatus)
	setEmployeeStringField("blood_group", &employee.BloodGroup, &basicInfo.BloodGroup)
	setEmployeeStringField("nationality", &employee.Nationality, &basicInfo.Nationality)
	setEmployeeStringField("personal_email", &employee.PersonalEmail, &basicInfo.PersonalEmail)
	setEmployeeStringField("mobile_number", &employee.PhoneNumber, &basicInfo.MobileNumber)
	setEmployeeStringField("alternate_number", &employee.AlternatePhone, &basicInfo.AlternateNumber)

	if err := s.basicInfoRepo.Update(basicInfo); err != nil {
		return fmt.Errorf("failed to update basic information: %w", err)
	}

	employee.UpdatedBy = updatedBy
	if err := s.employeeRepo.Update(employee); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	addresses, err := s.addressRepo.FindByEmployeeID(employee.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to load addresses: %w", err)
	}

	var currentAddr *models.EmployeeAddress
	var permanentAddr *models.EmployeeAddress
	for i := range addresses {
		addr := &addresses[i]
		if addr.AddressType == models.AddressTypeCurrent {
			currentAddr = addr
		}
		if addr.AddressType == models.AddressTypePermanent {
			permanentAddr = addr
		}
	}

	getString := func(key string) (*string, bool) {
		v, ok := data[key]
		if !ok {
			return nil, false
		}
		if v == nil {
			empty := ""
			return &empty, true
		}
		sVal, ok := v.(string)
		if !ok {
			return nil, false
		}
		return &sVal, true
	}

	addrLine, hasCurrentAddress := getString("current_address")
	_, hasCity := data["city"]
	_, hasState := data["state"]
	_, hasPostalCode := data["postal_code"]
	_, hasCountry := data["country"]
	if hasCurrentAddress || hasCity || hasState || hasPostalCode || hasCountry {
		if currentAddr == nil {
			currentAddr = &models.EmployeeAddress{
				EmployeeID:       &employee.ID,
				EmployeeIDString: &employee.EmployeeID,
				AddressType:      models.AddressTypeCurrent,
			}
		}
		if hasCurrentAddress && addrLine != nil {
			if *addrLine == "" {
				currentAddr.AddressLine1 = nil
			} else {
				currentAddr.AddressLine1 = addrLine
			}
			*updatedFields = append(*updatedFields, "current_address")
		}
		if city, ok := getString("city"); ok {
			if *city == "" {
				currentAddr.City = nil
			} else {
				currentAddr.City = city
			}
			*updatedFields = append(*updatedFields, "city")
		}
		if state, ok := getString("state"); ok {
			if *state == "" {
				currentAddr.State = nil
			} else {
				currentAddr.State = state
			}
			*updatedFields = append(*updatedFields, "state")
		}
		if pc, ok := getString("postal_code"); ok {
			if *pc == "" {
				currentAddr.PostalCode = nil
			} else {
				currentAddr.PostalCode = pc
			}
			*updatedFields = append(*updatedFields, "postal_code")
		}
		if country, ok := getString("country"); ok {
			if *country == "" {
				currentAddr.Country = nil
			} else {
				currentAddr.Country = country
			}
			*updatedFields = append(*updatedFields, "country")
		}
		if currentAddr.ID == 0 {
			if err := s.addressRepo.Create(currentAddr); err != nil {
				return fmt.Errorf("failed to save current address: %w", err)
			}
		} else if err := s.addressRepo.Update(currentAddr); err != nil {
			return fmt.Errorf("failed to update current address: %w", err)
		}
	}

	if permLine, ok := getString("permanent_address"); ok {
		if permanentAddr == nil {
			permanentAddr = &models.EmployeeAddress{
				EmployeeID:       &employee.ID,
				EmployeeIDString: &employee.EmployeeID,
				AddressType:      models.AddressTypePermanent,
			}
		}
		if *permLine == "" {
			permanentAddr.AddressLine1 = nil
		} else {
			permanentAddr.AddressLine1 = permLine
		}
		*updatedFields = append(*updatedFields, "permanent_address")
		if city, ok := getString("city"); ok {
			if *city == "" {
				permanentAddr.City = nil
			} else {
				permanentAddr.City = city
			}
		}
		if state, ok := getString("state"); ok {
			if *state == "" {
				permanentAddr.State = nil
			} else {
				permanentAddr.State = state
			}
		}
		if pc, ok := getString("postal_code"); ok {
			if *pc == "" {
				permanentAddr.PostalCode = nil
			} else {
				permanentAddr.PostalCode = pc
			}
		}
		if country, ok := getString("country"); ok {
			if *country == "" {
				permanentAddr.Country = nil
			} else {
				permanentAddr.Country = country
			}
		}
		if permanentAddr.ID == 0 {
			if err := s.addressRepo.Create(permanentAddr); err != nil {
				return fmt.Errorf("failed to save permanent address: %w", err)
			}
		} else if err := s.addressRepo.Update(permanentAddr); err != nil {
			return fmt.Errorf("failed to update permanent address: %w", err)
		}
	}

	return nil
}

func (s *OnboardingService) editEmployeeStep2(employee *models.Employee, data map[string]interface{}, updatedBy *uint, updatedFields *[]string) error {
	details, err := s.employmentDetailsRepo.FindByEmployeeID(employee.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			details = &models.EmployeeEmploymentDetails{
				EmployeeID:       &employee.ID,
				EmployeeIDString: &employee.EmployeeID,
			}
			if err := s.employmentDetailsRepo.Create(details); err != nil {
				return fmt.Errorf("failed to create employment details: %w", err)
			}
		} else {
			return fmt.Errorf("failed to load employment details: %w", err)
		}
	}

	if v, ok := data["official_email"]; ok {
		if v == nil {
			details.OfficialEmail = nil
			employee.WorkEmail = nil
			*updatedFields = append(*updatedFields, "official_email")
		} else if sVal, ok := v.(string); ok {
			details.OfficialEmail = &sVal
			employee.WorkEmail = &sVal
			*updatedFields = append(*updatedFields, "official_email")
		}
	}
	if v, ok := data["date_of_joining"]; ok {
		if v == nil {
			details.DateOfJoining = nil
			employee.HireDate = nil
			*updatedFields = append(*updatedFields, "date_of_joining")
		} else if sVal, ok := v.(string); ok {
			if t, ok := parseTimeFlexible(sVal); ok {
				details.DateOfJoining = &t
				employee.HireDate = &t
				*updatedFields = append(*updatedFields, "date_of_joining")
			}
		}
	}

	if id, ok := asUintPtr(data["department_id"]); ok {
		details.DepartmentID = id
		employee.DepartmentID = id
		*updatedFields = append(*updatedFields, "department_id")
	}
	if id, ok := asUintPtr(data["position_id"]); ok {
		details.PositionID = id
		employee.PositionID = id
		*updatedFields = append(*updatedFields, "position_id")
	}
	if id, ok := asUintPtr(data["reporting_manager_id"]); ok {
		details.ReportingManagerID = id
		employee.ReportsToID = id
		*updatedFields = append(*updatedFields, "reporting_manager_id")
	}
	if id, ok := asUintPtr(data["location_id"]); ok {
		details.LocationID = id
		employee.LocationID = id
		*updatedFields = append(*updatedFields, "location_id")
	}

	if v, ok := data["grade"]; ok {
		if v == nil {
			details.Grade = nil
			employee.Grade = nil
			*updatedFields = append(*updatedFields, "grade")
		} else if sVal, ok := v.(string); ok {
			details.Grade = &sVal
			employee.Grade = &sVal
			*updatedFields = append(*updatedFields, "grade")
		}
	}
	if v, ok := data["employment_type"]; ok {
		if v == nil {
			details.EmploymentType = nil
			employee.EmploymentType = nil
			*updatedFields = append(*updatedFields, "employment_type")
		} else if sVal, ok := v.(string); ok {
			details.EmploymentType = &sVal
			employee.EmploymentType = &sVal
			*updatedFields = append(*updatedFields, "employment_type")
		}
	}
	if v, ok := data["shift"]; ok {
		if v == nil {
			details.Shift = nil
			employee.Shift = nil
			*updatedFields = append(*updatedFields, "shift")
		} else if sVal, ok := v.(string); ok {
			details.Shift = &sVal
			employee.Shift = &sVal
			*updatedFields = append(*updatedFields, "shift")
		}
	}
	if v, ok := data["work_phone"]; ok {
		if v == nil {
			details.WorkPhone = nil
			employee.WorkPhone = nil
			*updatedFields = append(*updatedFields, "work_phone")
		} else if sVal, ok := v.(string); ok {
			details.WorkPhone = &sVal
			employee.WorkPhone = &sVal
			*updatedFields = append(*updatedFields, "work_phone")
		}
	}
	if days, ok := asIntPtr(data["probation_period_days"]); ok {
		details.ProbationPeriodDays = days
		employee.ProbationPeriodDays = days
		*updatedFields = append(*updatedFields, "probation_period_days")
	}
	if v, ok := data["expected_confirmation_date"]; ok {
		if v == nil {
			details.ExpectedConfirmationDate = nil
			employee.ExpectedConfirmationDate = nil
			*updatedFields = append(*updatedFields, "expected_confirmation_date")
		} else if sVal, ok := v.(string); ok {
			if t, ok := parseTimeFlexible(sVal); ok {
				details.ExpectedConfirmationDate = &t
				employee.ExpectedConfirmationDate = &t
				*updatedFields = append(*updatedFields, "expected_confirmation_date")
			}
		}
	}

	if err := s.employmentDetailsRepo.Update(details); err != nil {
		return fmt.Errorf("failed to update employment details: %w", err)
	}

	employee.UpdatedBy = updatedBy
	if err := s.employeeRepo.Update(employee); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	return nil
}

func (s *OnboardingService) editEmployeeStep3(employee *models.Employee, data map[string]interface{}, updatedBy *uint, updatedFields *[]string) error {
	salary, err := s.salaryRepo.FindByEmployeeID(employee.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			salary = &models.EmployeeSalaryComponent{
				EmployeeID:       &employee.ID,
				EmployeeIDString: &employee.EmployeeID,
				Currency:         "TZS",
			}
			if err := s.salaryRepo.Create(salary); err != nil {
				return fmt.Errorf("failed to create salary: %w", err)
			}
		} else {
			return fmt.Errorf("failed to load salary: %w", err)
		}
	}

	if v, ok := data["annual_ctc"]; ok {
		if f, ok := asFloat64Ptr(v); ok {
			salary.AnnualCTC = f
			*updatedFields = append(*updatedFields, "annual_ctc")
		}
	}
	if v, ok := data["ctc_effective_date"]; ok {
		if v == nil {
			salary.CTCEffectiveDate = nil
			*updatedFields = append(*updatedFields, "ctc_effective_date")
		} else if sVal, ok := v.(string); ok {
			if t, ok := parseTimeFlexible(sVal); ok {
				salary.CTCEffectiveDate = &t
				*updatedFields = append(*updatedFields, "ctc_effective_date")
			}
		}
	}
	if v, ok := data["currency"]; ok {
		if sVal, ok := v.(string); ok && sVal != "" {
			salary.Currency = sVal
			employee.Currency = sVal
			*updatedFields = append(*updatedFields, "currency")
		}
	}

	updateFloat := func(key string, target **float64) {
		if v, ok := data[key]; ok {
			if f, ok := asFloat64Ptr(v); ok {
				*target = f
				*updatedFields = append(*updatedFields, key)
			}
		}
	}

	updateFloat("basic_salary", &salary.BasicSalary)
	updateFloat("house_rent_allowance", &salary.HouseRentAllowance)
	updateFloat("transport_allowance", &salary.TransportAllowance)
	updateFloat("special_allowance", &salary.SpecialAllowance)
	updateFloat("other_allowances", &salary.OtherAllowances)
	updateFloat("income_tax", &salary.IncomeTax)
	updateFloat("provident_fund", &salary.ProvidentFund)
	updateFloat("professional_tax", &salary.ProfessionalTax)
	updateFloat("other_deductions", &salary.OtherDeductions)

	salary.CalculateGrossSalary()
	salary.CalculateNetSalary()

	if err := s.salaryRepo.Update(salary); err != nil {
		return fmt.Errorf("failed to update salary: %w", err)
	}

	if salary.BasicSalary != nil {
		employee.Salary = salary.BasicSalary
		*updatedFields = append(*updatedFields, "salary")
	}
	employee.UpdatedBy = updatedBy
	if err := s.employeeRepo.Update(employee); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	return nil
}

func (s *OnboardingService) editEmployeeStep4(employee *models.Employee, data map[string]interface{}, updatedBy *uint, updatedFields *[]string) error {
	banks, err := s.bankRepo.FindByEmployeeID(employee.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			banks = nil
		} else {
			return fmt.Errorf("failed to load bank details: %w", err)
		}
	}

	var bank *models.EmployeeBankAccount
	for i := range banks {
		if banks[i].IsPrimary {
			bank = &banks[i]
			break
		}
	}
	if bank == nil && len(banks) > 0 {
		bank = &banks[0]
	}
	if bank == nil {
		bank = &models.EmployeeBankAccount{
			EmployeeID:       &employee.ID,
			EmployeeIDString: &employee.EmployeeID,
			IsPrimary:        true,
		}
		if err := s.bankRepo.Create(bank); err != nil {
			return fmt.Errorf("failed to create bank details: %w", err)
		}
	}

	setString := func(key string, target **string) {
		if v, ok := data[key]; ok {
			if v == nil {
				*target = nil
				*updatedFields = append(*updatedFields, key)
			} else if sVal, ok := v.(string); ok {
				*target = &sVal
				*updatedFields = append(*updatedFields, key)
			}
		}
	}

	setString("bank_name", &bank.BankName)
	setString("account_holder_name", &bank.AccountHolderName)
	setString("account_number", &bank.AccountNumber)
	setString("account_type", &bank.AccountType)
	setString("branch_name", &bank.BranchName)
	setString("swift_code", &bank.SWIFTCode)

	if err := s.bankRepo.Update(bank); err != nil {
		return fmt.Errorf("failed to update bank details: %w", err)
	}

	employee.UpdatedBy = updatedBy
	if err := s.employeeRepo.Update(employee); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	return nil
}

func (s *OnboardingService) editEmployeeStep5(employee *models.Employee, data map[string]interface{}, updatedBy *uint, updatedFields *[]string) error {
	stat, err := s.statutoryRepo.FindByEmployeeID(employee.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			stat = &models.EmployeeStatutoryInfo{
				EmployeeID:       &employee.ID,
				EmployeeIDString: &employee.EmployeeID,
			}
			if err := s.statutoryRepo.Create(stat); err != nil {
				return fmt.Errorf("failed to create statutory details: %w", err)
			}
		} else {
			return fmt.Errorf("failed to load statutory details: %w", err)
		}
	}

	setString := func(key string, target **string) {
		if v, ok := data[key]; ok {
			if v == nil {
				*target = nil
				*updatedFields = append(*updatedFields, key)
			} else if sVal, ok := v.(string); ok {
				*target = &sVal
				*updatedFields = append(*updatedFields, key)
			}
		}
	}

	setString("tin_number", &stat.TINNumber)
	setString("nssf_number", &stat.NSSFNumber)
	setString("nhif_number", &stat.NHIFNumber)
	setString("wcf_number", &stat.WCFNumber)
	setString("sdl_number", &stat.SDLNumber)
	setString("passport_number", &stat.PassportNumber)
	setString("work_permit_number", &stat.WorkPermitNumber)

	if v, ok := data["passport_expiry_date"]; ok {
		if v == nil {
			stat.PassportExpiryDate = nil
			*updatedFields = append(*updatedFields, "passport_expiry_date")
		} else if sVal, ok := v.(string); ok {
			if t, ok := parseTimeFlexible(sVal); ok {
				stat.PassportExpiryDate = &t
				*updatedFields = append(*updatedFields, "passport_expiry_date")
			}
		}
	}
	if v, ok := data["work_permit_expiry_date"]; ok {
		if v == nil {
			stat.WorkPermitExpiryDate = nil
			*updatedFields = append(*updatedFields, "work_permit_expiry_date")
		} else if sVal, ok := v.(string); ok {
			if t, ok := parseTimeFlexible(sVal); ok {
				stat.WorkPermitExpiryDate = &t
				*updatedFields = append(*updatedFields, "work_permit_expiry_date")
			}
		}
	}

	if err := s.statutoryRepo.Update(stat); err != nil {
		return fmt.Errorf("failed to update statutory details: %w", err)
	}

	employee.UpdatedBy = updatedBy
	if err := s.employeeRepo.Update(employee); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	return nil
}

func (s *OnboardingService) editEmployeeStep6(employee *models.Employee, data map[string]interface{}, updatedBy *uint, updatedFields *[]string) error {
	if docsRaw, ok := data["documents"]; ok {
		docs, ok := docsRaw.([]interface{})
		if !ok {
			return errors.New("documents must be an array")
		}

		if err := s.documentRepo.DeleteByEmployeeID(employee.ID); err != nil {
			return fmt.Errorf("failed to replace documents: %w", err)
		}

		for _, docData := range docs {
			docMap, ok := docData.(map[string]interface{})
			if !ok {
				continue
			}
			docTypeStr, _ := docMap["document_type"].(string)
			fileName, _ := docMap["file_name"].(string)
			fileURL, _ := docMap["file_url"].(string)

			if docTypeStr == "" || fileName == "" || fileURL == "" {
				continue
			}

			docType := models.DocumentType(docTypeStr)
			doc := &models.EmployeeDocument{
				EmployeeID:       &employee.ID,
				EmployeeIDString: &employee.EmployeeID,
				DocumentType:     docType,
				FileName:         fileName,
				FileURL:          fileURL,
			}
			if fs, ok := docMap["file_size"].(float64); ok {
				val := int64(fs)
				doc.FileSize = &val
			}
			if mt, ok := docMap["mime_type"].(string); ok {
				doc.MimeType = &mt
			}
			if desc, ok := docMap["description"].(string); ok {
				doc.Description = &desc
			}

			if err := s.documentRepo.Create(doc); err != nil {
				return fmt.Errorf("failed to save document: %w", err)
			}
		}

		*updatedFields = append(*updatedFields, "documents")
	}

	employee.UpdatedBy = updatedBy
	if err := s.employeeRepo.Update(employee); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}
	return nil
}

func (s *OnboardingService) editEmployeeStep7(employee *models.Employee, data map[string]interface{}, updatedBy *uint, updatedFields *[]string) error {
	assetsRaw, ok := data["assets"]
	if !ok {
		return nil
	}
	assets, ok := assetsRaw.([]interface{})
	if !ok {
		return errors.New("assets must be an array")
	}

	if err := s.assetRepo.DeleteByEmployeeID(employee.ID); err != nil {
		return fmt.Errorf("failed to replace assets: %w", err)
	}

	for _, assetData := range assets {
		assetMap, ok := assetData.(map[string]interface{})
		if !ok {
			continue
		}

		assetType, _ := assetMap["asset_type"].(string)
		assetName, _ := assetMap["asset_name"].(string)
		if assetType == "" || assetName == "" {
			continue
		}

		asset := &models.EmployeeAsset{
			EmployeeID:       &employee.ID,
			EmployeeIDString: &employee.EmployeeID,
			AssetType:        assetType,
			AssetName:        assetName,
		}

		if v, ok := assetMap["serial_number"].(string); ok {
			asset.SerialNumber = &v
		}
		if v, ok := assetMap["asset_tag"].(string); ok {
			asset.AssetTag = &v
		}
		if v, ok := assetMap["assigned_date"].(string); ok {
			if t, ok := parseTimeFlexible(v); ok {
				asset.AssignedDate = &t
			}
		}
		if v, ok := assetMap["expected_return_date"].(string); ok {
			if t, ok := parseTimeFlexible(v); ok {
				asset.ExpectedReturnDate = &t
			}
		}
		if v, ok := assetMap["condition"].(string); ok {
			asset.Condition = &v
		}
		if v, ok := assetMap["notes"].(string); ok {
			asset.Notes = &v
		}

		if err := s.assetRepo.Create(asset); err != nil {
			return fmt.Errorf("failed to save asset: %w", err)
		}
	}

	*updatedFields = append(*updatedFields, "assets")

	employee.UpdatedBy = updatedBy
	if err := s.employeeRepo.Update(employee); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}
	return nil
}

func (s *OnboardingService) editEmployeeStep8(employee *models.Employee, data map[string]interface{}, updatedBy *uint, updatedFields *[]string) error {
	policy, err := s.policyRepo.FindByEmployeeID(employee.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			policy = &models.EmployeePolicy{
				EmployeeID:       &employee.ID,
				EmployeeIDString: &employee.EmployeeID,
			}
			if err := s.policyRepo.Create(policy); err != nil {
				return fmt.Errorf("failed to create policy: %w", err)
			}
		} else {
			return fmt.Errorf("failed to load policy: %w", err)
		}
	}

	if id, ok := asUintPtr(data["leave_policy_id"]); ok {
		policy.LeavePolicyID = id
		*updatedFields = append(*updatedFields, "leave_policy_id")
	}
	if id, ok := asUintPtr(data["attendance_policy_id"]); ok {
		policy.AttendancePolicyID = id
		*updatedFields = append(*updatedFields, "attendance_policy_id")
	}
	if v, ok := data["weekly_off_days"]; ok {
		if v == nil {
			policy.WeeklyOffDays = nil
			*updatedFields = append(*updatedFields, "weekly_off_days")
		} else if sVal, ok := v.(string); ok {
			policy.WeeklyOffDays = &sVal
			*updatedFields = append(*updatedFields, "weekly_off_days")
		}
	}
	if v, ok := data["effective_date"]; ok {
		if v == nil {
			policy.EffectiveDate = nil
			*updatedFields = append(*updatedFields, "effective_date")
		} else if sVal, ok := v.(string); ok {
			if t, ok := parseTimeFlexible(sVal); ok {
				policy.EffectiveDate = &t
				*updatedFields = append(*updatedFields, "effective_date")
			}
		}
	}

	if err := s.policyRepo.Update(policy); err != nil {
		return fmt.Errorf("failed to update policy: %w", err)
	}

	employee.UpdatedBy = updatedBy
	if err := s.employeeRepo.Update(employee); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}
	return nil
}

func (s *OnboardingService) editEmployeeStep9(employee *models.Employee, data map[string]interface{}, updatedBy *uint, updatedFields *[]string) error {
	contactsRaw, ok := data["contacts"]
	if !ok {
		return nil
	}
	contacts, ok := contactsRaw.([]interface{})
	if !ok {
		return errors.New("contacts must be an array")
	}

	if err := s.contactRepo.DeleteByEmployeeID(employee.ID); err != nil {
		return fmt.Errorf("failed to replace contacts: %w", err)
	}

	for _, contactData := range contacts {
		contactMap, ok := contactData.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := contactMap["contact_name"].(string)
		phone, _ := contactMap["phone_number"].(string)
		if name == "" || phone == "" {
			continue
		}
		contact := &models.EmployeeEmergencyContact{
			EmployeeID:       &employee.ID,
			EmployeeIDString: &employee.EmployeeID,
			ContactName:      name,
			PhoneNumber:      phone,
		}
		if v, ok := contactMap["relationship"].(string); ok {
			contact.Relationship = &v
		}
		if v, ok := contactMap["alternate_phone"].(string); ok {
			contact.AlternatePhone = &v
		}
		if v, ok := contactMap["email"].(string); ok {
			contact.Email = &v
		}
		if v, ok := contactMap["address"].(string); ok {
			contact.Address = &v
		}
		if v, ok := contactMap["is_primary"].(bool); ok {
			contact.IsPrimary = v
		}
		if err := s.contactRepo.Create(contact); err != nil {
			return fmt.Errorf("failed to save contact: %w", err)
		}
	}

	*updatedFields = append(*updatedFields, "contacts")

	employee.UpdatedBy = updatedBy
	if err := s.employeeRepo.Update(employee); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}
	return nil
}

func (s *OnboardingService) editEmployeeStep10(employee *models.Employee, data map[string]interface{}, updatedBy *uint, updatedFields *[]string) error {
	if v, ok := data["notes"]; ok {
		if v == nil {
			employee.Notes = ""
			*updatedFields = append(*updatedFields, "notes")
		} else if sVal, ok := v.(string); ok {
			employee.Notes = sVal
			*updatedFields = append(*updatedFields, "notes")
		}
	}
	employee.UpdatedBy = updatedBy
	if err := s.employeeRepo.Update(employee); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}
	return nil
}

func asUintPtr(v interface{}) (*uint, bool) {
	if v == nil {
		return nil, true
	}
	switch t := v.(type) {
	case float64:
		u := uint(t)
		return &u, true
	case int:
		u := uint(t)
		return &u, true
	case uint:
		return &t, true
	default:
		return nil, false
	}
}

func asIntPtr(v interface{}) (*int, bool) {
	if v == nil {
		return nil, true
	}
	switch t := v.(type) {
	case float64:
		i := int(t)
		return &i, true
	case int:
		return &t, true
	default:
		return nil, false
	}
}

func asFloat64Ptr(v interface{}) (*float64, bool) {
	if v == nil {
		return nil, true
	}
	switch t := v.(type) {
	case float64:
		return &t, true
	case int:
		f := float64(t)
		return &f, true
	default:
		return nil, false
	}
}

func parseTimeFlexible(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, true
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, true
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, true
	}
	return time.Time{}, false
}

// GetDraft retrieves a draft with all step data by draft ID
func (s *OnboardingService) GetDraft(draftID uint) (*models.GetDraftResponse, error) {
	draft, err := s.draftRepo.FindByID(draftID)
	if err != nil {
		return nil, errors.New("draft not found")
	}

	response := &models.GetDraftResponse{
		DraftID:                  draft.ID,
		EmployeeID:               draft.EmployeeID,
		LinkedUserID:             draft.LinkedUserID,
		IsExistingUserOnboarding: draft.LinkedUserID != nil,
		Progress:                 draft.Progress,
		CompletedSteps:           draft.GetCompletedStepsList(),
		IsCompleted:              draft.IsCompleted,
		Steps:                    make(map[string]interface{}),
		CreatedAt:                draft.CreatedAt,
		UpdatedAt:                draft.UpdatedAt,
	}
	if draft.LinkedUserID != nil {
		linkedUser, _ := s.userRepo.FindByID(*draft.LinkedUserID)
		if linkedUser != nil {
			response.LinkedUserEmail = &linkedUser.Email
		}
	}

	// Load all step data
	s.loadStepData(draft, response)

	// Add finished and unfinished steps information
	allSteps := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	finishedSteps := make([]int, 0)
	unfinishedSteps := make([]int, 0)

	for _, step := range allSteps {
		if draft.HasStepCompleted(models.OnboardingStep(step)) {
			finishedSteps = append(finishedSteps, step)
		} else {
			unfinishedSteps = append(unfinishedSteps, step)
		}
	}

	// Add to response (we'll extend GetDraftResponse model)
	response.FinishedSteps = finishedSteps
	response.UnfinishedSteps = unfinishedSteps

	return response, nil
}

// GetDraftByEmployeeID retrieves a draft with all step data by employee ID
func (s *OnboardingService) GetDraftByEmployeeID(employeeID string, tenantID *uint) (*models.GetDraftResponse, error) {
	draft, err := s.draftRepo.FindByEmployeeIDString(employeeID, tenantID)
	if err != nil {
		return nil, errors.New("draft not found for this employee ID")
	}
	return s.GetDraft(draft.ID)
}

// GetDraftSummary returns a simplified draft response for create/update operations
func (s *OnboardingService) GetDraftSummary(draftID uint) (*models.EmployeeOnboardingDraft, error) {
	return s.draftRepo.FindByID(draftID)
}

// ListDrafts lists all drafts for a tenant
func (s *OnboardingService) ListDrafts(tenantID *uint) ([]models.EmployeeOnboardingDraft, error) {
	return s.draftRepo.FindByTenantID(tenantID)
}

// ListNonEmployeeUsers returns users who are not yet linked to any employee (for onboarding existing users).
// Use the returned user id as linked_user_id when creating an onboarding draft.
func (s *OnboardingService) ListNonEmployeeUsers(tenantID *uint, page, pageSize int, search string) ([]userModels.User, int64, error) {
	return s.userRepo.ListNonEmployeeUsers(tenantID, page, pageSize, search)
}

// ListDraftEmployees lists all incomplete draft employees with their filled details
func (s *OnboardingService) ListDraftEmployees(tenantID *uint) ([]models.DraftEmployeeListItem, error) {
	// Get incomplete drafts
	drafts, err := s.draftRepo.FindIncompleteDrafts(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch drafts: %w", err)
	}
	return s.buildDraftEmployeeList(drafts)
}

// ListDraftEmployeesWithPagination lists incomplete draft employees with pagination
func (s *OnboardingService) ListDraftEmployeesWithPagination(tenantID *uint, page, pageSize int) ([]models.DraftEmployeeListItem, int64, error) {
	// Get incomplete drafts with pagination
	drafts, total, err := s.draftRepo.FindIncompleteDraftsWithPagination(tenantID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch drafts: %w", err)
	}

	draftEmployees, err := s.buildDraftEmployeeList(drafts)
	if err != nil {
		return nil, 0, err
	}

	return draftEmployees, total, nil
}

// GetCompletedEmployeeOnboarding retrieves all onboarding data for a completed employee
func (s *OnboardingService) GetCompletedEmployeeOnboarding(employeeID string, tenantID *uint) (*models.CompletedEmployeeOnboardingResponse, error) {
	// Find employee by employee_id string
	employee, err := s.employeeRepo.FindByEmployeeID(employeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Tenant ownership check removed (single-tenant)

	response := &models.CompletedEmployeeOnboardingResponse{
		EmployeeID:     employee.EmployeeID,
		EmployeeDBID:   employee.ID,
		Progress:       100.0,                                // Always 100% for completed employees
		CompletedSteps: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, // All steps completed
		IsCompleted:    true,
		Steps:          make(map[string]interface{}),
		CreatedAt:      employee.CreatedAt,
		UpdatedAt:      employee.UpdatedAt,
	}

	// Load all step data from employee record and related tables
	s.loadCompletedEmployeeStepData(employee, response)

	return response, nil
}

// loadCompletedEmployeeStepData loads all step data for a completed employee
func (s *OnboardingService) loadCompletedEmployeeStepData(employee *models.Employee, response *models.CompletedEmployeeOnboardingResponse) {
	// Step 1: Personal Info & Addresses
	basicInfo, _ := s.basicInfoRepo.FindByEmployeeID(employee.ID)
	addresses, _ := s.addressRepo.FindByEmployeeID(employee.ID)

	step1Response := make(map[string]interface{})
	if basicInfo != nil {
		step1Response["photo_url"] = basicInfo.PhotoURL
		step1Response["first_name"] = basicInfo.FirstName
		step1Response["middle_name"] = basicInfo.MiddleName
		step1Response["last_name"] = basicInfo.LastName
		step1Response["date_of_birth"] = basicInfo.DateOfBirth
		step1Response["gender"] = basicInfo.Gender
		step1Response["marital_status"] = basicInfo.MaritalStatus
		step1Response["blood_group"] = basicInfo.BloodGroup
		step1Response["nationality"] = basicInfo.Nationality
		step1Response["personal_email"] = basicInfo.PersonalEmail
		step1Response["mobile_number"] = basicInfo.MobileNumber
		step1Response["alternate_number"] = basicInfo.AlternateNumber
	} else {
		// Fallback to employee table if basic info not found
		step1Response["photo_url"] = employee.PhotoURL
		step1Response["first_name"] = employee.FirstName
		step1Response["middle_name"] = employee.MiddleName
		step1Response["last_name"] = employee.LastName
		step1Response["date_of_birth"] = employee.DateOfBirth
		step1Response["gender"] = employee.Gender
		step1Response["marital_status"] = employee.MaritalStatus
		step1Response["blood_group"] = employee.BloodGroup
		step1Response["nationality"] = employee.Nationality
		step1Response["personal_email"] = employee.PersonalEmail
		step1Response["mobile_number"] = employee.PhoneNumber
		step1Response["alternate_number"] = employee.AlternatePhone
	}
	if len(addresses) > 0 {
		step1Response["addresses"] = addresses
	}

	if len(step1Response) > 0 {
		response.Steps["1"] = step1Response
	}

	// Step 2: Employment Details
	employmentDetails, _ := s.employmentDetailsRepo.FindByEmployeeID(employee.ID)
	if employmentDetails != nil {
		step2Response := make(map[string]interface{})
		step2Response["employee_id"] = employee.EmployeeID
		step2Response["official_email"] = employmentDetails.OfficialEmail
		step2Response["date_of_joining"] = employmentDetails.DateOfJoining
		step2Response["department_id"] = employmentDetails.DepartmentID
		step2Response["position_id"] = employmentDetails.PositionID
		step2Response["grade"] = employmentDetails.Grade
		step2Response["reporting_manager_id"] = employmentDetails.ReportingManagerID
		step2Response["employment_type"] = employmentDetails.EmploymentType
		step2Response["location_id"] = employmentDetails.LocationID
		step2Response["shift"] = employmentDetails.Shift
		step2Response["work_phone"] = employmentDetails.WorkPhone
		step2Response["probation_period_days"] = employmentDetails.ProbationPeriodDays
		step2Response["expected_confirmation_date"] = employmentDetails.ExpectedConfirmationDate
		response.Steps["2"] = step2Response
	} else {
		// Fallback to employee table
		step2Response := make(map[string]interface{})
		step2Response["employee_id"] = employee.EmployeeID
		step2Response["official_email"] = employee.WorkEmail
		step2Response["date_of_joining"] = employee.HireDate
		step2Response["department_id"] = employee.DepartmentID
		step2Response["position_id"] = employee.PositionID
		step2Response["grade"] = employee.Grade
		step2Response["reporting_manager_id"] = employee.ReportsToID
		step2Response["employment_type"] = employee.EmploymentType
		step2Response["location_id"] = employee.LocationID
		step2Response["shift"] = employee.Shift
		step2Response["work_phone"] = employee.WorkPhone
		step2Response["probation_period_days"] = employee.ProbationPeriodDays
		step2Response["expected_confirmation_date"] = employee.ExpectedConfirmationDate
		response.Steps["2"] = step2Response
	}

	// Step 3: Salary
	salary, _ := s.salaryRepo.FindByEmployeeID(employee.ID)
	if salary != nil {
		response.Steps["3"] = salary
	}

	// Step 4: Bank
	bankAccounts, _ := s.bankRepo.FindByEmployeeID(employee.ID)
	if len(bankAccounts) > 0 {
		response.Steps["4"] = bankAccounts
	}

	// Step 5: Statutory
	statutory, _ := s.statutoryRepo.FindByEmployeeID(employee.ID)
	if statutory != nil {
		response.Steps["5"] = statutory
	}

	// Step 6: Documents
	documents, _ := s.documentRepo.FindByEmployeeID(employee.ID)
	if len(documents) > 0 {
		response.Steps["6"] = documents
	}

	// Step 7: Assets
	assets, _ := s.assetRepo.FindByEmployeeID(employee.ID)
	if len(assets) > 0 {
		response.Steps["7"] = assets
	}

	// Step 8: Policies
	policy, _ := s.policyRepo.FindByEmployeeID(employee.ID)
	if policy != nil {
		response.Steps["8"] = policy
	}

	// Step 9: Emergency Contacts
	contacts, _ := s.contactRepo.FindByEmployeeID(employee.ID)
	if len(contacts) > 0 {
		response.Steps["9"] = contacts
	} else {
		// Fallback to employee table legacy fields
		if employee.EmergencyContactName != "" {
			response.Steps["9"] = []map[string]interface{}{
				{
					"contact_name": employee.EmergencyContactName,
					"phone_number": employee.EmergencyContactPhone,
					"relationship": employee.EmergencyContactRelation,
					"is_primary":   true,
				},
			}
		}
	}

	// Step 10: Notes
	if employee.Notes != "" {
		response.Steps["10"] = map[string]interface{}{
			"notes": employee.Notes,
		}
	}
}

// buildDraftEmployeeList builds the draft employee list from drafts
func (s *OnboardingService) buildDraftEmployeeList(drafts []models.EmployeeOnboardingDraft) ([]models.DraftEmployeeListItem, error) {

	var draftEmployees []models.DraftEmployeeListItem

	for _, draft := range drafts {
		item := models.DraftEmployeeListItem{
			DraftID:        draft.ID,
			EmployeeID:     draft.EmployeeID,
			Progress:       draft.Progress,
			CompletedSteps: draft.GetCompletedStepsList(),
			IsCompleted:    draft.IsCompleted,
			CreatedAt:      draft.CreatedAt,
			UpdatedAt:      draft.UpdatedAt,
		}

		// Calculate finished and unfinished steps
		allSteps := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		finishedSteps := make([]int, 0)
		unfinishedSteps := make([]int, 0)

		for _, step := range allSteps {
			if draft.HasStepCompleted(models.OnboardingStep(step)) {
				finishedSteps = append(finishedSteps, step)
			} else {
				unfinishedSteps = append(unfinishedSteps, step)
			}
		}

		item.FinishedSteps = finishedSteps
		item.UnfinishedSteps = unfinishedSteps

		// Load step 1 data (personal information)
		basicInfo, _ := s.basicInfoRepo.FindByDraftID(draft.ID)
		if basicInfo != nil {
			item.FirstName = &basicInfo.FirstName
			item.LastName = &basicInfo.LastName
			fullName := basicInfo.FirstName
			if basicInfo.MiddleName != nil && *basicInfo.MiddleName != "" {
				fullName += " " + *basicInfo.MiddleName
			}
			fullName += " " + basicInfo.LastName
			item.FullName = &fullName
			item.PersonalEmail = basicInfo.PersonalEmail
			item.MobileNumber = basicInfo.MobileNumber
		}

		// Load step 2 data (employment details) for department and position
		employmentDetails, _ := s.employmentDetailsRepo.FindByDraftID(draft.ID)
		if employmentDetails != nil {
			if employmentDetails.DepartmentID != nil {
				department, err := s.departmentRepo.FindByID(*employmentDetails.DepartmentID)
				if err == nil && department != nil {
					item.DepartmentName = &department.Name
				}
			}
			if employmentDetails.PositionID != nil {
				position, err := s.positionRepo.FindByID(*employmentDetails.PositionID)
				if err == nil && position != nil {
					item.PositionName = &position.Title
				}
			}
		}

		draftEmployees = append(draftEmployees, item)
	}

	return draftEmployees, nil
}

// CompleteOnboarding finalizes the onboarding and creates the employee by draft ID
func (s *OnboardingService) CompleteOnboarding(draftID uint, updatedBy *uint) (*models.Employee, *UserCredentials, error) {
	draft, err := s.draftRepo.FindByID(draftID)
	if err != nil {
		return nil, nil, errors.New("draft not found")
	}
	return s.completeOnboardingForDraft(draft, updatedBy)
}

// CompleteOnboardingByEmployeeID finalizes the onboarding and creates the employee by employee ID
func (s *OnboardingService) CompleteOnboardingByEmployeeID(employeeID string, tenantID *uint, updatedBy *uint) (*models.Employee, *UserCredentials, error) {
	draft, err := s.draftRepo.FindByEmployeeIDString(employeeID, tenantID)
	if err != nil {
		return nil, nil, errors.New("draft not found for this employee ID")
	}
	return s.completeOnboardingForDraft(draft, updatedBy)
}

// CompleteOnboardingResponse includes employee and credentials
type CompleteOnboardingResponse struct {
	Employee    *models.Employee `json:"employee"`
	Credentials *UserCredentials `json:"credentials,omitempty"`
}

// UserCredentials contains login credentials for the new user
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// completeOnboardingForDraft is the internal method that completes onboarding
func (s *OnboardingService) completeOnboardingForDraft(draft *models.EmployeeOnboardingDraft, updatedBy *uint) (*models.Employee, *UserCredentials, error) {

	if draft.IsCompleted {
		return nil, nil, errors.New("onboarding already completed")
	}

	// Validate required steps are completed
	if !s.validateRequiredSteps(draft) {
		return nil, nil, errors.New("required steps are not completed")
	}

	// Create employee from draft data
	employee, err := s.createEmployeeFromDraft(draft)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create employee: %w", err)
	}

	// Migrate all related data from draft to employee
	if err := s.migrateDraftDataToEmployee(draft.ID, employee.ID); err != nil {
		return nil, nil, fmt.Errorf("failed to migrate draft data: %w", err)
	}

	// Create user account for the employee
	// Check if the user performing onboarding is admin - if so, the employee might also be admin
	var isAdminEmployee bool
	if updatedBy != nil {
		existingUser, err := s.userRepo.FindByID(*updatedBy)
		if err == nil && existingUser != nil && existingUser.IsAdmin() {
			// The person completing onboarding is admin, check if employee should also be admin
			// For now, we'll check if user already exists and preserve their role
			// Or we can add a flag in the future to explicitly mark admin employees
			isAdminEmployee = false // Default to false, can be enhanced with explicit flag
		}
	}

	credentials, err := s.createUserForEmployee(employee, draft, updatedBy, isAdminEmployee)
	if err != nil {
		// Log error but don't fail onboarding - user can be created later
		// return nil, nil, fmt.Errorf("failed to create user account: %w", err)
		// For now, we'll continue even if user creation fails
		credentials = nil
	}

	// Send welcome email with credentials if credentials were created
	if credentials != nil {
		// Use employee's personal email for sending credentials (already available in employee object)
		if employee.PersonalEmail != nil && *employee.PersonalEmail != "" {
			// Send email with credentials to personal email
			go s.sendWelcomeEmailWithCredentials(*employee.PersonalEmail, employee, credentials)
		} else {
			// Fallback to work email if personal email is not available
			if employee.WorkEmail != nil && *employee.WorkEmail != "" {
				// Send email with credentials to work email
				go s.sendWelcomeEmailWithCredentials(*employee.WorkEmail, employee, credentials)
			}
		}
	}

	// Mark draft as completed
	draft.IsCompleted = true
	draft.EmployeeIDFinal = &employee.ID
	draft.Progress = 100.0
	draft.UpdatedBy = updatedBy

	if err := s.draftRepo.Update(draft); err != nil {
		return nil, nil, fmt.Errorf("failed to update draft: %w", err)
	}

	return employee, credentials, nil
}

// Helper methods for each step
func (s *OnboardingService) saveStep1PersonalInfo(draft *models.EmployeeOnboardingDraft, data map[string]interface{}) error {
	// Validate required fields
	if firstName, ok := data["first_name"].(string); !ok || firstName == "" {
		return fmt.Errorf("first_name is required")
	}
	if lastName, ok := data["last_name"].(string); !ok || lastName == "" {
		return fmt.Errorf("last_name is required")
	}

	// Generate or use provided employee_id if not already set
	if draft.EmployeeID == nil || *draft.EmployeeID == "" {
		var employeeID string
		// Check if employee_id is provided in the data
		if empID, ok := data["employee_id"].(string); ok && empID != "" {
			// Normalize employee_id (uppercase, trimmed)
			employeeID = strings.ToUpper(strings.TrimSpace(empID))
		} else {
			// Auto-generate employee ID
			employeeID = s.generateEmployeeID()
		}
		draft.EmployeeID = &employeeID
		// Update draft with employee_id immediately - use Save to ensure all fields are persisted
		if err := s.draftRepo.Update(draft); err != nil {
			return fmt.Errorf("failed to save employee_id: %w", err)
		}
	} else {
		// Normalize existing employee_id (uppercase, trimmed) to ensure consistency
		normalizedID := strings.ToUpper(strings.TrimSpace(*draft.EmployeeID))
		if *draft.EmployeeID != normalizedID {
			draft.EmployeeID = &normalizedID
			if err := s.draftRepo.Update(draft); err != nil {
				return fmt.Errorf("failed to normalize employee_id: %w", err)
			}
		}
	}

	// Parse Step1PersonalInfoRequest from data
	var req models.Step1PersonalInfoRequest
	jsonData, _ := json.Marshal(data)
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return fmt.Errorf("invalid step 1 data: %w", err)
	}

	// Check if basic information already exists for this draft
	existingInfo, err := s.basicInfoRepo.FindByDraftID(draft.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing basic information: %w", err)
	}

	// Create or update basic information
	basicInfo := &models.EmployeeBasicInformation{
		DraftID:          &draft.ID,
		EmployeeIDString: draft.EmployeeID, // Set employee_id_string from draft
		PhotoURL:         req.PhotoURL,
		FirstName:        req.FirstName,
		MiddleName:       req.MiddleName,
		LastName:         req.LastName,
		DateOfBirth:      req.DateOfBirth,
		Gender:           req.Gender,
		MaritalStatus:    req.MaritalStatus,
		BloodGroup:       req.BloodGroup,
		Nationality:      req.Nationality,
		PersonalEmail:    req.PersonalEmail,
		MobileNumber:     req.MobileNumber,
		AlternateNumber:  req.AlternateNumber,
	}

	if existingInfo != nil {
		// Update existing record
		basicInfo.ID = existingInfo.ID
		if err := s.basicInfoRepo.Update(basicInfo); err != nil {
			return fmt.Errorf("failed to update basic information: %w", err)
		}
	} else {
		// Create new record
		if err := s.basicInfoRepo.Create(basicInfo); err != nil {
			return fmt.Errorf("failed to save basic information: %w", err)
		}
	}

	// Save addresses (current and permanent)
	// Delete existing addresses for this draft
	s.addressRepo.DeleteByDraftID(draft.ID)

	// Save current address
	if req.CurrentAddress != nil || req.City != nil {
		currentAddr := &models.EmployeeAddress{
			DraftID:          &draft.ID,
			EmployeeIDString: draft.EmployeeID, // Set employee_id_string from draft
			AddressType:      models.AddressTypeCurrent,
			AddressLine1:     req.CurrentAddress,
			City:             req.City,
			State:            req.State,
			PostalCode:       req.PostalCode,
			Country:          req.Country,
		}
		if err := s.addressRepo.Create(currentAddr); err != nil {
			return fmt.Errorf("failed to save current address: %w", err)
		}
	}

	// Save permanent address
	if req.PermanentAddress != nil {
		permanentAddr := &models.EmployeeAddress{
			DraftID:          &draft.ID,
			EmployeeIDString: draft.EmployeeID, // Set employee_id_string from draft
			AddressType:      models.AddressTypePermanent,
			AddressLine1:     req.PermanentAddress,
			City:             req.City,
			State:            req.State,
			PostalCode:       req.PostalCode,
			Country:          req.Country,
		}
		if err := s.addressRepo.Create(permanentAddr); err != nil {
			return fmt.Errorf("failed to save permanent address: %w", err)
		}
	}

	return nil
}

func (s *OnboardingService) saveStep2Employment(draft *models.EmployeeOnboardingDraft, data map[string]interface{}) error {
	var req models.Step2EmploymentRequest
	jsonData, _ := json.Marshal(data)
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return fmt.Errorf("invalid step 2 data: %w", err)
	}

	// Validate department, position, location if provided
	if req.DepartmentID != nil {
		department, err := s.departmentRepo.FindByID(*req.DepartmentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("department with ID %d not found", *req.DepartmentID)
			}
			return fmt.Errorf("failed to validate department: %w", err)
		}
		// Tenant check removed (single-tenant)
		// Check if department is active
		if !department.IsActive {
			return fmt.Errorf("department with ID %d is not active", *req.DepartmentID)
		}

		// Auto-suggest reporting manager from department head if not explicitly provided.
		// The department's ManagerID is the head of department — use it as the default
		// reporting manager. This is optional: if the user explicitly provides a
		// reporting_manager_id (even null), we respect that choice.
		if req.ReportingManagerID == nil && department.ManagerID != nil {
			req.ReportingManagerID = department.ManagerID
		}
	}
	// Position ID is optional - only validate if provided
	if req.PositionID != nil {
		position, err := s.positionRepo.FindByID(*req.PositionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("position with ID %d not found", *req.PositionID)
			}
			return fmt.Errorf("failed to validate position: %w", err)
		}
		// Tenant check removed (single-tenant)
		// Check if position is active
		if !position.IsActive {
			return fmt.Errorf("position with ID %d is not active", *req.PositionID)
		}
	}
	if req.LocationID != nil {
		location, err := s.locationRepo.FindByID(*req.LocationID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("location with ID %d not found", *req.LocationID)
			}
			return fmt.Errorf("failed to validate location: %w", err)
		}
		// Tenant check removed (single-tenant)
		// Check if location is active
		if !location.IsActive {
			return fmt.Errorf("location with ID %d is not active", *req.LocationID)
		}
	}

	// Validate reporting manager if provided (either explicitly or auto-suggested from department)
	if req.ReportingManagerID != nil {
		manager, err := s.employeeRepo.FindByID(*req.ReportingManagerID)
		if err != nil || manager == nil {
			// If the auto-suggested manager is not found, silently clear it instead of failing
			req.ReportingManagerID = nil
		} else if manager.Status != models.StatusActive || !manager.IsActive {
			// If the manager is not active, silently clear
			req.ReportingManagerID = nil
		}
	}

	// Handle employee ID - preserve existing one from draft if not provided in step 2
	// Only update if a new employee_id is explicitly provided in step 2
	if req.EmployeeID != nil && *req.EmployeeID != "" {
		// Normalize the provided employee ID
		normalizedID := strings.ToUpper(strings.TrimSpace(*req.EmployeeID))
		// Only update if it's different from what's already in the draft
		if draft.EmployeeID == nil || *draft.EmployeeID != normalizedID {
			draft.EmployeeID = &normalizedID
			// Save the employee_id immediately to ensure it's persisted
			if err := s.draftRepo.Update(draft); err != nil {
				return fmt.Errorf("failed to save employee_id: %w", err)
			}
		}
	} else {
		// If employee_id is not provided in step 2, ensure draft has one from step 1
		// If draft doesn't have employee_id, generate one (shouldn't happen if step 1 was saved)
		if draft.EmployeeID == nil || *draft.EmployeeID == "" {
			generatedID := s.generateEmployeeID()
			draft.EmployeeID = &generatedID
			// Save the generated employee_id immediately
			if err := s.draftRepo.Update(draft); err != nil {
				return fmt.Errorf("failed to save generated employee_id: %w", err)
			}
		} else {
			// Normalize existing employee_id to ensure consistency
			normalizedID := strings.ToUpper(strings.TrimSpace(*draft.EmployeeID))
			if *draft.EmployeeID != normalizedID {
				draft.EmployeeID = &normalizedID
				if err := s.draftRepo.Update(draft); err != nil {
					return fmt.Errorf("failed to normalize employee_id: %w", err)
				}
			}
		}
	}

	// Check if employment details already exist for this draft
	existingDetails, err := s.employmentDetailsRepo.FindByDraftID(draft.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing employment details: %w", err)
	}

	// Create or update employment details
	employmentDetails := &models.EmployeeEmploymentDetails{
		DraftID:                  &draft.ID,
		EmployeeIDString:         draft.EmployeeID, // Use draft's employee_id (already set above)
		OfficialEmail:            req.OfficialEmail,
		DateOfJoining:            req.DateOfJoining,
		DepartmentID:             req.DepartmentID,
		PositionID:               req.PositionID,
		Grade:                    req.Grade,
		ReportingManagerID:       req.ReportingManagerID,
		EmploymentType:           req.EmploymentType,
		LocationID:               req.LocationID,
		Shift:                    req.Shift,
		WorkPhone:                req.WorkPhone,
		ProbationPeriodDays:      req.ProbationPeriodDays,
		ExpectedConfirmationDate: req.ExpectedConfirmationDate,
	}

	if existingDetails != nil {
		// Update existing record
		employmentDetails.ID = existingDetails.ID
		if err := s.employmentDetailsRepo.Update(employmentDetails); err != nil {
			return fmt.Errorf("failed to update employment details: %w", err)
		}
	} else {
		// Create new record
		if err := s.employmentDetailsRepo.Create(employmentDetails); err != nil {
			return fmt.Errorf("failed to save employment details: %w", err)
		}
	}

	return nil
}

func (s *OnboardingService) saveStep3Salary(draft *models.EmployeeOnboardingDraft, data map[string]interface{}) error {
	var req models.Step3SalaryRequest
	jsonData, _ := json.Marshal(data)
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return fmt.Errorf("invalid step 3 data: %w", err)
	}

	// Check if salary already exists for this draft
	existingSalary, err := s.salaryRepo.FindByDraftID(draft.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing salary: %w", err)
	}

	// Create or update salary component
	salary := &models.EmployeeSalaryComponent{
		DraftID:            &draft.ID,
		EmployeeIDString:   draft.EmployeeID, // Set employee_id_string from draft
		AnnualCTC:          req.AnnualCTC,
		CTCEffectiveDate:   req.CTCEffectiveDate,
		Currency:           req.Currency,
		BasicSalary:        req.BasicSalary,
		HouseRentAllowance: req.HouseRentAllowance,
		TransportAllowance: req.TransportAllowance,
		SpecialAllowance:   req.SpecialAllowance,
		OtherAllowances:    req.OtherAllowances,
		IncomeTax:          req.IncomeTax,
		ProvidentFund:      req.ProvidentFund,
		ProfessionalTax:    req.ProfessionalTax,
		OtherDeductions:    req.OtherDeductions,
	}

	// Calculate gross and net salary
	salary.CalculateGrossSalary()
	salary.CalculateNetSalary()

	if existingSalary != nil {
		// Update existing record
		salary.ID = existingSalary.ID
		if err := s.salaryRepo.Update(salary); err != nil {
			return fmt.Errorf("failed to update salary: %w", err)
		}
	} else {
		// Create new record
		if err := s.salaryRepo.Create(salary); err != nil {
			return fmt.Errorf("failed to save salary: %w", err)
		}
	}

	return nil
}

func (s *OnboardingService) saveStep4Bank(draft *models.EmployeeOnboardingDraft, data map[string]interface{}) error {
	var req models.Step4BankRequest
	jsonData, _ := json.Marshal(data)
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return fmt.Errorf("invalid step 4 data: %w", err)
	}

	// Check if bank account already exists for this draft
	existingBank, err := s.bankRepo.FindByDraftID(draft.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing bank account: %w", err)
	}

	// Create or update bank account
	bank := &models.EmployeeBankAccount{
		DraftID:           &draft.ID,
		EmployeeIDString:  draft.EmployeeID, // Set employee_id_string from draft
		BankName:          req.BankName,
		AccountHolderName: req.AccountHolderName,
		AccountNumber:     req.AccountNumber,
		AccountType:       req.AccountType,
		BranchName:        req.BranchName,
		SWIFTCode:         req.SWIFTCode,
		IsPrimary:         true,
	}

	if existingBank != nil {
		// Update existing record
		bank.ID = existingBank.ID
		if err := s.bankRepo.Update(bank); err != nil {
			return fmt.Errorf("failed to update bank account: %w", err)
		}
	} else {
		// Create new record
		if err := s.bankRepo.Create(bank); err != nil {
			return fmt.Errorf("failed to save bank account: %w", err)
		}
	}

	return nil
}

func (s *OnboardingService) saveStep5Statutory(draft *models.EmployeeOnboardingDraft, data map[string]interface{}) error {
	var req models.Step5StatutoryRequest
	jsonData, _ := json.Marshal(data)
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return fmt.Errorf("invalid step 5 data: %w", err)
	}

	// Check if statutory info already exists for this draft
	existingStatutory, err := s.statutoryRepo.FindByDraftID(draft.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing statutory info: %w", err)
	}

	// Create or update statutory info
	statutory := &models.EmployeeStatutoryInfo{
		DraftID:              &draft.ID,
		EmployeeIDString:     draft.EmployeeID, // Set employee_id_string from draft
		TINNumber:            req.TINNumber,
		NSSFNumber:           req.NSSFNumber,
		NHIFNumber:           req.NHIFNumber,
		WCFNumber:            req.WCFNumber,
		SDLNumber:            req.SDLNumber,
		PassportNumber:       req.PassportNumber,
		PassportExpiryDate:   req.PassportExpiryDate,
		WorkPermitNumber:     req.WorkPermitNumber,
		WorkPermitExpiryDate: req.WorkPermitExpiryDate,
	}

	if existingStatutory != nil {
		// Update existing record
		statutory.ID = existingStatutory.ID
		if err := s.statutoryRepo.Update(statutory); err != nil {
			return fmt.Errorf("failed to update statutory info: %w", err)
		}
	} else {
		// Create new record
		if err := s.statutoryRepo.Create(statutory); err != nil {
			return fmt.Errorf("failed to save statutory info: %w", err)
		}
	}

	return nil
}

func (s *OnboardingService) saveStep6Documents(draft *models.EmployeeOnboardingDraft, data map[string]interface{}) error {
	var req models.Step6DocumentsRequest
	jsonData, _ := json.Marshal(data)
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return fmt.Errorf("invalid step 6 data: %w", err)
	}

	// Delete existing documents for this draft
	s.documentRepo.DeleteByDraftID(draft.ID)

	// Create documents
	for _, docInput := range req.Documents {
		doc := &models.EmployeeDocument{
			DraftID:          &draft.ID,
			EmployeeIDString: draft.EmployeeID, // Set employee_id_string from draft
			DocumentType:     docInput.DocumentType,
			FileName:         docInput.FileName,
			FileURL:          docInput.FileURL,
			FileSize:         docInput.FileSize,
			MimeType:         docInput.MimeType,
			Description:      docInput.Description,
		}
		if err := s.documentRepo.Create(doc); err != nil {
			return fmt.Errorf("failed to save document: %w", err)
		}
	}

	return nil
}

func (s *OnboardingService) saveStep7Assets(draft *models.EmployeeOnboardingDraft, data map[string]interface{}) error {
	// Parse assets data - expects array of assets
	assetsData, ok := data["assets"].([]interface{})
	if !ok {
		// Try direct array format
		if assetsArray, ok := data["assets"].([]map[string]interface{}); ok {
			assetsData = make([]interface{}, len(assetsArray))
			for i, asset := range assetsArray {
				assetsData[i] = asset
			}
		} else {
			return fmt.Errorf("assets must be an array")
		}
	}

	// Delete existing assets for this draft
	s.assetRepo.DeleteByDraftID(draft.ID)

	// Create asset records
	for _, assetData := range assetsData {
		assetMap, ok := assetData.(map[string]interface{})
		if !ok {
			continue
		}

		asset := &models.EmployeeAsset{
			DraftID:          &draft.ID,
			EmployeeIDString: draft.EmployeeID, // Set employee_id_string from draft
		}

		if assetType, ok := assetMap["asset_type"].(string); ok {
			asset.AssetType = assetType
		} else {
			return fmt.Errorf("asset_type is required for each asset")
		}

		if assetName, ok := assetMap["asset_name"].(string); ok {
			asset.AssetName = assetName
		} else {
			return fmt.Errorf("asset_name is required for each asset")
		}

		if serialNumber, ok := assetMap["serial_number"].(string); ok {
			asset.SerialNumber = &serialNumber
		}
		if assetTag, ok := assetMap["asset_tag"].(string); ok {
			asset.AssetTag = &assetTag
		}
		if assignedDateStr, ok := assetMap["assigned_date"].(string); ok {
			if assignedDate, err := time.Parse(time.RFC3339, assignedDateStr); err == nil {
				asset.AssignedDate = &assignedDate
			}
		}
		if expectedReturnDateStr, ok := assetMap["expected_return_date"].(string); ok {
			if expectedReturnDate, err := time.Parse(time.RFC3339, expectedReturnDateStr); err == nil {
				asset.ExpectedReturnDate = &expectedReturnDate
			}
		}
		if condition, ok := assetMap["condition"].(string); ok {
			asset.Condition = &condition
		}
		if notes, ok := assetMap["notes"].(string); ok {
			asset.Notes = &notes
		}

		if err := s.assetRepo.Create(asset); err != nil {
			return fmt.Errorf("failed to save asset: %w", err)
		}
	}

	return nil
}

func (s *OnboardingService) saveStep8Policies(draft *models.EmployeeOnboardingDraft, data map[string]interface{}) error {
	var req models.Step8PoliciesRequest
	jsonData, _ := json.Marshal(data)
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return fmt.Errorf("invalid step 8 data: %w", err)
	}

	// Delete existing policy for this draft
	s.policyRepo.DeleteByDraftID(draft.ID)

	// Create policy assignment
	policy := &models.EmployeePolicy{
		DraftID:            &draft.ID,
		EmployeeIDString:   draft.EmployeeID, // Set employee_id_string from draft
		LeavePolicyID:      req.LeavePolicyID,
		AttendancePolicyID: req.AttendancePolicyID,
		WeeklyOffDays:      req.WeeklyOffDays,
		EffectiveDate:      req.EffectiveDate,
	}

	if err := s.policyRepo.Create(policy); err != nil {
		return fmt.Errorf("failed to save policy: %w", err)
	}

	return nil
}

func (s *OnboardingService) saveStep9EmergencyContacts(draft *models.EmployeeOnboardingDraft, data map[string]interface{}) error {
	var req models.Step9EmergencyContactsRequest
	jsonData, _ := json.Marshal(data)
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return fmt.Errorf("invalid step 9 data: %w", err)
	}

	// Delete existing emergency contacts for this draft
	s.contactRepo.DeleteByDraftID(draft.ID)

	// Create emergency contacts
	for _, contactInput := range req.Contacts {
		contact := &models.EmployeeEmergencyContact{
			DraftID:          &draft.ID,
			EmployeeIDString: draft.EmployeeID, // Set employee_id_string from draft
			ContactName:      contactInput.ContactName,
			Relationship:     contactInput.Relationship,
			PhoneNumber:      contactInput.PhoneNumber,
			AlternatePhone:   contactInput.AlternatePhone,
			Email:            contactInput.Email,
			Address:          contactInput.Address,
			IsPrimary:        contactInput.IsPrimary,
		}
		if err := s.contactRepo.Create(contact); err != nil {
			return fmt.Errorf("failed to save emergency contact: %w", err)
		}
	}

	return nil
}

func (s *OnboardingService) saveStep10Notes(draft *models.EmployeeOnboardingDraft, data map[string]interface{}) error {
	var req models.Step10NotesRequest
	jsonData, _ := json.Marshal(data)
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return fmt.Errorf("invalid step 10 data: %w", err)
	}

	// Store step 10 data in draft's StepData field
	stepData := make(map[string]interface{})
	if draft.StepData != "" {
		json.Unmarshal([]byte(draft.StepData), &stepData)
	}
	stepData["step10"] = req
	stepDataBytes, _ := json.Marshal(stepData)
	draft.StepData = string(stepDataBytes)

	return nil
}

// loadStepData loads all step data into the response from database tables
func (s *OnboardingService) loadStepData(draft *models.EmployeeOnboardingDraft, response *models.GetDraftResponse) {
	// Step 1: Personal Info & Addresses
	basicInfo, _ := s.basicInfoRepo.FindByDraftID(draft.ID)
	addresses, _ := s.addressRepo.FindByDraftID(draft.ID)

	step1Response := make(map[string]interface{})
	if basicInfo != nil {
		// Convert basic info to map for response
		step1Response["photo_url"] = basicInfo.PhotoURL
		step1Response["first_name"] = basicInfo.FirstName
		step1Response["middle_name"] = basicInfo.MiddleName
		step1Response["last_name"] = basicInfo.LastName
		step1Response["date_of_birth"] = basicInfo.DateOfBirth
		step1Response["gender"] = basicInfo.Gender
		step1Response["marital_status"] = basicInfo.MaritalStatus
		step1Response["blood_group"] = basicInfo.BloodGroup
		step1Response["nationality"] = basicInfo.Nationality
		step1Response["personal_email"] = basicInfo.PersonalEmail
		step1Response["mobile_number"] = basicInfo.MobileNumber
		step1Response["alternate_number"] = basicInfo.AlternateNumber
	}
	if len(addresses) > 0 {
		step1Response["addresses"] = addresses
	}

	if len(step1Response) > 0 {
		response.Steps["1"] = step1Response
	}

	// Step 2: Employment Details
	employmentDetails, _ := s.employmentDetailsRepo.FindByDraftID(draft.ID)
	if employmentDetails != nil {
		step2Response := make(map[string]interface{})
		step2Response["employee_id"] = employmentDetails.EmployeeIDString
		step2Response["official_email"] = employmentDetails.OfficialEmail
		step2Response["date_of_joining"] = employmentDetails.DateOfJoining
		step2Response["department_id"] = employmentDetails.DepartmentID
		step2Response["position_id"] = employmentDetails.PositionID
		step2Response["grade"] = employmentDetails.Grade
		step2Response["reporting_manager_id"] = employmentDetails.ReportingManagerID
		step2Response["employment_type"] = employmentDetails.EmploymentType
		step2Response["location_id"] = employmentDetails.LocationID
		step2Response["shift"] = employmentDetails.Shift
		step2Response["work_phone"] = employmentDetails.WorkPhone
		step2Response["probation_period_days"] = employmentDetails.ProbationPeriodDays
		step2Response["expected_confirmation_date"] = employmentDetails.ExpectedConfirmationDate

		// Also include employee_id from draft if set
		if draft.EmployeeID != nil {
			step2Response["employee_id"] = *draft.EmployeeID
		}

		response.Steps["2"] = step2Response
	} else if draft.EmployeeID != nil {
		// Fallback: if no employment details but employee_id is set
		response.Steps["2"] = map[string]interface{}{
			"employee_id": *draft.EmployeeID,
		}
	}

	// Step 3: Salary
	salary, _ := s.salaryRepo.FindByDraftID(draft.ID)
	if salary != nil {
		response.Steps["3"] = salary
	}

	// Step 4: Bank
	bank, _ := s.bankRepo.FindByDraftID(draft.ID)
	if bank != nil {
		response.Steps["4"] = bank
	}

	// Step 5: Statutory
	statutory, _ := s.statutoryRepo.FindByDraftID(draft.ID)
	if statutory != nil {
		response.Steps["5"] = statutory
	}

	// Step 6: Documents
	documents, _ := s.documentRepo.FindByDraftID(draft.ID)
	if len(documents) > 0 {
		response.Steps["6"] = documents
	}

	// Step 7: Assets
	assets, _ := s.assetRepo.FindByDraftID(draft.ID)
	if len(assets) > 0 {
		response.Steps["7"] = assets
	}

	// Step 8: Policies
	policy, _ := s.policyRepo.FindByDraftID(draft.ID)
	if policy != nil {
		response.Steps["8"] = policy
	}

	// Step 9: Emergency Contacts
	contacts, _ := s.contactRepo.FindByDraftID(draft.ID)
	if len(contacts) > 0 {
		response.Steps["9"] = contacts
	}

	// Step 10: Notes (will be in employee record)
	response.Steps["10"] = map[string]interface{}{"note": "Notes stored in employee record"}
}

// validateRequiredSteps checks if required steps are completed
func (s *OnboardingService) validateRequiredSteps(draft *models.EmployeeOnboardingDraft) bool {
	// Steps 1, 2, and 9 are required (Personal Info, Employment, Emergency Contacts)
	requiredSteps := []models.OnboardingStep{
		models.StepPersonalInfo,
		models.StepEmployment,
		models.StepEmergency,
	}

	for _, step := range requiredSteps {
		if !draft.HasStepCompleted(step) {
			return false
		}
	}

	return true
}

// createEmployeeFromDraft creates an employee record from draft data
func (s *OnboardingService) createEmployeeFromDraft(draft *models.EmployeeOnboardingDraft) (*models.Employee, error) {
	// Load step 1 data from employee_basic_information table
	basicInfo, err := s.basicInfoRepo.FindByDraftID(draft.ID)
	if err != nil {
		return nil, fmt.Errorf("step 1 (personal information) is required: %w", err)
	}

	// Load step 2 data from employee_employment_details table
	employmentDetails, err := s.employmentDetailsRepo.FindByDraftID(draft.ID)
	if err != nil {
		return nil, fmt.Errorf("step 2 (employment details) is required: %w", err)
	}

	// Load step 10 data (notes) from step_data JSON (can be moved to separate table later)
	var step10Notes *string
	var stepData map[string]interface{}
	if draft.StepData != "" {
		json.Unmarshal([]byte(draft.StepData), &stepData)
		if step10Raw, ok := stepData["step10"]; ok {
			if step10Map, ok := step10Raw.(map[string]interface{}); ok {
				if notes, ok := step10Map["notes"].(string); ok {
					step10Notes = &notes
				}
			}
		}
	}

	// Generate employee ID if not provided
	employeeID := "EMP"
	if draft.EmployeeID != nil && *draft.EmployeeID != "" {
		employeeID = *draft.EmployeeID
	} else if employmentDetails.EmployeeIDString != nil && *employmentDetails.EmployeeIDString != "" {
		employeeID = *employmentDetails.EmployeeIDString
	} else {
		// Auto-generate employee ID using helper function
		employeeID = s.generateEmployeeID()
	}

	// Check if employee ID already exists, regenerate if needed
	maxAttempts := 10
	attempts := 0
	for s.employeeRepo.ExistsByEmployeeID(employeeID) && attempts < maxAttempts {
		employeeID = s.generateEmployeeID()
		attempts++
	}
	if attempts >= maxAttempts {
		return nil, errors.New("failed to generate unique employee ID after multiple attempts")
	}

	// Get department to extract organization_id
	var organizationID *uint
	var organizationUnitID *uint
	if employmentDetails.DepartmentID != nil {
		department, err := s.departmentRepo.FindByID(*employmentDetails.DepartmentID)
		if err == nil {
			organizationID = department.OrganizationID
			organizationUnitID = department.OrganizationUnitID
		}
	}

	// Build employee from table data
	employee := &models.Employee{
		EmployeeID:               employeeID,
		FirstName:                basicInfo.FirstName,
		MiddleName:               basicInfo.MiddleName,
		LastName:                 basicInfo.LastName,
		DateOfBirth:              basicInfo.DateOfBirth,
		Gender:                   basicInfo.Gender,
		MaritalStatus:            basicInfo.MaritalStatus,
		BloodGroup:               basicInfo.BloodGroup,
		Nationality:              basicInfo.Nationality,
		PhotoURL:                 basicInfo.PhotoURL,
		PersonalEmail:            basicInfo.PersonalEmail,
		WorkEmail:                employmentDetails.OfficialEmail,
		PhoneNumber:              basicInfo.MobileNumber,
		AlternatePhone:           basicInfo.AlternateNumber,
		DepartmentID:             employmentDetails.DepartmentID,
		PositionID:               employmentDetails.PositionID,
		LocationID:               employmentDetails.LocationID,
		HireDate:                 employmentDetails.DateOfJoining,
		EmploymentType:           employmentDetails.EmploymentType,
		ReportsToID:              employmentDetails.ReportingManagerID,
		Grade:                    employmentDetails.Grade,
		Shift:                    employmentDetails.Shift,
		WorkPhone:                employmentDetails.WorkPhone,
		ProbationPeriodDays:      employmentDetails.ProbationPeriodDays,
		ExpectedConfirmationDate: employmentDetails.ExpectedConfirmationDate,
		OrganizationID:           organizationID,
		OrganizationUnitID:       organizationUnitID,
		Status:                   models.StatusActive,
		IsActive:                 true,
		Notes:                    "",
	}

	// Set notes from step 10
	if step10Notes != nil {
		employee.Notes = *step10Notes
	}

	// Set default currency
	employee.Currency = "TZS"

	// Load salary to set basic salary
	salary, _ := s.salaryRepo.FindByDraftID(draft.ID)
	if salary != nil && salary.BasicSalary != nil {
		employee.Salary = salary.BasicSalary
		if salary.Currency != "" {
			employee.Currency = salary.Currency
		}
	}

	if err := s.employeeRepo.Create(employee); err != nil {
		return nil, fmt.Errorf("failed to create employee: %w", err)
	}

	return employee, nil
}

// migrateDraftDataToEmployee migrates all draft data to employee
func (s *OnboardingService) migrateDraftDataToEmployee(draftID uint, employeeID uint) error {
	// Migrate all related data from tables
	if err := s.basicInfoRepo.MigrateToEmployee(draftID, employeeID); err != nil {
		return fmt.Errorf("failed to migrate basic information: %w", err)
	}
	if err := s.employmentDetailsRepo.MigrateToEmployee(draftID, employeeID); err != nil {
		return fmt.Errorf("failed to migrate employment details: %w", err)
	}
	if err := s.addressRepo.MigrateToEmployee(draftID, employeeID); err != nil {
		return fmt.Errorf("failed to migrate addresses: %w", err)
	}
	if err := s.salaryRepo.MigrateToEmployee(draftID, employeeID); err != nil {
		return fmt.Errorf("failed to migrate salary: %w", err)
	}
	if err := s.bankRepo.MigrateToEmployee(draftID, employeeID); err != nil {
		return fmt.Errorf("failed to migrate bank: %w", err)
	}
	if err := s.statutoryRepo.MigrateToEmployee(draftID, employeeID); err != nil {
		return fmt.Errorf("failed to migrate statutory: %w", err)
	}
	if err := s.documentRepo.MigrateToEmployee(draftID, employeeID); err != nil {
		return fmt.Errorf("failed to migrate documents: %w", err)
	}
	if err := s.assetRepo.MigrateToEmployee(draftID, employeeID); err != nil {
		return fmt.Errorf("failed to migrate assets: %w", err)
	}
	if err := s.contactRepo.MigrateToEmployee(draftID, employeeID); err != nil {
		return fmt.Errorf("failed to migrate contacts: %w", err)
	}
	if err := s.policyRepo.MigrateToEmployee(draftID, employeeID); err != nil {
		return fmt.Errorf("failed to migrate policy: %w", err)
	}

	return nil
}

// generateEmployeeID generates a unique employee ID starting with "EMP"
// Format: EMP + 6-digit sequential number (e.g., EMP000001, EMP000002)
func (s *OnboardingService) generateEmployeeID() string {
	// Get the current highest employee ID number
	maxAttempts := 100
	for i := 0; i < maxAttempts; i++ {
		// Generate ID using timestamp + random component to ensure uniqueness
		timestamp := time.Now().Unix()
		// Use last 6 digits of timestamp + microsecond component
		idNum := (timestamp % 1000000) + int64(i)
		employeeID := fmt.Sprintf("EMP%06d", idNum)

		// Check if this ID already exists
		if !s.employeeRepo.ExistsByEmployeeID(employeeID) {
			return employeeID
		}
	}

	// Fallback: use timestamp with nanoseconds if all attempts fail
	return fmt.Sprintf("EMP%06d", time.Now().UnixNano()%1000000)
}

// generateSecurePassword generates a secure random password
func (s *OnboardingService) generateSecurePassword() (string, error) {
	// Generate 16 random bytes
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random password: %w", err)
	}
	// Convert to base64 and take first 12 characters, then add some complexity
	password := base64.URLEncoding.EncodeToString(bytes)[:12]
	// Ensure password has at least one uppercase, lowercase, number, and special char
	// For simplicity, we'll use a pattern: 8 random chars + 1 uppercase + 1 number + 1 special
	password = password[:8] + "A" + "1" + "!"
	return password, nil
}

// createUserForEmployee creates a user account for the employee
// isAdminEmployee indicates if the employee should have admin role (for admin employees)
func (s *OnboardingService) createUserForEmployee(employee *models.Employee, draft *models.EmployeeOnboardingDraft, updatedBy *uint, isAdminEmployee bool) (*UserCredentials, error) {
	// Get employee's work email from employment details
	employmentDetails, err := s.employmentDetailsRepo.FindByEmployeeID(employee.ID)
	if err != nil {
		// Try to get from draft if not migrated yet
		employmentDetails, err = s.employmentDetailsRepo.FindByDraftID(draft.ID)
		if err != nil {
			return nil, fmt.Errorf("employment details not found: %w", err)
		}
	}

	if employmentDetails.OfficialEmail == nil || *employmentDetails.OfficialEmail == "" {
		return nil, errors.New("employee work email is required to create user account")
	}

	email := *employmentDetails.OfficialEmail

	// Check if user already exists with this email (existing user being onboarded as employee)
	existingUser, err := s.userRepo.FindByEmail(email)
	if err == nil && existingUser != nil {
		// User already exists - ensure they are not already linked to another employee
		existingEmp, _ := s.employeeRepo.FindByUserID(existingUser.ID)
		if existingEmp != nil && existingEmp.ID != employee.ID {
			return nil, fmt.Errorf("user is already onboarded as an employee (employee_id: %s)", existingEmp.EmployeeID)
		}
		// Link employee to existing user; no new credentials
		employee.UserID = &existingUser.ID
		if err := s.employeeRepo.Update(employee); err != nil {
			_ = err
		}
		// Assign "employee" role from roles table for ALL existing users (admin, hr, or user) so they have employee access
		employeeRole, err := s.roleRepo.FindByCode("employee")
		if err == nil && employeeRole != nil {
			userRoles, _ := s.roleRepo.GetUserRoles(existingUser.ID)
			hasEmployeeRole := false
			for _, role := range userRoles {
				if role.Code == "employee" {
					hasEmployeeRole = true
					break
				}
			}
			if !hasEmployeeRole {
				_ = s.roleRepo.AssignUserRole(existingUser.ID, employeeRole.ID, updatedBy)
			}
		}
		// Return nil credentials (user already exists; they use their existing login)
		return nil, nil
	}

	// Get basic info for name
	basicInfo, _ := s.basicInfoRepo.FindByEmployeeID(employee.ID)
	if basicInfo == nil {
		basicInfo, _ = s.basicInfoRepo.FindByDraftID(draft.ID)
	}

	firstName := employee.FirstName
	lastName := employee.LastName
	if basicInfo != nil {
		firstName = basicInfo.FirstName
		lastName = basicInfo.LastName
	}

	// Generate username from email
	emailParts := strings.Split(email, "@")
	username := emailParts[0]
	username = strings.ToLower(strings.ReplaceAll(username, ".", ""))
	username = strings.ReplaceAll(username, "_", "")
	username = strings.ReplaceAll(username, "-", "")

	// Generate secure password
	plainPassword, err := s.generateSecurePassword()
	if err != nil {
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Determine user role
	// If isAdminEmployee is true, assign admin role; otherwise default to user role
	userRole := userModels.RoleUser
	if isAdminEmployee {
		userRole = userModels.RoleAdmin
	}

	// Create user with determined role
	user := &userModels.User{
		Username:      username,
		Email:         email,
		Password:      string(hashedPassword),
		FirstName:     firstName,
		LastName:      lastName,
		Role:          userRole, // Admin if isAdminEmployee is true, otherwise user
		Status:        "active",
		EmailVerified: false,
		IsActive:      true,
		UpdatedBy:     updatedBy,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Link user to employee
	employee.UserID = &user.ID
	if err := s.employeeRepo.Update(employee); err != nil {
		// Log error but don't fail - user is created
		_ = err
	}

	// Assign "employee" role from roles table (if exists)
	// This allows users to have multiple roles (e.g., admin + employee)
	employeeRole, err := s.roleRepo.FindByCode("employee")
	if err == nil && employeeRole != nil {
		// Assign role via user_roles table
		if err := s.roleRepo.AssignUserRole(user.ID, employeeRole.ID, updatedBy); err != nil {
			// Log but don't fail - user is created with their primary role
			_ = err
		}
	}

	// If user is admin, also try to assign "system_admin" or "admin" role from roles table
	// This allows admin employees to have both admin (UserRole enum) and system_admin (roles table) roles
	if userRole == userModels.RoleAdmin {
		// Try to find and assign system_admin role from roles table
		systemAdminRole, err := s.roleRepo.FindByCode("system_admin")
		if err == nil && systemAdminRole != nil {
			_ = s.roleRepo.AssignUserRole(user.ID, systemAdminRole.ID, updatedBy)
		}
		// Also try "admin" role from roles table
		adminRole, err := s.roleRepo.FindByCode("admin")
		if err == nil && adminRole != nil {
			_ = s.roleRepo.AssignUserRole(user.ID, adminRole.ID, updatedBy)
		}
	}

	// Return credentials
	return &UserCredentials{
		Email:    email,
		Password: plainPassword,
		Username: username,
	}, nil
}

// sendWelcomeEmailWithCredentials sends a welcome email with login credentials to the employee
func (s *OnboardingService) sendWelcomeEmailWithCredentials(recipientEmail string, employee *models.Employee, credentials *UserCredentials) {
	// Create email service instance
	emailService := email.NewEmailService()

	// Build email subject and body
	subject := fmt.Sprintf("Welcome to ScoopWorks - Your Login Credentials")

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to ScoopWorks</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; text-align: center; border-radius: 5px; }
        .content { background-color: #fff; padding: 30px; border-radius: 5px; margin-top: 20px; border: 1px solid #e9ecef; }
        .credentials { background-color: #f8f9fa; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .credential-item { margin: 10px 0; }
        .label { font-weight: bold; color: #495057; }
        .value { background-color: #fff; padding: 8px 12px; border: 1px solid #dee2e6; border-radius: 3px; font-family: monospace; }
        .footer { margin-top: 30px; text-align: center; color: #6c757d; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to ScoopWorks!</h1>
        </div>
        
        <div class="content">
            <h2>Hello %s %s,</h2>
            
            <p>Your employee onboarding has been successfully completed. Below are your login credentials for the ScoopWorks HR Management System:</p>
            
            <div class="credentials">
                <div class="credential-item">
                    <span class="label">Employee ID:</span>
                    <div class="value">%s</div>
                </div>
                <div class="credential-item">
                    <span class="label">Login URL:</span>
                    <div class="value">https://hrms.scoopworks.com</div>
                </div>
                <div class="credential-item">
                    <span class="label">Username/Email:</span>
                    <div class="value">%s</div>
                </div>
                <div class="credential-item">
                    <span class="label">Password:</span>
                    <div class="value">%s</div>
                </div>
            </div>
            
            <p><strong>Important Security Notes:</strong></p>
            <ul>
                <li>This is a temporary password - please change it after your first login</li>
                <li>Never share your credentials with anyone</li>
                <li>If you did not request this account, please contact HR immediately</li>
            </ul>
            
            <p>We recommend logging in as soon as possible to familiarize yourself with the system and update your password.</p>
            
            <p>Best regards,<br>
            <strong>ScoopWorks HR Team</strong></p>
        </div>
        
        <div class="footer">
            <p>This is an automated message. Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
`,
		employee.FirstName,
		employee.LastName,
		employee.EmployeeID,
		credentials.Email,
		credentials.Password)

	// Send email
	err := emailService.SendEmail(recipientEmail, subject, body)
	if err != nil {
		fmt.Printf("Failed to send welcome email to %s: %v\n", recipientEmail, err)
	}
}
