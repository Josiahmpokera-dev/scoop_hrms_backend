package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/auth/models"
	organizationModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organizations/models"
	organizationRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organizations/repositories"
	roleModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/models"
	roleRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/repositories"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	userRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// OnboardingService handles the complete onboarding flow
type OnboardingService struct {
	userRepo         *userRepos.UserRepository
	organizationRepo *organizationRepos.OrganizationRepository
	roleRepo         *roleRepos.RoleRepository
}

// NewOnboardingService creates a new onboarding service
func NewOnboardingService() *OnboardingService {
	return &OnboardingService{
		userRepo:         userRepos.NewUserRepository(),
		organizationRepo: organizationRepos.NewOrganizationRepository(),
		roleRepo:         roleRepos.NewRoleRepository(),
	}
}

// CompleteOnboarding handles the complete onboarding flow:
// 1. Create user account
// 2. Create organization
// 3. Assign super admin role to user
func (s *OnboardingService) CompleteOnboarding(req *models.OnboardingRequest) (*models.OnboardingResponse, error) {
	// Step 1: Check if user already exists
	if s.userRepo.ExistsByEmail(req.Email) {
		return nil, errors.New("user with this email already exists")
	}

	// Step 2: Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Step 3: Generate username from email if not provided
	username := req.Username
	if username == "" {
		emailParts := strings.Split(req.Email, "@")
		if len(emailParts) > 0 {
			username = emailParts[0]
		} else {
			username = strings.ToLower(req.FirstName + req.LastName)
		}
		username = strings.ToLower(strings.ReplaceAll(username, ".", ""))
		username = strings.ReplaceAll(username, "_", "")
		username = strings.ReplaceAll(username, "-", "")
	}

	// Step 4: Create user
	user := &userModels.User{
		Username:     username,
		Email:        req.Email,
		Password:     string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         userModels.RoleAdmin, // Set as admin initially
		IsActive:     true,
		Status:       "active",
		EmailVerified: false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Step 5: Generate organization code
	orgCode := s.generateOrganizationCode(req.Name)

	// Step 6: Create organization
	organization := &organizationModels.Organization{
		Code:              orgCode,
		Name:              req.Name,
		LegalName:         req.LegalName,
		RegistrationNumber: req.RegistrationNumber,
		Industry:          req.Industry,
		CompanySize:       req.CompanySize,
		Country:           req.Country,
		Status:            organizationModels.OrganizationStatusActive,
		IsActive:          true,
		UpdatedBy:         &user.ID,
	}

	// Set defaults
	if req.Timezone != nil {
		organization.Timezone = req.Timezone
	} else {
		defaultTimezone := "Africa/Dar_es_Salaam"
		organization.Timezone = &defaultTimezone
	}

	if req.CurrencyCode != nil {
		organization.CurrencyCode = req.CurrencyCode
	} else {
		defaultCurrency := "TZS"
		organization.CurrencyCode = &defaultCurrency
	}

	if err := s.organizationRepo.Create(organization); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// Step 7: Assign Super Admin role to user
	// Find or create super admin role
	superAdminRole, err := s.ensureSuperAdminRole(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to assign super admin role: %w", err)
	}

	// Assign role to user via user_roles table
	if err := s.assignRoleToUser(user.ID, superAdminRole.ID, user.ID); err != nil {
		// Log but don't fail - user is created and can be assigned role later
		fmt.Printf("Warning: Failed to assign role to user: %v\n", err)
	}

	// Step 9: Generate JWT token
	token, expiresIn, err := s.generateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &models.OnboardingResponse{
		User: &models.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
		},
		Organization: &models.OrganizationInfo{
			ID:     organization.ID,
			Code:   organization.Code,
			Name:   organization.Name,
			Status: string(organization.Status),
		},
		AccessToken:  token,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		SetupComplete: false, // Setup wizard will be next step
	}, nil
}

