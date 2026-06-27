package token

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"project/backend/config"
	"project/backend/internal/errors"
	"project/backend/services/jwt"
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
	accessClaims := jwt.Claims{
		UserID:   userID,
		Role:     role,
		DeviceID: deviceID,
		Type:     "access",
	}

	jwtSvc := jwt.NewService(config.GetConfig().JWT)
	accessToken, expiresAt, err := jwtSvc.GenerateToken(accessClaims, config.GetConfig().JWT.AccessExpire)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshClaims := jwt.Claims{
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
	// Redis不可用时不允许登录成功，因为无法验证后续请求
	if m.rdb == nil {
		return errors.NewAppError(errors.InternalError, "Redis unavailable")
	}

	accessKey := fmt.Sprintf(accessTokenKey, info.UserID, info.DeviceID)

	userKey := fmt.Sprintf(userTokensKey, info.UserID)

	pipe := m.rdb.Pipeline()

	pipe.Set(ctx, accessKey, accessToken, time.Until(info.ExpiresAt))
	pipe.Set(ctx, fmt.Sprintf(refreshTokenKey, info.UserID, info.DeviceID), refreshToken, config.GetConfig().JWT.RefreshExpire)
	pipe.SAdd(ctx, userKey, info.DeviceID)
	pipe.Expire(ctx, userKey, config.GetConfig().JWT.RefreshExpire)

	_, err := pipe.Exec(ctx)
	return err
}

// CheckTokenExists verifies if a token exists and optionally matches the provided token value
func (m *Manager) CheckTokenExists(ctx context.Context, userID, deviceID, tokenType string, tokenValue ...string) (bool, error) {
	// Redis不可用时不允许认证通过
	if m.rdb == nil {
		return false, errors.NewAppError(errors.InternalError, "Redis unavailable")
	}

	key := fmt.Sprintf(accessTokenKey, userID, deviceID)

	storedToken, err := m.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	// 如果提供了token值，校验是否匹配，防止Token复用/重放
	if len(tokenValue) > 0 && storedToken != tokenValue[0] {
		return false, nil
	}

	return true, nil
}

// InvalidateTokens removes tokens for a specific device
func (m *Manager) InvalidateTokens(ctx context.Context, userID, deviceID string) error {
	// 如果Redis不可用，直接返回成功，不影响用户登出
	if m.rdb == nil {
		return nil
	}

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

func (m *Manager) ValidateRefreshToken(token string) (*jwt.Claims, error) {
	jwtSvc := jwt.NewService(config.GetConfig().JWT)
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
