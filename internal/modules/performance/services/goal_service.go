package services

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/repositories"
)

type GoalService struct {
	repo *repositories.GoalRepository
}

func NewGoalService() *GoalService {
	return &GoalService{repo: repositories.NewGoalRepository()}
}

func (s *GoalService) CreateGoal(ownerID uint, ownerName string, req *models.CreateGoalRequest) (*models.Goal, error) {
	goal := &models.Goal{
		GoalCode:       s.repo.GetNextGoalCode(),
		Title:          req.Title,
		Description:    req.Description,
		Type:           req.Type,
		Level:          req.Level,
		OwnerID:        ownerID,
		OwnerName:      ownerName,
		Department:     req.Department,
		ParentGoalID:   req.ParentGoalID,
		Weight:         req.Weight,
		DueDate:        req.DueDate,
		Status:         "Not Started",
		Progress:       0,
		Visibility:     "Manager",
		ApprovalStatus: "Draft",
	}
	if goal.Weight == 0 {
		goal.Weight = 100
	}
	if req.Visibility != "" {
		goal.Visibility = req.Visibility
	}
	if err := s.repo.Create(goal); err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *GoalService) GetGoal(id uint) (*models.Goal, error) {
	return s.repo.FindByID(id)
}

func (s *GoalService) UpdateGoal(id uint, req *models.UpdateGoalRequest) (*models.Goal, error) {
	goal, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		goal.Title = *req.Title
	}
	if req.Description != nil {
		goal.Description = req.Description
	}
	if req.Type != nil {
		goal.Type = *req.Type
	}
	if req.Level != nil {
		goal.Level = *req.Level
	}
	if req.Department != nil {
		goal.Department = req.Department
	}
	if req.ParentGoalID != nil {
		goal.ParentGoalID = req.ParentGoalID
	}
	if req.Weight != nil {
		goal.Weight = *req.Weight
	}
	if req.DueDate != nil {
		goal.DueDate = req.DueDate
	}
	if req.Status != nil {
		goal.Status = *req.Status
	}
	if req.Progress != nil {
		goal.Progress = *req.Progress
	}
	if req.Visibility != nil {
		goal.Visibility = *req.Visibility
	}
	if err := s.repo.Update(goal); err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *GoalService) DeleteGoal(id uint) error {
	return s.repo.Delete(id)
}

func (s *GoalService) ListGoals(ownerID *uint, level, status, department, search string, page, pageSize int) ([]models.Goal, int64, error) {
	return s.repo.List(ownerID, level, status, department, search, page, pageSize)
}

func (s *GoalService) GetGoalStats(ownerID *uint) (map[string]interface{}, error) {
	counts, err := s.repo.CountByStatus(ownerID)
	if err != nil {
		return nil, err
	}
	goals, _, err := s.repo.List(ownerID, "", "", "", "", 1, 10000)
	if err != nil {
		return nil, err
	}
	var sum float64
	for _, g := range goals {
		sum += g.Progress
	}
	avgProgress := 0.0
	if len(goals) > 0 {
		avgProgress = sum / float64(len(goals))
	}
	out := make(map[string]interface{})
	for k, v := range counts {
		out[k] = v
	}
	out["avg_progress"] = avgProgress
	return out, nil
}

func (s *GoalService) CreateCheckIn(goalID uint, userID uint, userName string, req *models.CreateCheckInRequest) (*models.GoalCheckIn, error) {
	goal, err := s.repo.FindByID(goalID)
	if err != nil {
		return nil, err
	}
	checkIn := &models.GoalCheckIn{
		GoalID:        goalID,
		Progress:      req.Progress,
		Status:        req.Status,
		Notes:         req.Notes,
		EvidenceURL:   req.EvidenceURL,
		UpdatedByID:   userID,
		UpdatedByName: userName,
	}
	if err := s.repo.CreateCheckIn(checkIn); err != nil {
		return nil, err
	}
	goal.Progress = req.Progress
	if req.Status == "Behind" {
		goal.Status = "At Risk"
	} else if req.Status == "On Track" || req.Status == "At Risk" {
		goal.Status = req.Status
	}
	_ = s.repo.Update(goal)
	return checkIn, nil
}

func (s *GoalService) CreateKeyResult(goalID uint, req *models.CreateKeyResultRequest) (*models.KeyResult, error) {
	kr := &models.KeyResult{
		KRCode:      s.repo.GetNextKRCode(),
		GoalID:      goalID,
		Title:       req.Title,
		TargetValue: req.TargetValue,
		Unit:        req.Unit,
		DueDate:     req.DueDate,
	}
	if err := s.repo.CreateKeyResult(kr); err != nil {
		return nil, err
	}
	return kr, nil
}

