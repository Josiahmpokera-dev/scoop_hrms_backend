package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	empRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/realtime"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InterviewWorkflowService implements the stage-driven interview process (exclusive assignee, immutable completed evaluations).
type InterviewWorkflowService struct {
	wf  *repositories.InterviewWorkflowRepository
	rec *repositories.RecruitmentRepository
	emp *empRepos.EmployeeRepository
	hub *realtime.Hub
}

func NewInterviewWorkflowService() *InterviewWorkflowService {
	return &InterviewWorkflowService{
		wf:  repositories.NewInterviewWorkflowRepository(),
		rec: repositories.NewRecruitmentRepository(),
		emp: empRepos.NewEmployeeRepository(),
		hub: realtime.GetHub(),
	}
}

func (s *InterviewWorkflowService) appendAudit(processID, attemptID, eventType, actor, notes string, payload interface{}) {
	var payloadJSON string
	if payload != nil {
		b, _ := json.Marshal(payload)
		payloadJSON = string(b)
	}
	_ = s.wf.AppendAudit(&models.InterviewWorkflowAuditEvent{
		ProcessID:       processID,
		AttemptID:       attemptID,
		EventType:       eventType,
		ActorEmployeeID: actor,
		Notes:           notes,
		PayloadJSON:     payloadJSON,
		CreatedAt:       time.Now(),
	})
}

func (s *InterviewWorkflowService) marshalDimensions(dim []models.ScoringDimension) (string, error) {
	b, err := json.Marshal(dim)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *InterviewWorkflowService) parseDimensionsJSON(js string) ([]models.ScoringDimension, error) {
	if strings.TrimSpace(js) == "" {
		return nil, errors.New("scoring dimensions missing on snapshot")
	}
	var dims []models.ScoringDimension
	if err := json.Unmarshal([]byte(js), &dims); err != nil {
		return nil, err
	}
	return dims, nil
}

func (s *InterviewWorkflowService) validateScores(dims []models.ScoringDimension, scores map[string]float64) error {
	if len(scores) == 0 {
		return errors.New("scores required")
	}
	for _, d := range dims {
		v, ok := scores[d.Key]
		if d.Required && !ok {
			return fmt.Errorf("missing score for %s", d.Key)
		}
		if !ok {
			continue
		}
		if v < d.Min || v > d.Max {
			return fmt.Errorf("score %s out of range [%.1f, %.1f]", d.Key, d.Min, d.Max)
		}
	}
	for k := range scores {
		found := false
		for _, d := range dims {
			if d.Key == k {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unknown score key: %s", k)
		}
	}
	return nil
}

func (s *InterviewWorkflowService) requireAssignee(actorEmployeeID string, attempt *models.InterviewWorkflowStageAttempt) error {
	if strings.TrimSpace(actorEmployeeID) == "" {
		return errors.New("employee context required")
	}
	if strings.TrimSpace(attempt.AssignedEmployeeID) != strings.TrimSpace(actorEmployeeID) {
		return errors.New("only the assigned interviewer may perform this action")
	}
	return nil
}

func (s *InterviewWorkflowService) evaluationLocked(a *models.InterviewWorkflowStageAttempt) bool {
	return a.EvaluatedAt != nil && !a.EvaluatedAt.IsZero()
}

// --- Definitions (HR) ---

func (s *InterviewWorkflowService) CreateDefinition(req models.CreateInterviewWorkflowDefinitionRequest) (*models.InterviewWorkflowDefinition, error) {
	if _, err := s.rec.GetJobOpeningByID(req.JobOpeningID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("job opening not found")
		}
		return nil, err
	}
	def := &models.InterviewWorkflowDefinition{
		ID:                 uuid.New().String(),
		JobOpeningID:       strings.TrimSpace(req.JobOpeningID),
		Name:               strings.TrimSpace(req.Name),
		IncludeCEOApproval: req.IncludeCEOApproval,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	if err := s.wf.CreateDefinition(def); err != nil {
		return nil, err
	}
	stages, err := s.buildDefinitionStages(def.ID, req.Stages)
	if err != nil {
		return nil, err
	}
	if err := s.wf.ReplaceDefinitionStages(def.ID, stages); err != nil {
		return nil, err
	}
	return s.wf.GetDefinitionByID(def.ID)
}

func (s *InterviewWorkflowService) buildDefinitionStages(definitionID string, dtos []models.InterviewWorkflowDefinitionStageDTO) ([]models.InterviewWorkflowDefinitionStage, error) {
	var out []models.InterviewWorkflowDefinitionStage
	seenSeq := map[int]bool{}
	for _, st := range dtos {
		if seenSeq[st.Sequence] {
			return nil, fmt.Errorf("duplicate stage sequence: %d", st.Sequence)
		}
		seenSeq[st.Sequence] = true
		if len(st.ScoringDimensions) == 0 {
			return nil, errors.New("each stage requires scoringDimensions")
		}
		jsonStr, err := s.marshalDimensions(st.ScoringDimensions)
		if err != nil {
			return nil, err
		}
		out = append(out, models.InterviewWorkflowDefinitionStage{
			ID:                    uuid.New().String(),
			DefinitionID:          definitionID,
			Sequence:              st.Sequence,
			Title:                 strings.TrimSpace(st.Title),
			Description:           st.Description,
			ScoringDimensionsJSON: jsonStr,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		})
	}
	return out, nil
}

func (s *InterviewWorkflowService) UpdateDefinition(definitionID string, req models.CreateInterviewWorkflowDefinitionRequest) (*models.InterviewWorkflowDefinition, error) {
	def, err := s.wf.GetDefinitionByID(definitionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("definition not found")
		}
		return nil, err
	}
	if def.JobOpeningID != req.JobOpeningID {
		return nil, errors.New("jobOpeningId cannot be changed")
	}
	def.Name = strings.TrimSpace(req.Name)
	def.IncludeCEOApproval = req.IncludeCEOApproval
	def.UpdatedAt = time.Now()
	if err := s.wf.UpdateDefinitionMeta(def); err != nil {
		return nil, err
	}
	stages, err := s.buildDefinitionStages(def.ID, req.Stages)
	if err != nil {
		return nil, err
	}
	if err := s.wf.ReplaceDefinitionStages(def.ID, stages); err != nil {
		return nil, err
	}
	return s.wf.GetDefinitionByID(def.ID)
}

