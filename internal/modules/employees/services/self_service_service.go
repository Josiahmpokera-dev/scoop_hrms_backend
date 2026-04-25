package services

import (
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
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	userRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories"
	"gorm.io/gorm"
)

// SelfServiceService handles self-service operations for employees
type SelfServiceService struct {
	employeeRepo          *employeeRepos.EmployeeRepository
	basicInfoRepo         *employeeRepos.EmployeeBasicInformationRepository
	employmentDetailsRepo *employeeRepos.EmployeeEmploymentDetailsRepository
	emergencyContactRepo  *employeeRepos.EmployeeEmergencyContactRepository
	addressRepo           *employeeRepos.EmployeeAddressRepository
	documentRepo          *employeeRepos.EmployeeDocumentRepository
	serviceRequestRepo    *employeeRepos.ServiceRequestRepository
	profileUpdateRepo     *employeeRepos.ProfileUpdateRequestRepository
	departmentRepo        *departmentRepos.DepartmentRepository
	positionRepo          *positionRepos.JobPositionRepository
	locationRepo          *locationRepos.LocationRepository
	userRepo              *userRepos.UserRepository
}

// NewSelfServiceService creates a new self-service service
func NewSelfServiceService() *SelfServiceService {
	return &SelfServiceService{
		employeeRepo:          employeeRepos.NewEmployeeRepository(),
		basicInfoRepo:         employeeRepos.NewEmployeeBasicInformationRepository(),
		employmentDetailsRepo: employeeRepos.NewEmployeeEmploymentDetailsRepository(),
		emergencyContactRepo:  employeeRepos.NewEmployeeEmergencyContactRepository(),
		addressRepo:           employeeRepos.NewEmployeeAddressRepository(),
		documentRepo:          employeeRepos.NewEmployeeDocumentRepository(),
		serviceRequestRepo:    employeeRepos.NewServiceRequestRepository(),
		profileUpdateRepo:     employeeRepos.NewProfileUpdateRequestRepository(),
		departmentRepo:        departmentRepos.NewDepartmentRepository(),
		positionRepo:          positionRepos.NewJobPositionRepository(),
		locationRepo:          locationRepos.NewLocationRepository(),
		userRepo:              userRepos.NewUserRepository(),
	}
}

