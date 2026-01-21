package services

import (
	"errors"
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/repositories"
	departmentRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
)

type JobPositionService struct {
	repo            *repositories.JobPositionRepository
	departmentRepo  *departmentRepos.DepartmentRepository
}

func NewJobPositionService() *JobPositionService {
	return &JobPositionService{
		repo:           repositories.NewJobPositionRepository(),
		departmentRepo: departmentRepos.NewDepartmentRepository(),
	}
}

// CreateJobPosition creates a new job position
func (s *JobPositionService) CreateJobPosition(req *models.CreateJobPositionRequest, tenantID *uint, updatedBy *uint) (*models.JobPosition, error) {
	// Check if code already exists
	if s.repo.ExistsByCode(req.Code) {
		return nil, errors.New("job position with this code already exists")
	}

	// Validate department if provided
	if req.DepartmentID != nil {
		department, err := s.departmentRepo.FindByID(*req.DepartmentID)
		if err != nil {
			return nil, fmt.Errorf("department with ID %d not found", *req.DepartmentID)
		}
		// Check if department belongs to the same tenant (if tenant_id is set)
		if tenantID != nil && department.TenantID != nil {
			if *tenantID != *department.TenantID {
				return nil, fmt.Errorf("department with ID %d does not belong to your tenant", *req.DepartmentID)
			}
		}
		// Check if department is active
		if !department.IsActive {
			return nil, fmt.Errorf("department with ID %d is not active", *req.DepartmentID)
		}
	}

	position := &models.JobPosition{
		TenantID:            tenantID,
		Code:                req.Code,
		Title:               req.Title,
		Grade:               req.Grade,
		DepartmentID:        req.DepartmentID,
		ReportsToPositionID: req.ReportsToPositionID,
		BudgetedHeadcount:   req.BudgetedHeadcount,
		CurrentHeadcount:    req.CurrentHeadcount,
		EmploymentType:      req.EmploymentType,
		KeyCompetencies:     req.KeyCompetencies,
		IsActive:            true,
		UpdatedBy:           updatedBy,
	}

	if req.IsActive != nil {
		position.IsActive = *req.IsActive
	}

	if err := s.repo.Create(position); err != nil {
		return nil, fmt.Errorf("failed to create job position: %w", err)
	}

	return position, nil
}

// GetJobPositionByID retrieves a job position by ID
func (s *JobPositionService) GetJobPositionByID(id uint) (*models.JobPosition, error) {
	return s.repo.FindByID(id)
}

// UpdateJobPosition updates an existing job position
func (s *JobPositionService) UpdateJobPosition(id uint, req *models.UpdateJobPositionRequest, updatedBy *uint) (*models.JobPosition, error) {
	position, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("job position not found")
	}

	// Check code uniqueness if code is being updated
	if req.Code != nil && *req.Code != position.Code {
		if s.repo.ExistsByCode(*req.Code) {
			return nil, errors.New("job position with this code already exists")
		}
		position.Code = *req.Code
	}

	if req.Title != nil {
		position.Title = *req.Title
	}
	if req.Grade != nil {
		position.Grade = req.Grade
	}
	if req.DepartmentID != nil {
		// Validate department if being updated
		department, err := s.departmentRepo.FindByID(*req.DepartmentID)
		if err != nil {
			return nil, fmt.Errorf("department with ID %d not found", *req.DepartmentID)
		}
		// Check if department belongs to the same tenant (if tenant_id is set)
		if position.TenantID != nil && department.TenantID != nil {
			if *position.TenantID != *department.TenantID {
				return nil, fmt.Errorf("department with ID %d does not belong to your tenant", *req.DepartmentID)
			}
		}
		// Check if department is active
		if !department.IsActive {
			return nil, fmt.Errorf("department with ID %d is not active", *req.DepartmentID)
		}
		position.DepartmentID = req.DepartmentID
	}
	if req.ReportsToPositionID != nil {
		position.ReportsToPositionID = req.ReportsToPositionID
	}
	if req.BudgetedHeadcount != nil {
		position.BudgetedHeadcount = req.BudgetedHeadcount
	}
	if req.CurrentHeadcount != nil {
		position.CurrentHeadcount = req.CurrentHeadcount
	}
	if req.EmploymentType != nil {
		position.EmploymentType = req.EmploymentType
	}
	if req.KeyCompetencies != nil {
		position.KeyCompetencies = req.KeyCompetencies
	}
	if req.IsActive != nil {
		position.IsActive = *req.IsActive
	}

	position.UpdatedBy = updatedBy

	if err := s.repo.Update(position); err != nil {
		return nil, fmt.Errorf("failed to update job position: %w", err)
	}

	return position, nil
}

// DeleteJobPosition soft deletes a job position
func (s *JobPositionService) DeleteJobPosition(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("job position not found")
	}

	return s.repo.Delete(id)
}

// ListJobPositions lists job positions with pagination and filters
func (s *JobPositionService) ListJobPositions(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.JobPosition, int64, error) {
	return s.repo.List(tenantID, page, pageSize, filters)
}

// GetPositionsByDepartment retrieves all active positions for a specific department
func (s *JobPositionService) GetPositionsByDepartment(departmentID uint, tenantID *uint) ([]models.JobPosition, error) {
	// Validate department exists and belongs to tenant
	department, err := s.departmentRepo.FindByID(departmentID)
	if err != nil {
		return nil, fmt.Errorf("department with ID %d not found", departmentID)
	}

	// Check if department belongs to the same tenant (if tenant_id is set)
	if tenantID != nil && department.TenantID != nil {
		if *tenantID != *department.TenantID {
			return nil, fmt.Errorf("department with ID %d does not belong to your tenant", departmentID)
		}
	}

	// Get positions for this department
	positions, err := s.repo.FindByDepartmentID(departmentID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve positions for department: %w", err)
	}

	return positions, nil
}