func (s *InterviewWorkflowService) ListDefinitionsByJob(jobOpeningID string) ([]models.InterviewWorkflowDefinition, error) {
	return s.wf.ListDefinitionsByJobOpening(jobOpeningID)
}

func (s *InterviewWorkflowService) GetDefinition(id string) (*models.InterviewWorkflowDefinition, error) {
	return s.wf.GetDefinitionByID(id)
}

// --- Process ---

func (s *InterviewWorkflowService) maxSnapshotSequence(snaps []models.InterviewWorkflowProcessStageSnapshot) int {
	m := 0
	for _, sn := range snaps {
		if sn.Sequence > m {
			m = sn.Sequence
		}
	}
	return m
}

func (s *InterviewWorkflowService) StartProcess(actorEmployeeID string, req models.StartInterviewWorkflowProcessRequest) (*models.InterviewWorkflowProcess, error) {
	app, err := s.rec.GetApplicationByID(strings.TrimSpace(req.ApplicationID))
	if err != nil || app == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("application not found")
		}
		if err != nil {
			return nil, err
		}
		return nil, errors.New("application not found")
	}
	def, err := s.wf.GetDefinitionByID(strings.TrimSpace(req.DefinitionID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("workflow definition not found")
		}
		return nil, err
	}
	if def.JobOpeningID != app.JobOpeningID {
		return nil, errors.New("definition does not belong to this application's job opening")
	}
	if _, err := s.wf.GetProcessByApplicationID(app.ID); err == nil {
		return nil, errors.New("an interview workflow process already exists for this application")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if def.IncludeCEOApproval && strings.TrimSpace(req.FinalApprovals.CEOEmployeeID) == "" {
		return nil, errors.New("ceoEmployeeId is required when definition includes CEO approval")
	}
	if _, err := s.emp.FindByEmployeeID(strings.TrimSpace(req.FinalApprovals.HREmployeeID)); err != nil {
		return nil, errors.New("finalApprovals.hrEmployeeId not found")
	}
	if _, err := s.emp.FindByEmployeeID(strings.TrimSpace(req.FinalApprovals.ManagerEmployeeID)); err != nil {
		return nil, errors.New("finalApprovals.managerEmployeeId not found")
	}
	if def.IncludeCEOApproval {
		if _, err := s.emp.FindByEmployeeID(strings.TrimSpace(req.FinalApprovals.CEOEmployeeID)); err != nil {
			return nil, errors.New("finalApprovals.ceoEmployeeId not found")
		}
	}
	if _, err := s.emp.FindByEmployeeID(strings.TrimSpace(req.FirstStage.AssignedEmployeeID)); err != nil {
		return nil, errors.New("firstStage.assignedEmployeeId not found")
	}

	proc := &models.InterviewWorkflowProcess{
		ID:                        uuid.New().String(),
		ApplicationID:             app.ID,
		DefinitionID:              def.ID,
		JobOpeningID:              app.JobOpeningID,
		Status:                    models.IWProcessInProgress,
		IncludeCEOApproval:        def.IncludeCEOApproval,
		ApprovalHREmployeeID:      strings.TrimSpace(req.FinalApprovals.HREmployeeID),
		ApprovalManagerEmployeeID: strings.TrimSpace(req.FinalApprovals.ManagerEmployeeID),
		ApprovalCEOEmployeeID:     strings.TrimSpace(req.FinalApprovals.CEOEmployeeID),
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}
	if err := s.wf.CreateProcess(proc); err != nil {
		return nil, err
	}
	var snaps []models.InterviewWorkflowProcessStageSnapshot
	for _, st := range def.Stages {
		snaps = append(snaps, models.InterviewWorkflowProcessStageSnapshot{
			ID:                    uuid.New().String(),
			ProcessID:             proc.ID,
			Sequence:              st.Sequence,
			Title:                 st.Title,
			Description:           st.Description,
			ScoringDimensionsJSON: st.ScoringDimensionsJSON,
			CreatedAt:             time.Now(),
		})
	}
	if err := s.wf.CreateSnapshots(snaps); err != nil {
		return nil, err
	}

	sd, err := parseFlexibleDate(strings.TrimSpace(req.FirstStage.ScheduledDate))
	if err != nil || sd.IsZero() {
		return nil, errors.New("invalid firstStage.scheduledDate")
	}
	att := &models.InterviewWorkflowStageAttempt{
		ID:                 uuid.New().String(),
		ProcessID:          proc.ID,
		StageSequence:      1,
		AttemptNumber:      1,
		AssignedEmployeeID: strings.TrimSpace(req.FirstStage.AssignedEmployeeID),
		ScheduledDate:      sd,
		ScheduledTime:      strings.TrimSpace(req.FirstStage.ScheduledTime),
		DurationMin:        req.FirstStage.DurationMin,
		Mode:               req.FirstStage.Mode,
		Status:             models.IWAttemptStatusScheduled,
		AttendanceStatus:   models.IWAttendancePending,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	if err := s.wf.CreateAttempt(att); err != nil {
		return nil, err
	}
	s.appendAudit(proc.ID, att.ID, models.IWAuditProcessStarted, actorEmployeeID, "process started", map[string]string{"applicationId": app.ID})
	s.appendAudit(proc.ID, att.ID, models.IWAuditStageScheduled, actorEmployeeID, "stage 1 scheduled", nil)

	s.hub.Broadcast("interview_workflow.process_started", map[string]interface{}{
		"processId":     proc.ID,
		"applicationId": app.ID,
		"jobOpeningId":  app.JobOpeningID,
	})

	return s.wf.GetProcessByID(proc.ID)
}

func (s *InterviewWorkflowService) GetProcess(id string) (*models.InterviewWorkflowProcess, error) {
	return s.wf.GetProcessByID(id)
}

func (s *InterviewWorkflowService) GetProcessByApplication(applicationID string) (*models.InterviewWorkflowProcess, error) {
	return s.wf.GetProcessByApplicationID(applicationID)
}

func (s *InterviewWorkflowService) GetTimeline(processID string) (map[string]interface{}, error) {
	p, err := s.wf.GetProcessByID(processID)
	if err != nil {
		return nil, err
	}
	evs, err := s.wf.ListAuditByProcess(processID)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"process": p,
		"audit":   evs,
	}, nil
}

