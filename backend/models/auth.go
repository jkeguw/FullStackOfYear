package models

import (
	"time"
)

type LoginType string

const (
	EmailLogin     LoginType = "email"
	GoogleLogin    LoginType = "google"
	AppleLogin     LoginType = "apple"
	TwoFactorLogin LoginType = "2fa"
)

// UserRole 用户角色
type UserRole struct {
	Type string `bson:"type" json:"type"`
}

// Status 状态信息
type Status struct {
	EmailVerified bool `bson:"emailVerified" json:"emailVerified"`
}

// UserStats 用户统计信息
type UserStats struct {
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	LastLoginAt time.Time `bson:"lastLoginAt" json:"lastLoginAt"`
	LastLoginIP string    `bson:"lastLoginIP,omitempty" json:"lastLoginIP,omitempty"`
}

// Role 类型不再在此定义，改用 constants.go 中的常量

// OAuthInfo 存储 OAuth 相关信息
type OAuthInfo struct {
	Google *GoogleOAuth `bson:"google,omitempty" json:"google,omitempty"`
}

// GoogleOAuth 存储 Google OAuth 特定信息
type GoogleOAuth struct {
	ID          string    `bson:"id" json:"id"`
	Email       string    `bson:"email" json:"email"`
	Connected   bool      `bson:"connected" json:"connected"`
	ConnectedAt time.Time `bson:"connectedAt" json:"connectedAt"`
}

// LoginRecord 代表单次登录记录 (simplified without device tracking)
type LoginRecord struct {
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	IP        string    `bson:"ip" json:"ip"`
	UserAgent string    `bson:"userAgent" json:"userAgent"`
	Success   bool      `bson:"success" json:"success"`
}

// SecurityLog 记录安全相关的操作日志
type SecurityLog struct {
	Action      string    `bson:"action" json:"action"`
	Timestamp   time.Time `bson:"timestamp" json:"timestamp"`
	IP          string    `bson:"ip" json:"ip"`
	Description string    `bson:"description" json:"description"`
	DeviceInfo  string    `bson:"deviceInfo,omitempty" json:"deviceInfo,omitempty"`
}

// TwoFactorAuth 存储与二因素认证相关的信息
type TwoFactorAuth struct {
	Enabled     bool      `bson:"enabled" json:"enabled"`
	Secret      string    `bson:"secret" json:"-"` // TOTP 密钥，不通过 JSON 返回
	VerifiedAt  time.Time `bson:"verifiedAt,omitempty" json:"verifiedAt,omitempty"`
	BackupCodes []string  `bson:"backupCodes,omitempty" json:"-"` // 备份恢复码，不通过 JSON 返回
}

// Note: Device tracking has been removed
