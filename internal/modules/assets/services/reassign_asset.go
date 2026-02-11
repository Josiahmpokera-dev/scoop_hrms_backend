package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/assets/models"
	departmentRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/departments/repositories"
)

// ReassignAsset reassigns an asset from one employee to another
func (s *AssetService) ReassignAsset(req *models.ReassignAssetRequest, tenantID *uint, updatedBy *uint) (*models.Asset, error) {
	asset, err := s.repo.FindByID(req.AssetID)
	if err != nil {
		return nil, errors.New("asset not found")
	}


	// Check if asset is currently assigned
	if !asset.IsAssigned() {
		return nil, errors.New("asset is not currently assigned to any employee")
	}

	// Check if reassigning to the same employee
	if asset.EmployeeID != nil && *asset.EmployeeID == req.NewEmployeeID {
		return nil, errors.New("asset is already assigned to this employee")
	}

	// Find new employee
	newEmployee, err := s.employeeRepo.FindByEmployeeID(req.NewEmployeeID)
	if err != nil {
		return nil, errors.New("new employee not found")
	}

	if !newEmployee.IsActive {
		return nil, errors.New("new employee is not active")
	}

	// Parse dates
	reassignDate := time.Now()
	if req.ReassignDate != nil && *req.ReassignDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.ReassignDate)
		if err != nil {
			return nil, errors.New("invalid reassign_date format, expected YYYY-MM-DD")
		}
		reassignDate = parsed
	}

	// Store previous assignment info for history
	previousEmployeeID := asset.EmployeeID

	// Update condition if provided
	if req.Condition != nil && *req.Condition != "" {
		validConditions := []string{"excellent", "good", "fair", "poor"}
		isValidCondition := false
		for _, c := range validConditions {
			if strings.ToLower(*req.Condition) == c {
				asset.Condition = strings.ToLower(*req.Condition)
				isValidCondition = true
				break
			}
		}
		if !isValidCondition {
			return nil, errors.New("invalid condition, must be one of: excellent, good, fair, poor")
		}
	}

	// Reassign to new employee
	newFullName := fmt.Sprintf("%s %s", newEmployee.FirstName, newEmployee.LastName)
	asset.AssignedTo = &newFullName
	asset.EmployeeID = &newEmployee.EmployeeID
	
	// Get department name for new employee
	if newEmployee.DepartmentID != nil {
		deptRepo := departmentRepos.NewDepartmentRepository()
		dept, err := deptRepo.FindByID(*newEmployee.DepartmentID)
		if err == nil && dept != nil {
			asset.Department = &dept.Name
		} else {
			// Fallback to department ID as string
			deptIDStr := fmt.Sprintf("%d", *newEmployee.DepartmentID)
			asset.Department = &deptIDStr
		}
	} else {
		asset.Department = nil
	}
	
	asset.AssignedDate = &reassignDate
	asset.ReturnDate = nil // Clear return date on new assignment
	asset.Status = string(models.AssetStatusInUse)
	
	// Update notes
	if req.Notes != nil {
		existingNotes := ""
		if asset.Notes != nil {
			existingNotes = *asset.Notes + "\n"
		}
		previousEmpID := "Unknown"
		if previousEmployeeID != nil {
			previousEmpID = *previousEmployeeID
		}
		newNote := fmt.Sprintf("Reassigned from %s to %s on %s", 
			previousEmpID,
			req.NewEmployeeID,
			reassignDate.Format("2006-01-02"))
		if *req.Notes != "" {
			newNote += fmt.Sprintf(" - %s", *req.Notes)
		}
		combinedNotes := existingNotes + newNote
		asset.Notes = &combinedNotes
	}
	
	asset.UpdatedBy = updatedBy
	asset.UpdatedAt = time.Now()

	if err := s.repo.Update(asset); err != nil {
		return nil, fmt.Errorf("failed to reassign asset: %w", err)
	}

	return asset, nil
}
