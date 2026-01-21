package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
	"gorm.io/gorm"
)

// AttachmentRepository handles attachment database operations
type AttachmentRepository struct {
	db *gorm.DB
}

// NewAttachmentRepository creates a new attachment repository
func NewAttachmentRepository() *AttachmentRepository {
	return &AttachmentRepository{
		db: database.GetDB(),
	}
}

// Create creates a new attachment
func (r *AttachmentRepository) Create(attachment *models.Attachment) error {
	return r.db.Create(attachment).Error
}

// FindByTicketID finds all attachments for a ticket
func (r *AttachmentRepository) FindByTicketID(ticketID uint) ([]models.Attachment, error) {
	var attachments []models.Attachment
	err := r.db.Where("ticket_id = ?", ticketID).Order("uploaded_at DESC").Find(&attachments).Error
	return attachments, err
}

// FindByID finds an attachment by ID
func (r *AttachmentRepository) FindByID(id uint) (*models.Attachment, error) {
	var attachment models.Attachment
	err := r.db.Where("id = ?", id).First(&attachment).Error
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}

// Delete deletes an attachment
func (r *AttachmentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Attachment{}, id).Error
}
