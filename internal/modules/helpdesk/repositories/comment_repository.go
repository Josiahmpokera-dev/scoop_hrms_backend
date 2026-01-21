package repositories

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
	"gorm.io/gorm"
)

// CommentRepository handles comment database operations
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository creates a new comment repository
func NewCommentRepository() *CommentRepository {
	return &CommentRepository{
		db: database.GetDB(),
	}
}

// Create creates a new comment
func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

// FindByTicketID finds all comments for a ticket
func (r *CommentRepository) FindByTicketID(ticketID uint, includePrivate bool, userID *uint) ([]models.Comment, error) {
	var comments []models.Comment
	query := r.db.Where("ticket_id = ?", ticketID)

	// Filter private comments - only show if user is the author or if includePrivate is true and user is agent/admin
	if !includePrivate {
		query = query.Where("is_private = ?", false)
	} else if userID != nil {
		// Show private comments only to the author or agents/admin (we'll filter in service layer)
		query = query.Where("is_private = ? OR (is_private = ? AND author_id = ?)", false, true, *userID)
	} else {
		query = query.Where("is_private = ?", false)
	}

	err := query.Order("timestamp ASC").Find(&comments).Error
	return comments, err
}

// CountByTicketID counts comments for a ticket
func (r *CommentRepository) CountByTicketID(ticketID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Comment{}).Where("ticket_id = ? AND is_private = ?", ticketID, false).Count(&count).Error
	return count, err
}
