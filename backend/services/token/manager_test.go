package token

import (
	"project/backend/config"
	"project/backend/services/jwt"
	"project/backend/types/auth"
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TokenManagerSuite 定义令牌管理器测试套件
type TokenManagerSuite struct {
	suite.Suite
	ctx        context.Context
	miniRedis  *miniredis.Miniredis
	redisClient *redis.Client
	manager    *Manager
}

// SetupSuite 在所有测试前设置测试环境
func (s *TokenManagerSuite) SetupSuite() {
	// 初始化配置
	if config.Cfg == nil {
		config.Cfg = &config.Config{
			JWT: config.JWTConfig{
				Secret:        "test_secret_key_for_jwt_token_testing",
				Issuer:        "test_issuer",
				AccessExpire:  time.Hour,
				RefreshExpire: 24 * time.Hour,
			},
		}
	}

	// 创建 miniredis 实例
	mr, err := miniredis.Run()
	if err != nil {
		s.T().Fatalf("Failed to create miniredis: %v", err)
	}
	s.miniRedis = mr

	// 创建 Redis 客户端
	s.redisClient = redis.NewClient(&redis.Options{
		Addr: s.miniRedis.Addr(),
	})

	// 创建 TokenManager 实例
	s.manager = NewManager(s.redisClient)
	s.ctx = context.Background()
}

// TearDownSuite 在所有测试后清理资源
func (s *TokenManagerSuite) TearDownSuite() {
	s.redisClient.Close()
	s.miniRedis.Close()
}

// TestGenerateTokenPair 测试生成令牌对
func (s *TokenManagerSuite) TestGenerateTokenPair() {
	// 测试场景: 生成有效的令牌对
	userID := "user123"
	role := "admin"
	deviceID := "device456"

	// 生成令牌对
	accessToken, refreshToken, err := s.manager.GenerateTokenPair(userID, role, deviceID)

	// 验证结果
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), accessToken)
	assert.NotEmpty(s.T(), refreshToken)

	// 验证令牌存储在 Redis 中
	accessExists, err := s.manager.CheckTokenExists(s.ctx, userID, deviceID, "access")
	assert.NoError(s.T(), err)
	assert.True(s.T(), accessExists)

	refreshExists, err := s.manager.CheckTokenExists(s.ctx, userID, deviceID, "refresh")
	assert.NoError(s.T(), err)
	assert.True(s.T(), refreshExists)
}

// TestStoreTokens 测试令牌存储功能
func (s *TokenManagerSuite) TestStoreTokens() {
	// 测试数据
	userID := "user456"
	role := "user"
	deviceID := "device789"
	accessToken := "fake_access_token"
	refreshToken := "fake_refresh_token"
	
	// 创建令牌信息
	info := auth.TokenInfo{
		UserID:    userID,
		DeviceID:  deviceID,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	// 存储令牌
	err := s.manager.StoreTokens(s.ctx, accessToken, refreshToken, info)
	assert.NoError(s.T(), err)

	// 验证令牌存储成功
	accessKey := "token:access:" + userID + ":" + deviceID
	refreshKey := "token:refresh:" + userID + ":" + deviceID
	userTokensKey := "user:tokens:" + userID

	// 检查访问令牌
	storedAccessToken, err := s.redisClient.Get(s.ctx, accessKey).Result()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), accessToken, storedAccessToken)

	// 检查刷新令牌
	storedRefreshToken, err := s.redisClient.Get(s.ctx, refreshKey).Result()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), refreshToken, storedRefreshToken)

	// 检查用户设备集合
	deviceExists, err := s.redisClient.SIsMember(s.ctx, userTokensKey, deviceID).Result()
	assert.NoError(s.T(), err)
	assert.True(s.T(), deviceExists)
}

