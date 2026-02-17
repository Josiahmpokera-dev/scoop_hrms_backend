package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/repositories"
)

type DepartmentTargetService struct {
	repo *repositories.DepartmentTargetRepository
}

func NewDepartmentTargetService() *DepartmentTargetService {
	return &DepartmentTargetService{repo: repositories.NewDepartmentTargetRepository()}
}

func (s *DepartmentTargetService) GetByID(id uint) (*models.DepartmentTarget, error) {
	return s.repo.FindByID(id)
}

func (s *DepartmentTargetService) Create(createdByID uint, createdByName string, req *models.CreateDepartmentTargetRequest) (*models.DepartmentTarget, error) {
	_ = createdByName // model has no CreatedByName field
	target := &models.DepartmentTarget{
		Code:         s.repo.GetNextCode(),
		Title:        req.Title,
		Description:  req.Description,
		Department:   req.Department,
		Category:     req.Category,
		Metric:       req.Metric,
		TargetValue:  req.TargetValue,
		Unit:         req.Unit,
		Period:       req.Period,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Priority:     req.Priority,
		CreatedByID:  createdByID,
		Status:       "Not Started",
	}
	if err := s.repo.Create(target); err != nil {
		return nil, err
	}
	for _, m := range req.Milestones {
		milestone := &models.DepartmentTargetMilestone{
			DepartmentTargetID: target.ID,
			Title:              m.Title,
			DueDate:            m.DueDate,
			Status:             "Not Started",
		}
		if err := s.repo.CreateMilestone(milestone); err != nil {
			return nil, err
		}
	}
	return s.repo.FindByID(target.ID)
}

func (s *DepartmentTargetService) List(department, category, status, period string, page, pageSize int) ([]models.DepartmentTarget, int64, error) {
	return s.repo.List(department, category, status, period, "", page, pageSize)
}

func (s *DepartmentTargetService) Update(id uint, req *models.UpdateDepartmentTargetRequest) (*models.DepartmentTarget, error) {
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
	if req.Category != nil {
		target.Category = *req.Category
	}
	if req.Metric != nil {
		target.Metric = *req.Metric
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
	if req.StartDate != nil {
		target.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		target.EndDate = req.EndDate
	}
	if req.Status != nil {
		target.Status = *req.Status
	}
	if req.Priority != nil {
		target.Priority = *req.Priority
	}
	if err := s.repo.Update(target); err != nil {
		return nil, err
	}
	return target, nil
}

func (s *DepartmentTargetService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *DepartmentTargetService) UpdateProgress(id uint, currentValue float64, status, notes string) (*models.DepartmentTarget, error) {
	target, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	target.CurrentValue = currentValue
	if target.TargetValue > 0 {
		progress := (currentValue / target.TargetValue) * 100
		if progress > 100 {
			progress = 100
		}
		if status == "" && progress >= 100 {
			target.Status = "Completed"
		} else if status == "" {
			target.Status = "In Progress"
		}
	}
	if status != "" {
		target.Status = status
	}
	_ = notes // model has no notes field
	if err := s.repo.Update(target); err != nil {
		return nil, err
	}
	return target, nil
}

func (s *DepartmentTargetService) CompleteMilestone(targetID, milestoneID uint) (*models.DepartmentTargetMilestone, error) {
	m, err := s.repo.FindMilestoneByID(milestoneID)
	if err != nil {
		return nil, err
	}
	if m.DepartmentTargetID != targetID {
		return nil, errors.New("milestone does not belong to this target")
	}
	now := time.Now()
	m.Status = "Completed"
	m.CompletedDate = &now
	if err := s.repo.UpdateMilestone(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *DepartmentTargetService) LinkGoal(targetID uint, goalID string) (*models.DepartmentTarget, error) {
	target, err := s.repo.FindByID(targetID)
	if err != nil {
		return nil, err
	}
	var ids []string
	if target.LinkedGoalIDs != nil && *target.LinkedGoalIDs != "" {
		_ = json.Unmarshal([]byte(*target.LinkedGoalIDs), &ids)
	}
	ids = append(ids, goalID)
	b, _ := json.Marshal(ids)
	str := string(b)
	target.LinkedGoalIDs = &str
	if err := s.repo.Update(target); err != nil {
		return nil, err
	}
	return target, nil
}

func (s *DepartmentTargetService) LinkProject(targetID uint, req *models.LinkProjectRequest) (*models.DepartmentTarget, error) {
	target, err := s.repo.FindByID(targetID)
	if err != nil {
		return nil, err
	}
	target.LinkedProjectID = &req.ProjectID
	target.LinkedProjectCode = &req.ProjectCode
	target.LinkedProjectName = &req.ProjectName
	if err := s.repo.Update(target); err != nil {
		return nil, err
	}
	return target, nil
}

func (s *DepartmentTargetService) UnlinkProject(targetID uint) (*models.DepartmentTarget, error) {
	target, err := s.repo.FindByID(targetID)
	if err != nil {
		return nil, err
	}
	target.LinkedProjectID = nil
	target.LinkedProjectCode = nil
	target.LinkedProjectName = nil
	if err := s.repo.Update(target); err != nil {
		return nil, err
	}
	return target, nil
}
