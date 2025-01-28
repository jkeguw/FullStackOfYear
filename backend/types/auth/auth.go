package auth

import "time"

// TokenInfo represents the stored token information
type TokenInfo struct {
	UserID    string    `json:"userId"`
	DeviceID  string    `json:"deviceId"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// Claims represents the JWT claims structure
type Claims struct {
	UserID   string `json:"uid"`
	Role     string `json:"role"`
	DeviceID string `json:"deviceId"`
	Type     string `json:"type"` // "access" or "refresh"
}
