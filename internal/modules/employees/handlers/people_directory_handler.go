package handlers

import (
	"strconv"
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/middleware"
	departmentRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
	employeeRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/services"
	locationRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/repositories"
	positionRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/positions/repositories"
	teamRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/teams/repositories"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// PeopleDirectoryHandler handles People Directory HTTP requests
type PeopleDirectoryHandler struct {
	employeeRepo   *employeeRepos.EmployeeRepository
	departmentRepo *departmentRepos.DepartmentRepository
	positionRepo   *positionRepos.JobPositionRepository
	locationRepo   *locationRepos.LocationRepository
	teamRepo       *teamRepos.TeamRepository
	selfService    *services.SelfServiceService
}

// NewPeopleDirectoryHandler creates a new People Directory handler
func NewPeopleDirectoryHandler() *PeopleDirectoryHandler {
	return &PeopleDirectoryHandler{
		employeeRepo:   employeeRepos.NewEmployeeRepository(),
		departmentRepo: departmentRepos.NewDepartmentRepository(),
		positionRepo:   positionRepos.NewJobPositionRepository(),
		locationRepo:   locationRepos.NewLocationRepository(),
		teamRepo:       teamRepos.NewTeamRepository(),
		selfService:    services.NewSelfServiceService(),
	}
}

// SearchDirectory handles searching/filtering employees in the People Directory
// @Summary Search People Directory
// @Description Search and filter employees in the organization directory with comprehensive filters
// @Tags People Directory
// @Produce json
// @Param search query string false "Search term (name, employee ID, email, phone)"
// @Param department_id query int false "Filter by department ID"
// @Param position_id query int false "Filter by position/designation ID"
// @Param location_id query int false "Filter by location ID"
// @Param team_id query int false "Filter by team ID"
// @Param manager_id query int false "Filter by reporting manager ID"
// @Param employment_type query string false "Filter by employment type (full_time, part_time, contract, intern)"
// @Param status query string false "Filter by status (active, inactive, on_leave, suspended, terminated)"
// @Param letter query string false "Filter by first letter of first name (A-Z)"
// @Param sort_by query string false "Sort field (first_name, last_name, employee_id, department, hire_date)"
// @Param sort_order query string false "Sort direction (asc, desc)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(20)
// @Success 200 {object} response.APIResponse
// @Router /api/v1/people-directory [get]
func (h *PeopleDirectoryHandler) SearchDirectory(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	// Build filters from query parameters
	filters := employeeRepos.PeopleDirectoryFilters{
		TenantID:  tenantID,
		Page:      1,
		PageSize:  20,
		SortBy:    "first_name",
		SortOrder: "asc",
	}

	// Search
	if search := c.Query("search"); search != "" {
		filters.Search = &search
	}

	// Department
	if dID := c.Query("department_id"); dID != "" {
		if parsed, err := strconv.ParseUint(dID, 10, 32); err == nil {
			id := uint(parsed)
			filters.DepartmentID = &id
		}
	}

	// Position
	if pID := c.Query("position_id"); pID != "" {
		if parsed, err := strconv.ParseUint(pID, 10, 32); err == nil {
			id := uint(parsed)
			filters.PositionID = &id
		}
	}

	// Location
	if lID := c.Query("location_id"); lID != "" {
		if parsed, err := strconv.ParseUint(lID, 10, 32); err == nil {
			id := uint(parsed)
			filters.LocationID = &id
		}
	}

	// Team
	if tID := c.Query("team_id"); tID != "" {
		if parsed, err := strconv.ParseUint(tID, 10, 32); err == nil {
			id := uint(parsed)
			filters.TeamID = &id
		}
	}

	// Manager
	if mID := c.Query("manager_id"); mID != "" {
		if parsed, err := strconv.ParseUint(mID, 10, 32); err == nil {
			id := uint(parsed)
			filters.ManagerID = &id
		}
	}

	// Employment type
	if et := c.Query("employment_type"); et != "" {
		filters.EmploymentType = &et
	}

	// Status
	if status := c.Query("status"); status != "" {
		filters.Status = &status
	}

	// Letter filter
	if letter := c.Query("letter"); letter != "" {
		upperLetter := strings.ToUpper(letter[:1])
		filters.Letter = &upperLetter
	}

	// Sort
	if sortBy := c.Query("sort_by"); sortBy != "" {
		filters.SortBy = sortBy
	}
	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		filters.SortOrder = strings.ToLower(sortOrder)
	}

	// Pagination
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			filters.Page = parsed
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			filters.PageSize = parsed
		}
	}

	// Execute search
	employees, total, err := h.employeeRepo.SearchPeopleDirectory(filters)
	if err != nil {
		response.InternalServerError(c, "Failed to search directory", err.Error())
		return
	}

	// Enrich results with related data
	results := make([]gin.H, 0, len(employees))
	for _, emp := range employees {
		entry := gin.H{
			"id":          emp.ID,
			"employee_id": emp.EmployeeID,
			"first_name":  emp.FirstName,
			"last_name":   emp.LastName,
			"full_name":   emp.FullName(),
			"initials":    getInitials(emp.FirstName, emp.LastName),
			"status":      string(emp.Status),
		}

		// Photo
		if emp.PhotoURL != nil {
			entry["photo_url"] = *emp.PhotoURL
		}

		// Email (prefer work email)
		if emp.WorkEmail != nil {
			entry["email"] = *emp.WorkEmail
		} else if emp.PersonalEmail != nil {
			entry["email"] = *emp.PersonalEmail
		}

		// Phone
		if emp.PhoneNumber != nil {
			entry["phone"] = *emp.PhoneNumber
		}
		if emp.WorkPhone != nil {
			entry["work_phone"] = *emp.WorkPhone
		}

		// Employment type
		if emp.EmploymentType != nil {
			entry["employment_type"] = *emp.EmploymentType
		}

		// Department
		if emp.DepartmentID != nil {
			entry["department_id"] = *emp.DepartmentID
			dept, dErr := h.departmentRepo.FindByID(*emp.DepartmentID)
			if dErr == nil && dept != nil {
				entry["department"] = gin.H{
					"id":   dept.ID,
					"name": dept.Name,
					"code": dept.Code,
				}
			}
		}

		// Position / Designation
		if emp.PositionID != nil {
			entry["position_id"] = *emp.PositionID
			pos, pErr := h.positionRepo.FindByID(*emp.PositionID)
			if pErr == nil && pos != nil {
				entry["position"] = gin.H{
					"id":    pos.ID,
					"title": pos.Title,
					"code":  pos.Code,
				}
				entry["designation"] = pos.Title
			}
		}

		// Location
		if emp.LocationID != nil {
			entry["location_id"] = *emp.LocationID
			loc, lErr := h.locationRepo.FindByID(*emp.LocationID)
			if lErr == nil && loc != nil {
				entry["location"] = gin.H{
					"id":   loc.ID,
					"name": loc.Name,
				}
			}
		}

		// Team
		if emp.TeamID != nil {
			entry["team_id"] = *emp.TeamID
			team, tErr := h.teamRepo.FindByID(*emp.TeamID)
			if tErr == nil && team != nil {
				entry["team"] = gin.H{
					"id":   team.ID,
					"name": team.Name,
				}
			}
		}

		// Reporting Manager
		if emp.ReportsToID != nil {
			manager, mErr := h.employeeRepo.FindByID(*emp.ReportsToID)
			if mErr == nil && manager != nil {
				entry["reporting_manager"] = gin.H{
					"id":          manager.ID,
					"employee_id": manager.EmployeeID,
					"full_name":   manager.FullName(),
				}
			}
		}

		// Hire date
		if emp.HireDate != nil {
			entry["hire_date"] = emp.HireDate.Format("2006-01-02")
		}

		results = append(results, entry)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))
	meta := &response.Meta{
		Page:       filters.Page,
		PerPage:    filters.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	response.SuccessWithMeta(c, "People directory retrieved successfully", results, meta)
}

