package services

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/repositories"
)

type Feedback360Service struct {
	repo *repositories.Feedback360Repository
}

func NewFeedback360Service() *Feedback360Service {
	return &Feedback360Service{repo: repositories.NewFeedback360Repository()}
}

func (s *Feedback360Service) Launch(createdByID uint, req *models.LaunchFeedback360Request) (*models.Feedback360Campaign, error) {
	_ = createdByID
	campaign := &models.Feedback360Campaign{
		Code:        s.repo.GetNextCode(),
		Name:        req.Name,
		Description: req.Description,
		Department:  req.Department,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      "Active",
	}
	if err := s.repo.Create(campaign); err != nil {
		return nil, err
	}
	for _, rg := range req.RaterGroups {
		group := &models.Feedback360RaterGroup{
			CampaignID: campaign.ID,
			Name:       rg.Name,
			MinRaters:  rg.MinRaters,
		}
		if group.MinRaters == 0 {
			group.MinRaters = 1
		}
		if err := s.repo.CreateRaterGroup(group); err != nil {
			return nil, err
		}
	}
	return s.repo.FindByID(campaign.ID)
}

func (s *Feedback360Service) GetByID(id uint) (*models.Feedback360Campaign, error) {
	return s.repo.FindByID(id)
}

func (s *Feedback360Service) List(status, department string, employeeID *uint, page, pageSize int) ([]models.Feedback360Campaign, int64, error) {
	return s.repo.List(status, department, employeeID, page, pageSize)
}
