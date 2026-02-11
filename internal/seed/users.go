package seed

import (
	"log"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	userRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories"
	"golang.org/x/crypto/bcrypt"
)

// defaultPassword is the shared password for all seeded users.
// Change this in production or after first login.
const defaultPassword = "Admin@2024!"

// RunAdminUser creates the initial Admin user if it doesn't exist.
// On every run it also resets the password and unblocks the account so the
// credentials always work during development.
func RunAdminUser() {
	userRepo := userRepos.NewUserRepository()

	email := "admin@hrms.com"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Warning: Failed to hash password for admin user: %v", err)
		return
	}

	existing, _ := userRepo.FindByEmail(email)
	if existing != nil {
		// Reset password, unblock, and ensure active on every restart
		existing.Password = string(hashedPassword)
		existing.Status = "active"
		existing.IsActive = true
		existing.FailedLoginCount = 0
		existing.LockedUntil = nil
		if err := userRepo.Update(existing); err != nil {
			log.Printf("Warning: Failed to reset admin user: %v", err)
		} else {
			log.Printf("✅ Admin user reset: %s (password: %s)", email, defaultPassword)
		}
		return
	}

	admin := &models.User{
		Username:  "admin",
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: "Admin",
		LastName:  "User",
		Role:      models.RoleAdmin,
		Status:    "active",
		IsActive:  true,
	}

	if err := userRepo.Create(admin); err != nil {
		log.Printf("Warning: Failed to create admin user: %v", err)
		return
	}

	log.Printf("✅ Admin user created: %s (password: %s)", email, defaultPassword)
}

// RunHRUser creates the initial HR user if it doesn't exist.
// Same password reset behaviour as RunAdminUser.
func RunHRUser() {
	userRepo := userRepos.NewUserRepository()

	email := "hr@hrms.com"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Warning: Failed to hash password for HR user: %v", err)
		return
	}

	existing, _ := userRepo.FindByEmail(email)
	if existing != nil {
		existing.Password = string(hashedPassword)
		existing.Status = "active"
		existing.IsActive = true
		existing.FailedLoginCount = 0
		existing.LockedUntil = nil
		if err := userRepo.Update(existing); err != nil {
			log.Printf("Warning: Failed to reset HR user: %v", err)
		} else {
			log.Printf("✅ HR user reset: %s (password: %s)", email, defaultPassword)
		}
		return
	}

	hrUser := &models.User{
		Username:  "hrmanager",
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: "HR",
		LastName:  "Manager",
		Role:      models.RoleHR,
		Status:    "active",
		IsActive:  true,
	}

	if err := userRepo.Create(hrUser); err != nil {
		log.Printf("Warning: Failed to create HR user: %v", err)
		return
	}

	log.Printf("✅ HR user created: %s (password: %s)", email, defaultPassword)
}
