package models

// OnboardingRequest represents the complete onboarding request (signup + organization)
type OnboardingRequest struct {
	// User signup fields
	Username  string `json:"username,omitempty" binding:"omitempty,min=3"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required,min=2"`
	LastName  string `json:"last_name" binding:"required,min=2"`
	
	// Organization fields
	Name              string  `json:"name" binding:"required,min=2,max=255"`
	LegalName         *string `json:"legal_name,omitempty" binding:"omitempty,max=255"`
	RegistrationNumber *string `json:"registration_number,omitempty" binding:"omitempty,max=100"`
	Industry          *string `json:"industry,omitempty" binding:"omitempty,max=100"`
	CompanySize       *string `json:"company_size,omitempty" binding:"omitempty,oneof=startup small medium large enterprise"`
	Country           *string `json:"country,omitempty" binding:"omitempty,max=100"`
	Timezone          *string `json:"timezone,omitempty" binding:"omitempty,max=100"`
	CurrencyCode      *string `json:"currency_code,omitempty" binding:"omitempty,len=3"`
}

// OnboardingResponse represents the complete onboarding response
type OnboardingResponse struct {
	User         *UserInfo `json:"user"`
	Tenant       *TenantInfo `json:"tenant,omitempty"`
	Organization *OrganizationInfo `json:"organization,omitempty"`
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	SetupComplete bool     `json:"setup_complete"`
}

// TenantInfo represents tenant information in response
type TenantInfo struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Domain   *string `json:"domain,omitempty"`
	Status   string `json:"status"`
}

// OrganizationInfo represents organization information in response
type OrganizationInfo struct {
	ID       uint   `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

// SetupWizardStep represents a step in the setup wizard
type SetupWizardStep struct {
	Step        string `json:"step"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	Required    bool   `json:"required"`
}

// SetupWizardStatus represents the current status of the setup wizard
type SetupWizardStatus struct {
	Completed      bool              `json:"completed"`
	CurrentStep    string            `json:"current_step"`
	Steps          []SetupWizardStep `json:"steps"`
	Progress       int               `json:"progress"` // 0-100
}
