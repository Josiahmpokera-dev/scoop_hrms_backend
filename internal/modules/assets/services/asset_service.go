package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/repositories"
	employeeModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"gorm.io/gorm"
)

type AssetService struct {
	repo            *repositories.AssetRepository
	employeeRepo    *employeeRepos.EmployeeRepository
	assignmentRepo  *repositories.AssetAssignmentHistoryRepository
}

func NewAssetService() *AssetService {
	return &AssetService{
		repo:           repositories.NewAssetRepository(),
		employeeRepo:   employeeRepos.NewEmployeeRepository(),
		assignmentRepo: repositories.NewAssetAssignmentHistoryRepository(),
	}
}

func (s *AssetService) GetOverview(req *models.AssetOverviewRequest, tenantID *uint) (*models.AssetOverviewResponse, error) {
	filters := map[string]interface{}{}
	if req.AssetType != nil && *req.AssetType != "" {
		filters["asset_type"] = strings.ToLower(strings.TrimSpace(*req.AssetType))
	}
	if req.Status != nil && *req.Status != "" {
		filters["status"] = strings.ToLower(strings.TrimSpace(*req.Status))
	}
	if req.Department != nil && *req.Department != "" {
		filters["department"] = strings.TrimSpace(*req.Department)
	}

	row, err := s.repo.Overview(tenantID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset overview: %w", err)
	}

	currency := "TZS"
	if req.Currency != nil && strings.TrimSpace(*req.Currency) != "" {
		currency = strings.TrimSpace(*req.Currency)
	}

	return &models.AssetOverviewResponse{
		Filters: models.AssetOverviewFilters{
			AssetType:  req.AssetType,
			Department: req.Department,
			Status:     req.Status,
			Currency:   &currency,
		},
		Totals: models.AssetOverviewTotals{
			TotalAssets: row.TotalAssets,
			InUse:       row.InUse,
			Available:   row.Available,
			UnderRepair: row.UnderRepair,
			Retired:     row.Retired,
			TotalValue:  row.TotalValue,
		},
		Meta: models.AssetOverviewMeta{
			ValueUnit: "decimal",
			Currency:  currency,
		},
	}, nil
}

// GetAssetTypePrefix returns the prefix for asset code generation
func (s *AssetService) GetAssetTypePrefix(assetType string) string {
	prefixMap := map[string]string{
		"laptop":          "LAP",
		"mobile_phone":     "PHN",
		"desktop":          "DSK",
		"tablet":           "TAB",
		"id_card":         "IDC",
		"access_card":      "ACC",
		"vehicle":          "VEH",
		"office_equipment": "EQP",
	}
	if prefix, ok := prefixMap[strings.ToLower(assetType)]; ok {
		return prefix
	}
	return "AST" // Default prefix
}

// GenerateAssetCode generates a unique asset code
func (s *AssetService) GenerateAssetCode(assetType string, tenantID *uint) (string, error) {
	prefix := s.GetAssetTypePrefix(assetType)
	nextNumber, err := s.repo.GetNextAssetCodeNumber(prefix, tenantID)
	if err != nil {
		return "", fmt.Errorf("failed to generate asset code: %w", err)
	}
	return fmt.Sprintf("%s-%03d", prefix, nextNumber), nil
}

