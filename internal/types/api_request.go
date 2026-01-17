package types

// APIRequest represents the standardized POST-only request format
type APIRequest struct {
	Action     string                 `json:"action" binding:"required,oneof=read create update delete list"`
	Resource   string                 `json:"resource,omitempty"` // Optional, can be inferred from endpoint
	ID         *string                `json:"id,omitempty"`      // For read, update, delete actions
	Data       map[string]interface{} `json:"data,omitempty"`     // For create, update actions
	Filters    map[string]interface{} `json:"filters,omitempty"` // For list action
	Pagination *PaginationRequest     `json:"pagination,omitempty"` // For list action
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int `json:"page,omitempty"`     // Page number (default: 1)
	PageSize int `json:"page_size,omitempty"` // Items per page (default: 20, max: 100)
}

// GetPage returns the page number, defaulting to 1
func (p *PaginationRequest) GetPage() int {
	if p == nil || p.Page < 1 {
		return 1
	}
	return p.Page
}

// GetPageSize returns the page size, defaulting to 20, max 100
func (p *PaginationRequest) GetPageSize() int {
	if p == nil || p.PageSize < 1 {
		return 20
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}