// ScheduleStage schedules the next stage attempt (HR). sequence must be unlocked by previous stage completion with proceed.
func (s *InterviewWorkflowService) ScheduleStage(processID string, actorEmployeeID string, req models.ScheduleInterviewWorkflowStageRequest) (*models.InterviewWorkflowStageAttempt, error) {
	proc, err := s.wf.GetProcessByID(processID)
	if err != nil {
		return nil, err
	}
	if proc.Status != models.IWProcessInProgress {
		return nil, errors.New("process is not accepting new stage schedules")
	}
	maxSeq := s.maxSnapshotSequence(proc.Snapshots)
	if req.StageSequence < 1 || req.StageSequence > maxSeq {
		return nil, errors.New("invalid stageSequence for this process")
	}
	if req.StageSequence > 1 && !s.previousStageProceeded(proc, req.StageSequence-1) {
		return nil, errors.New("previous stage must be completed with outcome proceed before scheduling this stage")
	}
	if s.hasOpenAttemptForSequence(proc, req.StageSequence) {
		return nil, errors.New("an open attempt already exists for this stage; complete or resolve it first")
	}
	if _, err := s.emp.FindByEmployeeID(strings.TrimSpace(req.AssignedEmployeeID)); err != nil {
		return nil, errors.New("assignedEmployeeId not found")
	}
	sd, err := parseFlexibleDate(strings.TrimSpace(req.ScheduledDate))
	if err != nil || sd.IsZero() {
		return nil, errors.New("invalid scheduledDate")
	}
	maxN, err := s.wf.MaxAttemptNumber(proc.ID, req.StageSequence)
	if err != nil {
		return nil, err
	}
	nextN := maxN + 1
	if nextN < 1 {
		nextN = 1
	}
	att := &models.InterviewWorkflowStageAttempt{
		ID:                 uuid.New().String(),
		ProcessID:          proc.ID,
		StageSequence:      req.StageSequence,
		AttemptNumber:      nextN,
		AssignedEmployeeID: strings.TrimSpace(req.AssignedEmployeeID),
		ScheduledDate:      sd,
		ScheduledTime:      strings.TrimSpace(req.ScheduledTime),
		DurationMin:        req.DurationMin,
		Mode:               req.Mode,
		Status:             models.IWAttemptStatusScheduled,
		AttendanceStatus:   models.IWAttendancePending,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	if err := s.wf.CreateAttempt(att); err != nil {
		return nil, err
	}
	s.appendAudit(proc.ID, att.ID, models.IWAuditStageScheduled, actorEmployeeID, fmt.Sprintf("stage %d attempt %d scheduled", req.StageSequence, nextN), nil)
	s.hub.Broadcast("interview_workflow.stage_scheduled", map[string]interface{}{
		"processId": proc.ID, "attemptId": att.ID, "stageSequence": req.StageSequence, "attemptNumber": nextN,
	})
	return att, nil
}

func (s *InterviewWorkflowService) hasOpenAttemptForSequence(proc *models.InterviewWorkflowProcess, seq int) bool {
	for _, a := range proc.Attempts {
		if a.StageSequence != seq {
			continue
		}
		if a.Status == models.IWAttemptStatusSuperseded {
			continue
		}
		if a.Status == models.IWAttemptStatusScheduled || a.Status == models.IWAttemptStatusInProgress {
			return true
		}
	}
	return false
}

func (s *InterviewWorkflowService) previousStageProceeded(proc *models.InterviewWorkflowProcess, prevSeq int) bool {
	var latest *models.InterviewWorkflowStageAttempt
	for i := range proc.Attempts {
		a := &proc.Attempts[i]
		if a.StageSequence != prevSeq || a.Status == models.IWAttemptStatusSuperseded {
			continue
		}
		if latest == nil || a.AttemptNumber > latest.AttemptNumber {
			latest = a
		}
	}
	return latest != nil && latest.Status == models.IWAttemptStatusCompleted && latest.StageOutcome == models.IWStageOutcomeProceed
}

// RecordAttendance — assignee only; cannot change after evaluation is locked.
func (s *InterviewWorkflowService) RecordAttendance(attemptID string, actorEmployeeID string, req models.RecordInterviewWorkflowAttendanceRequest) error {
	a, err := s.wf.GetAttemptByID(attemptID)
	if err != nil {
		return err
	}
	if err := s.requireAssignee(actorEmployeeID, a); err != nil {
		return err
	}
	if s.evaluationLocked(a) {
		return errors.New("stage evaluation is finalized; attendance cannot be changed")
	}
	if a.Status != models.IWAttemptStatusScheduled && a.Status != models.IWAttemptStatusInProgress {
		return errors.New("attendance can only be recorded for scheduled or in-progress attempts")
	}
	st := strings.TrimSpace(req.Status)
	switch st {
	case models.IWAttendancePresent, models.IWAttendanceAbsent, models.IWAttendanceExcused:
	default:
		return errors.New("status must be present, absent, or excused")
	}
	if a.AttendanceStatus != models.IWAttendancePending {
		return errors.New("attendance was already recorded for this attempt")
	}
	now := time.Now()
	a.AttendanceStatus = st
	a.AttendanceAt = &now
	a.AttendanceBy = actorEmployeeID
	a.UpdatedAt = now
	if err := s.wf.SaveAttempt(a); err != nil {
		return err
	}
	s.appendAudit(a.ProcessID, a.ID, models.IWAuditAttendanceRecorded, actorEmployeeID, req.Notes, map[string]string{"status": st})
	s.hub.Broadcast("interview_workflow.attendance", map[string]interface{}{
		"processId": a.ProcessID, "attemptId": a.ID, "status": st,
	})
	return nil
}

// StartInProgress — assignee only; requires present attendance.
func (s *InterviewWorkflowService) StartInProgress(attemptID string, actorEmployeeID string) error {
	a, err := s.wf.GetAttemptByID(attemptID)
	if err != nil {
		return err
	}
	if err := s.requireAssignee(actorEmployeeID, a); err != nil {
		return err
	}
	if s.evaluationLocked(a) {
		return errors.New("stage already finalized")
	}
	if a.Status != models.IWAttemptStatusScheduled {
		return errors.New("only scheduled attempts can be started")
	}
	if a.AttendanceStatus != models.IWAttendancePresent {
		return errors.New("candidate must be marked present before starting the interview")
	}
	now := time.Now()
	a.Status = models.IWAttemptStatusInProgress
	a.StartedAt = &now
	a.StartedByEmployeeID = actorEmployeeID
	a.UpdatedAt = now
	if err := s.wf.SaveAttempt(a); err != nil {
		return err
	}
	s.appendAudit(a.ProcessID, a.ID, models.IWAuditStageStarted, actorEmployeeID, "", nil)
	return nil
}

// SubmitEvaluation — assignee only; mandatory remarks + structured scores; locks stage.
func (s *InterviewWorkflowService) SubmitEvaluation(attemptID string, actorEmployeeID string, req models.SubmitInterviewWorkflowEvaluationRequest) error {
	a, err := s.wf.GetAttemptByID(attemptID)
	if err != nil {
		return err
	}
	if err := s.requireAssignee(actorEmployeeID, a); err != nil {
		return err
	}
	if s.evaluationLocked(a) {
		return errors.New("evaluation already submitted for this attempt")
	}
	if a.Status != models.IWAttemptStatusInProgress {
		return errors.New("interview must be in progress to submit evaluation")
	}
	if a.AttendanceStatus != models.IWAttendancePresent {
		return errors.New("evaluation requires attendance marked present")
	}
	if strings.TrimSpace(req.Remarks) == "" {
		return errors.New("remarks are mandatory")
	}
	out := strings.TrimSpace(req.Outcome)
	switch out {
	case models.IWStageOutcomeProceed, models.IWStageOutcomeRepeatStage, models.IWStageOutcomeFail:
	default:
		return errors.New("invalid outcome")
	}

	proc, err := s.wf.GetProcessByID(a.ProcessID)
	if err != nil {
		return err
	}
	if proc.Status != models.IWProcessInProgress {
		return errors.New("process is not accepting evaluations")
	}
	dims, err := s.dimensionsForProcessStage(proc, a.StageSequence)
	if err != nil {
		return err
	}
	if err := s.validateScores(dims, req.Scores); err != nil {
		return err
	}
	scoresJSON, err := json.Marshal(req.Scores)
	if err != nil {
		return err
	}

	now := time.Now()
	a.EvaluationScoresJSON = string(scoresJSON)
	a.EvaluationRemarks = strings.TrimSpace(req.Remarks)
	a.StageOutcome = out
	a.EvaluatedAt = &now
	a.EvaluatedByEmployeeID = actorEmployeeID
	a.Status = models.IWAttemptStatusCompleted
	a.UpdatedAt = now

	switch out {
	case models.IWStageOutcomeFail:
		if err := s.wf.SaveAttempt(a); err != nil {
			return err
		}
		proc.Status = models.IWProcessFailed
		proc.UpdatedAt = now
		_ = s.wf.SaveProcess(proc)
		s.appendAudit(proc.ID, a.ID, models.IWAuditEvaluationSubmitted, actorEmployeeID, "fail", map[string]string{"outcome": out})
		_ = s.failApplication(proc.ApplicationID)
		s.hub.Broadcast("interview_workflow.evaluation", map[string]interface{}{"processId": proc.ID, "attemptId": a.ID, "outcome": out})

	case models.IWStageOutcomeRepeatStage:
		if strings.TrimSpace(req.NextScheduledDate) == "" || strings.TrimSpace(req.NextScheduledTime) == "" {
			return errors.New("nextScheduledDate and nextScheduledTime are required when outcome is repeat_stage")
		}
		sd, err := parseFlexibleDate(strings.TrimSpace(req.NextScheduledDate))
		if err != nil || sd.IsZero() {
			return errors.New("invalid nextScheduledDate")
		}
		if err := s.wf.SaveAttempt(a); err != nil {
			return err
		}
		nextN := a.AttemptNumber + 1
		na := &models.InterviewWorkflowStageAttempt{
			ID:                 uuid.New().String(),
			ProcessID:          proc.ID,
			StageSequence:      a.StageSequence,
			AttemptNumber:      nextN,
			AssignedEmployeeID: a.AssignedEmployeeID,
			ScheduledDate:      sd,
			ScheduledTime:      strings.TrimSpace(req.NextScheduledTime),
			DurationMin:        req.DurationMin,
			Mode:               firstNonEmpty(req.Mode, a.Mode),
			Status:             models.IWAttemptStatusScheduled,
			AttendanceStatus:   models.IWAttendancePending,
			CreatedAt:          now,
			UpdatedAt:          now,
		}
		if err := s.wf.CreateAttempt(na); err != nil {
			return err
		}
		s.appendAudit(proc.ID, a.ID, models.IWAuditEvaluationSubmitted, actorEmployeeID, "repeat_stage", map[string]string{"nextAttemptId": na.ID})
		s.hub.Broadcast("interview_workflow.evaluation", map[string]interface{}{"processId": proc.ID, "attemptId": a.ID, "outcome": out})

	case models.IWStageOutcomeProceed:
		if err := s.wf.SaveAttempt(a); err != nil {
			return err
		}
		maxSeq := s.maxSnapshotSequence(proc.Snapshots)
		if a.StageSequence >= maxSeq {
			if err := s.spawnFinalApprovals(proc, now); err != nil {
				return err
			}
		}
		s.appendAudit(proc.ID, a.ID, models.IWAuditEvaluationSubmitted, actorEmployeeID, "proceed", nil)
		s.hub.Broadcast("interview_workflow.evaluation", map[string]interface{}{"processId": proc.ID, "attemptId": a.ID, "outcome": out})
	}
	return nil
}

func firstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}