// CreateAsset creates a new asset
func (s *AssetService) CreateAsset(req *models.CreateAssetRequest, tenantID *uint, updatedBy *uint) (*models.Asset, error) {
	// Validate asset type
	validTypes := []string{"laptop", "mobile_phone", "desktop", "tablet", "id_card", "access_card", "vehicle", "office_equipment"}
	isValidType := false
	for _, t := range validTypes {
		if strings.ToLower(req.AssetType) == t {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return nil, errors.New("invalid asset type")
	}

	// Check if serial number already exists
	if s.repo.ExistsBySerialNumber(req.SerialNumber) {
		return nil, errors.New("asset with this serial number already exists")
	}

	// Generate asset code if not provided
	assetCode := ""
	if req.AssetCode != nil && *req.AssetCode != "" {
		// Check if provided asset code already exists
		if s.repo.ExistsByAssetCode(*req.AssetCode) {
			return nil, errors.New("asset with this code already exists")
		}
		assetCode = *req.AssetCode
	} else {
		// Auto-generate asset code
		generatedCode, err := s.GenerateAssetCode(req.AssetType, tenantID)
		if err != nil {
			return nil, err
		}
		assetCode = generatedCode
	}

	// Parse dates
	var purchaseDate *time.Time
	if req.PurchaseDate != nil && *req.PurchaseDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.PurchaseDate)
		if err != nil {
			return nil, errors.New("invalid purchase_date format, expected YYYY-MM-DD")
		}
		purchaseDate = &parsed
	}

	var warrantyExpiry *time.Time
	if req.WarrantyExpiry != nil && *req.WarrantyExpiry != "" {
		parsed, err := time.Parse("2006-01-02", *req.WarrantyExpiry)
		if err != nil {
			return nil, errors.New("invalid warranty_expiry format, expected YYYY-MM-DD")
		}
		warrantyExpiry = &parsed
	}

	// Set default condition
	condition := "excellent"
	if req.Condition != nil && *req.Condition != "" {
		validConditions := []string{"excellent", "good", "fair", "poor"}
		isValidCondition := false
		for _, c := range validConditions {
			if strings.ToLower(*req.Condition) == c {
				isValidCondition = true
				condition = strings.ToLower(*req.Condition)
				break
			}
		}
		if !isValidCondition {
			return nil, errors.New("invalid condition, must be one of: excellent, good, fair, poor")
		}
	}

	// Create asset
	asset := &models.Asset{
		AssetCode:     assetCode,
		AssetType:     strings.ToLower(req.AssetType),
		Brand:          req.Brand,
		Model:          req.Model,
		SerialNumber:   req.SerialNumber,
		Status:         string(models.AssetStatusAvailable),
		Condition:      condition,
		PurchaseDate:   purchaseDate,
		WarrantyExpiry: warrantyExpiry,
		Value:          req.Value,
		Notes:          req.Notes,
		UpdatedBy:      updatedBy,
	}

	// If employee ID is provided, assign immediately
	if req.AssignedToEmployeeID != nil {
		var employee *employeeModels.Employee
		var err error

		// Handle both numeric ID (database ID) and string ID (employee_id)
		if req.AssignedToEmployeeID.IsNumeric {
			// Look up by database ID
			employee, err = s.employeeRepo.FindByID(req.AssignedToEmployeeID.NumericID)
		} else {
			// Look up by employee_id string (e.g., "EMP001")
			if req.AssignedToEmployeeID.StringID == "" {
				// Empty string means no assignment
				employee = nil
			} else {
				employee, err = s.employeeRepo.FindByEmployeeID(req.AssignedToEmployeeID.StringID)
			}
		}

		if err != nil {
			return nil, errors.New("employee not found")
		}

		if employee != nil {
			if !employee.IsActive {
				return nil, errors.New("employee is not active")
			}

			// Assign asset
			now := time.Now()
			fullName := fmt.Sprintf("%s %s", employee.FirstName, employee.LastName)
			asset.AssignedTo = &fullName
			asset.EmployeeID = &employee.EmployeeID
			if employee.DepartmentID != nil {
				// Store department ID as string, handler will convert to name
				deptIDStr := fmt.Sprintf("%d", *employee.DepartmentID)
				asset.Department = &deptIDStr
			}
			asset.AssignedDate = &now
			asset.Status = string(models.AssetStatusInUse)
		}
	}

	if err := s.repo.Create(asset); err != nil {
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}

	return asset, nil
}

