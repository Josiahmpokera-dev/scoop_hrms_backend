package handlers

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/types"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// ActionHandler handles POST-only action-based API requests
type ActionHandler struct {
	Resource string
}

// NewActionHandler creates a new action handler
func NewActionHandler(resource string) *ActionHandler {
	return &ActionHandler{
		Resource: resource,
	}
}

// Handle processes action-based requests
func (h *ActionHandler) Handle(
	createHandler func(c *gin.Context, data map[string]interface{}) error,
	readHandler func(c *gin.Context, id string) error,
	updateHandler func(c *gin.Context, id string, data map[string]interface{}) error,
	deleteHandler func(c *gin.Context, id string) error,
	listHandler func(c *gin.Context, filters map[string]interface{}, pagination *types.PaginationRequest) error,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.APIRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "Invalid request format", err.Error())
			return
		}

		// Set resource if not provided
		if req.Resource == "" {
			req.Resource = h.Resource
		}

		switch req.Action {
		case "create":
			if createHandler == nil {
				response.BadRequest(c, "Create action not supported for this resource", nil)
				return
			}
			if err := createHandler(c, req.Data); err != nil {
				response.BadRequest(c, "Failed to create resource", err.Error())
				return
			}

		case "read":
			if readHandler == nil {
				response.BadRequest(c, "Read action not supported for this resource", nil)
				return
			}
			if req.ID == nil {
				response.BadRequest(c, "ID is required for read action", nil)
				return
			}
			if err := readHandler(c, *req.ID); err != nil {
				response.NotFound(c, "Resource not found")
				return
			}

		case "update":
			if updateHandler == nil {
				response.BadRequest(c, "Update action not supported for this resource", nil)
				return
			}
			if req.ID == nil {
				response.BadRequest(c, "ID is required for update action", nil)
				return
			}
			if err := updateHandler(c, *req.ID, req.Data); err != nil {
				response.BadRequest(c, "Failed to update resource", err.Error())
				return
			}

		case "delete":
			if deleteHandler == nil {
				response.BadRequest(c, "Delete action not supported for this resource", nil)
				return
			}
			if req.ID == nil {
				response.BadRequest(c, "ID is required for delete action", nil)
				return
			}
			if err := deleteHandler(c, *req.ID); err != nil {
				response.BadRequest(c, "Failed to delete resource", err.Error())
				return
			}

		case "list":
			if listHandler == nil {
				response.BadRequest(c, "List action not supported for this resource", nil)
				return
			}
			if err := listHandler(c, req.Filters, req.Pagination); err != nil {
				response.BadRequest(c, "Failed to list resources", err.Error())
				return
			}

		default:
			response.BadRequest(c, "Invalid action. Must be one of: create, read, update, delete, list", nil)
		}
	}
}
