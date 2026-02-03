package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// KnowledgeBaseHandler handles knowledge base HTTP requests
type KnowledgeBaseHandler struct {
	service *services.KnowledgeBaseService
}

// NewKnowledgeBaseHandler creates a new knowledge base handler
func NewKnowledgeBaseHandler() *KnowledgeBaseHandler {
	return &KnowledgeBaseHandler{
		service: services.NewKnowledgeBaseService(),
	}
}

// ListCategories returns KB categories with article counts (Browse by Category)
// GET /helpdesk/knowledge-base/categories
func (h *KnowledgeBaseHandler) ListCategories(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	categories, err := h.service.ListCategories(tenantID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Doc shape: [{ name, count }, ...]
	response.Success(c, "Categories retrieved successfully", categories)
}

// ListArticles returns KB articles with pagination and filters
// GET /helpdesk/knowledge-base/articles
func (h *KnowledgeBaseHandler) ListArticles(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	pageSize := 20
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	filters := map[string]interface{}{}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if featured := c.Query("featured"); featured == "true" || featured == "1" {
		filters["featured"] = true
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if sortBy := c.Query("sort_by"); sortBy != "" {
		switch sortBy {
		case "updatedAt":
			filters["sort_by"] = "updated_at"
		case "views":
			filters["sort_by"] = "views"
		default:
			filters["sort_by"] = sortBy
		}
	}
	if sortDir := c.Query("sort_dir"); sortDir != "" {
		filters["sort_dir"] = sortDir
	}

	articles, total, err := h.service.ListArticles(tenantID, page, pageSize, filters)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	formatted := make([]map[string]interface{}, 0, len(articles))
	for i := range articles {
		a := &articles[i]
		formatted = append(formatted, formatKBArticleForList(a))
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	payload := map[string]interface{}{
		"data": formatted,
		"meta": map[string]interface{}{
			"page":        page,
			"per_page":    pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	}
	response.Success(c, "Articles retrieved successfully", payload)
}

// GetArticle returns a single KB article (and increments views)
// GET /helpdesk/knowledge-base/articles/:article_id
func (h *KnowledgeBaseHandler) GetArticle(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	articleIDStr := c.Param("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid article ID", nil)
		return
	}

	article, err := h.service.GetArticleByID(uint(articleID), tenantID, true)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Article retrieved successfully", formatKBArticleDetail(article))
}

// SubmitFeedback records helpful/not helpful for an article
// POST /helpdesk/knowledge-base/articles/:article_id/feedback
func (h *KnowledgeBaseHandler) SubmitFeedback(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	userObj, ok := user.(*userModels.User)
	if !ok {
		response.Unauthorized(c, "Invalid user context")
		return
	}
	tenantID := middleware.GetTenantID(c)

	articleIDStr := c.Param("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid article ID", nil)
		return
	}

	var req struct {
		IsHelpful bool `json:"is_helpful"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	helpful, notHelpful, err := h.service.SubmitFeedback(uint(articleID), userObj.ID, tenantID, req.IsHelpful)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Feedback recorded successfully", map[string]interface{}{
		"article_id": articleID,
		"helpful":    helpful,
		"notHelpful": notHelpful,
	})
}

func formatKBArticleForList(a *models.KnowledgeBaseArticle) map[string]interface{} {
	preview := ""
	if a.ContentPreview != nil {
		preview = *a.ContentPreview
	} else if len(a.Content) > 200 {
		preview = a.Content[:200] + "..."
	} else {
		preview = a.Content
	}
	var tags []string
	if a.TagsJSON != nil && *a.TagsJSON != "" {
		_ = json.Unmarshal([]byte(*a.TagsJSON), &tags)
	}
	return map[string]interface{}{
		"id":              a.ID,
		"title":           a.Title,
		"content_preview": preview,
		"category":        a.Category,
		"tags":            tags,
		"views":           a.Views,
		"helpful":         a.Helpful,
		"notHelpful":      a.NotHelpful,
		"author":          a.Author,
		"createdAt":       a.CreatedAt,
		"updatedAt":       a.UpdatedAt,
		"status":          string(a.Status),
		"featured":        a.Featured,
	}
}

func formatKBArticleDetail(a *models.KnowledgeBaseArticle) map[string]interface{} {
	var tags []string
	if a.TagsJSON != nil && *a.TagsJSON != "" {
		_ = json.Unmarshal([]byte(*a.TagsJSON), &tags)
	}
	return map[string]interface{}{
		"id":         a.ID,
		"title":      a.Title,
		"content":    a.Content,
		"category":   a.Category,
		"tags":       tags,
		"views":      a.Views,
		"helpful":    a.Helpful,
		"notHelpful": a.NotHelpful,
		"author":     a.Author,
		"createdAt":  a.CreatedAt,
		"updatedAt":  a.UpdatedAt,
		"status":     string(a.Status),
		"featured":   a.Featured,
	}
}
