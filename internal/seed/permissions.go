package seed

import (
	"log"

	roleRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/repositories"
)

// RunPermissions seeds default permissions for the system.
func RunPermissions() {
	permRepo := roleRepos.NewPermissionRepository()
	if err := permRepo.SeedDefaultPermissions(); err != nil {
		log.Printf("Warning: Failed to seed default permissions: %v", err)
		return
	}
	log.Println("✅ Default permissions seeded successfully")
}
