package auth

import "time"

// LoginResponse represents a unified login response
type LoginResponse struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ExpiresIn    int64  `json:"expiresIn,omitempty"`
	TokenType    string `json:"tokenType,omitempty"`

	// User info
	UserID    string    `json:"userId"`
	Email     string    `json:"email,omitempty"`
	Username  string    `json:"username,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`

	// OAuth specific fields
	OAuthConnected bool   `json:"oauthConnected,omitempty"`
	OAuthProvider  string `json:"oauthProvider,omitempty"`
	
	// Two-factor auth specific fields
	RequireTwoFactor bool   `json:"requireTwoFactor,omitempty"`
	TwoFactorToken   string `json:"twoFactorToken,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// TwoFactorSetupResponse 两因素认证设置响应
type TwoFactorSetupResponse struct {
	Secret        string   `json:"secret"`
	QRCode        []byte   `json:"qrCode"`
	RecoveryCodes []string `json:"recoveryCodes"`
}

// TwoFactorStatusResponse 两因素认证状态响应
type TwoFactorStatusResponse struct {
	Enabled bool `json:"enabled"`
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	OS           string    `json:"os"`
	Browser      string    `json:"browser"`
	IsCurrent    bool      `json:"isCurrent"`
	LastUsedAt   time.Time `json:"lastUsedAt"`
	LastUsedIP   string    `json:"lastUsedIp"`
	LastLocation string    `json:"lastLocation,omitempty"`
	Trusted      bool      `json:"trusted"`
}

// DeviceListResponse 设备列表响应
type DeviceListResponse struct {
	Devices []DeviceInfo `json:"devices"`
}

// SecurityLogResponse 安全日志响应
type SecurityLogResponse struct {
	Logs []SecurityLogEntry `json:"logs"`
}

// SecurityLogEntry 安全日志条目
type SecurityLogEntry struct {
	Action      string    `json:"action"`
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
	IP          string    `json:"ip,omitempty"`
	Location    string    `json:"location,omitempty"`
	DeviceInfo  string    `json:"deviceInfo,omitempty"`
}
