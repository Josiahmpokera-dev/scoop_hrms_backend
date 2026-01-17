package models

// OrgChartNode represents a node in the organization chart
type OrgChartNode struct {
	ID            string          `json:"id"`            // Format: "emp-{employee_id}" or "pos-{position_id}" for vacant
	EmpID         *string         `json:"empId"`         // Employee ID (null for vacant positions)
	Name          string          `json:"name"`          // Employee name or "Open Position"
	Designation   string          `json:"designation"`   // Position title
	Department    *string         `json:"department"`    // Department name
	Photo         *string         `json:"photo"`         // Photo URL
	Email         *string         `json:"email"`         // Email (null for vacant)
	Phone         *string         `json:"phone"`         // Phone (null for vacant)
	ManagerEmpID  *string         `json:"managerEmpId"` // Manager's employee ID (null for root)
	Level         int             `json:"level"`         // Hierarchy level (0 = root)
	IsVacant      bool            `json:"isVacant"`      // True if position has no employee
	Subordinates  []OrgChartNode  `json:"subordinates"` // Child nodes
}
