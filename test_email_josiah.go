package main

import (
	"fmt"
	"os"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/email"
)

func main() {
	// Load configuration (this loads .env file)
	_, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	fmt.Println("🔧 Testing Email Delivery to josiahgeofreyfx@gmail.com")
	fmt.Println("====================================================")

	// Display SMTP configuration
	fmt.Printf("SMTP_HOST: %s\n", os.Getenv("SMTP_HOST"))
	fmt.Printf("SMTP_PORT: %s\n", os.Getenv("SMTP_PORT"))
	fmt.Printf("SMTP_USER: %s\n", os.Getenv("SMTP_USER"))
	fmt.Printf("SMTP_FROM_EMAIL: %s\n", os.Getenv("SMTP_FROM_EMAIL"))
	fmt.Println("")

	// Test email service
	emailService := email.NewEmailService()

	// Test email with professional content
	subject := "✅ HRMS Email Test - ScoopWorks System Verification"

	body := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Email Test - ScoopWorks HRMS</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; text-align: center; border-radius: 5px; }
        .content { background-color: #fff; padding: 30px; border-radius: 5px; margin-top: 20px; border: 1px solid #e9ecef; }
        .success { color: #28a745; font-weight: bold; }
        .info { background-color: #e9ecef; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .footer { margin-top: 30px; text-align: center; color: #6c757d; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✅ Email Test Successful</h1>
        </div>
        
        <div class="content">
            <h2>Hello Josiah,</h2>
            
            <p class="success">This email confirms that your ScoopWorks HR Management System email service is working correctly!</p>
            
            <div class="info">
                <h3>📧 Email Configuration Details:</h3>
                <p><strong>SMTP Server:</strong> smtp.gmail.com</p>
                <p><strong>Port:</strong> 587 (STARTTLS)</p>
                <p><strong>From Address:</strong> system@scoopworks.com</p>
                <p><strong>Test Recipient:</strong> josiahgeofreyfx@gmail.com</p>
            </div>
            
            <p>The system will automatically send welcome emails with login credentials to new employees when they complete their onboarding process.</p>
            
            <p><strong>Next Steps:</strong></p>
            <ul>
                <li>Employee onboarding will trigger automatic welcome emails</li>
                <li>Credentials will be sent securely via email</li>
                <li>Professional HTML templates are configured</li>
                <li>Fallback to simulation mode if SMTP fails</li>
            </ul>
            
            <p>Best regards,<br>
            <strong>ScoopWorks HR Team</strong></p>
        </div>
        
        <div class="footer">
            <p>This is an automated test message. Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
`

	fmt.Printf("📤 Sending test email to: josiahgeofreyfx@gmail.com\n")
	fmt.Printf("📝 Subject: %s\n", subject)
	fmt.Println("")

	// Send email
	err = emailService.SendEmail(
		"josiahgeofreyfx@gmail.com",
		subject,
		body,
	)

	if err != nil {
		fmt.Printf("❌ Error sending email: %v\n", err)
		fmt.Println("")
		fmt.Println("🔍 Troubleshooting Tips:")
		fmt.Println("1. Check SMTP credentials in .env file")
		fmt.Println("2. Verify Gmail app password is correct")
		fmt.Println("3. Ensure less secure apps is enabled (if using Gmail)")
		fmt.Println("4. Check firewall/network restrictions")
	} else {
		fmt.Println("✅ Email sent successfully!")
		fmt.Println("")
		fmt.Println("📨 Please check your inbox at josiahgeofreyfx@gmail.com")
		fmt.Println("   (Also check spam/junk folder if not in inbox)")
	}

	fmt.Println("")
	fmt.Println("====================================================")
}
