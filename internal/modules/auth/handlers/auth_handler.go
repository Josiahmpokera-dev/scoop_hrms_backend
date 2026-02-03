package handlers

import (
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/auth/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/auth/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration data"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 422 {object} response.APIResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	authResponse, err := h.authService.Register(&req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "User registered successfully", authResponse)
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and get access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", err.Error())
		return
	}

	authResponse, err := h.authService.Login(&req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, "Login successful", authResponse)
}

// Refresh handles refresh token request and returns new access and refresh tokens
// @Summary Refresh token
// @Description Exchange a valid refresh token for new access and refresh tokens to stay logged in
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req models.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Validation failed", "refresh_token is required")
		return
	}

	authResponse, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, "Tokens refreshed successfully", authResponse)
}

// GetProfile handles getting user profile
// @Summary Get user profile
// @Description Get authenticated user's profile
// @Tags Authentication
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	profile, err := h.authService.GetProfile(userID.(uint))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Profile retrieved successfully", profile)
}

// Logout handles user logout.
// The client should discard the access token (e.g. remove from storage). JWTs are stateless;
// this endpoint returns success so the client can clear the token and optionally log the action in audit.
//
// @Summary Logout
// @Description Logout the current user. Client must discard the access token after calling this.
// @Tags Authentication
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Optional: could invalidate token in a blacklist here if you add one later
	response.Success(c, "Logged out successfully. Please discard the access token on the client.", gin.H{"logged_out": true})
}
