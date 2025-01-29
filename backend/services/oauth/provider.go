package oauth

import (
	"FullStackOfYear/backend/types/auth"
	"context"
)

// Provider defines the interface for OAuth providers
type Provider interface {
	// GenerateAuthURL generates the OAuth authorization URL
	GenerateAuthURL(state string) string

	// ExchangeCode exchanges OAuth code for tokens
	ExchangeCode(ctx context.Context, code string) (*auth.OAuthToken, error)

	// GetUserInfo fetches user information using OAuth token
	GetUserInfo(ctx context.Context, token *auth.OAuthToken) (*auth.OAuthUserInfo, error)

	// RefreshToken refreshes an expired access token
	RefreshToken(ctx context.Context, refreshToken string) (*auth.OAuthToken, error)
}