func (s *InterviewWorkflowService) dimensionsForProcessStage(proc *models.InterviewWorkflowProcess, seq int) ([]models.ScoringDimension, error) {
	for _, sn := range proc.Snapshots {
		if sn.Sequence == seq {
			return s.parseDimensionsJSON(sn.ScoringDimensionsJSON)
		}
	}
	return nil, errors.New("stage snapshot not found")
}

func (s *InterviewWorkflowService) spawnFinalApprovals(proc *models.InterviewWorkflowProcess, now time.Time) error {
	approvals := []models.InterviewWorkflowFinalApproval{
		{
			ID: uuid.New().String(), ProcessID: proc.ID, StepOrder: 1, RoleCode: models.IWApprovalRoleHR,
			AssignedEmployeeID: proc.ApprovalHREmployeeID, Status: models.IWApprovalPending, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: uuid.New().String(), ProcessID: proc.ID, StepOrder: 2, RoleCode: models.IWApprovalRoleManager,
			AssignedEmployeeID: proc.ApprovalManagerEmployeeID, Status: models.IWApprovalPending, CreatedAt: now, UpdatedAt: now,
		},
	}
	if proc.IncludeCEOApproval && strings.TrimSpace(proc.ApprovalCEOEmployeeID) != "" {
		approvals = append(approvals, models.InterviewWorkflowFinalApproval{
			ID: uuid.New().String(), ProcessID: proc.ID, StepOrder: 3, RoleCode: models.IWApprovalRoleCEO,
			AssignedEmployeeID: strings.TrimSpace(proc.ApprovalCEOEmployeeID), Status: models.IWApprovalPending, CreatedAt: now, UpdatedAt: now,
		})
	}
	for i := range approvals {
		if err := s.wf.CreateFinalApproval(&approvals[i]); err != nil {
			return err
		}
	}
	proc.Status = models.IWProcessPendingApprovals
	proc.UpdatedAt = now
	return s.wf.SaveProcess(proc)
}

