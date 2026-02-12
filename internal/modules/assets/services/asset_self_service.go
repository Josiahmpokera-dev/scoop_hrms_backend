package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/models"
	assetRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"gorm.io/gorm"
)

// AssetSelfService handles self-service asset operations for employees
type AssetSelfService struct {
	assetRepo        *assetRepos.AssetRepository
	assetRequestRepo *assetRepos.AssetRequestRepository
	assetIssueRepo   *assetRepos.AssetIssueRepository
	employeeRepo     *employeeRepos.EmployeeRepository
}

// NewAssetSelfService creates a new asset self-service
func NewAssetSelfService() *AssetSelfService {
	return &AssetSelfService{
		assetRepo:        assetRepos.NewAssetRepository(),
		assetRequestRepo: assetRepos.NewAssetRequestRepository(),
		assetIssueRepo:   assetRepos.NewAssetIssueRepository(),
		employeeRepo:     employeeRepos.NewEmployeeRepository(),
	}
}

// GetMyAssignedAssets retrieves assets assigned to the authenticated employee
func (s *AssetSelfService) GetMyAssignedAssets(userID uint, tenantID *uint) ([]models.Asset, error) {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("employee record not found for this user")
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	// Get assets assigned to this employee
	assets, err := s.assetRepo.FindByEmployeeID(employee.EmployeeID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assigned assets: %w", err)
	}

	return assets, nil
}

// CreateAssetRequest creates a new asset request
func (s *AssetSelfService) CreateAssetRequest(userID uint, tenantID *uint, requestType, assetType, justification string, brand, model *string, priority string, replacingAssetID *uint, notes *string, updatedBy *uint) (*models.AssetRequest, error) {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("employee record not found for this user")
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	// Validate request type
	validType := models.AssetRequestType(requestType)
	if validType != models.AssetRequestTypeNew && validType != models.AssetRequestTypeReplacement && validType != models.AssetRequestTypeAdditional {
		return nil, errors.New("invalid request type. Must be 'new', 'replacement', or 'additional'")
	}

	// Validate asset type
	validAssetTypes := []string{"laptop", "mobile_phone", "desktop", "tablet", "id_card", "access_card", "vehicle", "office_equipment"}
	isValidAssetType := false
	for _, t := range validAssetTypes {
		if assetType == t {
			isValidAssetType = true
			break
		}
	}
	if !isValidAssetType {
		return nil, errors.New("invalid asset type")
	}

	// Validate priority
	validPriority := models.AssetRequestPriority(priority)
	if validPriority != models.AssetRequestPriorityLow && validPriority != models.AssetRequestPriorityMedium && validPriority != models.AssetRequestPriorityHigh && validPriority != models.AssetRequestPriorityUrgent {
		validPriority = models.AssetRequestPriorityMedium // Default
	}

	// If replacing, verify the asset exists and is assigned to the employee
	var replacingAssetCode *string
	if replacingAssetID != nil {
		asset, err := s.assetRepo.FindByID(*replacingAssetID)
		if err != nil {
			return nil, errors.New("replacing asset not found")
		}
		if asset.EmployeeID == nil || *asset.EmployeeID != employee.EmployeeID {
			return nil, errors.New("replacing asset is not assigned to you")
		}
		replacingAssetCode = &asset.AssetCode
	}

	// Generate request number
	requestNumber, err := s.assetRequestRepo.GenerateRequestNumber(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate request number: %w", err)
	}

	// Create request
	request := &models.AssetRequest{
		EmployeeID:       employee.EmployeeID,
		RequestNumber:    requestNumber,
		RequestType:      validType,
		AssetType:        assetType,
		Brand:            brand,
		Model:            model,
		Priority:         validPriority,
		Status:           models.AssetRequestStatusPending,
		Justification:    justification,
		Notes:            notes,
		ReplacingAssetID: replacingAssetID,
		ReplacingAssetCode: replacingAssetCode,
		RequestedDate:    time.Now(),
		UpdatedBy:       updatedBy,
	}

	if err := s.assetRequestRepo.Create(request); err != nil {
		return nil, fmt.Errorf("failed to create asset request: %w", err)
	}

	return request, nil
}

// ListAssetRequests retrieves asset requests for the authenticated employee
func (s *AssetSelfService) ListAssetRequests(userID uint, tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.AssetRequest, int64, error) {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errors.New("employee record not found for this user")
		}
		return nil, 0, fmt.Errorf("failed to get employee: %w", err)
	}

	// Get requests
	requests, total, err := s.assetRequestRepo.FindByEmployeeID(employee.EmployeeID, tenantID, page, pageSize, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get asset requests: %w", err)
	}

	return requests, total, err
}

// GetAssetRequestDetails retrieves details of a specific asset request
func (s *AssetSelfService) GetAssetRequestDetails(userID uint, requestID uint) (*models.AssetRequest, error) {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("employee record not found for this user")
	}

	// Get request
	request, err := s.assetRequestRepo.FindByID(requestID)
	if err != nil {
		return nil, errors.New("asset request not found")
	}

	// Verify ownership
	if request.EmployeeID != employee.EmployeeID {
		return nil, errors.New("unauthorized: this request does not belong to you")
	}

	return request, nil
}

