package services

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/repositories"
)

type EmployeeTargetService struct {
	repo *repositories.EmployeeTargetRepository
}

func NewEmployeeTargetService() *EmployeeTargetService {
	return &EmployeeTargetService{repo: repositories.NewEmployeeTargetRepository()}
}

func (s *EmployeeTargetService) Create(assignerID uint, assignerName string, req *models.CreateEmployeeTargetRequest) (*models.EmployeeTarget, error) {
	_ = assignerName
	target := &models.EmployeeTarget{
		Code:               s.repo.GetNextCode(),
		EmployeeID:         req.EmployeeID,
		DepartmentTargetID: req.DepartmentTargetID,
		Title:              req.Title,
		Description:        req.Description,
		TargetValue:        req.TargetValue,
		Unit:               req.Unit,
		Period:             req.Period,
		Department:         req.Department,
		DueDate:            req.DueDate,
		AssignedByID:       assignerID,
		Status:             "Not Started",
	}
	if err := s.repo.Create(target); err != nil {
		return nil, err
	}
	return target, nil
}

func (s *EmployeeTargetService) BulkAssign(assignerID uint, assignerName string, req *models.BulkAssignTargetRequest) ([]*models.EmployeeTarget, error) {
	_ = assignerName
	targets := make([]*models.EmployeeTarget, 0, len(req.EmployeeIDs))
	for _, empID := range req.EmployeeIDs {
		t := &models.EmployeeTarget{
			Code:         s.repo.GetNextCode(),
			EmployeeID:   empID,
			Title:        req.Title,
			Description:  req.Description,
			TargetValue:  req.TargetValue,
			Unit:         req.Unit,
			Period:       req.Period,
			Department:   req.Department,
			DueDate:      req.DueDate,
			AssignedByID: assignerID,
			Status:       "Not Started",
		}
		if err := s.repo.Create(t); err != nil {
			return nil, err
		}
		targets = append(targets, t)
	}
	return targets, nil
}

func (s *EmployeeTargetService) GetByID(id uint) (*models.EmployeeTarget, error) {
	return s.repo.FindByID(id)
}

func (s *EmployeeTargetService) Update(id uint, req *models.UpdateEmployeeTargetRequest) (*models.EmployeeTarget, error) {
	target, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		target.Title = *req.Title
	}
	if req.Description != nil {
		target.Description = req.Description
	}
	if req.TargetValue != nil {
		target.TargetValue = *req.TargetValue
	}
	if req.Unit != nil {
		target.Unit = *req.Unit
	}
	if req.Period != nil {
		target.Period = *req.Period
	}
	if req.Department != nil {
		target.Department = req.Department
	}
	if req.Status != nil {
		target.Status = *req.Status
	}
	if req.DueDate != nil {
		target.DueDate = req.DueDate
	}
	if err := s.repo.Update(target); err != nil {
		return nil, err
	}
	return target, nil
}

func (s *EmployeeTargetService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *EmployeeTargetService) List(employeeID, department, deptTargetID, status, period, assignedBy string, page, pageSize int) ([]models.EmployeeTarget, int64, error) {
	return s.repo.List(employeeID, department, deptTargetID, status, period, assignedBy, page, pageSize)
}

func (s *EmployeeTargetService) UpdateProgress(id uint, currentValue float64, status, notes string) (*models.EmployeeTarget, error) {
	target, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	target.CurrentValue = currentValue
	if target.TargetValue > 0 {
		progress := (currentValue / target.TargetValue) * 100
		if progress >= 100 && (status == "" || status == "Completed") {
			target.Status = "Completed"
		} else if status != "" {
			target.Status = status
		} else {
			target.Status = "In Progress"
		}
	} else if status != "" {
		target.Status = status
	}
	_ = notes
	if err := s.repo.Update(target); err != nil {
		return nil, err
	}
	return target, nil
}