func (s *InterviewWorkflowService) failApplication(applicationID string) error {
	app, err := s.rec.GetApplicationByID(applicationID)
	if err != nil || app == nil {
		return err
	}
	app.Stage = models.StageRejected
	app.LastStatusDate = time.Now()
	app.UpdatedAt = time.Now()
	return s.rec.UpdateApplication(app)
}

// AbsenceAction — assignee only after marking absent.
func (s *InterviewWorkflowService) AbsenceAction(attemptID string, actorEmployeeID string, req models.InterviewWorkflowAbsenceActionRequest) error {
	a, err := s.wf.GetAttemptByID(attemptID)
	if err != nil {
		return err
	}
	if err := s.requireAssignee(actorEmployeeID, a); err != nil {
		return err
	}
	if s.evaluationLocked(a) {
		return errors.New("stage already finalized")
	}
	if a.AttendanceStatus != models.IWAttendanceAbsent && a.AttendanceStatus != models.IWAttendanceExcused {
		return errors.New("absence actions are only available when attendance is absent or excused")
	}
	if a.Status != models.IWAttemptStatusScheduled && a.Status != models.IWAttemptStatusInProgress {
		return errors.New("invalid attempt status for absence action")
	}
	if strings.TrimSpace(req.Notes) == "" {
		return errors.New("notes are mandatory")
	}
	proc, err := s.wf.GetProcessByID(a.ProcessID)
	if err != nil {
		return err
	}
	now := time.Now()
	act := strings.TrimSpace(req.Action)
	switch act {
	case "fail_candidate":
		a.Status = models.IWAttemptStatusMissed
		a.UpdatedAt = now
		if err := s.wf.SaveAttempt(a); err != nil {
			return err
		}
		proc.Status = models.IWProcessFailed
		proc.UpdatedAt = now
		_ = s.wf.SaveProcess(proc)
		s.appendAudit(proc.ID, a.ID, models.IWAuditAbsenceFail, actorEmployeeID, req.Notes, nil)
		_ = s.failApplication(proc.ApplicationID)
		s.hub.Broadcast("interview_workflow.absence", map[string]interface{}{"processId": proc.ID, "attemptId": a.ID, "action": act})
	case "reschedule":
		if strings.TrimSpace(req.NextScheduledDate) == "" || strings.TrimSpace(req.NextScheduledTime) == "" {
			return errors.New("nextScheduledDate and nextScheduledTime are required for reschedule")
		}
		sd, err := parseFlexibleDate(strings.TrimSpace(req.NextScheduledDate))
		if err != nil || sd.IsZero() {
			return errors.New("invalid nextScheduledDate")
		}
		a.Status = models.IWAttemptStatusSuperseded
		a.UpdatedAt = now
		if err := s.wf.SaveAttempt(a); err != nil {
			return err
		}
		nextN := a.AttemptNumber + 1
		na := &models.InterviewWorkflowStageAttempt{
			ID:                 uuid.New().String(),
			ProcessID:          proc.ID,
			StageSequence:      a.StageSequence,
			AttemptNumber:      nextN,
			AssignedEmployeeID: a.AssignedEmployeeID,
			ScheduledDate:      sd,
			ScheduledTime:      strings.TrimSpace(req.NextScheduledTime),
			DurationMin:        req.DurationMin,
			Mode:               firstNonEmpty(req.Mode, a.Mode),
			Status:             models.IWAttemptStatusScheduled,
			AttendanceStatus:   models.IWAttendancePending,
			CreatedAt:          now,
			UpdatedAt:          now,
		}
		if err := s.wf.CreateAttempt(na); err != nil {
			return err
		}
		s.appendAudit(proc.ID, a.ID, models.IWAuditAbsenceReschedule, actorEmployeeID, req.Notes, map[string]string{"newAttemptId": na.ID})
		s.hub.Broadcast("interview_workflow.absence", map[string]interface{}{"processId": proc.ID, "attemptId": na.ID, "action": act})
	default:
		return errors.New("invalid action: use fail_candidate or reschedule")
	}
	return nil
}

