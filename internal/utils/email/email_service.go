package email

import (
	"fmt"
)

// EmailService handles email notifications
type EmailService struct {
	// SMTP config or API keys would go here
}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
	return &EmailService{}
}

// SendEmail sends an email (Simulated)
func (s *EmailService) SendEmail(to, subject, body string) error {
	// In a real implementation, this would use SMTP or an API (SendGrid, SES, etc.)
	// For now, we simulate it by logging to stdout
	fmt.Printf("--------------------------------------------------\n")
	fmt.Printf("SIMULATED EMAIL SENT:\n")
	fmt.Printf("To: %s\n", to)
	fmt.Printf("Subject: %s\n", subject)
	fmt.Printf("Body: %s\n", body)
	fmt.Printf("--------------------------------------------------\n")
	return nil
}