// GetEmployeeProfile retrieves complete profile for the authenticated employee
func (s *SelfServiceService) GetEmployeeProfile(userID uint, tenantID *uint) (map[string]interface{}, error) {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No employee record found, return user profile instead
			return s.getUserProfile(userID, tenantID)
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	// Get basic information
	basicInfo, _ := s.basicInfoRepo.FindByEmployeeID(employee.ID)

	// Get employment details
	employmentDetails, _ := s.employmentDetailsRepo.FindByEmployeeID(employee.ID)

	// Get emergency contacts
	emergencyContacts, _ := s.emergencyContactRepo.FindByEmployeeID(employee.ID)

	// Get addresses
	addresses, _ := s.addressRepo.FindByEmployeeID(employee.ID)

	// Get documents
	documents, _ := s.documentRepo.FindByEmployeeID(employee.ID)

	// Get profile update status
	pendingRequests, _ := s.profileUpdateRepo.FindPendingByEmployeeID(employee.ID, tenantID)

	// Build response
	response := map[string]interface{}{
		"employee_id": employee.EmployeeID,
	}

	// Personal information
	personal := map[string]interface{}{
		"first_name":     employee.FirstName,
		"last_name":      employee.LastName,
		"full_name":      employee.FullName(),
		"date_of_birth":  nil,
		"gender":         nil,
		"marital_status": nil,
		"blood_group":    nil,
		"nationality":    nil,
		"personal_email": nil,
		"personal_phone": nil,
	}

	if basicInfo != nil {
		if basicInfo.PhotoURL != nil {
			personal["photo"] = *basicInfo.PhotoURL
		}
		if basicInfo.DateOfBirth != nil {
			personal["date_of_birth"] = basicInfo.DateOfBirth.Format("2006-01-02")
		}
		if basicInfo.Gender != nil {
			personal["gender"] = *basicInfo.Gender
		}
		if basicInfo.MaritalStatus != nil {
			personal["marital_status"] = *basicInfo.MaritalStatus
		}
		if basicInfo.BloodGroup != nil {
			personal["blood_group"] = *basicInfo.BloodGroup
		}
		if basicInfo.Nationality != nil {
			personal["nationality"] = *basicInfo.Nationality
		}
		if basicInfo.PersonalEmail != nil {
			personal["personal_email"] = *basicInfo.PersonalEmail
		}
		if basicInfo.MobileNumber != nil {
			personal["personal_phone"] = *basicInfo.MobileNumber
		}
	}

	// Addresses
	var currentAddress, permanentAddress map[string]interface{}
	for _, addr := range addresses {
		addrMap := map[string]interface{}{
			"address_line1": addr.AddressLine1,
			"address_line2": addr.AddressLine2,
			"city":          addr.City,
			"state":         addr.State,
			"country":       addr.Country,
			"postal_code":   addr.PostalCode,
		}
		if addr.AddressType == models.AddressTypeCurrent {
			currentAddress = addrMap
		} else if addr.AddressType == models.AddressTypePermanent {
			permanentAddress = addrMap
		}
	}
	personal["current_address"] = currentAddress
	personal["permanent_address"] = permanentAddress
	response["personal"] = personal

	// Job information
	job := map[string]interface{}{
		"employee_id":     employee.EmployeeID,
		"official_email":  nil,
		"date_of_joining": nil,
		"grade":           employee.Grade,
		"employment_type": employee.EmploymentType,
		"work_phone":      employee.WorkPhone,
		"shift":           employee.Shift,
	}

	if employmentDetails != nil {
		if employmentDetails.OfficialEmail != nil {
			job["official_email"] = *employmentDetails.OfficialEmail
		}
		if employmentDetails.DateOfJoining != nil {
			job["date_of_joining"] = employmentDetails.DateOfJoining.Format("2006-01-02")
		}
		if employmentDetails.Grade != nil {
			job["grade"] = *employmentDetails.Grade
		}
		if employmentDetails.EmploymentType != nil {
			job["employment_type"] = *employmentDetails.EmploymentType
		}
		if employmentDetails.WorkPhone != nil {
			job["work_phone"] = *employmentDetails.WorkPhone
		}
		if employmentDetails.Shift != nil {
			job["shift"] = *employmentDetails.Shift
		}

		// Department
		if employmentDetails.DepartmentID != nil {
			dept, _ := s.departmentRepo.FindByID(*employmentDetails.DepartmentID)
			if dept != nil {
				job["department"] = map[string]interface{}{
					"id":   dept.ID,
					"name": dept.Name,
					"code": dept.Code,
				}
			}
		}

		// Position
		if employmentDetails.PositionID != nil {
			pos, _ := s.positionRepo.FindByID(*employmentDetails.PositionID)
			if pos != nil {
				job["position"] = map[string]interface{}{
					"id":   pos.ID,
					"name": pos.Title,
					"code": pos.Code,
				}
			}
		}

		// Location
		if employmentDetails.LocationID != nil {
			loc, _ := s.locationRepo.FindByID(*employmentDetails.LocationID)
			if loc != nil {
				// Build address from location fields
				addressParts := []string{}
				if loc.AddressLine1 != nil && *loc.AddressLine1 != "" {
					addressParts = append(addressParts, *loc.AddressLine1)
				}
				if loc.AddressLine2 != nil && *loc.AddressLine2 != "" {
					addressParts = append(addressParts, *loc.AddressLine2)
				}
				if loc.City != nil && *loc.City != "" {
					addressParts = append(addressParts, *loc.City)
				}
				if loc.State != nil && *loc.State != "" {
					addressParts = append(addressParts, *loc.State)
				}
				if loc.Country != nil && *loc.Country != "" {
					addressParts = append(addressParts, *loc.Country)
				}
				address := strings.Join(addressParts, ", ")

				job["work_location"] = map[string]interface{}{
					"id":      loc.ID,
					"name":    loc.Name,
					"address": address,
				}
			}
		}

		// Reporting Manager
		if employmentDetails.ReportingManagerID != nil {
			manager, _ := s.employeeRepo.FindByID(*employmentDetails.ReportingManagerID)
			if manager != nil {
				managerEmail := ""
				if manager.WorkEmail != nil {
					managerEmail = *manager.WorkEmail
				}
				job["reporting_manager"] = map[string]interface{}{
					"id":          manager.ID,
					"employee_id": manager.EmployeeID,
					"full_name":   manager.FullName(),
					"email":       managerEmail,
				}
			}
		}
	}
	response["job"] = job

	// Emergency contacts
	contactsList := []map[string]interface{}{}
	for _, contact := range emergencyContacts {
		contactMap := map[string]interface{}{
			"id":           contact.ID,
			"name":         contact.ContactName,
			"relationship": contact.Relationship,
			"phone":        contact.PhoneNumber,
			"email":        contact.Email,
			"address":      contact.Address,
			"is_primary":   contact.IsPrimary,
		}
		contactsList = append(contactsList, contactMap)
	}
	response["emergency_contacts"] = contactsList

	// Documents
	documentsList := []map[string]interface{}{}
	for _, doc := range documents {
		docMap := map[string]interface{}{
			"id":            doc.ID,
			"document_type": doc.DocumentType,
			"file_name":     doc.FileName,
			"file_url":      doc.FileURL,
			"uploaded_at":   doc.CreatedAt,
		}
		documentsList = append(documentsList, docMap)
	}
	response["documents"] = documentsList

	// Profile update status
	hasPending := len(pendingRequests) > 0
	updateStatus := map[string]interface{}{
		"has_pending_updates":      hasPending,
		"last_update_request_date": nil,
		"last_update_status":       nil,
	}
	if len(pendingRequests) > 0 {
		latest := pendingRequests[0]
		updateStatus["last_update_request_date"] = latest.SubmittedAt
		updateStatus["last_update_status"] = string(latest.Status)
	}
	response["profile_update_status"] = updateStatus

	return response, nil
}

