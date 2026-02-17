package services

import (
	"errors"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/repositories"
)

type AppraisalService struct {
	repo *repositories.AppraisalRepository
}

func NewAppraisalService() *AppraisalService {
	return &AppraisalService{repo: repositories.NewAppraisalRepository()}
}

func (s *AppraisalService) CreateCycle(createdByID uint, req *models.CreateAppraisalCycleRequest) (*models.AppraisalCycle, error) {
	_ = createdByID
	cycle := &models.AppraisalCycle{
		Code:      s.repo.GetNextCycleCode(),
		Name:      req.Name,
		CycleType: req.CycleType,
		Year:      req.Year,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    "Draft",
	}
	if err := s.repo.CreateCycle(cycle); err != nil {
		return nil, err
	}
	for _, st := range req.Steps {
		step := &models.AppraisalWorkflowStep{
			AppraisalCycleID: cycle.ID,
			StepOrder:        st.StepOrder,
			Name:             st.Name,
			Description:      st.Description,
			DueDate:          st.DueDate,
		}
		if err := s.repo.CreateWorkflowStep(step); err != nil {
			return nil, err
		}
	}
	return s.repo.FindCycleByID(cycle.ID)
}

func (s *AppraisalService) GetCycleByID(id uint) (*models.AppraisalCycle, error) {
	return s.repo.FindCycleByID(id)
}

func (s *AppraisalService) UpdateCycle(id uint, req *models.UpdateAppraisalCycleRequest) (*models.AppraisalCycle, error) {
	cycle, err := s.repo.FindCycleByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		cycle.Name = *req.Name
	}
	if req.CycleType != nil {
		cycle.CycleType = *req.CycleType
	}
	if req.Year != nil {
		cycle.Year = *req.Year
	}
	if req.StartDate != nil {
		cycle.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		cycle.EndDate = *req.EndDate
	}
	if err := s.repo.UpdateCycle(cycle); err != nil {
		return nil, err
	}
	return cycle, nil
}

func (s *AppraisalService) ChangeCycleStatus(id uint, status string) (*models.AppraisalCycle, error) {
	cycle, err := s.repo.FindCycleByID(id)
	if err != nil {
		return nil, err
	}
	// Validate transitions: Draft -> Active -> Closed
	switch cycle.Status {
	case "Draft":
		if status != "Active" && status != "Draft" {
			return nil, errors.New("from Draft only Active or Draft is allowed")
		}
	case "Active":
		if status != "Closed" && status != "Active" {
			return nil, errors.New("from Active only Closed or Active is allowed")
		}
	case "Closed":
		return nil, errors.New("cannot change status of a Closed cycle")
	}
	cycle.Status = status
	if err := s.repo.UpdateCycle(cycle); err != nil {
		return nil, err
	}
	return cycle, nil
}

func (s *AppraisalService) DeleteCycle(id uint) error {
	cycle, err := s.repo.FindCycleByID(id)
	if err != nil {
		return err
	}
	if cycle.Status != "Draft" {
		return errors.New("only Draft cycles can be deleted")
	}
	return s.repo.DeleteCycle(id)
}

func (s *AppraisalService) ListCycles(status, cycleType string, year int, page, pageSize int) ([]models.AppraisalCycle, int64, error) {
	return s.repo.ListCycles(status, cycleType, year, page, pageSize)
}

func (s *AppraisalService) GetActiveCycle() (*models.AppraisalCycle, error) {
	return s.repo.FindActiveCycle()
}

func (s *AppraisalService) ListAppraisals(employeeID *uint, cycleID *uint, status, department string, mine bool, page, pageSize int) ([]models.Appraisal, int64, error) {
	return s.repo.ListAppraisals(employeeID, cycleID, status, department, mine, page, pageSize)
}

func (s *AppraisalService) GetAppraisalByID(id uint) (*models.Appraisal, error) {
	return s.repo.FindAppraisalByID(id)
}

func (s *AppraisalService) GetAppraisalSummary(employeeID uint) (map[string]interface{}, error) {
	eID := &employeeID
	counts, err := s.repo.CountAppraisalsByStatus(eID)
	if err != nil {
		return nil, err
	}
	appraisals, _, _ := s.repo.ListAppraisals(eID, nil, "Completed", "", true, 1, 1)
	lastRating := interface{}(nil)
	if len(appraisals) > 0 && appraisals[0].OverallRating != nil {
		lastRating = *appraisals[0].OverallRating
	}
	return map[string]interface{}{
		"count_by_status": counts,
		"last_rating":     lastRating,
	}, nil
}

func (s *AppraisalService) SubmitSelfReview(appraisalID uint, req *models.SubmitSelfReviewRequest, submitterName string) (*models.Appraisal, error) {
	_ = submitterName
	appraisal, err := s.repo.FindAppraisalByID(appraisalID)
	if err != nil {
		return nil, err
	}
	appraisal.SelfRating = &req.SelfRating
	appraisal.Status = "Manager Review"
	if err := s.repo.UpdateAppraisal(appraisal); err != nil {
		return nil, err
	}
	return appraisal, nil
}

func (s *AppraisalService) SubmitManagerReview(appraisalID uint, req *models.SubmitManagerReviewRequest, submitterName string) (*models.Appraisal, error) {
	_ = submitterName
	appraisal, err := s.repo.FindAppraisalByID(appraisalID)
	if err != nil {
		return nil, err
	}
	appraisal.ManagerRating = &req.ManagerRating
	appraisal.Status = "Calibration"
	if err := s.repo.UpdateAppraisal(appraisal); err != nil {
		return nil, err
	}
	return appraisal, nil
}

func (s *AppraisalService) SubmitCalibration(appraisalID uint, req *models.SubmitCalibrationRequest) (*models.Appraisal, error) {
	appraisal, err := s.repo.FindAppraisalByID(appraisalID)
	if err != nil {
		return nil, err
	}
	appraisal.OverallRating = &req.OverallRating
	appraisal.Status = "Completed"
	if err := s.repo.UpdateAppraisal(appraisal); err != nil {
		return nil, err
	}
	return appraisal, nil
}

func (s *AppraisalService) SendBack(appraisalID uint, req *models.SendBackRequest) (*models.Appraisal, error) {
	appraisal, err := s.repo.FindAppraisalByID(appraisalID)
	if err != nil {
		return nil, err
	}
	appraisal.Status = req.ToStep
	if err := s.repo.UpdateAppraisal(appraisal); err != nil {
		return nil, err
	}
	return appraisal, nil
}

func (s *AppraisalService) Finalize(appraisalID uint, req *models.FinalizeRequest) (*models.Appraisal, error) {
	appraisal, err := s.repo.FindAppraisalByID(appraisalID)
	if err != nil {
		return nil, err
	}
	appraisal.OverallRating = &req.OverallRating
	appraisal.Status = "Completed"
	if err := s.repo.UpdateAppraisal(appraisal); err != nil {
		return nil, err
	}
	_ = req.Outcomes // model has no Outcomes field
	return appraisal, nil
}
