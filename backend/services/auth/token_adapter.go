package auth

import (
	"context"
	"project/backend/internal/database"
	"project/backend/services/jwt"
	"project/backend/services/token"
	"project/backend/types/claims"
	"time"
)

// SimpleTokenGenerator is a minimal implementation of the TokenGenerator interface
type SimpleTokenGenerator struct {
	jwtService jwt.Service
}

// NewSimpleTokenGenerator creates a new SimpleTokenGenerator
func NewSimpleTokenGenerator(jwtService jwt.Service) *SimpleTokenGenerator {
	return &SimpleTokenGenerator{
		jwtService: jwtService,
	}
}

// GenerateTokenPair generates a pair of access and refresh tokens and stores them in Redis
func (g *SimpleTokenGenerator) GenerateTokenPair(userID, role, deviceID string) (string, string, error) {
	// Generate access token (1 hour expiration)
	accessClaims := jwt.Claims{
		UserID:   userID,
		Role:     role,
		DeviceID: deviceID,
		Type:     "access",
	}

	accessToken, accessExpiresAt, err := g.jwtService.GenerateToken(accessClaims, 3600*time.Second)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token (7 days expiration)
	refreshClaims := jwt.Claims{
		UserID:   userID,
		Role:     role,
		DeviceID: deviceID,
		Type:     "refresh",
	}

	refreshToken, _, err := g.jwtService.GenerateToken(refreshClaims, 7*24*3600*time.Second)
	if err != nil {
		return "", "", err
	}

	// Store tokens in Redis for revocation support
	tokenManager := token.NewManager(database.RedisClient)
	info := token.TokenInfo{
		UserID:    userID,
		DeviceID:  deviceID,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiresAt: accessExpiresAt,
	}
	if err := tokenManager.StoreTokens(context.Background(), accessToken, refreshToken, info); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateRefreshToken validates a refresh token and returns the claims
func (g *SimpleTokenGenerator) ValidateRefreshToken(token string) (*claims.Claims, error) {
	jwtClaims, err := g.jwtService.ParseToken(token)
	if err != nil {
		return nil, err
	}
	
	return &claims.Claims{
		UserID:   jwtClaims.UserID,
		Role:     jwtClaims.Role,
		DeviceID: jwtClaims.DeviceID,
		Type:     jwtClaims.Type,
	}, nil
}

// RevokeTokens revokes all tokens for a user and device
func (g *SimpleTokenGenerator) RevokeTokens(userID, deviceID string) error {
	tokenManager := token.NewManager(database.RedisClient)
	return tokenManager.RevokeTokens(userID, deviceID)
}