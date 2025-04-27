package token

import (
	"project/backend/config"
	"project/backend/internal/errors"
	jwtService "project/backend/services/jwt"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	accessTokenKey  = "token:access:%s:%s"  // userID:deviceID
	refreshTokenKey = "token:refresh:%s:%s" // userID:deviceID
	userTokensKey   = "user:tokens:%s"      // userID
)

// TokenInfo 存储token的相关信息
type TokenInfo struct {
	UserID    string
	DeviceID  string
	Role      string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// Manager handles token operations
type Manager struct {
	rdb *redis.Client
}

// NewManager creates a new token manager
func NewManager(rdb *redis.Client) *Manager {
	return &Manager{rdb: rdb}
}

// GenerateTokenPair creates both access and refresh tokens
func (m *Manager) GenerateTokenPair(userID, role, deviceID string) (string, string, error) {
	// Generate access token
	accessClaims := jwtService.Claims{
		UserID:   userID,
		Role:     role,
		DeviceID: deviceID,
		Type:     "access",
	}

	jwtSvc := jwtService.NewService(config.GetConfig().JWT)
	accessToken, expiresAt, err := jwtSvc.GenerateToken(accessClaims, config.GetConfig().JWT.AccessExpire)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshClaims := jwtService.Claims{
		UserID:   userID,
		Role:     role,
		DeviceID: deviceID,
		Type:     "refresh",
	}

	refreshToken, _, err := jwtSvc.GenerateToken(refreshClaims, config.GetConfig().JWT.RefreshExpire)
	if err != nil {
		return "", "", err
	}

	// Store tokens
	info := TokenInfo{
		UserID:    userID,
		DeviceID:  deviceID,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiresAt: expiresAt,
	}

	if err := m.StoreTokens(context.Background(), accessToken, refreshToken, info); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// StoreTokens saves tokens to Redis
func (m *Manager) StoreTokens(ctx context.Context, accessToken, refreshToken string, info TokenInfo) error {
	accessKey := fmt.Sprintf(accessTokenKey, info.UserID, info.DeviceID)
	refreshKey := fmt.Sprintf(refreshTokenKey, info.UserID, info.DeviceID)
	userKey := fmt.Sprintf(userTokensKey, info.UserID)

	pipe := m.rdb.Pipeline()

	pipe.Set(ctx, accessKey, accessToken, time.Until(info.ExpiresAt))
	pipe.Set(ctx, refreshKey, refreshToken, config.GetConfig().JWT.RefreshExpire)
	pipe.SAdd(ctx, userKey, info.DeviceID)
	pipe.Expire(ctx, userKey, config.GetConfig().JWT.RefreshExpire)

	_, err := pipe.Exec(ctx)
	return err
}

// CheckTokenExists verifies if a token exists
func (m *Manager) CheckTokenExists(ctx context.Context, userID, deviceID, tokenType string) (bool, error) {
	var key string
	if tokenType == "access" {
		key = fmt.Sprintf(accessTokenKey, userID, deviceID)
	} else {
		key = fmt.Sprintf(refreshTokenKey, userID, deviceID)
	}

	exists, err := m.rdb.Exists(ctx, key).Result()
	return exists == 1, err
}

// InvalidateTokens removes tokens for a specific device
func (m *Manager) InvalidateTokens(ctx context.Context, userID, deviceID string) error {
	// Get all related keys
	accessKey := fmt.Sprintf(accessTokenKey, userID, deviceID)
	refreshKey := fmt.Sprintf(refreshTokenKey, userID, deviceID)
	userTokensKey := fmt.Sprintf(userTokensKey, userID)

	// Use pipeline to ensure atomic operation
	pipe := m.rdb.Pipeline()

	// Remove access token
	pipe.Del(ctx, accessKey)

	// Remove refresh token
	pipe.Del(ctx, refreshKey)

	// Remove device from user's token set
	pipe.SRem(ctx, userTokensKey, deviceID)

	// Check if this was the last device
	remainingDevices := pipe.SCard(ctx, userTokensKey)

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to invalidate tokens: %w", err)
	}

	// If no devices left, remove the user's token set
	count, err := remainingDevices.Result()
	if err != nil {
		return fmt.Errorf("failed to get remaining devices count: %w", err)
	}

	if count == 0 {
		if err := m.rdb.Del(ctx, userTokensKey).Err(); err != nil {
			return fmt.Errorf("failed to remove user tokens key: %w", err)
		}
	}

	return nil
}

func (m *Manager) ValidateRefreshToken(token string) (*jwtService.Claims, error) {
	jwtSvc := jwtService.NewService(config.GetConfig().JWT)
	claims, err := jwtSvc.ParseToken(token)
	if err != nil {
		return nil, err
	}
	if claims.Type != "refresh" {
		return nil, errors.NewAppError(errors.Unauthorized, "Invalid token type")
	}
	return claims, nil
}

// RevokeTokens 吊销用户某个设备的所有令牌
func (m *Manager) RevokeTokens(userID, deviceID string) error {
	return m.InvalidateTokens(context.Background(), userID, deviceID)
}