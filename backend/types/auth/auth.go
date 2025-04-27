package auth

import (
	"project/backend/types/claims"
	"time"
)

type Claims = claims.Claims

// TokenInfo represents the stored token information
type TokenInfo struct {
	UserID    string    `json:"userId"`
	DeviceID  string    `json:"deviceId"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}
