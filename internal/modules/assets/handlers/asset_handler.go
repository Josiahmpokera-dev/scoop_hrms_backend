package handlers

import (
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/services"
	departmentRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	service       *services.AssetService
	employeeRepo  *employeeRepos.EmployeeRepository
	departmentRepo *departmentRepos.DepartmentRepository
}

func NewAssetHandler() *AssetHandler {
	return &AssetHandler{
		service:        services.NewAssetService(),
		employeeRepo:   employeeRepos.NewEmployeeRepository(),
		departmentRepo: departmentRepos.NewDepartmentRepository(),
	}
}

// toAssetResponse converts Asset model to AssetResponse with employee photo and department name
func (h *AssetHandler) toAssetResponse(asset *models.Asset) *models.AssetResponse {
	resp := &models.AssetResponse{
		ID:                  asset.ID,
		AssetCode:           asset.AssetCode,
		AssetType:           asset.AssetType,
		Brand:               asset.Brand,
		Model:               asset.Model,
		SerialNumber:        asset.SerialNumber,
		AssignedTo:          asset.AssignedTo,
		EmployeeID:          asset.EmployeeID,
		Department:         asset.Department, // This is stored as department ID reference, we'll convert it
		AssignedDate:        asset.AssignedDate,
		ReturnDate:          asset.ReturnDate,
		Status:              asset.Status,
		Condition:           asset.Condition,
		PurchaseDate:        asset.PurchaseDate,
		WarrantyExpiry:      asset.WarrantyExpiry,
		Value:               asset.Value,
		Notes:               asset.Notes,
		RepairReason:        asset.RepairReason,
		RepairNotes:         asset.RepairNotes,
		RepairDate:          asset.RepairDate,
		RepairCost:          asset.RepairCost,
		RepairCompletedDate: asset.RepairCompletedDate,
		RetirementReason:    asset.RetirementReason,
		RetirementDate:      asset.RetirementDate,
		CreatedAt:           asset.CreatedAt,
		UpdatedAt:           asset.UpdatedAt,
	}

	// Get employee photo if employee ID exists
	if asset.EmployeeID != nil {
		employee, err := h.employeeRepo.FindByEmployeeID(*asset.EmployeeID)
		if err == nil && employee != nil {
			if employee.PhotoURL != nil && *employee.PhotoURL != "" {
				resp.EmployeePhoto = employee.PhotoURL
			} else {
				defaultPhoto := "/img/avatars/default.jpg"
				resp.EmployeePhoto = &defaultPhoto
			}
		}
	}

	// Convert department ID to department name
	// Note: In the asset model, Department is stored as *string (department ID as string)
	// We need to convert it to department name
	if asset.Department != nil {
		// Try to parse as uint (department ID)
		if deptID, err := strconv.ParseUint(*asset.Department, 10, 32); err == nil {
			dept, err := h.departmentRepo.FindByID(uint(deptID))
			if err == nil && dept != nil {
				resp.Department = &dept.Name
			}
		} else {
			// If it's already a name, use it directly
			resp.Department = asset.Department
		}
	}

	return resp
}

