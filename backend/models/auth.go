package models

import "time"

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

// LoginRecord 代表单次登录记录
type LoginRecord struct {
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	IP        string    `bson:"ip" json:"ip"`
	UserAgent string    `bson:"userAgent" json:"userAgent"`
	Location  string    `bson:"location,omitempty" json:"location,omitempty"`
	Success   bool      `bson:"success" json:"success"`
	DeviceID  string    `bson:"deviceId" json:"deviceId"`
}

// SecurityLog 记录安全相关的操作日志
type SecurityLog struct {
	Action      string    `bson:"action" json:"action"`
	Timestamp   time.Time `bson:"timestamp" json:"timestamp"`
	IP          string    `bson:"ip" json:"ip"`
	Description string    `bson:"description" json:"description"`
	DeviceInfo  string    `bson:"deviceInfo,omitempty" json:"deviceInfo,omitempty"`
}