// TestCheckTokenExists 测试检查令牌存在性
func (s *TokenManagerSuite) TestCheckTokenExists() {
	// 测试数据
	userID := "user_check"
	deviceID := "device_check"
	
	// 设置测试令牌
	accessKey := "token:access:" + userID + ":" + deviceID
	refreshKey := "token:refresh:" + userID + ":" + deviceID
	
	s.redisClient.Set(s.ctx, accessKey, "fake_token", time.Hour)
	
	// 测试存在的令牌
	exists, err := s.manager.CheckTokenExists(s.ctx, userID, deviceID, "access")
	assert.NoError(s.T(), err)
	assert.True(s.T(), exists)
	
	// 测试不存在的令牌
	exists, err = s.manager.CheckTokenExists(s.ctx, userID, deviceID, "refresh")
	assert.NoError(s.T(), err)
	assert.False(s.T(), exists)
	
	// 设置刷新令牌并再次检查
	s.redisClient.Set(s.ctx, refreshKey, "fake_refresh", time.Hour)
	exists, err = s.manager.CheckTokenExists(s.ctx, userID, deviceID, "refresh")
	assert.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

// TestInvalidateTokens 测试令牌失效功能
func (s *TokenManagerSuite) TestInvalidateTokens() {
	// 场景1: 用户有多个设备
	userID := "multi_device_user"
	device1 := "device1"
	device2 := "device2"
	userTokensKey := "user:tokens:" + userID
	
	// 设置测试数据
	s.redisClient.Set(s.ctx, "token:access:"+userID+":"+device1, "token1", time.Hour)
	s.redisClient.Set(s.ctx, "token:refresh:"+userID+":"+device1, "refresh1", time.Hour)
	s.redisClient.Set(s.ctx, "token:access:"+userID+":"+device2, "token2", time.Hour)
	s.redisClient.Set(s.ctx, "token:refresh:"+userID+":"+device2, "refresh2", time.Hour)
	s.redisClient.SAdd(s.ctx, userTokensKey, device1, device2)
	
	// 使一个设备的令牌失效
	err := s.manager.InvalidateTokens(s.ctx, userID, device1)
	assert.NoError(s.T(), err)
	
	// 验证设备1的令牌已被删除
	exists, _ := s.redisClient.Exists(s.ctx, "token:access:"+userID+":"+device1).Result()
	assert.Equal(s.T(), int64(0), exists)
	exists, _ = s.redisClient.Exists(s.ctx, "token:refresh:"+userID+":"+device1).Result()
	assert.Equal(s.T(), int64(0), exists)
	
	// 验证设备2的令牌仍然存在
	exists, _ = s.redisClient.Exists(s.ctx, "token:access:"+userID+":"+device2).Result()
	assert.Equal(s.T(), int64(1), exists)
	
	// 验证用户令牌集合仍然存在
	exists, _ = s.redisClient.Exists(s.ctx, userTokensKey).Result()
	assert.Equal(s.T(), int64(1), exists)
	
	// 场景2: 用户最后一个设备
	singleUserID := "single_device_user"
	singleDevice := "single_device"
	singleUserTokensKey := "user:tokens:" + singleUserID
	
	// 设置测试数据
	s.redisClient.Set(s.ctx, "token:access:"+singleUserID+":"+singleDevice, "token", time.Hour)
	s.redisClient.Set(s.ctx, "token:refresh:"+singleUserID+":"+singleDevice, "refresh", time.Hour)
	s.redisClient.SAdd(s.ctx, singleUserTokensKey, singleDevice)
	
	// 使最后一个设备的令牌失效
	err = s.manager.InvalidateTokens(s.ctx, singleUserID, singleDevice)
	assert.NoError(s.T(), err)
	
	// 验证令牌已被删除
	exists, _ = s.redisClient.Exists(s.ctx, "token:access:"+singleUserID+":"+singleDevice).Result()
	assert.Equal(s.T(), int64(0), exists)
	
	// 验证用户令牌集合也被删除
	exists, _ = s.redisClient.Exists(s.ctx, singleUserTokensKey).Result()
	assert.Equal(s.T(), int64(0), exists)
}

// TestValidateRefreshToken 测试刷新令牌验证
func (s *TokenManagerSuite) TestValidateRefreshToken() {
	// 生成有效的刷新令牌
	claims := auth.Claims{
		UserID:   "refresh_user",
		Role:     "user",
		DeviceID: "refresh_device",
		Type:     "refresh",
	}
	
	refreshToken, _, err := jwt.GenerateToken(claims, time.Hour)
	assert.NoError(s.T(), err)
	
	// 验证有效的刷新令牌
	parsedClaims, err := s.manager.ValidateRefreshToken(refreshToken)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), claims.UserID, parsedClaims.UserID)
	assert.Equal(s.T(), claims.Role, parsedClaims.Role)
	assert.Equal(s.T(), claims.DeviceID, parsedClaims.DeviceID)
	assert.Equal(s.T(), claims.Type, parsedClaims.Type)
	
	// 验证无效的令牌类型
	accessClaims := auth.Claims{
		UserID:   "access_user",
		Role:     "user",
		DeviceID: "access_device",
		Type:     "access",
	}
	accessToken, _, err := jwt.GenerateToken(accessClaims, time.Hour)
	assert.NoError(s.T(), err)
	
	_, err = s.manager.ValidateRefreshToken(accessToken)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "Invalid token type")
}

// TestRevokeTokens 测试撤销令牌
func (s *TokenManagerSuite) TestRevokeTokens() {
	// 设置测试数据
	userID := "revoke_user"
	deviceID := "revoke_device"
	
	s.redisClient.Set(s.ctx, "token:access:"+userID+":"+deviceID, "token", time.Hour)
	s.redisClient.Set(s.ctx, "token:refresh:"+userID+":"+deviceID, "refresh", time.Hour)
	s.redisClient.SAdd(s.ctx, "user:tokens:"+userID, deviceID)
	
	// 撤销令牌
	err := s.manager.RevokeTokens(userID, deviceID)
	assert.NoError(s.T(), err)
	
	// 验证令牌已被撤销
	accessExists, err := s.manager.CheckTokenExists(s.ctx, userID, deviceID, "access")
	assert.NoError(s.T(), err)
	assert.False(s.T(), accessExists)
	
	refreshExists, err := s.manager.CheckTokenExists(s.ctx, userID, deviceID, "refresh")
	assert.NoError(s.T(), err)
	assert.False(s.T(), refreshExists)
}

// TestTokenManagerSuite 启动测试套件
func TestTokenManagerSuite(t *testing.T) {
	suite.Run(t, new(TokenManagerSuite))
}