// ListAssets handles listing assets with pagination and filters
func (h *AssetHandler) ListAssets(c *gin.Context) {
	var req models.ListAssetsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	assets, total, err := h.service.ListAssets(&req, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Convert to response format
	assetResponses := make([]models.AssetResponse, len(assets))
	for i, asset := range assets {
		assetResponses[i] = *h.toAssetResponse(&asset)
	}

	// Calculate pagination metadata
	page := req.Page
	if page == 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 20
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	response.SuccessWithMeta(c, "Assets retrieved successfully", assetResponses, &response.Meta{
		Page:       page,
		PerPage:    pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h *AssetHandler) GetAssetOverview(c *gin.Context) {
	var req models.AssetOverviewRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	data, err := h.service.GetOverview(&req, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Asset overview retrieved successfully", data)
}

// GetAsset handles getting an asset by ID
func (h *AssetHandler) GetAsset(c *gin.Context) {
	var req models.GetAssetRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		// Try to get ID from query parameter
		idStr := c.Query("id")
		if idStr == "" {
			response.ValidationError(c, "Validation failed", "id is required")
			return
		}
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.ValidationError(c, "Validation failed", "invalid id format")
			return
		}
		req.ID = uint(id)
	} else if req.ID == 0 {
		response.ValidationError(c, "Validation failed", "id is required")
		return
	}

	tenantID := middleware.GetTenantID(c)
	asset, err := h.service.GetAsset(req.ID, tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Asset retrieved successfully", h.toAssetResponse(asset))
}

// CreateAsset handles asset creation
func (h *AssetHandler) CreateAsset(c *gin.Context) {
	var req models.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if user != nil {
		if u, ok := user.(*userModels.User); ok {
			updatedBy = &u.ID
		}
	}

	asset, err := h.service.CreateAsset(&req, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Asset created successfully", h.toAssetResponse(asset))
}

// UpdateAsset handles asset update
func (h *AssetHandler) UpdateAsset(c *gin.Context) {
	var req models.UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if user != nil {
		if u, ok := user.(*userModels.User); ok {
			updatedBy = &u.ID
		}
	}

	asset, err := h.service.UpdateAsset(&req, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Asset updated successfully", h.toAssetResponse(asset))
}

// AssignAsset handles asset assignment to employee
func (h *AssetHandler) AssignAsset(c *gin.Context) {
	var req models.AssignAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if user != nil {
		if u, ok := user.(*userModels.User); ok {
			updatedBy = &u.ID
		}
	}

	asset, err := h.service.AssignAsset(&req, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Asset assigned successfully", h.toAssetResponse(asset))
}

// ReturnAsset handles asset return
func (h *AssetHandler) ReturnAsset(c *gin.Context) {
	var req models.ReturnAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if user != nil {
		if u, ok := user.(*userModels.User); ok {
			updatedBy = &u.ID
		}
	}

	asset, err := h.service.ReturnAsset(&req, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Asset returned successfully", h.toAssetResponse(asset))
}

// ReassignAsset handles asset reassignment to another employee
func (h *AssetHandler) ReassignAsset(c *gin.Context) {
	var req models.ReassignAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if user != nil {
		if u, ok := user.(*userModels.User); ok {
			updatedBy = &u.ID
		}
	}

	asset, err := h.service.ReassignAsset(&req, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Asset reassigned successfully", h.toAssetResponse(asset))
}

// MarkForRepair handles marking asset for repair
func (h *AssetHandler) MarkForRepair(c *gin.Context) {
	var req models.MarkForRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if user != nil {
		if u, ok := user.(*userModels.User); ok {
			updatedBy = &u.ID
		}
	}

	asset, err := h.service.MarkForRepair(&req, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Asset marked for repair successfully", h.toAssetResponse(asset))
}

// CompleteRepair handles completing asset repair
func (h *AssetHandler) CompleteRepair(c *gin.Context) {
	var req models.CompleteRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if user != nil {
		if u, ok := user.(*userModels.User); ok {
			updatedBy = &u.ID
		}
	}

	asset, err := h.service.CompleteRepair(&req, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Asset repair completed successfully", h.toAssetResponse(asset))
}

// RetireAsset handles asset retirement
func (h *AssetHandler) RetireAsset(c *gin.Context) {
	var req models.RetireAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var updatedBy *uint
	if user != nil {
		if u, ok := user.(*userModels.User); ok {
			updatedBy = &u.ID
		}
	}

	asset, err := h.service.RetireAsset(&req, tenantID, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Asset retired successfully", h.toAssetResponse(asset))
}

// DeleteAsset handles asset deletion
func (h *AssetHandler) DeleteAsset(c *gin.Context) {
	var req models.DeleteAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	if err := h.service.DeleteAsset(req.ID, tenantID); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Asset deleted successfully", nil)
}

// GetAssetTypes handles getting asset types
func (h *AssetHandler) GetAssetTypes(c *gin.Context) {
	assetTypes := h.service.GetAssetTypes()
	response.Success(c, "Asset types retrieved successfully", assetTypes)
}
