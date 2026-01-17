package types

import (
	"time"
)

// APIResponse represents the standardized API response format
type APIResponse struct {
	Success bool                   `json:"success"`
	Data    interface{}           `json:"data,omitempty"`
	Message string                `json:"message,omitempty"`
	Error   interface{}           `json:"error,omitempty"`
	Meta    *ResponseMeta         `json:"meta,omitempty"`
}

// ResponseMeta contains metadata about the response
type ResponseMeta struct {
	RequestID  string              `json:"request_id,omitempty"`
	Pagination *PaginationMeta    `json:"pagination,omitempty"`
	Timestamp  time.Time           `json:"timestamp"`
}

// PaginationMeta contains pagination information
type PaginationMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int `json:"total_pages"`
}

// CalculateTotalPages calculates total pages from total items and page size
func (p *PaginationMeta) CalculateTotalPages() {
	if p.PageSize > 0 {
		p.TotalPages = int((p.Total + int64(p.PageSize) - 1) / int64(p.PageSize))
	} else {
		p.TotalPages = 1
	}
}
