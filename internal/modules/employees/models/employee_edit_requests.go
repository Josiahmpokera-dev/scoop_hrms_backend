package models

import "time"

type EmployeeEditRequest struct {
	Step int                    `json:"step" binding:"required,min=1,max=10"`
	Data map[string]interface{} `json:"data" binding:"required"`
}

type EmployeeEditResponse struct {
	EmployeeID    string    `json:"employee_id"`
	Step          int       `json:"step"`
	UpdatedFields []string  `json:"updated_fields"`
	Message       string    `json:"message"`
	UpdatedAt     time.Time `json:"updated_at"`
}
