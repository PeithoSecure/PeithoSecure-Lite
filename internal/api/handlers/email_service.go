package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/peithosecure/peitho-backend/internal/db/sqlite"
)

type EmailServiceConfig struct {
	SMTPHost      string
	SMTPPort      string
	SMTPUsername  string
	SMTPPassword  string
	FromEmail     string
	FrontendURL   string
	AppLinkScheme string
}

var emailConfig EmailServiceConfig

func InitEmailService() {
	emailConfig = EmailServiceConfig{
		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      os.Getenv("SMTP_PORT"),
		SMTPUsername:  os.Getenv("SMTP_USERNAME"),
		SMTPPassword:  os.Getenv("SMTP_PASSWORD"),
		FromEmail:     os.Getenv("FROM_EMAIL"),
		FrontendURL:   strings.TrimSuffix(os.Getenv("FRONTEND_URL"), "/"),
		AppLinkScheme: strings.TrimSuffix(os.Getenv("APP_LINK_SCHEME"), "/"),
	}

	if emailConfig.SMTPHost == "" || emailConfig.SMTPPort == "" {
		log.Println("⚠️ Email service is misconfigured: missing SMTP_HOST or SMTP_PORT")
	}
	if emailConfig.FromEmail == "" {
		emailConfig.FromEmail = "no-reply@peithosecure.io"
	}
	if emailConfig.FrontendURL == "" {
		emailConfig.FrontendURL = "http://localhost:3000"
	}
}

func SendVerificationEmail(username, email string) error {
	token := generateToken(16)
	link := generateDeepLink("verify", token)

	fmt.Printf("📩 [EMAIL TOKEN] Generated for %s → %s\n", email, token)

	subject := "Verify your PeithoSecure Email"
	body := fmt.Sprintf(`
<html>
  <body style="font-family: Arial, sans-serif; color: #333;">
    <p>Hello %s,</p>
    <p>Please verify your email by clicking the button below:</p>
    <p style="margin: 20px 0;">
      <a href="%s" target="_blank" style="background-color: #2563eb; color: #fff; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block;">
        Verify Email
      </a>
    </p>
    <p>If the button doesn't work, copy and paste this link into your browser:</p>
    <p><a href="%s">%s</a></p>
    <p>Thank you,<br/>PeithoSecure Team</p>
  </body>
</html>`, username, link, link, link)

	if err := sqlite.InsertEmailToken(email, token, "verify"); err != nil {
		log.Printf("❌ Failed to insert email token: %v", err)
		return fmt.Errorf("token_insert_fail: %w", err)
	}
	return sendEmail(email, subject, body)
}

func SendPasswordResetEmail(email, token string) error {
	link := generateDeepLink("reset", token)

	subject := "Reset Your PeithoSecure Password"
	body := fmt.Sprintf(`
<html>
  <body style="font-family: Arial, sans-serif; color: #333;">
    <p>You requested to reset your PeithoSecure password.</p>
    <p style="margin: 20px 0;">
      <a href="%s" target="_blank" style="background-color: #e11d48; color: #fff; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block;">
        Reset Password
      </a>
    </p>
    <p>If the button doesn't work, copy and paste this link into your browser:</p>
    <p><a href="%s">%s</a></p>
    <p>If you didn’t request this, you can safely ignore this email.</p>
  </body>
</html>`, link, link, link)

	return sendEmail(email, subject, body)
}

func generateDeepLink(linkType, token string) string {
	if emailConfig.AppLinkScheme != "" {
		if strings.HasPrefix(emailConfig.AppLinkScheme, "http") {
			return fmt.Sprintf("%s?type=%s&token=%s", emailConfig.AppLinkScheme, linkType, token)
		}
		return fmt.Sprintf("%s/%s?token=%s", emailConfig.AppLinkScheme, linkType, token)
	}
	return fmt.Sprintf("%s/%s?token=%s", emailConfig.FrontendURL, linkType, token)
}

func sendEmail(to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", emailConfig.SMTPUsername, emailConfig.SMTPPassword, emailConfig.SMTPHost)

	headers := map[string]string{
		"From":         emailConfig.FromEmail,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
	}

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n" + htmlBody)

	err := smtp.SendMail(
		emailConfig.SMTPHost+":"+emailConfig.SMTPPort,
		auth,
		emailConfig.FromEmail,
		[]string{to},
		[]byte(msg.String()),
	)

	if err != nil {
		log.Printf("❌ Failed to send email to %s: %v", to, err)
		return fmt.Errorf("email_send_fail: %w", err)
	}
	return nil
}

func generateToken(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
