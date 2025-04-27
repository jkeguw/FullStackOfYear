package email

import (
	"fmt"
	"log"
)

// Service 提供电子邮件功能
type Service struct {
}

// NewService 创建一个新的邮件服务
func NewService(config interface{}) *Service {
	return &Service{}
}

// SendEmail 发送一封电子邮件
func (s *Service) SendEmail(to, subject, body string) error {
	log.Printf("发送邮件到 %s，主题：%s", to, subject)
	// 在实际实现中，这里会使用SMTP发送邮件
	return nil
}

// SendVerificationEmail 发送邮箱验证邮件
func (s *Service) SendVerificationEmail(to, username, token string) error {
	subject := "请验证您的电子邮件地址"
	verificationURL := fmt.Sprintf("http://localhost:8080/api/auth/verify-email?token=%s", token)
	
	// 构建邮件内容
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
	resetURL := fmt.Sprintf("http://localhost:8080/reset-password?token=%s", token)
	
	// 构建邮件内容
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
	return nil
}