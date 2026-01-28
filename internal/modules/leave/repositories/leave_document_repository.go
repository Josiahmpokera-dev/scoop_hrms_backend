package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"gorm.io/gorm"
)

type LeaveDocumentRepository struct {
	db *gorm.DB
}

func NewLeaveDocumentRepository() *LeaveDocumentRepository {
	return &LeaveDocumentRepository{
		db: database.GetDB(),
	}
}

// Create creates a new leave document
func (r *LeaveDocumentRepository) Create(document *models.LeaveDocument) error {
	return r.db.Create(document).Error
}

// CreateBatch creates multiple documents
func (r *LeaveDocumentRepository) CreateBatch(documents []models.LeaveDocument) error {
	if len(documents) == 0 {
		return nil
	}
	return r.db.Create(&documents).Error
}

// FindByID finds a leave document by ID
func (r *LeaveDocumentRepository) FindByID(id uint) (*models.LeaveDocument, error) {
	var document models.LeaveDocument
	err := r.db.First(&document, id).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

// FindByLeaveRequestID finds all documents for a leave request
func (r *LeaveDocumentRepository) FindByLeaveRequestID(requestID uint) ([]models.LeaveDocument, error) {
	var documents []models.LeaveDocument
	err := r.db.Where("leave_request_id = ?", requestID).Find(&documents).Error
	return documents, err
}

// Delete deletes a leave document
func (r *LeaveDocumentRepository) Delete(id uint) error {
	return r.db.Delete(&models.LeaveDocument{}, id).Error
}

// DeleteByLeaveRequestID deletes all documents for a leave request
func (r *LeaveDocumentRepository) DeleteByLeaveRequestID(requestID uint) error {
	return r.db.Where("leave_request_id = ?", requestID).Delete(&models.LeaveDocument{}).Error
}
