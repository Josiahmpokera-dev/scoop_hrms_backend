package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/models"
	assetRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
)

// AssetHRService handles HR/Admin asset management operations
type AssetHRService struct {
	assetRepo        *assetRepos.AssetRepository
	assetRequestRepo *assetRepos.AssetRequestRepository
	assetIssueRepo   *assetRepos.AssetIssueRepository
	employeeRepo     *employeeRepos.EmployeeRepository
	assetService     *AssetService
}

// NewAssetHRService creates a new asset HR service
func NewAssetHRService() *AssetHRService {
	return &AssetHRService{
		assetRepo:        assetRepos.NewAssetRepository(),
		assetRequestRepo: assetRepos.NewAssetRequestRepository(),
		assetIssueRepo:   assetRepos.NewAssetIssueRepository(),
		employeeRepo:     employeeRepos.NewEmployeeRepository(),
		assetService:     NewAssetService(),
	}
}

// ListAssetRequests lists all asset requests (for HR/Admin)
func (s *AssetHRService) ListAssetRequests(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.AssetRequest, int64, error) {
	requests, total, err := s.assetRequestRepo.ListAll(tenantID, page, pageSize, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get asset requests: %w", err)
	}
	return requests, total, err
}

// GetAssetRequestDetails gets details of a specific asset request
func (s *AssetHRService) GetAssetRequestDetails(requestID uint, tenantID *uint) (*models.AssetRequest, error) {
	request, err := s.assetRequestRepo.FindByID(requestID)
	if err != nil {
		return nil, errors.New("asset request not found")
	}


	return request, nil
}

// ApproveAssetRequest approves an asset request
func (s *AssetHRService) ApproveAssetRequest(requestID uint, tenantID *uint, approvedBy *uint) (*models.AssetRequest, error) {
	request, err := s.assetRequestRepo.FindByID(requestID)
	if err != nil {
		return nil, errors.New("asset request not found")
	}


	// Check if can be approved
	if request.Status != models.AssetRequestStatusPending {
		return nil, fmt.Errorf("cannot approve request with status: %s", request.Status)
	}

	// Approve
	now := time.Now()
	request.Status = models.AssetRequestStatusApproved
	request.ApprovedBy = approvedBy
	request.ApprovedAt = &now

	if err := s.assetRequestRepo.Update(request); err != nil {
		return nil, fmt.Errorf("failed to approve asset request: %w", err)
	}

	return request, nil
}

// RejectAssetRequest rejects an asset request
func (s *AssetHRService) RejectAssetRequest(requestID uint, tenantID *uint, rejectionReason string, rejectedBy *uint) (*models.AssetRequest, error) {
	request, err := s.assetRequestRepo.FindByID(requestID)
	if err != nil {
		return nil, errors.New("asset request not found")
	}


	// Check if can be rejected
	if request.Status != models.AssetRequestStatusPending {
		return nil, fmt.Errorf("cannot reject request with status: %s", request.Status)
	}

	// Reject
	now := time.Now()
	request.Status = models.AssetRequestStatusRejected
	request.RejectedAt = &now
	request.RejectionReason = &rejectionReason
	request.UpdatedBy = rejectedBy

	if err := s.assetRequestRepo.Update(request); err != nil {
		return nil, fmt.Errorf("failed to reject asset request: %w", err)
	}

	return request, nil
}

// FulfillAssetRequest fulfills an approved asset request by assigning an asset
func (s *AssetHRService) FulfillAssetRequest(requestID uint, assetID uint, tenantID *uint, fulfilledBy *uint) (*models.AssetRequest, *models.Asset, error) {
	request, err := s.assetRequestRepo.FindByID(requestID)
	if err != nil {
		return nil, nil, errors.New("asset request not found")
	}

	// Check if can be fulfilled
	if request.Status != models.AssetRequestStatusApproved {
		return nil, nil, fmt.Errorf("can only fulfill approved requests. Current status: %s", request.Status)
	}

	// Get asset
	asset, err := s.assetRepo.FindByID(assetID)
	if err != nil {
		return nil, nil, errors.New("asset not found")
	}


	// Check if asset is available
	if !asset.CanBeAssigned() {
		return nil, nil, errors.New("asset is not available for assignment")
	}

	// Get employee
	employee, err := s.employeeRepo.FindByEmployeeID(request.EmployeeID)
	if err != nil {
		return nil, nil, errors.New("employee not found")
	}

	if !employee.IsActive {
		return nil, nil, errors.New("employee is not active")
	}

	// Assign asset to employee
	assignedDate := time.Now()
	fullName := fmt.Sprintf("%s %s", employee.FirstName, employee.LastName)
	asset.AssignedTo = &fullName
	asset.EmployeeID = &employee.EmployeeID
	if employee.DepartmentID != nil {
		deptIDStr := fmt.Sprintf("%d", *employee.DepartmentID)
		asset.Department = &deptIDStr
	}
	asset.AssignedDate = &assignedDate
	asset.Status = string(models.AssetStatusInUse)
	asset.UpdatedBy = fulfilledBy
	asset.UpdatedAt = time.Now()

	if err := s.assetRepo.Update(asset); err != nil {
		return nil, nil, fmt.Errorf("failed to assign asset: %w", err)
	}

	// Mark request as fulfilled
	now := time.Now()
	request.Status = models.AssetRequestStatusFulfilled
	request.FulfilledAt = &now
	request.AssignedAssetID = &asset.ID
	assetCode := asset.AssetCode
	request.AssignedAssetCode = &assetCode
	request.UpdatedBy = fulfilledBy

	if err := s.assetRequestRepo.Update(request); err != nil {
		return nil, nil, fmt.Errorf("failed to update request status: %w", err)
	}

	return request, asset, nil
}

// ReassignAsset reassigns an asset from one employee to another
func (s *AssetHRService) ReassignAsset(assetID uint, newEmployeeID string, tenantID *uint, reassignedBy *uint, notes *string) (*models.Asset, error) {
	// Get asset
	asset, err := s.assetRepo.FindByID(assetID)
	if err != nil {
		return nil, errors.New("asset not found")
	}

	// Check if asset is currently assigned
	if !asset.IsAssigned() {
		return nil, errors.New("asset is not currently assigned to any employee")
	}

	// Check if reassigning to the same employee
	if asset.EmployeeID != nil && *asset.EmployeeID == newEmployeeID {
		return nil, errors.New("asset is already assigned to this employee")
	}

	// Get new employee
	newEmployee, err := s.employeeRepo.FindByEmployeeID(newEmployeeID)
	if err != nil {
		return nil, errors.New("new employee not found")
	}

	if !newEmployee.IsActive {
		return nil, errors.New("new employee is not active")
	}

	// Reassign asset
	reassignDate := time.Now()
	fullName := fmt.Sprintf("%s %s", newEmployee.FirstName, newEmployee.LastName)
	asset.AssignedTo = &fullName
	asset.EmployeeID = &newEmployee.EmployeeID
	if newEmployee.DepartmentID != nil {
		deptIDStr := fmt.Sprintf("%d", *newEmployee.DepartmentID)
		asset.Department = &deptIDStr
	}
	asset.AssignedDate = &reassignDate
	asset.UpdatedBy = reassignedBy
	asset.UpdatedAt = time.Now()
	if notes != nil {
		asset.Notes = notes
	}

	if err := s.assetRepo.Update(asset); err != nil {
		return nil, fmt.Errorf("failed to reassign asset: %w", err)
	}

	return asset, nil
}