func (s *GoalService) UpdateKeyResult(krID uint, req *models.UpdateKeyResultRequest) (*models.KeyResult, error) {
	kr, err := s.repo.FindKeyResultByID(krID)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		kr.Title = *req.Title
	}
	if req.TargetValue != nil {
		kr.TargetValue = *req.TargetValue
	}
	if req.CurrentValue != nil {
		kr.CurrentValue = *req.CurrentValue
	}
	if req.Unit != nil {
		kr.Unit = *req.Unit
	}
	if req.DueDate != nil {
		kr.DueDate = req.DueDate
	}
	if req.Status != nil {
		kr.Status = *req.Status
	}
	if err := s.repo.UpdateKeyResult(kr); err != nil {
		return nil, err
	}
	return kr, nil
}

func (s *GoalService) DeleteKeyResult(krID uint) error {
	return s.repo.DeleteKeyResult(krID)
}

func (s *GoalService) GetAlignmentMap() (map[string]interface{}, error) {
	goals, err := s.repo.ListAllForAlignment()
	if err != nil {
		return nil, err
	}
	nodes := make([]map[string]interface{}, 0, len(goals))
	edges := make([]map[string]interface{}, 0)
	for _, g := range goals {
		nodes = append(nodes, map[string]interface{}{
			"id": g.GoalCode, "title": g.Title, "level": g.Level, "owner_id": g.OwnerID,
		})
	}
	for _, g := range goals {
		if g.ParentGoalID != nil {
			var parentCode string
			for _, pg := range goals {
				if pg.ID == *g.ParentGoalID {
					parentCode = pg.GoalCode
					break
				}
			}
			if parentCode != "" {
				edges = append(edges, map[string]interface{}{"from": g.GoalCode, "to": parentCode})
			}
		}
	}
	return map[string]interface{}{"nodes": nodes, "edges": edges}, nil
}

