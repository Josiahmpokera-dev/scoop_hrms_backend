package services

import (
	"errors"
	"fmt"

	departmentRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/teams/models"
	teamRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/teams/repositories"
)

type TeamService struct {
	repo           *teamRepos.TeamRepository
	departmentRepo *departmentRepos.DepartmentRepository
}

func NewTeamService() *TeamService {
	return &TeamService{
		repo:           teamRepos.NewTeamRepository(),
		departmentRepo: departmentRepos.NewDepartmentRepository(),
	}
}

// CreateTeam creates a new team
func (s *TeamService) CreateTeam(req *models.CreateTeamRequest, tenantID *uint, updatedBy *uint) (*models.Team, error) {
	// Check if code already exists
	if s.repo.ExistsByCode(req.Code) {
		return nil, errors.New("team with this code already exists")
	}

	// Validate department if provided
	if req.DepartmentID != nil {
		_, err := s.departmentRepo.FindByID(*req.DepartmentID)
		if err != nil {
			return nil, errors.New("department not found")
		}
	}

	team := &models.Team{
		TenantID:     tenantID,
		DepartmentID: req.DepartmentID,
		Code:         req.Code,
		Name:         req.Name,
		Description:  req.Description,
		TeamLeadID:   req.TeamLeadID,
		TeamType:     req.TeamType,
		ProjectID:    req.ProjectID,
		MaxMembers:   req.MaxMembers,
		LocationID:   req.LocationID,
		IsActive:     true,
		UpdatedBy:    updatedBy,
	}

	// Set default max members if not provided
	if team.MaxMembers == nil {
		defaultMaxMembers := 10
		team.MaxMembers = &defaultMaxMembers
	}

	if req.IsActive != nil {
		team.IsActive = *req.IsActive
	}

	if err := s.repo.Create(team); err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	return team, nil
}

// GetTeamByID retrieves a team by ID
func (s *TeamService) GetTeamByID(id uint) (*models.Team, error) {
	return s.repo.FindByID(id)
}

// UpdateTeam updates an existing team
func (s *TeamService) UpdateTeam(id uint, req *models.UpdateTeamRequest, updatedBy *uint) (*models.Team, error) {
	team, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("team not found")
	}

	// Check code uniqueness if code is being updated
	if req.Code != nil && *req.Code != team.Code {
		if s.repo.ExistsByCode(*req.Code) {
			return nil, errors.New("team with this code already exists")
		}
		team.Code = *req.Code
	}

	// Validate department if being updated
	if req.DepartmentID != nil {
		_, err := s.departmentRepo.FindByID(*req.DepartmentID)
		if err != nil {
			return nil, errors.New("department not found")
		}
		team.DepartmentID = req.DepartmentID
	}

	if req.Name != nil {
		team.Name = *req.Name
	}
	if req.Description != nil {
		team.Description = req.Description
	}
	if req.TeamLeadID != nil {
		team.TeamLeadID = req.TeamLeadID
	}
	if req.TeamType != nil {
		team.TeamType = req.TeamType
	}
	if req.ProjectID != nil {
		team.ProjectID = req.ProjectID
	}
	if req.MaxMembers != nil {
		team.MaxMembers = req.MaxMembers
	}
	if req.LocationID != nil {
		team.LocationID = req.LocationID
	}
	if req.IsActive != nil {
		team.IsActive = *req.IsActive
	}

	team.UpdatedBy = updatedBy

	if err := s.repo.Update(team); err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}

	return team, nil
}

// DeleteTeam soft deletes a team
func (s *TeamService) DeleteTeam(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("team not found")
	}

	return s.repo.Delete(id)
}

// ListTeams lists teams with pagination and filters
func (s *TeamService) ListTeams(tenantID, departmentID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Team, int64, error) {
	return s.repo.List(tenantID, departmentID, page, pageSize, filters)
}

// GetTeamsByDepartment gets all teams in a department
func (s *TeamService) GetTeamsByDepartment(departmentID uint) ([]models.Team, error) {
	return s.repo.FindByDepartment(departmentID)
}
