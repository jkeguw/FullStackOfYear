package middleware

import (
	"project/backend/config"
	"project/backend/internal/database"
	"project/backend/internal/errors"
	"project/backend/services/jwt"
	"project/backend/types/auth"
	"context"
	"encoding/json"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// AuthMiddlewareSuite 是Auth中间件的测试套件
type AuthMiddlewareSuite struct {
	suite.Suite
	router     *gin.Engine
	miniRedis  *miniredis.Miniredis
	testServer *httptest.Server
}

// SetupSuite 在所有测试前设置测试环境
func (s *AuthMiddlewareSuite) SetupSuite() {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 初始化配置
	if config.Cfg == nil {
		config.Cfg = &config.Config{
			JWT: config.JWTConfig{
				Secret:        "test_secret_key_for_middleware_testing",
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

	// 设置 Redis 客户端
	database.RedisClient = redis.NewClient(&redis.Options{
		Addr: s.miniRedis.Addr(),
	})

	// 创建路由器
	s.router = gin.New()
	s.router.Use(gin.Recovery())

	// 创建测试路由
	s.setupTestRoutes()
}

// TearDownSuite 在所有测试后清理资源
func (s *AuthMiddlewareSuite) TearDownSuite() {
	if s.miniRedis != nil {
		s.miniRedis.Close()
	}
	if database.RedisClient != nil {
		database.RedisClient.Close()
	}
}

// setupTestRoutes 设置测试路由
func (s *AuthMiddlewareSuite) setupTestRoutes() {
	// 公开路由
	s.router.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "public route"})
	})

	// 需要认证的路由
	authRoutes := s.router.Group("/auth")
	authRoutes.Use(Auth())
	{
		authRoutes.GET("/user", func(c *gin.Context) {
			userId := c.GetString("userId")
			role := c.GetString("userRole")
			c.JSON(http.StatusOK, gin.H{
				"status":   "success",
				"userId":   userId,
				"userRole": role,
			})
		})
	}

	// 需要特定角色的路由
	adminRoutes := s.router.Group("/admin")
	adminRoutes.Use(Auth(), RequireRoles("admin"))
	{
		adminRoutes.GET("/dashboard", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "success", "message": "admin dashboard"})
		})
	}
}

// 生成有效令牌并存储在Redis中
func (s *AuthMiddlewareSuite) generateValidToken(userID, role, deviceID string) string {
	claims := auth.Claims{
		UserID:   userID,
		Role:     role,
		DeviceID: deviceID,
		Type:     "access",
	}

	token, expiresAt, err := jwt.GenerateToken(claims, time.Hour)
	if err != nil {
		s.T().Fatalf("Failed to generate token: %v", err)
	}

	// 存储令牌到Redis
	accessKey := fmt.Sprintf("token:access:%s:%s", userID, deviceID)
	userTokensKey := fmt.Sprintf("user:tokens:%s", userID)

	database.RedisClient.Set(context.Background(), accessKey, token, time.Until(expiresAt))
	database.RedisClient.SAdd(context.Background(), userTokensKey, deviceID)

	return token
}

// TestPublicRoute 测试公开路由
func (s *AuthMiddlewareSuite) TestPublicRoute() {
	req, _ := http.NewRequest("GET", "/public", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "success", response["status"])
	assert.Equal(s.T(), "public route", response["message"])
}

// TestAuthenticatedRouteWithValidToken 测试有效令牌的认证路由
func (s *AuthMiddlewareSuite) TestAuthenticatedRouteWithValidToken() {
	// 生成有效令牌
	userID := "test_user_id"
	role := "user"
	deviceID := "test_device_id"
	token := s.generateValidToken(userID, role, deviceID)

	// 创建请求
	req, _ := http.NewRequest("GET", "/auth/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(s.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "success", response["status"])
	assert.Equal(s.T(), userID, response["userId"])
	assert.Equal(s.T(), role, response["userRole"])
}

// TestAuthenticatedRouteWithoutToken 测试没有令牌的认证路由
func (s *AuthMiddlewareSuite) TestAuthenticatedRouteWithoutToken() {
	req, _ := http.NewRequest("GET", "/auth/user", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), float64(errors.Unauthorized), response["code"])
	assert.Contains(s.T(), response["message"], "header required")
}