// CancelAssetRequest cancels an asset request
func (s *AssetSelfService) CancelAssetRequest(userID uint, requestID uint, reason string) error {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("employee record not found for this user")
	}

	// Get request
	request, err := s.assetRequestRepo.FindByID(requestID)
	if err != nil {
		return errors.New("asset request not found")
	}

	// Verify ownership
	if request.EmployeeID != employee.EmployeeID {
		return errors.New("unauthorized: this request does not belong to you")
	}

	// Check if can be cancelled
	if request.Status == models.AssetRequestStatusFulfilled || request.Status == models.AssetRequestStatusCancelled {
		return errors.New("cannot cancel a fulfilled or already cancelled request")
	}

	// Cancel
	now := time.Now()
	request.Status = models.AssetRequestStatusCancelled
	request.CancelledDate = &now
	request.CancelledReason = &reason

	if err := s.assetRequestRepo.Update(request); err != nil {
		return fmt.Errorf("failed to cancel asset request: %w", err)
	}

	return nil
}

// CreateAssetIssue creates a new asset issue report
func (s *AssetSelfService) CreateAssetIssue(userID uint, tenantID *uint, assetID uint, issueType, title, description string, priority string, incidentDate *time.Time, incidentLocation, policeReportNumber *string, attachmentURLs []string, updatedBy *uint) (*models.AssetIssue, error) {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("employee record not found for this user")
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	// Get asset and verify it's assigned to the employee
	asset, err := s.assetRepo.FindByID(assetID)
	if err != nil {
		return nil, errors.New("asset not found")
	}

	if asset.EmployeeID == nil || *asset.EmployeeID != employee.EmployeeID {
		return nil, errors.New("this asset is not assigned to you")
	}

	// Validate issue type
	validType := models.AssetIssueType(issueType)
	if validType != models.AssetIssueTypeMalfunction && validType != models.AssetIssueTypeDamage && validType != models.AssetIssueTypeLost && validType != models.AssetIssueTypeStolen && validType != models.AssetIssueTypeOther {
		return nil, errors.New("invalid issue type. Must be 'malfunction', 'damage', 'lost', 'stolen', or 'other'")
	}

	// Validate priority
	validPriority := models.AssetIssuePriority(priority)
	if validPriority != models.AssetIssuePriorityLow && validPriority != models.AssetIssuePriorityMedium && validPriority != models.AssetIssuePriorityHigh && validPriority != models.AssetIssuePriorityUrgent {
		validPriority = models.AssetIssuePriorityMedium // Default
	}

	// Generate issue number
	issueNumber, err := s.assetIssueRepo.GenerateIssueNumber(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate issue number: %w", err)
	}

	// Store attachment URLs as JSON
	var attachmentURLsJSON *string
	if len(attachmentURLs) > 0 {
		urlsJSON, _ := json.Marshal(attachmentURLs)
		urlsStr := string(urlsJSON)
		attachmentURLsJSON = &urlsStr
	}

	// Create issue
	issue := &models.AssetIssue{
		EmployeeID:        employee.EmployeeID,
		IssueNumber:       issueNumber,
		AssetID:           assetID,
		AssetCode:         asset.AssetCode,
		IssueType:         validType,
		Priority:          validPriority,
		Status:            models.AssetIssueStatusReported,
		Title:             title,
		Description:       description,
		ReportedDate:      time.Now(),
		IncidentDate:      incidentDate,
		IncidentLocation:  incidentLocation,
		PoliceReportNumber: policeReportNumber,
		AttachmentURLsJSON: attachmentURLsJSON,
		UpdatedBy:        updatedBy,
	}

	if err := s.assetIssueRepo.Create(issue); err != nil {
		return nil, fmt.Errorf("failed to create asset issue: %w", err)
	}

	return issue, nil
}

// ListAssetIssues retrieves asset issues for the authenticated employee
func (s *AssetSelfService) ListAssetIssues(userID uint, tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.AssetIssue, int64, error) {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errors.New("employee record not found for this user")
		}
		return nil, 0, fmt.Errorf("failed to get employee: %w", err)
	}

	// Get issues
	issues, total, err := s.assetIssueRepo.FindByEmployeeID(employee.EmployeeID, tenantID, page, pageSize, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get asset issues: %w", err)
	}

	return issues, total, err
}

// GetAssetIssueDetails retrieves details of a specific asset issue
func (s *AssetSelfService) GetAssetIssueDetails(userID uint, issueID uint) (*models.AssetIssue, error) {
	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("employee record not found for this user")
	}

	// Get issue
	issue, err := s.assetIssueRepo.FindByID(issueID)
	if err != nil {
		return nil, errors.New("asset issue not found")
	}

	// Verify ownership
	if issue.EmployeeID != employee.EmployeeID {
		return nil, errors.New("unauthorized: this issue does not belong to you")
	}

	return issue, nil
}
