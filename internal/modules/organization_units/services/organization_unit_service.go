package services

import (
	"errors"
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization_units/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization_units/repositories"
	organizationRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organizations/repositories"
)

type OrganizationUnitService struct {
	repo           *repositories.OrganizationUnitRepository
	organizationRepo *organizationRepos.OrganizationRepository
}

func NewOrganizationUnitService() *OrganizationUnitService {
	return &OrganizationUnitService{
		repo:           repositories.NewOrganizationUnitRepository(),
		organizationRepo: organizationRepos.NewOrganizationRepository(),
	}
}

// CreateOrganizationUnit creates a new organization unit
func (s *OrganizationUnitService) CreateOrganizationUnit(req *models.CreateOrganizationUnitRequest, tenantID *uint, updatedBy *uint) (*models.OrganizationUnit, error) {
	// Check if code already exists
	if s.repo.ExistsByCode(req.Code) {
		return nil, errors.New("organization unit with this code already exists")
	}

	// Validate organization if provided
	if req.OrganizationID != nil {
		_, err := s.organizationRepo.FindByID(*req.OrganizationID)
		if err != nil {
			return nil, errors.New("organization not found")
		}
	}

	// Validate parent unit if provided
	if req.ParentUnitID != nil {
		_, err := s.repo.FindByID(*req.ParentUnitID)
		if err != nil {
			return nil, errors.New("parent organization unit not found")
		}
	}

	unit := &models.OrganizationUnit{
		TenantID:       tenantID,
		OrganizationID: req.OrganizationID,
		Code:           req.Code,
		Name:           req.Name,
		Description:    req.Description,
		UnitType:       req.UnitType,
		ParentUnitID:   req.ParentUnitID,
		HeadID:         req.HeadID,
		BudgetAllocated: req.BudgetAllocated,
		BudgetCurrency: req.BudgetCurrency,
		LocationID:     req.LocationID,
		IsActive:       true,
		UpdatedBy:      updatedBy,
	}

	// Set default budget currency if not provided
	if unit.BudgetCurrency == nil && unit.BudgetAllocated != nil {
		defaultCurrency := "TZS"
		unit.BudgetCurrency = &defaultCurrency
	}

	if req.IsActive != nil {
		unit.IsActive = *req.IsActive
	}

	if err := s.repo.Create(unit); err != nil {
		return nil, fmt.Errorf("failed to create organization unit: %w", err)
	}

	return unit, nil
}

// GetOrganizationUnitByID retrieves an organization unit by ID
func (s *OrganizationUnitService) GetOrganizationUnitByID(id uint) (*models.OrganizationUnit, error) {
	return s.repo.FindByID(id)
}

// UpdateOrganizationUnit updates an existing organization unit
func (s *OrganizationUnitService) UpdateOrganizationUnit(id uint, req *models.UpdateOrganizationUnitRequest, updatedBy *uint) (*models.OrganizationUnit, error) {
	unit, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("organization unit not found")
	}

	// Check code uniqueness if code is being updated
	if req.Code != nil && *req.Code != unit.Code {
		if s.repo.ExistsByCode(*req.Code) {
			return nil, errors.New("organization unit with this code already exists")
		}
		unit.Code = *req.Code
	}

	// Validate organization if being updated
	if req.OrganizationID != nil {
		_, err := s.organizationRepo.FindByID(*req.OrganizationID)
		if err != nil {
			return nil, errors.New("organization not found")
		}
		unit.OrganizationID = req.OrganizationID
	}

	// Validate parent unit if being updated
	if req.ParentUnitID != nil {
		if *req.ParentUnitID != 0 {
			_, err := s.repo.FindByID(*req.ParentUnitID)
			if err != nil {
				return nil, errors.New("parent organization unit not found")
			}
		}
		unit.ParentUnitID = req.ParentUnitID
	}

	if req.Name != nil {
		unit.Name = *req.Name
	}
	if req.Description != nil {
		unit.Description = req.Description
	}
	if req.UnitType != nil {
		unit.UnitType = req.UnitType
	}
	if req.HeadID != nil {
		unit.HeadID = req.HeadID
	}
	if req.BudgetAllocated != nil {
		unit.BudgetAllocated = req.BudgetAllocated
	}
	if req.BudgetCurrency != nil {
		unit.BudgetCurrency = req.BudgetCurrency
	}
	if req.LocationID != nil {
		unit.LocationID = req.LocationID
	}
	if req.IsActive != nil {
		unit.IsActive = *req.IsActive
	}

	unit.UpdatedBy = updatedBy

	if err := s.repo.Update(unit); err != nil {
		return nil, fmt.Errorf("failed to update organization unit: %w", err)
	}

	return unit, nil
}

// DeleteOrganizationUnit soft deletes an organization unit
func (s *OrganizationUnitService) DeleteOrganizationUnit(id uint) error {
	unit, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("organization unit not found")
	}

	// Check if unit has sub-units
	if len(unit.SubUnits) > 0 {
		return errors.New("cannot delete organization unit with sub-units")
	}

	return s.repo.Delete(id)
}

// ListOrganizationUnits lists organization units with pagination and filters
func (s *OrganizationUnitService) ListOrganizationUnits(tenantID, organizationID *uint, page, pageSize int, filters map[string]interface{}) ([]models.OrganizationUnit, int64, error) {
	return s.repo.List(tenantID, organizationID, page, pageSize, filters)
}

// GetRootOrganizationUnits gets all root organization units
func (s *OrganizationUnitService) GetRootOrganizationUnits(organizationID *uint) ([]models.OrganizationUnit, error) {
	return s.repo.FindRootUnits(organizationID)
}
