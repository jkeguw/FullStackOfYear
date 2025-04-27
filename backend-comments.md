# 后端代码注释汇总

## 身份认证处理器
### `/backend/handlers/auth/auth.go`
- `// Register handles user registration`
- `// Login handles user login`
- `// RefreshToken handles token refresh requests`
- `// Logout handles user logout`

### `/backend/handlers/auth/oauth.go`
- `// OAuthHandler 处理OAuth认证`
- `// HandleOAuthLogin 处理OAuth登录请求`
- `// HandleOAuthCallback 处理OAuth回调`

## 中间件
### `/backend/middleware/auth.go`
- `// Auth validates the JWT token and adds claims to context`
- `// Use a context with timeout for token operations`
- `// 使用注入的JWT服务`
- `// 使用全局JWT服务（向后兼容）`
- `// Verify token exists in Redis`
- `// Set claims to context`
- `// RequireRoles validates user roles`

### `/backend/middleware/recovery.go`
- Recovery middleware comments focus on error handling

### `/backend/middleware/limiter.go`
- `// RateLimit applies rate limiting to routes`

## 配置
### `/backend/config/logger.go`
- `// CreateLogger creates and returns a new zap logger`

### `/backend/config/init.go`
- `// InitConfig initializes the configuration from the given path`
- `// GetConfig returns the current configuration`
- `// InitLogger initializes the logger`

## 服务
### `/backend/services/auth/service.go`
- `// UpdateUser 更新用户信息`
- `// GetTwoFactorStatus 获取两因素认证状态`
- `// 简化实现，只返回认证错误`
- `// 直接返回来自OAuth Provider的错误`
- Multiple comments about two-factor authentication implementation
- `// 生成两因素认证临时Token`
- `// 解码两因素认证临时Token`
- `// 提取设备信息`
- `// 记录登录活动`
- `// 不要中断登录流程，但记录错误`

### `/backend/services/auth/security.go`
- `// SecurityService 处理安全相关功能`
- `// GetTwoFactorStatus 获取两因素认证状态`
- `// ListDevices 列出用户设备`
- `// RemoveDevice 移除设备`
- `// UpdateDevice 更新设备信息`
- `// VerifyTwoFactorCode 验证两因素认证码`
- `// VerifyAndEnableTwoFactor 验证并启用两因素认证`
- `// SetupTwoFactor 设置两因素认证`
- `// 辅助方法` (Helper methods)
- `// generateRecoveryCodes 生成备份码`
- `// generateRandomCode 生成随机字符串`

## 模型
### `/backend/models/auth.go`
- `// UserRole 用户角色`
- `// Status 状态信息`
- `// UserStats 用户统计信息`
- `// Role 类型不再在此定义，改用 constants.go 中的常量`
- `// OAuthInfo 存储 OAuth 相关信息`
- `// GoogleOAuth 存储 Google OAuth 特定信息`
- `// LoginRecord 代表单次登录记录`
- `// SecurityLog 记录安全相关的操作日志`
- `// TwoFactorAuth 存储与二因素认证相关的信息`

### `/backend/models/user.go`
- `// User 代表系统中的用户实体`
- `// NewUser 创建新用户`
- `// AddLoginRecord 添加登录记录`
- `// AddSecurityLog 添加安全日志`
- `// UpdateOAuthInfo 更新 OAuth 信息`
- `// AddDevice 添加设备到活跃设备列表`
- `// RemoveDevice 从活跃设备列表中移除设备`
- `// EnableTwoFactor 启用两因素认证`
- `// DisableTwoFactor 禁用两因素认证`
- `// UseBackupCode 使用备份恢复码`
- `// UpdateUserTwoFactorPending 更新用户两因素认证待激活状态`