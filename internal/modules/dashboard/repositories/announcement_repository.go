package repositories

import (
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/dashboard/models"
	"gorm.io/gorm"
)

// AnnouncementRepository handles announcement database operations
type AnnouncementRepository struct {
	db *gorm.DB
}

// NewAnnouncementRepository creates a new announcement repository
func NewAnnouncementRepository() *AnnouncementRepository {
	return &AnnouncementRepository{
		db: database.DB,
	}
}

// GetAnnouncements retrieves paginated announcements
func (r *AnnouncementRepository) GetAnnouncements(tenantID *uint, page, pageSize int, announcementType string, activeOnly bool) ([]models.Announcement, int64, error) {
	var announcements []models.Announcement
	var total int64

	query := r.db.Model(&models.Announcement{})

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if announcementType != "" {
		query = query.Where("type = ?", announcementType)
	}

	if activeOnly {
		query = query.Where("is_active = ?", true)
		query = query.Where("expires_at IS NULL OR expires_at > ?", time.Now())
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := query.Preload("Attachments").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&announcements).Error; err != nil {
		return nil, 0, err
	}

	return announcements, total, nil
}

// GetByID retrieves an announcement by ID
func (r *AnnouncementRepository) GetByID(id uint, tenantID *uint) (*models.Announcement, error) {
	var announcement models.Announcement
	query := r.db.Preload("Attachments")

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if err := query.First(&announcement, id).Error; err != nil {
		return nil, err
	}
	return &announcement, nil
}

// Create creates a new announcement
func (r *AnnouncementRepository) Create(announcement *models.Announcement) error {
	return r.db.Create(announcement).Error
}

// Update updates an announcement
func (r *AnnouncementRepository) Update(announcement *models.Announcement) error {
	return r.db.Save(announcement).Error
}

// Delete soft deletes an announcement
func (r *AnnouncementRepository) Delete(id uint) error {
	return r.db.Delete(&models.Announcement{}, id).Error
}

// MarkAsRead marks an announcement as read by a user
func (r *AnnouncementRepository) MarkAsRead(announcementID, userID uint) error {
	read := models.AnnouncementRead{
		AnnouncementID: announcementID,
		UserID:         userID,
		ReadAt:         time.Now(),
	}
	// Use FirstOrCreate to avoid duplicates
	return r.db.Where("announcement_id = ? AND user_id = ?", announcementID, userID).
		FirstOrCreate(&read).Error
}

// IsRead checks if an announcement is read by a user
func (r *AnnouncementRepository) IsRead(announcementID, userID uint) bool {
	var count int64
	r.db.Model(&models.AnnouncementRead{}).
		Where("announcement_id = ? AND user_id = ?", announcementID, userID).
		Count(&count)
	return count > 0
}

// GetUnreadCount returns the count of unread announcements for a user
func (r *AnnouncementRepository) GetUnreadCount(tenantID *uint, userID uint) (int64, error) {
	var count int64

	subQuery := r.db.Model(&models.AnnouncementRead{}).
		Select("announcement_id").
		Where("user_id = ?", userID)

	query := r.db.Model(&models.Announcement{}).
		Where("is_active = ?", true).
		Where("expires_at IS NULL OR expires_at > ?", time.Now()).
		Where("id NOT IN (?)", subQuery)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetReadAt returns the read timestamp for an announcement by a user
func (r *AnnouncementRepository) GetReadAt(announcementID, userID uint) *time.Time {
	var read models.AnnouncementRead
	err := r.db.Where("announcement_id = ? AND user_id = ?", announcementID, userID).
		First(&read).Error
	if err != nil {
		return nil
	}
	return &read.ReadAt
}
