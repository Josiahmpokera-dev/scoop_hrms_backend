package services

import (
	"errors"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
	helpdeskRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/repositories"
)

// KnowledgeBaseService handles knowledge base business logic
type KnowledgeBaseService struct {
	kbRepo *helpdeskRepos.KnowledgeBaseRepository
}

// NewKnowledgeBaseService creates a new knowledge base service
func NewKnowledgeBaseService() *KnowledgeBaseService {
	return &KnowledgeBaseService{
		kbRepo: helpdeskRepos.NewKnowledgeBaseRepository(),
	}
}

// ListCategories returns KB categories with article counts
func (s *KnowledgeBaseService) ListCategories(tenantID *uint) ([]map[string]interface{}, error) {
	return s.kbRepo.GetCategories(tenantID)
}

// ListArticles returns KB articles with filters and pagination
func (s *KnowledgeBaseService) ListArticles(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.KnowledgeBaseArticle, int64, error) {
	return s.kbRepo.ListArticles(tenantID, page, pageSize, filters)
}

// GetArticleByID returns a KB article by ID and optionally increments views
func (s *KnowledgeBaseService) GetArticleByID(articleID uint, tenantID *uint, incrementViews bool) (*models.KnowledgeBaseArticle, error) {
	article, err := s.kbRepo.FindByID(articleID)
	if err != nil {
		return nil, errors.New("article not found")
	}
	if tenantID != nil && article.TenantID != nil && *article.TenantID != *tenantID {
		return nil, errors.New("article does not belong to your tenant")
	}
	if article.Status != models.KBArticleStatusPublished {
		return nil, errors.New("article not found")
	}
	if incrementViews {
		_ = s.kbRepo.IncrementViews(articleID)
		article.Views++
	}
	return article, nil
}

// SubmitFeedback records helpful/not helpful feedback for an article
func (s *KnowledgeBaseService) SubmitFeedback(articleID uint, userID uint, tenantID *uint, isHelpful bool) (helpful, notHelpful int, err error) {
	article, err := s.kbRepo.FindByID(articleID)
	if err != nil {
		return 0, 0, errors.New("article not found")
	}
	if tenantID != nil && article.TenantID != nil && *article.TenantID != *tenantID {
		return 0, 0, errors.New("article does not belong to your tenant")
	}

	existing, _ := s.kbRepo.GetFeedbackByUserAndArticle(articleID, userID)
	if existing != nil {
		_ = s.kbRepo.UpdateFeedback(articleID, userID, isHelpful)
	} else {
		_ = s.kbRepo.CreateFeedback(&models.KBArticleFeedback{
			TenantID:  tenantID,
			ArticleID: articleID,
			UserID:    userID,
			IsHelpful: isHelpful,
		})
	}
	_ = s.kbRepo.UpdateArticleFeedbackCounts(articleID)
	updated, _ := s.kbRepo.FindByID(articleID)
	return updated.Helpful, updated.NotHelpful, nil
}
