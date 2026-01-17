package models

// CreateTeamRequest represents the request to create a team
type CreateTeamRequest struct {
	DepartmentID *uint   `json:"department_id,omitempty"`
	Code         string  `json:"code" binding:"required,min=2,max=20"`
	Name         string  `json:"name" binding:"required,min=2,max=100"`
	Description  *string `json:"description,omitempty"`
	TeamLeadID   *uint   `json:"team_lead_id,omitempty"`
	TeamType     *string `json:"team_type,omitempty" binding:"omitempty,oneof=project functional cross-functional virtual"`
	ProjectID    *uint   `json:"project_id,omitempty"`
	MaxMembers   *int    `json:"max_members,omitempty"`
	LocationID   *uint   `json:"location_id,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

// UpdateTeamRequest represents the request to update a team
type UpdateTeamRequest struct {
	DepartmentID *uint   `json:"department_id,omitempty"`
	Code         *string `json:"code,omitempty" binding:"omitempty,min=2,max=20"`
	Name         *string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Description  *string `json:"description,omitempty"`
	TeamLeadID   *uint   `json:"team_lead_id,omitempty"`
	TeamType     *string `json:"team_type,omitempty" binding:"omitempty,oneof=project functional cross-functional virtual"`
	ProjectID    *uint   `json:"project_id,omitempty"`
	MaxMembers   *int    `json:"max_members,omitempty"`
	LocationID   *uint   `json:"location_id,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}
