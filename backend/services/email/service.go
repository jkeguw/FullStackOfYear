package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"html/template"
)

// Config represents email service configuration
type Config struct {
	SMTP struct {
		Host     string
		Port     int
		Username string
		Password string
	}
	From      string
	BaseURL   string
	Templates map[string]string
}

type Dialer interface {
	DialAndSend(m ...*gomail.Message) error
}

// Service implements email sending functionality
type Service struct {
	config    *Config
	dialer    Dialer
	templates map[string]*template.Template
	logger    *zap.Logger
}

// EmailData represents the data needed for email templates
type EmailData struct {
	Username    string
	VerifyLink  string
	ExpiresIn   int64
	SupportMail string
}

// NewEmailService creates a new email service instance
func NewEmailService(config *Config, logger *zap.Logger) (*Service, error) {
	dialer := gomail.NewDialer(
		config.SMTP.Host,
		config.SMTP.Port,
		config.SMTP.Username,
		config.SMTP.Password,
	)

	// Configure TLS
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}

	// Load templates
	templates := make(map[string]*template.Template)
	for name, path := range config.Templates {
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %v", name, err)
		}
		templates[name] = tmpl
	}

	return &Service{
		config:    config,
		dialer:    dialer,
		templates: templates,
		logger:    logger,
	}, nil
}

// SendVerificationEmail sends verification email to user
func (s *Service) SendVerificationEmail(to, username, token string) error {
	return s.withRetry(func() error {
		// Prepare template data
		data := EmailData{
			Username:    username,
			VerifyLink:  fmt.Sprintf("%s/verify-email?token=%s", s.config.BaseURL, token),
			ExpiresIn:   24,
			SupportMail: s.config.From,
		}

		// Execute template
		var body bytes.Buffer
		tmpl := s.templates["verifyEmail"]
		if tmpl == nil {
			return fmt.Errorf("verify email template not found")
		}

		if err := tmpl.Execute(&body, data); err != nil {
			return fmt.Errorf("failed to execute template: %v", err)
		}

		// Create email message
		m := gomail.NewMessage()
		m.SetHeader("From", s.config.From)
		m.SetHeader("To", to)
		m.SetHeader("Subject", "Verify Your Email Address")
		m.SetBody("text/html", body.String())

		// Send email
		if err := s.dialer.DialAndSend(m); err != nil {
			// Wrap network related errors as retryable
			if isNetworkError(err) {
				return &RetryableError{Err: err}
			}
			return err
		}

		return nil
	})
}

// isNetworkError checks if the error is network related
func isNetworkError(err error) bool {
	// Add specific error type checks based on your needs
	// For example: timeout, connection refused, etc.
	return true // For now, consider all errors retryable
}

// SendPasswordResetEmail sends password reset email to user
func (s *Service) SendPasswordResetEmail(to, username, token string) error {
	// Similar to SendVerificationEmail but with different template
	// Will implement when we add password reset functionality
	return nil
}

// TestConnection tests the email configuration by sending a test email
func (s *Service) TestConnection() error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", s.config.From) // Send to self
	m.SetHeader("Subject", "Test Email")
	m.SetBody("text/plain", "This is a test email to verify SMTP configuration.")

	return s.dialer.DialAndSend(m)
}