// UpdateAsset updates an asset
func (s *AssetService) UpdateAsset(req *models.UpdateAssetRequest, tenantID *uint, updatedBy *uint) (*models.Asset, error) {
	asset, err := s.repo.FindByID(req.ID)
	if err != nil {
		return nil, errors.New("asset not found")
	}


	// Check if asset is retired (cannot be updated)
	if asset.Status == string(models.AssetStatusRetired) {
		return nil, errors.New("retired assets cannot be updated")
	}

	// Update fields if provided
	if req.AssetType != nil {
		validTypes := []string{"laptop", "mobile_phone", "desktop", "tablet", "id_card", "access_card", "vehicle", "office_equipment"}
		isValidType := false
		for _, t := range validTypes {
			if strings.ToLower(*req.AssetType) == t {
				isValidType = true
				asset.AssetType = strings.ToLower(*req.AssetType)
				break
			}
		}
		if !isValidType {
			return nil, errors.New("invalid asset type")
		}
	}

	if req.Brand != nil {
		asset.Brand = *req.Brand
	}

	if req.Model != nil {
		asset.Model = *req.Model
	}

	if req.SerialNumber != nil {
		// Check if serial number already exists (excluding current asset)
		existing, err := s.repo.FindBySerialNumber(*req.SerialNumber)
		if err == nil && existing.ID != asset.ID {
			return nil, errors.New("asset with this serial number already exists")
		}
		asset.SerialNumber = *req.SerialNumber
	}

	if req.Condition != nil {
		validConditions := []string{"excellent", "good", "fair", "poor"}
		isValidCondition := false
		for _, c := range validConditions {
			if strings.ToLower(*req.Condition) == c {
				isValidCondition = true
				asset.Condition = strings.ToLower(*req.Condition)
				break
			}
		}
		if !isValidCondition {
			return nil, errors.New("invalid condition, must be one of: excellent, good, fair, poor")
		}
	}

	if req.Value != nil {
		asset.Value = req.Value
	}

	if req.PurchaseDate != nil && *req.PurchaseDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.PurchaseDate)
		if err != nil {
			return nil, errors.New("invalid purchase_date format, expected YYYY-MM-DD")
		}
		asset.PurchaseDate = &parsed
	}

	if req.WarrantyExpiry != nil && *req.WarrantyExpiry != "" {
		parsed, err := time.Parse("2006-01-02", *req.WarrantyExpiry)
		if err != nil {
			return nil, errors.New("invalid warranty_expiry format, expected YYYY-MM-DD")
		}
		asset.WarrantyExpiry = &parsed
	}

	if req.Notes != nil {
		asset.Notes = req.Notes
	}

	asset.UpdatedBy = updatedBy
	asset.UpdatedAt = time.Now()

	if err := s.repo.Update(asset); err != nil {
		return nil, fmt.Errorf("failed to update asset: %w", err)
	}

	return asset, nil
}

// AssignAsset assigns an asset to an employee
func (s *AssetService) AssignAsset(req *models.AssignAssetRequest, tenantID *uint, updatedBy *uint) (*models.Asset, error) {
	asset, err := s.repo.FindByID(req.AssetID)
	if err != nil {
		return nil, errors.New("asset not found")
	}


	// Check if asset is available
	if !asset.CanBeAssigned() {
		return nil, errors.New("asset is not available for assignment")
	}

	// Find employee
	employee, err := s.employeeRepo.FindByEmployeeID(req.EmployeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	if !employee.IsActive {
		return nil, errors.New("employee is not active")
	}

	// Parse assigned date
	assignedDate := time.Now()
	if req.AssignedDate != nil && *req.AssignedDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.AssignedDate)
		if err != nil {
			return nil, errors.New("invalid assigned_date format, expected YYYY-MM-DD")
		}
		assignedDate = parsed
	}

	// Assign asset
	fullName := fmt.Sprintf("%s %s", employee.FirstName, employee.LastName)
	asset.AssignedTo = &fullName
	asset.EmployeeID = &employee.EmployeeID
	if employee.DepartmentID != nil {
		// Get department name - store department ID as string, handler will convert to name
		deptIDStr := fmt.Sprintf("%d", *employee.DepartmentID)
		asset.Department = &deptIDStr
	}
	asset.AssignedDate = &assignedDate
	asset.Status = string(models.AssetStatusInUse)
	if req.Notes != nil {
		asset.Notes = req.Notes
	}
	asset.UpdatedBy = updatedBy
	asset.UpdatedAt = time.Now()

	if err := s.repo.Update(asset); err != nil {
		return nil, fmt.Errorf("failed to assign asset: %w", err)
	}
	fromEmployeeName := ""
	toEmployeeName := ""
	if asset.AssignedTo != nil {
		toEmployeeName = *asset.AssignedTo
	}
	if err := s.assignmentRepo.Create(&models.AssetAssignmentHistory{
		TenantID:       tenantID,
		AssetID:        asset.ID,
		AssetCode:      asset.AssetCode,
		Action:         models.AssetAssignmentActionAssign,
		FromEmployeeID: nil,
		ToEmployeeID:   asset.EmployeeID,
		FromEmployee:   &fromEmployeeName,
		ToEmployee:     &toEmployeeName,
		Department:     asset.Department,
		Notes:          req.Notes,
		OccurredAt:     assignedDate,
		ActorUserID:    updatedBy,
	}); err != nil {
		return nil, fmt.Errorf("failed to write assignment history: %w", err)
	}

	return asset, nil
}

