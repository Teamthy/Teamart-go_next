package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/teamart/commerce-api/pkg/logger"
)

// EmailProvider defines different email service providers
type EmailProvider string

const (
	ProviderSMTP     EmailProvider = "smtp"
	ProviderSES      EmailProvider = "ses"
	ProviderSendGrid EmailProvider = "sendgrid"
)

// EmailConfig represents email service configuration
type EmailConfig struct {
	Provider    EmailProvider
	FromAddress string
	FromName    string

	// SMTP Configuration
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string

	// AWS SES Configuration
	SESRegion    string
	SESAccessKey string
	SESSecretKey string

	// SendGrid Configuration
	SendGridAPIKey string

	// Email templates
	OTPTemplate           string
	PasswordResetTemplate string
	LoginAlertTemplate    string
	OnboardingTemplate    string
	KYCTemplate           string
}

// EmailService handles email operations
type EmailService struct {
	config *EmailConfig
	logger *logger.Logger
}

// NewEmailService creates a new email service
func NewEmailService(config *EmailConfig, logger *logger.Logger) *EmailService {
	return &EmailService{
		config: config,
		logger: logger,
	}
}

// SendEmailInput represents input for sending an email
type SendEmailInput struct {
	To         []string
	CC         []string
	BCC        []string
	Subject    string
	Body       string
	HTMLBody   string
	IsHTML     bool
	ReplyTo    string
	References map[string]interface{} // For tracking/audit
}

// SendEmailOutput represents the result of sending an email
type SendEmailOutput struct {
	MessageID string
	Status    string
}

// SendEmail sends an email using the configured provider
func (es *EmailService) SendEmail(ctx context.Context, input *SendEmailInput) (*SendEmailOutput, error) {
	if len(input.To) == 0 {
		return nil, fmt.Errorf("at least one recipient is required")
	}

	if input.Subject == "" {
		return nil, fmt.Errorf("subject is required")
	}

	if input.Body == "" && input.HTMLBody == "" {
		return nil, fmt.Errorf("email body is required")
	}

	switch es.config.Provider {
	case ProviderSMTP:
		return es.sendViaSMTP(ctx, input)
	case ProviderSES:
		return es.sendViaSES(ctx, input)
	case ProviderSendGrid:
		return es.sendViaSendGrid(ctx, input)
	default:
		return nil, fmt.Errorf("unsupported email provider: %s", es.config.Provider)
	}
}

// sendViaSMTP sends email via SMTP
func (es *EmailService) sendViaSMTP(ctx context.Context, input *SendEmailInput) (*SendEmailOutput, error) {
	// Build SMTP auth
	auth := smtp.PlainAuth("", es.config.SMTPUsername, es.config.SMTPPassword, es.config.SMTPHost)

	// Build email headers
	headers := buildEmailHeaders(es.config, input)

	// Combine headers and body
	var body string
	if input.IsHTML && input.HTMLBody != "" {
		body = headers + "\r\n" + input.HTMLBody
	} else {
		body = headers + "\r\n" + input.Body
	}

	// Send email
	addr := fmt.Sprintf("%s:%d", es.config.SMTPHost, es.config.SMTPPort)
	err := smtp.SendMail(addr, auth, es.config.FromAddress, input.To, []byte(body))
	if err != nil {
		es.logger.Errorf("failed to send email via SMTP: %v", err)
		return nil, err
	}

	messageID := generateMessageID()
	es.logger.Infof("email sent via SMTP to %v (message_id: %s)", input.To, messageID)

	return &SendEmailOutput{
		MessageID: messageID,
		Status:    "sent",
	}, nil
}

// sendViaSES sends email via AWS SES (stub for future implementation)
func (es *EmailService) sendViaSES(ctx context.Context, input *SendEmailInput) (*SendEmailOutput, error) {
	// TODO: Implement AWS SES integration
	// - Initialize SES client
	// - Send email
	// - Return message ID
	es.logger.Warnf("AWS SES provider not yet implemented, falling back to SMTP")
	return es.sendViaSMTP(ctx, input)
}

// sendViaSendGrid sends email via SendGrid (stub for future implementation)
func (es *EmailService) sendViaSendGrid(ctx context.Context, input *SendEmailInput) (*SendEmailOutput, error) {
	// TODO: Implement SendGrid integration
	// - Initialize SendGrid client
	// - Send email
	// - Return message ID
	es.logger.Warnf("SendGrid provider not yet implemented, falling back to SMTP")
	return es.sendViaSMTP(ctx, input)
}

