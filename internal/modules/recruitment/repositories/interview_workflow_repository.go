package repositories

import (
	"errors"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/models"
	"gorm.io/gorm"
)

// InterviewWorkflowRepository persists structured interview workflows.
type InterviewWorkflowRepository struct {
	db *gorm.DB
}

func NewInterviewWorkflowRepository() *InterviewWorkflowRepository {
	return &InterviewWorkflowRepository{db: database.GetDB()}
}

func (r *InterviewWorkflowRepository) CreateDefinition(def *models.InterviewWorkflowDefinition) error {
	return r.db.Create(def).Error
}

func (r *InterviewWorkflowRepository) ReplaceDefinitionStages(definitionID string, stages []models.InterviewWorkflowDefinitionStage) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("definition_id = ?", definitionID).Delete(&models.InterviewWorkflowDefinitionStage{}).Error; err != nil {
			return err
		}
		for i := range stages {
			if err := tx.Create(&stages[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *InterviewWorkflowRepository) GetDefinitionByID(id string) (*models.InterviewWorkflowDefinition, error) {
	var def models.InterviewWorkflowDefinition
	err := r.db.Preload("Stages", func(db *gorm.DB) *gorm.DB {
		return db.Order("sequence ASC")
	}).Where("id = ?", id).First(&def).Error
	if err != nil {
		return nil, err
	}
	return &def, nil
}

func (r *InterviewWorkflowRepository) ListDefinitionsByJobOpening(jobOpeningID string) ([]models.InterviewWorkflowDefinition, error) {
	var rows []models.InterviewWorkflowDefinition
	err := r.db.Preload("Stages", func(db *gorm.DB) *gorm.DB {
		return db.Order("sequence ASC")
	}).Where("job_opening_id = ?", jobOpeningID).Order("created_at DESC").Find(&rows).Error
	return rows, err
}

func (r *InterviewWorkflowRepository) UpdateDefinitionMeta(def *models.InterviewWorkflowDefinition) error {
	return r.db.Model(def).Updates(map[string]interface{}{
		"name":                 def.Name,
		"include_ceo_approval": def.IncludeCEOApproval,
		"updated_at":           def.UpdatedAt,
	}).Error
}

func (r *InterviewWorkflowRepository) CreateProcess(proc *models.InterviewWorkflowProcess) error {
	return r.db.Create(proc).Error
}

func (r *InterviewWorkflowRepository) CreateSnapshots(snaps []models.InterviewWorkflowProcessStageSnapshot) error {
	for i := range snaps {
		if err := r.db.Create(&snaps[i]).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *InterviewWorkflowRepository) GetProcessByID(id string) (*models.InterviewWorkflowProcess, error) {
	var p models.InterviewWorkflowProcess
	err := r.db.Preload("Snapshots", func(db *gorm.DB) *gorm.DB {
		return db.Order("sequence ASC")
	}).Preload("Attempts", func(db *gorm.DB) *gorm.DB {
		return db.Order("stage_sequence ASC, attempt_number ASC")
	}).Preload("Approvals", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_order ASC")
	}).Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *InterviewWorkflowRepository) GetProcessByApplicationID(applicationID string) (*models.InterviewWorkflowProcess, error) {
	var p models.InterviewWorkflowProcess
	err := r.db.Preload("Snapshots", func(db *gorm.DB) *gorm.DB {
		return db.Order("sequence ASC")
	}).Preload("Attempts", func(db *gorm.DB) *gorm.DB {
		return db.Order("stage_sequence ASC, attempt_number ASC")
	}).Preload("Approvals", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_order ASC")
	}).Where("application_id = ?", applicationID).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *InterviewWorkflowRepository) SaveProcess(proc *models.InterviewWorkflowProcess) error {
	return r.db.Save(proc).Error
}

func (r *InterviewWorkflowRepository) CreateAttempt(a *models.InterviewWorkflowStageAttempt) error {
	return r.db.Create(a).Error
}

func (r *InterviewWorkflowRepository) SaveAttempt(a *models.InterviewWorkflowStageAttempt) error {
	return r.db.Save(a).Error
}

func (r *InterviewWorkflowRepository) GetAttemptByID(id string) (*models.InterviewWorkflowStageAttempt, error) {
	var a models.InterviewWorkflowStageAttempt
	err := r.db.Where("id = ?", id).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// MaxAttemptNumber returns highest attempt_number for a stage sequence on a process.
func (r *InterviewWorkflowRepository) MaxAttemptNumber(processID string, stageSequence int) (int, error) {
	var a models.InterviewWorkflowStageAttempt
	err := r.db.Where("process_id = ? AND stage_sequence = ?", processID, stageSequence).
		Order("attempt_number DESC").First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return a.AttemptNumber, nil
}

func (r *InterviewWorkflowRepository) AppendAudit(ev *models.InterviewWorkflowAuditEvent) error {
	return r.db.Create(ev).Error
}

func (r *InterviewWorkflowRepository) ListAuditByProcess(processID string) ([]models.InterviewWorkflowAuditEvent, error) {
	var rows []models.InterviewWorkflowAuditEvent
	err := r.db.Where("process_id = ?", processID).Order("created_at ASC, id ASC").Find(&rows).Error
	return rows, err
}

func (r *InterviewWorkflowRepository) CreateFinalApproval(a *models.InterviewWorkflowFinalApproval) error {
	return r.db.Create(a).Error
}

func (r *InterviewWorkflowRepository) GetFinalApprovalByID(id string) (*models.InterviewWorkflowFinalApproval, error) {
	var a models.InterviewWorkflowFinalApproval
	err := r.db.Where("id = ?", id).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *InterviewWorkflowRepository) SaveFinalApproval(a *models.InterviewWorkflowFinalApproval) error {
	return r.db.Save(a).Error
}

func (r *InterviewWorkflowRepository) ListFinalApprovals(processID string) ([]models.InterviewWorkflowFinalApproval, error) {
	var rows []models.InterviewWorkflowFinalApproval
	err := r.db.Where("process_id = ?", processID).Order("step_order ASC").Find(&rows).Error
	return rows, err
}
