package models

// CreateJobPositionRequest represents the request to create a job position
type CreateJobPositionRequest struct {
	Code               string  `json:"code" binding:"required,min=2,max=20"` // Position Code * (required)
	Title              string  `json:"title" binding:"required,min=2,max=100"` // Position Title * (required)
	Grade              *string `json:"grade,omitempty" binding:"omitempty,max=20"` // Grade Level * (required)
	DepartmentID       *uint   `json:"department_id,omitempty"` // Department * (required)
	ReportsToPositionID *uint  `json:"reports_to_position_id,omitempty"` // Reports To (optional)
	BudgetedHeadcount  *int    `json:"budgeted_headcount,omitempty"` // Budgeted Headcount (optional)
	CurrentHeadcount   *int    `json:"current_headcount,omitempty"` // Current Headcount (optional)
	EmploymentType     *string `json:"employment_type,omitempty" binding:"omitempty,oneof=full_time part_time contract intern"` // Employment Type (optional)
	KeyCompetencies    *string `json:"key_competencies,omitempty"` // Key Competencies - comma separated (optional)
	IsActive           *bool   `json:"is_active,omitempty"`
}

// UpdateJobPositionRequest represents the request to update a job position
type UpdateJobPositionRequest struct {
	Code               *string `json:"code,omitempty" binding:"omitempty,min=2,max=20"` // Position Code
	Title              *string `json:"title,omitempty" binding:"omitempty,min=2,max=100"` // Position Title
	Grade              *string `json:"grade,omitempty" binding:"omitempty,max=20"` // Grade Level
	DepartmentID       *uint   `json:"department_id,omitempty"` // Department
	ReportsToPositionID *uint  `json:"reports_to_position_id,omitempty"` // Reports To
	BudgetedHeadcount  *int    `json:"budgeted_headcount,omitempty"` // Budgeted Headcount
	CurrentHeadcount   *int    `json:"current_headcount,omitempty"` // Current Headcount
	EmploymentType     *string `json:"employment_type,omitempty" binding:"omitempty,oneof=full_time part_time contract intern"` // Employment Type
	KeyCompetencies    *string `json:"key_competencies,omitempty"` // Key Competencies - comma separated
	IsActive           *bool   `json:"is_active,omitempty"`
}