// GetEmployeeProfile returns detailed profile for a specific employee in the directory
// @Summary Get employee directory profile
// @Description Get detailed information about a specific employee
// @Tags People Directory
// @Produce json
// @Param id path int true "Employee ID (database ID)"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/people-directory/:id [get]
func (h *PeopleDirectoryHandler) GetEmployeeProfile(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	idStr := c.Param("id")

	// Try as numeric ID first, then as employee_id string
	var employeeIDStr string
	if _, err := strconv.ParseUint(idStr, 10, 32); err != nil {
		// Not a number, use as employee_id string
		employeeIDStr = idStr
	} else {
		employeeIDStr = idStr
	}

	details, err := h.selfService.GetEmployeeDirectoryDetails(employeeIDStr, tenantID)
	if err != nil {
		// Try finding by numeric ID
		numID, numErr := strconv.ParseUint(idStr, 10, 32)
		if numErr != nil {
			response.NotFound(c, "Employee not found")
			return
		}
		emp, findErr := h.employeeRepo.FindByID(uint(numID))
		if findErr != nil || emp == nil {
			response.NotFound(c, "Employee not found")
			return
		}
		// Retry with the employee_id string
		details, err = h.selfService.GetEmployeeDirectoryDetails(emp.EmployeeID, tenantID)
		if err != nil {
			response.NotFound(c, "Employee not found")
			return
		}
	}

	// Add initials
	if fn, ok := details["first_name"].(string); ok {
		if ln, ok := details["last_name"].(string); ok {
			details["initials"] = getInitials(fn, ln)
		}
	}

	response.Success(c, "Employee profile retrieved successfully", details)
}

