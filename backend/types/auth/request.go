package auth

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
	DeviceID        string `json:"deviceId" binding:"required"`
	DeviceName      string `json:"deviceName,omitempty"`
	DeviceType      string `json:"deviceType,omitempty"`
	DeviceOS        string `json:"deviceOS,omitempty"`
	DeviceBrowser   string `json:"deviceBrowser,omitempty"`
	IP              string `json:"ip,omitempty"`
}

// LoginRequest represents a unified login request
type LoginRequest struct {
	// Common fields
	LoginType LoginType `json:"loginType" binding:"-"` // 移除required，在登录时自动设置
	DeviceID  string    `json:"deviceId" binding:"required"`

	// Email login fields
	Email    string `json:"email" binding:"required_if=LoginType email"`
	Password string `json:"password" binding:"required_if=LoginType email"`

	// OAuth login fields
	Code  string `json:"code,omitempty"`
	State string `json:"state,omitempty"`

	// Two-factor auth fields
	TwoFactorToken string `json:"twoFactorToken,omitempty"`
	TwoFactorCode  string `json:"twoFactorCode,omitempty"`

	// Device info fields
	DeviceName    string `json:"deviceName,omitempty"`
	DeviceType    string `json:"deviceType,omitempty"`
	DeviceOS      string `json:"deviceOS,omitempty"`
	DeviceBrowser string `json:"deviceBrowser,omitempty"`
	IP            string `json:"ip,omitempty"`
}

// Validate performs login type specific validation
func (r *LoginRequest) Validate() bool {
	switch r.LoginType {
	case EmailLogin:
		return r.Email != "" && r.Password != ""
	case GoogleLogin:
		return r.Code != "" && r.State != ""
	case TwoFactorLogin:
		return r.TwoFactorToken != "" && r.TwoFactorCode != ""
	default:
		return false
	}
}

// TwoFactorVerifyRequest 两因素认证验证请求
type TwoFactorVerifyRequest struct {
	Code string `json:"code" binding:"required"`
}

// PasswordChangeRequest 密码修改请求
type PasswordChangeRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

// DeviceUpdateRequest 设备更新请求
type DeviceUpdateRequest struct {
	Name    string `json:"name"`
	Trusted bool   `json:"trusted"`
}
