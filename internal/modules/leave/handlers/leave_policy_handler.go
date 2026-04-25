package handlers

import (
	"fmt"
	"strconv"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/services"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type LeavePolicyHandler struct {
	service *services.LeavePolicyService
}

func NewLeavePolicyHandler() *LeavePolicyHandler {
	return &LeavePolicyHandler{
		service: services.NewLeavePolicyService(),
	}
}

// ListLeavePolicies handles listing leave policies
func (h *LeavePolicyHandler) ListLeavePolicies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize > 100 {
		pageSize = 100
	}

	tenantID := middleware.GetTenantID(c)

	// Build filters
	filters := make(map[string]interface{})
	if country := c.Query("country"); country != "" {
		filters["country"] = country
	}
	if leaveTypeCode := c.Query("leave_type_code"); leaveTypeCode != "" {
		filters["leave_type_code"] = leaveTypeCode
	}
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			filters["is_active"] = isActive
		}
	}

	policies, total, err := h.service.ListLeavePolicies(tenantID, page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve leave policies", err.Error())
		return
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	if totalPages == 0 {
		totalPages = 1
	}

	response.Success(c, "Leave policies retrieved successfully", map[string]interface{}{
		"data": policies,
		"meta": map[string]interface{}{
			"page":        page,
			"per_page":    pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetLeavePolicy handles getting a leave policy by ID
func (h *LeavePolicyHandler) GetLeavePolicy(c *gin.Context) {
	idStr := c.Param("policy_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid policy ID", nil)
		return
	}

	policy, err := h.service.GetLeavePolicyByID(uint(id))
	if err != nil {
		response.NotFound(c, "Leave policy not found")
		return
	}

	response.Success(c, "Leave policy retrieved successfully", policy)
}

// GetPolicyGuidelines handles getting policy guidelines
func (h *LeavePolicyHandler) GetPolicyGuidelines(c *gin.Context) {
	leaveTypeCode := c.Query("leave_type_code")
	if leaveTypeCode == "" {
		response.BadRequest(c, "leave_type_code is required", nil)
		return
	}

	country := c.DefaultQuery("country", "Tanzania")
	tenantID := middleware.GetTenantID(c)

	policy, err := h.service.GetPolicyGuidelines(leaveTypeCode, country, tenantID)
	if err != nil {
		response.NotFound(c, "Policy guidelines not found")
		return
	}

	// Format as guidelines response
	guidelines := []string{}
	if policy.MinimumNoticeDays > 0 {
		guidelines = append(guidelines, fmt.Sprintf("Minimum %d days notice required", policy.MinimumNoticeDays))
	}
	if policy.CarryForward && policy.CarryForwardLimit != nil {
		guidelines = append(guidelines, fmt.Sprintf("Carry forward max %d days", *policy.CarryForwardLimit))
	}
	if policy.HalfDayAllowed {
		guidelines = append(guidelines, "Can be applied in half-days")
	}

	response.Success(c, "Policy guidelines retrieved successfully", map[string]interface{}{
		"leave_type":               leaveTypeCode,
		"minimum_notice_days":      policy.MinimumNoticeDays,
		"maximum_days_per_request": policy.MaximumDaysPerRequest,
		"carry_forward_limit":      policy.CarryForwardLimit,
		"carry_forward_expiry":     policy.CarryForwardExpiry,
		"half_day_allowed":         policy.HalfDayAllowed,
		"requires_documentation":   policy.LeaveType.RequiresDocumentation,
		"documentation_types":      []string{},
		"guidelines":               guidelines,
	})
}

// CreateLeavePolicy handles creating a leave policy
func (h *LeavePolicyHandler) CreateLeavePolicy(c *gin.Context) {
	var req models.CreateLeavePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	user, _ := c.Get("user")
	var createdBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		createdBy = &userObj.ID
	}

	policy, err := h.service.CreateLeavePolicy(&req, tenantID, createdBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Leave policy created successfully", policy)
}

// UpdateLeavePolicy handles updating a leave policy
func (h *LeavePolicyHandler) UpdateLeavePolicy(c *gin.Context) {
	idStr := c.Param("policy_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid policy ID", nil)
		return
	}

	var req models.UpdateLeavePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, _ := c.Get("user")
	var updatedBy *uint
	if userObj, ok := user.(*userModels.User); ok {
		updatedBy = &userObj.ID
	}

	policy, err := h.service.UpdateLeavePolicy(uint(id), &req, updatedBy)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Leave policy updated successfully", policy)
}

// DeleteLeavePolicy handles deleting a leave policy
func (h *LeavePolicyHandler) DeleteLeavePolicy(c *gin.Context) {
	idStr := c.Param("policy_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid policy ID", nil)
		return
	}

	if err := h.service.DeleteLeavePolicy(uint(id)); err != nil {
		if err.Error() == "Cannot delete leave policy in use" {
			response.BadRequest(c, "Cannot delete leave policy in use", map[string]string{"policy_id": "assigned_to_employee"})
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Leave policy deleted successfully", nil)
}