// getUserProfile retrieves profile for non-employee users (admins, HR, etc.)
func (s *SelfServiceService) getUserProfile(userID uint, tenantID *uint) (map[string]interface{}, error) {
	// Get user by ID
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Build user profile response
	response := map[string]interface{}{
		"user_id": user.ID,
		"role":    string(userModels.UserRole(user.Role)),
		"status":  user.Status,
	}

	// Personal information
	personal := map[string]interface{}{
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"full_name":    user.FullName(),
		"username":     user.Username,
		"email":        user.Email,
		"phone_number": nil,
	}

	if user.PhoneNumber != nil {
		personal["phone_number"] = *user.PhoneNumber
	}

	response["personal"] = personal

	// Account information
	account := map[string]interface{}{
		"email_verified": user.EmailVerified,
		"is_active":      user.IsActive,
		"last_login":     nil,
		"created_at":     user.CreatedAt,
		"updated_at":     user.UpdatedAt,
	}

	if user.LastLogin != nil {
		account["last_login"] = user.LastLogin
	}

	response["account"] = account

	// Job information (minimal for non-employees)
	job := map[string]interface{}{
		"role": string(user.Role),
	}
	response["job"] = job

	// Empty arrays for consistency
	response["emergency_contacts"] = []map[string]interface{}{}
	response["documents"] = []map[string]interface{}{}

	// Profile update status (not applicable for non-employees)
	response["profile_update_status"] = map[string]interface{}{
		"has_pending_updates":      false,
		"last_update_request_date": nil,
		"last_update_status":       nil,
	}

	return response, nil
}

// GetEmployeeDocuments retrieves documents for the authenticated employee
func (s *SelfServiceService) GetEmployeeDocuments(userID uint, tenantID *uint, documentType *string, page, pageSize int) ([]models.EmployeeDocument, int64, error) {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return nil, 0, errors.New("employee record not found for this user")
	}

	// Get documents
	documents, total, err := s.documentRepo.FindByEmployeeIDWithPagination(employee.ID, documentType, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get documents: %w", err)
	}

	return documents, total, nil
}