// TestAuthenticatedRouteWithInvalidTokenFormat 测试无效格式令牌的认证路由
func (s *AuthMiddlewareSuite) TestAuthenticatedRouteWithInvalidTokenFormat() {
	req, _ := http.NewRequest("GET", "/auth/user", nil)
	req.Header.Set("Authorization", "InvalidFormat token123")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), float64(errors.Unauthorized), response["code"])
	assert.Contains(s.T(), response["message"], "Invalid authorization format")
}

// TestAuthenticatedRouteWithInvalidToken 测试无效令牌的认证路由
func (s *AuthMiddlewareSuite) TestAuthenticatedRouteWithInvalidToken() {
	req, _ := http.NewRequest("GET", "/auth/user", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.format")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
}

// TestAuthenticatedRouteWithRevokedToken 测试已撤销令牌的认证路由
func (s *AuthMiddlewareSuite) TestAuthenticatedRouteWithRevokedToken() {
	// 生成有效令牌但不存储到Redis
	claims := auth.Claims{
		UserID:   "revoked_user",
		Role:     "user",
		DeviceID: "revoked_device",
		Type:     "access",
	}

	token, _, err := jwt.GenerateToken(claims, time.Hour)
	assert.NoError(s.T(), err)

	req, _ := http.NewRequest("GET", "/auth/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), float64(errors.Unauthorized), response["code"])
	assert.Contains(s.T(), response["message"], "revoked")
}

// TestRoleBasedAuthWithCorrectRole 测试具有正确角色的角色认证
func (s *AuthMiddlewareSuite) TestRoleBasedAuthWithCorrectRole() {
	// 生成管理员角色的令牌
	userID := "admin_user_id"
	role := "admin"
	deviceID := "admin_device_id"
	token := s.generateValidToken(userID, role, deviceID)

	req, _ := http.NewRequest("GET", "/admin/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "success", response["status"])
	assert.Equal(s.T(), "admin dashboard", response["message"])
}

// TestRoleBasedAuthWithIncorrectRole 测试具有错误角色的角色认证
func (s *AuthMiddlewareSuite) TestRoleBasedAuthWithIncorrectRole() {
	// 生成普通用户角色的令牌
	userID := "regular_user_id"
	role := "user" // 非管理员角色
	deviceID := "regular_device_id"
	token := s.generateValidToken(userID, role, deviceID)

	req, _ := http.NewRequest("GET", "/admin/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusForbidden, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), float64(errors.Forbidden), response["code"])
	assert.Contains(s.T(), response["message"], "权限不足")
}

// TestRequireRolesWithMultipleRoles 测试允许多个角色访问
func (s *AuthMiddlewareSuite) TestRequireRolesWithMultipleRoles() {
	// 创建一个临时路由用于测试多角色
	router := gin.New()
	router.GET("/multi-role", Auth(), RequireRoles("admin", "editor", "moderator"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	// 测试每个角色
	roles := []string{"admin", "editor", "moderator"}
	for _, role := range roles {
		token := s.generateValidToken("user_"+role, role, "device_"+role)

		req, _ := http.NewRequest("GET", "/multi-role", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusOK, w.Code, "Role '%s' should have access", role)
	}

	// 测试无权限角色
	token := s.generateValidToken("user_viewer", "viewer", "device_viewer")
	req, _ := http.NewRequest("GET", "/multi-role", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusForbidden, w.Code, "Role 'viewer' should not have access")
}

// TestAuthMiddlewareSuite 启动测试套件
func TestAuthMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareSuite))
}