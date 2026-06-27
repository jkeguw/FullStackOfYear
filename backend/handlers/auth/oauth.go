package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"project/backend/config"
	"project/backend/internal/errors"
	"project/backend/services/auth"
	authtypes "project/backend/types/auth"
	"time"

	"github.com/gin-gonic/gin"
)

// OAuthHandler 处理OAuth认证
type OAuthHandler struct {
	Service auth.Service
}

// NewOAuthHandler 创建OAuth处理器
func NewOAuthHandler(service auth.Service) *OAuthHandler {
	return &OAuthHandler{
		Service: service,
	}
}

// 生成随机状态字符串
func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// HandleOAuthLogin 处理OAuth登录请求
func (h *OAuthHandler) HandleOAuthLogin(c *gin.Context) {
	// 从配置中读取OAuth信息，避免使用不存在的context key导致panic
	cfg := config.GetConfig()
	oauthConfig := authtypes.GoogleOAuthConfig{
		ClientID:     cfg.OAuth.Google.ClientID,
		ClientSecret: cfg.OAuth.Google.ClientSecret,
		RedirectURL:  cfg.OAuth.Google.RedirectURL,
		Scopes:       []string{"email", "profile"},
	}

	provider := auth.NewGoogleOAuthProvider(oauthConfig)

	// 生成随机状态参数，用于防止CSRF攻击
	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalServerError, "Failed to generate state"))
		return
	}

	// 设置Cookie存储状态
	// 生产环境应启用HTTPS并将Secure设为true
	secureCookie := os.Getenv("SECURE_COOKIE") == "true" || os.Getenv("ENV") == "production"
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("oauth_state", state, int(5*time.Minute.Seconds()), "/", "", secureCookie, true)

	// 获取认证URL
	authURL := provider.GetAuthURL(state)

	// 重定向到认证URL
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// HandleOAuthCallback 处理OAuth回调
func (h *OAuthHandler) HandleOAuthCallback(c *gin.Context) {
	// 获取认证码
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, errors.NewAppError(333, "Missing code parameter"))
		return
	}

	// 获取并验证状态参数
	receivedState := c.Query("state")
	storedState, err := c.Cookie("oauth_state")
	if err != nil || receivedState != storedState {
		c.JSON(http.StatusBadRequest, errors.NewAppError(333, "Invalid state parameter"))
		return
	}

	// 从配置中读取OAuth信息
	cfg := config.GetConfig()
	oauthConfig := authtypes.GoogleOAuthConfig{
		ClientID:     cfg.OAuth.Google.ClientID,
		ClientSecret: cfg.OAuth.Google.ClientSecret,
		RedirectURL:  cfg.OAuth.Google.RedirectURL,
		Scopes:       []string{"email", "profile"},
	}

	provider := auth.NewGoogleOAuthProvider(oauthConfig)

	// 交换授权码获取token
	token, err := provider.ExchangeCode(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalServerError, fmt.Sprintf("Failed to exchange token: %v", err)))
		return
	}

	// 获取用户信息
	userInfo, err := provider.GetUserInfo(context.Background(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalServerError, fmt.Sprintf("Failed to get user info: %v", err)))
		return
	}

	// 使用获取到的用户信息进行OAuth登录或注册
	resp, err := h.Service.HandleOAuthLogin(c.Request.Context(), userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, fmt.Sprintf("登录失败: %v", err)))
		return
	}

	// 设置认证Cookie
	secureCookie := os.Getenv("SECURE_COOKIE") == "true" || os.Getenv("ENV") == "production"
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", resp.AccessToken, 3600, "/", "", secureCookie, true)
	c.SetCookie("refresh_token", resp.RefreshToken, 86400*7, "/", "", secureCookie, true)

	// 使用JavaScript关闭窗口并发送消息到父窗口
	html := `
	<html>
	<body>
		<script>
		if (window.opener) {
			// 发送登录成功消息到父窗口
			window.opener.postMessage({
				type: 'oauth-callback',
				provider: 'google',
				success: true,
				data: {
					userId: "` + resp.UserID + `",
					email: "` + resp.Email + `",
					username: "` + resp.Username + `",
					role: "` + resp.Role + `",
					accessToken: "` + resp.AccessToken + `",
					refreshToken: "` + resp.RefreshToken + `"
				}
			}, "*");
			window.close();
		} else {
			document.write('Authentication successful. You can now close this window and return to the app.');
		}
		</script>
	</body>
	</html>
	`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}
