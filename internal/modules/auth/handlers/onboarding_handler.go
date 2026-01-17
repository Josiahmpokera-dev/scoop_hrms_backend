package handlers

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/auth/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/auth/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// OnboardingHandler handles onboarding HTTP requests
type OnboardingHandler struct {
	onboardingService *services.OnboardingService
}

// NewOnboardingHandler creates a new onboarding handler
func NewOnboardingHandler() *OnboardingHandler {
	return &OnboardingHandler{
		onboardingService: services.NewOnboardingService(),
	}
}

// CompleteOnboarding handles the complete onboarding flow
// POST /api/v1/auth/onboard
func (h *OnboardingHandler) CompleteOnboarding(c *gin.Context) {
	var req models.OnboardingRequest

	// Bind complete onboarding request (signup + organization)
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	onboardingResponse, err := h.onboardingService.CompleteOnboarding(&req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Onboarding completed successfully. Welcome to HRMS!", onboardingResponse)
}

// GetSetupWizardStatus returns the current status of the setup wizard
// GET /api/v1/auth/setup-wizard/status
func (h *OnboardingHandler) GetSetupWizardStatus(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	status, err := h.onboardingService.GetSetupWizardStatus(userID.(uint))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Setup wizard status retrieved successfully", status)
}
