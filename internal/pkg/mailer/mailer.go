package mailer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime"
	"net/smtp"
	"path/filepath"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
)

// Attachment represents an email attachment
type Attachment struct {
	Filename string
	Content  []byte
}

// Mailer handles sending emails via SMTP
type Mailer struct {
	config config.SMTPConfig
}

// NewMailer creates a new mailer
func NewMailer() *Mailer {
	return &Mailer{
		config: config.AppConfig.SMTP,
	}
}

// SendEmail sends an email with optional attachments
func (m *Mailer) SendEmail(to []string, subject, body string, attachments []Attachment) error {
	if m.config.Host == "" || m.config.Username == "" {
		return fmt.Errorf("SMTP configuration is incomplete")
	}

	delimiter := "boundary_delimiter_1234567890"
	
	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", m.config.FromName, m.config.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to[0])) // Simplification: Send to first recipient in header
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", delimiter))
	msg.WriteString("\r\n")

	// Body
	msg.WriteString(fmt.Sprintf("--%s\r\n", delimiter))
	msg.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)
	msg.WriteString("\r\n")

	// Attachments
	for _, attr := range attachments {
		msg.WriteString(fmt.Sprintf("--%s\r\n", delimiter))
		msg.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", 
			mime.TypeByExtension(filepath.Ext(attr.Filename)), attr.Filename))
		msg.WriteString("Content-Transfer-Encoding: base64\r\n")
		msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", attr.Filename))
		msg.WriteString("\r\n")
		
		b := make([]byte, base64.StdEncoding.EncodedLen(len(attr.Content)))
		base64.StdEncoding.Encode(b, attr.Content)
		msg.Write(b)
		msg.WriteString("\r\n")
	}

	msg.WriteString(fmt.Sprintf("--%s--\r\n", delimiter))

	auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)
	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)

	return smtp.SendMail(addr, auth, m.config.From, to, msg.Bytes())
}
