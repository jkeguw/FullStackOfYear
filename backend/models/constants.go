package models

import "time"

const (
	// Common constants
	MaxLoginHistory       = 50
	MaxConcurrentSessions = 5
	TwoFactorRecoveryCodeCount = 10 // 两因素认证恢复码数量

	// Token related
	VerifyTokenExpiration = 24 * time.Hour
)

// Collection names constants - Not duplicated in other files
const (
	UsersCollection        = "users"
	CartCollection         = "carts"
	ReviewsCollection      = "reviews"
)
