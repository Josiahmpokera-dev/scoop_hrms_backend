package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func stringifyErr(err interface{}) interface{} {
	if err == nil {
		return nil
	}
	switch e := err.(type) {
	case error:
		return e.Error()
	default:
		return fmt.Sprintf("%v", e)
	}
}

// APIResponse represents a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta contains pagination and metadata information
type Meta struct {
	Page       int   `json:"page,omitempty"`
	PerPage    int   `json:"per_page,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// Success sends a successful response
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta sends a successful response with pagination metadata
func SuccessWithMeta(c *gin.Context, message string, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Created sends a 201 Created response
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string, err interface{}) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, APIResponse{
		Success: false,
		Message: message,
		Error:   "Unauthorized access",
	})
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, APIResponse{
		Success: false,
		Message: message,
		Error:   "Forbidden",
	})
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Message: message,
		Error:   "Resource not found",
	})
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, message string, err interface{}) {
	c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Message: message,
		Error:   stringifyErr(err),
	})
}

// ServiceUnavailable sends a 503 Service Unavailable response
func ServiceUnavailable(c *gin.Context, message string, err interface{}) {
	c.JSON(http.StatusServiceUnavailable, APIResponse{
		Success: false,
		Message: message,
		Error:   stringifyErr(err),
	})
}

// ValidationError sends a 422 Unprocessable Entity response for validation errors
func ValidationError(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusUnprocessableEntity, APIResponse{
		Success: false,
		Message: message,
		Error:   errors,
	})
}