// GetSetupWizardStatus returns the current status of the setup wizard
func (s *OnboardingService) GetSetupWizardStatus(userID uint) (*models.SetupWizardStatus, error) {
	// Validate user exists
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check what's been set up
	steps := []models.SetupWizardStep{
		{
			Step:        "organization",
			Title:       "Organization Setup",
			Description: "Create your organization",
			Completed:   false,
			Required:    true,
		},
		{
			Step:        "departments",
			Title:       "Departments",
			Description: "Create your first department",
			Completed:   false,
			Required:    false,
		},
		{
			Step:        "locations",
			Title:       "Locations",
			Description: "Add your office locations",
			Completed:   false,
			Required:    false,
		},
		{
			Step:        "job_positions",
			Title:       "Job Positions",
			Description: "Define job positions",
			Completed:   false,
			Required:    false,
		},
		{
			Step:        "employees",
			Title:       "Employees",
			Description: "Add your first employee",
			Completed:   false,
			Required:    false,
		},
	}

	// Check organization exists
	orgs, _, err := s.organizationRepo.List(nil, 1, 1, nil)
	if err == nil && len(orgs) > 0 {
		steps[0].Completed = true
	}

	// TODO: Check other steps (departments, locations, etc.)
	// This would require querying those repositories

	// Calculate progress
	completedCount := 0
	for _, step := range steps {
		if step.Completed {
			completedCount++
		}
	}
	progress := (completedCount * 100) / len(steps)

	// Determine current step
	currentStep := "organization"
	for _, step := range steps {
		if !step.Completed {
			currentStep = step.Step
			break
		}
	}

	allCompleted := progress == 100

	return &models.SetupWizardStatus{
		Completed:   allCompleted,
		CurrentStep: currentStep,
		Steps:       steps,
		Progress:    progress,
	}, nil
}

// Helper functions

// generateOrganizationCode generates a unique organization code
func (s *OnboardingService) generateOrganizationCode(orgName string) string {
	// Take first 3-4 letters of each word, uppercase
	words := strings.Fields(orgName)
	code := ""
	for _, word := range words {
		if len(word) >= 3 {
			code += strings.ToUpper(word[:3])
		} else {
			code += strings.ToUpper(word)
		}
		if len(code) >= 10 {
			break
		}
	}
	
	// Add year suffix
	year := time.Now().Year()
	code = fmt.Sprintf("%s-%d", code, year)
	
	// Ensure uniqueness
	counter := 1
	originalCode := code
	for s.organizationRepo.ExistsByCode(code) {
		code = fmt.Sprintf("%s-%d", originalCode, counter)
		counter++
	}
	
	return code
}

// ensureSuperAdminRole ensures a super admin role exists
func (s *OnboardingService) ensureSuperAdminRole(createdBy uint) (*roleModels.Role, error) {
	// Try to find existing super admin role
	role, err := s.roleRepo.FindByCode("super_admin")
	if err == nil && role != nil {
		return role, nil
	}
	
	// Create super admin role if it doesn't exist
	superAdminRole := &roleModels.Role{
		Code:        "super_admin",
		Name:        "Super Administrator",
		Description: func() *string { desc := "Full system access with all permissions for the organization"; return &desc }(),
		UpdatedBy:   &createdBy,
	}

	if err := s.roleRepo.Create(superAdminRole); err != nil {
		return nil, fmt.Errorf("failed to create super admin role: %w", err)
	}

	// Assign all permissions to super admin role
	permRepo := roleRepos.NewPermissionRepository()
	allPerms, _, err := permRepo.List(1, 1000) // Get all permissions
	if err == nil && len(allPerms) > 0 {
		var permIDs []uint
		for _, perm := range allPerms {
			permIDs = append(permIDs, perm.ID)
		}
		if err := s.roleRepo.AssignPermissions(superAdminRole.ID, permIDs); err != nil {
			// Log but don't fail - role is created
			fmt.Printf("Warning: Failed to assign all permissions to super admin role: %v\n", err)
		}
	}

	return superAdminRole, nil
}

// assignRoleToUser assigns a role to a user via user_roles table
func (s *OnboardingService) assignRoleToUser(userID, roleID, assignedBy uint) error {
	userRole := &roleModels.UserRole{
		UserID:     userID,
		RoleID:     roleID,
		AssignedBy: &assignedBy,
		AssignedAt: time.Now(),
	}

	// Use database connection directly
	db := database.GetDB()
	return db.Create(userRole).Error
}

// generateToken generates a JWT token for the user
func (s *OnboardingService) generateToken(user *userModels.User) (string, int64, error) {
	expiryDuration, err := time.ParseDuration(config.AppConfig.JWT.Expiry)
	if err != nil {
		expiryDuration = 24 * time.Hour
	}

	expiresAt := time.Now().Add(expiryDuration)
	expiresIn := int64(expiryDuration.Seconds())

	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"email":     user.Email,
		"role":      string(user.Role),
		"exp":       expiresAt.Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWT.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresIn, nil
}
