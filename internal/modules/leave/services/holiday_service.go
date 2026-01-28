package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/repositories"
)

type HolidayService struct {
	repo *repositories.HolidayRepository
}

func NewHolidayService() *HolidayService {
	return &HolidayService{
		repo: repositories.NewHolidayRepository(),
	}
}

// CreateHoliday creates a new holiday
func (s *HolidayService) CreateHoliday(req *models.CreateHolidayRequest, tenantID *uint, createdBy *uint) (*models.Holiday, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format. Use YYYY-MM-DD")
	}

	// Serialize location and applicable_for arrays to JSON
	locationJSON, err := json.Marshal(req.Location)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize location: %w", err)
	}

	applicableForJSON, err := json.Marshal(req.ApplicableFor)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize applicable_for: %w", err)
	}

	holiday := &models.Holiday{
		TenantID:      tenantID,
		Name:          req.Name,
		Date:          date,
		Type:          req.Type,
		IsFloater:     false,
		Location:      string(locationJSON),
		ApplicableFor: string(applicableForJSON),
		Description:   req.Description,
		CreatedBy:     createdBy,
		UpdatedBy:     createdBy,
	}

	if req.IsFloater != nil {
		holiday.IsFloater = *req.IsFloater
	}

	if err := s.repo.Create(holiday); err != nil {
		return nil, fmt.Errorf("failed to create holiday: %w", err)
	}

	return s.repo.FindByID(holiday.ID)
}

// GetHolidayByID retrieves a holiday by ID
func (s *HolidayService) GetHolidayByID(id uint) (*models.Holiday, error) {
	return s.repo.FindByID(id)
}

// UpdateHoliday updates a holiday
func (s *HolidayService) UpdateHoliday(id uint, req *models.UpdateHolidayRequest, updatedBy *uint) (*models.Holiday, error) {
	holiday, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("holiday not found")
	}

	if req.Name != nil {
		holiday.Name = *req.Name
	}
	if req.Date != nil {
		date, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return nil, errors.New("invalid date format. Use YYYY-MM-DD")
		}
		holiday.Date = date
	}
	if req.Type != nil {
		holiday.Type = *req.Type
	}
	if req.IsFloater != nil {
		holiday.IsFloater = *req.IsFloater
	}
	if req.Location != nil {
		locationJSON, err := json.Marshal(req.Location)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize location: %w", err)
		}
		holiday.Location = string(locationJSON)
	}
	if req.ApplicableFor != nil {
		applicableForJSON, err := json.Marshal(req.ApplicableFor)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize applicable_for: %w", err)
		}
		holiday.ApplicableFor = string(applicableForJSON)
	}
	if req.Description != nil {
		holiday.Description = req.Description
	}

	holiday.UpdatedBy = updatedBy

	if err := s.repo.Update(holiday); err != nil {
		return nil, fmt.Errorf("failed to update holiday: %w", err)
	}

	return s.repo.FindByID(id)
}

// DeleteHoliday deletes a holiday
func (s *HolidayService) DeleteHoliday(id uint) error {
	return s.repo.Delete(id)
}

// ListHolidays lists holidays with pagination and filters
func (s *HolidayService) ListHolidays(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Holiday, int64, error) {
	return s.repo.List(tenantID, page, pageSize, filters)
}

// GetHolidaysByDateRange gets holidays within a date range
func (s *HolidayService) GetHolidaysByDateRange(startDate, endDate time.Time, tenantID *uint) ([]models.Holiday, error) {
	return s.repo.FindByDateRange(startDate, endDate, tenantID)
}

// GetHolidaysByDate gets holidays on a specific date
func (s *HolidayService) GetHolidaysByDate(date time.Time, tenantID *uint) ([]models.Holiday, error) {
	return s.repo.FindByDate(date, tenantID)
}

// IsHoliday checks if a date is a holiday
func (s *HolidayService) IsHoliday(date time.Time, tenantID *uint) (bool, error) {
	holidays, err := s.repo.FindByDate(date, tenantID)
	if err != nil {
		return false, err
	}
	return len(holidays) > 0, nil
}
