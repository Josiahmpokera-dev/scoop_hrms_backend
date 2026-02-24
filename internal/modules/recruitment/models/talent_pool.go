package models

import (
	"time"

	"gorm.io/gorm"
)

type TalentPoolCandidate struct {
	ID            string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name          string         `json:"name" gorm:"not null"`
	Email         string         `json:"email" gorm:"unique;not null"`
	Skills        []string       `json:"skills" gorm:"serializer:json"`
	PreferredRole []string       `json:"preferredRole" gorm:"serializer:json"`
	ResumeURL     string         `json:"resumeUrl"`
	ExperienceMin int            `json:"experienceMin"` // Years
	Location      string         `json:"location"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type AddToTalentPoolRequest struct {
	Name          string   `json:"name" binding:"required"`
	Email         string   `json:"email" binding:"required,email"`
	Skills        []string `json:"skills"`
	PreferredRole []string `json:"preferredRole"`
	ResumeURL     string   `json:"resumeUrl"`
}
