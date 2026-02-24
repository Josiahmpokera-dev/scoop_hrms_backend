package models

import (
	"time"

	"gorm.io/gorm"
)

type JobOpening struct {
	ID               string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	RequisitionID    string         `json:"requisitionId" gorm:"not null"`
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
	ExpiryDate       time.Time      `json:"expiryDate"`
	Status           string         `json:"status" gorm:"default:'Draft'"` // Draft, Published, Closed
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
	Education   []string `json:"education"`
	SalaryRange struct {
		Min      float64 `json:"min"`
		Max      float64 `json:"max"`
		Currency string  `json:"currency"`
	} `json:"salaryRange"`
	ExpiryDate string `json:"expiryDate"` // Parse to time.Time
}

type PublishJobOpeningRequest struct {
	Platforms  []string `json:"platforms" binding:"required"`
	ExpiryDate string   `json:"expiryDate"`
}
