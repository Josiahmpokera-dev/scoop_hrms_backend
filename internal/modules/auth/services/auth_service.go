package services

import (
	"errors"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/auth/models"
	roleRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/repositories"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repositories.UserRepository
	roleRepo *roleRepos.RoleRepository
}

// NewAuthService creates a new authentication service
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repositories.NewUserRepository(),
		roleRepo: roleRepos.NewRoleRepository(),
	}
}

// buildUserRoles returns a deduplicated slice of role strings: legacy user_type plus assigned role codes.
// Always includes "user" as a base role since all system users are treated as users.
func (s *AuthService) buildUserRoles(user *userModels.User) []string {
	seen := make(map[string]bool)
	var roles []string

	// Add the legacy role first
	legacy := strings.ToLower(string(user.Role))
	if legacy != "" && !seen[legacy] {
		seen[legacy] = true
		roles = append(roles, legacy)
	}

	// Add RBAC roles from the database
	codes, err := s.roleRepo.GetUserRoleCodes(user.ID)
	if err == nil {
		for _, code := range codes {
			c := strings.ToLower(strings.TrimSpace(code))
			if c != "" && !seen[c] {
				seen[c] = true
				roles = append(roles, c)
			}
		}
	}

	// Always ensure "user" is in the roles array - all system users are treated as users
	if !seen["user"] {
		roles = append(roles, "user")
	}

	return roles
}

// primaryRole returns the first role for backward compatibility.
func (s *AuthService) primaryRole(user *userModels.User) string {
	roles := s.buildUserRoles(user)
	if len(roles) > 0 {
		return roles[0]
	}
	return string(user.Role)
}

// Register creates a new user account
func (s *AuthService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	if s.userRepo.ExistsByEmail(req.Email) {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Generate username from email if not provided
	username := req.Username
	if username == "" {
		// Extract username from email (part before @)
		emailParts := strings.Split(req.Email, "@")
		if len(emailParts) > 0 {
			username = emailParts[0]
		} else {
			// Fallback: use first name + last name
			username = strings.ToLower(req.FirstName + req.LastName)
		}
		// Remove any special characters and make it lowercase
		username = strings.ToLower(strings.ReplaceAll(username, ".", ""))
		username = strings.ReplaceAll(username, "_", "")
		username = strings.ReplaceAll(username, "-", "")
	}

	// Determine role
	role := userModels.RoleUser
	if req.Role == "admin" {
		role = userModels.RoleAdmin
	} else if req.Role == "hr" {
		role = userModels.RoleHR
	}

	// Create user
	user := &userModels.User{
		Username:  username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      role,
		IsActive:  true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate access and refresh tokens
	accessToken, expiresIn, err := s.generateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	refreshToken, refreshExpiresIn, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	roles := s.buildUserRoles(user)
	return &models.AuthResponse{
		User: &models.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      s.primaryRole(user),
			Roles:     roles,
			IsActive:  user.IsActive,
		},
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        expiresIn,
		RefreshExpiresIn: refreshExpiresIn,
	}, nil
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user can login (active, not suspended or blocked)
	if !user.CanLogin() {
		if !user.IsActive {
			return nil, errors.New("account is deactivated")
		}
		switch user.Status {
		case "suspended":
			return nil, errors.New("account is suspended")
		case "blocked":
			return nil, errors.New("account is blocked")
		default:
			return nil, errors.New("account is not active")
		}
	}

	// Temporarily locked (rate limit lockout window)
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		return nil, errors.New("account is temporarily locked due to too many failed login attempts")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		maxAttempts := 5
		lockoutMinutes := 0
		if config.AppConfig != nil {
			maxAttempts = config.AppConfig.LoginRateLimit.MaxAttempts
			if maxAttempts <= 0 {
				maxAttempts = 5
			}
			lockoutMinutes = config.AppConfig.LoginRateLimit.LockoutMinutes
		}
		user.FailedLoginCount++
		if user.FailedLoginCount >= maxAttempts {
			user.Status = userModels.UserStatusBlocked
			if lockoutMinutes > 0 {
				until := time.Now().Add(time.Duration(lockoutMinutes) * time.Minute)
				user.LockedUntil = &until
			}
		}
		_ = s.userRepo.Update(user)
		return nil, errors.New("invalid email or password")
	}

	// Success: clear failed attempt count and lock
	user.FailedLoginCount = 0
	user.LockedUntil = nil
	_ = s.userRepo.Update(user)

	// Update last login
	_ = s.userRepo.UpdateLastLogin(user.ID)

	// Generate access and refresh tokens
	accessToken, expiresIn, err := s.generateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	refreshToken, refreshExpiresIn, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	roles := s.buildUserRoles(user)
	return &models.AuthResponse{
		User: &models.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      s.primaryRole(user),
			Roles:     roles,
			IsActive:  user.IsActive,
		},
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        expiresIn,
		RefreshExpiresIn: refreshExpiresIn,
	}, nil
}

