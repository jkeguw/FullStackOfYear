package auth

import (
	"FullStackOfYear/backend/internal/errors"
	"FullStackOfYear/backend/services/auth"
	"FullStackOfYear/backend/services/oauth"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	ErrStateInvalid    = "Invalid state parameter"
	ErrCodeExchange    = "Failed to exchange authorization code"
	ErrUserInfo        = "Failed to get user info"
	ErrDatabaseOp      = "Database operation failed"
	ErrTokenGeneration = "Failed to generate tokens"
)

type OAuthHandler struct {
	stateManager *oauth.StateManager
	provider     oauth.Provider
	authService  auth.Service
}

var OAuthInstance *OAuthHandler

func NewOAuthHandler(sm *oauth.StateManager, provider oauth.Provider, as auth.Service) *OAuthHandler {
	return &OAuthHandler{
		stateManager: sm,
		provider:     provider,
		authService:  as,
	}
}

// InitiateOAuth starts OAuth flow
func (h *OAuthHandler) InitiateOAuth(c *gin.Context) {
	state, err := h.stateManager.GenerateState(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	authURL := h.provider.GenerateAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// HandleCallback processes OAuth callback
func (h *OAuthHandler) HandleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	// 验证必要参数
	if code == "" {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "Authorization code is required"))
		return
	}

	if state == "" {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "State parameter is required"))
		return
	}

	// 验证state
	if !h.stateManager.ValidateState(c, state) {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, ErrStateInvalid))
		return
	}

	// 交换token
	token, err := h.provider.ExchangeCode(c, code)
	if err != nil {
		log.Printf("OAuth code exchange failed: %v", err)
		c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, ErrCodeExchange))
		return
	}

	// 获取用户信息
	userInfo, err := h.provider.GetUserInfo(c, token)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, ErrUserInfo))
		return
	}

	// 处理用户登录
	response, err := h.authService.HandleOAuthLogin(c, userInfo)
	if err != nil {
		log.Printf("OAuth login failed: %v", err)
		c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}

// InitOAuthHandler initializes the OAuth handler instance
func InitOAuthHandler(sm *oauth.StateManager, provider oauth.Provider, as auth.Service) {
	OAuthInstance = NewOAuthHandler(sm, provider, as)
}
