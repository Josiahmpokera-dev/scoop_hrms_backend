package database

import (
	recruitmentModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/models"
)

// EnsureRecruitmentSchema runs AutoMigrate for recruitment tables so existing databases
// pick up new columns (e.g. job_openings.apply_token) without a separate migrator run.
func EnsureRecruitmentSchema() error {
	models := []interface{}{
		&recruitmentModels.JobRequisition{},
		&recruitmentModels.ApprovalStep{},
		&recruitmentModels.JobOpening{},
		&recruitmentModels.Candidate{},
		&recruitmentModels.JobApplication{},
		&recruitmentModels.ApplicationSubmissionQueue{},
		&recruitmentModels.Interview{},
		&recruitmentModels.Feedback{},
		&recruitmentModels.InterviewWorkflowDefinition{},
		&recruitmentModels.InterviewWorkflowDefinitionStage{},
		&recruitmentModels.InterviewWorkflowProcess{},
		&recruitmentModels.InterviewWorkflowProcessStageSnapshot{},
		&recruitmentModels.InterviewWorkflowStageAttempt{},
		&recruitmentModels.InterviewWorkflowFinalApproval{},
		&recruitmentModels.InterviewWorkflowAuditEvent{},
		&recruitmentModels.Offer{},
		&recruitmentModels.TalentPoolCandidate{},
	}
	return Migrate(models...)
}
