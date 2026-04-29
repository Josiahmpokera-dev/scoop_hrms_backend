package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	"gorm.io/gorm"
)

type EmailLogRepository struct {
	db *gorm.DB
}

func NewEmailLogRepository() *EmailLogRepository {
	return &EmailLogRepository{
		db: database.GetDB(),
	}
}

func (r *EmailLogRepository) CreateLog(log *models.EmailLog) error {
	return r.db.Create(log).Error
}

func (r *EmailLogRepository) GetLogs(page, pageSize int) ([]models.EmailLog, int64, error) {
	var logs []models.EmailLog
	var total int64

	query := r.db.Model(&models.EmailLog{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
