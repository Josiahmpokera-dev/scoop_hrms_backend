package seed

import (
	"log"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
)

// Run runs all seeders when the app is in development.
// Only creates permissions, roles, and two starter accounts (Admin + HR).
func Run() {
	env := config.AppConfig.Server.Env
	if env != "development" {
		log.Println("Skipping seed (not development environment)")
		return
	}

	log.Println("Running seeders...")

	// System setup
	RunPermissions()
	RunRoles()

	// Starter accounts (Admin + HR) — same password, reset on every restart
	RunAdminUser()
	RunHRUser()

	// Reference data
	RunTanzaniaHolidays()

	log.Println("────────────────────────────────────────────────────────")
	log.Printf("📋 SEED CREDENTIALS  (password: %s)", defaultPassword)
	log.Println("────────────────────────────────────────────────────────")
	log.Println("  Admin:  admin@hrms.com")
	log.Println("  HR:     hr@hrms.com")
	log.Println("────────────────────────────────────────────────────────")

	log.Println("All seeders completed.")
}