// SendOTPEmail sends an OTP email
func (es *EmailService) SendOTPEmail(ctx context.Context, email string, otp string) error {
	subject := "Your Teamart Verification Code"
	body := fmt.Sprintf(`Your Teamart OTP is: %s

This code expires in 10 minutes.
If you didn't request this, please ignore this email.`, otp)

	_, err := es.SendEmail(ctx, &SendEmailInput{
		To:      []string{email},
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	})

	return err
}

// SendPasswordResetEmail sends a password reset email
func (es *EmailService) SendPasswordResetEmail(ctx context.Context, email string, resetToken string, resetLink string) error {
	subject := "Reset Your Teamart Password"
	body := fmt.Sprintf(`Click the link below to reset your password:

%s

This link expires in 24 hours.
If you didn't request this, please ignore this email.`, resetLink)

	_, err := es.SendEmail(ctx, &SendEmailInput{
		To:      []string{email},
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	})

	return err
}

// SendSuspiciousLoginAlert sends a suspicious login alert email
func (es *EmailService) SendSuspiciousLoginAlert(ctx context.Context, email string, location string, ipAddress string) error {
	subject := "Suspicious Login Activity on Your Teamart Account"
	body := fmt.Sprintf(`We detected a login attempt from an unusual location:

Location: %s
IP Address: %s
Time: %s

If this wasn't you, please change your password immediately.
Security is our priority.`, location, ipAddress, "now")

	_, err := es.SendEmail(ctx, &SendEmailInput{
		To:      []string{email},
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	})

	return err
}

// SendOnboardingEmail sends an onboarding email
func (es *EmailService) SendOnboardingEmail(ctx context.Context, email string, userName string) error {
	subject := "Welcome to Teamart!"
	body := fmt.Sprintf(`Welcome to Teamart, %s!

You're all set to start creating content and building your audience.

Next steps:
1. Complete your profile
2. Set up your store
3. Create your first stream

Let's get started: https://teamart.app/onboarding

Happy streaming!`, userName)

	_, err := es.SendEmail(ctx, &SendEmailInput{
		To:      []string{email},
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	})

	return err
}

// SendKYCVerificationEmail sends a KYC verification email
func (es *EmailService) SendKYCVerificationEmail(ctx context.Context, email string, verificationLink string) error {
	subject := "Complete Your KYC Verification"
	body := fmt.Sprintf(`To enable payouts, you need to complete your KYC verification.

Click the link below to start:

%s

This process usually takes less than 10 minutes.
Your account will be reviewed within 24 hours.`, verificationLink)

	_, err := es.SendEmail(ctx, &SendEmailInput{
		To:      []string{email},
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	})

	return err
}

// ===== Helper Functions =====

// buildEmailHeaders builds email headers
func buildEmailHeaders(config *EmailConfig, input *SendEmailInput) string {
	var headers []string

	// From header
	from := fmt.Sprintf("%s <%s>", config.FromName, config.FromAddress)
	headers = append(headers, fmt.Sprintf("From: %s", from))

	// To header
	headers = append(headers, fmt.Sprintf("To: %s", strings.Join(input.To, ", ")))

	// CC header
	if len(input.CC) > 0 {
		headers = append(headers, fmt.Sprintf("Cc: %s", strings.Join(input.CC, ", ")))
	}

	// BCC header (not typically included in message)
	if len(input.BCC) > 0 {
		headers = append(headers, fmt.Sprintf("Bcc: %s", strings.Join(input.BCC, ", ")))
	}

	// Subject
	headers = append(headers, fmt.Sprintf("Subject: %s", input.Subject))

	// Reply-To
	if input.ReplyTo != "" {
		headers = append(headers, fmt.Sprintf("Reply-To: %s", input.ReplyTo))
	}

	// Content-Type
	if input.IsHTML {
		headers = append(headers, "Content-Type: text/html; charset=UTF-8")
	} else {
		headers = append(headers, "Content-Type: text/plain; charset=UTF-8")
	}

	// MIME version
	headers = append(headers, "MIME-Version: 1.0")

	return strings.Join(headers, "\r\n")
}

// generateMessageID generates a unique message ID
func generateMessageID() string {
	// In production, this should be a more robust implementation
	return fmt.Sprintf("teamart-%d@teamart.app", time.Now().UnixNano())
}
