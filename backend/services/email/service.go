package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"project/backend/config"
)

// Service 提供电子邮件功能
type Service struct {
	config  config.EmailConfig
	baseURL string
}

// NewService 创建一个新的邮件服务
func NewService(cfg config.EmailConfig) *Service {
	return &Service{
		config:  cfg,
		baseURL: strings.TrimRight(cfg.BaseURL, "/"),
	}
}

// SendEmail 使用SMTP发送一封电子邮件
func (s *Service) SendEmail(to, subject, body string) error {
	if s.config.SMTP.Host == "" || s.config.SMTP.Username == "" || s.config.SMTP.Password == "" {
		return fmt.Errorf("email service is not configured")
	}

	from := s.config.From
	if from == "" {
		from = s.config.SMTP.Username
	}

	addr := fmt.Sprintf("%s:%d", s.config.SMTP.Host, s.config.SMTP.Port)

	headers := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n",
		from, to, subject,
	)
	msg := []byte(headers + body)

	auth := smtp.PlainAuth("", s.config.SMTP.Username, s.config.SMTP.Password, s.config.SMTP.Host)

	// 587端口使用STARTTLS
	if s.config.SMTP.Port == 587 {
		client, err := smtp.Dial(addr)
		if err != nil {
			log.Printf("Failed to connect to SMTP server: %v", err)
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		defer client.Close()

		if err := client.Hello("localhost"); err != nil {
			return err
		}

		if err := client.StartTLS(&tls.Config{ServerName: s.config.SMTP.Host}); err != nil {
			log.Printf("Failed to start TLS: %v", err)
			return fmt.Errorf("failed to start TLS: %w", err)
		}

		if err := client.Auth(auth); err != nil {
			log.Printf("SMTP authentication failed: %v", err)
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}

		if err := client.Mail(from); err != nil {
			return err
		}
		if err := client.Rcpt(to); err != nil {
			return err
		}

		w, err := client.Data()
		if err != nil {
			return err
		}
		if _, err := w.Write(msg); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}

		return client.Quit()
	}

	// 465端口使用TLS直接连接
	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// SendVerificationEmail 发送邮箱验证邮件
func (s *Service) SendVerificationEmail(to, username, token string) error {
	subject := "请验证您的电子邮件地址"
	verificationURL := fmt.Sprintf("%s/api/auth/verify-email?token=%s", s.baseURL, token)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>邮箱验证</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .button { display: inline-block; padding: 10px 20px; background-color: #4CAF50; color: white; 
                 text-decoration: none; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>请验证您的电子邮件地址</h1>
        <p>您好，%s! 感谢您注册！请点击下面的按钮验证您的电子邮件地址：</p>
        <p><a href="%s" class="button">验证邮箱</a></p>
        <p>或者，您可以复制并粘贴以下链接到您的浏览器：</p>
        <p>%s</p>
        <p>如果您没有注册账号，请忽略此邮件。</p>
    </div>
</body>
</html>
`, username, verificationURL, verificationURL)

	return s.SendEmail(to, subject, body)
}

// SendPasswordResetEmail 发送密码重置邮件
func (s *Service) SendPasswordResetEmail(to, username, token string) error {
	subject := "密码重置请求"
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.baseURL, token)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>密码重置</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .button { display: inline-block; padding: 10px 20px; background-color: #4CAF50; color: white; 
                 text-decoration: none; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>密码重置请求</h1>
        <p>%s, 我们收到了重置您密码的请求。如果这是您申请的，请点击下面的按钮重置密码：</p>
        <p><a href="%s" class="button">重置密码</a></p>
        <p>或者，您可以复制并粘贴以下链接到您的浏览器：</p>
        <p>%s</p>
        <p>如果您没有申请重置密码，请忽略此邮件。</p>
    </div>
</body>
</html>
`, username, resetURL, resetURL)

	return s.SendEmail(to, subject, body)
}

// TestConnection 测试邮件配置
func (s *Service) TestConnection() error {
	if s.config.SMTP.Host == "" {
		return fmt.Errorf("SMTP host not configured")
	}
	return nil
}