// RefreshToken validates the refresh token and returns new access and refresh tokens
func (s *AuthService) RefreshToken(refreshTokenString string) (*models.AuthResponse, error) {
	userID, err := s.validateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.CanLogin() {
		if !user.IsActive {
			return nil, errors.New("account is deactivated")
		}
		switch user.Status {
		case "suspended":
			return nil, errors.New("account is suspended")
		case "blocked":
			return nil, errors.New("account is blocked")
		default:
			return nil, errors.New("account is not active")
		}
	}

	accessToken, expiresIn, err := s.generateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	newRefreshToken, refreshExpiresIn, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	roles := s.buildUserRoles(user)
	return &models.AuthResponse{
		User: &models.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      s.primaryRole(user),
			Roles:     roles,
			IsActive:  user.IsActive,
		},
		AccessToken:      accessToken,
		RefreshToken:     newRefreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        expiresIn,
		RefreshExpiresIn: refreshExpiresIn,
	}, nil
}

// GetProfile returns user profile
func (s *AuthService) GetProfile(userID uint) (*models.UserInfo, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	roles := s.buildUserRoles(user)
	return &models.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      s.primaryRole(user),
		Roles:     roles,
		IsActive:  user.IsActive,
	}, nil
}

// generateToken generates a JWT token for the user
func (s *AuthService) generateToken(user *userModels.User) (string, int64, error) {
	// Parse expiry duration
	expiryDuration, err := time.ParseDuration(config.AppConfig.JWT.Expiry)
	if err != nil {
		expiryDuration = 24 * time.Hour // Default to 24 hours
	}

	expiresAt := time.Now().Add(expiryDuration)
	expiresIn := int64(expiryDuration.Seconds())

	roles := s.buildUserRoles(user)
	primary := s.primaryRole(user)

	// Create claims (role = primary for backward compat, roles = full array)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    primary,
		"roles":   roles,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWT.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresIn, nil
}

// generateRefreshToken generates a JWT refresh token (longer-lived, type "refresh")
func (s *AuthService) generateRefreshToken(user *userModels.User) (string, int64, error) {
	expiryStr := config.AppConfig.JWT.RefreshExpiry
	if expiryStr == "" {
		expiryStr = "168h" // 7 days default
	}
	expiryDuration, err := time.ParseDuration(expiryStr)
	if err != nil {
		expiryDuration = 168 * time.Hour // 7 days
	}

	expiresAt := time.Now().Add(expiryDuration)
	expiresIn := int64(expiryDuration.Seconds())

	roles := s.buildUserRoles(user)
	primary := s.primaryRole(user)

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    primary,
		"roles":   roles,
		"type":    "refresh",
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWT.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresIn, nil
}

// validateRefreshToken parses and validates a refresh token; returns user ID
func (s *AuthService) validateRefreshToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return 0, errors.New("invalid or expired refresh token")
	}

	if !token.Valid {
		return 0, errors.New("refresh token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	tokenType, _ := claims["type"].(string)
	if tokenType != "refresh" {
		return 0, errors.New("invalid token type")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user ID in token")
	}

	return uint(userID), nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *AuthService) ValidateToken(tokenString string) (uint, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return 0, "", errors.New("invalid token")
	}

	if !token.Valid {
		return 0, "", errors.New("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", errors.New("invalid user ID in token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", errors.New("invalid role in token")
	}

	return uint(userID), role, nil
}