// CreateProfileUpdateRequest creates a profile update request
func (s *SelfServiceService) CreateProfileUpdateRequest(userID uint, tenantID *uint, section string, updates map[string]interface{}, reason *string, updatedBy *uint) (*models.ProfileUpdateRequest, error) {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("employee record not found for this user")
	}

	// Validate section
	validSection := models.ProfileUpdateSection(section)
	if validSection != models.ProfileSectionPersonal && validSection != models.ProfileSectionEmergencyContacts && validSection != models.ProfileSectionAddress {
		return nil, errors.New("invalid section. Must be 'personal', 'emergency_contacts', or 'address'")
	}

	// Convert updates to JSON
	updatesJSON, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize updates: %w", err)
	}

	// Generate request ID
	requestID, err := s.profileUpdateRepo.GenerateUpdateRequestID(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate request ID: %w", err)
	}

	// Create request
	request := &models.ProfileUpdateRequest{
		EmployeeID:      employee.ID,
		UpdateRequestID: requestID,
		Section:         validSection,
		Status:          models.ProfileUpdateStatusPending,
		UpdatesJSON:     string(updatesJSON),
		Reason:          reason,
		SubmittedAt:     time.Now(),
		UpdatedBy:       updatedBy,
	}

	// Get manager for approval (from employment details)
	employmentDetails, _ := s.employmentDetailsRepo.FindByEmployeeID(employee.ID)
	if employmentDetails != nil && employmentDetails.ReportingManagerID != nil {
		request.ApproverID = employmentDetails.ReportingManagerID
	}

	if err := s.profileUpdateRepo.Create(request); err != nil {
		return nil, fmt.Errorf("failed to create update request: %w", err)
	}

	return request, nil
}

// GetProfileUpdateStatus retrieves profile update request status
func (s *SelfServiceService) GetProfileUpdateStatus(userID uint, tenantID *uint) (map[string]interface{}, error) {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("employee record not found for this user")
	}

	// Get all requests
	allRequests, err := s.profileUpdateRepo.FindByEmployeeID(employee.ID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get update requests: %w", err)
	}

	// Separate pending and recent
	var pending []map[string]interface{}
	var recent []map[string]interface{}

	for _, req := range allRequests {
		reqMap := map[string]interface{}{
			"update_request_id": req.UpdateRequestID,
			"section":           string(req.Section),
			"status":            string(req.Status),
			"submitted_at":      req.SubmittedAt,
		}

		if req.Status == models.ProfileUpdateStatusPending {
			// Get approver info
			if req.ApproverID != nil {
				approver, _ := s.employeeRepo.FindByID(*req.ApproverID)
				if approver != nil {
					reqMap["approver"] = map[string]interface{}{
						"id":    approver.ID,
						"name":  approver.FullName(),
						"email": approver.WorkEmail,
					}
				}
			}
			pending = append(pending, reqMap)
		} else {
			if req.ApprovedAt != nil {
				reqMap["approved_at"] = *req.ApprovedAt
			}
			if req.RejectedAt != nil {
				reqMap["rejected_at"] = *req.RejectedAt
			}
			if req.ApproverID != nil {
				approver, _ := s.employeeRepo.FindByID(*req.ApproverID)
				if approver != nil {
					reqMap["approved_by"] = map[string]interface{}{
						"id":   approver.ID,
						"name": approver.FullName(),
					}
				}
			}
			recent = append(recent, reqMap)
		}
	}

	return map[string]interface{}{
		"pending_requests": pending,
		"recent_updates":   recent,
	}, nil
}

// ListServiceRequests retrieves service requests for the authenticated employee
func (s *SelfServiceService) ListServiceRequests(userID uint, tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.ServiceRequest, int64, error) {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return nil, 0, errors.New("employee record not found for this user")
	}

	// Get requests
	requests, total, err := s.serviceRequestRepo.FindByEmployeeID(employee.ID, tenantID, page, pageSize, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get service requests: %w", err)
	}

	return requests, total, nil
}

// GetServiceRequestDetails retrieves details of a specific service request
func (s *SelfServiceService) GetServiceRequestDetails(userID uint, requestID uint) (*models.ServiceRequest, error) {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("employee record not found for this user")
	}

	// Get request
	request, err := s.serviceRequestRepo.FindByID(requestID)
	if err != nil {
		return nil, errors.New("service request not found")
	}

	// Verify ownership
	if request.EmployeeID != employee.ID {
		return nil, errors.New("unauthorized: this request does not belong to you")
	}

	return request, nil
}

