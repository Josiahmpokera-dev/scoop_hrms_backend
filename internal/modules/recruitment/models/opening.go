package models

import (
	"time"

	"gorm.io/gorm"
)

// Job opening lifecycle statuses (use Active for published/live roles).
const (
	JobOpeningStatusDraft  = "Draft"
	JobOpeningStatusActive = "Active"
	JobOpeningStatusClosed = "Closed"
	// LegacyStatusPublished is accepted when reading/filtering for backward compatibility with older rows.
	LegacyStatusPublished = "Published"
)

type JobOpening struct {
	ID               string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	RequisitionID    string         `json:"requisitionId" gorm:"not null;index"`
	Requisition      *JobRequisition `json:"requisition,omitempty" gorm:"foreignKey:RequisitionID;references:ID"`
	ApplyToken       string         `json:"applyToken" gorm:"size:32;uniqueIndex"` // public apply link token
	JobTitle         string         `json:"jobTitle"`
	JobDescription   string         `json:"jobDescription"`
	Responsibilities []string       `json:"responsibilities" gorm:"serializer:json"`
	MustHaveSkills   []string       `json:"mustHaveSkills" gorm:"serializer:json"`
	NiceToHaveSkills []string       `json:"niceToHaveSkills" gorm:"serializer:json"`
	ExperienceMin    int            `json:"experienceMin"`
	ExperienceMax    int            `json:"experienceMax"`
	Education        []string       `json:"education" gorm:"serializer:json"`
	SalaryMin        float64        `json:"salaryMin"`
	SalaryMax        float64        `json:"salaryMax"`
	SalaryCurrency   string         `json:"salaryCurrency" gorm:"default:'TZS'"`
	ApplicantCount   int64          `json:"applicantCount" gorm:"-:all"`
	ExpiryDate       time.Time      `json:"expiryDate"`
	Status           string         `json:"status" gorm:"default:'Draft';index"` // Draft, Active, Closed (Published legacy)
	Platforms        []string       `json:"platforms" gorm:"serializer:json"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateJobOpeningRequest struct {
	RequisitionID    string   `json:"requisitionId" binding:"required"`
	JobTitle         string   `json:"jobTitle" binding:"required"`
	JobDescription   string   `json:"jobDescription" binding:"required"`
	Responsibilities []string `json:"responsibilities"`
	MustHaveSkills   []string `json:"mustHaveSkills"`
	NiceToHaveSkills []string `json:"niceToHaveSkills"`
	Experience       struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"experience"`
	Education  []string `json:"education"`
	ExpiryDate string `json:"expiryDate"` // YYYY-MM-DD or RFC3339; optional while Draft
}

type UpdateJobOpeningRequest struct {
	JobTitle         *string  `json:"jobTitle"`
	JobDescription   *string  `json:"jobDescription"`
	Responsibilities []string `json:"responsibilities"`
	MustHaveSkills   []string `json:"mustHaveSkills"`
	NiceToHaveSkills []string `json:"niceToHaveSkills"`
	ExperienceMin    *int     `json:"experienceMin"`
	ExperienceMax    *int     `json:"experienceMax"`
	Education        []string `json:"education"`
	SalaryMin        *float64 `json:"salaryMin"`
	SalaryMax        *float64 `json:"salaryMax"`
	SalaryCurrency   *string  `json:"salaryCurrency"`
	ExpiryDate       *string  `json:"expiryDate"` // YYYY-MM-DD or RFC3339
	Platforms        []string `json:"platforms"`
}

type PublishJobOpeningRequest struct {
	Platforms  []string `json:"platforms"`  // optional; channels where the role is advertised
	ExpiryDate string   `json:"expiryDate"` // required to go live unless already set on opening
}

type JobOpeningShareLinkResponse struct {
	OpeningID       string `json:"opening_id"`
	ApplyToken      string `json:"apply_token"`
	CareersApplyURL string `json:"careers_apply_url,omitempty"` // when PUBLIC_RECRUITMENT_CAREERS_URL is set
	PublicJobAPIURL string `json:"public_job_api_url"`          // GET JSON for this opening (token)
	PublicApplyURL  string `json:"public_apply_url"`            // hint: POST /recruitment/public/apply same host
	Instructions    string `json:"instructions"`
}
