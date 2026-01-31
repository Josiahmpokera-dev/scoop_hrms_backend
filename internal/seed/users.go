package seed

import (
	"log"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	userRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories"
	"golang.org/x/crypto/bcrypt"
)

// RunAdminUser creates the initial admin user if it doesn't exist.
func RunAdminUser() {
	userRepo := userRepos.NewUserRepository()

	username := "admin"
	email := "admin@hrms.com"
	password := "admin123"
	firstName := "Admin"
	lastName := "User"

	if userRepo.ExistsByEmail(email) {
		log.Println("Admin user already exists, skipping seed")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Warning: Failed to hash password for admin user: %v", err)
		return
	}

	admin := &models.User{
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
		Role:      models.RoleAdmin,
		IsActive:  true,
	}

	if err := userRepo.Create(admin); err != nil {
		log.Printf("Warning: Failed to create admin user: %v", err)
		return
	}

	log.Printf("✅ Admin user created successfully: %s (%s)", username, email)
}

// NonEmployeeUserSeed defines a non-employee user to seed (for "onboard existing user" flow).
type NonEmployeeUserSeed struct {
	Username  string
	Email     string
	Password  string
	FirstName string
	LastName  string
	Role      models.UserRole
}

// RunNonEmployeeUsers creates sample users who are NOT yet employees (IT and HR samples).
// They appear in GET /api/v1/employees/onboarding/non-employee-users and can be onboarded as employees.
// Credentials are documented in docs/SEED_NON_EMPLOYEE_USERS.md.
func RunNonEmployeeUsers() {
	userRepo := userRepos.NewUserRepository()

	users := []NonEmployeeUserSeed{
		{
			Username:  "ituser",
			Email:     "it.user@hrms.com",
			Password:  "ItUser123!",
			FirstName: "IT",
			LastName:  "User",
			Role:      models.RoleIT,
		},
		{
			Username:  "hruser",
			Email:     "hr.user@hrms.com",
			Password:  "HrUser123!",
			FirstName: "HR",
			LastName:  "User",
			Role:      models.RoleHR,
		},
	}

	for _, u := range users {
		existing, _ := userRepo.FindByEmail(u.Email)
		if existing != nil {
			// If IT user exists with old role "user", update to "it" so login returns role "it"
			if u.Email == "it.user@hrms.com" && existing.Role == models.RoleUser {
				existing.Role = models.RoleIT
				if err := userRepo.Update(existing); err != nil {
					log.Printf("Warning: Failed to update IT user role: %v", err)
				} else {
					log.Printf("✅ IT user role updated to 'it': %s", u.Email)
				}
			} else {
				log.Printf("Non-employee user already exists, skipping: %s", u.Email)
			}
			continue
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Warning: Failed to hash password for %s: %v", u.Email, err)
			continue
		}

		user := &models.User{
			Username:  u.Username,
			Email:     u.Email,
			Password:  string(hashedPassword),
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      u.Role,
			Status:    "active",
			IsActive:  true,
		}

		if err := userRepo.Create(user); err != nil {
			log.Printf("Warning: Failed to create non-employee user %s: %v", u.Email, err)
			continue
		}

		log.Printf("✅ Non-employee user created: %s (%s) - role %s", u.Username, u.Email, string(u.Role))
	}

	log.Println("✅ Non-employee users seed completed (see docs/SEED_NON_EMPLOYEE_USERS.md for credentials)")
}
