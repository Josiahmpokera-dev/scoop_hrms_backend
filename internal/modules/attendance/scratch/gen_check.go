//go:build ignore
// +build ignore

package main

import (
	"bufio"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("internal/modules/attendance/scratch/models.txt")
	if err != nil { panic(err) }
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var models []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "&") {
			line = strings.TrimSuffix(line, ",")
			models = append(models, line)
		}
	}

	out, _ := os.Create("internal/modules/attendance/scratch/check_schema.go")
	defer out.Close()

	out.WriteString(`package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	roleModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/roles/models"
	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	organizationModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization/models"
	organizationUnitModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organization_units/models"
	departmentModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/models"
	teamModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/teams/models"
	positionModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/models"
	locationModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/models"
	costCenterModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/cost_centers/models"
	assetModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/models"
	employeeModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
	helpdeskModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
	biometricModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/models"
	shiftModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/shifts/models"
	leaveModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	auditModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/audit/models"
	payrollModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/models"
	dashboardModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/dashboard/models"
	attendanceModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
	projectModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/projects/models"
	performanceModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/performance/models"
	settingsModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/settings/models"
	recruitmentModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/recruitment/models"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=scoop password=my_scoop_@_2026 dbname=scoop_db_v001 port=5432 sslmode=disable TimeZone=UTC"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	modelsToTest := []interface{}{
`)

	for _, m := range models {
		out.WriteString("\t\t" + m + ",\n")
	}

	out.WriteString(`	}

	for i, model := range modelsToTest {
		log.Printf("Migrating model index %d (%T)...", i, model)
		err = db.AutoMigrate(model)
		if err != nil {
			log.Fatalf("Migration failed for %T: %v", model, err)
		}
	}
	log.Println("All models migrated successfully")
}
`)
}