// CreateServiceRequest creates a new service request
func (s *SelfServiceService) CreateServiceRequest(userID uint, tenantID *uint, requestType, category, subject string, description *string, priority string, additionalData map[string]interface{}, updatedBy *uint) (*models.ServiceRequest, error) {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("employee record not found for this user")
	}

	// Validate request type
	validType := models.ServiceRequestType(requestType)
	if validType != models.ServiceRequestTypeHRLetter && validType != models.ServiceRequestTypeITRequest && validType != models.ServiceRequestTypeFacilities {
		return nil, errors.New("invalid request type. Must be 'hr_letter', 'it_request', or 'facilities'")
	}

	// Validate priority
	validPriority := models.ServiceRequestPriority(priority)
	if validPriority != models.ServiceRequestPriorityLow && validPriority != models.ServiceRequestPriorityMedium && validPriority != models.ServiceRequestPriorityHigh && validPriority != models.ServiceRequestPriorityUrgent {
		validPriority = models.ServiceRequestPriorityMedium // Default
	}

	// Generate request number
	requestNumber, err := s.serviceRequestRepo.GenerateRequestNumber(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate request number: %w", err)
	}

	// Determine SLA based on type
	var slaHours *int
	defaultSLA := 48 // Default 48 hours
	slaHours = &defaultSLA
	if validType == models.ServiceRequestTypeITRequest {
		defaultSLA = 24 // IT requests: 24 hours
		slaHours = &defaultSLA
	}

	// Create request
	request := &models.ServiceRequest{
		EmployeeID:    employee.ID,
		RequestNumber: requestNumber,
		Type:          validType,
		Category:      category,
		Subject:       subject,
		Description:   description,
		Status:        models.ServiceRequestStatusSubmitted,
		Priority:      validPriority,
		SLAHours:      slaHours,
		RequestedDate: time.Now(),
		UpdatedBy:     updatedBy,
	}

	// Set type-specific fields
	if validType == models.ServiceRequestTypeHRLetter {
		if letterType, ok := additionalData["letter_type"].(string); ok {
			request.LetterType = &letterType
		}
		if purpose, ok := additionalData["purpose"].(string); ok {
			request.Purpose = &purpose
		}
		if addressedTo, ok := additionalData["addressed_to"].(string); ok {
			request.AddressedTo = &addressedTo
		}
		if notes, ok := additionalData["additional_notes"].(string); ok {
			request.AdditionalNotes = &notes
		}
	} else {
		// IT Request or Facilities - store requested items as JSON
		if items, ok := additionalData["requested_items"].([]interface{}); ok {
			itemsJSON, _ := json.Marshal(items)
			itemsStr := string(itemsJSON)
			request.RequestedItemsJSON = &itemsStr
		}
	}

	// Calculate estimated completion date
	if slaHours != nil {
		estimatedDate := time.Now().Add(time.Duration(*slaHours) * time.Hour)
		request.EstimatedCompletionDate = &estimatedDate
	}

	// Assign to appropriate team (simplified - can be enhanced)
	if validType == models.ServiceRequestTypeHRLetter {
		assignedTo := "HR Support Team"
		request.AssignedTo = &assignedTo
	} else if validType == models.ServiceRequestTypeITRequest {
		assignedTo := "IT Support Team"
		request.AssignedTo = &assignedTo
	} else {
		assignedTo := "Facilities Team"
		request.AssignedTo = &assignedTo
	}

	if err := s.serviceRequestRepo.Create(request); err != nil {
		return nil, fmt.Errorf("failed to create service request: %w", err)
	}

	return request, nil
}

// CancelServiceRequest cancels a service request
func (s *SelfServiceService) CancelServiceRequest(userID uint, requestID uint, reason string) error {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("employee record not found for this user")
	}

	// Get request
	request, err := s.serviceRequestRepo.FindByID(requestID)
	if err != nil {
		return errors.New("service request not found")
	}

	// Verify ownership
	if request.EmployeeID != employee.ID {
		return errors.New("unauthorized: this request does not belong to you")
	}

	// Check if can be cancelled
	if request.Status == models.ServiceRequestStatusCompleted || request.Status == models.ServiceRequestStatusCancelled {
		return errors.New("cannot cancel a completed or already cancelled request")
	}

	// Cancel
	now := time.Now()
	request.Status = models.ServiceRequestStatusCancelled
	request.CancelledDate = &now
	request.CancelledReason = &reason

	if err := s.serviceRequestRepo.Update(request); err != nil {
		return fmt.Errorf("failed to cancel service request: %w", err)
	}

	return nil
}

