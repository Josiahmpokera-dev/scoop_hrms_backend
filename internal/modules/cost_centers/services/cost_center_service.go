package services

import (
	"errors"
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/cost_centers/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/cost_centers/repositories"
	departmentRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
	organizationRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organizations/repositories"
)

type CostCenterService struct {
	repo           *repositories.CostCenterRepository
	organizationRepo *organizationRepos.OrganizationRepository
	departmentRepo *departmentRepos.DepartmentRepository
}

func NewCostCenterService() *CostCenterService {
	return &CostCenterService{
		repo:           repositories.NewCostCenterRepository(),
		organizationRepo: organizationRepos.NewOrganizationRepository(),
		departmentRepo: departmentRepos.NewDepartmentRepository(),
	}
}

// CreateCostCenter creates a new cost center
func (s *CostCenterService) CreateCostCenter(req *models.CreateCostCenterRequest, tenantID *uint, updatedBy *uint) (*models.CostCenter, error) {
	// Check if code already exists
	if s.repo.ExistsByCode(req.Code) {
		return nil, errors.New("cost center with this code already exists")
	}

	// Validate organization if provided
	if req.OrganizationID != nil {
		_, err := s.organizationRepo.FindByID(*req.OrganizationID)
		if err != nil {
			return nil, errors.New("organization not found")
		}
	}

	// Validate department if provided
	if req.DepartmentID != nil {
		_, err := s.departmentRepo.FindByID(*req.DepartmentID)
		if err != nil {
			return nil, errors.New("department not found")
		}
	}

	// Validate parent cost center if provided
	if req.ParentCostCenterID != nil {
		_, err := s.repo.FindByID(*req.ParentCostCenterID)
		if err != nil {
			return nil, errors.New("parent cost center not found")
		}
	}

	costCenter := &models.CostCenter{
		TenantID:          tenantID,
		OrganizationID:    req.OrganizationID,
		Code:              req.Code,
		Name:              req.Name,
		Description:       req.Description,
		CostCenterType:    req.CostCenterType,
		ParentCostCenterID: req.ParentCostCenterID,
		BudgetAllocated:   req.BudgetAllocated,
		BudgetCurrency:   req.BudgetCurrency,
		ManagerID:         req.ManagerID,
		DepartmentID:      req.DepartmentID,
		IsActive:          true,
		UpdatedBy:         updatedBy,
	}

	// Set default budget currency if not provided
	if costCenter.BudgetCurrency == nil && costCenter.BudgetAllocated != nil {
		defaultCurrency := "TZS"
		costCenter.BudgetCurrency = &defaultCurrency
	}

	if req.IsActive != nil {
		costCenter.IsActive = *req.IsActive
	}

	if err := s.repo.Create(costCenter); err != nil {
		return nil, fmt.Errorf("failed to create cost center: %w", err)
	}

	return costCenter, nil
}

// GetCostCenterByID retrieves a cost center by ID
func (s *CostCenterService) GetCostCenterByID(id uint) (*models.CostCenter, error) {
	return s.repo.FindByID(id)
}

// UpdateCostCenter updates an existing cost center
func (s *CostCenterService) UpdateCostCenter(id uint, req *models.UpdateCostCenterRequest, updatedBy *uint) (*models.CostCenter, error) {
	costCenter, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("cost center not found")
	}

	// Check code uniqueness if code is being updated
	if req.Code != nil && *req.Code != costCenter.Code {
		if s.repo.ExistsByCode(*req.Code) {
			return nil, errors.New("cost center with this code already exists")
		}
		costCenter.Code = *req.Code
	}

	// Validate organization if being updated
	if req.OrganizationID != nil {
		_, err := s.organizationRepo.FindByID(*req.OrganizationID)
		if err != nil {
			return nil, errors.New("organization not found")
		}
		costCenter.OrganizationID = req.OrganizationID
	}

	// Validate department if being updated
	if req.DepartmentID != nil {
		_, err := s.departmentRepo.FindByID(*req.DepartmentID)
		if err != nil {
			return nil, errors.New("department not found")
		}
		costCenter.DepartmentID = req.DepartmentID
	}

	// Validate parent cost center if being updated
	if req.ParentCostCenterID != nil {
		if *req.ParentCostCenterID != 0 {
			_, err := s.repo.FindByID(*req.ParentCostCenterID)
			if err != nil {
				return nil, errors.New("parent cost center not found")
			}
		}
		costCenter.ParentCostCenterID = req.ParentCostCenterID
	}

	if req.Name != nil {
		costCenter.Name = *req.Name
	}
	if req.Description != nil {
		costCenter.Description = req.Description
	}
	if req.CostCenterType != nil {
		costCenter.CostCenterType = req.CostCenterType
	}
	if req.BudgetAllocated != nil {
		costCenter.BudgetAllocated = req.BudgetAllocated
	}
	if req.BudgetCurrency != nil {
		costCenter.BudgetCurrency = req.BudgetCurrency
	}
	if req.ManagerID != nil {
		costCenter.ManagerID = req.ManagerID
	}
	if req.IsActive != nil {
		costCenter.IsActive = *req.IsActive
	}

	costCenter.UpdatedBy = updatedBy

	if err := s.repo.Update(costCenter); err != nil {
		return nil, fmt.Errorf("failed to update cost center: %w", err)
	}

	return costCenter, nil
}

// DeleteCostCenter soft deletes a cost center
func (s *CostCenterService) DeleteCostCenter(id uint) error {
	costCenter, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("cost center not found")
	}

	// Check if cost center has sub-cost centers
	if len(costCenter.SubCostCenters) > 0 {
		return errors.New("cannot delete cost center with sub-cost centers")
	}

	return s.repo.Delete(id)
}

// ListCostCenters lists cost centers with pagination and filters
func (s *CostCenterService) ListCostCenters(tenantID, organizationID *uint, page, pageSize int, filters map[string]interface{}) ([]models.CostCenter, int64, error) {
	return s.repo.List(tenantID, organizationID, page, pageSize, filters)
}

// GetRootCostCenters gets all root cost centers
func (s *CostCenterService) GetRootCostCenters(organizationID *uint) ([]models.CostCenter, error) {
	return s.repo.FindRootCostCenters(organizationID)
}
