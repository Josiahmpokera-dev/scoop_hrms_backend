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

	// Phase 1: Core system setup (permissions, roles, admin)
	RunPermissions()
	RunRoles()
	RunAdminUser()
	RunNonEmployeeUsers()

	// Phase 2: Development test data (locations, departments, positions, teams, employees, leave)
	RunDevelopmentData()

	// Phase 3: Relationship assignments (department heads, team leads — must run after employees exist)
	RunDepartmentHeads()
	RunTeamLeads()

	log.Println("All seeders completed.")
}
