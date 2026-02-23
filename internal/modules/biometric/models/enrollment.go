package models

import (
	"time"

	"gorm.io/gorm"
)

// BiometricEnrollment links an HRMS employee record to a biometric device user.
type BiometricEnrollment struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	EmployeeID     uint           `json:"employee_id" gorm:"uniqueIndex;not null;index"`
	EmpCode        string         `json:"emp_code" gorm:"uniqueIndex;not null;size:100;index"`
	DeviceUserName string         `json:"device_user_name" gorm:"size:255"`
	Status         string         `json:"status" gorm:"type:varchar(20);default:'linked';index"` // linked, unlinked
	LinkedAt       time.Time      `json:"linked_at"`
	LinkedBy       *uint          `json:"linked_by,omitempty" gorm:"index"`
	Notes          *string        `json:"notes,omitempty" gorm:"type:text"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (BiometricEnrollment) TableName() string { return "biometric_enrollments" }

// ---------- Request DTOs ----------

// LinkEmployeeRequest links a single employee to a biometric emp_code.
type LinkEmployeeRequest struct {
	EmployeeID uint   `json:"employee_id" binding:"required"`
	EmpCode    string `json:"emp_code" binding:"required"`
	Notes      string `json:"notes,omitempty"`
}

// BulkLinkRequest links multiple employees at once.
type BulkLinkRequest struct {
	Mappings []LinkEmployeeRequest `json:"mappings" binding:"required,min=1,dive"`
}

// UnlinkRequest — no body needed, employee_id from URL.

// ---------- Response DTOs ----------

// EnrollmentDetail is a rich response combining enrollment + employee + device data.
type EnrollmentDetail struct {
	ID             uint    `json:"id"`
	EmployeeID     uint    `json:"employee_id"`
	EmployeeCode   string  `json:"employee_code"`
	EmployeeName   string  `json:"employee_name"`
	Department     *string `json:"department,omitempty"`
	Position       *string `json:"position,omitempty"`
	EmpCode        string  `json:"emp_code"`
	DeviceUserName string  `json:"device_user_name"`
	Status         string  `json:"status"`
	LinkedAt       string  `json:"linked_at"`
}

// DeviceUser represents a distinct user from the biometric device data.
type DeviceUser struct {
	EmpCode    string `json:"emp_code"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Department string `json:"department"`
	Position   string `json:"position"`
	Linked     bool   `json:"linked"`
	EmployeeID *uint  `json:"employee_id,omitempty"`
}

// UnlinkedEmployee is an employee not yet linked to any biometric device user.
type UnlinkedEmployee struct {
	ID           uint    `json:"id"`
	EmployeeCode string  `json:"employee_code"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Department   *string `json:"department,omitempty"`
	Position     *string `json:"position,omitempty"`
	SuggestedEmpCode *string `json:"suggested_emp_code,omitempty"`
}

// MergedAttendance combines biometric punch data with HRMS employee profile.
type MergedAttendance struct {
	EmployeeID     uint    `json:"employee_id"`
	EmployeeCode   string  `json:"employee_code"`
	EmployeeName   string  `json:"employee_name"`
	Department     *string `json:"department,omitempty"`
	Position       *string `json:"position,omitempty"`
	PhotoURL       *string `json:"photo_url,omitempty"`
	EmpCode        string  `json:"emp_code"`
	Date           string  `json:"date"`
	CheckIn        *string `json:"check_in,omitempty"`
	CheckOut       *string `json:"check_out,omitempty"`
	PunchCount     int     `json:"punch_count"`
	WorkingHours   *string `json:"working_hours,omitempty"`
	Status         string  `json:"status"` // Present, Absent, Late, Half Day
}

// AutoLinkResult reports outcomes of the auto-linking operation.
type AutoLinkResult struct {
	TotalEmployees int                `json:"total_employees"`
	TotalDeviceUsers int              `json:"total_device_users"`
	Linked         int                `json:"linked"`
	AlreadyLinked  int                `json:"already_linked"`
	NoMatch        int                `json:"no_match"`
	Details        []AutoLinkDetail   `json:"details,omitempty"`
}

type AutoLinkDetail struct {
	EmployeeID   uint   `json:"employee_id"`
	EmployeeCode string `json:"employee_code"`
	EmployeeName string `json:"employee_name"`
	EmpCode      string `json:"emp_code"`
	Action       string `json:"action"` // linked, already_linked, no_match
}
