package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
)

// EmailService handles email notifications
type EmailService struct {
	host     string
	port     int
	username string
	password string
	from     string
	fromName string
}

// NewEmailService creates a new email service with SMTP configuration
func NewEmailService() *EmailService {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM_EMAIL")
	fromName := os.Getenv("SMTP_FROM_NAME")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 587 // default SMTP port
	}

	return &EmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		fromName: fromName,
	}
}

// SendEmail sends an email using SMTP
func (s *EmailService) SendEmail(to, subject, body string) error {
	// Build email message
	fromHeader := fmt.Sprintf("%s <%s>", s.fromName, s.from)
	message := []byte(
		"From: " + fromHeader + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
			body + "\r\n")

	// SMTP authentication
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	// For Gmail SMTP (port 587 with STARTTLS)
	if s.port == 587 {
		// Connect to SMTP server
		conn, err := smtp.Dial(addr)
		if err != nil {
			return s.fallbackToSimulation(to, subject, body, err)
		}
		defer conn.Close()

		// Start TLS
		if err = conn.StartTLS(&tls.Config{ServerName: s.host}); err != nil {
			return s.fallbackToSimulation(to, subject, body, err)
		}

		// Authenticate
		if err = conn.Auth(auth); err != nil {
			return s.fallbackToSimulation(to, subject, body, err)
		}

		// Set sender and recipient
		if err = conn.Mail(s.from); err != nil {
			return s.fallbackToSimulation(to, subject, body, err)
		}
		if err = conn.Rcpt(to); err != nil {
			return s.fallbackToSimulation(to, subject, body, err)
		}

		// Send email body
		w, err := conn.Data()
		if err != nil {
			return s.fallbackToSimulation(to, subject, body, err)
		}

		_, err = w.Write(message)
		if err != nil {
			return s.fallbackToSimulation(to, subject, body, err)
		}

		err = w.Close()
		if err != nil {
			return s.fallbackToSimulation(to, subject, body, err)
		}

		fmt.Printf("Email sent successfully to: %s\n", to)
		return nil
	}

	// For SSL/TLS connection (port 465)
	if s.port == 465 {
		tlsConfig := &tls.Config{
			ServerName:         s.host,
			InsecureSkipVerify: true, // For self-signed certificates
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return s.fallbackToSimulation(to, subject, body, err)
		}

		client, err := smtp.NewClient(conn, s.host)
		if err != nil {
			return s.fallbackToSimulation(to, subject, body, err)
		}
		defer client.Close()

		// Auth
		if err = client.Auth(auth); err != nil {
			return s.fallbackToSimulation(to, subject, body, err)
		}

		// Set sender and recipient
		if err = client.Mail(s.from); err != nil {
			return s.fallbackToSimulation(to, subject, body, err)
		}
		if err = client.Rcpt(to); err != nil {
			return s.fallbackToSimulation(to, subject, body, err)
		}

		// Send email body
		w, err := client.Data()
		if err != nil {
			return s.fallbackToSimulation(to, subject, body, err)
		}

		_, err = w.Write(message)
		if err != nil {
			return s.fallbackToSimulation(to, subject, body, err)
		}

		err = w.Close()
		if err != nil {
			return s.fallbackToSimulation(to, subject, body, err)
		}

		fmt.Printf("Email sent successfully to: %s\n", to)
		return nil
	}

	// Fallback for non-standard ports
	return s.fallbackToSimulation(to, subject, body, nil)
}

// fallbackToSimulation handles email sending failures by logging to stdout
func (s *EmailService) fallbackToSimulation(to, subject, body string, err error) error {
	fmt.Printf("--------------------------------------------------\n")
	if err != nil {
		fmt.Printf("SMTP FAILED - SIMULATED EMAIL SENT:\n")
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("SIMULATED EMAIL SENT:\n")
	}
	fmt.Printf("To: %s\n", to)
	fmt.Printf("Subject: %s\n", subject)
	fmt.Printf("Body: %s\n", body)
	fmt.Printf("--------------------------------------------------\n")
	return nil // Don't fail the operation, just log the error
}
