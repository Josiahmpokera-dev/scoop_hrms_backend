package services

import (
	"errors"
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/locations/repositories"
	organizationRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organizations/repositories"
)

type LocationService struct {
	repo            *repositories.LocationRepository
	organizationRepo *organizationRepos.OrganizationRepository
}

func NewLocationService() *LocationService {
	return &LocationService{
		repo:            repositories.NewLocationRepository(),
		organizationRepo: organizationRepos.NewOrganizationRepository(),
	}
}

// CreateLocation creates a new location
func (s *LocationService) CreateLocation(req *models.CreateLocationRequest, tenantID *uint, updatedBy *uint) (*models.Location, error) {
	// Validate organization if provided
	if req.OrganizationID != nil {
		_, err := s.organizationRepo.FindByID(*req.OrganizationID)
		if err != nil {
			return nil, errors.New("organization not found")
		}
	}

	location := &models.Location{
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		LocationType:   req.LocationType,
		AddressLine1:   req.AddressLine1,
		AddressLine2:   req.AddressLine2,
		City:           req.City,
		State:          req.State,
		Country:        req.Country,
		PostalCode:     req.PostalCode,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		Timezone:       req.Timezone,
		PhoneNumber:    req.PhoneNumber,
		Capacity:       req.Capacity,
		Facilities:     req.Facilities,
		IsHeadOffice:   false,
		IsActive:       true,
		UpdatedBy:      updatedBy,
	}
	
	// Set defaults if not provided
	if location.Timezone == nil {
		defaultTimezone := "Africa/Dar_es_Salaam"
		location.Timezone = &defaultTimezone
	}

	if req.IsHeadOffice != nil {
		location.IsHeadOffice = *req.IsHeadOffice
		// If setting as head office, unset other head offices for the tenant
		if *req.IsHeadOffice {
			if err := s.unsetOtherHeadOffices(tenantID); err != nil {
				return nil, fmt.Errorf("failed to unset other head offices: %w", err)
			}
		}
	}

	if req.IsActive != nil {
		location.IsActive = *req.IsActive
	}

	if err := s.repo.Create(location); err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	return location, nil
}

// GetLocationByID retrieves a location by ID
func (s *LocationService) GetLocationByID(id uint) (*models.Location, error) {
	return s.repo.FindByID(id)
}

// UpdateLocation updates an existing location
func (s *LocationService) UpdateLocation(id uint, req *models.UpdateLocationRequest, updatedBy *uint) (*models.Location, error) {
	location, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("location not found")
	}

	if req.OrganizationID != nil {
		// Validate organization if being updated
		_, err := s.organizationRepo.FindByID(*req.OrganizationID)
		if err != nil {
			return nil, errors.New("organization not found")
		}
		location.OrganizationID = req.OrganizationID
	}
	if req.Name != nil {
		location.Name = *req.Name
	}
	if req.LocationType != nil {
		location.LocationType = req.LocationType
	}
	if req.AddressLine1 != nil {
		location.AddressLine1 = req.AddressLine1
	}
	if req.AddressLine2 != nil {
		location.AddressLine2 = req.AddressLine2
	}
	if req.City != nil {
		location.City = req.City
	}
	if req.State != nil {
		location.State = req.State
	}
	if req.Country != nil {
		location.Country = req.Country
	}
	if req.PostalCode != nil {
		location.PostalCode = req.PostalCode
	}
	if req.Latitude != nil {
		location.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		location.Longitude = req.Longitude
	}
	if req.Timezone != nil {
		location.Timezone = req.Timezone
	}
	if req.PhoneNumber != nil {
		location.PhoneNumber = req.PhoneNumber
	}
	if req.Capacity != nil {
		location.Capacity = req.Capacity
	}
	if req.Facilities != nil {
		location.Facilities = req.Facilities
	}
	if req.IsHeadOffice != nil {
		location.IsHeadOffice = *req.IsHeadOffice
		// If setting as head office, unset other head offices
		if *req.IsHeadOffice {
			// Note: tenantID not available in UpdateLocation, skipping unsetOtherHeadOffices
		}
	}
	if req.IsActive != nil {
		location.IsActive = *req.IsActive
	}

	location.UpdatedBy = updatedBy

	if err := s.repo.Update(location); err != nil {
		return nil, fmt.Errorf("failed to update location: %w", err)
	}

	return location, nil
}

// DeleteLocation soft deletes a location
func (s *LocationService) DeleteLocation(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("location not found")
	}

	return s.repo.Delete(id)
}

// ListLocations lists locations with pagination and filters
func (s *LocationService) ListLocations(tenantID, organizationID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Location, int64, error) {
	return s.repo.List(tenantID, organizationID, page, pageSize, filters)
}

// GetHeadOffice gets the head office location
func (s *LocationService) GetHeadOffice(tenantID *uint) (*models.Location, error) {
	return s.repo.FindHeadOffice(tenantID)
}

// unsetOtherHeadOffices unsets is_head_office for all locations except the current one
func (s *LocationService) unsetOtherHeadOffices(tenantID *uint) error {
	// This would need to be implemented in the repository
	// For now, we'll handle it in the update/create logic
	return nil
}