// DecideFinalApproval — only the assigned approver for that step; sequential chain.
func (s *InterviewWorkflowService) DecideFinalApproval(approvalID string, actorEmployeeID string, req models.InterviewWorkflowFinalApprovalDecisionRequest) error {
	ap, err := s.wf.GetFinalApprovalByID(approvalID)
	if err != nil {
		return err
	}
	if strings.TrimSpace(ap.AssignedEmployeeID) != strings.TrimSpace(actorEmployeeID) {
		return errors.New("only the assigned approver may decide this step")
	}
	if ap.Status != models.IWApprovalPending {
		return errors.New("this approval step is no longer pending")
	}
	proc, err := s.wf.GetProcessByID(ap.ProcessID)
	if err != nil {
		return err
	}
	if proc.Status != models.IWProcessPendingApprovals {
		return errors.New("process is not in final approval state")
	}
	all, err := s.wf.ListFinalApprovals(proc.ID)
	if err != nil {
		return err
	}
	for _, row := range all {
		if row.StepOrder < ap.StepOrder {
			if row.Status != models.IWApprovalApproved {
				return errors.New("previous approval steps must be completed before acting on this step")
			}
		}
	}
	dec := strings.TrimSpace(req.Decision)
	if dec != "approve" && dec != "reject" {
		return errors.New("decision must be approve or reject")
	}
	if strings.TrimSpace(req.Remarks) == "" {
		return errors.New("remarks are mandatory")
	}
	now := time.Now()
	switch dec {
	case "reject":
		ap.Status = models.IWApprovalRejected
		ap.Remarks = strings.TrimSpace(req.Remarks)
		ap.DecidedAt = &now
		ap.DecidedByEmployeeID = actorEmployeeID
		ap.UpdatedAt = now
		if err := s.wf.SaveFinalApproval(ap); err != nil {
			return err
		}
		proc.Status = models.IWProcessRejected
		proc.UpdatedAt = now
		_ = s.wf.SaveProcess(proc)
		_ = s.failApplication(proc.ApplicationID)
		s.appendAudit(proc.ID, "", models.IWAuditFinalApprovalDecided, actorEmployeeID, "rejected", map[string]string{"approvalId": ap.ID, "role": ap.RoleCode})
	case "approve":
		ap.Status = models.IWApprovalApproved
		ap.Remarks = strings.TrimSpace(req.Remarks)
		ap.DecidedAt = &now
		ap.DecidedByEmployeeID = actorEmployeeID
		ap.UpdatedAt = now
		if err := s.wf.SaveFinalApproval(ap); err != nil {
			return err
		}
		s.appendAudit(proc.ID, "", models.IWAuditFinalApprovalDecided, actorEmployeeID, "approved", map[string]string{"approvalId": ap.ID, "role": ap.RoleCode})
		// If last pending in chain -> process approved
		lastOrder := 0
		for _, row := range all {
			if row.StepOrder > lastOrder {
				lastOrder = row.StepOrder
			}
		}
		if ap.StepOrder == lastOrder {
			proc.Status = models.IWProcessApproved
			proc.UpdatedAt = now
			if err := s.wf.SaveProcess(proc); err != nil {
				return err
			}
			app, err := s.rec.GetApplicationByID(proc.ApplicationID)
			if err == nil && app != nil {
				app.Stage = models.StageInterviewCompleted
				app.LastStatusDate = now
				app.UpdatedAt = now
				_ = s.rec.UpdateApplication(app)
			}
		}
	}
	s.hub.Broadcast("interview_workflow.final_approval", map[string]interface{}{
		"processId": proc.ID, "approvalId": ap.ID, "decision": dec,
	})
	return nil
}
