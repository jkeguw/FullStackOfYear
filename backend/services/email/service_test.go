// services/email/service_test.go
package email

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

// mockGoMailDialer 模拟 gomail.Dialer
type mockGoMailDialer struct {
	lastMessage *gomail.Message
	sendCount   int
	shouldFail  bool
}

func (m *mockGoMailDialer) DialAndSend(messages ...*gomail.Message) error {
	m.sendCount++
	if m.shouldFail {
		return &RetryableError{Err: assert.AnError}
	}
	if len(messages) > 0 {
		m.lastMessage = messages[0]
	}
	return nil
}

func TestEmailService(t *testing.T) {
	// 创建测试配置
	config := &Config{
		SMTP: struct {
			Host     string
			Port     int
			Username string
			Password string
		}{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "test@example.com",
			Password: "test_password",
		},
		From:    "noreply@example.com",
		BaseURL: "http://localhost:3000",
	}

	// 创建测试用的模板
	tmpl, err := template.New("verifyEmail").Parse(`
        <!DOCTYPE html>
        <html>
        <body>
            <h1>验证邮箱</h1>
            <p>你好 {{.Username}},</p>
            <p>请点击以下链接验证你的邮箱：</p>
            <a href="{{.VerifyLink}}">验证邮箱</a>
            <p>链接将在 {{.ExpiresIn}} 小时后过期</p>
            <p>如有问题请联系: {{.SupportMail}}</p>
        </body>
        </html>
    `)
	require.NoError(t, err)

	// 创建测试用的 templates map
	templates := map[string]*template.Template{
		"verifyEmail": tmpl,
	}

	// 创建 logger
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	// 创建mock dialer
	mockDialer := &mockGoMailDialer{}

	service := &Service{
		config:    config,
		dialer:    mockDialer,
		templates: templates,
		logger:    logger,
	}

	t.Run("测试发送验证邮件", func(t *testing.T) {
		to := "user@example.com"
		username := "testuser"
		token := "verify_token_123"

		err := service.SendVerificationEmail(to, username, token)
		assert.NoError(t, err)

		// 验证发送的邮件内容
		require.NotNil(t, mockDialer.lastMessage)
		assert.Equal(t, []string{to}, mockDialer.lastMessage.GetHeader("To"))
		assert.Equal(t, []string{config.From}, mockDialer.lastMessage.GetHeader("From"))
		assert.Equal(t, []string{"Verify Your Email Address"}, mockDialer.lastMessage.GetHeader("Subject"))

		// 获取邮件内容
		var body bytes.Buffer
		mockDialer.lastMessage.WriteTo(&body)
		bodyStr := body.String()

		// 验证邮件内容包含必要信息
		assert.Contains(t, bodyStr, username)
		assert.Contains(t, bodyStr, token)
		assert.Contains(t, bodyStr, config.BaseURL)
	})

	t.Run("测试邮件模板渲染失败", func(t *testing.T) {
		// 创建一个无效的模板
		invalidTmpl, err := template.New("invalid").Parse("{{.InvalidField}}")
		require.NoError(t, err)

		invalidService := &Service{
			config: config,
			dialer: mockDialer,
			templates: map[string]*template.Template{
				"verifyEmail": invalidTmpl,
			},
			logger: logger,
		}

		err = invalidService.SendVerificationEmail("test@example.com", "user", "token")
		assert.Error(t, err)
	})

	t.Run("测试发送重试机制", func(t *testing.T) {
		// 设置mock dialer为失败模式
		failingDialer := &mockGoMailDialer{shouldFail: true}
		failingService := &Service{
			config:    config,
			dialer:    failingDialer,
			templates: templates,
			logger:    logger,
		}

		err := failingService.SendVerificationEmail("test@example.com", "user", "token")
		assert.Error(t, err)
		assert.GreaterOrEqual(t, failingDialer.sendCount, 3) // 至少重试3次
	})
}

// ... TestEmailData_Rendering 保持不变 ...

func TestEmailService_TestConnection(t *testing.T) {
	t.Run("测试成功连接", func(t *testing.T) {
		mockDialer := &mockGoMailDialer{}
		service := &Service{
			config: &Config{
				From: "test@example.com",
			},
			dialer: mockDialer,
		}

		err := service.TestConnection()
		assert.NoError(t, err)
		assert.Equal(t, 1, mockDialer.sendCount)
	})

	t.Run("测试连接失败", func(t *testing.T) {
		mockDialer := &mockGoMailDialer{shouldFail: true}
		service := &Service{
			config: &Config{
				From: "test@example.com",
			},
			dialer: mockDialer,
		}

		err := service.TestConnection()
		assert.Error(t, err)
	})
}
