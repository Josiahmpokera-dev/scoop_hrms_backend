package handlers

import (
	employeeModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/employees/models"
)

// Helper functions shared across handlers

// getEmployeeFullName gets the full name from an employee object
func getEmployeeFullName(employee interface{}) string {
	if employee == nil {
		return ""
	}
	if emp, ok := employee.(*employeeModels.Employee); ok {
		return emp.FullName()
	}
	return ""
}

// getEmployeePhoto gets the photo URL from an employee object
func getEmployeePhoto(employee interface{}) string {
	if employee == nil {
		return ""
	}
	if emp, ok := employee.(*employeeModels.Employee); ok && emp.PhotoURL != nil {
		return *emp.PhotoURL
	}
	return ""
}

// getEmployeeDepartment gets the department name from an employee object
// Note: This is a placeholder - department name would need to be fetched from department repository
func getEmployeeDepartment(employee interface{}) string {
	if employee == nil {
		return ""
	}
	// Department name would need to be fetched from department repository
	// For now, return empty string
	return ""
}