// GetPendingHRRequests gets all pending service requests assigned to HR with employee information
func (s *SelfServiceService) GetPendingHRRequests(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]map[string]interface{}, int64, error) {
	// Build base query for HR requests
	queryFilters := map[string]interface{}{
		"type":   string(models.ServiceRequestTypeHRLetter),
		"status": string(models.ServiceRequestStatusSubmitted),
	}

	// Merge additional filters
	for key, value := range filters {
		queryFilters[key] = value
	}

	// Get all HR requests
	requests, totalCount, err := s.serviceRequestRepo.FindAll(tenantID, page, pageSize, queryFilters)
	if err != nil {
		return nil, 0, err
	}

	// Enrich requests with employee information
	enrichedRequests := make([]map[string]interface{}, 0, len(requests))

	for _, request := range requests {
		// Get employee information
		employee, err := s.employeeRepo.FindByID(request.EmployeeID)
		if err != nil {
			// Skip if employee not found but continue with other requests
			continue
		}

		// Get employee basic information
		basicInfo, err := s.basicInfoRepo.FindByEmployeeID(request.EmployeeID)
		if err != nil {
			// Skip if basic info not found but continue with other requests
			continue
		}

		// Get employment details for department information
		employmentDetails, err := s.employmentDetailsRepo.FindByEmployeeID(request.EmployeeID)
		var departmentName, positionTitle *string
		if err == nil && employmentDetails != nil {
			if employmentDetails.DepartmentID != nil {
				department, err := s.departmentRepo.FindByID(*employmentDetails.DepartmentID)
				if err == nil {
					departmentName = &department.Name
				}
			}

			if employmentDetails.PositionID != nil {
				position, err := s.positionRepo.FindByID(*employmentDetails.PositionID)
				if err == nil {
					positionTitle = &position.Title
				}
			}
		}

		// Create enriched response object
		enrichedRequest := map[string]interface{}{
			"request": request,
			"employee": map[string]interface{}{
				"id":           employee.ID,
				"employee_id":  employee.EmployeeID,
				"first_name":   basicInfo.FirstName,
				"last_name":    basicInfo.LastName,
				"full_name":    fmt.Sprintf("%s %s", basicInfo.FirstName, basicInfo.LastName),
				"email":        employee.Email,
				"department":   departmentName,
				"position":     positionTitle,
				"phone_number": basicInfo.MobileNumber,
			},
		}

		enrichedRequests = append(enrichedRequests, enrichedRequest)
	}

	return enrichedRequests, totalCount, nil
}

// ApproveServiceRequest approves a service request (HR only)
func (s *SelfServiceService) ApproveServiceRequest(requestID uint, approverUserID uint, notes *string) (*models.ServiceRequest, error) {
	// Get request
	request, err := s.serviceRequestRepo.FindByID(requestID)
	if err != nil {
		return nil, errors.New("service request not found")
	}

	// Check if request can be approved
	if request.Status != models.ServiceRequestStatusSubmitted && request.Status != models.ServiceRequestStatusInProgress {
		return nil, errors.New("only submitted or in-progress requests can be approved")
	}

	// Update request status
	request.Status = models.ServiceRequestStatusApproved
	request.UpdatedBy = &approverUserID

	// Set completion date
	now := time.Now()
	request.CompletedDate = &now

	if err := s.serviceRequestRepo.Update(request); err != nil {
		return nil, fmt.Errorf("failed to approve service request: %w", err)
	}

	return request, nil
}

// RejectServiceRequest rejects a service request (HR only)
func (s *SelfServiceService) RejectServiceRequest(requestID uint, approverUserID uint, reason string) (*models.ServiceRequest, error) {
	// Get request
	request, err := s.serviceRequestRepo.FindByID(requestID)
	if err != nil {
		return nil, errors.New("service request not found")
	}

	// Check if request can be rejected
	if request.Status != models.ServiceRequestStatusSubmitted && request.Status != models.ServiceRequestStatusInProgress {
		return nil, errors.New("only submitted or in-progress requests can be rejected")
	}

	// Update request status
	request.Status = models.ServiceRequestStatusRejected
	request.UpdatedBy = &approverUserID

	// Set completion date and rejection reason
	now := time.Now()
	request.CompletedDate = &now
	request.CancelledReason = &reason

	if err := s.serviceRequestRepo.Update(request); err != nil {
		return nil, fmt.Errorf("failed to reject service request: %w", err)
	}

	return request, nil
}

