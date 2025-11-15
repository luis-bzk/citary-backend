package services

import (
	"bytes"
	"citary-backend/internal/infrastructure/config"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"time"
)

// SMTPEmailService implements the EmailService interface using SMTP
type SMTPEmailService struct {
	config *config.Config
	auth   smtp.Auth
}

// NewSMTPEmailService creates a new SMTP email service
func NewSMTPEmailService(cfg *config.Config) *SMTPEmailService {
	auth := smtp.PlainAuth("", cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPHost)

	return &SMTPEmailService{
		config: cfg,
		auth:   auth,
	}
}

// SendVerificationEmail sends an email verification link to the user
func (s *SMTPEmailService) SendVerificationEmail(ctx context.Context, email, token string) error {
	start := time.Now()
	log.Printf("[SMTPEmailService] SendVerificationEmail: email=%s", email)

	verificationLink := fmt.Sprintf("%s/auth/verify-email?token=%s", s.config.FrontendURL, token)

	subject := "Verify Your Email Address"
	htmlBody, err := s.renderVerificationEmailTemplate(verificationLink)
	if err != nil {
		log.Printf("[SMTPEmailService] SendVerificationEmail ERROR: failed to render template, email=%s, error=%v", email, err)
		return fmt.Errorf("failed to render email template: %w", err)
	}

	err = s.sendEmail(email, subject, htmlBody)
	duration := time.Since(start)

	if err != nil {
		log.Printf("[SMTPEmailService] SendVerificationEmail ERROR: email=%s, error=%v, duration=%v", email, err, duration)
		return err
	}

	log.Printf("[SMTPEmailService] SendVerificationEmail: success, email=%s, duration=%v", email, duration)
	return nil
}

// sendEmail sends an email using SMTP
func (s *SMTPEmailService) sendEmail(to, subject, htmlBody string) error {
	from := fmt.Sprintf("%s <%s>", s.config.SMTPFromName, s.config.SMTPFromEmail)

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + htmlBody

	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	err := smtp.SendMail(addr, s.auth, s.config.SMTPFromEmail, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// renderVerificationEmailTemplate renders the HTML template for verification email
func (s *SMTPEmailService) renderVerificationEmailTemplate(verificationLink string) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email</title>
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f4;">
    <table width="100%" cellpadding="0" cellspacing="0" style="background-color: #f4f4f4; padding: 20px;">
        <tr>
            <td align="center">
                <table width="600" cellpadding="0" cellspacing="0" style="background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding: 40px 20px; text-align: center;">
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: bold;">Welcome to Citary!</h1>
                        </td>
                    </tr>

                    <!-- Body -->
                    <tr>
                        <td style="padding: 40px 30px;">
                            <h2 style="color: #333333; margin: 0 0 20px 0; font-size: 24px;">Verify Your Email Address</h2>
                            <p style="color: #666666; line-height: 1.6; margin: 0 0 20px 0; font-size: 16px;">
                                Thank you for signing up! To complete your registration and start using Citary,
                                please verify your email address by clicking the button below.
                            </p>
                            <p style="color: #666666; line-height: 1.6; margin: 0 0 30px 0; font-size: 16px;">
                                This verification link will expire in 24 hours.
                            </p>

                            <!-- Button -->
                            <table width="100%" cellpadding="0" cellspacing="0">
                                <tr>
                                    <td align="center" style="padding: 20px 0;">
                                        <a href="{{.VerificationLink}}"
                                           style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
                                                  color: #ffffff;
                                                  text-decoration: none;
                                                  padding: 15px 40px;
                                                  border-radius: 5px;
                                                  font-size: 16px;
                                                  font-weight: bold;
                                                  display: inline-block;">
                                            Verify Email Address
                                        </a>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #666666; line-height: 1.6; margin: 30px 0 0 0; font-size: 14px;">
                                If the button doesn't work, copy and paste this link into your browser:
                            </p>
                            <p style="color: #667eea; line-height: 1.6; margin: 10px 0 0 0; font-size: 14px; word-break: break-all;">
                                {{.VerificationLink}}
                            </p>
                        </td>
                    </tr>

                    <!-- Footer -->
                    <tr>
                        <td style="background-color: #f8f9fa; padding: 30px; text-align: center; border-top: 1px solid #eeeeee;">
                            <p style="color: #999999; margin: 0 0 10px 0; font-size: 14px;">
                                If you didn't create an account with Citary, you can safely ignore this email.
                            </p>
                            <p style="color: #999999; margin: 0; font-size: 12px;">
                                &copy; {{.Year}} Citary. All rights reserved.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`

	t, err := template.New("verification").Parse(tmpl)
	if err != nil {
		return "", err
	}

	data := struct {
		VerificationLink string
		Year             int
	}{
		VerificationLink: verificationLink,
		Year:             time.Now().Year(),
	}

	var buffer bytes.Buffer
	err = t.Execute(&buffer, data)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