// ReturnAsset returns an asset from an employee
func (s *AssetService) ReturnAsset(req *models.ReturnAssetRequest, tenantID *uint, updatedBy *uint) (*models.Asset, error) {
	asset, err := s.repo.FindByID(req.AssetID)
	if err != nil {
		return nil, errors.New("asset not found")
	}


	// Check if asset can be returned
	if !asset.CanBeReturned() {
		return nil, errors.New("asset is not currently assigned")
	}

	// Parse return date
	returnDate := time.Now()
	if req.ReturnDate != nil && *req.ReturnDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.ReturnDate)
		if err != nil {
			return nil, errors.New("invalid return_date format, expected YYYY-MM-DD")
		}
		returnDate = parsed
	}

	// Return asset
	prevEmployeeID := asset.EmployeeID
	prevEmployeeName := asset.AssignedTo
	asset.ReturnDate = &returnDate
	asset.AssignedTo = nil
	asset.EmployeeID = nil
	asset.Department = nil
	asset.AssignedDate = nil
	asset.Status = string(models.AssetStatusAvailable)

	if req.Condition != nil {
		validConditions := []string{"excellent", "good", "fair", "poor"}
		isValidCondition := false
		for _, c := range validConditions {
			if strings.ToLower(*req.Condition) == c {
				isValidCondition = true
				asset.Condition = strings.ToLower(*req.Condition)
				break
			}
		}
		if !isValidCondition {
			return nil, errors.New("invalid condition, must be one of: excellent, good, fair, poor")
		}
	}

	if req.Notes != nil {
		asset.Notes = req.Notes
	}
	asset.UpdatedBy = updatedBy
	asset.UpdatedAt = time.Now()

	if err := s.repo.Update(asset); err != nil {
		return nil, fmt.Errorf("failed to return asset: %w", err)
	}
	blankName := ""
	if err := s.assignmentRepo.Create(&models.AssetAssignmentHistory{
		TenantID:       tenantID,
		AssetID:        asset.ID,
		AssetCode:      asset.AssetCode,
		Action:         models.AssetAssignmentActionReturn,
		FromEmployeeID: prevEmployeeID,
		ToEmployeeID:   nil,
		FromEmployee:   prevEmployeeName,
		ToEmployee:     &blankName,
		Department:     nil,
		Notes:          req.Notes,
		OccurredAt:     returnDate,
		ActorUserID:    updatedBy,
	}); err != nil {
		return nil, fmt.Errorf("failed to write return history: %w", err)
	}

	return asset, nil
}

// MarkForRepair marks an asset for repair
func (s *AssetService) MarkForRepair(req *models.MarkForRepairRequest, tenantID *uint, updatedBy *uint) (*models.Asset, error) {
	asset, err := s.repo.FindByID(req.AssetID)
	if err != nil {
		return nil, errors.New("asset not found")
	}


	// Check if asset can be repaired
	if !asset.CanBeRepaired() {
		return nil, errors.New("asset cannot be marked for repair (already under repair or retired)")
	}

	// Mark for repair
	now := time.Now()
	asset.RepairDate = &now
	asset.Status = string(models.AssetStatusUnderRepair)
	if req.RepairReason != nil {
		asset.RepairReason = req.RepairReason
	}
	if req.RepairNotes != nil {
		asset.RepairNotes = req.RepairNotes
	}
	asset.UpdatedBy = updatedBy
	asset.UpdatedAt = time.Now()

	if err := s.repo.Update(asset); err != nil {
		return nil, fmt.Errorf("failed to mark asset for repair: %w", err)
	}

	return asset, nil
}

// CompleteRepair completes asset repair
func (s *AssetService) CompleteRepair(req *models.CompleteRepairRequest, tenantID *uint, updatedBy *uint) (*models.Asset, error) {
	asset, err := s.repo.FindByID(req.AssetID)
	if err != nil {
		return nil, errors.New("asset not found")
	}


	// Check if asset is under repair
	if asset.Status != string(models.AssetStatusUnderRepair) {
		return nil, errors.New("asset is not under repair")
	}

	// Complete repair
	now := time.Now()
	asset.RepairCompletedDate = &now
	asset.Status = string(models.AssetStatusAvailable) // Return to available status

	if req.RepairCost != nil {
		asset.RepairCost = req.RepairCost
	}

	if req.Condition != nil {
		validConditions := []string{"excellent", "good", "fair", "poor"}
		isValidCondition := false
		for _, c := range validConditions {
			if strings.ToLower(*req.Condition) == c {
				isValidCondition = true
				asset.Condition = strings.ToLower(*req.Condition)
				break
			}
		}
		if !isValidCondition {
			return nil, errors.New("invalid condition, must be one of: excellent, good, fair, poor")
		}
	}

	if req.Notes != nil {
		asset.Notes = req.Notes
	}
	asset.UpdatedBy = updatedBy
	asset.UpdatedAt = time.Now()

	if err := s.repo.Update(asset); err != nil {
		return nil, fmt.Errorf("failed to complete repair: %w", err)
	}

	return asset, nil
}