// GetFilters returns available filter options for the People Directory
// @Summary Get directory filter options
// @Description Returns available departments, positions, locations, teams, employment types for directory filters
// @Tags People Directory
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /api/v1/people-directory/filters [get]
func (h *PeopleDirectoryHandler) GetFilters(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	// Departments
	departments, _, _ := h.departmentRepo.List(tenantID, 1, 500, map[string]interface{}{"is_active": true})
	deptOptions := make([]gin.H, 0, len(departments))
	for _, d := range departments {
		deptOptions = append(deptOptions, gin.H{
			"id":   d.ID,
			"name": d.Name,
			"code": d.Code,
		})
	}

	// Positions
	positions, _, _ := h.positionRepo.List(tenantID, 1, 500, map[string]interface{}{"is_active": true})
	posOptions := make([]gin.H, 0, len(positions))
	for _, p := range positions {
		posOptions = append(posOptions, gin.H{
			"id":    p.ID,
			"title": p.Title,
			"code":  p.Code,
		})
	}

	// Locations
	locations, _, _ := h.locationRepo.List(tenantID, nil, 1, 500, map[string]interface{}{"is_active": true})
	locOptions := make([]gin.H, 0, len(locations))
	for _, l := range locations {
		locOptions = append(locOptions, gin.H{
			"id":   l.ID,
			"name": l.Name,
		})
	}

	// Teams
	teams, _, _ := h.teamRepo.List(tenantID, nil, 1, 500, map[string]interface{}{"is_active": true})
	teamOptions := make([]gin.H, 0, len(teams))
	for _, t := range teams {
		teamOptions = append(teamOptions, gin.H{
			"id":   t.ID,
			"name": t.Name,
		})
	}

	// Employment types
	employmentTypes := []gin.H{
		{"value": "full_time", "label": "Full Time"},
		{"value": "part_time", "label": "Part Time"},
		{"value": "contract", "label": "Contract"},
		{"value": "intern", "label": "Intern"},
	}

	// Statuses
	statuses := []gin.H{
		{"value": "active", "label": "Active"},
		{"value": "inactive", "label": "Inactive"},
		{"value": "on_leave", "label": "On Leave"},
		{"value": "suspended", "label": "Suspended"},
		{"value": "terminated", "label": "Terminated"},
	}

	// Alphabet counts
	alphabetCounts, _ := h.employeeRepo.GetAlphabetCounts(tenantID)
	alphabet := make([]gin.H, 0, len(alphabetCounts))
	for _, a := range alphabetCounts {
		alphabet = append(alphabet, gin.H{
			"letter": a.Letter,
			"count":  a.Count,
		})
	}

	response.Success(c, "Directory filter options retrieved successfully", gin.H{
		"departments":      deptOptions,
		"positions":        posOptions,
		"locations":        locOptions,
		"teams":            teamOptions,
		"employment_types": employmentTypes,
		"statuses":         statuses,
		"alphabet":         alphabet,
	})
}

// GetStatistics returns summary statistics for the People Directory
// @Summary Get directory statistics
// @Description Returns employee count breakdowns by department, location, employment type
// @Tags People Directory
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /api/v1/people-directory/stats [get]
func (h *PeopleDirectoryHandler) GetStatistics(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	// Total active employees
	totalFilters := employeeRepos.PeopleDirectoryFilters{
		TenantID: tenantID,
		Page:     1,
		PageSize: 1,
	}
	_, totalActive, _ := h.employeeRepo.SearchPeopleDirectory(totalFilters)

	// Department breakdown
	deptCounts, _ := h.employeeRepo.GetDepartmentEmployeeCounts(tenantID)
	departmentBreakdown := make([]gin.H, 0, len(deptCounts))
	for _, dc := range deptCounts {
		entry := gin.H{
			"department_id": dc.DepartmentID,
			"count":         dc.Count,
		}
		dept, err := h.departmentRepo.FindByID(dc.DepartmentID)
		if err == nil && dept != nil {
			entry["department_name"] = dept.Name
			entry["department_code"] = dept.Code
		}
		departmentBreakdown = append(departmentBreakdown, entry)
	}

	// Location breakdown
	locCounts, _ := h.employeeRepo.GetLocationEmployeeCounts(tenantID)
	locationBreakdown := make([]gin.H, 0, len(locCounts))
	for _, lc := range locCounts {
		entry := gin.H{
			"location_id": lc.LocationID,
			"count":       lc.Count,
		}
		loc, err := h.locationRepo.FindByID(lc.LocationID)
		if err == nil && loc != nil {
			entry["location_name"] = loc.Name
		}
		locationBreakdown = append(locationBreakdown, entry)
	}

	// Employment type breakdown
	etCounts, _ := h.employeeRepo.GetEmploymentTypeCounts(tenantID)
	employmentTypeBreakdown := make([]gin.H, 0, len(etCounts))
	for _, ec := range etCounts {
		employmentTypeBreakdown = append(employmentTypeBreakdown, gin.H{
			"employment_type": ec.EmploymentType,
			"count":           ec.Count,
		})
	}

	response.Success(c, "Directory statistics retrieved successfully", gin.H{
		"total_active_employees":    totalActive,
		"by_department":             departmentBreakdown,
		"by_location":               locationBreakdown,
		"by_employment_type":        employmentTypeBreakdown,
	})
}

// getInitials returns the initials from first and last name
func getInitials(firstName, lastName string) string {
	initials := ""
	if len(firstName) > 0 {
		initials += strings.ToUpper(firstName[:1])
	}
	if len(lastName) > 0 {
		initials += strings.ToUpper(lastName[:1])
	}
	return initials
}
