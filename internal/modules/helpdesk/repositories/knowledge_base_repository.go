package repositories

import (
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
	"gorm.io/gorm"
)

// KnowledgeBaseRepository handles knowledge base database operations
type KnowledgeBaseRepository struct {
	db *gorm.DB
}

// NewKnowledgeBaseRepository creates a new knowledge base repository
func NewKnowledgeBaseRepository() *KnowledgeBaseRepository {
	return &KnowledgeBaseRepository{
		db: database.GetDB(),
	}
}

// Create creates a new KB article
func (r *KnowledgeBaseRepository) Create(article *models.KnowledgeBaseArticle) error {
	return r.db.Create(article).Error
}

// FindByID finds a KB article by ID
func (r *KnowledgeBaseRepository) FindByID(id uint) (*models.KnowledgeBaseArticle, error) {
	var article models.KnowledgeBaseArticle
	err := r.db.Where("id = ?", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// ListArticles lists KB articles with filters
func (r *KnowledgeBaseRepository) ListArticles(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.KnowledgeBaseArticle, int64, error) {
	var articles []models.KnowledgeBaseArticle
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.KnowledgeBaseArticle{})

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Only show published articles for non-admin users
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	} else if _, isAdmin := filters["is_admin"]; !isAdmin {
		query = query.Where("status = ?", string(models.KBArticleStatusPublished))
	}

	// Apply filters
	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	if featured, ok := filters["featured"].(bool); ok {
		query = query.Where("featured = ?", featured)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortBy := "updated_at"
	sortDir := "DESC"
	if sb, ok := filters["sort_by"].(string); ok && sb != "" {
		sortBy = sb
	}
	if sd, ok := filters["sort_dir"].(string); ok && sd != "" {
		sortDir = sd
	}
	orderBy := fmt.Sprintf("%s %s", sortBy, sortDir)

	// Get paginated results
	err := query.Order(orderBy).Offset(offset).Limit(pageSize).Find(&articles).Error
	return articles, total, err
}

// Update updates a KB article
func (r *KnowledgeBaseRepository) Update(article *models.KnowledgeBaseArticle) error {
	return r.db.Save(article).Error
}

// IncrementViews increments the view count
func (r *KnowledgeBaseRepository) IncrementViews(id uint) error {
	return r.db.Model(&models.KnowledgeBaseArticle{}).Where("id = ?", id).UpdateColumn("views", gorm.Expr("views + 1")).Error
}

// GetCategories gets list of categories with article counts
func (r *KnowledgeBaseRepository) GetCategories(tenantID *uint) ([]map[string]interface{}, error) {
	var results []struct {
		Category string
		Count    int64
	}

	query := r.db.Model(&models.KnowledgeBaseArticle{}).
		Select("category, COUNT(*) as count").
		Where("status = ?", string(models.KBArticleStatusPublished))

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	err := query.Group("category").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	categories := make([]map[string]interface{}, len(results))
	for i, r := range results {
		categories[i] = map[string]interface{}{
			"name":  r.Category,
			"count": r.Count,
		}
	}

	return categories, nil
}

// GenerateArticleNumber generates a unique article number
func (r *KnowledgeBaseRepository) GenerateArticleNumber(tenantID *uint) (string, error) {
	var count int64
	query := r.db.Model(&models.KnowledgeBaseArticle{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Count(&count).Error; err != nil {
		return "", err
	}

	// Format: KB-XXX (e.g., KB-001)
	sequence := int(count) + 1
	return "KB-" + formatKBSequence(sequence), nil
}

// formatKBSequence formats sequence number with leading zeros
func formatKBSequence(seq int) string {
	if seq < 10 {
		return "00" + fmt.Sprintf("%d", seq)
	} else if seq < 100 {
		return "0" + fmt.Sprintf("%d", seq)
	}
	return fmt.Sprintf("%d", seq)
}

// CreateFeedback creates feedback for an article
func (r *KnowledgeBaseRepository) CreateFeedback(feedback *models.KBArticleFeedback) error {
	return r.db.Create(feedback).Error
}

// UpdateFeedback updates feedback (if user already gave feedback)
func (r *KnowledgeBaseRepository) UpdateFeedback(articleID, userID uint, isHelpful bool) error {
	return r.db.Model(&models.KBArticleFeedback{}).
		Where("article_id = ? AND user_id = ?", articleID, userID).
		Updates(map[string]interface{}{
			"is_helpful": isHelpful,
			"updated_at": time.Now(),
		}).Error
}

// GetFeedbackByUserAndArticle gets feedback for a specific user and article
func (r *KnowledgeBaseRepository) GetFeedbackByUserAndArticle(articleID, userID uint) (*models.KBArticleFeedback, error) {
	var feedback models.KBArticleFeedback
	err := r.db.Where("article_id = ? AND user_id = ?", articleID, userID).First(&feedback).Error
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

// UpdateArticleFeedbackCounts updates helpful/not helpful counts
func (r *KnowledgeBaseRepository) UpdateArticleFeedbackCounts(articleID uint) error {
	var helpfulCount, notHelpfulCount int64
	r.db.Model(&models.KBArticleFeedback{}).Where("article_id = ? AND is_helpful = ?", articleID, true).Count(&helpfulCount)
	r.db.Model(&models.KBArticleFeedback{}).Where("article_id = ? AND is_helpful = ?", articleID, false).Count(&notHelpfulCount)

	return r.db.Model(&models.KnowledgeBaseArticle{}).Where("id = ?", articleID).
		Updates(map[string]interface{}{
			"helpful":    helpfulCount,
			"not_helpful": notHelpfulCount,
		}).Error
}
