package auth

import (
	"FullStackOfYear/backend/services/auth"
	"FullStackOfYear/backend/services/oauth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OAuthHandler struct {
	stateManager *oauth.StateManager
	provider     oauth.Provider
	authService  *auth.AuthService
}

var OAuthInstance *OAuthHandler

func NewOAuthHandler(sm *oauth.StateManager, provider oauth.Provider, as *auth.AuthService) *OAuthHandler {
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

	if !h.stateManager.ValidateState(c, state) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	token, err := h.provider.ExchangeCode(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code"})
		return
	}

	userInfo, err := h.provider.GetUserInfo(c, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	response, err := h.authService.HandleOAuthLogin(c, userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// InitOAuthHandler initializes the OAuth handler instance
func InitOAuthHandler(sm *oauth.StateManager, provider oauth.Provider, as *auth.AuthService) {
	OAuthInstance = NewOAuthHandler(sm, provider, as)
}