// SearchEmployees searches employees in the directory (for People Directory)
func (s *SelfServiceService) SearchEmployees(tenantID *uint, search *string, departmentID, positionID, locationID *uint, status *string, page, pageSize int) ([]models.Employee, int64, error) {
	return s.employeeRepo.SearchEmployees(tenantID, search, departmentID, positionID, locationID, status, page, pageSize)
}

// GetEmployeeDirectoryDetails gets detailed information for a specific employee (for People Directory)
func (s *SelfServiceService) GetEmployeeDirectoryDetails(employeeID string, tenantID *uint) (map[string]interface{}, error) {
	// Get employee
	employee, err := s.employeeRepo.FindByEmployeeID(employeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Tenant check removed (single-tenant)

	// Get employment details
	employmentDetails, _ := s.employmentDetailsRepo.FindByEmployeeID(employee.ID)

	// Build response
	response := map[string]interface{}{
		"id":          employee.ID,
		"employee_id": employee.EmployeeID,
		"full_name":   employee.FullName(),
		"first_name":  employee.FirstName,
		"last_name":   employee.LastName,
		"status":      string(employee.Status),
	}

	// Photo
	if employee.PhotoURL != nil {
		response["photo"] = *employee.PhotoURL
	}

	// Employment details
	if employmentDetails != nil {
		if employmentDetails.OfficialEmail != nil {
			response["official_email"] = *employmentDetails.OfficialEmail
		}
		if employmentDetails.DateOfJoining != nil {
			response["date_of_joining"] = employmentDetails.DateOfJoining.Format("2006-01-02")
		}

		// Department
		if employmentDetails.DepartmentID != nil {
			dept, _ := s.departmentRepo.FindByID(*employmentDetails.DepartmentID)
			if dept != nil {
				response["department"] = map[string]interface{}{
					"id":   dept.ID,
					"name": dept.Name,
					"code": dept.Code,
				}
			}
		}

		// Position
		if employmentDetails.PositionID != nil {
			pos, _ := s.positionRepo.FindByID(*employmentDetails.PositionID)
			if pos != nil {
				response["position"] = map[string]interface{}{
					"id":   pos.ID,
					"name": pos.Title,
					"code": pos.Code,
				}
				response["designation"] = pos.Title
			}
		}

		// Location
		if employmentDetails.LocationID != nil {
			loc, _ := s.locationRepo.FindByID(*employmentDetails.LocationID)
			if loc != nil {
				// Build address from location fields
				addressParts := []string{}
				if loc.AddressLine1 != nil && *loc.AddressLine1 != "" {
					addressParts = append(addressParts, *loc.AddressLine1)
				}
				if loc.AddressLine2 != nil && *loc.AddressLine2 != "" {
					addressParts = append(addressParts, *loc.AddressLine2)
				}
				if loc.City != nil && *loc.City != "" {
					addressParts = append(addressParts, *loc.City)
				}
				if loc.State != nil && *loc.State != "" {
					addressParts = append(addressParts, *loc.State)
				}
				if loc.Country != nil && *loc.Country != "" {
					addressParts = append(addressParts, *loc.Country)
				}
				address := strings.Join(addressParts, ", ")

				response["location"] = map[string]interface{}{
					"id":      loc.ID,
					"name":    loc.Name,
					"address": address,
				}
			}
		}

		// Reporting Manager
		if employmentDetails.ReportingManagerID != nil {
			manager, _ := s.employeeRepo.FindByID(*employmentDetails.ReportingManagerID)
			if manager != nil {
				managerEmail := ""
				if manager.WorkEmail != nil {
					managerEmail = *manager.WorkEmail
				}
				response["reporting_manager"] = map[string]interface{}{
					"id":          manager.ID,
					"employee_id": manager.EmployeeID,
					"full_name":   manager.FullName(),
					"designation": nil,
					"email":       managerEmail,
				}
			}
		}

		// Work phone
		if employmentDetails.WorkPhone != nil {
			response["work_phone"] = *employmentDetails.WorkPhone
		}
	}

	return response, nil
}
