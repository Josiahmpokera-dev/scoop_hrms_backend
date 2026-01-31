package seed

import (
	"log"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
)

// Run runs all seeders when the app is in development.
// Call this from main.go after migrations: if config.AppConfig.Server.Env == "development" { seed.Run() }
func Run() {
	env := config.AppConfig.Server.Env
	if env != "development" {
		log.Println("Skipping seed (not development environment)")
		return
	}

	log.Println("Running seeders...")

	RunPermissions()
	RunRoles()
	RunAdminUser()
	RunNonEmployeeUsers()

	log.Println("All seeders completed.")
}
