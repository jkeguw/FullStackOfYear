//go:generate mockgen -source=types.go -destination=../../tests/mocks/auth_mocks.go -package=mocks
package auth

import (
	"FullStackOfYear/backend/models"
	"FullStackOfYear/backend/types/auth"
	"FullStackOfYear/backend/types/claims"
	"context"
)

type Service interface {
	// Auth
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
	ValidateEmailPassword(ctx context.Context, email, password string) (*models.User, error)

	// OAuth
	HandleOAuthLogin(ctx context.Context, userInfo *auth.OAuthUserInfo) (*auth.LoginResponse, error)

	// Email
	SendVerificationEmail(ctx context.Context, userID string) error
	VerifyEmail(ctx context.Context, token string) error
	GenerateEmailVerificationToken(ctx context.Context, userID string) (string, error)
	GenerateEmailChangeToken(user *models.User, newEmail string) (string, error)

	// Token
	GenerateTokenPair(ctx context.Context, userID, role, deviceID string) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	RevokeTokens(ctx context.Context, userID, deviceID string) error

	// User
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

// TokenGenerator interface
type TokenGenerator interface {
	GenerateTokenPair(userID, role, deviceID string) (string, string, error)
	ValidateRefreshToken(token string) (*claims.Claims, error)
	RevokeTokens(userID, deviceID string) error
}

// EmailSender interface
type EmailSender interface {
	SendVerificationEmail(to, username, token string) error
}

// OAuthProvider interface
type OAuthProvider interface {
	ExchangeCode(ctx context.Context, code string) (*auth.OAuthToken, error)
	GetUserInfo(ctx context.Context, token *auth.OAuthToken) (*auth.OAuthUserInfo, error)
}
