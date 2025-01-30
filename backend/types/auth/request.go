package auth

// LoginRequest represents a unified login request
type LoginRequest struct {
	// Common fields
	LoginType LoginType `json:"loginType" binding:"required"`
	DeviceID  string    `json:"deviceId" binding:"required"`

	// Email login fields
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`

	// OAuth login fields
	Code  string `json:"code,omitempty"`
	State string `json:"state,omitempty"`
}

// Validate performs login type specific validation
func (r *LoginRequest) Validate() bool {
	switch r.LoginType {
	case EmailLogin:
		return r.Email != "" && r.Password != ""
	case GoogleLogin:
		return r.Code != "" && r.State != ""
	default:
		return false
	}
}
