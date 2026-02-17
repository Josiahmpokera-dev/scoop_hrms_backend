package services

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/repositories"
)

type TalentReviewService struct {
	repo *repositories.TalentReviewRepository
}

func NewTalentReviewService() *TalentReviewService {
	return &TalentReviewService{repo: repositories.NewTalentReviewRepository()}
}

// ratingToBucket: Low=1-2.5->0, Mid=2.5-3.5->1, High=3.5-5->2
func ratingToBucket(r float64) int {
	if r < 2.5 {
		return 0
	}
	if r < 3.5 {
		return 1
	}
	return 2
}

// computeBox returns 9-box (1-9) from performance and potential ratings.
func computeBox(perf, potential float64) int {
	p := ratingToBucket(perf)
	pot := ratingToBucket(potential)
	box := pot*3 + p + 1
	if box < 1 {
		box = 1
	}
	if box > 9 {
		box = 9
	}
	return box
}

func (s *TalentReviewService) List(department string, box *int, criticalRole *bool, riskOfLoss string, page, pageSize int) ([]models.TalentReview, int64, error) {
	return s.repo.List(department, box, criticalRole, riskOfLoss, page, pageSize)
}

func (s *TalentReviewService) GetByEmployeeID(employeeID uint) (*models.TalentReview, error) {
	return s.repo.FindByEmployeeID(employeeID)
}

func (s *TalentReviewService) Update(employeeID uint, req *models.UpdateTalentReviewRequest) (*models.TalentReview, error) {
	review, err := s.repo.FindByEmployeeID(employeeID)
	if err != nil {
		return nil, err
	}
	if req.PerformanceRating != nil {
		review.PerformanceRating = *req.PerformanceRating
	}
	if req.PotentialRating != nil {
		review.PotentialRating = *req.PotentialRating
	}
	if req.Box != nil {
		review.Box = *req.Box
	}
	if req.BoxLabel != nil {
		review.BoxLabel = *req.BoxLabel
	}
	if req.CriticalRole != nil {
		review.CriticalRole = *req.CriticalRole
	}
	if req.Tenure != nil {
		review.Tenure = req.Tenure
	}
	if req.RiskOfLoss != nil {
		review.RiskOfLoss = *req.RiskOfLoss
	}
	if req.Readiness != nil {
		review.Readiness = req.Readiness
	}
	if req.DevelopmentPriority != nil {
		review.DevelopmentPriority = req.DevelopmentPriority
	}
	if req.SuccessorFor != nil {
		review.SuccessorFor = req.SuccessorFor
	}
	if req.Notes != nil {
		review.Notes = req.Notes
	}
	// Compute Box from performance + potential if both set and Box not explicitly overridden
	if req.PerformanceRating != nil || req.PotentialRating != nil {
		if req.Box == nil {
			review.Box = computeBox(review.PerformanceRating, review.PotentialRating)
			review.BoxLabel = boxLabel(review.Box)
		}
	}
	if err := s.repo.Update(review); err != nil {
		return nil, err
	}
	return review, nil
}

func boxLabel(box int) string {
	labels := map[int]string{
		1: "Low Potential / Low Performance", 2: "Low Potential / Mid Performance", 3: "Low Potential / High Performance",
		4: "Mid Potential / Low Performance", 5: "Mid Potential / Mid Performance", 6: "Mid Potential / High Performance",
		7: "High Potential / Low Performance", 8: "High Potential / Mid Performance", 9: "High Potential / High Performance",
	}
	if l, ok := labels[box]; ok {
		return l
	}
	return ""
}

func (s *TalentReviewService) CreateCalibrationSession(createdByID uint, req *models.CreateCalibrationSessionRequest) (*models.CalibrationSession, error) {
	session := &models.CalibrationSession{
		Name:           req.Name,
		Department:     req.Department,
		ParticipantIDs: req.ParticipantIDs,
		CreatedByID:    createdByID,
	}
	if err := s.repo.CreateCalibrationSession(session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *TalentReviewService) AddSuccessionPlan(employeeID uint, createdByID uint, req *models.AddSuccessionPlanRequest) (*models.SuccessionPlan, error) {
	plan := &models.SuccessionPlan{
		EmployeeID:      employeeID,
		Role:            req.Role,
		Readiness:       req.Readiness,
		DevelopmentPlan: req.DevelopmentPlan,
		CreatedByID:     createdByID,
	}
	if err := s.repo.CreateSuccessionPlan(plan); err != nil {
		return nil, err
	}
	return plan, nil
}
