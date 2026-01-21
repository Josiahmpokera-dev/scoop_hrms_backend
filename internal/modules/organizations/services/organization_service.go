package services

import (
	"errors"
	"fmt"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organizations/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/organizations/repositories"
)

type OrganizationService struct {
	repo *repositories.OrganizationRepository
}

func NewOrganizationService() *OrganizationService {
	return &OrganizationService{
		repo: repositories.NewOrganizationRepository(),
	}
}

// CreateOrganization creates a new organization
func (s *OrganizationService) CreateOrganization(req *models.CreateOrganizationRequest, tenantID *uint, updatedBy *uint) (*models.Organization, error) {
	// Check if code already exists
	if s.repo.ExistsByCode(req.Code) {
		return nil, errors.New("organization with this code already exists")
	}

	organization := &models.Organization{
		TenantID:          tenantID,
		Code:              req.Code,
		Name:              req.Name,
		LegalName:         req.LegalName,
		RegistrationNumber: req.RegistrationNumber,
		TaxID:             req.TaxID,
		Industry:          req.Industry,
		CompanySize:       req.CompanySize,
		Description:       req.Description,
		Website:           req.Website,
		Domain:            req.Domain,
		LogoURL:           req.LogoURL,
		AddressLine1:      req.AddressLine1,
		AddressLine2:      req.AddressLine2,
		City:              req.City,
		State:             req.State,
		Country:           req.Country,
		PostalCode:        req.PostalCode,
		PhoneNumber:       req.PhoneNumber,
		Email:             req.Email,
		FoundedDate:       req.FoundedDate,
		FiscalYearStart:   req.FiscalYearStart,
		Timezone:          req.Timezone,
		CurrencyCode:      req.CurrencyCode,
		Status:            models.OrganizationStatusActive,
		IsActive:          true,
		UpdatedBy:         updatedBy,
	}
	
	// Set defaults if not provided
	if organization.FiscalYearStart == nil {
		defaultFiscalYear := 1
		organization.FiscalYearStart = &defaultFiscalYear
	}
	if organization.Timezone == nil {
		defaultTimezone := "Africa/Dar_es_Salaam"
		organization.Timezone = &defaultTimezone
	}
	if organization.CurrencyCode == nil {
		defaultCurrency := "TZS"
		organization.CurrencyCode = &defaultCurrency
	}

	if req.Status != nil {
		organization.Status = models.OrganizationStatus(*req.Status)
	}
	if req.IsActive != nil {
		organization.IsActive = *req.IsActive
	}

	if err := s.repo.Create(organization); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	return organization, nil
}

// GetOrganizationByID retrieves an organization by ID
func (s *OrganizationService) GetOrganizationByID(id uint) (*models.Organization, error) {
	return s.repo.FindByID(id)
}

// UpdateOrganization updates an existing organization
func (s *OrganizationService) UpdateOrganization(id uint, req *models.UpdateOrganizationRequest, updatedBy *uint) (*models.Organization, error) {
	organization, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("organization not found")
	}

	// Check code uniqueness if code is being updated
	if req.Code != nil && *req.Code != organization.Code {
		if s.repo.ExistsByCode(*req.Code) {
			return nil, errors.New("organization with this code already exists")
		}
		organization.Code = *req.Code
	}

	if req.Name != nil {
		organization.Name = *req.Name
	}
	if req.LegalName != nil {
		organization.LegalName = req.LegalName
	}
	if req.RegistrationNumber != nil {
		organization.RegistrationNumber = req.RegistrationNumber
	}
	if req.TaxID != nil {
		organization.TaxID = req.TaxID
	}
	if req.Industry != nil {
		organization.Industry = req.Industry
	}
	if req.CompanySize != nil {
		organization.CompanySize = req.CompanySize
	}
	if req.Description != nil {
		organization.Description = req.Description
	}
	if req.Website != nil {
		organization.Website = req.Website
	}
	if req.Domain != nil {
		organization.Domain = req.Domain
	}
	if req.LogoURL != nil {
		organization.LogoURL = req.LogoURL
	}
	if req.AddressLine1 != nil {
		organization.AddressLine1 = req.AddressLine1
	}
	if req.AddressLine2 != nil {
		organization.AddressLine2 = req.AddressLine2
	}
	if req.City != nil {
		organization.City = req.City
	}
	if req.State != nil {
		organization.State = req.State
	}
	if req.Country != nil {
		organization.Country = req.Country
	}
	if req.PostalCode != nil {
		organization.PostalCode = req.PostalCode
	}
	if req.PhoneNumber != nil {
		organization.PhoneNumber = req.PhoneNumber
	}
	if req.Email != nil {
		organization.Email = req.Email
	}
	if req.FoundedDate != nil {
		organization.FoundedDate = req.FoundedDate
	}
	if req.FiscalYearStart != nil {
		organization.FiscalYearStart = req.FiscalYearStart
	}
	if req.Timezone != nil {
		organization.Timezone = req.Timezone
	}
	if req.CurrencyCode != nil {
		organization.CurrencyCode = req.CurrencyCode
	}
	if req.Status != nil {
		organization.Status = models.OrganizationStatus(*req.Status)
	}
	if req.IsActive != nil {
		organization.IsActive = *req.IsActive
	}

	organization.UpdatedBy = updatedBy

	if err := s.repo.Update(organization); err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	return organization, nil
}

// DeleteOrganization soft deletes an organization
func (s *OrganizationService) DeleteOrganization(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("organization not found")
	}

	return s.repo.Delete(id)
}

// ListOrganizations lists organizations with pagination and filters
func (s *OrganizationService) ListOrganizations(tenantID *uint, page, pageSize int, filters map[string]interface{}) ([]models.Organization, int64, error) {
	return s.repo.List(tenantID, page, pageSize, filters)
}
