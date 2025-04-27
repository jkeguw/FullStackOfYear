package auth

import (
	"project/backend/types/auth"
	"context"
)

// MockOAuthProvider is a simple implementation of the OAuthProvider interface
type MockOAuthProvider struct{}

// ExchangeCode exchanges an authorization code for an OAuth token
func (p *MockOAuthProvider) ExchangeCode(ctx context.Context, code string) (*auth.OAuthToken, error) {
	// This is a mock implementation
	return &auth.OAuthToken{
		AccessToken:  "mock_access_token",
		TokenType:    "Bearer",
		RefreshToken: "mock_refresh_token",
		ExpiresIn:    3600,
	}, nil
}

// GetUserInfo fetches user information using an OAuth token
func (p *MockOAuthProvider) GetUserInfo(ctx context.Context, token *auth.OAuthToken) (*auth.OAuthUserInfo, error) {
	// This is a mock implementation
	return &auth.OAuthUserInfo{
		ID:            "mock_user_id",
		Email:         "mock_user@example.com",
		VerifiedEmail: true,
		Name:          "Mock User",
		Picture:       "https://example.com/avatar.jpg",
	}, nil
}

// NewMockOAuthProvider creates a new mock OAuth provider
func NewMockOAuthProvider() *MockOAuthProvider {
	return &MockOAuthProvider{}
}