func (s *GoalService) SubmitForApproval(goalID uint, comments string) (*models.Goal, error) {
	goal, err := s.repo.FindByID(goalID)
	if err != nil {
		return nil, err
	}
	goal.ApprovalStatus = "Pending Approval"
	if err := s.repo.Update(goal); err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *GoalService) ApproveGoal(goalID uint, approverID uint, approverName string, action string, comments string) (*models.Goal, error) {
	goal, err := s.repo.FindByID(goalID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	if action == "approve" || action == "approved" {
		goal.ApprovalStatus = "Approved"
		goal.ApprovedByID = &approverID
		goal.ApprovedByName = &approverName
		goal.ApprovedDate = &now
	} else {
		goal.ApprovalStatus = "Rejected"
	}
	if err := s.repo.Update(goal); err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *GoalService) RequestCompletion(goalID uint, evidence string, evidenceURL string, finalProgress float64) (*models.Goal, error) {
	goal, err := s.repo.FindByID(goalID)
	if err != nil {
		return nil, err
	}
	status := "Pending Verification"
	goal.CompletionStatus = &status
	if evidence != "" {
		goal.CompletionEvidence = &evidence
	}
	if evidenceURL != "" {
		goal.CompletionEvidenceURL = &evidenceURL
	}
	goal.Progress = finalProgress
	if err := s.repo.Update(goal); err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *GoalService) VerifyCompletion(goalID uint, verifierID uint, verifierName string, action string, finalRating *float64, comments string) (*models.Goal, error) {
	goal, err := s.repo.FindByID(goalID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	if action == "verify" || action == "verified" {
		status := "Verified"
		goal.CompletionStatus = &status
		goal.Status = "Completed"
		goal.Progress = 100
		goal.VerifiedByID = &verifierID
		goal.VerifiedByName = &verifierName
		goal.VerifiedDate = &now
		if finalRating != nil {
			goal.FinalRating = finalRating
		}
	} else {
		rejected := "Rejected"
		goal.CompletionStatus = &rejected
	}
	if err := s.repo.Update(goal); err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *GoalService) AssignGoal(assignerID uint, assignerName string, req *models.AssignGoalRequest) (*models.Goal, error) {
	source, err := s.repo.FindByID(req.GoalID)
	if err != nil {
		return nil, err
	}
	goal := &models.Goal{
		GoalCode:       s.repo.GetNextGoalCode(),
		Title:          source.Title,
		Description:    source.Description,
		Type:           source.Type,
		Level:          source.Level,
		OwnerID:        req.AssignedTo,
		OwnerName:      "", // caller may set separately
		Department:     source.Department,
		ParentGoalID:   source.ParentGoalID,
		Weight:         source.Weight,
		DueDate:        source.DueDate,
		Status:         "In Progress",
		Progress:        0,
		Visibility:     source.Visibility,
		ApprovalStatus: "Approved",
		ApprovedByID:   &assignerID,
		ApprovedByName: &assignerName,
		AssignedToID:   &req.AssignedTo,
	}
	now := time.Now()
	goal.ApprovedDate = &now
	if err := s.repo.Create(goal); err != nil {
		return nil, err
	}
	for _, krInput := range req.KeyResults {
		kr := &models.KeyResult{
			KRCode:      s.repo.GetNextKRCode(),
			GoalID:      goal.ID,
			Title:       krInput.Title,
			TargetValue: krInput.TargetValue,
			Unit:        krInput.Unit,
			DueDate:     krInput.DueDate,
		}
		_ = s.repo.CreateKeyResult(kr)
	}
	return s.repo.FindByID(goal.ID)
}

func (s *GoalService) LinkProject(goalID uint, req *models.LinkProjectRequest) (*models.Goal, error) {
	goal, err := s.repo.FindByID(goalID)
	if err != nil {
		return nil, err
	}
	goal.LinkedProjectID = &req.ProjectID
	goal.LinkedProjectCode = &req.ProjectCode
	goal.LinkedProjectName = &req.ProjectName
	if err := s.repo.Update(goal); err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *GoalService) UnlinkProject(goalID uint) (*models.Goal, error) {
	goal, err := s.repo.FindByID(goalID)
	if err != nil {
		return nil, err
	}
	goal.LinkedProjectID = nil
	goal.LinkedProjectCode = nil
	goal.LinkedProjectName = nil
	if err := s.repo.Update(goal); err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *GoalService) GetDashboardStats(userID uint) (map[string]interface{}, error) {
	ownerID := &userID
	goalCounts, err := s.repo.CountByStatus(ownerID)
	if err != nil {
		return nil, err
	}
	goals, _, err := s.repo.List(ownerID, "", "", "", "", 1, 1000)
	if err != nil {
		return nil, err
	}
	var sum float64
	for _, g := range goals {
		sum += g.Progress
	}
	avgProgress := 0.0
	if len(goals) > 0 {
		avgProgress = sum / float64(len(goals))
	}
	appraisalRepo := repositories.NewAppraisalRepository()
	appraisalCounts, _ := appraisalRepo.CountAppraisalsByStatus(&userID)
	activeCycle, _ := appraisalRepo.FindActiveCycle()
	appraisalStatus := ""
	if activeCycle != nil {
		appraisalStatus = activeCycle.Status
	}
	lastRating := interface{}(nil)
	appraisals, _, _ := appraisalRepo.ListAppraisals(&userID, nil, "Completed", "", true, 1, 1)
	if len(appraisals) > 0 && appraisals[0].OverallRating != nil {
		lastRating = *appraisals[0].OverallRating
	}
	return map[string]interface{}{
		"goals": map[string]interface{}{
			"total":       goalCounts["Completed"] + goalCounts["In Progress"] + goalCounts["Not Started"] + goalCounts["At Risk"] + goalCounts["Cancelled"],
			"completed":   goalCounts["Completed"],
			"in_progress": goalCounts["In Progress"],
			"at_risk":     goalCounts["At Risk"],
			"avg_progress": avgProgress,
		},
		"appraisal": map[string]interface{}{
			"status":       appraisalStatus,
			"count_by_status": appraisalCounts,
			"last_rating": lastRating,
		},
		"quick_stats": map[string]interface{}{
			"goals_on_track": goalCounts["In Progress"],
			"last_overall_rating": lastRating,
		},
	}, nil
}

func (s *GoalService) GetUpcomingActions(userID uint) ([]map[string]interface{}, error) {
	actions := make([]map[string]interface{}, 0)
	ownerID := &userID
	goals, _, err := s.repo.List(ownerID, "", "In Progress", "", "", 1, 50)
	if err != nil {
		return actions, err
	}
	now := time.Now()
	for _, g := range goals {
		if g.DueDate != nil && g.DueDate.Before(now.AddDate(0, 0, 14)) && g.DueDate.After(now) {
			actions = append(actions, map[string]interface{}{
				"type": "checkin_due", "goal_id": g.ID, "goal_code": g.GoalCode, "title": g.Title, "due_date": g.DueDate,
			})
		}
	}
	pendingApproval, _, _ := s.repo.List(nil, "", "", "", "", 1, 100)
	for _, g := range pendingApproval {
		if g.ApprovalStatus == "Pending Approval" {
			actions = append(actions, map[string]interface{}{
				"type": "pending_approval", "goal_id": g.ID, "goal_code": g.GoalCode, "title": g.Title,
			})
		}
	}
	appraisalRepo := repositories.NewAppraisalRepository()
	appraisals, _, _ := appraisalRepo.ListAppraisals(&userID, nil, "Not Started", "", true, 1, 10)
	for _, a := range appraisals {
		actions = append(actions, map[string]interface{}{
			"type": "self_review", "appraisal_id": a.ID, "code": a.Code, "status": a.Status,
		})
	}
	return actions, nil
}