// RetireAsset retires an asset
func (s *AssetService) RetireAsset(req *models.RetireAssetRequest, tenantID *uint, updatedBy *uint) (*models.Asset, error) {
	asset, err := s.repo.FindByID(req.AssetID)
	if err != nil {
		return nil, errors.New("asset not found")
	}


	// Check if asset can be retired
	if !asset.CanBeRetired() {
		return nil, errors.New("asset is already retired")
	}

	// Retire asset
	retirementDate := time.Now()
	if req.RetirementDate != nil && *req.RetirementDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.RetirementDate)
		if err != nil {
			return nil, errors.New("invalid retirement_date format, expected YYYY-MM-DD")
		}
		retirementDate = parsed
	}

	asset.RetirementDate = &retirementDate
	asset.Status = string(models.AssetStatusRetired)
	if req.RetirementReason != nil {
		asset.RetirementReason = req.RetirementReason
	}
	if req.Notes != nil {
		asset.Notes = req.Notes
	}
	asset.UpdatedBy = updatedBy
	asset.UpdatedAt = time.Now()

	// Clear assignment if any
	asset.AssignedTo = nil
	asset.EmployeeID = nil
	asset.Department = nil
	asset.AssignedDate = nil

	if err := s.repo.Update(asset); err != nil {
		return nil, fmt.Errorf("failed to retire asset: %w", err)
	}

	return asset, nil
}

// DeleteAsset permanently deletes an asset
func (s *AssetService) DeleteAsset(id uint, tenantID *uint) error {
	asset, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("asset not found")
	}

	// Check if asset can be deleted (only retired assets)
	if !asset.CanBeDeleted() {
		return errors.New("only retired assets can be deleted")
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}

	return nil
}

// GetAsset retrieves an asset by ID
func (s *AssetService) GetAsset(id uint, tenantID *uint) (*models.Asset, error) {
	asset, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("asset not found")
		}
		return nil, err
	}


	return asset, nil
}

// ListAssets lists assets with pagination and filters
func (s *AssetService) ListAssets(req *models.ListAssetsRequest, tenantID *uint) ([]models.Asset, int64, error) {
	// Set defaults
	page := 1
	if req.Page > 0 {
		page = req.Page
	}

	pageSize := 20
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}

	// Build filters
	filters := make(map[string]interface{})
	if req.Search != nil && *req.Search != "" {
		filters["search"] = *req.Search
	}
	if req.Status != nil && *req.Status != "" {
		filters["status"] = *req.Status
	}
	if req.AssetType != nil && *req.AssetType != "" {
		filters["asset_type"] = *req.AssetType
	}
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		filters["assigned_to"] = *req.AssignedTo
	}
	if req.Department != nil && *req.Department != "" {
		filters["department"] = *req.Department
	}

	return s.repo.List(tenantID, page, pageSize, filters)
}

// GetAssetTypes returns the list of available asset types
func (s *AssetService) GetAssetTypes() []models.AssetTypeResponse {
	return []models.AssetTypeResponse{
		{ID: 1, Name: "Laptop", Code: "laptop", Description: "Portable computers"},
		{ID: 2, Name: "Mobile Phone", Code: "mobile_phone", Description: "Mobile communication devices"},
		{ID: 3, Name: "Desktop", Code: "desktop", Description: "Desktop computers"},
		{ID: 4, Name: "Tablet", Code: "tablet", Description: "Tablet devices"},
		{ID: 5, Name: "ID Card", Code: "id_card", Description: "Employee identification cards"},
		{ID: 6, Name: "Access Card", Code: "access_card", Description: "Access control cards"},
		{ID: 7, Name: "Vehicle", Code: "vehicle", Description: "Company vehicles"},
		{ID: 8, Name: "Office Equipment", Code: "office_equipment", Description: "Office equipment and supplies"},
	}
}
