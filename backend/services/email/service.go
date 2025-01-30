package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
	"time"
)

// Service handles email operations
type Service struct {
	config    *Config
	dialer    *gomail.Dialer
	templates map[string]*template.Template
}

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

// EmailData represents the data needed for email templates
type EmailData struct {
	Username    string
	VerifyLink  string
	ExpiresIn   time.Duration
	SupportMail string
}

// NewEmailService creates a new email service instance
func NewEmailService(config *Config) (*Service, error) {
	dialer := gomail.NewDialer(
		config.SMTP.Host,
		config.SMTP.Port,
		config.SMTP.Username,
		config.SMTP.Password,
	)

	// Enable SSL/TLS
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}

	// Initialize templates
	templates := make(map[string]*template.Template)
	for name, path := range config.Templates {
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return nil, err
		}
		templates[name] = tmpl
	}

	return &Service{
		config:    config,
		dialer:    dialer,
		templates: templates,
	}, nil
}

// SendVerificationEmail sends email verification link
func (s *Service) SendVerificationEmail(to, username, token string) error {
	data := EmailData{
		Username:    username,
		VerifyLink:  s.buildVerificationLink(token),
		ExpiresIn:   24 * time.Hour, // Token expiration time
		SupportMail: s.config.From,
	}

	var body bytes.Buffer
	if err := s.templates["verifyEmail"].Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Verify Your Email Address")
	m.SetBody("text/html", body.String())

	return s.dialer.DialAndSend(m)
}

// buildVerificationLink generates the verification URL
func (s *Service) buildVerificationLink(token string) string {
	return fmt.Sprintf("%s/verify-email?token=%s", s.config.BaseURL, token)
}

// SendTestEmail sends a test email to verify SMTP configuration
func (s *Service) SendTestEmail(to string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Test Email")
	m.SetBody("text/plain", "This is a test email from the system.")

	return s.dialer.DialAndSend(m)
}
